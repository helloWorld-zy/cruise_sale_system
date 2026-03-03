package handler

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// TestTask19_NoEmptyFunctionBodiesInProduction enforces that production code under internal/
// does not use empty function bodies as fake implementations.
func TestTask19_NoEmptyFunctionBodiesInProduction(t *testing.T) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve caller path")
	}

	internalDir := filepath.Join(filepath.Dir(currentFile), "..")
	emptyFuncRe := regexp.MustCompile(`func\s*(\([^\)]*\))?\s*[A-Za-z_][A-Za-z0-9_]*\s*\([^\)]*\)\s*(\([^\)]*\)|[A-Za-z0-9_\*\.\[\]]+)?\s*\{\s*\}`)

	var violations []string
	err := filepath.WalkDir(internalDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == "vendor" || name == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if emptyFuncRe.Match(content) {
			rel, relErr := filepath.Rel(internalDir, path)
			if relErr != nil {
				rel = path
			}
			violations = append(violations, rel)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scan internal code failed: %v", err)
	}
	if len(violations) > 0 {
		t.Fatalf("found empty function bodies in production code: %v", violations)
	}
}
