package vr

import "math"

// PrepareOK handle PrepareOK from backup to primary
func (srv *PBServer) PrepareOK(args PrepareOKArgs, reply *PrepareOKReply) error {
	// The primary waits for f PREPAREOK messages from different backups
	srv.mu.Lock()
	defer srv.mu.Unlock()

	peerIdx := args.From

	srv.debug("received prepare ok from", peerIdx)

	if args.View != srv.currentView {
		srv.debug("prepare reply from", peerIdx, "dropped")
		return nil
	}

	srv.prepareOKTable[args.Index] = append(srv.prepareOKTable[args.Index], peerIdx)

	// The primary waits for f PREPAREOK messages from different backups;
	// at this point it considers the operation (and all earlier ones) to be commit- ted.
	if srv.shouldCommit(args.Index) {
		srv.debug("command", args.Index, "received f prepare response, increment commit index from", srv.commitIndex, "to", args.Index)
		srv.commitIndex = int(math.Max(float64(srv.commitIndex), float64(args.Index)))
		srv.sendCommits(args.Index)
	}
	return nil
}

// StateTransfer accept a state transfer request
func (srv *PBServer) StateTransfer(args StateTransArgs, reply *StateTransReply) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.status != NORMAL ||
		srv.currentView != args.View ||
		len(srv.log) <= args.Index {
		*reply = StateTransReply{
			Success: false,
		}
	}

	*reply = StateTransReply{
		View:         srv.currentView,
		Logs:         srv.log[args.Index:],
		CommitNumber: srv.commitIndex,
		Success:      true,
	}
	return nil
}

// Commit is the RPC handler for the Commit RPC
func (srv *PBServer) Commit(args CommitArgs, reply *CommitReply) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.currentView != args.View {
		return nil
	}
	if srv.commitIndex < args.PrimaryCommit {
		srv.debug("commiting from", srv.commitIndex, args.PrimaryCommit)
		srv.commitIndex = args.PrimaryCommit
	} else {
		srv.debug("dropped commit", srv.commitIndex, args.PrimaryCommit)
	}
	*reply = CommitReply{
		View:    srv.currentView,
		Success: true,
	}
	return nil
}

// Prepare is the RPC handler for the Prepare RPC
func (srv *PBServer) Prepare(args PrepareArgs, reply *PrepareReply) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	srv.debug("received prepare", args)

	*reply = PrepareReply{
		View:    srv.currentView,
		Success: true,
	}

	if srv.currentView < args.View {
		srv.currentView = args.View
		go srv.requireRecovery()
		return nil
	}

	// caller need recover, drop it's request
	if srv.currentView > args.View {
		*reply = PrepareReply{
			View:    srv.currentView,
			Success: false,
		}
		return nil
	}

	if srv.needStateTransfer(args.Index) ||
		srv.status != NORMAL {
		srv.debug("cannot process prepare, doing state transfer", srv.commitIndex != args.PrimaryCommit)
		go srv.doStateTransferIfNot(args.Index, func(success bool) {
			srv.mu.Lock()
			defer srv.mu.Unlock()

			srv.debug("state transfer for prepare", success)
			if success {
				go srv.backupAppendPrepareCommand(args)
			}
		})
	} else {
		srv.debug("processing prepare")
		go srv.backupAppendPrepareCommand(args)
	}
	return nil
}
