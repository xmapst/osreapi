//go:build !windows

package engine

import (
	"fmt"
	"syscall"
)

func (c *Cmd) ExecScript(execID string) (exitCode int, content string) {
	outputChan := make(chan string)
	errorChan := make(chan error)
	defer c.clear()
	if !c.initCmd(execID) {
		return 255, "command_type not found"
	}
	defer c.cancelFunc()
	c.exec.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	//cmd.Dir, _ = os.UserHomeDir()
	go func() {
		output, err := c.exec.CombinedOutput()
		if err != nil {
			errorChan <- fmt.Errorf(string(output) + err.Error())
			return
		}
		outputChan <- string(output)
	}()
	select {
	// execute timeout signal
	case <-c.context.Done():
		// If you use cmd.Process.Kill() directly, only the child process is killed,
		//but the grandchild process is not killed
		//err := cmd.Process.Kill()
		if c.exec.Process != nil {
			_ = syscall.Kill(-c.exec.Process.Pid, syscall.SIGKILL)
		}
		return 255, "exec time out"
	// execution result output
	case output := <-outputChan:
		code := 0
		if c.exec.ProcessState != nil {
			code = c.exec.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
		}
		return code, output
	// Execute exception output
	case err := <-errorChan:
		code := 255
		if c.exec.ProcessState != nil {
			code = c.exec.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
		}
		return code, err.Error()
	}
}
