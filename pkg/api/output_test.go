// Copyright 2024 daniel pickens
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Copyright 2024 daniel pickens
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
// limitations under the License

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testOutput1 = &Output{
	Name:      "some name one",
	Namespace: "lamb-namespace",
	FilePath:  "path-to-file",
	APIVersion: &Version{
		Name:                   "extensions/v1beta1",
		Kind:                   "Deployment",
		DeprecatedIn:           "v1.0.0",
		RemovedIn:              "v1.0.0",
		ReplacementAPI:         "apps/v1",
		ReplacementAvailableIn: "v1.10.0",
		Component:              "foo",
	},
}
var testOutput2 = &Output{
	Name: "some name two",
	APIVersion: &Version{
		Name:                   "extensions/v1beta1",
		Kind:                   "Deployment",
		DeprecatedIn:           "v1.0.0",
		RemovedIn:              "v1.0.0",
		ReplacementAPI:         "apps/v1",
		ReplacementAvailableIn: "v1.10.0",
		Component:              "foo",
	},
}

var testOutputNoOutput = &Output{
	Name: "not a deprecated object",
	APIVersion: &Version{
		Name:                   "apps/v1",
		Kind:                   "Deployment",
		DeprecatedIn:           "",
		RemovedIn:              "",
		ReplacementAPI:         "",
		ReplacementAvailableIn: "",
		Component:              "foo",
	},
}

var testOutputDeprecatedNotRemoved = &Output{
	Name: "deprecated not removed",
	APIVersion: &Version{
		Name:           "apps/v1",
		Kind:           "Deployment",
		DeprecatedIn:   "v1.0.0",
		RemovedIn:      "",
		ReplacementAPI: "none",
		Component:      "foo",
	},
	ReplacementAvailable: false,
}

func init() {
	padChar = byte('-')
}

func ExampleInstance_DisplayOutput_normal() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs: []*Output{
			testOutput1,
			testOutput2,
			testOutputDeprecatedNotRemoved,
		},
		OutputFormat: "normal",
		Components:   []string{"foo"},
	}
	_ = instance.DisplayOutput()

	// Output:
	// NAME-------------------- KIND-------- VERSION------------- REPLACEMENT-- REMOVED-- DEPRECATED-- REPL AVAIL--
	// some name one----------- Deployment-- extensions/v1beta1-- apps/v1------ true----- true-------- true--------
	// some name two----------- Deployment-- extensions/v1beta1-- apps/v1------ true----- true-------- true--------
	// deprecated not removed-- Deployment-- apps/v1------------- none--------- false---- true-------- false-------
}

func ExampleInstance_DisplayOutput_onlyShowRemoved() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		OnlyShowRemoved: true,
		Outputs: []*Output{
			testOutput1,
			testOutput2,
			testOutputDeprecatedNotRemoved,
		},
		OutputFormat: "normal",
		Components:   []string{"foo"},
	}
	_ = instance.DisplayOutput()

	// Output:
	// NAME----------- KIND-------- VERSION------------- REPLACEMENT-- REMOVED-- DEPRECATED-- REPL AVAIL--
	// some name one-- Deployment-- extensions/v1beta1-- apps/v1------ true----- true-------- true--------
	// some name two-- Deployment-- extensions/v1beta1-- apps/v1------ true----- true-------- true--------
}

func ExampleInstance_DisplayOutput_wide() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat: "wide",
		Components:   []string{"foo"},
	}
	_ = instance.DisplayOutput()

	// Output:
	// NAME----------- NAMESPACE-------- KIND-------- VERSION------------- REPLACEMENT-- DEPRECATED-- DEPRECATED IN-- REMOVED-- REMOVED IN-- REPL AVAIL-- REPL AVAIL IN--
	// some name one-- lamb-namespace-- Deployment-- extensions/v1beta1-- apps/v1------ true-------- v1.0.0--------- true----- v1.0.0----- true-------- v1.10.0--------
	// some name two-- <EXAMPLE>-------- Deployment-- extensions/v1beta1-- apps/v1------ true-------- v1.0.0--------- true----- v1.0.0----- true-------- v1.10.0--------
}

func ExampleInstance_DisplayOutput_custom() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat:  "custom",
		Components:    []string{"foo"},
		CustomColumns: []string{"NAMESPACE", "NAME", "DEPRECATED IN", "DEPRECATED", "REPLACEMENT", "VERSION", "KIND", "COMPONENT", "FILEPATH"},
	}
	_ = instance.DisplayOutput()

	// Output:
	// NAMESPACE-------- NAME----------- DEPRECATED IN-- DEPRECATED-- REPLACEMENT-- VERSION------------- KIND-------- COMPONENT-- FILEPATH------
	// lamb-namespace-- some name one-- v1.0.0--------- true-------- apps/v1------ extensions/v1beta1-- Deployment-- foo-------- path-to-file--
	// <EXAMPLE>-------- some name two-- v1.0.0--------- true-------- apps/v1------ extensions/v1beta1-- Deployment-- foo-------- <EXAMPLE>-----
}

func ExampleInstance_DisplayOutput_markdown() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat: "markdown",
		Components:   []string{"foo"},
	}
	_ = instance.DisplayOutput()

	// Output:
	// |     NAME      |    NAMESPACE    |    KIND    |      VERSION       | REPLACEMENT | DEPRECATED | DEPRECATED IN | REMOVED | REMOVED IN | REPL AVAIL | REPL AVAIL IN |
	// |---------------|-----------------|------------|--------------------|-------------|------------|---------------|---------|------------|------------|---------------|
	// | some name one | lamb-namespace | Deployment | extensions/v1beta1 | apps/v1     | true       | v1.0.0        | true    | v1.0.0    | true       | v1.10.0       |
	// | some name two | <EXAMPLE>       | Deployment | extensions/v1beta1 | apps/v1     | true       | v1.0.0        | true    | v1.0.0    | true       | v1.10.0       |
}

func ExampleInstance_DisplayOutput_markdown_customcolumns() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat:  "markdown",
		Components:    []string{"foo"},
		CustomColumns: []string{"NAMESPACE", "NAME", "DEPRECATED IN", "DEPRECATED", "REPLACEMENT", "VERSION", "KIND", "COMPONENT", "FILEPATH"},
	}
	_ = instance.DisplayOutput()

	// Output:
	// |    NAMESPACE    |     NAME      | DEPRECATED IN | DEPRECATED | REPLACEMENT |      VERSION       |    KIND    | COMPONENT |   FILEPATH   |
	// |-----------------|---------------|---------------|------------|-------------|--------------------|------------|-----------|--------------|
	// | lamb-namespace | some name one | v1.0.0        | true       | apps/v1     | extensions/v1beta1 | Deployment | foo       | path-to-file |
	// | <EXAMPLE>       | some name two | v1.0.0        | true       | apps/v1     | extensions/v1beta1 | Deployment | foo       | <EXAMPLE>    |
}

func ExampleInstance_DisplayOutput_json() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		OutputFormat: "json",
		Components:   []string{"foo"},
	}
	_ = instance.DisplayOutput()

	// Output:
	// {"items":[{"name":"some name one","filePath":"path-to-file","namespace":"lamb-namespace","api":{"version":"extensions/v1beta1","kind":"Deployment","deprecated-in":"v1.0.0","removed-in":"v1.0.0","replacement-api":"apps/v1","replacement-available-in":"v1.10.0","component":"foo"},"deprecated":true,"removed":true,"replacementAvailable":true},{"name":"some name two","api":{"version":"extensions/v1beta1","kind":"Deployment","deprecated-in":"v1.0.0","removed-in":"v1.0.0","replacement-api":"apps/v1","replacement-available-in":"v1.10.0","component":"foo"},"deprecated":true,"removed":true,"replacementAvailable":true}],"target-versions":{"foo":"v1.0.0"}}
}

func ExampleInstance_DisplayOutput_yaml() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		Components:   []string{"foo"},
		OutputFormat: "yaml",
	}
	_ = instance.DisplayOutput()

	// Output:
	// items:
	//     - name: some name one
	//       filePath: path-to-file
	//       namespace: lamb-namespace
	//       api:
	//         version: extensions/v1beta1
	//         kind: Deployment
	//         deprecated-in: v1.0.0
	//         removed-in: v1.0.0
	//         replacement-api: apps/v1
	//         replacement-available-in: v1.10.0
	//         component: foo
	//       deprecated: true
	//       removed: true
	//       replacementAvailable: true
	//     - name: some name two
	//       api:
	//         version: extensions/v1beta1
	//         kind: Deployment
	//         deprecated-in: v1.0.0
	//         removed-in: v1.0.0
	//         replacement-api: apps/v1
	//         replacement-available-in: v1.10.0
	//         component: foo
	//       deprecated: true
	//       removed: true
	//       replacementAvailable: true
	// target-versions:
	//     foo: v1.0.0
}

func ExampleInstance_DisplayOutput_csv() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		Components:   []string{"foo"},
		OutputFormat: "csv",
	}
	_ = instance.DisplayOutput()

	// Output:
	// NAME,NAMESPACE,KIND,VERSION,REPLACEMENT,DEPRECATED,DEPRECATED IN,REMOVED,REMOVED IN,REPL AVAIL,REPL AVAIL IN
	// some name one,lamb-namespace,Deployment,extensions/v1beta1,apps/v1,true,v1.0.0,true,v1.0.0,true,v1.10.0
	// some name two,<EXAMPLE>,Deployment,extensions/v1beta1,apps/v1,true,v1.0.0,true,v1.0.0,true,v1.10.0
}

func ExampleInstance_DisplayOutput_csv_customcolumns() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		Components:    []string{"foo"},
		OutputFormat:  "csv",
		CustomColumns: []string{"NAMESPACE", "NAME", "DEPRECATED IN", "DEPRECATED", "REPLACEMENT", "VERSION", "KIND", "COMPONENT", "FILEPATH"},
	}
	_ = instance.DisplayOutput()

	// Output:
	// NAMESPACE,NAME,DEPRECATED IN,DEPRECATED,REPLACEMENT,VERSION,KIND,COMPONENT,FILEPATH
	// lamb-namespace,some name one,v1.0.0,true,apps/v1,extensions/v1beta1,Deployment,foo,path-to-file
	// <EXAMPLE>,some name two,v1.0.0,true,apps/v1,extensions/v1beta1,Deployment,foo,<EXAMPLE>
}

func ExampleInstance_DisplayOutput_csv_noHeaders() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs: []*Output{
			testOutput1,
			testOutput2,
		},
		Components:   []string{"foo"},
		OutputFormat: "csv",
		NoHeaders:    true,
	}
	_ = instance.DisplayOutput()

	// Output:
	// some name one,lamb-namespace,Deployment,extensions/v1beta1,apps/v1,true,v1.0.0,true,v1.0.0,true,v1.10.0
	// some name two,<EXAMPLE>,Deployment,extensions/v1beta1,apps/v1,true,v1.0.0,true,v1.0.0,true,v1.10.0
}

func ExampleInstance_DisplayOutput_noOutput() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs: []*Output{
			testOutputNoOutput,
		},
		OutputFormat: "normal",
		Components:   []string{"foo"},
	}
	_ = instance.DisplayOutput()

	// Output: No output to display
}

func ExampleInstance_DisplayOutput_zeroLength() {
	instance := &Instance{
		TargetVersions: map[string]string{
			"foo": "v1.0.0",
		},
		Outputs:      []*Output{},
		OutputFormat: "normal",
	}
	_ = instance.DisplayOutput()

	// Output: There were no resources found with known deprecated apiVersions.
}

func TestGetReturnCode(t *testing.T) {

	type args struct {
		outputs                      []*Output
		ignoreDeprecations           bool
		ignoreRemovals               bool
		ignoreReplacementUnavailable bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "empty return zero",
			args: args{
				outputs:                      []*Output{},
				ignoreDeprecations:           false,
				ignoreRemovals:               false,
				ignoreReplacementUnavailable: false,
			},
			want: 0,
		},
		{
			name: "version is deprecated return two",
			args: args{
				outputs: []*Output{
					{
						APIVersion: &Version{
							DeprecatedIn: "v1.0.0",
							RemovedIn:    "v1.0.0",
							Component:    "foo",
						},
					},
				},
				ignoreDeprecations:           false,
				ignoreRemovals:               false,
				ignoreReplacementUnavailable: false,
			},
			want: 4,
		},
		{
			name: "version is deprecated ignore deprecations",
			args: args{
				outputs: []*Output{
					{
						APIVersion: &Version{
							DeprecatedIn: "v1.0.0",
							RemovedIn:    "v1.0.0",
							Component:    "foo",
						},
					},
				},
				ignoreDeprecations:           true,
				ignoreRemovals:               false,
				ignoreReplacementUnavailable: false,
			},
			want: 4,
		},
		{
			name: "version is removed",
			args: args{
				outputs: []*Output{
					{
						APIVersion: &Version{
							RemovedIn:    "v1.0.0",
							DeprecatedIn: "v1.12.0",
							Component:    "foo",
						},
					},
				},
				ignoreDeprecations:           false,
				ignoreRemovals:               false,
				ignoreReplacementUnavailable: false,
			},
			want: 4,
		},
		{
			name: "version is removed and replacement is unavailable",
			args: args{
				outputs: []*Output{
					{
						APIVersion: &Version{
							DeprecatedIn:           "v1.0.0",
							RemovedIn:              "v1.0.0",
							ReplacementAvailableIn: "v1.0.0",
							Component:              "foo",
						},
					},
				},
				ignoreDeprecations:           false,
				ignoreRemovals:               false,
				ignoreReplacementUnavailable: true,
			},
			want: 3,
		},
		{
			name: "version is deprecated and replacement is unavailable",
			args: args{
				outputs: []*Output{
					{
						APIVersion: &Version{
							DeprecatedIn:           "v1.0.0",
							RemovedIn:              "v1.0.0",
							ReplacementAvailableIn: "v1.0.0",
							Component:              "foo",
						},
					},
				},
				ignoreDeprecations: false,
				ignoreRemovals:     false,
			},
			want: 4,
		},
		{
			name: "version is deprecated and replacement is unavailable but ignored",
			args: args{
				outputs: []*Output{
					{
						APIVersion: &Version{
							DeprecatedIn:           "v1.0.0",
							RemovedIn:              "v1.0.0",
							ReplacementAvailableIn: "v1.0.0",
							Component:              "foo",
						},
					},
				},
				ignoreDeprecations:           false,
				ignoreRemovals:               false,
				ignoreReplacementUnavailable: true,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := &Instance{
				TargetVersions: map[string]string{
					"foo": "v1.0.0",
				},
				IgnoreDeprecations:            tt.args.ignoreDeprecations,
				IgnoreRemovals:                tt.args.ignoreRemovals,
				IgnoreUnavailableReplacements: tt.args.ignoreReplacementUnavailable,
				Outputs:                       tt.args.outputs,
			}
			got := instance.GetReturnCode()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOutput_GetReturnType(t *testing.T) {
	type fields struct {
		Name                   string
		Namespace              string
		FilePath               string
		APIType            	  *Type
		Deprecated             bool
		Removed                bool
		ReplacementAvailable   bool
		ReplacementAvailableIn string
		Component              string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "deprecated",
			fields: fields{
				APIVersion: &Version{
					DeprecatedIn: "v1.0.0",
				},
			},
			want: 4,
		},
		{
			name: "removed",
			fields: fields{
				APIVersion: &Version{
					RemovedIn: "v1.0.0",
				},
			},
			want: 4,
		},
		{
			name: "replacement unavailable",
			fields: fields{
				APIVersion: &Version{
					DeprecatedIn: "v1.0.0",
					RemovedIn:    "v1.0.0",
				},
			},
			want: 4,
		},
		{
			name: "replacement available",
			fields: fields{
				APIVersion: &Version{
					DeprecatedIn:           "v1.0.0",
					RemovedIn:              "v1.0.0",
					ReplacementAvailableIn: "v1.0.0",
				},
			},
			want: 4,
		},
		{
			name: "no issues",
			fields: fields{
				APIVersion: &Version{
					DeprecatedIn:           "",
					RemovedIn:              "",
					ReplacementAvailableIn: "",
				},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Output{
				Name:                   tt.fields.Name,
				Namespace:              tt.fields.Namespace,
				FilePath:               tt.fields.FilePath,
				APIVersion:             tt.fields.APIVersion,
				Deprecated:             tt.fields.Deprecated,
				Removed:                tt.fields.Removed,
				ReplacementAvailable:   tt.fields.ReplacementAvailable,
				ReplacementAvailableIn: tt.fields.Replacement
			}
			got := o.GetReturnType()
			assert.Equal(t, tt.want, got)
		})
	}
}
