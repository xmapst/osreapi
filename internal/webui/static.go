package webui

import (
	"embed"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/xmapst/osreapi/internal/logx"
)

//go:embed static/*
var fileSystem embed.FS

const (
	dir   = "static"
	index = "index.html"
)

// File 使用go.16新的特性embed 到包前端编译后的代码. 替代nginx.   one binary rules them all
func File(prefixs ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := strings.TrimSpace(c.Request.URL.Path)
		for _, prefix := range prefixs {
			if strings.HasPrefix(path, prefix) {
				// 处理包含指定前缀的
				fileHandler(c, prefix)
			}
		}
	}
}

func fileHandler(c *gin.Context, prefix string) {
	path := strings.TrimSpace(c.Request.URL.Path)
	var err error
	var content []byte
	defer func() {
		// 统一处理异常
		if err != nil {
			// 交由前端处理路由
			content, err = fileSystem.ReadFile(filepath.ToSlash(filepath.Join(dir, index)))
			if err != nil {
				logx.Warnln(err)
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.Header("Cache-Control", "no-cache")
			c.Header("Content-Type", "text/html; charset=utf-8")
			_, _ = c.Writer.Write(content)
			c.Abort()
			return
		}
	}()
	path = filepath.ToSlash(filepath.Join(dir, strings.TrimPrefix(path, prefix)))
	if path == "" {
		path = index
	}
	// 尝试打开文件
	var file fs.File
	file, err = fileSystem.Open(path)
	if err != nil {
		file, err = fileSystem.Open(filepath.ToSlash(filepath.Dir(path)))
		if err != nil {
			logx.Warnln(err)
			return
		}
	}
	var fileInfo fs.FileInfo
	fileInfo, err = file.Stat()
	if err != nil {
		logx.Warnln(err)
		return
	}
	// 如果是目录则返回index.html, 没有index.html则返回404
	if fileInfo.IsDir() {
		content, err = fileSystem.ReadFile(filepath.ToSlash(filepath.Join(path, index)))
		if err != nil {
			logx.Warnln(err)
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Header("Content-Type", "text/html; charset=utf-8")
		_, _ = c.Writer.Write(content)
		c.Abort()
		return
	}
	content, err = fileSystem.ReadFile(path)
	if err != nil {
		logx.Warnln(err)
		return
	}
	// 设置响应类型
	c.Header("Cache-Control", "no-cache")
	ctype := mime.TypeByExtension(filepath.Ext(path))
	if ctype == "" {
		ctype = http.DetectContentType(content)
	}
	c.Header("Content-Type", ctype)

	c.Status(200)
	_, _ = c.Writer.Write(content)
	c.Abort()
}
