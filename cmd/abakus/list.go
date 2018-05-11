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
	"os"
	"text/tabwriter"
	"time"

	"github.com/andybug/abakus/pkg/snapshot"
	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List snapshots in the repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		root := getRoot()

		store, err := snapshot.GetStore(root)
		exitError(err)
		defer store.Close()

		w := tabwriter.NewWriter(os.Stdout, 4, 0, 4, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "ID\tTIME\tMERKLE\tFILES\tSIZE")

		metadataList := store.GetAllMetadata()
		for _, metadata := range metadataList {
			fmt.Fprintf(w, "%d\t%s\t%x\t%s\t%s\n",
				metadata.Id,
				humanize.Time(time.Unix(metadata.Timestamp, 0)),
				metadata.MerkleRoot[:4],
				humanize.Comma(int64(metadata.FileCount)),
				humanize.Bytes(metadata.Size))
		}
		w.Flush()
	},
}
