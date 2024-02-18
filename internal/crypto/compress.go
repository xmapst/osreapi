package crypto

import (
	"bytes"
	"compress/gzip"
	"io"
)

func Compress(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err := gz.Write(input)
	if err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decompress(compressed []byte) ([]byte, error) {
	buf := bytes.NewReader(compressed)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	decompressed, err := io.ReadAll(gz)
	if err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return decompressed, nil
}
