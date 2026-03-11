package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	dir := t.TempDir()
	ws, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if ws.Root() != dir {
		// On Windows, TempDir may not match abs exactly
		abs, _ := filepath.Abs(dir)
		if ws.Root() != abs {
			t.Errorf("Root = %q, want %q", ws.Root(), dir)
		}
	}
}

func TestWriteAndReadFile(t *testing.T) {
	ws, _ := New(t.TempDir())
	content := "package main\n\nfunc main() {}\n"

	err := ws.WriteFile("myproject/main.go", content)
	if err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	got, err := ws.ReadFile("myproject/main.go")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if got != content {
		t.Errorf("ReadFile = %q, want %q", got, content)
	}
}

func TestProductDir(t *testing.T) {
	ws, _ := New(t.TempDir())
	dir := ws.ProductDir("alpha")

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("ProductDir not created: %v", err)
	}
	if !info.IsDir() {
		t.Error("ProductDir is not a directory")
	}
}

func TestFileExists(t *testing.T) {
	ws, _ := New(t.TempDir())

	if ws.FileExists("nonexistent.go") {
		t.Error("FileExists should return false for missing file")
	}

	ws.WriteFile("exists.go", "package x")
	if !ws.FileExists("exists.go") {
		t.Error("FileExists should return true for existing file")
	}
}

func TestReadSourceFilesSkipsHiveAndProducts(t *testing.T) {
	dir := t.TempDir()
	prod := &Product{Name: "test", Dir: dir}

	// Create source files that should be included
	writeFile(t, dir, "main.go", "package main")
	writeFile(t, dir, "lib/util.go", "package lib")

	// Create .hive telemetry files that should be skipped
	writeFile(t, dir, ".hive/telemetry.json", `{"tokens": 100}`)
	writeFile(t, dir, ".hive/run-001.json", `{"status": "done"}`)

	// Create products files that should be skipped
	writeFile(t, dir, "products/app/main.go", "package main // generated")
	writeFile(t, dir, "products/app/go.mod", "module app")

	files, err := prod.ReadSourceFiles()
	if err != nil {
		t.Fatalf("ReadSourceFiles: %v", err)
	}

	// Should include source files
	if _, ok := files["main.go"]; !ok {
		t.Error("missing main.go")
	}
	if _, ok := files["lib/util.go"]; !ok {
		t.Error("missing lib/util.go")
	}

	// Should NOT include .hive or products files
	for path := range files {
		if strings.HasPrefix(path, ".hive") {
			t.Errorf("should skip .hive file: %s", path)
		}
		if strings.HasPrefix(path, "products") {
			t.Errorf("should skip products file: %s", path)
		}
	}

	if len(files) != 2 {
		t.Errorf("ReadSourceFiles = %d files, want 2; got: %v", len(files), keys(files))
	}
}

func writeFile(t *testing.T, base, rel, content string) {
	t.Helper()
	path := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func TestListFiles(t *testing.T) {
	ws, _ := New(t.TempDir())
	ws.WriteFile("proj/main.go", "package main")
	ws.WriteFile("proj/lib/util.go", "package lib")

	files, err := ws.ListFiles("proj")
	if err != nil {
		t.Fatalf("ListFiles: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("ListFiles = %d files, want 2", len(files))
	}
}
