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
	"bytes"
)

// FileListDiff contains the paths of files that differ
// between two file lists
type FileListDiff struct {
	Added    []string
	Modified []string
	Deleted  []string
}

// Diff returns a FileListDiff structure that contains the paths
// of files that were added, modified, or deleted between the old
// file list and the new
func Diff(old *FileList, new *FileList) *FileListDiff {
	var added []string
	var modified []string
	var deleted []string

	it := old.Files.Iterator()
	for it.Next() {
		relPath := it.Key().(string)
		_, inNew := new.Files.Get(relPath)
		if !inNew {
			deleted = append(deleted, relPath)
		}
	}

	it = new.Files.Iterator()
	for it.Next() {
		relPath := it.Key().(string)
		oldMetadataInterface, inOld := old.Files.Get(relPath)
		if !inOld {
			added = append(added, relPath)
			continue
		}

		oldMetadata := oldMetadataInterface.(*FileMetadata)
		newMetadata := it.Value().(*FileMetadata)

		if bytes.Compare(oldMetadata.Hash, newMetadata.Hash) != 0 {
			modified = append(modified, relPath)
		} else if oldMetadata.Mode != newMetadata.Mode {
			modified = append(modified, relPath)
		}
	}

	diff := FileListDiff{
		Added:    added,
		Modified: modified,
		Deleted:  deleted,
	}

	return &diff
}
