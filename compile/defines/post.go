package defines

import (
	"time"

	"github.com/teawithsand/handmd/util/fsal"
)

type Post struct {
	PostMetadata PostMetadata
	Content      PostContent

	Dir string  // directory name discovered by loader
	FS  fsal.FS // FS from directory discovered
}

type PostMetadata struct {
	UnstableID string `json:"unstableId"`

	Path string `json:"path"`

	Slug string `json:"slug"`

	Title string   `json:"title"`
	Tags  []string `json:"tags"`

	CreatedAt    time.Time  `json:"createdAt"`
	LastEditedAt *time.Time `json:"lastEditedAt"`

	DirName string `json:"dirName"`
}
