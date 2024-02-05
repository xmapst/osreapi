package base

import (
	"archive/tar"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/utils"
)

// SaveToTarFile 保存为tar压缩文件
func (g *Gin) SaveToTarFile(path string) error {
	form, err := g.MultipartForm()
	if err != nil {
		logx.Errorln(err)
		return err
	}
	files := form.File["files"]
	if len(files) == 0 {
		return fmt.Errorf("files is null")
	}
	fw, err := os.Create(path)
	if err != nil {
		logx.Errorln(err)
		return err
	}

	tw := tar.NewWriter(fw)
	defer func() {
		_ = tw.Close()
		_ = fw.Close()
	}()
	for _, f := range files {
		// file.Filename does not contain the directory path
		// RFC 7578, Section 4.2 requires that if a filename is provided, the
		// directory path information must not be used.
		// 🙂🙂🙂🙂🙂🙂
		v := f.Header.Get("Content-Disposition")
		_, dispositionParams, err := mime.ParseMediaType(v)
		if err != nil {
			logx.Errorln(err)
			return err
		}
		fileName, ok := dispositionParams["filename"]
		if !ok {
			return fmt.Errorf("filename does not exist")
		}

		hdr := &tar.Header{
			Name: fileName,
			Mode: 0644,
			Size: f.Size,
		}
		if err = tw.WriteHeader(hdr); err != nil {
			logx.Errorln(err)
			return err
		}
		_f, err := f.Open()
		if err != nil {
			logx.Errorln(err)
			return err
		}
		if _, err = io.Copy(tw, _f); err != nil {
			logx.Errorln(err)
			return err
		}
	}
	return nil
}

func (g *Gin) SaveFiles(path string) error {
	form, err := g.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	if len(files) == 0 {
		return fmt.Errorf("files is null")
	}
	var wg = new(sync.WaitGroup)
	var errs []error
	var lock sync.Mutex

	// ensure dir exist
	if err = utils.EnsureDirExist(path); err != nil {
		logx.Errorln(err)
		return err
	}

	for _, f := range files {
		wg.Add(1)
		go func(path string, file *multipart.FileHeader) {
			defer wg.Done()
			// file.Filename does not contain the directory path
			// RFC 7578, Section 4.2 requires that if a filename is provided, the
			// directory path information must not be used.
			// 🙂🙂🙂🙂🙂🙂
			v := file.Header.Get("Content-Disposition")
			_, dispositionParams, err := mime.ParseMediaType(v)
			if err != nil {
				logx.Errorln(err)
				return
			}
			fileName, ok := dispositionParams["filename"]
			if !ok {
				logx.Errorln("filename does not exist")
				return
			}

			file.Filename = fileName
			// Default save path
			uploadFileName := filepath.Base(file.Filename)
			uploadFPath := filepath.Dir(file.Filename)
			// Process folder upload
			if uploadFPath != "." {
				path = filepath.Join(path, uploadFPath)
				if err = utils.EnsureDirExist(path); err != nil {
					logx.Errorln(err)
					return
				}
			}

			// save file to local in _tp
			if err = g.SaveUploadedFile(file, filepath.Join(path, uploadFileName)); err != nil {
				lock.Lock()
				defer lock.Unlock()
				errs = append(errs, fmt.Errorf("%s %v", file.Filename, err))
			}
		}(path, f)
	}
	wg.Wait()
	if errs == nil {
		return nil
	}
	return fmt.Errorf("%v", errs)
}
