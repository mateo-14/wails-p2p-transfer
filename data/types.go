package data

import "github.com/mateo-14/wails-p2p-transfer/p2p"

type Peer struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	PeerID  string `json:"peerID"`
}

type File struct {
	ID   int64  `json:"id"`
	Path string `json:"path"`
	Size int64  `json:"size"`
	Hash string `json:"hash"`
	Name string `json:"name"`
}

type InitialData struct {
	HostData    p2p.HostData `json:"hostData"`
	SharedFiles []File       `json:"sharedFiles"`
}

type PeerFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	ID   int64  `json:"id"`
}
