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

	JobResult struct{
		Job		model.Job
		TestSet		model.Testset
	}
)

func (Result)GetViewControl(db *gorm.DB) *Result{
	var result Result
	db.Order("created_at desc").Find(&result.Job)
	db.Find(&result.Testset)
	return &result
}


func (JobResult)GetJobViewControl(db *gorm.DB, id uint) *JobResult{
	var jobResult JobResult
	db.Find(&jobResult.Job, "id = ?", id).Related(&jobResult.Job.WrkResult)
	db.Find(&jobResult.TestSet).Related(&jobResult.TestSet.Testcase)

	if jobResult.Job.ID == 0{
		return nil
	}else{
		return &jobResult
	}
}