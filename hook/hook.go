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
func Exec(ctx context.Context, command string) {
	if command == "" {
		return
	}

	go func() {
		errH := execute(ctx, command)
		if errH != nil {
			panic(errH)
		}
	}()
}

func execute(ctx context.Context, command string) error {
	ctxCmd, cancel := context.WithTimeout(ctx, 30*time.Second)
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
