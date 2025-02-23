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
	"text/tabwriter"

	"github.com/D7x7z49/lite-tie/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [<name>] [--simple]",
	Short: "List all symlinks",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		simple, _ := cmd.Flags().GetBool("simple")
		var name string
		if len(args) > 0 {
			name = args[0]
		}
		if err := handleList(name, simple); err != nil {
			fmt.Fprintf(os.Stderr, "List failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().Bool("simple", false, "Simplify output to a table")
}

func handleList(name string, simple bool) error {
	entries, err := config.GetEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return nil
	}

	if name != "" {
		if entry, exists := entries[name]; exists {
			if simple {
				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				fmt.Fprintln(w, "Name\tAvailable")
				fmt.Fprintln(w, "----\t---------")
				fmt.Fprintf(w, "%s\t%v\n", name, entry.Available)
				w.Flush()
			} else {
				fmt.Printf("Name: %s\n", name)
				fmt.Printf("Available: %v\n", entry.Available)
				fmt.Printf("Target: %s\n", entry.Source)
			}
		} else {
			fmt.Printf("Entry '%s' not found.\n", name)
		}
	} else {
		if simple {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "Name\tAvailable")
			fmt.Fprintln(w, "----\t---------")
			for alias, entry := range entries {
				fmt.Fprintf(w, "%s\t%v\n", alias, entry.Available)
			}
			w.Flush()
		} else {
			fmt.Println("---")
			for alias, entry := range entries {
				fmt.Printf("Name: %s\n", alias)
				fmt.Printf("Available: %v\n", entry.Available)
				fmt.Printf("Target: %s\n", entry.Source)
				fmt.Println("---")
			}
		}
	}
	return nil
}
