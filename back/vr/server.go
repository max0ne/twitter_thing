package vr

//
// This is a outline of primary-backup replication based on a simplifed version of Viewstamp replication.
//
//
//

import (
	"fmt"
	"net/rpc"
	"sync"
	//"fmt"
)

// the 3 possible server status
const (
	NORMAL = iota
	VIEWCHANGE
	RECOVERING
)

// PBServer defines the state of a replica server (either primary or backup)
type PBServer struct {
	mu             sync.Mutex // Lock to protect shared access to this peer's state
	peers          []*peer    // RPC end points of all peers
	me             int        // this peer's index into peers[]
	currentView    int        // what this peer believes to be the current active view
	status         int        // the server's current status (NORMAL, VIEWCHANGE or RECOVERING)
	lastNormalView int        // the latest view which had a NORMAL status

	log         []interface{} // the log of "commands"
	commitIndex int           // all log entries <= commitIndex are considered to have been committed.

	// process 1 command
	processCommand func(cmd interface{})
	// replace entire log
	replaceCommands func(cmds []interface{})
}

// Prepare defines the arguments for the Prepare RPC
// Note that all field names must start with a capital letter for an RPC args struct
type PrepareArgs struct {
	View          int         // the primary's current view
	PrimaryCommit int         // the primary's commitIndex
	Index         int         // the index position at which the log entry is to be replicated on backups
	Entry         interface{} // the log entry to be replicated
}

// PrepareReply defines the reply for the Prepare RPC
// Note that all field names must start with a capital letter for an RPC reply struct
type PrepareReply struct {
	View    int  // the backup's current view
	Success bool // whether the Prepare request has been accepted or rejected
}

// RecoverArgs defined the arguments for the Recovery RPC
type RecoveryArgs struct {
	View   int // the view that the backup would like to synchronize with
	Server int // the server sending the Recovery RPC (for debugging)
}

type RecoveryReply struct {
	View          int           // the view of the primary
	Entries       []interface{} // the primary's log including entries replicated up to and including the view.
	PrimaryCommit int           // the primary's commitIndex
	Success       bool          // whether the Recovery request has been accepted or rejected
}

type ViewChangeArgs struct {
	View int // the new view to be changed into
}

type ViewChangeReply struct {
	LastNormalView int           // the latest view which had a NORMAL status at the server
	Log            []interface{} // the log at the server
	Success        bool          // whether the ViewChange request has been accepted/rejected
}

type StartViewArgs struct {
	View int           // the new view which has completed view-change
	Log  []interface{} // the log associated with the new new
}

type StartViewReply struct {
}

type CommitArg struct {
	PrimaryCommit int
}

type CommitReply struct {
}

// GetPrimary is an auxilary function that returns the server index of the
// primary server given the view number (and the total number of replica servers)
func GetPrimary(view int, nservers int) int {
	return view % nservers
}

// isCommitted is called by tester to check whether an index position
// has been considered committed by this server
func (srv *PBServer) isCommitted(index int) (committed bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.commitIndex >= index {
		return true
	}
	return false
}

// viewStatus is called by tester to find out the current view of this server
// and whether this view has a status of NORMAL.
func (srv *PBServer) viewStatus() (currentView int, statusIsNormal bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.currentView, srv.status == NORMAL
}

// getEntryAtIndex is called by tester to return the command replicated at
// a specific log index. If the server's log is shorter than "index", then
// ok = false, otherwise, ok = true
func (srv *PBServer) getEntryAtIndex(index int) (ok bool, command interface{}) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if len(srv.log) > index {
		return true, srv.log[index]
	}
	return false, command
}

// kill is called by tester to clean up (e.g. stop the current server)
// before moving on to the next test
func (srv *PBServer) kill() {
	// Your code here, if necessary
}

// Make is called by tester to create and initalize a PBServer
// peers is the list of RPC endpoints to every server (including self)
// me is this server's index into peers.
// startingView is the initial view (set to be zero) that all servers start in
func Make(me int,
	processCommand func(cmd interface{}),
	replaceCommands func(cmds []interface{})) *PBServer {
	srv := &PBServer{
		me:              me,
		processCommand:  processCommand,
		replaceCommands: replaceCommands,
		status:          NORMAL,
	}
	// all servers' log are initialized with a dummy command at index 0
	var v interface{}
	srv.log = append(srv.log, v)
	return srv
}

func Start(srv *PBServer, peers []*rpc.Client) {
	srv.peers = makePeers(peers)
}

// PushCommand() is invoked by tester on some replica server to replicate a
// command.  Only the primary should process this request by appending
// the command to its log and then return *immediately* (while the log is being replicated to backup servers).
// if this server isn't the primary, returns false.
// Note that since the function returns immediately, there is no guarantee that this command
// will ever be committed upon return, since the primary
// may subsequently fail before replicating the command to all servers
//
// The first return value is the index that the command will appear at
// *if it's eventually committed*. The second return value is the current
// view. The third return value is true if this server believes it is
// the primary.
func (srv *PBServer) PushCommand(command interface{}, reply *bool) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	// if i am not the primary, forward message to primary
	if GetPrimary(srv.currentView, len(srv.peers)) != srv.me {
		srv.peers[GetPrimary(srv.currentView, len(srv.peers))].Call("PBServer.PushCommand", &command, nil)
		return nil
	}
	// do not process command if status is not NORMAL
	if srv.status != NORMAL {
		return fmt.Errorf("vr server status not NORMAL, %d", srv.status)
	}

	srv.doProcessCommand(command)

	go func(command interface{}, primaryServer *PBServer, log_length int) {
		cnt := 0
		for i := 0; i < len(primaryServer.peers); i++ {
			var reply PrepareReply
			args := PrepareArgs{
				View:          primaryServer.currentView,
				PrimaryCommit: primaryServer.commitIndex,
				Index:         log_length - 1,
				Entry:         command,
			}

			rpc_ok := primaryServer.sendPrepare(i, &args, &reply)
			if rpc_ok && reply.Success {
				cnt = cnt + 1
				// If the primary has received Success=true responses from a majority of servers (including itself)
				// it considers the corresponding log index as "committed".
				if len(primaryServer.peers)/2+1 == cnt {
					if primaryServer.commitIndex < log_length-1 {
						primaryServer.commitIndex = log_length - 1
					}
				}
			}
		}
	}(command, srv, len(srv.log))
	return nil
}

// exmple code to send an AppendEntries RPC to a server.
// server is the index of the target server in srv.peers[].
// expects RPC arguments in args.
// The RPC library fills in *reply with RPC reply, so caller should pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
func (srv *PBServer) sendPrepare(server int, args *PrepareArgs, reply *PrepareReply) bool {
	return srv.peers[server].Call("PBServer.Prepare", args, reply)
}

// Prepare is the RPC handler for the Prepare RPC
func (srv *PBServer) Prepare(args *PrepareArgs, reply *PrepareReply) error {
	// Your code here
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.status != NORMAL || srv.currentView > args.View {
		reply.Success = false
	} else if GetPrimary(args.View, len(srv.peers)) == srv.me {
		reply.Success = true
	} else if args.Index == len(srv.log) && args.View == srv.currentView {
		srv.doProcessCommand(args.Entry)
		srv.commitIndex = args.PrimaryCommit
		reply.Success = true
		// its view is smaller or its log is missing entries
	} else if len(srv.log) < args.Index || srv.currentView < args.View {
		primary_id := GetPrimary(args.View, len(srv.peers))
		reply.Success = false
		recover_arg := RecoveryArgs{
			View:   args.View, // the view that the backup would like to synchronize with
			Server: srv.me,    // the server sending the Recovery RPC (for debugging)
		}
		var recover_reply RecoveryReply
		// rpc call to the primary server
		ok := srv.peers[primary_id].Call("PBServer.Recovery", &recover_arg, &recover_reply)
		if ok && recover_reply.Success {
			srv.doReplaceCommands(recover_reply.Entries)
			srv.currentView = recover_reply.View
			srv.commitIndex = recover_reply.PrimaryCommit
			reply.Success = true
		}
	} else if len(srv.log) >= args.Index {
		reply.Success = true
	}
	return nil
}

// Recovery is the RPC handler for the Recovery RPC
func (srv *PBServer) Recovery(args *RecoveryArgs, reply *RecoveryReply) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	reply.View = srv.currentView
	reply.Entries = srv.log
	reply.PrimaryCommit = srv.commitIndex
	reply.Success = true
	/*
		View          int           // the view of the primary
		Entries       []interface{} // the primary's log including entries replicated up to and including the view.
		PrimaryCommit int           // the primary's commitIndex
		Success       bool          // whether the Recovery request has been accepted or rejected
	*/
	return nil
}

// Some external oracle prompts the primary of the newView to
// switch to the newView.
// promptViewChange just kicks start the view change protocol to move to the newView
// It does not block waiting for the view change process to complete.
func (srv *PBServer) promptViewChange(newView int) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	newPrimary := GetPrimary(newView, len(srv.peers))

	if newPrimary != srv.me { //only primary of newView should do view change
		return
	} else if newView <= srv.currentView {
		return
	}
	vcArgs := &ViewChangeArgs{
		View: newView,
	}
	vcReplyChan := make(chan *ViewChangeReply, len(srv.peers))
	// send ViewChange to all servers including myself
	for i := 0; i < len(srv.peers); i++ {
		go func(server int) {
			var reply ViewChangeReply
			ok := srv.peers[server].Call("PBServer.ViewChange", vcArgs, &reply)
			// fmt.Printf("node-%d (nReplies %d) received reply ok=%v reply=%v\n", srv.me, nReplies, ok, r.reply)
			if ok {
				vcReplyChan <- &reply
			} else {
				vcReplyChan <- nil
			}
		}(i)
	}

	// wait to receive ViewChange replies
	// if view change succeeds, send StartView RPC
	go func() {
		var successfulReplies []*ViewChangeReply
		var nReplies int
		majority := len(srv.peers)/2 + 1
		for r := range vcReplyChan {
			nReplies++
			if r != nil && r.Success {
				successfulReplies = append(successfulReplies, r)
			}
			if nReplies == len(srv.peers) || len(successfulReplies) == majority {
				break
			}
		}
		ok, log := srv.determineNewViewLog(successfulReplies)
		if !ok {
			return
		}
		svArgs := &StartViewArgs{
			View: vcArgs.View,
			Log:  log,
		}
		// send StartView to all servers including myself
		for i := 0; i < len(srv.peers); i++ {
			var reply StartViewReply
			go func(server int) {
				// fmt.Printf("node-%d sending StartView v=%d to node-%d\n", srv.me, svArgs.View, server)
				srv.peers[server].Call("PBServer.StartView", svArgs, &reply)
			}(i)
		}
	}()
}

// determineNewViewLog is invoked to determine the log for the newView based on
// the collection of replies for successful ViewChange requests.
// if a quorum of successful replies exist, then ok is set to true.
// otherwise, ok = false.
func (srv *PBServer) determineNewViewLog(successfulReplies []*ViewChangeReply) (
	ok bool, newViewLog []interface{}) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	currentView := 0
	if len(successfulReplies) != 0 {
		ok = true
	} else {
		ok = false
	}
	for i := 0; i < len(successfulReplies); i++ {
		if successfulReplies[i].LastNormalView > currentView {
			currentView = successfulReplies[i].LastNormalView
			newViewLog = successfulReplies[i].Log
		} else if len(successfulReplies[i].Log) > len(newViewLog) &&
			successfulReplies[i].LastNormalView == currentView {
			newViewLog = successfulReplies[i].Log
		}
	}
	return ok, newViewLog
}

// ViewChange is the RPC handler to process ViewChange RPC.
func (srv *PBServer) ViewChange(args *ViewChangeArgs, reply *ViewChangeReply) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.currentView < args.View {
		srv.status = VIEWCHANGE
		reply.LastNormalView = srv.currentView
		reply.Log = srv.log
		reply.Success = true
	} else {
		reply.Success = false
	}
	return nil
}

// StartView is the RPC handler to process StartView RPC.
func (srv *PBServer) StartView(args *StartViewArgs, reply *StartViewReply) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if args.View < srv.currentView {
		return nil
	} else if args.View == srv.currentView && len(args.Log) < len(srv.log) {
		return nil
	}
	srv.doReplaceCommands(args.Log)
	srv.currentView = args.View
	srv.status = NORMAL
	return nil
}
