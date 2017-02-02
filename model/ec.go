package model

import "github.com/jinzhu/gorm"

type ECJob struct {
	gorm.Model
	Url		string
	RequestPerSec	float64
	LowNumber	int
	HighNumber	int
	Estimate	int
	TimeoutError	int
	WriteError	int
	ReadError	int
	IsDone		int
}
