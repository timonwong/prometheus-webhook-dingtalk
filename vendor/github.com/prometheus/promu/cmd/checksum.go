// Copyright © 2017 Prometheus Team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	checksumsFilename = "sha256sums.txt"
)

var (
	checksumcmd      = app.Command("checksum", "Calculate the SHA256 checksum for each file in the given location")
	checksumLocation = checksumcmd.Arg("location", "Location to checksum").Default(".").Strings()
)

func runChecksum(path string) {
	checksums, err := calculateSHA256s(path)
	if err != nil {
		fatal(errors.Wrap(err, "Failed to calculate checksums"))
	}

	file, err := os.Create(filepath.Join(path, checksumsFilename))
	if err != nil {
		fatal(errors.Wrap(err, "Failed to create checksums file"))
	}
	defer file.Close()
	for _, c := range checksums {
		if _, err := fmt.Fprintf(file, "%x  %s\n", c.checksum, c.filename); err != nil {
			fatal(errors.Wrap(err, "Failed to write to checksums file"))
		}
	}
}

type checksumSHA256 struct {
	filename string
	checksum []byte
}

// calculateSHA256s calculates the sha256 checksum for each file in the given
// path and returns a checksumSHA256 type in the order returned of
// filepath.Walk.
func calculateSHA256s(path string) ([]checksumSHA256, error) {
	var checksums []checksumSHA256
	path = fmt.Sprintf("%s%c", filepath.Clean(path), filepath.Separator)
	calculateSHA256 := func(filepath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		file, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer file.Close()

		hash := sha256.New()
		if _, err = io.Copy(hash, file); err != nil {
			return err
		}
		checksums = append(checksums, checksumSHA256{
			filename: strings.TrimPrefix(filepath, path),
			checksum: hash.Sum(nil),
		})

		return nil
	}
	if err := filepath.Walk(path, calculateSHA256); err != nil {
		return nil, err
	}
	return checksums, nil
}
