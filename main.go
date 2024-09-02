package main

import (
	"context"
	"errors"
	"log"
	"os"
	"runtime"

	"github.com/eiannone/keyboard"
	"github.com/jakewan/dotfiler/internal/cmd"
)

func init() {
	// Omit date and time from log messages.
	log.SetFlags(0)
}

type dependencies struct{}

// GetSingleKey implements cmd.Dependencies.
func (d *dependencies) GetSingleKey() (ch rune, err error) {
	if ch, _, err := keyboard.GetSingleKey(); err != nil {
		return 0, err
	} else {
		return ch, nil
	}
}

// GetHomeDirectory implements cmd.Dependencies.
func (d *dependencies) GetHomeDirectory() (string, error) {
	return os.UserHomeDir()
}

// GetArch implements cmd.Dependencies.
func (d *dependencies) GetArch() string {
	return runtime.GOARCH
}

// GetOS implements cmd.Dependencies.
func (d *dependencies) GetOS() string {
	return runtime.GOOS
}

func main() {
	// Determine the root directory (e.g., the user's home directory).
	ctx := context.Background()
	deps := dependencies{}
	if err := cmd.Execute(ctx, os.Args[1:], &deps); err != nil {
		if errors.Is(err, cmd.ErrInternal) {
			os.Exit(1)
		} else {
			log.Fatalf("The operation resulted in an error: %s", err)
		}
	}
}
