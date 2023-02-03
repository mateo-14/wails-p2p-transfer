package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"os"

	"github.com/mateo-14/wails-p2p-transfer/data"
	"github.com/mateo-14/wails-p2p-transfer/p2p"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Handle peers messages
func (a *App) onMessage(mh *p2p.MessageHandler) {
	mh.HandleRequest(p2p.ReqGetFiles, func(req *p2p.Request) {
		defer req.Close()
		sfiles, err := data.GetSharedFiles()
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error getting shared files: %s", err.Error())
			return
		}

		files := make([]data.PeerFile, 0, len(sfiles))
		for _, f := range sfiles {
			files = append(files, data.PeerFile{
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
	})

	mh.HandleRequest(p2p.ReqDownloadFile, func(req *p2p.Request) {
		defer req.Close()
		idb := make([]byte, 8)
		_, err := req.Body.Read(idb)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error reading request body: %s", err.Error())
			return
		}
		id := binary.BigEndian.Uint64(idb)

		file, err := data.GetSharedFile(int64(id))
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error getting shared file: %s", err.Error())
			return
		}

		f, err := os.Open(file.Path)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error opening file: %s", err.Error())
			return
		}

		_, err = req.Write(f)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error writing response: %s", err.Error())
		}
	})
}
