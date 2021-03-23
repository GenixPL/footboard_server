package models

type Move struct {
	PerformedBy string `json:"performedBy"`

	SX int `json:"sX"`
	SY int `json:"sY"`
	EX int `json:"eX"`
	EY int `json:"eY"`
}
