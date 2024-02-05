// Code generated by 'yaegi extract unicode/utf16'. DO NOT EDIT.

//go:build go1.19 && !go1.20
// +build go1.19,!go1.20

package stdlib

import (
	"reflect"
	"unicode/utf16"
)

func init() {
	Symbols["unicode/utf16/utf16"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Decode":      reflect.ValueOf(utf16.Decode),
		"DecodeRune":  reflect.ValueOf(utf16.DecodeRune),
		"Encode":      reflect.ValueOf(utf16.Encode),
		"EncodeRune":  reflect.ValueOf(utf16.EncodeRune),
		"IsSurrogate": reflect.ValueOf(utf16.IsSurrogate),
	}
}
