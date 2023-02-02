package data

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Peer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	PeerID  string `json:"peerID"`
}

type File struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
	Size int64  `json:"size"`
	Hash string `json:"hash"`
	Name string `json:"name"`
}

var db *sql.DB

func Init() *sql.DB {
	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		fmt.Println(err)
	}

	res, err := db.Exec(`CREATE TABLE IF NOT EXISTS peers (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, address TEXT, peer_id TEXT, is_blocked INTEGER DEFAULT 0, UNIQUE(peer_id));
	CREATE TABLE IF NOT EXISTS shared_files (id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT, size INTEGER, hash TEXT);
	`)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

	return db
}

func AddPeer(name string, address string, peerID string) error {
	_, err := db.Exec("INSERT INTO peers (name, address, peer_id) VALUES (?, ?, ?)", name, address, peerID)
	return err
}

func GetPeers() ([]Peer, error) {
	rows, err := db.Query("SELECT * FROM peers")
	if err != nil {
		return nil, err
	}

	var peers []Peer
	for rows.Next() {
		var peer Peer
		err = rows.Scan(&peer.ID, &peer.Name, &peer.Address, &peer.PeerID)
		if err != nil {
			fmt.Println(err)
		} else {
			peers = append(peers, peer)
		}
	}

	return peers, nil
}

func AddSharedFile(path string, size int64, hash string) error {
	_, err := db.Exec("INSERT INTO shared_files (path, size, hash) VALUES (?, ?, ?)", path, size, hash)
	return err
}

func GetSharedFiles() ([]File, error) {
	rows, err := db.Query("SELECT * FROM shared_files")
	if err != nil {
		return nil, err
	}

	var files []File
	for rows.Next() {
		var file File
		err = rows.Scan(&file.ID, &file.Path, &file.Size, &file.Hash)
		if err != nil {
			fmt.Println(err)
		} else {
			file.Name = filepath.Base(file.Path)
			files = append(files, file)
		}
	}

	return files, nil
}
