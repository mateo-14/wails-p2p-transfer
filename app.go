package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mateo-14/wails-p2p-transfer/p2p"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	ReqGetFiles p2p.RequestID = iota
	ReqDownloadFile
)

// Errors
var ErrP2PAlreadyStarted = errors.New("P2P already started")

// App struct
type App struct {
	ctx context.Context
	p2p *p2p.P2P
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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

func (a *App) OnFrontendLoad() (*p2p.HostData, error) {
	return a.startP2P()
}

func (a *App) onMessage(mh *p2p.MessageHandler) {
	mh.HandleRequest(ReqGetFiles, func(req *p2p.Request) {
		homeDir, _ := os.UserHomeDir()
		entries, _ := os.ReadDir(path.Join(homeDir, "Downloads"))
		files := make([]PeerFile, 0, len(entries))

		for _, entry := range entries {
			filei, _ := entry.Info()
			files = append(files, PeerFile{
				Name: entry.Name(),
				Size: filei.Size(),
				Path: path.Join(homeDir, "Downloads", entry.Name()),
			})
		}

		enc := gob.NewEncoder(req.WriteCloser)
		enc.Encode(files)

		err := req.Close()
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error closing request: %s", err.Error())
		}

	})

	mh.HandleRequest(ReqDownloadFile, func(req *p2p.Request) {
		/* path, err := ioutil.ReadAll(req.Body)
		defer req.Close()

		if err != nil {
			runtime.LogErrorf(a.ctx, "Error reading request body: %s", err.Error())
			return
		}

		file, err := os.Open(string(path))
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error opening file: %s", err.Error())
			return
		}

		fileb, err := ioutil.ReadAll(file)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error reading file: %s", err.Error())
			return
		} */

		// req.Write(fileb)
	})
}

func (a *App) GetPeerSharedFiles(peerID string) ([]PeerFile, error) {
	res, err := a.p2p.SendMessage(a.ctx, peerID, ReqGetFiles, nil)
	if err != nil {
		return nil, err
	}

	var files []PeerFile
	dec := gob.NewDecoder(res.Body)
	err = dec.Decode(&files)

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (a *App) AddFiles() {
	files, _ := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})
	runtime.LogDebugf(a.ctx, "Selected files: %v", files)
}

type PeerFile struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Path string `json:"path"`
}

func (a *App) DownloadFile(peerID string, path string) {
	res, err := a.p2p.SendMessage(a.ctx, peerID, ReqDownloadFile, strings.NewReader(path))
	fmt.Println(res)

	if err != nil {
		runtime.LogErrorf(a.ctx, "Error sending request: %s", err.Error())
		return
	}

	var buf bytes.Buffer
	buf.ReadFrom(res.Body)
	runtime.LogInfof(a.ctx, "Response: %s", buf.Bytes())
}
