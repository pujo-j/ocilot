/*
 *    Copyright 2020 Josselin Pujo
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 *
 */

package ocilot

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/dustin/go-humanize"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Snapshot struct {
	Name     string
	Creation time.Time
	Files    map[string]*tar.Header
}

func (s *Snapshot) String() string {
	return s.Name + "@" + s.Creation.Format("15:04:05.000")
}

type countingWriter struct {
	Size int64
}

func (c *countingWriter) Write(p []byte) (n int, err error) {
	l := len(p)
	c.Size += int64(l)
	return l, nil
}

func (s *Snapshot) AsLayer(workFile string) (v1.Layer, error) {
	file, err := os.Create(workFile)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		return nil, err
	}
	hasher := sha256.New()
	multiWriter := io.MultiWriter(file, hasher)
	writer := tar.NewWriter(multiWriter)
	for path, header := range s.Files {
		err := writer.WriteHeader(header)
		if err != nil {
			return nil, err
		}
		if header.Typeflag != tar.TypeDir {
			f, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(writer, f)
			_ = f.Close()
			if err != nil {
				fmt.Printf("Error copying %v , type: %s", header, string(header.Typeflag))
				return nil, err
			}
		}
	}
	err = writer.Close()
	diffId := v1.Hash{
		Algorithm: "sha256",
		Hex:       hex.EncodeToString(hasher.Sum(make([]byte, 0, hasher.Size()))),
	}
	compressedFile, err := os.Create(workFile + ".gz")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = compressedFile.Close()
	}()
	counter := &countingWriter{}
	compressedHasher := sha256.New()
	newWriter, err := gzip.NewWriterLevel(io.MultiWriter(compressedFile, counter, compressedHasher), 2)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(newWriter, file)
	if err != nil {
		return nil, err
	}
	err = newWriter.Close()
	if err != nil {
		return nil, err
	}
	digest := v1.Hash{
		Algorithm: "sha256",
		Hex:       hex.EncodeToString(compressedHasher.Sum(make([]byte, 0, compressedHasher.Size()))),
	}
	return &SnapshotLayer{
		tarFile:        workFile,
		compressedFile: workFile + ".gz",
		diffID:         diffId,
		digest:         digest,
		compressedSize: counter.Size,
		mediaType:      types.DockerLayer,
	}, nil
}

func NewSnapshot(dirPath string) (*Snapshot, error) {
	humanize.Time(time.Now())
	res := &Snapshot{
		Name:  dirPath,
		Files: make(map[string]*tar.Header),
	}
	err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		fullPath, err := filepath.Abs(filePath)
		if err != nil {
			return err
		}
		fullPath = strings.Replace(fullPath, "\\", "/", -1)
		symlink, err := filepath.EvalSymlinks(fullPath)
		if err != nil {
			return err
		}
		symlink = strings.Replace(symlink, "\\", "/", -1)
		th := tar.Header{}
		th.Name = fullPath
		th.ModTime = info.ModTime()
		th.Size = info.Size()
		if symlink != th.Name {
			th.Typeflag = tar.TypeSymlink
			th.Linkname = symlink
		} else {
			if info.IsDir() {
				th.Typeflag = tar.TypeDir
			} else {
				th.Typeflag = tar.TypeReg
			}
		}
		res.Files[filePath] = &th
		return nil
	})
	return res, err
}

func Diff(old *Snapshot, new *Snapshot) (*Snapshot, error) {
	res := &Snapshot{
		Files: make(map[string]*tar.Header),
	}
	for name, th := range new.Files {
		old, ok := old.Files[name]
		if ok {
			if !th.ModTime.Equal(old.ModTime) || th.Size != old.Size {
				res.Files[name] = th
			}
		} else {
			res.Files[name] = th
		}
	}
	return res, nil
}

type SnapshotLayer struct {
	tarFile        string
	compressedFile string
	diffID         v1.Hash
	digest         v1.Hash
	compressedSize int64
	mediaType      types.MediaType
}

func (s *SnapshotLayer) Digest() (v1.Hash, error) {
	return s.digest, nil
}

func (s *SnapshotLayer) DiffID() (v1.Hash, error) {
	return s.diffID, nil
}

func (s *SnapshotLayer) Compressed() (io.ReadCloser, error) {
	return os.Open(s.compressedFile)
}

func (s *SnapshotLayer) Uncompressed() (io.ReadCloser, error) {
	return os.Open(s.tarFile)
}

func (s *SnapshotLayer) Size() (int64, error) {
	return s.compressedSize, nil
}

func (s SnapshotLayer) MediaType() (types.MediaType, error) {
	return s.mediaType, nil
}
