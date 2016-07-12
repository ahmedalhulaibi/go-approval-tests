package reporters

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
)

type clipboard struct{}

// NewClipboardReporter copies move file command to your clipboard
func NewClipboardReporter() Reporter {
	return &clipboard{}
}

func (s *clipboard) Report(approved, received string) bool {
	receivedFull, _ := filepath.Abs(received)
	approvedFull, _ := filepath.Abs(approved)

	switch runtime.GOOS {
	case "windows":
		move := fmt.Sprintf("move /Y \"%s\" \"%s\"", receivedFull, approvedFull)
		return copyToWindowsClipboard(move)
	default:
		move := fmt.Sprintf("mv %s %s", receivedFull, approvedFull)
		return copyToDarwinClipboard(move)
	}

	return false
}

func copyToWindowsClipboard(text string) bool {
	return pipeToProgram("clip", text)
}

func copyToDarwinClipboard(text string) bool {
	return pipeToProgram("pbcopy", text)
}

func pipeToProgram(programName, text string) bool {
	c := exec.Command(programName)
	pipe, err := c.StdinPipe()
	if err != nil {
		fmt.Printf("StdinPipe: err=%s", err)
		return false
	}
	pipe.Write([]byte(text))
	pipe.Close()

	c.Start()
	return true
}
