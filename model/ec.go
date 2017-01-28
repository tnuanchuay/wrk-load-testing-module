package model

import "github.com/jinzhu/gorm"

type ECJob struct {
	gorm.Model
	Url		string
	LowNumber	int
	HighNumber	int
	Estimate	int
	TimeoutError	int
	WriteError	int
	ReadError	int
	IsDone		int
}
