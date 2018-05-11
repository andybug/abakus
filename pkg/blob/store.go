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

package blob

import (
	"bufio"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/andybug/abakus/pkg/filelist"
	"github.com/andybug/abakus/pkg/repo"
	"github.com/peterbourgon/diskv"
)

// Store wraps the diskv handle
type Store struct {
	root     string
	blobsDir string
	handle   *diskv.Diskv
}

// GetStore returns a blob store object
func GetStore(root string) (*Store, error) {
	blobsDir := repo.GetBlobsDir(root)
	handle := diskv.New(diskv.Options{
		BasePath:     blobsDir,
		Transform:    func(s string) []string { return []string{} },
		CacheSizeMax: 1024 * 1024,
		Compression:  diskv.NewZlibCompressionLevel(6),
	})

	store := Store{
		root:     root,
		blobsDir: blobsDir,
		handle:   handle,
	}

	return &store, nil
}

// AddFiles will check each file in the file list to ensure that it is
// in the blob store; if not, it will be added. returns the number of
// files added to the store and the number that were already present
func (store *Store) AddFiles(fl *filelist.FileList) (uint64, uint64, error) {
	var newFiles uint64 = 0
	var existingFiles uint64 = 0

	it := fl.Files.Iterator()
	for it.Next() {
		metadata := it.Value().(*filelist.FileMetadata)
		hashString := hex.EncodeToString(metadata.Hash)
		if store.handle.Has(hashString) {
			existingFiles += 1
			continue
		}

		absPath := filepath.Join(store.root, it.Key().(string))
		stream, err := os.Open(absPath)
		if err != nil {
			return newFiles, existingFiles, err
		}

		reader := bufio.NewReader(stream)
		err = store.handle.WriteStream(hashString, reader, true)
		if err != nil {
			return newFiles, existingFiles, err
		}
		newFiles += 1
		stream.Close()
	}

	return newFiles, existingFiles, nil
}
