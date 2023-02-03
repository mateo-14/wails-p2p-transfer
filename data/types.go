package data

import "github.com/mateo-14/wails-p2p-transfer/p2p"

type Peer struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	PeerID  string `json:"peerID"`
}

type File struct {
	Path string `json:"path"`
	PeerFile
}

type InitialData struct {
	HostData    p2p.HostData `json:"hostData"`
	SharedFiles []File       `json:"sharedFiles"`
}

type PeerFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	ID   int64  `json:"id"`
	Hash string `json:"hash"`
}
