package vfs4

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/twpayne/go-vfs"
)

func TestVfsScoped_BasicScoping(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	// Create test files
	if err := os.MkdirAll(filepath.Join(tmpDir, "dir/subdir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "dir/file.txt"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"file in root", "file.txt", false},
		{"file in subdir", "dir/file.txt", false},
		{"file in nested subdir", "dir/subdir/file.txt", true}, // doesn't exist
		{"empty path resolves to root", "", false},
		{"dot resolves to root", ".", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := v.Open(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				f.Close()
			}
		})
	}
}

func TestVfsScoped_EscapePrevention(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(tmpDir, "dir"), 0755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"parent dir", "..", true},
		{"sibling escape", "../ sibling", true},
		{"deep escape", "dir/../../other", true},
		{"root parent", tmpDir + "/..", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := v.Open(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

func TestVfsScoped_AbsolutePathHandling(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"absolute inside root", tmpDir + "/file.txt", false},
		{"absolute inside nested", tmpDir + "/dir/file.txt", true}, // doesn't exist
		{"absolute outside root", "/home/user/other", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := v.Open(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				f.Close()
			}
		})
	}
}

func TestVfsScoped_AbsolutePathEscape(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"escape via absolute", tmpDir + "/../other", true},
		{"escape deep via absolute", tmpDir + "/dir/../../other", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := v.Open(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

func TestVfsScoped_SymlinkEscapes(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	// Create a file outside scope
	outsideDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(outsideDir, "secret.txt"), []byte("secret"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create symlink inside scope pointing outside
	err = os.Symlink(outsideDir, filepath.Join(tmpDir, "link_to_outside"))
	if err != nil {
		t.Fatal(err)
	}

	// Symlink inside scope pointing to file inside scope
	if err := os.MkdirAll(filepath.Join(tmpDir, "dir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "dir/target.txt"), []byte("target"), 0644); err != nil {
		t.Fatal(err)
	}
	err = os.Symlink(filepath.Join(tmpDir, "dir/target.txt"), filepath.Join(tmpDir, "link_to_inside"))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"symlink to outside", "link_to_outside", true}, // target escapes scope
		{"symlink to inside", "link_to_inside", false},  // target inside scope
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := v.Open(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				f.Close()
			}
		})
	}
}

func TestVfsScoped_AllOperations(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := v.Chmod("test.txt", 0600); err != nil {
		t.Errorf("Chmod error = %v", err)
	}

	if err := v.Chtimes("test.txt", time.Now(), time.Now()); err != nil {
		t.Errorf("Chtimes error = %v", err)
	}

	if _, err := v.Create("newfile.txt"); err != nil {
		t.Errorf("Create error = %v", err)
	}

	if _, err := v.Stat("test.txt"); err != nil {
		t.Errorf("Stat error = %v", err)
	}

	if _, err := v.Lstat("test.txt"); err != nil {
		t.Errorf("Lstat error = %v", err)
	}

	if data, err := v.ReadFile("test.txt"); err != nil {
		t.Errorf("ReadFile error = %v", err)
	} else if string(data) != "hello" {
		t.Errorf("ReadFile content = %q, want %q", string(data), "hello")
	}

	if entries, err := v.ReadDir("."); err != nil {
		t.Errorf("ReadDir error = %v", err)
	} else if len(entries) == 0 {
		t.Error("ReadDir returned no entries")
	}

	if err := v.WriteFile("written.txt", []byte("world"), 0644); err != nil {
		t.Errorf("WriteFile error = %v", err)
	}

	if err := v.Truncate("written.txt", 0); err != nil {
		t.Errorf("Truncate error = %v", err)
	}

	if err := v.Rename("written.txt", "renamed.txt"); err != nil {
		t.Errorf("Rename error = %v", err)
	}

	if err := v.Remove("renamed.txt"); err != nil {
		t.Errorf("Remove error = %v", err)
	}

	if err := v.RemoveAll("newfile.txt"); err != nil {
		t.Errorf("RemoveAll error = %v", err)
	}

	if _, err := v.OpenFile("opened.txt", os.O_RDWR|os.O_CREATE, 0644); err != nil {
		t.Errorf("OpenFile error = %v", err)
	}

	if _, err := v.Glob("*.txt"); err != nil {
		t.Errorf("Glob error = %v", err)
	}

	if ps := v.PathSeparator(); ps != '/' {
		t.Errorf("PathSeparator = %v, want '/'", ps)
	}

	if _, err := v.RawPath("test.txt"); err != nil {
		// RawPath may not be supported by all filesystems
		t.Logf("RawPath error (may be expected): %v", err)
	}
}

func TestVfsScoped_MkdirAll(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	if err := MkdirAllScoped(v, "a/b/c", 0755); err != nil {
		t.Errorf("MkdirAllScoped error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "a/b/c")); err != nil {
		t.Errorf("Directory not created: %v", err)
	}
}

func TestVfsScoped_Contains(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(tmpDir, "dir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "dir/file.txt"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	if ok, err := ContainsScoped(v, "dir", "file.txt"); err != nil {
		t.Errorf("Contains error = %v", err)
	} else if !ok {
		t.Errorf("Contains returned false, want true")
	}

	if ok, err := ContainsScoped(v, "dir", "nonexistent.txt"); err != nil {
		t.Errorf("Contains error = %v", err)
	} else if ok {
		t.Errorf("Contains returned true, want false")
	}
}

func TestVfsScoped_Walk(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(tmpDir, "dir/subdir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "dir/file2.txt"), []byte("2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "dir/subdir/file3.txt"), []byte("3"), 0644); err != nil {
		t.Fatal(err)
	}

	var found []string
	err = WalkScoped(v, ".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		found = append(found, path)
		return nil
	})
	if err != nil {
		t.Errorf("Walk error = %v", err)
	}

	if len(found) == 0 {
		t.Error("Walk found no files")
	}
}

func TestVfsScoped_InvalidRoot(t *testing.T) {
	_, err := NewVfsScoped("relative/path", vfs.OSFS)
	if err == nil {
		t.Error("NewVfsScoped accepted relative path, want error")
	}
}

func TestVfsScoped_AllowLevel(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(tmpDir, "l1dir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, "l1dir/l2dir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, "l1dir/l2dir/l3dir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "l1dir/file.txt"), []byte("l1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "l1dir/l2dir/file.txt"), []byte("l2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "l1dir/l2dir/l3dir/file.txt"), []byte("l3"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		maxLevel int
		path     string
		wantErr  bool
	}{
		{"unlimited level 0 - root", 0, "", false},
		{"unlimited level 0 - l1 file", 0, "l1dir/file.txt", false},
		{"unlimited level 0 - l2 file", 0, "l1dir/l2dir/file.txt", false},
		{"unlimited level 0 - l3 file", 0, "l1dir/l2dir/l3dir/file.txt", false},

		{"level 1 - root", 1, "", false},
		{"level 1 - root file", 1, "root.txt", false},
		{"level 1 - l1 dir", 1, "l1dir", true},
		{"level 1 - l1 file", 1, "l1dir/file.txt", true},
		{"level 1 - l2 file", 1, "l1dir/l2dir/file.txt", true},

		{"level 2 - root", 2, "", false},
		{"level 2 - root file", 2, "root.txt", false},
		{"level 2 - l1 dir", 2, "l1dir", false},
		{"level 2 - l1 file", 2, "l1dir/file.txt", false},
		{"level 2 - l2 dir", 2, "l1dir/l2dir", true},
		{"level 2 - l2 file", 2, "l1dir/l2dir/file.txt", true},

		{"level 3 - root", 3, "", false},
		{"level 3 - root file", 3, "root.txt", false},
		{"level 3 - l1 dir", 3, "l1dir", false},
		{"level 3 - l1 file", 3, "l1dir/file.txt", false},
		{"level 3 - l2 dir", 3, "l1dir/l2dir", false},
		{"level 3 - l2 file", 3, "l1dir/l2dir/file.txt", false},
		{"level 3 - l3 dir", 3, "l1dir/l2dir/l3dir", true},
		{"level 3 - l3 file", 3, "l1dir/l2dir/l3dir/file.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewVfsScoped(tmpDir, vfs.OSFS, WithAllowLevel(tt.maxLevel))
			if err != nil {
				t.Fatal(err)
			}
			_, err = v.Open(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

func TestVfsScoped_AllowLevel_ReadDir(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(tmpDir, "l1dir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, "l1dir/l2dir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "l1dir/file.txt"), []byte("l1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "l1dir/l2dir/file.txt"), []byte("l2"), 0644); err != nil {
		t.Fatal(err)
	}

	v, _ := NewVfsScoped(tmpDir, vfs.OSFS, WithAllowLevel(1))
	entries, err := v.ReadDir("")
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0)
	for _, e := range entries {
		names = append(names, e.Name())
	}
	if len(names) != 1 || names[0] != "root.txt" {
		t.Errorf("ReadDir level 1: got %v, want [root.txt]", names)
	}

	v2, _ := NewVfsScoped(tmpDir, vfs.OSFS, WithAllowLevel(2))
	entries2, err := v2.ReadDir("")
	if err != nil {
		t.Fatal(err)
	}
	names2 := make([]string, 0)
	for _, e := range entries2 {
		names2 = append(names2, e.Name())
	}
	if len(names2) != 2 {
		t.Errorf("ReadDir level 2: got %v, want [l1dir root.txt]", names2)
	}
}

func TestVfsScoped_AllowLevel_Glob(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "other.go"), []byte("code"), 0644); err != nil {
		t.Fatal(err)
	}

	v, err := NewVfsScoped(tmpDir, vfs.OSFS, WithAllowLevel(1))
	if err != nil {
		t.Fatal(err)
	}

	// Glob with metacharacters must not fail when maxLevel > 0
	matches, err := v.Glob("*.txt")
	if err != nil {
		t.Fatalf("Glob(*.txt) with maxLevel=1 returned error: %v", err)
	}
	if len(matches) != 1 || matches[0] != "root.txt" {
		t.Errorf("Glob(*.txt) = %v, want [root.txt]", matches)
	}

	// Glob with exact path at allowed level
	matches, err = v.Glob("root.txt")
	if err != nil {
		t.Fatalf("Glob(root.txt) with maxLevel=1 returned error: %v", err)
	}
	if len(matches) != 1 || matches[0] != "root.txt" {
		t.Errorf("Glob(root.txt) = %v, want [root.txt]", matches)
	}

	// Glob pattern matching nothing — must not crash
	matches, err = v.Glob("*.nonexistent")
	if err != nil {
		t.Fatalf("Glob(*.nonexistent) with maxLevel=1 returned error: %v", err)
	}
	if len(matches) != 0 {
		t.Errorf("Glob(*.nonexistent) = %v, want []", matches)
	}
}

func TestVfsScoped_AllowLevel_Walk(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(tmpDir, "l1dir/l2dir/l3dir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "l1dir/file.txt"), []byte("l1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "l1dir/l2dir/file.txt"), []byte("l2"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "l1dir/l2dir/l3dir/file.txt"), []byte("l3"), 0644); err != nil {
		t.Fatal(err)
	}

	var found []string
	v, _ := NewVfsScoped(tmpDir, vfs.OSFS, WithAllowLevel(2))
	if err := WalkScoped(v, ".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		found = append(found, path)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	l3Found := false
	for _, f := range found {
		if strings.Contains(f, "l3dir") {
			l3Found = true
			break
		}
	}
	if l3Found {
		t.Errorf("Walk should not have found l3dir contents with level 2, found: %v", found)
	}
}

func TestVfsScoped_SymlinkOldnameEscape(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	outsideDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(outsideDir, "secret.txt"), []byte("secret"), 0644); err != nil {
		t.Fatal(err)
	}

	err = v.Symlink(filepath.Join(outsideDir, "secret.txt"), "link_to_outside")
	if err == nil {
		t.Fatal("BUG: Symlink accepted oldname outside scope without error")
	}
}

func TestVfsScoped_SymlinkAbsoluteOldnameEscape(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	err = v.Symlink("/etc/passwd", "link_to_etc")
	if err == nil {
		t.Fatal("BUG: Symlink accepted absolute oldname outside scope without error")
	}
}

func TestVfsScoped_RawPathEscape(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	outsideDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(outsideDir, "secret.txt"), []byte("secret"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(tmpDir, "dir"), 0755); err != nil {
		t.Fatal(err)
	}
	err = os.Symlink(outsideDir, filepath.Join(tmpDir, "dir/link"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = v.RawPath("dir/link")
	if err == nil {
		t.Fatal("BUG: RawPath should reject symlink with target outside scope")
	}
}

func TestVfsScoped_HasRootPrefixEdgeCase(t *testing.T) {
	otherDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(otherDir, "foobar"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(otherDir, "foobar", "file.txt"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	v, err := NewVfsScoped(filepath.Join(otherDir, "foo"), vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	_, err = v.Open("bar/file.txt")
	if err == nil {
		t.Error("Open should fail: path /foo/bar is outside root /foo but hasRootPrefix does string prefix match")
	}
}

func TestVfsScoped_WalkSymlinkToOutside(t *testing.T) {
	tmpDir := t.TempDir()
	v, err := NewVfsScoped(tmpDir, vfs.OSFS)
	if err != nil {
		t.Fatal(err)
	}

	outsideDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(outsideDir, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outsideDir, "subdir", "file.txt"), []byte("secret"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(tmpDir, "dir"), 0755); err != nil {
		t.Fatal(err)
	}
	err = os.Symlink(outsideDir, filepath.Join(tmpDir, "dir/escape"))
	if err != nil {
		t.Fatal(err)
	}

	var visited []string
	if err := WalkScoped(v, ".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		visited = append(visited, path)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	for _, p := range visited {
		if strings.Contains(p, "..") || strings.Contains(p, outsideDir) {
			t.Errorf("Walk leaked path outside root: %q", p)
		}
	}
}
