package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mateo-14/wails-p2p-transfer/p2p"

	"github.com/libp2p/go-libp2p/core/crypto"
)

// App struct
type App struct {
	ctx context.Context
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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func loadPrivateKey() (crypto.PrivKey, error) {
	file, err := os.ReadFile("private.key")
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPrivateKey(file)
}

func storePrivateKey(privk crypto.PrivKey) error {
	privkb, err := crypto.MarshalPrivateKey(privk)
	if err != nil {
		return err
	}

	return os.WriteFile("private.key", privkb, 0644)
}

func generatePrivateKey() (crypto.PrivKey, error) {
	privk, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	return privk, err
}

type HostData struct {
	Addr string `json:"address"`
	Id   string `json:"id"`
}

func (a *App) StartP2P() (*HostData, error) {
	privk, err := loadPrivateKey()
	if err != nil {
		privk, err = generatePrivateKey()
		if err != nil {
			panic(err)
		}

		err = storePrivateKey(privk)
		if err != nil {
			panic(err)
		}
	}

	p2p, err := p2p.New(privk)
	if err != nil {
		return nil, err
	}

	return &HostData{
		Addr: p2p.Addrs()[0].String(),
		Id:   p2p.Id(),
	}, nil

}
