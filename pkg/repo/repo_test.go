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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHomeDir(t *testing.T) {
	root := "/a/b/c/d/"
	expected := "/a/b/c/d/.abakus"

	home := GetHomeDir(root)
	assert.Equal(t, expected, home)
}

func TestGetObjectsDir(t *testing.T) {
	root := "a/b"
	expected := "a/b/.abakus/objects"

	objects := GetObjectsDir(root)
	assert.Equal(t, expected, objects)
}

func TestGetSnapshotsDbPath(t *testing.T) {
	root := "/a/b/c"
	expected := "/a/b/c/.abakus/snapshots.db"

	snapshotDb := GetSnapshotsDbPath(root)
	assert.Equal(t, expected, snapshotDb)
}

func TestFindRoot(t *testing.T) {
	dir1, _ := ioutil.TempDir("", "TestFindRootSuccess")
	defer os.RemoveAll(dir1)

	os.Mkdir(filepath.Join(dir1, HOME_DIR), 0755)
	os.Mkdir(filepath.Join(dir1, "a"), 0755)
	os.Mkdir(filepath.Join(dir1, "a", "b"), 0755)
	os.Mkdir(filepath.Join(dir1, "a", "b", "c"), 0755)

	root, err := FindRoot(filepath.Join(dir1, "a", "b", "c"))
	assert.Equal(t, dir1, root)
	assert.Nil(t, err)

	dir2, _ := ioutil.TempDir("", "TestFindRootFailure")
	defer os.RemoveAll(dir2)

	root2, err2 := FindRoot(dir2)
	assert.Equal(t, "", root2)
	assert.NotNil(t, err2)
}

func TestCreateSuccess(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestCreateSuccess")
	defer os.RemoveAll(dir)

	fmt.Printf("TestCreateSuccess: %s\n", dir)

	home, err := Create(dir)
	assert.Nil(t, err)

	_, err = os.Stat(home)
	assert.Nil(t, err)

	objects := GetObjectsDir(dir)
	_, err = os.Stat(objects)
	assert.Nil(t, err)

	snapshots_db := GetSnapshotsDbPath(dir)
	_, err = os.Stat(snapshots_db)
	assert.Nil(t, err)
}

func TestCreateMissingDir(t *testing.T) {
	_, err := Create("/hopefully/not/a/real/dir")
	assert.NotNil(t, err)
}

func TestCreateTwice(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestCreateTwice")
	defer os.RemoveAll(dir)

	fmt.Printf("TestCreateTwice: %s\n", dir)

	_, err := Create(dir)
	assert.Nil(t, err)

	_, err = Create(dir)
	assert.NotNil(t, err)
}
