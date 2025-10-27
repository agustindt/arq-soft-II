package models

// Ajusta los campos según tu entidad real. Usamos tags JSON para indexar en Solr.
type Activity struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Category  string   `json:"category"`
	Location  string   `json:"location"`
	OwnerID   string   `json:"owner_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}
