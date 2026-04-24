package vfs4

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/twpayne/go-vfs"
)

// VfsScoped is a virtual filesystem wrapper that restricts all operations to a specific
// directory tree (the "root") and optionally limits the maximum depth of accessible paths.
//
// # Path Scoping
//
// Every path passed to VfsScoped operations is interpreted as relative to root.
// Attempts to access paths outside root (via ".." or absolute paths not starting
// with root) are rejected with os.ErrNotExist.
//
// # Absolute Path Handling
//
// Absolute paths are accepted if they start with root; the root prefix is stripped
// and the remaining path is treated as relative. Absolute paths outside root are
// rejected.
//
// # Level Restrictions
//
// By default (maxLevel=0), there is no limit on path depth. Use WithAllowLevel(n)
// to restrict access to paths at or below level n.
//
// # Level Semantics
//
// Level represents the depth of an item in the directory tree:
//   - Level 0: root itself ("", ".")
//   - Level 1: files directly in root (root.txt)
//   - Level 2: directories at depth 1 under root AND files in those directories
//     (folder-a/, folder-a/file.txt)
//   - Level 3: directories at depth 2 AND their contents
//     (folder-a/grand-A/, folder-a/grand-A/granda.txt)
//
// Directories count as one level deeper than files at the same physical depth.
//
// # Directory Visibility
//
// When a directory is hidden due to level restrictions, all its contents are
// also hidden (they cannot be listed, traversed, or accessed).
//
// # Examples
//
// Basic scoped filesystem (no depth limit):
//
//	v, err := NewVfsScoped("/home/user/project", vfs.OSFS)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// Open("file.txt") reads /home/user/project/file.txt
//	// Open("../secret") returns os.ErrNotExist
//
// Scoped to root-level files only (no directories visible):
//
//	v, err := NewVfsScoped("/home/user/project", vfs.OSFS, WithAllowLevel(1))
//
// Scoped to root + one level of subdirectories:
//
//	v, err := NewVfsScoped("/home/user/project", vfs.OSFS, WithAllowLevel(2))
type OptionFunc func(*VfsScoped)

func WithAllowLevel(n int) OptionFunc {
	return func(v *VfsScoped) {
		v.maxLevel = n
	}
}

// VfsScoped restricts filesystem operations to a specific root directory and
// optionally limits the maximum depth of accessible paths.
type VfsScoped struct {
	root     string
	maxLevel int
	inner    vfs.FS
}

// NewVfsScoped creates a new scoped filesystem with the given root directory.
// The root must be an absolute path. If maxLevel is 0 (default), there is no
// depth limit. Use WithAllowLevel to restrict path depth.
//
// Examples:
//
// Basic scoped filesystem (no depth limit):
//
//	v, err := NewVfsScoped("/home/user/project", vfs.OSFS)
//
// Scoped to root-level files only (no directories visible):
//
//	v, err := NewVfsScoped("/home/user/project", vfs.OSFS, WithAllowLevel(1))
//
// Scoped to root + one level of subdirectories:
//
//	v, err := NewVfsScoped("/home/user/project", vfs.OSFS, WithAllowLevel(2))
func NewVfsScoped(root string, inner vfs.FS, opts ...OptionFunc) (*VfsScoped, error) {
	if !filepath.IsAbs(root) {
		return nil, errors.New("root must be an absolute path")
	}
	root = filepath.Clean(root)
	v := &VfsScoped{root: root, maxLevel: 0, inner: inner}
	for _, opt := range opts {
		opt(v)
	}
	return v, nil
}

func (s *VfsScoped) levelFromPath(path string, isDir bool) int {
	rel, err := filepath.Rel(s.root, filepath.Clean(path))
	if err != nil || rel == "." {
		return 0
	}
	slashes := strings.Count(rel, "/")
	if isDir {
		return slashes + 2
	}
	return slashes + 1
}

func (s *VfsScoped) scopePath(name string) (string, error) {
	if filepath.IsAbs(name) {
		if !hasRootPrefix(name, s.root) {
			return "", os.ErrNotExist
		}
		name = name[len(s.root):]
		if name == "" {
			name = "."
		}
	}
	scoped := filepath.Join(s.root, name)
	scoped = filepath.Clean(scoped)
	if !hasRootPrefix(scoped, s.root) {
		return "", os.ErrNotExist
	}
	if s.maxLevel > 0 {
		info, err := s.inner.Lstat(scoped)
		if err != nil {
			return "", err
		}
		if s.levelFromPath(scoped, info.IsDir()) > s.maxLevel {
			return "", os.ErrNotExist
		}
	}
	return scoped, nil
}

func hasRootPrefix(path, root string) bool {
	path = filepath.Clean(path)
	root = filepath.Clean(root)
	if path == root {
		return true
	}
	if len(path) > len(root) && path[:len(root)] == root {
		return path[len(root)] == '/'
	}
	return false
}

func (s *VfsScoped) Chmod(name string, mode os.FileMode) error {
	scoped, err := s.scopePath(name)
	if err != nil {
		return err
	}
	return s.inner.Chmod(scoped, mode)
}

func (s *VfsScoped) Chown(name string, uid, gid int) error {
	scoped, err := s.scopePath(name)
	if err != nil {
		return err
	}
	return s.inner.Chown(scoped, uid, gid)
}

func (s *VfsScoped) Chtimes(name string, atime, mtime time.Time) error {
	scoped, err := s.scopePath(name)
	if err != nil {
		return err
	}
	return s.inner.Chtimes(scoped, atime, mtime)
}

func (s *VfsScoped) Create(name string) (*os.File, error) {
	scoped, err := s.scopePath(name)
	if err != nil {
		return nil, err
	}
	return s.inner.Create(scoped)
}

func (s *VfsScoped) Glob(pattern string) ([]string, error) {
	scoped := filepath.Join(s.root, pattern)
	scoped = filepath.Clean(scoped)
	if !hasRootPrefix(scoped, s.root) {
		return nil, os.ErrNotExist
	}
	matches, err := s.inner.Glob(scoped)
	if err != nil {
		return nil, err
	}
	var filtered []string
	for _, match := range matches {
		if hasRootPrefix(match, s.root) {
			info, err := s.inner.Lstat(match)
			if err != nil {
				continue
			}
			if s.maxLevel > 0 && s.levelFromPath(match, info.IsDir()) > s.maxLevel {
				continue
			}
			rel, err := filepath.Rel(s.root, match)
			if err != nil {
				continue
			}
			filtered = append(filtered, rel)
		}
	}
	return filtered, nil
}

func (s *VfsScoped) Lchown(name string, uid, gid int) error {
	scoped, err := s.scopePath(name)
	if err != nil {
		return err
	}
	return s.inner.Lchown(scoped, uid, gid)
}

func (s *VfsScoped) Mkdir(name string, perm os.FileMode) error {
	scoped, err := s.scopePath(name)
	if err != nil {
		return err
	}
	return s.inner.Mkdir(scoped, perm)
}

func (s *VfsScoped) Lstat(name string) (os.FileInfo, error) {
	scoped, err := s.scopePath(name)
	if err != nil {
		return nil, err
	}
	return s.inner.Lstat(scoped)
}

func (s *VfsScoped) Open(name string) (*os.File, error) {
	scoped, err := s.scopePath(name)
	if err != nil {
		return nil, err
	}
	info, err := s.inner.Lstat(scoped)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := s.inner.Readlink(scoped)
		if err != nil {
			return nil, err
		}
		var resolved string
		if filepath.IsAbs(target) {
			resolved = target
		} else {
			resolved = filepath.Join(filepath.Dir(scoped), target)
		}
		resolved = filepath.Clean(resolved)
		if !hasRootPrefix(resolved, s.root) {
			return nil, os.ErrNotExist
		}
	}
	return s.inner.Open(scoped)
}

func (s *VfsScoped) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	scoped, err := s.scopePath(name)
	if err != nil {
		return nil, err
	}
	return s.inner.OpenFile(scoped, flag, perm)
}

func (s *VfsScoped) PathSeparator() rune {
	return '/'
}

func (s *VfsScoped) RawPath(name string) (string, error) {
	scoped, err := s.scopePath(name)
	if err != nil {
		return "", err
	}
	info, err := s.inner.Lstat(scoped)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := s.inner.Readlink(scoped)
		if err != nil {
			return "", err
		}
		var resolved string
		if filepath.IsAbs(target) {
			resolved = target
		} else {
			resolved = filepath.Join(filepath.Dir(scoped), target)
		}
		resolved = filepath.Clean(resolved)
		if !hasRootPrefix(resolved, s.root) {
			return "", os.ErrNotExist
		}
	}
	return s.inner.RawPath(scoped)
}

func (s *VfsScoped) ReadDir(dirname string) ([]os.FileInfo, error) {
	scoped, err := s.scopePath(dirname)
	if err != nil {
		return nil, err
	}
	entries, err := s.inner.ReadDir(scoped)
	if err != nil {
		return nil, err
	}
	if s.maxLevel > 0 {
		var filtered []os.FileInfo
		for _, entry := range entries {
			entryPath := filepath.Join(scoped, entry.Name())
			if s.levelFromPath(entryPath, entry.IsDir()) > s.maxLevel {
				continue
			}
			filtered = append(filtered, entry)
		}
		return filtered, nil
	}
	return entries, nil
}

func (s *VfsScoped) ReadFile(filename string) ([]byte, error) {
	scoped, err := s.scopePath(filename)
	if err != nil {
		return nil, err
	}
	return s.inner.ReadFile(scoped)
}

func (s *VfsScoped) Readlink(name string) (string, error) {
	scoped, err := s.scopePath(name)
	if err != nil {
		return "", err
	}
	target, err := s.inner.Readlink(scoped)
	if err != nil {
		return "", err
	}
	if filepath.IsAbs(target) {
		if !hasRootPrefix(target, s.root) {
			return "", os.ErrNotExist
		}
		rel, err := filepath.Rel(s.root, target)
		if err != nil {
			return "", os.ErrNotExist
		}
		return rel, nil
	}
	resolved := filepath.Join(filepath.Dir(scoped), target)
	resolved = filepath.Clean(resolved)
	if !hasRootPrefix(resolved, s.root) {
		return "", os.ErrNotExist
	}
	rel, _ := filepath.Rel(s.root, resolved)
	return rel, nil
}

func (s *VfsScoped) Remove(name string) error {
	scoped, err := s.scopePath(name)
	if err != nil {
		return err
	}
	return s.inner.Remove(scoped)
}

func (s *VfsScoped) RemoveAll(name string) error {
	scoped, err := s.scopePath(name)
	if err != nil {
		return err
	}
	return s.inner.RemoveAll(scoped)
}

func (s *VfsScoped) Rename(oldpath, newpath string) error {
	scopedOld, err := s.scopePath(oldpath)
	if err != nil {
		return err
	}
	scopedNew, err := s.scopePath(newpath)
	if err != nil {
		return err
	}
	return s.inner.Rename(scopedOld, scopedNew)
}

func (s *VfsScoped) Stat(name string) (os.FileInfo, error) {
	scoped, err := s.scopePath(name)
	if err != nil {
		return nil, err
	}
	return s.inner.Stat(scoped)
}

func (s *VfsScoped) Symlink(oldname, newname string) error {
	scopedNew, err := s.scopePath(newname)
	if err != nil {
		return err
	}
	var resolved string
	if filepath.IsAbs(oldname) {
		resolved = oldname
	} else {
		resolved = filepath.Join(filepath.Dir(scopedNew), oldname)
	}
	resolved = filepath.Clean(resolved)
	if !hasRootPrefix(resolved, s.root) {
		return os.ErrNotExist
	}
	return s.inner.Symlink(oldname, scopedNew)
}

func (s *VfsScoped) Truncate(name string, size int64) error {
	scoped, err := s.scopePath(name)
	if err != nil {
		return err
	}
	return s.inner.Truncate(scoped, size)
}

func (s *VfsScoped) WriteFile(filename string, data []byte, perm os.FileMode) error {
	scoped, err := s.scopePath(filename)
	if err != nil {
		return err
	}
	return s.inner.WriteFile(scoped, data, perm)
}

func MkdirAllScoped(v *VfsScoped, name string, perm os.FileMode) error {
	scoped, err := v.scopePath(name)
	if err != nil {
		return err
	}
	return vfs.MkdirAll(v.inner, scoped, perm)
}

func ContainsScoped(v *VfsScoped, dir, name string) (bool, error) {
	if filepath.IsAbs(name) {
		return false, nil
	}
	scopedDir, err := v.scopePath(dir)
	if err != nil {
		return false, err
	}
	scopedTarget, err := v.scopePath(filepath.Join(dir, name))
	if err != nil {
		return false, nil
	}
	_, err = v.inner.Lstat(scopedTarget)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	targetRel, err := filepath.Rel(scopedDir, scopedTarget)
	if err != nil {
		return false, err
	}
	if targetRel == "." {
		return true, nil
	}
	if !strings.HasPrefix(targetRel, "..") {
		return true, nil
	}
	return false, nil
}

func WalkScoped(v *VfsScoped, root string, fn filepath.WalkFunc) error {
	scoped, err := v.scopePath(root)
	if err != nil {
		return err
	}
	return vfs.Walk(v.inner, scoped, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fn(path, info, err)
		}
		if v.maxLevel > 0 && v.levelFromPath(path, info.IsDir()) > v.maxLevel {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(v.root, path)
		if err != nil {
			return fn(path, info, err)
		}
		return fn(rel, info, nil)
	})
}
