package workspace

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/router/base"
	"github.com/xmapst/osreapi/internal/router/types"
	"github.com/xmapst/osreapi/internal/server/config"
	"github.com/xmapst/osreapi/internal/utils"
)

// Get
// @Summary task workspace list
// @description task workspace list
// @Tags Task
// @Accept application/json
// @Accept application/toml
// @Accept application/x-yaml
// @Accept multipart/form-data
// @Produce application/json
// @Produce application/x-yaml
// @Produce application/toml
// @Param task path string true "task id"
// @Param path query string false "dir path"
// @Success 200 {object} types.BaseRes
// @Failure 500 {object} types.BaseRes
// @Router /api/v1/task/{task}/workspace [get]
func Get(c *gin.Context) {
	render := base.Gin{Context: c}
	task := c.Param("task")
	if task == "" {
		render.SetError(base.CodeErrNoData, errors.New("task does not exist"))
		return
	}
	prefix := filepath.Join(config.App.WorkSpace, task)
	if !utils.FileOrPathExist(prefix) {
		render.SetError(base.CodeErrNoData, errors.New("task does not exist"))
		return
	}
	path := filepath.Join(prefix, utils.PathEscape(c.Query("path")))
	file, err := os.Open(path)
	if err != nil {
		logx.Errorln(err)
		render.SetError(base.CodeErrNoData, err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		logx.Errorln(err)
		render.SetError(base.CodeErrNoData, err)
		return
	}
	if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
		finalPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			logx.Errorln(err)
			render.SetError(base.CodeErrNoData, err)
			return
		}
		fileInfo, err = os.Lstat(finalPath)
		if err != nil {
			logx.Errorln(err)
			render.SetError(base.CodeErrNoData, err)
			return
		}
	}
	if !fileInfo.IsDir() {
		ctype := mime.TypeByExtension(fileInfo.Name())
		if ctype == "" {
			ctype = "application/octet-stream"
		}
		c.Header("Content-Type", ctype)
		c.Header("Transfer-Encoding", "chunked")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileInfo.Name()))
		_, _ = io.Copy(c.Writer, file)
		return
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		logx.Errorln(err)
		render.SetError(base.CodeErrNoData, err)
		return
	}
	var scheme = "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	var infos = new(types.FileListRes)
	var uriPrefix = fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, strings.TrimSuffix(c.Request.URL.Path, "/"))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		_path := filepath.Join(path, info.Name())
		isDir := info.IsDir()
		size := info.Size()
		if info.Mode()&os.ModeSymlink != 0 {
			s, err := os.Stat(_path)
			if err == nil {
				isDir = s.IsDir()
				size = s.Size()
			}
		}
		_path = strings.TrimPrefix(filepath.ToSlash(_path), filepath.ToSlash(prefix))
		infos.Files = append(infos.Files, &types.FileRes{
			URL:     fmt.Sprintf("%s?path=%s", uriPrefix, _path),
			Name:    info.Name(),
			Path:    _path,
			Size:    size,
			Mode:    info.Mode().String(),
			ModTime: info.ModTime().Unix(),
			IsDir:   isDir,
		})
	}
	infos.Total = len(infos.Files)
	render.SetRes(infos)
}
