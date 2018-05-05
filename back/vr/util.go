package vr

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

var debugging = true

func (srv *PBServer) debugf(format string, params ...interface{}) {
	srv.debug(fmt.Sprintf(format, params...))
}

func (srv *PBServer) debug(params ...interface{}) {
	if !debugging {
		return
	}
	printFunc := func(params ...interface{}) {
		if debugging {
			// fmt.Printf(format+"\n", params...)

			color.New([]color.Attribute{
				color.FgRed,
				color.FgGreen,
				color.FgBlue,
				color.FgYellow,
				color.FgMagenta,
				color.FgCyan,
			}[srv.me]).Print(params...)
		}
	}

	prefix := fmt.Sprintf("%s: [%d@%d] %s log=%d cmt=%d ",
		time.Now().Format("05.000"),
		srv.me, srv.currentView,
		[]string{
			"NORMAL",
			"VIEWCHANGE",
			"RECOVERING",
			"STATETRANSFER",
		}[srv.status], len(srv.log), srv.commitIndex)

	params = append([]interface{}{prefix}, params...)
	params = append(params, "\n")
	printFunc(params...)
}

func (srv *PBServer) doProcessCommand(cmds ...interface{}) {
	for _, cmd := range cmds {
		srv.processCommand(cmd)
		srv.log = append(srv.log, cmd)
	}
}

func (srv *PBServer) doReplaceCommands(cmd []interface{}) {
	srv.log = cmd
	srv.replaceCommands(cmd)
}

func (srv *PBServer) callPeer(server int, command string, args interface{}, reply interface{}) bool {
	srv.debugf("calling peer %d <%s> %s", server, command, args)
	err := srv.peers[server].Call(command, args, reply)
	srv.debugf("got reply from peer %d <%s> %s err=%s", server, command, args, reply, err)
	return err
}
