package model

import "github.com/jinzhu/gorm"

type (
	Testcase struct{
		gorm.Model
		Thread		string
		Connection	string
		Duration	string
		TestsetID	uint
	}

	Testset struct{
		gorm.Model
		Name		string
		Testcase	[]Testcase
	}
)
