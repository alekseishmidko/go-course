package domain

import "time"

// WebFile contains a public asset without depending on HTTP or filesystem APIs.
type WebFile struct {
	Name    string
	Content []byte
	ModTime time.Time
}
