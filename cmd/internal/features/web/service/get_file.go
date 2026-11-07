package web_service

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
	core_errors "github.com/alekseishmidko/go-course/cmd/internal/core/error"
)

func (s *WebService) GetFile(ctx context.Context, urlPath string) (domain.WebFile, error) {
	name := strings.TrimPrefix(urlPath, "/")
	if name == "" {
		name = "index.html"
	}
	// Only relative asset paths are valid; never normalize traversal into a valid path.
	if !fs.ValidPath(name) || strings.Contains(name, "\\") {
		return domain.WebFile{}, fmt.Errorf("invalid public file path: %w", core_errors.ErrNotFound)
	}
	file, err := s.webRepository.GetFile(ctx, name)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrPermission) {
			return domain.WebFile{}, fmt.Errorf("public file %q: %w", name, core_errors.ErrNotFound)
		}
		return domain.WebFile{}, fmt.Errorf("get public file: %w", err)
	}
	return file, nil
}
