package model

type PlatformTree struct {
	ID     uint      `json:"id"`
	Name   string    `json:"name"`
	Region []Regions `json:"region"`
}

type JsonTagFields struct {
	ID    int    `json:"id"`
	Field string `json:"field"`
	Name  string `json:"name"`
}
