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
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"log"
)

type uberLinter struct {
	fs   *token.FileSet
	errs []error
}

func (u *uberLinter) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	// These two cases are necessary to recurse into actual code
	case *ast.File:
		return u
	case *ast.Package:
		return u
	// This is the stuff we care about
	case *ast.GenDecl:
		if n.Tok == token.CONST {
			for _, spec := range n.Specs {
				if t, ok := spec.(*ast.ValueSpec); !ok {
					log.Fatal("Unknown type: ", spec)
				} else {
					for i, name := range t.Names {
						u.lintConst(name, t.Values[i])
					}
				}
			}
		}
	}
	return nil
}

func (u *uberLinter) lintConst(id *ast.Ident, val ast.Expr) {
	name := id.Name
	if len(name) == 0 {
		u.addError("zero length constant name?", id.Pos())
	} else if string(name[0]) != "_" {
		u.addError(fmt.Sprintf(`const "%s" should be "_%s" to avoid shadowing`, name, name), id.Pos())
	}
}

func (u *uberLinter) addError(msg string, pos token.Pos) {
	p := u.fs.Position(pos)
	u.errs = append(u.errs, fmt.Errorf("%s:%d %s", p.Filename, p.Line, msg))
}

func (u *uberLinter) summarize(out io.Writer) int {
	if len(u.errs) == 0 { // sall good man
		return 0
	}
	for _, err := range u.errs {
		fmt.Fprintf(out, err.Error()+"\n")
	}
	return 1
}
