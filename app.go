package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jaevor/go-nanoid"
	"github.com/mateo-14/wails-p2p-transfer/data"
	"github.com/mateo-14/wails-p2p-transfer/p2p"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	p2p            *p2p.P2P
	db             *sql.DB
	frontendLoaded bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

var idgen func() string

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.db = data.Init()

	p2p, err := p2p.NewP2P(a.ctx, a.onMessage)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error starting P2P: %s\n", err.Error())
		return
	}

	a.p2p = p2p
	p2p.Notify(&AppNotifiee{
		ctx: a.ctx,
	})

	idgen, _ = nanoid.Standard(6)
}

func (a *App) shutdown(ctx context.Context) {
	a.db.Close()
}

// Frontend methods
func (a *App) ConnectToNode(addr string, id string) error {
	nodeaddr := fmt.Sprintf("%s/%s", addr, id)

	err := a.p2p.Connect(a.ctx, nodeaddr)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error connecting to node %s: %s\n", nodeaddr, err.Error())
		return err
	}

	return nil
}

func (a *App) OnFrontendLoad() (*data.InitialData, error) {
	a.frontendLoaded = true
	sharedFiles, err := data.GetSharedFiles()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error getting shared files: %s\n", err.Error())
		return nil, err
	}

	peers, err := data.GetPeers()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error getting peers: %s\n", err.Error())
		return nil, err
	}

	initialData := &data.InitialData{
		SharedFiles: sharedFiles,
		Peers:       peers,
	}

	if a.p2p != nil {
		initialData.HostData = a.p2p.GetHostData()
	}

	return initialData, nil
}

func (a *App) GetPeerSharedFiles(peerID string) ([]data.PeerFile, error) {
	res, err := a.p2p.SendMessage(a.ctx, peerID, p2p.ReqGetFiles, nil)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var files []data.PeerFile
	dec := gob.NewDecoder(res.Body)
	err = dec.Decode(&files)

	if err != nil {
		runtime.LogPrintf(a.ctx, "Error decoding response: %s", err.Error())
		return nil, err
	}

	return files, nil
}

func (a *App) AddFiles() (string, error) {
	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error opening file dialog: %s", err.Error())
		return "", err
	}

	if (len(files)) == 0 {
		return "", nil
	}

	eventid := idgen()

	go func() {
		filesd := make([]data.File, 0, len(files))
		for _, path := range files {
			_, err := os.Stat(path)
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

			n, err := io.Copy(hash, file)
			if err != nil {
				runtime.LogErrorf(a.ctx, "Error hashing file: %s", err.Error())
				continue
			}

			hashs := hex.EncodeToString(hash.Sum(nil))

			var filed data.File
			filed.Hash = hashs
			filed.Name = filepath.Base(path)
			filed.Size = n
			filed.Path = path

			id, err := data.AddSharedFile(path, filed.Size, hashs)
			if err == nil {
				filed.ID = id
				filesd = append(filesd, filed)
			}
		}

		runtime.EventsEmit(a.ctx, eventid, filesd)
	}()

	return eventid, nil
}

func (a *App) RemoveSharedFile(id int64) error {
	err := data.RemoveSharedFile(id)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error removing shared file: %s", err.Error())
		return err
	}

	return nil
}

func (a *App) DownloadFile(peerID string, id uint64) error {
	fmt.Println("Downloading file", id, "from", peerID)
	idb := make([]byte, 8)
	binary.LittleEndian.PutUint64(idb, id)
	res, err := a.p2p.SendMessage(a.ctx, peerID, p2p.ReqDownloadFile, bytes.NewReader(idb))

	if err != nil {
		runtime.LogErrorf(a.ctx, "Error sending request: %s", err.Error())
		return err
	}

	filei := &data.PeerFile{}
	dec := gob.NewDecoder(res.Body)
	err = dec.Decode(&filei)

	if err != nil {
		runtime.LogErrorf(a.ctx, "Error decoding response: %s", err.Error())
		return err
	}

	defer res.Body.Close()

	os.Mkdir("downloads", os.ModePerm)

	path := filepath.Join("downloads", filei.Name)
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(path)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error creating file: %s", err.Error())
			return err
		}

		defer f.Close()

		_, err = io.CopyN(f, res.Body, filei.Size)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error writing file: %s", err.Error())
			return err
		}
	} else {
		runtime.LogErrorf(a.ctx, "File already exists: %s", path)
		return errors.New("file already exists")
	}

	return nil
}

func (a *App) AddPeer(name string, addr string) (*data.Peer, error) {
	lastslashi := strings.LastIndex(addr, "/") + 1
	id := addr[lastslashi:]
	addr = addr[:lastslashi-1]

	idi, err := data.AddPeer(name, addr, id)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Error adding peer: %s", err.Error())
		return nil, err
	}

	return &data.Peer{
		ID:      idi,
		Name:    name,
		Address: addr,
		PeerID:  id,
	}, nil
}
