package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/mateo-14/wails-p2p-transfer/p2p"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	mh.HandleRequest(p2p.GetFiles, func(req *p2p.MessageRequest) {
		homeDir, _ := os.UserHomeDir()
		entries, _ := os.ReadDir(path.Join(homeDir, "Downloads"))
		files := make([]string, 0, len(entries))

		for _, entry := range entries {
			files = append(files, entry.Name())
		}

		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		enc.Encode(files)

		req.Write(buf.Bytes())
	})
}

func (a *App) GetPeerSharedFiles(peerID string) ([]string, error) {
	res, err := a.p2p.SendMessage(a.ctx, peerID, p2p.GetFiles, nil)
	if err != nil {
		return nil, err
	}

	var files []string
	dec := gob.NewDecoder(bytes.NewReader(res.Payload))
	err = dec.Decode(&files)

	if err != nil {
		return nil, err
	}

	return files, nil

}
