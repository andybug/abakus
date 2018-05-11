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
	"errors"
	"fmt"

	"github.com/andybug/abakus/pkg/filelist"
	"github.com/andybug/abakus/pkg/repo"
)

// backend is an interface that snapshot storage mechanisms must implement
type backend interface {
	readMetadata(map[uint64]*SnapshotMetadata) (uint64, error)
	createSnapshot(*filelist.FileList, uint64) (*SnapshotMetadata, error)
	getSnapshotFiles(uint64) (*filelist.FileList, error)
	close()
}

// Store maintains a mapping of snapshot metadata for all snapshots, the
// database backend, and the latest snapshot
type Store struct {
	root     string
	backend  backend
	metadata map[uint64]*SnapshotMetadata
	latest   uint64
}

// GetStore returns a new snapshot Store
func GetStore(root string) (*Store, error) {
	dbPath := repo.GetSnapshotsDbPath(root)
	backend, err := newBoltBackend(dbPath)
	if err != nil {
		return nil, err
	}

	store := Store{
		root:     root,
		backend:  backend,
		metadata: make(map[uint64]*SnapshotMetadata),
		latest:   0,
	}

	latest, err := store.backend.readMetadata(store.metadata)
	if err != nil {
		return nil, err
	}

	store.latest = latest
	return &store, nil
}

// CreateSnapshot asks the backend to write a new snapshot with the given
// file list and id (latest + 1). The metadata for the new snapshot is added
// to the internal mapping and returned.
func (store *Store) CreateSnapshot(fl *filelist.FileList) (*SnapshotMetadata, error) {
	id := store.latest + 1

	snapshotMetadata, err := store.backend.createSnapshot(fl, id)
	if err != nil {
		return nil, err
	}

	store.metadata[id] = snapshotMetadata
	store.latest = id
	return snapshotMetadata, nil
}

// GetSnapshot returns a Snapshot that combines the metadata and data (file list).
// The backend is asked for the file list and it is combined with the metadata from
// the internal mapping.
func (store *Store) GetSnapshot(id uint64) (*Snapshot, error) {
	fl, err := store.backend.getSnapshotFiles(id)
	if err != nil {
		return nil, err
	}

	metadata := store.metadata[id]
	if metadata == nil {
		return nil, errors.New(fmt.Sprintf("No snapshot metadata with id %d", id))
	}

	snapshot := &Snapshot{
		Metadata: metadata,
		Files:    fl,
	}

	return snapshot, nil
}

// GetLatestSnapshot returns the Snapshot object associated with the latest
// snapshot
func (store *Store) GetLatestSnapshot() (*Snapshot, error) {
	snapshot, err := store.GetSnapshot(store.latest)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

// GetLatestId returns the id of the latest snapshot
func (store *Store) GetLatestId() uint64 {
	return store.latest
}

// GetAllMetadata returns a list containing the metadata of all snapshots
func (store *Store) GetAllMetadata() []*SnapshotMetadata {
	list := make([]*SnapshotMetadata, 0, len(store.metadata))
	for _, metadata := range store.metadata {
		list = append(list, metadata)
	}

	//FIXME: sort the list
	return list
}

// Close closes the backend
func (store *Store) Close() {
	store.backend.close()
}
