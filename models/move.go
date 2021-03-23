package models

type Move struct {
	PerformedBy string `json:"performedBy"`

	SP Point `json:"sP"`
	EP Point `json:"eP"`
}
