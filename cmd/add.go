/*
Copyright © 2025 D7x7z49 <85430783+D7x7z49@users.noreply.github.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/D7x7z49/lite-tie/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <exec_path> [--alias <name>]",
	Short: "Add a symlink for a portable executable",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		execPath := args[0]
		alias, _ := cmd.Flags().GetString("alias")
		if err := handleAdd(execPath, alias); err != nil {
			fmt.Fprintf(os.Stderr, "Add failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().String("alias", "", "Optional alias for the link")
}

func handleAdd(execPath, alias string) error {
	absPath, err := filepath.Abs(execPath)
	if err != nil {
		return fmt.Errorf("invalid path: %v", err)
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("invalid path: %v", err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("path is a directory")
	}

	exeName := filepath.Base(absPath)
	linkName := alias
	if linkName == "" {
		linkName = exeName
	}
	linkPath := filepath.Join(filepath.Dir(os.Args[0]), linkName)

	if _, err := os.Lstat(linkPath); err == nil {
		fmt.Printf("Warning: %s exists, overwriting\n", linkName)
		os.Remove(linkPath)
	}

	if runtime.GOOS == "windows" {
		err := os.Symlink(absPath, linkPath)
		if err != nil && strings.Contains(err.Error(), "privilege") {
			batPath := linkPath + ".bat"
			batContent := fmt.Sprintf("@\"%s\" %%*", absPath)
			if err := os.WriteFile(batPath, []byte(batContent), 0755); err != nil {
				return fmt.Errorf("create .bat failed: %v", err)
			}
			linkPath = batPath
			fmt.Printf("Note: Created .bat file due to permission limitations\n")
		} else if err != nil {
			return fmt.Errorf("symlink failed: %v", err)
		}
	} else {
		if err := os.Symlink(absPath, linkPath); err != nil {
			return fmt.Errorf("symlink failed: %v", err)
		}
	}

	if err := config.AddEntry(linkName, absPath, linkPath); err != nil {
		os.Remove(linkPath)
		return err
	}

	fmt.Printf("Added: %s -> %s\n", linkName, absPath)
	return nil
}
