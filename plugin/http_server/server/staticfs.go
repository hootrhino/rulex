package server

import (
	"embed"
	"io/fs"
	"net/http"
	"path"

	"github.com/gin-contrib/static"
)

// --------------------------------------------------------------------------------
// 前端静态资源打包目录
// --------------------------------------------------------------------------------
//
//go:embed  www/*
var files embed.FS

//go:embed  www/index.html
var indexHTML []byte

type WWWFS struct {
	http.FileSystem
}

func (f WWWFS) Exists(prefix string, filepath string) bool {
	_, err := f.Open(path.Join(prefix, filepath))
	return err == nil
}

func WWWRoot(dir string) static.ServeFileSystem {
	if sub, err := fs.Sub(files, path.Join("www", dir)); err == nil {
		return WWWFS{http.FS(sub)}
	}
	return nil
}
