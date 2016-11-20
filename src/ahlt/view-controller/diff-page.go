package view_controller

import (
	"ahlt/model"
	"github.com/jinzhu/gorm"
	"ahlt/compare"
	"fmt"
)

type(
	DiffPage struct{
		Job1		model.Job
		Job2		model.Job
		Benchmark1	compare.BenchmarkResult
		Benchmark2	compare.BenchmarkResult
		XAxis		[]string
	}

)

func (DiffPage) GetViewControl(db *gorm.DB, id1, id2 uint) *DiffPage{
	var job1 model.Job
	var job2 model.Job

	db.Where("id = ?", id1).Find(&job1)
	db.Where("id = ?", id2).Find(&job2)

	db.Model(&job1).Related(&job1.WrkResult)
	db.Model(&job2).Related(&job2.WrkResult)

	var diffPage DiffPage
	diffPage.Job1 = job1
	diffPage.Job2 = job2

	diffPage.Benchmark1.FromWrkResultToJobData(job1)
	diffPage.Benchmark2.FromWrkResultToJobData(job2)

	var testset model.Testset
	db.Where("id = ?", job1.Testset).Find(&testset)
	db.Model(&testset).Related(&testset.Testcase)

	for _, testcase := range testset.Testcase{
		diffPage.XAxis = append(diffPage.XAxis, testcase.Connection)
	}

	return &diffPage
}
