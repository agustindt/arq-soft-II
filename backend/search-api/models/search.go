package models

type SearchResult[T any] struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Items []T `json:"items"`
}
