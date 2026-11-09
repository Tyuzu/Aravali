package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func mmain() {
	// Set your target root directory here (e.g., "." for current directory)
	root := "."

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Look for files named "api.ts"
		if !d.IsDir() && d.Name() == "api.ts" {
			dir := filepath.Dir(path)
			parentDirName := filepath.Base(dir)

			// Construct the new filename: foldername + "API.ts"
			newFileName := parentDirName + "API.ts"
			targetPath := filepath.Join(dir, newFileName)

			fmt.Printf("Renaming: %s  ->  %s\n", path, targetPath)

			// Perform the rename operation
			err := os.Rename(path, targetPath)
			if err != nil {
				return fmt.Errorf("failed to rename %s: %w", path, err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking the directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nAll api.ts files successfully renamed!")
}
