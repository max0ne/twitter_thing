package vr

// promptViewChange Some external oracle prompts the primary of the newView to
// switch to the newView.
// promptViewChange just kicks start the view change protocol to move to the newView
// It does not block waiting for the view change process to complete.
func (srv *PBServer) promptViewChange(newView int) {
	srv.debug("promptViewChange", newView)

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
			ok := srv.callPeer(server, "PBServer.ViewChange", vcArgs, &reply)
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
		var successReplies []*ViewChangeReply
		var nReplies int
		majority := len(srv.peers)/2 + 1
		for r := range vcReplyChan {
			nReplies++
			if r != nil {
				srv.debugf("recieved viewchange reply from=%d r=%s", r.From, r)
			}
			if r != nil && r.Success {
				successReplies = append(successReplies, r)
			}
			if nReplies == len(srv.peers) || len(successReplies) == majority {
				srv.debugf("received majority viewchange reply")
				break
			}
		}
		ok, log := srv.determineNewViewLog(successReplies)
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
				srv.callPeer(server, "PBServer.StartView", svArgs, &reply)
			}(i)
		}
	}()
}

// determineNewViewLog is invoked to determine the log for the newView based on
// the collection of replies for successful ViewChange requests.
// if a quorum of successful replies exist, then ok is set to true.
// otherwise, ok = false.
func (srv *PBServer) determineNewViewLog(successReplies []*ViewChangeReply) (
	ok bool, newViewLog []interface{}) {
	maxNormalView := -1
	ok = false
	for _, reply := range successReplies {
		if !reply.Success {
			continue
		}
		if reply.LastNormalView > maxNormalView {
			ok = true
			maxNormalView = reply.LastNormalView
			newViewLog = reply.Log
		} else if reply.LastNormalView == maxNormalView && len(reply.Log) > len(newViewLog) {
			newViewLog = reply.Log
		}
	}
	srv.debugf("determined new log %s, %s", ok, newViewLog)
	return ok, newViewLog
}

// ViewChange is the RPC handler to process ViewChange RPC.
func (srv *PBServer) ViewChange(args ViewChangeArgs, reply *ViewChangeReply) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if args.View < srv.currentView {
		return nil
	}
	srv.status = VIEWCHANGE
	srv.currentView = args.View
	*reply = ViewChangeReply{
		LastNormalView: srv.lastNormalView,
		Log:            srv.log,
		Success:        true,
		From:           srv.me,
	}
	return nil
}

// StartView is the RPC handler to process StartView RPC.
func (srv *PBServer) StartView(args StartViewArgs, reply *StartViewReply) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	if srv.currentView > args.View {
		return nil
	}

	srv.status = NORMAL
	srv.currentView = args.View
	srv.lastNormalView = args.View
	srv.log = args.Log
	return nil
}
