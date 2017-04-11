// Copyright 2017 Aiden Scandella
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

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	_testdata      = "testdata"
	_expectedError = "expected.error"
)

func TestRealCases(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	require.NoError(t, err, "Unable to get working directory")

	filepath.Walk(filepath.Join(wd, _testdata), func(path string, info os.FileInfo, err error) error {
		require.NoError(t, err, "Unexpected error walking testdata")
		if strings.HasSuffix(path, _testdata) {
			// skip the TLD
			return nil
		}
		if info.IsDir() {
			t.Run(path, func(t *testing.T) {
				out := &bytes.Buffer{}
				result := runLinter(path)
				ret := result.summarize(out)
				maybeError := filepath.Join(path, _expectedError)

				if _, err := os.Stat(maybeError); err == nil && !os.IsNotExist(err) {
					expected, err := ioutil.ReadFile(maybeError)
					require.NoError(t, err, "Unable to read expected error file")
					outScrubbed := strings.Replace(out.String(), path, "", -1)
					assert.Equal(t, outScrubbed, string(expected))
					assert.NotEqual(t, 0, ret, "Expected non-zero exit code")
				} else {
					assert.Empty(t, out.String(), "Expected no lint erors")
					assert.Equal(t, 0, ret, "Expected zero exit code")
				}
			})
		}
		return nil
	})
}
