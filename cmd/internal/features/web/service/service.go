package web_service

import (
	"context"
	"github.com/alekseishmidko/go-course/cmd/internal/core/domain"
)

type WebService struct{ webRepository WebRepository }

type WebRepository interface {
	GetFile(ctx context.Context, name string) (domain.WebFile, error)
}

func NewWebService(webRepository WebRepository) *WebService {
	return &WebService{webRepository: webRepository}
}
