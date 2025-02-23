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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/D7x7z49/lite-tie/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [<name> ...] [--silent] [--clean]",
	Short: "Remove symlinks",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		silent, _ := cmd.Flags().GetBool("silent")
		clean, _ := cmd.Flags().GetBool("clean")

		if len(args) == 0 && !clean {
			fmt.Println("Error: Provide at least one <name> or use --clean")
			cmd.Usage()
			os.Exit(1)
		}

		if len(args) > 0 {
			if err := handleRemove(args, silent); err != nil {
				fmt.Fprintf(os.Stderr, "Remove failed: %v\n", err)
				os.Exit(1)
			}
		}

		if clean {
			if err := handleClean(); err != nil {
				fmt.Fprintf(os.Stderr, "Clean failed: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().Bool("silent", false, "Skip confirmation")
	removeCmd.Flags().Bool("clean", false, "Remove all unavailable entries")
}

func handleRemove(names []string, silent bool) error {
	entries, err := config.GetEntries()
	if err != nil {
		return err
	}

	var toDelete []string
	for _, name := range names {
		entry, exists := entries[name]
		if !exists {
			fmt.Printf("%s not found\n", name)
			continue
		}

		if entry.Available && !silent {
			fmt.Printf("Delete %s? (yes/no): ", name)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "yes" && response != "y" {
				fmt.Printf("Skipped %s\n", name)
				continue
			}
		}

		if err := os.Remove(entry.Link); err != nil && !os.IsNotExist(err) {
			fmt.Printf("Failed to remove %s: %v\n", name, err)
		} else {
			fmt.Printf("Removed %s\n", name)
			toDelete = append(toDelete, name)
		}
	}

	if len(toDelete) > 0 {
		if err := config.RemoveEntries(toDelete); err != nil {
			fmt.Printf("Registry update failed: %v\n", err)
			return err
		}
	}
	return nil
}

func handleClean() error {
	entries, err := config.GetEntries()
	if err != nil {
		return err
	}

	var toDelete []string
	for name, entry := range entries {
		if !entry.Available {
			if err := os.Remove(entry.Link); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Failed to remove %s: %v\n", name, err)
			} else {
				fmt.Printf("Cleaned %s\n", name)
				toDelete = append(toDelete, name)
			}
		}
	}

	if len(toDelete) > 0 {
		if err := config.RemoveEntries(toDelete); err != nil {
			fmt.Printf("Registry update failed: %v\n", err)
			return err
		}
	}
	return nil
}
