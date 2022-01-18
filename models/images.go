package models

type Image struct {
	Base
	Name string `json:"name"`
	Path string `json:"path"`
}
