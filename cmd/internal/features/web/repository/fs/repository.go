package web_fs_repository

type WebRepository struct {
	publicDir string
}

func NewWebRepository(publicDir string) *WebRepository {
	return &WebRepository{publicDir: publicDir}
}
