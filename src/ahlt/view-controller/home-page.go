package view_controller

import (
	"ahlt/model"
	"github.com/jinzhu/gorm"
)

type (
	Home struct{
		Testset		[]model.Testset
	}
)

func (Home) GetViewControl(db *gorm.DB) *Home{
	//find All Testset
	var home Home

	db.Find(&home.Testset)

	return &home
}
