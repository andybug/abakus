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

package snapshot

import (
	"github.com/andybug/abakus/pkg/filelist"
)

// SnapshotMetadata contains all of the metadata about a snapshot
// It only lacks the file list. The snapshot store maintains a mapping
// of all of the metadata.
type SnapshotMetadata struct {
	Id         uint64 `json:"-"`
	Timestamp  int64  `json:"timestamp"`
	MerkleRoot []byte `json:"merkle"`
	FileCount  uint64 `json:"files"`
	Size       uint64 `json:"size"`
}

// Snapshot contains the metadata and data of a snapshot
type Snapshot struct {
	Metadata *SnapshotMetadata
	Files    *filelist.FileList
}
