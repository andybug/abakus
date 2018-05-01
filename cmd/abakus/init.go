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
	//"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new abakus repository in the current directory",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()
		home := filepath.Join(cwd, ".abakus")
		objects := filepath.Join(home, "objects")

		if _, err := os.Stat(home); err == nil {
			log.WithFields(log.Fields{
				"home": home,
			}).Fatal("Abakus repository already exists")
		}

		log.WithFields(log.Fields{
			"home": home,
		}).Info("Initializing abakus repository")

		createDir(home)
		createDir(objects)

		createSnapshotDb(home)
	},
}

func createDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.WithFields(log.Fields{"dir": path}).Debug("Creating directory")
		err := os.Mkdir(path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createSnapshotDb(home string) {
	dbpath := filepath.Join(home, "snapshots.db")

	db, err := bolt.Open(dbpath, 0664, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"database": dbpath,
		}).Fatal("Failed to create snapshots db")
	}
	db.Close()

	log.WithFields(log.Fields{
		"database": dbpath,
	}).Debug("Created snapshot db")
}
