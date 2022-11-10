package models

type Secret struct {
	//int
	ID     string
	UserID string
	Value  interface{}
}
