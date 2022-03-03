package internal

// Tag is used to categorize articles.
type Tag struct {
	ID   int    `bun:",pk,autoincrement" json:"id"`
	Name string `json:"name"`
}
