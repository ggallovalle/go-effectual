# Fantastic Four

Gimmicky package name for "utils" or "common" -style packages.

## Packages

### vfs4

A logging wrapper around `github.com/twpayne/go-vfs` that implements the `vfs.FS` interface
and logs all operations at a specified log level.

#### LogVfs

```go
type LogVfs struct {
    logger *slog.Logger
    level  slog.Level
    inner  vfs.FS  // nil means no-op after logging
}
```

**Constructor:**
```go
func NewLogVfs(logger *slog.Logger, level slog.Level, inner vfs.FS) *LogVfs
```

**Logging format for each method:**
- `message` = Unix command equivalent (e.g., `rm "a/b/file.txt"`)
- `op` = method name (e.g., `Remove`)
- `arg.<name>` = method arguments

**Example for `Remove("/a/b/file.txt")`:**
```
message=rm "a/b/file.txt"
op=Remove
arg.name="a/b/file.txt"
```

**Methods implemented:** All 23 methods of `vfs.FS` interface.

**Utility functions:**
- `MkdirAll(v *LogVfs, name string, perm os.FileMode) error`
- `Contains(v *LogVfs, dir, name string) (bool, error)`
- `Walk(v *LogVfs, root string, fn filepath.WalkFunc) error`

**Usage:**

```go
// Dry-run mode (log only, no actual operations)
fs := vfs4.NewLogVfs(logger, level, nil)

// Normal mode (log + delegate to OSFS)
fs := vfs4.NewLogVfs(logger, level, vfs.OSFS)
```
