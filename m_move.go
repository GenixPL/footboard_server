package main

type Move struct {
	PerformedBy Client `json:"performedBy"`

	SX int `json:"sX"`
	SY int `json:"sY"`
	EX int `json:"eX"`
	EY int `json:"eY"`
}
