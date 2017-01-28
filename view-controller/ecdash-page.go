package view_controller

import (
	"github.com/tspn/wrk-load-testing-module/model"
	"github.com/jinzhu/gorm"
)

type(
	ECDashPage struct{
		Jobs		[]model.ECJob
	}
)

func (ECDashPage) GetPageViewControl(db *gorm.DB)(*ECDashPage){
	var ecdashPage ECDashPage
	db.Find(&ecdashPage.Jobs)
	return &ecdashPage
}
