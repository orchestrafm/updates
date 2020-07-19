package database

type Patch struct {
	ID          uint64 `db:"id" json:"id"`
	Application string `db:"app" json:"app"`
	Name        string `db:"name" json:"name"`
	Platform    string `db:"platform" json:"platform"`
	URL         string `db:"url" json:"url"`
}
