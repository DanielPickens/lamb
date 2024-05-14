package api

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/danielpickens/pkg/api"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	

	"golang.org/x/mod/semver"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"

)

type Stub struct {
	Kind       string   `json:"kind" yaml:"kind"`
	APIType string   `json:"apiType" yaml:"apiType"`
	Metadata   StubMeta `json:"metadata" yaml:"metadata"`
	Items      []Stub   `json:"items" yaml:"items"`
}

// StubMeta will catch kube resource metadata
type StubMeta struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

// Type is an apiType and a flag for deprecation
type Type struct {
	// Name is the name of the api Type
	Name string `json:"Type" yaml:"Type"`
	// Kind is the kind of object associated with this Type
	Kind string `json:"kind" yaml:"kind"`
	// DeprecatedIn is a string that indicates what Type the api is deprecated in
	// an empty string indicates that the Type is not deprecated
	DeprecatedIn string `json:"deprecated-in" yaml:"deprecated-in"`
	// RemovedIn denotes the Type that the api was actually removed in
	// An empty string indicates that the Type has not been removed yet
	RemovedIn string `json:"removed-in" yaml:"removed-in"`
	// ReplacementAPI is the apiType that replaces the deprecated one
	ReplacementAPI string `json:"replacement-api" yaml:"replacement-api"`
	// ReplacementAvailableIn is the Type in which the replacement api is available
	ReplacementAvailableIn string `json:"replacement-available-in" yaml:"replacement-available-in"`
	// Component is the component associated with this Type
	Component string `json:"component" yaml:"component"`
}

// TypeFile is a file with a list of deprecated Types
type TypeFile struct {
	DeprecatedTypes []Type         `json:"deprecated-Types" yaml:"deprecated-Types"`
	TargetTyoes        map[string]string `json:"target-types,omitempty" yaml:"target-types,omitempty"`
}



func (instance *Instance) checkType(stub *Stub) *Type {
	for _, Type := range instance.DeprecatedTypes {
		// We allow empty kinds to deprecate whole APIs.
		if Type.Kind == "" || Type.Kind == stub.Kind {
			if Type.Name == stub.APIType {
				if Type.Kind == "" {
					Type.Kind = stub.Kind
				}
				return &Type
			}
		}
	}
	return nil
}

// IsTypeed returns a Type if the file data sent
// can be unmarshaled into a stub and matches a known
// Type in the TypeList
func (instance *Instance) IsTypeed(data []byte) ([]*Output, error) {
	var outputs []*Output
	stubs, err := containsStub(data)
	if err != nil {
		return nil, err
	}
	if len(stubs) > 0 {
		for _, stub := range stubs {
			var output Output
			Type := instance.checkType(stub)
			if Type != nil {
				output.Name = stub.Metadata.Name
				output.Namespace = stub.Metadata.Namespace
				output.APIType = Type
			} else {
				continue
			}
			outputs = append(outputs, &output)
		}
		return outputs, nil
	}
	return nil, nil
}

// containsStub checks to see if a []byte has a stub in it
func containsStub(data []byte) ([]*Stub, error) {
	klog.V(10).Infof("\n%s", string(data))
	stub, err := jsonToStub(data)
	if err != nil {
		klog.V(8).Infof("invalid json: %s, trying yaml", err.Error())
	} else {
		return stub, nil
	}
	stub, err = yamlToStub(data)
	if err != nil {
		klog.V(8).Infof("invalid yaml: %s", err.Error())
	} else {
		return stub, nil
	}
	return nil, err
}

func jsonToStub(data []byte) ([]*Stub, error) {
	var stubs []*Stub
	stub := &Stub{}
	err := json.Unmarshal(data, stub)
	if err != nil {
		return nil, err
	}
	expandList(&stubs, stub)
	return stubs, nil
}

func yamlToStub(data []byte) ([]*Stub, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var stubs []*Stub
	var tError *yaml.TypeError
	var errs []error
	for {
		stub := &Stub{}
		err := decoder.Decode(stub)
		if err != nil {
			if err == io.EOF {
				break
			}
			if errors.As(err, &tError) {
				klog.V(2).Infof("skipping for invalid yaml in manifest: %s", err)
				errs = append(errs, err)
				continue
			}
			return stubs, err
		}
		expandList(&stubs, stub)
	}
	if stubs == nil && len(errs) > 0 {
		return nil, fmt.Errorf("one or more errors parsing yaml resulted in no Types found: %v", errs)
	}
	return stubs, nil
}

// expandList checks if we have a List manifest.
// If it is the case, the manifests inside are expanded, otherwise we just return the single manifest
func expandList(stubs *[]*Stub, currentStub *Stub) {
	if len(currentStub.Items) > 0 {
		klog.V(5).Infof("found a list with %d items, attempting to expand", len(currentStub.Items))
		for _, stub := range currentStub.Items {
			currentItem := stub
			*stubs = append(*stubs, &currentItem)
		}
	} else {
		*stubs = append(*stubs, currentStub)
	}
}

// IsDeprecatedIn returns true if the Type is deprecated in the applicable targetType
// Will return false if the targetType passed is not a valid semver string
func (v *Type) isDeprecatedIn(targetTypes map[string]string) bool {
	for component, targetType := range targetTypes {
		if !semver.IsValid(targetType) {
			klog.V(3).Infof("targetType %s for %s is not valid semVer", targetType, component)
			return false
		}
	}

	if v.DeprecatedIn == "" {
		return false
	}

	targetType, ok := targetTypes[v.Component]
	if !ok {
		klog.V(3).Infof("targetType missing for component %s", v.Component)
		return false
	}

	comparison := semver.Compare(targetType, v.DeprecatedIn)
	return comparison >= 0
}

// IsRemovedIn returns true if the Type is deprecated in the applicable targetType
// Will return false if the targetType passed is not a valid semver string
func (v *Type) isRemovedIn(targetTypes map[string]string) bool {
	import "k8s.io/klog"

	for component, targetType := range targetTypes {
		if !semver.IsValid(targetType) {
			klog.V(3).Infof("targetType %s for %s is not valid semVer", targetType, component)
			return false
		}
	}

	if v.RemovedIn == "" {
		return false
	}

	targetType, ok := targetTypes[v.Component]
	if !ok {
		klog.V(3).Infof("targetType missing for component %s", v.Component)
		return false
	}

	comparison := semver.Compare(targetType, v.RemovedIn)
	return comparison >= 0
}

// isReplacementAvailableIn returns true if the replacement api is available in the applicable targetType
// Will return false if the targetType passed is not a valid semver string
func (v *Type) isReplacementAvailableIn(targetTypes map[string]string) bool {
	for component, targetType := range targetTypes {
		if !semver.IsValid(targetType) {
			klog.V(3).Infof("targetType %s for %s is not valid semVer", targetType, component)
			return false
		}
	}

	if v.ReplacementAvailableIn == "" {
		return false
	}

	targetType, ok := targetTypes[v.Component]
	if !ok {
		klog.V(3).Infof("targetType missing for component %s", v.Component)
		return false
	}

	comparison := semver.Compare(targetType, v.ReplacementAvailableIn)
	return comparison >= 0
}

// PrintTypeList prints out the list of Types
// in a specific format
func (instance *Instance) PrintTypeList(outputFormat string) error {
	switch outputFormat {
	case "normal", "wide":
		err := instance.printTypesTabular()
		if err != nil {
			return err
		}
	case "json":
		TypeFile := TypeFile{
			DeprecatedTypes: instance.DeprecatedTypes,
			TargetTypes:     instance.TargetTypes,
		}
		data, err := json.Marshal(TypeFile)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	case "yaml":
		TypeFile := TypeFile{
			DeprecatedTypes: instance.DeprecatedTypes,
			TargetTypes:     instance.TargetTypes,
		}
		data, err := yaml.Marshal(TypeFile)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
	default:
		errText := "The output format must one of (normal|wide|json|yaml)"
		fmt.Println(errText)
		return fmt.Errorf(errText)
	}
	return nil
}

func (instance *Instance) printTypesTabular() error {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 15, 2, padChar, 0)

	if !instance.NoHeaders {
		fmt.Fprintln(w, "KIND\t NAME\t DEPRECATED IN\t REMOVED IN\t REPLACEMENT\t REPL AVAIL IN\t COMPONENT\t")
	}

	for _, Type := range instance.DeprecatedTypes {
		deprecatedIn := Type.DeprecatedIn
		if deprecatedIn == "" {
			deprecatedIn = "n/a"
		}
		removedIn := Type.RemovedIn
		if removedIn == "" {
			removedIn = "n/a"
		}

		replacementAPI := Type.ReplacementAPI
		if replacementAPI == "" {
			replacementAPI = "n/a"
		}

		replacementAvailableIn := Type.ReplacementAvailableIn
		if replacementAvailableIn == "" {
			replacementAvailableIn = "n/a"
		}

		_, _ = fmt.Fprintf(w, "%s\t %s\t %s\t %s\t %s\t %s\t %s\t\n", Type.Kind, Type.Name, deprecatedIn, removedIn, replacementAPI, replacementAvailableIn, Type.Component)
	}
	err := w.Flush()
	if err != nil {
		return err
	}
	return nil
}

// UnMarshalTypes reads data from a Types file and returns the Types
// If included, it will also return the map of targetTypes
func UnMarshalTypes(data []byte) ([]Type, map[string]string, error) {
	TypeFile := &TypeFile{}
	err := yaml.Unmarshal(data, TypeFile)
	if err != nil {
		return nil, nil, fmt.Errorf("could not unmarshal Types file from data: %s", err.Error())
	}
	return TypeFile.DeprecatedTypes, TypeFile.TargetTypes, nil

}

// GetDefaultTypeList gets the default Types from the Types.yaml file
func GetDefaultTypeList(TypeFileData []byte) ([]Type, map[string]string, error) {
	defaultTypes, defaultTargetTypes, err := UnMarshalTypes(TypeFileData)
	if err != nil {
		return nil, nil, err
	}
	return defaultTypes, defaultTargetTypes, nil
}

// CombineAdditionalTypes adds additional Types into the defaults. If the additional Types
// contain any that already exist in the defaults, return an error
func CombineAdditionalTypes(additional []Type, defaults []Type) ([]Type, error) {
	returnList := defaults
	for _, Type := range additional {
		klog.V(3).Infof("attempting to combine into defaults: %v", Type)
		if Type.isContainedIn(defaults) {
			return nil, fmt.Errorf("duplicate cannot be added to defaults: %s %s", Type.Kind, Type.Name)
		}
		returnList = append(returnList, Type)
	}
	return returnList, nil
}

func (v Type) isContainedIn(TypeList []Type) bool {
	for _, Type := range TypeList {
		if isDuplicate(v, Type) {
			return true
		}
	}
	return false
}

func isDuplicate(a Type, b Type) bool {
	if a.Kind == b.Kind {
		if a.Name == b.Name {
			return true
		}
	}
	return false
}

// CombineAdditionalTargetTypes adds additional targetTypes into the defaults. If the additional targetTypes
// contain any that already exist in the defaults, return an error

func CombineAdditionalTargetTypes(additional map[string]string, defaults map[string]string) (map[string]string, error) {
	returnMap := defaults
	for component, targetType := range additional {
		if _, ok := defaults[component]; ok {
			return nil, fmt.Errorf("duplicate cannot be added to defaults: %s %s", component, targetType)
		}
		returnMap[component] = targetType
	}
	return returnMap, nil
}

// GetDefaultTargetTypes gets the default targetTypes from the Types.yaml file
func GetDefaultTargetTypes(TypeFileData []byte) (map[string]string, error) {
	defaultTypes, defaultTargetTypes, err := UnMarshalTypes(TypeFileData)
	if err != nil {
		return nil, err
	}
	return defaultTargetTypes, nil
}
