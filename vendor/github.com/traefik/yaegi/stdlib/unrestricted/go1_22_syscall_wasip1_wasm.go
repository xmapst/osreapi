// Code generated by 'yaegi extract syscall'. DO NOT EDIT.

//go:build go1.22
// +build go1.22

package unrestricted

import (
	"reflect"
	"syscall"
)

func init() {
	Symbols["syscall/syscall"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Exit":         reflect.ValueOf(syscall.Exit),
		"Kill":         reflect.ValueOf(syscall.Kill),
		"ProcExit":     reflect.ValueOf(syscall.ProcExit),
		"RawSyscall":   reflect.ValueOf(syscall.RawSyscall),
		"RawSyscall6":  reflect.ValueOf(syscall.RawSyscall6),
		"Shutdown":     reflect.ValueOf(syscall.Shutdown),
		"StartProcess": reflect.ValueOf(syscall.StartProcess),
		"Syscall":      reflect.ValueOf(syscall.Syscall),
		"Syscall6":     reflect.ValueOf(syscall.Syscall6),
	}
}
