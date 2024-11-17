package models

type Priority string

const (
	Highest Priority = "highest"
	High    Priority = "high"
	Medium  Priority = "medium"
	Low     Priority = "low"
	Lowest  Priority = "lowest"
)

type Todo struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Priority `json:"priority,omitempty"`
}
