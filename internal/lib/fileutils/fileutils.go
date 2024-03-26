// fileutils отвечает за обработку путей
package fileutils

import "path/filepath"

// CreateFullPathFromRelative полный путь файла по относительному.
func CreateFullPathFromRelative(relativePath string) (fullPath string, err error) {

	validPath := filepath.FromSlash(relativePath)

	if filepath.IsAbs(validPath) {
		return validPath, nil
	}

	dir, err := filepath.Abs("")
	if err != nil {
		return "", err
	}

	return dir + validPath, nil
}
