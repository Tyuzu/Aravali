package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Comment syntax based on file extension
var commentStyles = map[string]string{
	".go":   "// File: %s\n\n",
	".js":   "// File: %s\n\n",
	".ts":   "// File: %s\n\n",
	".py":   "# File: %s\n\n",
	".sh":   "# File: %s\n\n",
	".yaml": "# File: %s\n\n",
	".yml":  "# File: %s\n\n",
	".css":  "/* File: %s */\n\n",
	".c":    "/* File: %s */\n\n",
	".cpp":  "/* File: %s */\n\n",
	".h":    "/* File: %s */\n\n",
	".html": "<!-- File: %s -->\n\n",
}

func main() {
	rootDir := "." // Target directory (change to your target path if needed)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files/folders (e.g., .git)
		if d.IsDir() || strings.HasPrefix(d.Name(), ".") || strings.Contains(path, "/.") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		format, exists := commentStyles[ext]
		if !exists {
			return nil // Skip unsupported file types
		}

		// Normalize path slashes for consistency across OS platforms
		normalizedPath := filepath.ToSlash(path)
		header := fmt.Sprintf(format, normalizedPath)

		// Read existing file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		// Avoid duplicate header insertion if script is re-run
		if bytes.HasPrefix(content, []byte(strings.TrimSpace(header))) {
			fmt.Printf("Skipped (already present): %s\n", normalizedPath)
			return nil
		}

		// Combine new header + original content
		var newContent bytes.Buffer
		newContent.WriteString(header)
		newContent.Write(content)

		// Overwrite the file safely with combined content
		info, err := d.Info()
		if err != nil {
			return err
		}

		err = os.WriteFile(path, newContent.Bytes(), info.Mode())
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}

		fmt.Printf("Updated: %s\n", normalizedPath)
		return nil
	})

	if err != nil {
		fmt.Printf("Error processing files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done!")
}
