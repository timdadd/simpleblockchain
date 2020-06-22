package node

import (
	"fmt"
	"net/http"
	"simpleblockchain/dao"
	"strconv"
	"time"
)

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash     dao.Hash     `json:"block_hash"`
	Balances dao.Balances `json:"balances"`
}

type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

type TxAddRes struct {
	Hash dao.Hash `json:"block_hash"`
}

type StatusRes struct {
	Hash        dao.Hash            `json:"block_hash"`
	BlockNumber uint64              `json:"block_number"`
	KnownPeers  map[string]PeerNode `json:"peers_known"`
}

type SyncRes struct {
	Blocks []dao.Block `json:"blocks"`
}

type AddPeerRes struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func listBalancesHandler(w http.ResponseWriter, r *http.Request, state *dao.State) {
	writeRes(w, BalancesRes{state.LatestBlockHash(), state.Balances})
}

func txAddHandler(w http.ResponseWriter, r *http.Request, state *dao.State) {
	req := TxAddReq{}
	err := readReq(r, &req)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	tx := dao.NewTx(dao.NewAccount(req.From), dao.NewAccount(req.To), req.Value, req.Data)

	block := dao.NewBlock(
		state.LatestBlockHash(),
		state.NextBlockNumber(),
		0,
		uint64(time.Now().Unix()),
		[]dao.Tx{tx},
	)

	hash, err := state.AddBlock(block)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, TxAddRes{hash})
}

func statusHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	res := StatusRes{
		Hash:        node.state.LatestBlockHash(),
		BlockNumber: node.state.LatestBlock().Header.BlockNumber,
		KnownPeers:  node.knownPeers,
	}

	writeRes(w, res)
}

func syncHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	reqHash := r.URL.Query().Get(endpointSyncQueryKeyFromBlock)

	hash := dao.Hash{}
	err := hash.UnmarshalText([]byte(reqHash))
	if err != nil {
		writeErrRes(w, err)
		return
	}

	blocks, err := dao.GetBlocksAfter(hash, node.state)
	if err != nil {
		writeErrRes(w, err)
		return
	}

	writeRes(w, SyncRes{Blocks: blocks})
}

func addPeerHandler(w http.ResponseWriter, r *http.Request, node *Node) {
	peerIP := r.URL.Query().Get(endpointAddPeerQueryKeyIP)
	peerPortRaw := r.URL.Query().Get(endpointAddPeerQueryKeyPort)

	peerPort, err := strconv.ParseUint(peerPortRaw, 10, 32)
	if err != nil {
		writeRes(w, AddPeerRes{false, err.Error()})
		return
	}

	peer := NewPeerNode(peerIP, peerPort, false, true)

	node.AddPeer(peer)

	fmt.Printf("Peer '%s' was added into KnownPeers\n", peer.TcpAddress())

	writeRes(w, AddPeerRes{true, ""})
}
