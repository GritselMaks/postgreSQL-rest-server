package model

type Articles struct {
	ID       int      `json:"id,omitempty"`
	Title    string   `json:"title,omitempty"`
	FullText string   `json:"fulltext,omitempty"`
	Price    int      `json:"price,omitempty"`
	URLFoto  []string `json:"urlfoto,omitempty"`
	Date     string   `json:"date,omitempty"`
}
