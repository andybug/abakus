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
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/andybug/abakus/pkg/filelist"
	"github.com/andybug/abakus/pkg/snapshot"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show files in a snapshot",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		root := getRoot()

		snapshotStore, err := snapshot.GetStore(root)
		exitError(err)
		defer snapshotStore.Close()

		if len(args) != 1 {
			exitError(errors.New("status requires and id argument"))
		}

		id, err := strconv.ParseUint(args[0], 10, 64)
		exitError(err)

		snapshot, err := snapshotStore.GetSnapshot(id)
		exitError(err)

		w := tabwriter.NewWriter(os.Stdout, 4, 0, 4, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "PATH\tHASH\tSIZE\tMODE")

		it := snapshot.Files.Files.Iterator()
		for it.Next() {
			metadata := it.Value().(*filelist.FileMetadata)
			fmt.Fprintf(w, "%s\t%x\t%s\t%o\n",
				it.Key().(string),
				metadata.Hash,
				humanize.Bytes(metadata.Size),
				metadata.Mode,
			)
		}
		w.Flush()
	},
}
