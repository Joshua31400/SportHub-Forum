package requêtepost

import ()

type Post struct {
	ID        int
	UserID    int
	Title     string
	Contenu   string
	ImageURL  *string // nullable
	CreatedAt string
}
