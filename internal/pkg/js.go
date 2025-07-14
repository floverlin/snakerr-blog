package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func UpdateJSImports(root string, version string) error {
	re := regexp.MustCompile(`(?m)(import\s+(?:[\w*\s{},]*?\s+from\s+)?["'])(\.{1,2}/[^"']+?\.js)(\?v=[^"']+)?(["'])`)

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".js") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		original := string(data)

		updated := re.ReplaceAllStringFunc(original, func(match string) string {
			return re.ReplaceAllString(match, fmt.Sprintf(`${1}${2}?v=%s$4`, version))
		})

		if updated != original {
			err = os.WriteFile(path, []byte(updated), 0644)
			if err != nil {
				return fmt.Errorf("error writing %s: %w", path, err)
			}
		}
		return nil
	})
}
