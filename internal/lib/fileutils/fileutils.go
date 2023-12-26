package fileutils

import "path/filepath"

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
