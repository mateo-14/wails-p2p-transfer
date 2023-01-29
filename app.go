package main

import (
	"context"
	"errors"
	"fmt"

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

	appn := &AppNotifiee{}
	p2p, err := p2p.New(appn)

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

	runtime.LogInfof(a.ctx, "Connected to node %s\n", nodeaddr)
	return nil
	/*
		 	time.Sleep(5 * time.Second)
			fmt.Println("Connected to node")
			go func() {
				time.Sleep(10 * time.Second)
				runtime.EventsEmit(a.ctx, "peer:disconnected", id)
			}()
			return nil
	*/
}

func (a *App) OnFrontendLoad() (*p2p.HostData, error) {
	return a.startP2P()
}
