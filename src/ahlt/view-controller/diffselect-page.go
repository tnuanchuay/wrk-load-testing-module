package view_controller

import (
	"ahlt/model"
	"github.com/jinzhu/gorm"
)

type(
	DiffSelect struct{
		SameTestsetJobs		[]model.Job
	}
)

func (DiffSelect) GetViewControl(db *gorm.DB, id uint) *DiffSelect{
	var diffSelect DiffSelect
	var firstJob model.Job
	db.Where("id = ?", id).First(&firstJob)
	db.Where("testset = ?", firstJob.Testset).Find(&diffSelect.SameTestsetJobs)
	return &diffSelect
}
