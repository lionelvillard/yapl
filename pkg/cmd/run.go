// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
	"gopkg.in/yaml.v2"
)

var filename string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the specified YAPL file",
	Long:  `Run the specified YAPL file`,
	Run: func(cmd *cobra.Command, args []string) {
		resolve.AllowGlobalReassign = true
		var src interface{}
		f, err := syntax.Parse(filename, src, 0)
		if err != nil {
			fatalE(err)
		}

		if len(f.Stmts) == 0 {
			fatal("empty source file")
		}

		exprStmt, ok := f.Stmts[len(f.Stmts)-1].(*syntax.ExprStmt)
		if !ok {
			fatal("top-level statement must be an expression")
		}

		f.Stmts[len(f.Stmts)-1] = &syntax.AssignStmt{LHS: &syntax.Ident{Name: "$OUT$"}, RHS: exprStmt.X, Op: syntax.EQ}

		predeclared := make(starlark.StringDict)
		prog, err := starlark.FileProgram(f, predeclared.Has)
		if err != nil {
			fatalE(err)
		}

		thread := &starlark.Thread{}
		g, err := prog.Init(thread, predeclared)

		if err != nil {
			fatalE(err)
		}
		g.Freeze()

		value, ok := g["$OUT$"].(*starlark.Dict)
		if !ok {
			fatal("expected YAML output")
		}

		var j interface{}
		if err = json.Unmarshal([]byte(value.String()), &j); err != nil {
			fatalE(err)
		}

		d, err := yaml.Marshal(j)
		fmt.Printf("%s", string(d))
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&filename, "filename", "f", "", "name of the file containing the program to run")

}
