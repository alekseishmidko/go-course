package web_fs_repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestGetFileCannotEscapePublicDirectory(t *testing.T) {
	dir := t.TempDir()
	public := filepath.Join(dir, "public")
	if err := os.Mkdir(public, 0700); err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(dir, "secret.txt")
	if err := os.WriteFile(secret, []byte("private"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(secret, filepath.Join(public, "link.txt")); err != nil {
		t.Fatal(err)
	}
	repo := NewWebRepository(public)
	for _, name := range []string{"../secret.txt", "link.txt", "."} {
		t.Run(name, func(t *testing.T) {
			if _, err := repo.GetFile(context.Background(), name); err == nil {
				t.Fatalf("unexpected access to %q", name)
			}
		})
	}
}
