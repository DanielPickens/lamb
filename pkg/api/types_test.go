package api

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test_data/types.yaml
var testtypesFile []byte

var mockInstance = Instance{
	Targettypes: map[string]string{
		"k8s":          "v1.0.0",
		"istio":        "1.0.0",
		"cert-manager": "v0.0.0",
	},
	Deprecatedtypes: []Type{
		testTypeDeployment,
	},
}

var testTypeDeploymentString = `deprecated-types:
- Type: extensions/v1beta1
  kind: Deployment
  deprecated-in: v1.0.0
  removed-in: v1.0.0
  replacement-api: apps/v1
  replacement-available-in: v1.0.0
  component: k8s
target-types:
  k8s: v1.0.0
  istio: v1.0.0
  cert-manager: v0.0.0`

var testTypeDeployment = Type{
	Name:                   "extensions/v1beta1",
	Kind:                   "Deployment",
	DeprecatedIn:           "v1.0.0",
	RemovedIn:              "v1.0.0",
	ReplacementAPI:         "apps/v1",
	ReplacementAvailableIn: "v1.0.0",
	Component:              "k8s",
}

func Test_jsonToStub(t *testing.T) {

	tests := []struct {
		name    string
		data    []byte
		want    []*Stub
		wantErr bool
	}{
		{
			name:    "json not stub",
			data:    []byte("{}"),
			want:    []*Stub{{}},
			wantErr: false,
		},
		{
			name:    "no data",
			data:    []byte(""),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "json is stub",
			data:    []byte(`{"kind": "foo", "apiType": "bar"}`),
			want:    []*Stub{{Kind: "foo", APIType: "bar"}},
			wantErr: false,
		},
		{
			name:    "json list is multiple stubs",
			data:    []byte(`{"kind": "List", "apiType": "v1", "items": [{"kind": "foo", "apiType": "bar"},{"kind": "bar", "apiType": "foo"}]}`),
			want:    []*Stub{{Kind: "foo", APIType: "bar"}, {Kind: "bar", APIType: "foo"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonToStub(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_yamlToStub(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []*Stub
		wantErr bool
	}{
		{
			name:    "yaml not stub",
			data:    []byte("foo: bar"),
			want:    []*Stub{{}},
			wantErr: false,
		},
		{
			name:    "not yaml",
			data:    []byte("*."),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "yaml is stub",
			data:    []byte("kind: foo\napiType: bar"),
			want:    []*Stub{{Kind: "foo", APIType: "bar"}},
			wantErr: false,
		},
		{
			name:    "yaml list is multiple stubs",
			data:    []byte("kind: List\napiType: v1\nitems:\n- kind: foo\n  apiType: bar\n- kind: bar\n  apiType: foo"),
			want:    []*Stub{{Kind: "foo", APIType: "bar"}, {Kind: "bar", APIType: "foo"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := yamlToStub(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_containsStub(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []*Stub
		wantErr bool
	}{
		{
			name:    "yaml not stub",
			data:    []byte("foo: bar"),
			want:    []*Stub{{}},
			wantErr: false,
		},
		{
			name:    "not yaml",
			data:    []byte("*."),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "yaml is stub",
			data:    []byte("kind: foo\napiType: bar"),
			want:    []*Stub{{Kind: "foo", APIType: "bar"}},
			wantErr: false,
		},
		{
			name:    "json not stub",
			data:    []byte("{}"),
			want:    []*Stub{{}},
			wantErr: false,
		},
		{
			name:    "empty string",
			data:    []byte(""),
			want:    nil,
			wantErr: false,
		},
		{
			name:    "no data",
			data:    []byte{},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "json is stub",
			data:    []byte(`{"kind": "foo", "apiType": "bar"}`),
			want:    []*Stub{{Kind: "foo", APIType: "bar"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := containsStub(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_IsTypeed(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []*Output
		wantErr bool
	}{
		{
			name:    "yaml no Type",
			data:    []byte("foo: bar"),
			want:    nil,
			wantErr: false,
		},
		{
			name:    "not json or yaml",
			data:    []byte("some text\nthat is not yaml"),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "yaml has Type",
			data:    []byte("kind: Deployment\napiType: extensions/v1beta1"),
			want:    []*Output{{APIType: &testTypeDeployment}},
			wantErr: false,
		},
		{
			name:    "yaml list has Type",
			data:    []byte("kind: List\napiType: v1\nitems:\n- kind: Deployment\n  apiType: extensions/v1beta1"),
			want:    []*Output{{APIType: &testTypeDeployment}},
			wantErr: false,
		},
		{
			name:    "json no Type",
			data:    []byte("{}"),
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty string",
			data:    []byte(""),
			want:    nil,
			wantErr: false,
		},
		{
			name:    "no data",
			data:    []byte{},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "json has Type",
			data:    []byte(`{"kind": "Deployment", "apiType": "extensions/v1beta1"}`),
			want:    []*Output{{APIType: &testTypeDeployment}},
			wantErr: false,
		},
		{
			name:    "json list has Type",
			data:    []byte(`{"kind": "List", "apiType": "v1", "items": [{"kind": "Deployment", "apiType": "extensions/v1beta1"}]}`),
			want:    []*Output{{APIType: &testTypeDeployment}},
			wantErr: false,
		},
		{
			name:    "not yaml",
			data:    []byte("*."),
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mockInstance.IsTypeed(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestType_IsDeprecatedIn(t *testing.T) {
	tests := []struct {
		name         string
		targettypes  map[string]string
		component    string
		want         bool
		deprecatedIn string
	}{
		{
			name:         "not deprecated yet 1.0.0",
			targettypes:  map[string]string{"foo": "v1.0.0"},
			component:    "foo",
			deprecatedIn: "v1.0.0",
			want:         false,
		},
		{
			name:         "equal values",
			targettypes:  map[string]string{"foo": "v1.0.0"},
			component:    "foo",
			deprecatedIn: "v1.0.0",
			want:         true,
		},
		{
			name:         "greater than",
			targettypes:  map[string]string{"foo": "v1.17.0"},
			component:    "foo",
			deprecatedIn: "v1.0.0",
			want:         true,
		},
		{
			name:         "Bad semVer",
			targettypes:  map[string]string{"foo": "foo"},
			component:    "foo",
			deprecatedIn: "v1.0.0",
			want:         false,
		},
		{
			name:         "blank deprecatedIn - not deprecated",
			targettypes:  map[string]string{"foo": "v1.0.0"},
			component:    "foo",
			deprecatedIn: "",
			want:         false,
		},
		{
			name:         "targetType not included",
			targettypes:  map[string]string{"one": "v1.0.0"},
			component:    "two",
			deprecatedIn: "v1.0.0",
			want:         false,
		},
	}
	for _, tt := range tests {
		deprecatedType := &Type{DeprecatedIn: tt.deprecatedIn, Component: tt.component}
		got := deprecatedType.isDeprecatedIn(tt.targettypes)
		assert.Equal(t, tt.want, got, "test failed: "+tt.name)
	}
}

func TestType_IsRemovedIn(t *testing.T) {

	tests := []struct {
		name        string
		targettypes map[string]string
		component   string
		want        bool
		removedIn   string
	}{
		{
			name:        "not removed yet 1.0.0",
			targettypes: map[string]string{"foo": "v1.0.0"},
			component:   "foo",
			removedIn:   "v1.0.0",
			want:        false,
		},
		{
			name:        "equal values",
			targettypes: map[string]string{"foo": "v1.0.0"},
			component:   "foo",
			removedIn:   "v1.0.0",
			want:        true,
		},
		{
			name:        "greater than",
			targettypes: map[string]string{"foo": "v1.17.0"},
			component:   "foo",
			removedIn:   "v1.0.0",
			want:        true,
		},
		{
			name:        "bad semVer",
			targettypes: map[string]string{"foo": "foo"},
			removedIn:   "v1.0.0",
			want:        false,
		},
		{
			name:        "blank removedIn - not removed",
			targettypes: map[string]string{"foo": "v1.0.0"},
			component:   "foo",
			removedIn:   "",
			want:        false,
		},
		{
			name:        "targettypes not included for component",
			targettypes: map[string]string{"one": "v1.0.0"},
			component:   "two",
			removedIn:   "v1.0.0",
			want:        false,
		},
	}
	for _, tt := range tests {
		removedType := &Type{RemovedIn: tt.removedIn, Component: tt.component}
		got := removedType.isRemovedIn(tt.targettypes)
		assert.Equal(t, tt.want, got, "test failed: "+tt.name)
	}
}

func TestType_isReplacementAvailableIn(t *testing.T) {
	tests := []struct {
		name                   string
		targettypes            map[string]string
		component              string
		want                   bool
		replacementAvailableIn string
	}{
		{
			name:                   "not available yet 1.0.0",
			targettypes:            map[string]string{"foo": "v1.0.0"},
			component:              "foo",
			replacementAvailableIn: "v1.0.0",
			want:                   false,
		},
		{
			name:                   "equal values",
			targettypes:            map[string]string{"foo": "v1.0.0"},
			component:              "foo",
			replacementAvailableIn: "v1.0.0",
			want:                   true,
		},
		{
			name:                   "greater than",
			targettypes:            map[string]string{"foo": "v1.17.0"},
			component:              "foo",
			replacementAvailableIn: "v1.0.0",
			want:                   true,
		},
		{
			name:                   "bad semVer",
			targettypes:            map[string]string{"foo": "foo"},
			replacementAvailableIn: "v1.0.0",
			want:                   false,
		},
		{
			name:                   "blank replacementAvailableIn - is available",
			targettypes:            map[string]string{"foo": "v1.0.0"},
			component:              "foo",
			replacementAvailableIn: "",
			want:                   false,
		},
		{
			name:                   "targettypes not included for component",
			targettypes:            map[string]string{"one": "v1.0.0"},
			component:              "two",
			replacementAvailableIn: "v1.0.0",
			want:                   false,
		},
	}
	for _, tt := range tests {
		removedType := &Type{ReplacementAvailableIn: tt.replacementAvailableIn, Component: tt.component}
		got := removedType.isReplacementAvailableIn(tt.targettypes)
		assert.Equal(t, tt.want, got, "test failed: "+tt.name)
	}
}

func ExampleInstance_printtypesTabular() {
	instance := Instance{
		Deprecatedtypes: []Type{
			testTypeDeployment,
			{Kind: "testkind", Name: "testname", DeprecatedIn: "", RemovedIn: "", ReplacementAvailableIn: "", Component: "custom"},
		},
	}
	_ = instance.printtypesTabular()

	// Output:
	// KIND-------- NAME---------------- DEPRECATED IN-- REMOVED IN-- REPLACEMENT-- REPL AVAIL IN-- COMPONENT--
	// Deployment-- extensions/v1beta1-- v1.0.0--------- v1.0.0----- apps/v1------ v1.0.0-------- k8s--------
	// testkind---- testname------------ n/a------------ n/a--------- n/a---------- n/a------------ custom-----
}

func ExampleInstance_printtypesTabular_noHeaders() {
	instance := Instance{
		Deprecatedtypes: []Type{
			testTypeDeployment,
			{Kind: "testkind", Name: "testname", DeprecatedIn: "", RemovedIn: "", ReplacementAvailableIn: "", Component: "custom"},
		},
		NoHeaders: true,
	}
	_ = instance.printtypesTabular()

	// Output:
	// Deployment-- extensions/v1beta1-- v1.0.0-- v1.0.0-- apps/v1-- v1.0.0-- k8s-----
	// testkind---- testname------------ n/a----- n/a------ n/a------ n/a------ custom--
}

func ExampleInstance_PrintTypeList_json() {
	instance := Instance{
		Deprecatedtypes: []Type{testTypeDeployment},
	}
	_ = instance.PrintTypeList("json")

	// Output:
	// {"deprecated-types":[{"Type":"extensions/v1beta1","kind":"Deployment","deprecated-in":"v1.0.0","removed-in":"v1.0.0","replacement-api":"apps/v1","replacement-available-in":"v1.0.0","component":"k8s"}]}
}

func ExampleInstance_PrintTypeList_yaml() {
	instance := Instance{
		Deprecatedtypes: []Type{testTypeDeployment},
	}
	_ = instance.PrintTypeList("yaml")

	// Output:
	// deprecated-types:
	//     - Type: extensions/v1beta1
	//       kind: Deployment
	//       deprecated-in: v1.0.0
	//       removed-in: v1.0.0
	//       replacement-api: apps/v1
	//       replacement-available-in: v1.0.0
	//       component: k8s
}

func ExampleInstance_PrintTypeList_normal() {
	instance := Instance{
		Deprecatedtypes: []Type{testTypeDeployment},
	}
	_ = instance.PrintTypeList("normal")

	// Output:
	// KIND-------- NAME---------------- DEPRECATED IN-- REMOVED IN-- REPLACEMENT-- REPL AVAIL IN-- COMPONENT--
	// Deployment-- extensions/v1beta1-- v1.0.0--------- v1.0.0----- apps/v1------ v1.0.0-------- k8s--------
}

func ExampleInstance_PrintTypeList_wide() {
	instance := Instance{
		Deprecatedtypes: []Type{testTypeDeployment},
	}
	_ = instance.PrintTypeList("wide")

	// Output:
	// KIND-------- NAME---------------- DEPRECATED IN-- REMOVED IN-- REPLACEMENT-- REPL AVAIL IN-- COMPONENT--
	// Deployment-- extensions/v1beta1-- v1.0.0--------- v1.0.0----- apps/v1------ v1.0.0-------- k8s--------
}

func ExampleInstance_PrintTypeList_badformat() {
	instance := Instance{
		Deprecatedtypes: []Type{testTypeDeployment},
	}
	_ = instance.PrintTypeList("foo")

	// Output:
	// The output format must one of (normal|wide|json|yaml)
}

func Test_isDuplicate(t *testing.T) {

	tests := []struct {
		name string
		a    Type
		b    Type
		want bool
	}{
		{
			name: "is duplicate",
			a:    Type{Kind: "Deployment", Name: "apps/v1"},
			b:    Type{Kind: "Deployment", Name: "apps/v1"},
			want: true,
		},
		{
			name: "is not duplicate",
			a:    Type{Kind: "Deployment", Name: "extensions/v1beta1"},
			b:    Type{Kind: "Deployment", Name: "apps/v1"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDuplicate(tt.a, tt.b)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestType_isContainedIn(t *testing.T) {
	tests := []struct {
		name     string
		Type     Type
		TypeList []Type
		want     bool
	}{
		{
			name: "true",
			Type: Type{Kind: "Deployment", Name: "extensions/v1beta1"},
			TypeList: []Type{
				{Kind: "Deployment", Name: "extensions/v1beta1"},
				{Kind: "Deployment", Name: "apps/v1"},
			},
			want: true,
		},
		{
			name: "false",
			Type: Type{Kind: "Deployment", Name: "extensions/v1beta1"},
			TypeList: []Type{
				{Kind: "Deployment", Name: "apps/v1"},
				{Kind: "Deployment", Name: "extensions/v1beta2"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.Type.isContainedIn(tt.TypeList)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCombineAdditionaltypes(t *testing.T) {
	type args struct {
		additional []Type
		defaults   []Type
	}
	tests := []struct {
		name     string
		args     args
		want     []Type
		wantErr  bool
		errorMsg string
	}{
		{
			name: "error combining due to duplicate",
			args: args{
				additional: []Type{
					{Kind: "Deployment", Name: "apps/v1"},
				},
				defaults: []Type{
					{Kind: "Deployment", Name: "apps/v1"},
				},
			},
			wantErr:  true,
			errorMsg: "duplicate cannot be added to defaults: Deployment apps/v1",
		},
		{
			name: "error combining due to duplicate",
			args: args{
				additional: []Type{
					{Kind: "Deployment", Name: "extensions/v1beta1"},
				},
				defaults: []Type{
					{Kind: "Deployment", Name: "apps/v1"},
				},
			},
			wantErr: false,
			want: []Type{
				{Kind: "Deployment", Name: "apps/v1"},
				{Kind: "Deployment", Name: "extensions/v1beta1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CombineAdditionaltypes(tt.args.additional, tt.args.defaults)
			if tt.wantErr {
				assert.EqualError(t, err, tt.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestMarshaltypes(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []Type
		wantErr bool
	}{
		{
			name:    "no error",
			data:    []byte(testTypeDeploymentString),
			want:    []Type{testTypeDeployment},
			wantErr: false,
		},
		{
			name:    "unmarshal error",
			data:    []byte(`foo`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := UnMarshaltypes(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.want, got)
			}
		})
	}
}

func TestGetDefaultTypeList(t *testing.T) {

	// This test will ensure that the types.yaml file is well-formed and doesn't break anything.
	defaulttypes, defaultTargettypes, err := GetDefaultTypeList(testtypesFile)
	assert.NoError(t, err)
	assert.NotNil(t, defaulttypes)
	assert.NotNil(t, defaultTargettypes)
}

func TestInstance_checkType(t *testing.T) {
	tests := []struct {
		name     string
		instance *Instance
		stub     *Stub
		want     *Type
	}{
		{
			name: "empty kind",
			instance: &Instance{
				Deprecatedtypes: []Type{
					{Kind: "", Name: "cert-manager.k8s.io", Component: "cert-manager"},
				},
			},
			stub: &Stub{
				Kind:    "any",
				APIType: "cert-manager.k8s.io",
				Metadata: StubMeta{
					Name:      "foo",
					Namespace: "foobar",
				},
			},
			want: &Type{
				Name:      "cert-manager.k8s.io",
				Kind:      "any",
				Component: "cert-manager",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.instance.checkType(tt.stub)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
