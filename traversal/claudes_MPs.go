// https://claude.ai/chat/08a5377a-e74a-4c74-a14e-525e6e229a21

package traversal

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func getMaterializedPaths(fsys fs.FS, rootPath string) ([]string, error) {
	// Clean and normalize the rootPath
	rootPath = filepath.Clean(rootPath)

	// Check if rootPath exists
	fileInfo, err := os.Stat(rootPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the path doesn't exist, return it as the only entry
			return []string{filepath.ToSlash(rootPath)}, nil
		}
		return nil, err
	}

if !fileInfo.IsDir() {
		// If rootPath is not a directory, return it as the only entry
		return []string{filepath.ToSlash(rootPath)}, nil
	}

	// Ensure rootPath ends with a separator for directories
	if !strings.HasSuffix(rootPath, string(filepath.Separator)) {
		rootPath += string(filepath.Separator)
	}

	paths := []string{filepath.ToSlash(rootPath)}  // Root directory with trailing slash

	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipDir
		}
		
		if path != "." {
			fullPath := filepath.Join(rootPath, path)
			normalizedPath := filepath.ToSlash(fullPath)
			if d.IsDir() && !strings.HasSuffix(normalizedPath, "/") {
				normalizedPath += "/"
			}
			paths = append(paths, normalizedPath)
		}
		
		return nil
	})

if err != nil {
		return nil, err
	}

	sort.Strings(paths)
	return paths, nil
}

