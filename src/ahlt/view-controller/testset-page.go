package view_controller

import (
	"ahlt/model"
	"github.com/jinzhu/gorm"
)

type(
	TestsetPage	struct {
		Testset		[]model.Testset
	}
)

func (TestsetPage) GetPageViewControl(db *gorm.DB)(*TestsetPage){
	var testsetPage TestsetPage
	db.Find(&testsetPage.Testset)
	for i := 0 ; i < len(testsetPage.Testset) ; i++ {
		db.Find(&testsetPage.Testset[i]).Related(&testsetPage.Testset[i].Testcase)
	}
	return &testsetPage
}