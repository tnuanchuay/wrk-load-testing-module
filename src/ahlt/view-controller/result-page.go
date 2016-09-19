package view_controller

import (
	"ahlt/model"
	"github.com/jinzhu/gorm"
)

type(
	Result struct{
		Job		[]model.Job
		Testset		[]model.Testset
	}
)

func (Result)GetViewControl(db *gorm.DB) *Result{
	var result Result
	db.Find(&result.Job)
	db.Find(&result.Testset)
	return &result
}
