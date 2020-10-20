// Copyright © 2020 The Tekton Authors.
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

package get

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/tektoncd/hub/api/pkg/cli/app"
	"github.com/tektoncd/hub/api/pkg/cli/flag"
	"github.com/tektoncd/hub/api/pkg/cli/hub"
	"github.com/tektoncd/hub/api/pkg/cli/printer"
)

type options struct {
	cli     app.CLI
	catalog string
	version string
	kind    string
	args    []string
}

func Command(cli app.CLI) *cobra.Command {

	opts := &options{cli: cli}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get resource manifest by its name, kind, catalog, and version",
		Long:  ``,
		Annotations: map[string]string{
			"commandType": "main",
		},
		SilenceUsage: true,
	}

	cmd.AddCommand(
		commandForKind("task", opts),
		commandForKind("pipeline", opts),
	)

	cmd.PersistentFlags().StringVar(&opts.catalog, "catalog", "tekton", "Name of Catalog to which resource belongs to.")
	cmd.PersistentFlags().StringVar(&opts.version, "version", "", "Version of Resource")

	return cmd
}

// commandForKind creates a cobra.Command that when run sets
// opts.Kind and opts.Args and invokes opts.run
func commandForKind(kind string, opts *options) *cobra.Command {

	return &cobra.Command{
		Use:          kind,
		Short:        "Get " + kind + " by name, catalog and version",
		Long:         ``,
		SilenceUsage: true,
		Annotations: map[string]string{
			"commandType": "main",
		},
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.kind = kind
			opts.args = args
			return opts.run()
		},
	}
}

func (opts *options) run() error {

	if err := opts.validate(); err != nil {
		return err
	}

	hubClient := opts.cli.Hub()

	resource := hubClient.GetResource(hub.ResourceOption{
		Name:    opts.name(),
		Catalog: opts.catalog,
		Kind:    opts.kind,
		Version: opts.version,
	})

	out := opts.cli.Stream().Out
	return printer.New(out).Raw(resource.Manifest())
}

func (opts *options) validate() error {
	return flag.ValidateVersion(opts.version)
}

func (opts *options) name() string {
	return strings.TrimSpace(opts.args[0])
}
