// Copyright 2021 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/open-policy-agent/opa/cmd/internal/env"
	fileurl "github.com/open-policy-agent/opa/internal/file/url"
	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/bundle"
	"github.com/open-policy-agent/opa/v1/format"
	"github.com/open-policy-agent/opa/v1/loader"
	"github.com/open-policy-agent/opa/v1/refactor"
)

type moveCommandParams struct {
	mapping      repeatedStringFlag
	ignore       []string
	overwrite    bool
	v0Compatible bool
	v1Compatible bool
}

func (m *moveCommandParams) regoVersion() ast.RegoVersion {
	// v0 takes precedence over v1
	if m.v0Compatible {
		return ast.RegoV0
	}
	if m.v1Compatible {
		return ast.RegoV1
	}
	return ast.DefaultRegoVersion
}

func initRefactor(root *cobra.Command, brand string) {
	executable := root.Name()

	var moveCommandParams moveCommandParams

	var refactorCommand = &cobra.Command{
		Use:    "refactor",
		Short:  "Refactor Rego file(s)",
		Hidden: true,
	}

	var moveCommand = &cobra.Command{
		Use:   "move [file-path [...]]",
		Short: "Rename packages and their references in Rego file(s)",
		Long: `Rename packages and their references in Rego file(s).

The 'move' command takes one or more Rego source file(s) and rewrites package paths and other references in them as per
the mapping defined by the '-p' option. At least one mapping should be provided and should be of the form:

	<from>:<to>

The 'move' command formats the Rego modules after renaming packages, etc. and prints the formatted modules to stdout by default.
If the '-w' option is supplied, the 'move' command will overwrite the source file instead.

Example:
--------

"policy.rego" contains the below policy:
 _ _ _ _ _ _ _ _ _ _ _ _ _
| package lib.foo         |
|                         |
| default allow = false   |
| _ _ _ _ _ _ _ _ _ _ _ _ |     
	
	$ ` + executable + ` refactor move -p data.lib.foo:data.baz.bar policy.rego

The 'move' command outputs the below policy to stdout with the package name rewritten as per the mapping:

 _ _ _ _ _ _ _ _ _ _ _ _ _
| package baz.bar         |
|                         |
| default allow = false   |
| _ _ _ _ _ _ _ _ _ _ _ _ | 
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := validateMoveArgs(args); err != nil {
				return err
			}
			return env.CmdFlags.CheckEnvironmentVariables(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true

			if err := doMove(moveCommandParams, args, os.Stdout); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				return err
			}
			return nil
		},
	}

	moveCommand.Flags().VarP(&moveCommandParams.mapping, "path", "p", "set the mapping that defines how references should be rewritten (ie. <from>:<to>). This flag can be repeated.")
	moveCommand.Flags().BoolVarP(&moveCommandParams.overwrite, "write", "w", false, "overwrite the original source file")
	addIgnoreFlag(moveCommand.Flags(), &moveCommandParams.ignore)
	refactorCommand.AddCommand(moveCommand)
	addV0CompatibleFlag(moveCommand.Flags(), &moveCommandParams.v0Compatible, false)
	addV1CompatibleFlag(moveCommand.Flags(), &moveCommandParams.v1Compatible, false)
	root.AddCommand(refactorCommand)
}

func doMove(params moveCommandParams, args []string, out io.Writer) error {
	if len(params.mapping.v) == 0 {
		return errors.New("specify at least one mapping of the form <from>:<to>")
	}

	srcDstMap, err := parseSrcDstMap(params.mapping.v)
	if err != nil {
		return err
	}

	result, err := loader.NewFileLoader().
		WithBundleLazyLoadingMode(bundle.HasExtension()).
		WithRegoVersion(params.regoVersion()).
		Filtered(args, ignored(params.ignore).Apply)
	if err != nil {
		return err
	}

	mq := refactor.MoveQuery{
		Modules:       result.ParsedModules(),
		SrcDstMapping: srcDstMap,
	}.WithValidation(true)

	movedModules, err := refactor.New().Move(mq)
	if err != nil {
		return err
	}

	for filename, mod := range movedModules.Result {
		filename, err = fileurl.Clean(filename)
		if err != nil {
			return err
		}

		formatted, err := format.AstWithOpts(mod, format.Opts{RegoVersion: params.regoVersion()})
		if err != nil {
			return newError("failed to parse Rego source file: %v", err)
		}

		if params.overwrite {
			info, err := os.Stat(filename)
			if err != nil {
				return err
			}

			outfile, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, info.Mode())
			if err != nil {
				return newError("failed to open file for writing: %v", err)
			}
			defer outfile.Close()
			out = outfile
		}

		_, err = out.Write(formatted)
		if err != nil {
			return newError("failed writing formatted contents: %v", err)
		}
	}

	return nil
}

func parseSrcDstMap(data []string) (map[string]string, error) {
	result := map[string]string{}

	for _, d := range data {
		term, err := ast.ParseTerm("{" + d + "}")
		if err != nil {
			return nil, newError("failed to parse mapping: %v", err)
		}
		obj, ok := term.Value.(ast.Object)
		if !ok {
			return nil, newError("expected mapping of the form <from>:<to>")
		}
		keys := obj.Keys()
		if len(keys) != 1 {
			return nil, newError("expected mapping of the form <from>:<to>")
		}
		result[keys[0].String()] = obj.Get(keys[0]).String()
	}
	return result, nil
}

func validateMoveArgs(args []string) error {
	if len(args) == 0 {
		return errors.New("specify at least one path containing policy files")
	}
	return nil
}
