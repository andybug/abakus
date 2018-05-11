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
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/andybug/abakus/pkg/filelist"
	"github.com/boltdb/bolt"
)

// BOLT_METADATA_KEY is the key in each snapshot bucket that contains
// a SnapshotMetadata object in json
const BOLT_METADATA_KEY = "__abakus.metadata"

// bolt_backend wraps the bolt db handle
type bolt_backend struct {
	dbPath string
	db     *bolt.DB
}

// newBoltBackend opens the snapshot db and returns the bolt_backend
func newBoltBackend(dbPath string) (backend, error) {
	db, err := bolt.Open(dbPath, 0644, nil)
	if err != nil {
		return nil, err
	}

	b := &bolt_backend{
		dbPath: dbPath,
		db:     db,
	}

	return b, nil
}

// readMetadata fills out the metadataMap with the metadata from all snapshots
// in the store. The mapping in id -> SnapshotMetadata. It returns the latest
// (i.e. the most recent) snapshot.
func (b bolt_backend) readMetadata(metadataMap map[uint64]*SnapshotMetadata) (uint64, error) {
	var latest uint64 = 0
	re := regexp.MustCompile(`snapshot:(?P<id>\d+)`)

	err := b.db.View(func(tx *bolt.Tx) error {
		// iterate over every bucket
		return tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			// match bucket name to expected format
			groups := re.FindStringSubmatch(string(name))
			if len(groups) != 2 {
				return errors.New(fmt.Sprintf("invalid bucket name '%s'", string(name)))
			}

			// get snapshot id from bucket name
			id, err := strconv.ParseUint(groups[1], 10, 64)
			if err != nil {
				return err
			}

			if id > latest {
				latest = id
			}

			metadataJson := bucket.Get([]byte(BOLT_METADATA_KEY))
			metadata := new(SnapshotMetadata)

			err = json.Unmarshal(metadataJson, metadata)
			if err != nil {
				return err
			}

			metadata.Id = id
			metadataMap[id] = metadata

			return nil
		})
	})

	if err != nil {
		return 0, err
	}

	return latest, nil
}

// createSnapshot takes a file list and id and writes a new snapshot to the
// database. each snapshot is in its own bucket. it returns the metadata for
// the created snapshot
func (b bolt_backend) createSnapshot(fl *filelist.FileList, id uint64) (*SnapshotMetadata, error) {
	timestamp := time.Now().Unix()
	merkle := fl.MerkleRoot()
	var size uint64 = 0
	var fileCount uint64 = 0
	var snapshotMetadata *SnapshotMetadata = nil

	err := b.db.Update(func(tx *bolt.Tx) error {
		bucketName := fmt.Sprintf("snapshot:%d", id)
		bucket, err := tx.CreateBucket([]byte(bucketName))
		if err != nil {
			return err
		}

		it := fl.Files.Iterator()
		for it.Next() {
			metadata := it.Value().(*filelist.FileMetadata)
			jsonMetadata, err := json.Marshal(metadata)
			if err != nil {
				return err
			}

			err = bucket.Put([]byte(it.Key().(string)), []byte(jsonMetadata))
			if err != nil {
				return err
			}

			fileCount += 1
			size += metadata.Size
		}

		snapshotMetadata = &SnapshotMetadata{
			Id:         id,
			Timestamp:  timestamp,
			MerkleRoot: merkle,
			FileCount:  fileCount,
			Size:       size,
		}

		jsonSnapshotMetadata, err := json.Marshal(snapshotMetadata)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(BOLT_METADATA_KEY), jsonSnapshotMetadata)

		return err
	})

	if err != nil {
		return nil, err
	}

	return snapshotMetadata, nil

}

// getSnapshotFiles returns the file list associated with the snapshot id. the
// metadata is not retrieved because the snapshot store maintains a list of
// all of the metadata
func (b bolt_backend) getSnapshotFiles(id uint64) (*filelist.FileList, error) {
	bucketName := fmt.Sprintf("snapshot:%d", id)
	var fl *filelist.FileList

	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("No snapshot with id %d", id))
		}

		var err error
		fl, err = bolt_readFileList(bucket)

		return err
	})

	if err != nil {
		return nil, err
	}

	return fl, nil
}

// bolt_readFileList iterates over the keys in a bucket, ignore the metadata
// key, and builds a file list from the rest
func bolt_readFileList(bucket *bolt.Bucket) (*filelist.FileList, error) {
	var fl = filelist.New()
	err := bucket.ForEach(func(key []byte, value []byte) error {
		path := string(key)
		if path == BOLT_METADATA_KEY {
			return nil
		}

		metadata := new(filelist.FileMetadata)

		err := json.Unmarshal(value, metadata)
		if err != nil {
			return err
		}

		fl.Add(path, metadata)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fl, nil
}

// close closes the bolt database
func (b bolt_backend) close() {
	b.db.Close()
}
