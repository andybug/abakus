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


package repo

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

// HOME_DIR is the name of the abakus repo directory from root
const HOME_DIR string = ".abakus"

// OBJECTS_DIR is the name of the objects directory inside HOME_DIR
const OBJECTS_DIR string = "objects"

// SNAPSHOTS_DB is the name of the local database in the home dir
const SNAPSHOTS_DB string = "snapshots.db"


// GetHomeDir returns the path to the home directory with root as the base
func GetHomeDir(root string) (home string) {
	home = filepath.Join(root, HOME_DIR)
	return
}

// GetObjectsDir returns the path to the objects directory with root as the base
func GetObjectsDir(root string) (objects string) {
	objects = filepath.Join(root, HOME_DIR, OBJECTS_DIR)
	return
}

// GetSnapshotsDbPath returns the path to the local snapshot db with root as the base
func GetSnapshotsDbPath(root string) (snapshots_db string) {
	snapshots_db = filepath.Join(root, HOME_DIR, SNAPSHOTS_DB)
	return
}

// Create makes a new abakus repo in root/HOME_DIR
// Returns the path to home directory and an error (or nil)
func Create(root string) (string, error) {
	home := GetHomeDir(root)
	objects := GetObjectsDir(root)

	// check if the repo already exists
	if _, err := os.Stat(home); err == nil {
		errorStr := "Abakus repo already exists"
		return home, errors.New(errorStr)
	}

	// create home directory
	if err := createDir(home); err != nil {
		return home, err
	}

	// create objects directory
	if err := createDir(objects); err != nil {
		return home, err
	}

	// create empty snapshots db
	if err := createSnapshotsDb(root); err != nil {
		return home, err
	}

	return home, nil
}

// createDir creates a directory at the specified path if it
// does not exist.
func createDir(path string) error {
	// make sure the path doesn't already exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if it doesn't, create the path
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// createSnapshotsDb saves a new boltdb to HOME_DIR/SNAPSHOTS_DB
func createSnapshotsDb(root string) error {
	dbpath := GetSnapshotsDbPath(root)

	// open a new database
	db, err := bolt.Open(dbpath, 0664, nil)
	if err != nil {
		return err
	}

	// ...and close it. this causes the empty db to be written to disk
	db.Close()

	return nil
}
