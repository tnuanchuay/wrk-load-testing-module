package view_controller

import (
	"ahlt/model"
	"github.com/jinzhu/gorm"
)

type(
	DiffSelect struct{
		OriginId		uint
		SameTestsetJobs		[]model.Job
	}
)

func (DiffSelect) GetViewControl(db *gorm.DB, id uint) *DiffSelect{
	var diffSelect DiffSelect
	var firstJob model.Job
	diffSelect.OriginId = id
	db.Where("id = ?", id).First(&firstJob)
	db.Where("testset = ?", firstJob.Testset).Order("created_at desc").Find(&diffSelect.SameTestsetJobs)
	return &diffSelect
}
