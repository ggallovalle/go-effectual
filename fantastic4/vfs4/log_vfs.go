package vfs4

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/twpayne/go-vfs"
)

type LogVfs struct {
	logger *slog.Logger
	level  slog.Level
	inner  vfs.FS
}

func NewLogVfs(logger *slog.Logger, level slog.Level, inner vfs.FS) *LogVfs {
	return &LogVfs{logger: logger, level: level, inner: inner}
}

func (l *LogVfs) log(op string, msg string, args ...slog.Attr) {
	args = append(args, slog.String("op", op))
	l.logger.LogAttrs(context.Background(), l.level, msg, args...)
}

func (l *LogVfs) Chmod(name string, mode os.FileMode) error {
	var err error
	if l.inner != nil {
		err = l.inner.Chmod(name, mode)
	}
	l.log("Chmod", fmt.Sprintf("chmod %o %s", mode, name), slog.String("arg.name", name), slog.Any("arg.mode", mode),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func (l *LogVfs) Chown(name string, uid, gid int) error {
	var err error
	if l.inner != nil {
		err = l.inner.Chown(name, uid, gid)
	}
	l.log("Chown", fmt.Sprintf("chown %d:%d %s", uid, gid, name), slog.String("arg.name", name), slog.Int("arg.uid", uid), slog.Int("arg.gid", gid),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func (l *LogVfs) Chtimes(name string, atime, mtime time.Time) error {
	var err error
	if l.inner != nil {
		err = l.inner.Chtimes(name, atime, mtime)
	}
	l.log("Chtimes", fmt.Sprintf("touch -t %s %s", mtime.Format("200601021504.05"), name),
		slog.String("arg.name", name),
		slog.Time("arg.atime", atime),
		slog.Time("arg.mtime", mtime),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func (l *LogVfs) Create(name string) (*os.File, error) {
	var file *os.File
	var err error
	if l.inner != nil {
		file, err = l.inner.Create(name)
	}
	l.log("Create", fmt.Sprintf("touch %s", name), slog.String("arg.name", name),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return file, err
}

func (l *LogVfs) Glob(pattern string) ([]string, error) {
	var matches []string
	var err error
	if l.inner != nil {
		matches, err = l.inner.Glob(pattern)
	}
	l.log("Glob", fmt.Sprintf("ls %s", pattern), slog.String("arg.pattern", pattern),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return matches, err
}

func (l *LogVfs) Lchown(name string, uid, gid int) error {
	var err error
	if l.inner != nil {
		err = l.inner.Lchown(name, uid, gid)
	}
	l.log("Lchown", fmt.Sprintf("lchown %d:%d %s", uid, gid, name), slog.String("arg.name", name), slog.Int("arg.uid", uid), slog.Int("arg.gid", gid),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func (l *LogVfs) Mkdir(name string, perm os.FileMode) error {
	var err error
	if l.inner != nil {
		err = l.inner.Mkdir(name, perm)
	}
	l.log("Mkdir", fmt.Sprintf("mkdir %s", name), slog.String("arg.name", name), slog.Any("arg.perm", perm),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func (l *LogVfs) Lstat(name string) (os.FileInfo, error) {
	var info os.FileInfo
	var err error
	if l.inner != nil {
		info, err = l.inner.Lstat(name)
	}
	l.log("Lstat", fmt.Sprintf("ls -l %s", name), slog.String("arg.name", name),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return info, err
}

func (l *LogVfs) Open(name string) (*os.File, error) {
	var file *os.File
	var err error
	if l.inner != nil {
		file, err = l.inner.Open(name)
	}
	l.log("Open", fmt.Sprintf("cat %s", name), slog.String("arg.name", name),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return file, err
}

func (l *LogVfs) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	var file *os.File
	var err error
	if l.inner != nil {
		file, err = l.inner.OpenFile(name, flag, perm)
	}
	l.log("OpenFile", fmt.Sprintf("touch %s", name), slog.String("arg.name", name), slog.Int("arg.flag", flag), slog.Any("arg.perm", perm),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return file, err
}

func (l *LogVfs) PathSeparator() rune {
	return '/'
}

func (l *LogVfs) RawPath(name string) (string, error) {
	var path string
	var err error
	if l.inner != nil {
		path, err = l.inner.RawPath(name)
	}
	l.log("RawPath", fmt.Sprintf("realpath %s", name), slog.String("arg.name", name),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return path, err
}

func (l *LogVfs) ReadDir(dirname string) ([]os.FileInfo, error) {
	var entries []os.FileInfo
	var err error
	if l.inner != nil {
		entries, err = l.inner.ReadDir(dirname)
	}
	l.log("ReadDir", fmt.Sprintf("ls %s/", dirname), slog.String("arg.dirname", dirname),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return entries, err
}

func (l *LogVfs) ReadFile(filename string) ([]byte, error) {
	var data []byte
	var err error
	if l.inner != nil {
		data, err = l.inner.ReadFile(filename)
	}
	l.log("ReadFile", fmt.Sprintf("cat %s", filename), slog.String("arg.filename", filename),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return data, err
}

func (l *LogVfs) Readlink(name string) (string, error) {
	var target string
	var err error
	if l.inner != nil {
		target, err = l.inner.Readlink(name)
	}
	l.log("Readlink", fmt.Sprintf("readlink %s", name), slog.String("arg.name", name),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return target, err
}

func (l *LogVfs) Remove(name string) error {
	err := l.inner.Remove(name)
	l.log("Remove", fmt.Sprintf("rm %s", name), slog.String("arg.name", name),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func (l *LogVfs) RemoveAll(name string) error {
	err := l.inner.RemoveAll(name)
	l.log("RemoveAll", fmt.Sprintf("rm -rf %s", name), slog.String("arg.name", name),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func (l *LogVfs) Rename(oldpath, newpath string) error {
	var err error
	if l.inner != nil {
		err = l.inner.Rename(oldpath, newpath)
	}
	l.log("Rename", fmt.Sprintf("mv %s %s", oldpath, newpath), slog.String("arg.oldpath", oldpath), slog.String("arg.newpath", newpath),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func (l *LogVfs) Stat(name string) (os.FileInfo, error) {
	var info os.FileInfo
	var err error
	if l.inner != nil {
		info, err = l.inner.Stat(name)
	}
	l.log("Stat", fmt.Sprintf("ls -l %s", name), slog.String("arg.name", name),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return info, err
}

func (l *LogVfs) Symlink(oldname, newname string) error {
	var err error
	if l.inner != nil {
		err = l.inner.Symlink(oldname, newname)
	}
	l.log("Symlink", fmt.Sprintf("ln -s %s %s", oldname, newname), slog.String("arg.oldname", oldname), slog.String("arg.newname", newname),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func (l *LogVfs) Truncate(name string, size int64) error {
	var err error
	if l.inner != nil {
		err = l.inner.Truncate(name, size)
	}
	l.log("Truncate", fmt.Sprintf("truncate -s %d %s", size, name), slog.String("arg.name", name), slog.Int64("arg.size", size),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func (l *LogVfs) WriteFile(filename string, data []byte, perm os.FileMode) error {
	var err error
	if l.inner != nil {
		err = l.inner.WriteFile(filename, data, perm)
	}
	l.log("WriteFile", fmt.Sprintf("tee %s", filename), slog.String("arg.filename", filename), slog.Any("arg.perm", perm), slog.Int("arg.data_len", len(data)),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func MkdirAll(v *LogVfs, name string, perm os.FileMode) error {
	var err error
	if v.inner != nil {
		err = vfs.MkdirAll(v.inner, name, perm)
	}
	v.log("MkdirAll", fmt.Sprintf("mkdir -p %s", name), slog.String("arg.name", name), slog.Any("arg.perm", perm),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}

func Contains(v *LogVfs, dir, name string) (bool, error) {
	var ok bool
	var err error
	if v.inner != nil {
		ok, err = vfs.Contains(v.inner, dir, name)
	}
	v.log("Contains", fmt.Sprintf("test -d %s && contains %s", dir, name), slog.String("arg.dir", dir), slog.String("arg.name", name),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return ok, err
}

func Walk(v *LogVfs, root string, fn filepath.WalkFunc) error {
	var err error
	if v.inner != nil {
		err = vfs.Walk(v.inner, root, fn)
	}
	v.log("Walk", fmt.Sprintf("find %s", root), slog.String("arg.root", root),
		slog.Bool("result.success", err == nil),
		slog.String("result.err", fmt.Sprintf("%v", err)),
	)
	return err
}
