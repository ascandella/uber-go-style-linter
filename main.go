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
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	run := runLinter(getDir())
	os.Exit(run.summarize(os.Stderr))
}

func runLinter(dir string) *uberLinter {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, dir, nil, parser.ParseComments)
	if err != nil {
		log.Fatal("Unable to parse dir:", err)
	}

	linter := &uberLinter{
		fs: fs,
	}

	for _, pkg := range pkgs {
		ast.Walk(linter, pkg)
	}

	return linter
}

func getDir() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	if dir, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		return dir
	}
}
