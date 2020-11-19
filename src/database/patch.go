package database

type Patch struct {
	ID            uint64 `db:"id" json:"id"`
	Application   string `db:"app" json:"app"`
	Name          string `db:"name" json:"name"`
	Platform      string `db:"platform" json:"platform"`
	Issuer        uint64 `db:"issuer" json:"issuer"`
	URL           string `db:"url" json:"url"`
	Hash          uint32 `db:"url_hash" json:"hash"`
	Signature     string `db:"sig_url" json:"sig"`
	SignatureHash uint32 `db:"sig_url_hash" json:"sig_hash"`
	Architecture  string `db:"arch" json:"arch"`
}
