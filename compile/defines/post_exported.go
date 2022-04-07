package defines

/*

type ExportedPostMetadata struct {
	UnstableID string `json:"unstableId"`

	Path string `json:"path"`

	Slug string `json:"slug"`

	Title string   `json:"title"`
	Tags  []string `json:"tags"`

	CreatedAt    time.Time  `json:"createdAt"`
	LastEditedAt *time.Time `json:"lastEditedAt"`

	Content string `json:"content"`
}

func (epm ExportedPostMetadata) Summary() SummaryExportedPostMetadata {
	partialContent := epm.Content
	const limit = 150
	if len(partialContent) > limit {
		partialContent = partialContent[:limit]
		partialContent += "..."
	}
	return SummaryExportedPostMetadata{
		UnstableID:     epm.UnstableID,
		Path:           epm.Path,
		Title:          epm.Title,
		Tags:           epm.Tags,
		CreatedAt:      epm.CreatedAt,
		LastEditedAt:   epm.LastEditedAt,
		PartialContent: partialContent,
	}
}

type SummaryExportedPostMetadata struct {
	UnstableID string `json:"unstableId"`

	Path  string   `json:"path"`
	Title string   `json:"title"`
	Tags  []string `json:"tags"`

	CreatedAt    time.Time  `json:"createdAt"`
	LastEditedAt *time.Time `json:"lastEditedAt"`

	PartialContent string `json:"partialContent"`
}
*/
