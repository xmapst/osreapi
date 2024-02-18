package crypto

import (
	"bytes"
	"errors"
)

// 循环移位加解密算法

var ErrEmptyKey = errors.New("empty key")

// Encrypt 加密
func Encrypt(key []byte, data []byte) ([]byte, error) {
	if key == nil {
		return nil, ErrEmptyKey
	}
	r := &rotate{
		data: data,
		key:  key,
		dL:   len(data),
		kL:   len(key),
	}
	return r.encrypt(), nil
}

// Decrypt 解密
func Decrypt(key []byte, data []byte) ([]byte, error) {
	if key == nil {
		return nil, ErrEmptyKey
	}
	r := &rotate{
		data: data,
		key:  key,
		dL:   len(data),
		kL:   len(key),
	}
	return r.decrypt(), nil
}

type rotate struct {
	data []byte // 数据
	key  []byte // 密钥
	dL   int    // 数据长度
	kL   int    // key长度
	pL   int    // 补齐的数据长度
}

// 补齐数据长度
func (r *rotate) complete() {
	paddedData := r.data
	if len(r.data)%r.kL > 0 {
		r.pL = r.kL - (len(r.data) % r.kL)
		paddedData = append(r.data, bytes.Repeat([]byte{0}, r.pL)...)
	}
	r.data = paddedData
	r.dL = len(paddedData)
}

func (r *rotate) encrypt() []byte {
	// 进行补全, 全部按密钥长度的倍数处理
	r.complete()
	result := make([]byte, r.dL)
	for i := 0; i < r.dL; i += r.kL {
		for j := 0; j < r.kL; j++ {
			if i+j < r.dL {
				// 异或操作
				result[i+j] = r.data[i+j] ^ r.key[j%r.kL]
			}
		}
		// 右移
		r.shiftRight(result[i : i+r.kL])
	}
	// 去除补全的字符
	return result[:len(result)-r.pL]
}

func (r *rotate) decrypt() []byte {
	// 进行补全, 全部按密钥长度的倍数处理
	r.complete()
	result := make([]byte, r.dL)
	for i := 0; i < r.dL; i += r.kL {
		// 左移
		r.shiftLeft(r.data[i : i+r.kL])
		for j := 0; j < r.kL; j++ {
			if i+j < r.dL {
				// 异或操作
				result[i+j] = r.data[i+j] ^ r.key[j%r.kL]
			}
		}
	}
	// 去除补全的字符
	return result[:len(result)-r.pL]
}

// 右移
func (r *rotate) shiftRight(data []byte) {
	for i := 0; i < len(data); i++ {
		shift := int(r.key[i%r.kL]) % 8
		data[i] = (data[i] >> shift) | (data[i] << (8 - shift))
	}
}

// 左移
func (r *rotate) shiftLeft(data []byte) {
	for i := 0; i < len(data); i++ {
		shift := int(r.key[i%r.kL]) % 8
		data[i] = (data[i] << shift) | (data[i] >> (8 - shift))
	}
}
