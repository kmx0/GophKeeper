package models

type Secret struct {
	//int
	ID     string
	UserID string
	Key    string
	Type   string
	Value  string
}
