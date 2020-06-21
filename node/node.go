package node

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"simpleblockchain/dao"
)

const DefaultIP = "127.0.0.1"
const DefaultHTTPort = 8080
const endpointStatus = "/node/status"

const endpointSync = "/node/sync"
const endpointSyncQueryKeyFromBlock = "fromBlock"

const endpointAddPeer = "/node/peer"
const endpointAddPeerQueryKeyIP = "ip"
const endpointAddPeerQueryKeyPort = "port"

const EndpointBalancesList = "/balances/list"
const EndpointTxAdd = "/tx/add"

type PeerNode struct {
	IP          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"is_bootstrap"`

	// Whenever this node already established connection, sync with this Peer
	connected bool
}

func (pn PeerNode) TcpAddress() string {
	return fmt.Sprintf("%s:%d", pn.IP, pn.Port)
}

type Node struct {
	ip   string
	port uint64

	state *dao.State

	knownPeers map[string]PeerNode
}

func New(s *dao.State, ip string, port uint64, bootstrap PeerNode) *Node {
	knownPeers := make(map[string]PeerNode)
	knownPeers[bootstrap.TcpAddress()] = bootstrap

	return &Node{
		state:      s,
		ip:         ip,
		port:       port,
		knownPeers: knownPeers,
	}
}

func NewPeerNode(ip string, port uint64, isBootstrap bool, connected bool) PeerNode {
	return PeerNode{ip, port, isBootstrap, connected}
}

func (n *Node) Run() error {
	ctx := context.Background()
	fmt.Println(fmt.Sprintf("Listening on: %s:%d", n.ip, n.port))

	go n.sync(ctx)

	http.HandleFunc(EndpointBalancesList, func(w http.ResponseWriter, r *http.Request) {
		listBalancesHandler(w, r, n.state)
	})

	http.HandleFunc(EndpointTxAdd, func(w http.ResponseWriter, r *http.Request) {
		txAddHandler(w, r, n.state)
	})

	http.HandleFunc(endpointStatus, func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, n)
	})

	http.HandleFunc(endpointSync, func(w http.ResponseWriter, r *http.Request) {
		syncHandler(w, r, n)
	})

	http.HandleFunc(endpointAddPeer, func(w http.ResponseWriter, r *http.Request) {
		addPeerHandler(w, r, n)
	})

	err := n.writeThisPeerNode()
	if err != nil {
		return fmt.Errorf("Error writing the node information: %w", err)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", n.port), nil)
}

func (n *Node) AddPeer(peer PeerNode) {
	n.knownPeers[peer.TcpAddress()] = peer
}

func (n *Node) RemovePeer(peer PeerNode) {
	delete(n.knownPeers, peer.TcpAddress())
}

func (n *Node) IsKnownPeer(peer PeerNode) bool {
	if peer.IP == n.ip && peer.Port == n.port {
		return true
	}

	_, isKnownPeer := n.knownPeers[peer.TcpAddress()]

	return isKnownPeer
}

func (n *Node) writeThisPeerNode() error {
	// Make a note of this node
	thisPeerNode := PeerNode{
		IP:          n.ip,
		Port:        n.port,
		IsBootstrap: false,
		connected:   false,
	}
	thisPeerNodeJson, err := json.Marshal(thisPeerNode)
	if err != nil {
		return fmt.Errorf("Cannot marshall this peer node to json: %w", err)
	}

	return ioutil.WriteFile(dao.GetThisPeerJsonFilePath(n.state.DataDir()), thisPeerNodeJson, 0644)
}

func LoadThisPeerNoce(dataDir string) (PeerNode, error) {
	content, err := ioutil.ReadFile(dao.GetThisPeerJsonFilePath(dataDir))
	if err != nil {
		return PeerNode{}, err
	}

	var loadedPeerNode PeerNode
	err = json.Unmarshal(content, &loadedPeerNode)
	if err != nil {
		return PeerNode{}, err
	}

	return loadedPeerNode, nil
}
