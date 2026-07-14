package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileOperationOptions controls common file operation behavior.
type FileOperationOptions struct {
	Overwrite bool
}

// EnsureDir creates dir and every missing parent directory.
func EnsureDir(dir string) error {
	if dir == "" {
		return fmt.Errorf("directory path is empty")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory %s: %w", dir, err)
	}

	return nil
}

// AtomicWriteFile writes data to path by writing a temporary file in the same
// directory and renaming it over the destination.
func AtomicWriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)

	if err := EnsureDir(dir); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary file for %s: %w", path, err)
	}

	tmpPath := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temporary file for %s: %w", path, err)
	}

	if err := tmp.Chmod(perm); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("chmod temporary file for %s: %w", path, err)
	}

	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temporary file for %s: %w", path, err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("rename temporary file to %s: %w", path, err)
	}

	cleanup = false
	return nil
}

// CopyFile copies a regular file from src to dst.
func CopyFile(src, dst string, opt FileOperationOptions) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat source file %s: %w", src, err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("source is not a regular file: %s", src)
	}

	if !opt.Overwrite {
		if _, err := os.Stat(dst); err == nil {
			return fmt.Errorf("destination already exists: %s", dst)
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat destination file %s: %w", dst, err)
		}
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file %s: %w", src, err)
	}
	defer in.Close()

	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	if !opt.Overwrite {
		flag |= os.O_EXCL
	}

	out, err := os.OpenFile(dst, flag, info.Mode().Perm())
	if err != nil {
		return fmt.Errorf("open destination file %s: %w", dst, err)
	}

	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()

	if copyErr != nil {
		return fmt.Errorf("copy %s to %s: %w", src, dst, copyErr)
	}

	if closeErr != nil {
		return fmt.Errorf("close destination file %s: %w", dst, closeErr)
	}

	return nil
}

// CopyDir recursively copies a directory tree from src to dst.
func CopyDir(src, dst string, opt FileOperationOptions) error {
	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	dstAbs, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	if isSubPath(dstAbs, srcAbs) {
		return fmt.Errorf("destination must not be inside source directory: %s", dst)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat source directory %s: %w", src, err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory: %s", src)
	}

	if !opt.Overwrite {
		if _, err := os.Stat(dst); err == nil {
			return fmt.Errorf("destination already exists: %s", dst)
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat destination directory %s: %w", dst, err)
		}
	}

	return filepath.WalkDir(src, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				return err
			}

			if err := os.MkdirAll(target, info.Mode().Perm()); err != nil {
				return fmt.Errorf("create directory %s: %w", target, err)
			}

			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		return CopyFile(path, target, opt)
	})
}

// MovePath moves a file or directory from src to dst.
func MovePath(src, dst string, opt FileOperationOptions) error {
	if !opt.Overwrite {
		if _, err := os.Stat(dst); err == nil {
			return fmt.Errorf("destination already exists: %s", dst)
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat destination %s: %w", dst, err)
		}
	}

	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}

	if opt.Overwrite {
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("remove destination %s: %w", dst, err)
		}
	}

	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat source %s: %w", src, err)
	}

	if info.IsDir() {
		if err := CopyDir(src, dst, opt); err != nil {
			return err
		}
	} else {
		if err := CopyFile(src, dst, opt); err != nil {
			return err
		}
	}

	return RemovePath(src)
}

// RemovePath removes a file or directory tree.
func RemovePath(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("remove %s: %w", path, err)
	}

	return nil
}

func isSubPath(path, parent string) bool {
	rel, err := filepath.Rel(parent, path)
	if err != nil {
		return false
	}

	return rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
