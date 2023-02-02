package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mateo-14/wails-p2p-transfer/data"
	"github.com/mateo-14/wails-p2p-transfer/p2p"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Enums
const (
	ReqGetFiles p2p.RequestID = iota
	ReqDownloadFile
)

// Errors
var ErrP2PAlreadyStarted = errors.New("P2P already started")

// Types
type InitialData struct {
	HostData    p2p.HostData `json:"hostData"`
	SharedFiles []data.File  `json:"sharedFiles"`
}

type PeerFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	ID   int    `json:"id"`
}

// App struct
type App struct {
	ctx context.Context
	p2p *p2p.P2P
	db  *sql.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.db = data.Init()
}

func (a *App) shutdown(ctx context.Context) {
	a.db.Close()
}

func (a *App) startP2P() (*p2p.HostData, error) {
	if a.p2p != nil {
		runtime.LogInfo(a.ctx, "P2P already started")
		return a.p2p.GetHostData(), nil
	}

	appn := &AppNotifiee{
		ctx: a.ctx,
	}
	p2p, err := p2p.NewP2P(a.ctx, appn, a.onMessage)

	if err != nil {
		return nil, err
	}

	a.p2p = p2p

	return a.p2p.GetHostData(), nil

}

func (a *App) ConnectToNode(addr string, id string) error {
	nodeaddr := fmt.Sprintf("%s/%s", addr, id)

	err := a.p2p.Connect(a.ctx, nodeaddr)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error connecting to node %s: %s\n", nodeaddr, err.Error())
		return err
	}

	return nil
}

func (a *App) OnFrontendLoad() (*InitialData, error) {
	hostd, err := a.startP2P()

	if err != nil {
		runtime.LogErrorf(a.ctx, "Error starting P2P: %s\n", err.Error())
		return nil, err
	}

	sharedFiles, err := data.GetSharedFiles()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error getting shared files: %s\n", err.Error())
		return nil, err
	}

	initialData := &InitialData{
		HostData:    *hostd,
		SharedFiles: sharedFiles,
	}

	return initialData, nil
}

func (a *App) onMessage(mh *p2p.MessageHandler) {
	mh.HandleRequest(ReqGetFiles, func(req *p2p.Request) {
		sfiles, err := data.GetSharedFiles()
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error getting shared files: %s", err.Error())
		}

		files := make([]PeerFile, 0, len(sfiles))
		for _, f := range sfiles {
			files = append(files, PeerFile{
				Name: f.Name,
				Size: f.Size,
				ID:   f.ID,
			})
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		enc.Encode(files)
		_, err = req.Write(&buf)

		if err != nil {
			runtime.LogErrorf(a.ctx, "Error writing response: %s", err.Error())
		}

		req.Close()
	})

	mh.HandleRequest(ReqDownloadFile, func(req *p2p.Request) {
		path, err := io.ReadAll(req.Body)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error reading request body: %s", err.Error())
			return
		}

		file, err := os.Open(string(path))
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error opening file: %s", err.Error())
			return
		}

		_, err = req.Write(file)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error streaming file: %s", err.Error())
			return
		}
	})
}

func (a *App) GetPeerSharedFiles(peerID string) ([]PeerFile, error) {
	res, err := a.p2p.SendMessage(a.ctx, peerID, ReqGetFiles, nil)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var files []PeerFile
	dec := gob.NewDecoder(res.Body)
	err = dec.Decode(&files)

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (a *App) AddFiles() ([]data.File, error) {
	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error opening file dialog: %s", err.Error())
		return nil, err
	}

	runtime.LogDebugf(a.ctx, "Selected files: %v", files)
	filesd := make([]data.File, 0, len(files))
	for _, path := range files {
		filei, err := os.Stat(path)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error getting file info: %s", err.Error())
			continue
		}

		hash := sha256.New()
		file, err := os.Open(path)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error opening file: %s", err.Error())
			continue
		}

		_, err = io.Copy(hash, file)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error hashing file: %s", err.Error())
			continue
		}

		hashs := hex.EncodeToString(hash.Sum(nil))

		var filed data.File
		filed.Hash = hashs
		filed.Name = filepath.Base(path)
		filed.Size = filei.Size()
		filed.Path = path

		data.AddSharedFile(path, filei.Size(), hashs)
		filesd = append(filesd, filed)
	}

	return filesd, nil
}

/* func (a *App) DownloadFile(peerID string, id int) {
	res, err := a.p2p.SendMessage(a.ctx, peerID, ReqDownloadFile, strings.NewReader(path))

	if err != nil {
		runtime.LogErrorf(a.ctx, "Error sending request: %s", err.Error())
		return
	}

	defer res.Body.Close()

	fname := filepath.Base(path)
	file, err := os.Create(fname)
	io.Copy(file, res.Body)
	file.Close()
} */
