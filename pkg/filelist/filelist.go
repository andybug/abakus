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

package filelist

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/andybug/abakus/pkg/repo"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/golang/crypto/blake2b"
)

// FileList maps the relative paths of files (from the root) to a FileMetadata
// structure that describes that file
type FileList struct {
	Files *treemap.Map
}

// FileMetadata describes a file in a FileList
// Hash - binary digest (blake2b)
// Size - size in bytes
// Mode - octal unix mode
// ModTime - unix time (seconds since epoch)
type FileMetadata struct {
	Hash    []byte `json:"hash"`
	Size    uint64 `json:"size"`
	Mode    uint32 `json:"mode"`
	ModTime uint64 `json:"-"`
}

// New creates an empty FileList
func New() *FileList {
	return &FileList{
		Files: treemap.NewWithStringComparator(),
	}
}

// NewFromRoot creates a FileList that includes all of the non-explicitly ignored
// files under the root of the repository
func NewFromRoot(root string) (*FileList, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	// create an exclusion rule for the home dir
	ignoreHome := newExcludeRules(root)
	ignoreHome.add(fmt.Sprintf("/%s", repo.HOME_DIR))
	esr := newExcludeRulesStack()
	esr.push(ignoreHome)

	fl := New()
	fl.addTree(root, root, esr)

	return fl, nil
}

// Add adds file at relative path to the file list with the given metadata
// the filelist maps path -> metadata
func (fl *FileList) Add(relPath string, metadata *FileMetadata) {
	fl.Files.Put(relPath, metadata)
}

// addTree adds all of the files under that point to the FileList
// root and dir must be absolute paths, and dir must be under root
// addTree will use the stack to keep track of what exclusions apply
// to different directories as it walks the file system
func (fl *FileList) addTree(root string, dir string, stack *excludeRulesStack) error {
	rules, err := readRules(dir)
	if err != nil {
		return err
	}
	stack.push(rules)
	defer stack.pop()

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		absFilePath := filepath.Join(dir, file.Name())
		relFilePath, _ := filepath.Rel(root, absFilePath)

		// check if this file matches an exclusion rule
		// ignore it if it does
		if stack.exclude(absFilePath) {
			continue
		}

		if file.IsDir() {
			err = fl.addTree(root, absFilePath, stack)
			if err != nil {
				return err
			}
		} else {
			hash, err := hashFile(absFilePath)
			if err != nil {
				return err
			}

			metadata := FileMetadata{
				Hash:    hash,
				Size:    uint64(file.Size()),
				Mode:    uint32(file.Mode()),
				ModTime: uint64(file.ModTime().Unix()),
			}

			fl.Add(relFilePath, &metadata)
		}
	}

	return nil
}

// MerkleRoot calculates the blake2b root hash of a tree
// built from the filelist (like bitcoin). The MerkleRoot
// function hashes each file path/content hash and adds them
// to an array. This array represents the leaves in the
// merkle tree. The array is passed to the merkleTree function
// to calculate the merkle hash of the subtree.
func (fl *FileList) MerkleRoot() []byte {
	hasher, _ := blake2b.New256(nil)
	var hashes [][]byte

	it := fl.Files.Iterator()
	for it.Next() {
		path := it.Key().(string)
		metadata := it.Value().(*FileMetadata)

		hasher.Write([]byte(path))
		hasher.Write(metadata.Hash)
		sum := hasher.Sum(nil)
		hasher.Reset()

		hashes = append(hashes, sum)
	}

	return merkleTree(hashes)
}

// merkleTree recursively calculates the merkle hash of a subtree
func merkleTree(hashes [][]byte) []byte {
	if len(hashes) == 1 {
		return hashes[0]
	}
	if len(hashes)%2 != 0 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	hasher, _ := blake2b.New256(nil)
	var newHashes [][]byte

	for i := 0; i < len(hashes); i += 2 {
		hasher.Write(hashes[i])
		hasher.Write(hashes[i+1])
		sum := hasher.Sum(nil)
		newHashes = append(newHashes, sum)
	}

	return merkleTree(newHashes)
}

// hashFile returns the blake2b hash of a file on disk
func hashFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b2b, _ := blake2b.New256(nil)
	if _, err = io.Copy(b2b, f); err != nil {
		return nil, err
	}

	out := b2b.Sum(nil)
	return out, nil
}
