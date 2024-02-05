package utils

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func FileOrPathExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func EnsureDirExist(name string) error {
	if !FileOrPathExist(name) {
		return os.MkdirAll(name, os.ModePerm)
	}
	return nil
}

func SliceToStrMap(s []string) map[string]string {
	m := make(map[string]string)
	for _, v := range s {
		slice := strings.Split(v, "=")
		switch {
		case len(slice) > 2:
			m[slice[0]] = strings.Join(slice[1:], "=")
		case len(slice) == 2:
			m[slice[0]] = slice[1]
		case len(slice) == 1:
			m[v] = ""
		}
	}
	return m
}

func ClearDir(path string) {
	_ = os.RemoveAll(path)
	_ = EnsureDirExist(path)
}

func TimeStr(nsec int64) string {
	if nsec == 0 {
		return ""
	}
	return time.Unix(0, nsec).Format(time.RFC3339)
}

func PathEscape(s string) string {
	s = filepath.Clean(strings.TrimPrefix(s, ".."))
	if s == ".." {
		return ""
	}
	if !strings.HasPrefix(s, "..") {
		return s
	}
	return PathEscape(s)
}
