package hook

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Exec Execute a command on a go routine.
func Exec(command string) {
	if command == "" {
		return
	}

	go func() {
		errH := execute(command)
		if errH != nil {
			panic(errH)
		}
	}()
}

func execute(command string) error {
	ctxCmd, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	parts := strings.Fields(os.ExpandEnv(command))
	output, err := exec.CommandContext(ctxCmd, parts[0], parts[1:]...).CombinedOutput()
	if len(output) > 0 {
		fmt.Println(string(output))
	}

	if errors.Is(ctxCmd.Err(), context.DeadlineExceeded) {
		return errors.New("hook timed out")
	}

	return err
}
