package web_fs_repository

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
)

func (r *WebRepository) GetFile(ctx context.Context, name string) (domain.WebFile, error) {
	if err := ctx.Err(); err != nil {
		return domain.WebFile{}, err
	}
	// OpenRoot prevents paths and symbolic links from escaping publicDir.
	root, err := os.OpenRoot(r.publicDir)
	if err != nil {
		return domain.WebFile{}, fmt.Errorf("open public directory: %w", err)
	}
	defer root.Close()
	file, err := root.Open(name)
	if err != nil {
		return domain.WebFile{}, fmt.Errorf("open public file: %w", err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return domain.WebFile{}, fmt.Errorf("stat public file: %w", err)
	}
	// Directory listings and special files are not public assets.
	if !info.Mode().IsRegular() {
		return domain.WebFile{}, fmt.Errorf("not a regular file: %w", fs.ErrNotExist)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return domain.WebFile{}, fmt.Errorf("read public file: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return domain.WebFile{}, err
	}
	return domain.WebFile{Name: info.Name(), Content: content, ModTime: info.ModTime()}, nil
}
