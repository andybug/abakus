// Copyright © 2018 Andrew Fields <andy@andybug.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"time"

	"github.com/andybug/abakus/pkg/filelist"
	"github.com/andybug/abakus/pkg/snapshot"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show changes to the workind directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		root := getRoot()

		store, err := snapshot.GetStore(root)
		exitError(err)
		defer store.Close()

		var latest_fl *filelist.FileList = nil

		// check if there are any snapshots
		if store.GetLatestId() == 0 {
			// if not, make latest be an empty file list
			latest_fl = filelist.New()
			fmt.Println("No Snapshots")
		} else {
			// if so, get the file list from the latest snapshot
			latest, err := store.GetLatestSnapshot()
			exitError(err)
			latest_fl = latest.Files
			fmt.Printf("Latest snapshot %d (%s)\n",
				latest.Metadata.Id,
				time.Unix(latest.Metadata.Timestamp, 0).String())
		}

		// get the file list for the working dir
		workdir, err := filelist.NewFromRoot(root)
		exitError(err)

		diff := filelist.Diff(latest_fl, workdir)

		// check if there are any changes
		if len(diff.Added) == 0 &&
			len(diff.Modified) == 0 &&
			len(diff.Deleted) == 0 {
			fmt.Println("No changes.")
			return
		}

		// output differences
		c := color.New(color.FgGreen)
		for _, added := range diff.Added {
			c.Printf("added:       %s\n", added)
		}

		c = color.New(color.FgRed)
		for _, modified := range diff.Modified {
			c.Printf("modified:    %s\n", modified)
		}

		c = color.New(color.FgRed)
		for _, deleted := range diff.Deleted {
			c.Printf("deleted:     %s\n", deleted)
		}
	},
}
