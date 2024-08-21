package model

type ListNews struct {
	Success bool `json:"success"`
	News    []NewsModel
}
