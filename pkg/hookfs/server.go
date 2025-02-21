package hookfs

import (
	"path/filepath"
	"time"

	// log "github.com/sirupsen/logrus"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/hanwen/go-fuse/v2/fuse/nodefs"
	"github.com/hanwen/go-fuse/v2/fuse/pathfs"
)

func newHookServer(hookfs *HookFs, directMount, debug bool) (*fuse.Server, error) {
	opts := &nodefs.Options{
		NegativeTimeout: time.Second,
		AttrTimeout:     time.Second,
		EntryTimeout:    time.Second,
	}
	pathFsOpts := &pathfs.PathNodeFsOptions{ClientInodes: true}
	pathFs := pathfs.NewPathNodeFs(hookfs, pathFsOpts)
	conn := nodefs.NewFileSystemConnector(pathFs.Root(), opts)
	originalAbs, _ := filepath.Abs(hookfs.Original)
	mOpts := &fuse.MountOptions{
		AllowOther:  true,
		Name:        hookfs.FsName,
		FsName:      originalAbs,
		DirectMount: directMount,
		Debug:       debug,
	}
	server, err := fuse.NewServer(conn.RawFS(), hookfs.Mountpoint, mOpts)
	if err != nil {
		return nil, err
	}

	if LogLevel() == LogLevelMax {
		server.SetDebug(true)
	}

	return server, nil
}
