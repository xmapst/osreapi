// Code generated by 'yaegi extract encoding/hex'. DO NOT EDIT.

//go:build go1.20
// +build go1.20

package stdlib

import (
	"encoding/hex"
	"reflect"
)

func init() {
	Symbols["encoding/hex/hex"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Decode":         reflect.ValueOf(hex.Decode),
		"DecodeString":   reflect.ValueOf(hex.DecodeString),
		"DecodedLen":     reflect.ValueOf(hex.DecodedLen),
		"Dump":           reflect.ValueOf(hex.Dump),
		"Dumper":         reflect.ValueOf(hex.Dumper),
		"Encode":         reflect.ValueOf(hex.Encode),
		"EncodeToString": reflect.ValueOf(hex.EncodeToString),
		"EncodedLen":     reflect.ValueOf(hex.EncodedLen),
		"ErrLength":      reflect.ValueOf(&hex.ErrLength).Elem(),
		"NewDecoder":     reflect.ValueOf(hex.NewDecoder),
		"NewEncoder":     reflect.ValueOf(hex.NewEncoder),

		// type definitions
		"InvalidByteError": reflect.ValueOf((*hex.InvalidByteError)(nil)),
	}
}
