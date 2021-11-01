package model

type Articles struct {
	ID       int
	Title    string `json:"title"`
	FullText string `json:"fulltext"`
	Price    int    `json:"prise"`
	URLFoto  string `json:"urlfoto"`
	Data     string `json:"data"`
}
