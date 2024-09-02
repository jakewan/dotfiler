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
		expectedOps := []string{"symlink"}

		// Display intended file operations for user confirmation.
		for _, cfg := range manifestData {
			if !slices.Contains(expectedOps, cfg.Op) {
				colorError.Printf(
					"Unexpected file operation: %s (expected operations: %s)\n",
					cfg.Op,
					expectedOps,
				)
				return ErrInternal
			}
			displayTargetFileOperation(
				deps,
				cfg,
				sourceDir,
				c.destRootDir,
			)
		}
		colorConfirmation.Println("Apply these operations? (y/n)")
		if char, err := deps.GetSingleKey(); err != nil {
			return fmt.Errorf("getting keystroke: %w", err)
		} else if char != 'Y' && char != 'y' {
			colorWarning.Println("User cancelled")
			return ErrInternal
		}

		// Check destinations.
		for _, cfg := range manifestData {
			switch op := cfg.Op; op {
			case "symlink":
				if err := checkSymlinkDestination(
					deps,
					cfg,
					c.destRootDir,
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unexpected operation: %s", op)
			}
		}

		// Process operations.
		for _, cfg := range manifestData {
			switch op := cfg.Op; op {
			case "symlink":
				if err := updateSymlink(
					deps,
					cfg,
					sourceDir,
					c.destRootDir,
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unexpected operation: %s", op)
			}
		}
	}
	return nil
}

func displayTargetFileOperation(
	deps Dependencies,
	f FileConfig,
	srcDir string,
	destRootDir string,
) {
	foundOS := slices.Contains(f.TargetOS, deps.GetOS())
	foundArch := slices.Contains(f.TargetArch, deps.GetArch())
	if foundOS && foundArch {
		src := filepath.Join(srcDir, f.SrcFilePath)
		if _, err := os.Stat(src); err != nil {
			colorWarning.Println("Error retrieving file info:", err)
		}
		dest := filepath.Join(destRootDir, f.DstFilePath)
		fmt.Printf("Creating symlink from %s to %s\n", src, dest)
	}
}

func checkSymlinkDestination(
	deps Dependencies,
	f FileConfig,
	destRootDir string,
) error {
	foundOS := slices.Contains(f.TargetOS, deps.GetOS())
	foundArch := slices.Contains(f.TargetArch, deps.GetArch())
	dest := filepath.Join(destRootDir, f.DstFilePath)
	parentDirs := []string{}
	if foundOS && foundArch {
		if fi, err := os.Lstat(dest); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				fmt.Println("Creating new symlink", dest)
			} else {
				return fmt.Errorf("retrieving file info: %w", err)
			}
		} else if fi.Mode()&fs.ModeSymlink == fs.ModeSymlink {
			fmt.Println("Recreating existing symlink", dest)
		} else {
			colorError.Printf("Existing destination file %s is not a symlink.\n", dest)
			return ErrInternal
		}

		parentDir := filepath.Dir(dest)
		if _, err := os.Stat(parentDir); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				if !slices.Contains(parentDirs, parentDir) {
					parentDirs = append(parentDirs, parentDir)
				}
			} else {
				return fmt.Errorf("retrieving info for parent directory of %s: %s", parentDir, err)
			}
		}
	}
	if len(parentDirs) > 0 {
		for _, d := range parentDirs {
			fmt.Println("Parent directory will be created:", d)
		}
	}
	return nil
}

func updateSymlink(
	deps Dependencies,
	f FileConfig,
	srcDir string,
	destRootDir string,
) error {
	foundOS := slices.Contains(f.TargetOS, deps.GetOS())
	foundArch := slices.Contains(f.TargetArch, deps.GetArch())
	src := filepath.Join(srcDir, f.SrcFilePath)
	dest := filepath.Join(destRootDir, f.DstFilePath)

	// Check destination directories
	if foundOS && foundArch {
		removeFile := true
		if _, err := os.Lstat(dest); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				removeFile = false
			} else {
				return err
			}
		}
		if removeFile {
			if err := os.Remove(dest); err != nil {
				return fmt.Errorf("removing existing file: %w", err)
			}
		}
		if err := os.Symlink(src, dest); err != nil {
			return fmt.Errorf("creating symlink: %w", err)
		}
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
