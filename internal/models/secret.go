package models

type Secret struct {
	//int
	ID     int
	UserID int
	Key    string
	Type   string
	Value  string
}
