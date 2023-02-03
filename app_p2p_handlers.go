package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
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
		var id uint64
		err := binary.Read(req.Body, binary.LittleEndian, &id)

		if err != nil {
			runtime.LogErrorf(a.ctx, "Error decoding file id: %s", err.Error())
			return
		}

		file, err := data.GetSharedFile(int64(id))
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error getting shared file: %s", err.Error())
			return
		}

		peerf := data.PeerFile{
			Name: file.Name,
			Size: file.Size,
			ID:   file.ID,
			Hash: file.Hash,
		}

		f, err := os.Open(file.Path)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error opening file: %s", err.Error())
			return
		}

		// Send file info
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		enc.Encode(peerf)

		_, err = req.Write(&buf)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error writing file info response: %s", err.Error())
			return
		}

		// req.Write(bytes.NewBufferString("-=#-=#-=#-=#-=#-=#-=#-=#"))
		// Send file
		w, err := req.Write(f)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Error writing file response: %s", err.Error())
		}
		fmt.Println(w, file.Size)
	})
}
