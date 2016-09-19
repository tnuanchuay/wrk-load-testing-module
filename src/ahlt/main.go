package main

import (
	"github.com/jinzhu/gorm"
	"os"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kataras/iris"
	"ahlt/model"
	"ahlt/view-controller"
	"fmt"
	"sync"
)

func main() {
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		panic("Cannot open database")
		os.Exit(-1)
	}
	defer db.Close()

	initializeMode(db)
	initializeTestset(db)
	var wrkChannel = make(chan *model.Job, 100)

	iris.Config.IsDevelopment = true

	iris.Static("/assets", "./static/assets", 1)

	iris.Get("/", func(ctx *iris.Context){
		home := view_controller.Home{}.GetViewControl(db)
		fmt.Println(home)
		ctx.Render("home.html", home)
	})

	iris.Post("/wrk", func(ctx *iris.Context){
		name := string(ctx.FormValue("name"))
		url := string(ctx.FormValue("url"))
		method := string(ctx.FormValue("method"))
		testset := string(ctx.FormValue("testset"))
		keys := ctx.FormValues("key")
		values := ctx.FormValues("value")

		var job model.Job
		job.Name = name
		job.RequestUrl = url
		job.RequestMethod = method

		var testSet model.Testset
		db.Find(&testSet, "name = ?", testset).Related(&testSet.Testcase)
		job.Testset = testSet.ID

		job.KeyValueToLoad(keys, values)
		db.Create(&job)

		wrkChannel <- &job
	})

	go func(){
		wg := sync.WaitGroup{}
		for{
			select{
			case job := <- wrkChannel:
				wg.Add(1)
				go func(){
					var testset model.Testset
					db.Find(&testset, "id = ?", job.Testset).Related(&testset.Testcase)
					scriptFile := job.GenerateScript(job.Name)
					for _, testcase := range testset.Testcase{
						job.RunWrk(testcase, "time", scriptFile)

					}
				}()
			}
		}
	}()

	iris.Listen(":2559")
}
func initializeTestset(db *gorm.DB) {
	var t1 model.Testset

	t1.Name = "simple testset"
	db.First(&t1, "name = ?", t1.Name).Related(&t1.Testcase)

	if t1.ID == 1{
		return;
	}

	t1.Testcase = append(t1.Testcase,
		model.Testcase{Thread:"1",
			Connection:"1",
			Duration:"30s",
		})
	t1.Testcase = append(t1.Testcase,
		model.Testcase{Thread:"4",
			Connection:"10",
			Duration:"30s",
		})
	t1.Testcase = append(t1.Testcase,
		model.Testcase{Thread:"4",
			Connection:"100",
			Duration:"30s",
		})
	t1.Testcase = append(t1.Testcase,
		model.Testcase{Thread:"4",
			Connection:"1k",
			Duration:"30s",
		})
	t1.Testcase = append(t1.Testcase,
		model.Testcase{Thread:"4",
			Connection:"10k",
			Duration:"30s",
		})
	t1.Testcase = append(t1.Testcase,
		model.Testcase{Thread:"4",
			Connection:"100k",
			Duration:"30s",
		})

	db.Save(&t1)
}

func initializeMode(db *gorm.DB){
	db.AutoMigrate(&model.Job{})
	db.AutoMigrate(&model.Testcase{})
	db.AutoMigrate(&model.Testset{})
}
