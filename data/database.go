package data

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init() *sql.DB {
	var err error
	db, err = sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS peers (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, address TEXT, peer_id TEXT, is_blocked INTEGER DEFAULT 0, UNIQUE(peer_id));
	CREATE TABLE IF NOT EXISTS shared_files (id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT, size INTEGER, hash TEXT, UNIQUE(hash, path));
	`)

	if err != nil {
		fmt.Println(err)
	}

	return db
}

func AddPeer(name string, address string, peerID string) (int64, error) {
	res, err := db.Exec("INSERT INTO peers (name, address, peer_id) VALUES (?, ?, ?)", name, address, peerID)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func GetPeers() ([]Peer, error) {
	rows, err := db.Query("SELECT * FROM peers")
	if err != nil {
		return nil, err
	}

	var peers []Peer
	var isBlocked bool
	for rows.Next() {
		var peer Peer
		err = rows.Scan(&peer.ID, &peer.Name, &peer.Address, &peer.PeerID, &isBlocked)
		if err != nil {
			fmt.Println(err)
		} else {
			peers = append(peers, peer)
		}
	}

	return peers, nil
}

func AddSharedFile(path string, size int64, hash string) (int64, error) {
	res, err := db.Exec("INSERT INTO shared_files (path, size, hash) VALUES (?, ?, ?)", path, size, hash)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func GetSharedFiles() ([]File, error) {
	rows, err := db.Query("SELECT * FROM shared_files")
	if err != nil {
		return nil, err
	}

	files := make([]File, 0)
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

func RemoveSharedFile(id int64) error {
	_, err := db.Exec("DELETE FROM shared_files WHERE id = ?", id)
	return err
}

func GetSharedFile(id int64) (File, error) {
	var file File
	err := db.QueryRow("SELECT * FROM shared_files WHERE id = ?", id).Scan(&file.ID, &file.Path, &file.Size, &file.Hash)
	if err != nil {
		return file, err
	}

	file.Name = filepath.Base(file.Path)
	return file, nil
}
