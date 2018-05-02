package vr

// Recovery is the RPC handler for the Recovery RPC
func (srv *PBServer) Recovery(args RecoveryArgs, reply *RecoveryReply) error {
	srv.debugf("received recover request from %d @ %d", args.Server, args.View)

	if srv.status != NORMAL {
		srv.debugf("cannot reply recover, status not normal")
		*reply = RecoveryReply{
			Success: false,
		}
		return nil
	}

	if srv.currentView < args.View {
		srv.debugf("view mismatch - cannot reply recover")
		*reply = RecoveryReply{
			Success: false,
		}
		return nil
	}

	*reply = RecoveryReply{
		View:          srv.currentView,
		Entries:       srv.log,
		PrimaryCommit: srv.commitIndex,
		Success:       true,
		From:          srv.me,
	}
	return nil
}

func (srv *PBServer) requireRecovery() {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.status == RECOVERING {
		return
	}

	srv.status = RECOVERING

	args := RecoveryArgs{
		View:   srv.currentView,
		Server: srv.me,
	}
	recoverChan := make(chan *RecoveryReply)
	for idx := range srv.peers {
		// dont send recover to me
		if idx == srv.me {
			continue
		}
		var reply RecoveryReply
		go func(idx int) {
			success := srv.callPeer(idx, "PBServer.Recovery", &args, &reply)
			if success && reply.Success {
				recoverChan <- &reply
			} else {
				recoverChan <- nil
			}
		}(idx)
	}

	responses := []RecoveryReply{}
	for reply := range recoverChan {
		if reply != nil {
			responses = append(responses, *reply)
		}
		if srv.receivedEnoughRecoverReply(responses) {
			go srv.recoverWithResponses(responses)
			break
		}
	}
}

func (srv *PBServer) findPrimaryFromRecoverResponse(replies []RecoveryReply) *RecoveryReply {
	for _, reply := range replies {
		if GetPrimary(reply.View, len(srv.peers)) == reply.From {
			return &reply
		}
	}
	return nil
}

func (srv *PBServer) receivedEnoughRecoverReply(replies []RecoveryReply) bool {
	return len(replies) == len(srv.peers)-1 ||
		(len(replies) >= srv.nF()+1 && srv.findPrimaryFromRecoverResponse(replies) != nil)
}

func (srv *PBServer) recoverWithResponses(replies []RecoveryReply) {
	primaryResponse := srv.findPrimaryFromRecoverResponse(replies)
	srv.debugf("recovering with primary response", primaryResponse)

	srv.currentView = primaryResponse.View
	srv.lastNormalView = primaryResponse.View
	srv.log = primaryResponse.Entries
	srv.commitIndex = primaryResponse.PrimaryCommit
	srv.status = NORMAL
}
