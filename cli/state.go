package cli

import (
	"fmt"
	"net"
	"os"
	"simpleblockchain/dao"
	"simpleblockchain/node"
	"time"
)

var state *dao.State
var thisPeerNode node.PeerNode
var conn net.Conn

// Establish the current state by starting at genesis
// and applying any existing transactions
func openState() {
	var err error
	// First of all are we running as a network service or not?
	thisPeerNode, _ = node.LoadThisPeerNoce(dataDir)
	if thisPeerNode.IP != "" {
		// Try and get a connection to the server
		timeout := 1 * time.Second
		conn, err = net.DialTimeout("tcp", thisPeerNode.TcpAddress(), timeout)
		if err != nil {
			thisPeerNode = node.PeerNode{}
			conn = nil
		} else {
			fmt.Printf("Routing command lines to the active node @ %s\n", thisPeerNode.TcpAddress())
			state = &dao.State{}
		}
	}

	if conn == nil {
		state, err = dao.LoadStateFromDisk(dataDir)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func closeState() {
	var err error
	if conn == nil {
		err = state.Close()
	} else {
		err = conn.Close()
	}
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
