// lambversionsfile makes the lamb versions.yaml file available to the lamb
// binary and others that import lamb packages.
package lambversionsfile

import (
	_ "embed"
)

var (
	//go:embed versions.yaml
	lambVersionsFileContent []byte
)

// Content returns the lamb versions.yaml file as a slice of bytes.
func Content() []byte {
	return lambVersionsFileContent
}
