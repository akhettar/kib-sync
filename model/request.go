package model

type QueryRequest struct {
	Size int `json:"size"`
	Query Query `json:"query"`
}
type Exists struct {
	Field string `json:"field"`
}
type Must struct {
	Exists Exists `json:"exists"`
}
type Bool struct {
	Must Must `json:"must"`
}
type Query struct {
	Bool Bool `json:"bool"`
}