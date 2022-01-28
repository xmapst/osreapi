//go:build windows

package engine

import (
	"fmt"
	"github.com/xmapst/osreapi/utils"
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
	c.exec.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	//cmd.Dir, _ = os.UserHomeDir()
	go func() {
		output, err := c.exec.CombinedOutput()
		if err != nil {
			errorChan <- fmt.Errorf(string(output) + err.Error())
			return
		}
		// windows decode output
		output = utils.GbkToUtf8(output)
		outputChan <- string(output)
	}()
	select {
	// execute timeout signal
	case <-c.context.Done():
		// If you use cmd.Process.Kill() directly, only the child process is killed,
		//but the grandchild process is not killed
		//_ = cmd.Process.Kill()
		msg := "exec time out"
		if c.exec.Process != nil {
			err := KillAll(c.exec.Process.Pid)
			if err != nil {
				msg += err.Error()
			}
		}
		return 255, msg
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
