package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type cmdFilesUpdate struct {
	flags            *flag.FlagSet
	manifestFilePath string
	destRootDir      string
}

// init implements subcommandRunner.
func (c *cmdFilesUpdate) init(args []string) (subcommandRunner, error) {
	if err := c.flags.Parse(args); err != nil {
		return nil, err
	}
	return c, nil
}

// name implements subcommandRunner.
func (c *cmdFilesUpdate) name() string {
	return "update"
}

// run implements subcommandRunner.
func (c *cmdFilesUpdate) run(_ context.Context, deps Dependencies) error {
	fmt.Println("Current OS:", deps.GetOS())
	fmt.Println("Current Arch:", deps.GetArch())

	// Determine the manifest file path.
	if c.manifestFilePath == "" {
		if n, err := filepath.Abs(manifestFileName); err != nil {
			return err
		} else {
			c.manifestFilePath = n
		}
	} else if n, err := filepath.Abs(c.manifestFilePath); err != nil {
		return err
	} else {
		c.manifestFilePath = n
	}
	if !strings.HasSuffix(c.manifestFilePath, manifestFileName) {
		colorWarning.Println(
			"The given manifest file path was not terminated with the manifest filename. We assume it is the containing directory.",
		)
		c.manifestFilePath = filepath.Join(c.manifestFilePath, manifestFileName)
	}
	fmt.Println("Manifest file path:", c.manifestFilePath)
	if _, err := os.Stat(c.manifestFilePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			colorError.Println(
				"A manifest file was not found at the given location.",
			)
			return ErrInternal
		} else {
			return fmt.Errorf("reading manifest file info: %w", err)
		}
	}

	// Determine the destination root directory
	if c.destRootDir == "" {
		if d, err := deps.GetHomeDirectory(); err != nil {
			return fmt.Errorf("determining user home directory: %w", err)
		} else {
			c.destRootDir = d
		}
	}
	if c.destRootDir == "" {
		colorError.Println("Could not determine the destination directory.")
		return ErrInternal
	}
	fmt.Println("Destination root directory:", c.destRootDir)

	// Process the manifest.
	if f, err := os.Open(c.manifestFilePath); err != nil {
		return fmt.Errorf("opening manifest file: %s", err)
	} else if data, err := io.ReadAll(f); err != nil {
		return fmt.Errorf("reading manifest file: %s", err)
	} else {
		var manifestData []FileConfig
		if err := yaml.Unmarshal(data, &manifestData); err != nil {
			return fmt.Errorf("deserializing manifest file: %w", err)
		}
		sourceDir := filepath.Dir(c.manifestFilePath)
		for _, cfg := range manifestData {
			if err := displayTargetFileOperation(
				deps,
				cfg,
				sourceDir,
				c.destRootDir,
			); err != nil {
				colorError.Printf(
					"An error occurred processing file configuration: %+v\n\nProcessing halted.\n",
					cfg,
				)
				return ErrInternal
			}
		}
		colorConfirmation.Println("Apply these operations? (y/n)")
		if char, err := deps.GetSingleKey(); err != nil {
			return fmt.Errorf("getting keystroke: %w", err)
		} else if char != 'Y' && char != 'y' {
			colorWarning.Println("User cancelled")
			return ErrInternal
		}
	}
	return nil
}

func displayTargetFileOperation(
	deps Dependencies,
	f FileConfig,
	srcDir string,
	destRootDir string,
) error {
	foundOS := slices.Contains(f.TargetOS, deps.GetOS())
	foundArch := slices.Contains(f.TargetArch, deps.GetArch())
	if foundOS && foundArch {
		fmt.Printf("Processing %v\n", f)
	}
	return nil
}

func newCmdFilesUpdate() subcommandRunner {
	result := &cmdFilesUpdate{
		flags: flag.NewFlagSet("update", flag.ExitOnError),
	}
	result.flags.StringVar(&result.manifestFilePath, "manifest", "", flagUsageManifest)
	result.flags.StringVar(&result.manifestFilePath, "m", "", flagUsageManifest)
	result.flags.StringVar(&result.destRootDir, "destinationrootdir", "", flagUsageDestinationRootDir)
	result.flags.StringVar(&result.destRootDir, "destination", "", flagUsageDestinationRootDir)
	result.flags.StringVar(&result.destRootDir, "d", "", flagUsageDestinationRootDir)
	return result
}
