package vr

//
// This is a outline of primary-backup replication based on a simplifed version of Viewstamp replication.
//
//
//

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"time"
)

// ProcessCommand process command function
type ProcessCommand func(cmd interface{})

// the 3 possible server status
const (
	NORMAL = iota
	VIEWCHANGE
	RECOVERING
	STATETRANSFER
)

const debugging = false

// PBServer defines the state of a replica server (either primary or backup)
type PBServer struct {
	mu             sync.Mutex    // Lock to protect shared access to this peer's state
	peers          []*rpc.Client // RPC end points of all peers
	me             int           // this peer's index into peers[]
	currentView    int           // what this peer believes to be the current active view
	status         int           // the server's current status (NORMAL, VIEWCHANGE or RECOVERING)
	lastNormalView int           // the latest view which had a NORMAL status

	log []interface{} // the log of "commands"

	// `commit-number`
	commitIndex int // all log entries <= commitIndex are considered to have been committed.

	// key being request index
	// value being list of peer ids
	prepareOKTable map[int][]int

	processCommand ProcessCommand
}

// PrepareArgs Prepare defines the arguments for the Prepare RPC
// Note that all field names must start with a capital letter for an RPC args struct
type PrepareArgs struct {
	View          int         // the primary's current view
	PrimaryCommit int         // the primary's commitIndex
	Index         int         // the index position at which the log entry is to be replicated on backups
	Entry         interface{} // the log entry to be replicated
	From          int         // primary machine index
}

// PrepareReply defines the reply for the Prepare RPC
// Note that all field names must start with a capital letter for an RPC reply struct
type PrepareReply struct {
	View    int  // the backup's current view
	Success bool // whether the Prepare request has been accepted or rejected
}

// CommitArgs Prepare defines the arguments for the Prepare RPC
// Note that all field names must start with a capital letter for an RPC args struct
type CommitArgs struct {
	View          int // the primary's current view
	PrimaryCommit int // the primary's commitIndex
}

// CommitReply defines the reply for the Prepare RPC
// Note that all field names must start with a capital letter for an RPC reply struct
type CommitReply struct {
	View    int  // the backup's current view
	Success bool // whether the Prepare request has been accepted or rejected
}

// RecoveryArgs defined the arguments for the Recovery RPC
type RecoveryArgs struct {
	View   int // the view that the backup would like to synchronize with
	Server int // the server sending the Recovery RPC (for debugging)
}

// RecoveryReply - -
type RecoveryReply struct {
	View          int           // the view of the primary
	Entries       []interface{} // the primary's log including entries replicated up to and including the view.
	PrimaryCommit int           // the primary's commitIndex
	Success       bool          // whether the Recovery request has been accepted or rejected
	From          int
}

// ViewChangeArgs - -
type ViewChangeArgs struct {
	View int // the new view to be changed into
}

// ViewChangeReply - -
type ViewChangeReply struct {
	LastNormalView int           // the latest view which had a NORMAL status at the server
	Log            []interface{} // the log at the server
	Success        bool          // whether the ViewChange request has been accepted/rejected
	From           int
}

// StartViewArgs - -
type StartViewArgs struct {
	View int           // the new view which has completed view-change
	Log  []interface{} // the log associated with the new new
}

// StartViewReply - -
type StartViewReply struct {
}

// StateTransArgs state transfer arguments
type StateTransArgs struct {
	View  int
	Index int
}

// StateTransReply state transfer reply
type StateTransReply struct {
	View         int
	Logs         []interface{}
	CommitNumber int
	Success      bool
}

// GetPrimary is an auxilary function that returns the server index of the
// primary server given the view number (and the total number of replica servers)
func GetPrimary(view int, nservers int) int {
	return view % nservers
}

func (srv *PBServer) nF() int {
	return len(srv.peers)/2 + 1
}

// isCommitted is called by tester to check whether an index position
// has been considered committed by this server
func (srv *PBServer) isCommitted(index int) (committed bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.debug("isCommitted", srv.commitIndex, index)
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
func Make(processCommand ProcessCommand) *PBServer {
	srv := &PBServer{
		currentView:    0,
		lastNormalView: 0,
		status:         NORMAL,
	}
	// all servers' log are initialized with a dummy command at index 0
	var v interface{}
	srv.log = append(srv.log, v)

	srv.prepareOKTable = map[int][]int{}
	srv.processCommand = processCommand

	// Your other initialization code here, if there's any
	return srv
}

// Start by assigning a list of connected peers
// this should only be called once after `Make`
func (srv *PBServer) Start(peers []*rpc.Client, me int) {
	srv.peers = peers
	srv.me = me
}

// PushCommand is invoked by tester on some replica server to replicate a
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
func (srv *PBServer) PushCommand(command interface{}) (int, int, bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	// do not process command if status is not NORMAL
	// and if i am not the primary in the current view
	if srv.status != NORMAL {
		return -1, srv.currentView, false
	} else if GetPrimary(srv.currentView, len(srv.peers)) != srv.me {
		return -1, srv.currentView, false
	}

	// adds the request to the end of the log
	srv.debug("appending command log", command, srv.log, len(srv.log))
	commandInsertIndex := len(srv.log)

	srv.appendLogEntry(command)

	// Then it sends a ⟨PREPARE v, m, n, k⟩ message to the other replicas
	// where v is the current view-number
	// m is the message it received from the client
	// n is the op-number it assigned to the request
	// and k is the commit-number
	srv.primarySendPrepare(srv.currentView, command, srv.commitIndex, commandInsertIndex)

	return commandInsertIndex, srv.currentView, true
}

func (srv *PBServer) primarySendPrepare(
	viewNumber int,
	msg interface{},
	commitNumber int,
	commandInsertIndex int) {
	prepareArg := &PrepareArgs{
		View:          viewNumber,
		PrimaryCommit: commitNumber,
		Index:         commandInsertIndex,
		Entry:         msg,
		From:          srv.me,
	}
	for idx := range srv.peers {
		go srv.sendPrepare(idx, prepareArg, &PrepareReply{})
	}
}

// PrepareOKArgs = =
type PrepareOKArgs struct {
	View int

	// commandInsertIndex
	Index int

	// machine index
	From int
}

// PrepareOKReply = =
type PrepareOKReply struct {
	Success bool
}

func (srv *PBServer) shouldCommit(index int) bool {
	for idx := 0; idx <= index; idx++ {
		if len(srv.prepareOKTable[index]) < srv.nF() {
			return false
		}
	}
	return true
}

func (srv *PBServer) sendCommits(index int) {
	for idx := 0; idx <= index; idx++ {
		for _, peerIdx := range srv.prepareOKTable[idx] {
			commitArgs := CommitArgs{
				View:          srv.currentView,
				PrimaryCommit: idx,
			}
			go srv.primarySendCommit(peerIdx, &commitArgs, &CommitReply{})
		}
	}
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
	return srv.callPeer(server, "PBServer.Prepare", args, reply)
}

func (srv *PBServer) sendStateTransferRequest(server int, args *StateTransArgs, reply *StateTransReply) bool {
	return srv.callPeer(server, "PBServer.StateTransfer", args, reply)
}

func (srv *PBServer) needStateTransfer(targetIndex int) bool {
	return len(srv.log) < targetIndex
}

// doStateTransferIfNot do state transfer if it
func (srv *PBServer) doStateTransferIfNot(targetIndex int, completion func(success bool)) {
	srv.mu.Lock()
	defer func() {
		srv.status = NORMAL
		srv.mu.Unlock()
	}()

	srv.debug("preparing state transfer, target=%d", targetIndex)

	if !srv.needStateTransfer(targetIndex) {
		srv.debug("state transfer skipped")
		go completion(true)
	}

	if srv.status != NORMAL {
		srv.debug("not doing state transfer because status =", srv.status)
		return
	}
	srv.status = STATETRANSFER
	srv.debug("start state transferring")
	for idx := range srv.peers {
		var reply StateTransReply
		successed := srv.sendStateTransferRequest(idx, &StateTransArgs{
			View:  srv.currentView,
			Index: len(srv.log),
		}, &reply)
		srv.debug("peer", idx, "responsed", successed && reply.Success, "to state transfer request")
		if !successed || !reply.Success {
			continue
		}
		if reply.View == srv.currentView {
			srv.debug("appending logs", srv.log, "with", reply.Logs)

			srv.appendLogEntry(reply.Logs...)
			srv.commitIndex = reply.CommitNumber

			if len(srv.log) >= targetIndex {
				go completion(true)
				return
			}
		}
	}
	srv.debug("failed to do state transfer from", len(srv.log), "to", targetIndex)
	go completion(false)
}

func (srv *PBServer) primarySendCommit(server int, args *CommitArgs, reply *CommitReply) bool {
	return srv.callPeer(server, "PBServer.Commit", args, reply)
}

func (srv *PBServer) backupAppendPrepareCommand(args PrepareArgs) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.currentView != args.View {
		return
	}

	srv.debug("processed prepare request", args.Index)

	if args.Index < len(srv.log) {
		srv.debug("no need to append log item because already contained in log")
	} else if args.Index == len(srv.log) {
		srv.debug("appending log item", args.Entry)
		srv.appendLogEntry(args.Entry)
	} else {
		log.Fatalln("unable to append entry - missing log item before", len(srv.log), args.Index)
	}

	var okReply PrepareOKReply
	go srv.callPeer(args.From, "PBServer.PrepareOK", &PrepareOKArgs{
		View:  srv.currentView,
		Index: args.Index,
		From:  srv.me,
	}, &okReply)
}

func (srv *PBServer) debugf(format string, params ...interface{}) {
	srv.debug(fmt.Sprintf(format, params...))
}

func (srv *PBServer) debug(params ...interface{}) {
	if !debugging {
		return
	}
	print := func(format string, params ...interface{}) {
		if debugging {
			fmt.Printf(format+"\n", params...)

			// []func(format string, a ...interface{}){
			// 	color.Black,
			// 	color.Red,
			// 	color.Blue,
			// }[srv.me](format, params...)
		}
	}

	print("%s: [%d@%d] %s log=%d cmt=%d %s",
		time.Now().Format("05.000"),
		srv.me, srv.currentView,
		[]string{
			"NORMAL",
			"VIEWCHANGE",
			"RECOVERING",
			"STATETRANSFER",
		}[srv.status], len(srv.log), srv.commitIndex, fmt.Sprintf("%s", params))
}

func (srv *PBServer) appendLogEntry(cmds ...interface{}) {
	for _, cmd := range cmds {
		srv.processCommand(cmd)
		srv.log = append(srv.log, cmd)
	}
}

func (srv *PBServer) callPeer(server int, command string, args interface{}, reply interface{}) bool {
	srv.debugf("calling peer %d <%s> %s", server, command, args)
	err := srv.peers[server].Call(command, args, reply)
	srv.debugf("got reply from peer %d <%s> %s err=%s", server, command, args, reply, err.Error())
	return err == nil
}
