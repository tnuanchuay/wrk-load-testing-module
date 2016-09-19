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
	"time"
	"github.com/googollee/go-socket.io"
	"log"
	"strconv"
)

func main() {
	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		panic("Cannot open database")
		os.Exit(-1)
	}
	defer db.Close()

	initializeModel(db)
	initializeTestset(db)
	initializeTestset(db)

	var jobProgress map[uint]float64 = make(map[uint]float64)
	var wrkChannel = make(chan *model.Job, 100)

	iris.Config.IsDevelopment = true

	iris.Static("/assets", "./static/assets", 1)

	iris.Get("/", func(ctx *iris.Context){
		ctx.Redirect("/run")
	})

	iris.Get("/run", func(ctx *iris.Context){
		home := view_controller.Home{}.GetViewControl(db)
		fmt.Println(home)
		ctx.Render("home.html", home)
	})

	iris.Get("/result", func(ctx *iris.Context){
		resultViewControl := view_controller.Result{}.GetViewControl(db)
		ctx.Render("result.html", resultViewControl)
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
		job.ExitInterrupt = true
		db.Create(&job)
		jobProgress[job.ID] = 1;
		wrkChannel <- &job
		ctx.Redirect("/")
	})

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	server.On("connection", func(so socketio.Socket){
		so.Join("real-time")
		server.On("get-progress", func(msg string){
			i, _ := strconv.Atoi(msg)
			progress := jobProgress[uint(i)]

			var job model.Job
			db.Find(&job, "id = ?", i)

			if progress == 0{
				server.BroadcastTo("real-time", "_" + strconv.Itoa(int(i)), 100.00)
			}else {
				server.BroadcastTo("real-time", "_" + strconv.Itoa(int(i)), fmt.Sprintf("%.2f",progress))
			}
		})
	})

	server.On("error", func(so socketio.Socket, err error){
		log.Fatal(err)
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
					for i, testcase := range testset.Testcase{
						job.RunWrk(testcase, "time", scriptFile, db)
						jobProgress[job.ID] = float64(i+1) / float64(len(testset.Testcase)) *100.0
						server.BroadcastTo("real-time", "_" + strconv.Itoa(int(job.ID)), fmt.Sprintf("%.2f",jobProgress[job.ID]))
						time.Sleep(10 * time.Second)
					}
					job.ExitInterrupt = false
					db.Save(&job)
				}()
			}
		}
	}()

	iris.Handle(iris.MethodGet, "/socket.io/", iris.ToHandler(server))
	iris.Handle(iris.MethodPost, "/socket.io/", iris.ToHandler(server))

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

func initializeModel(db *gorm.DB){
	db.AutoMigrate(&model.Job{})
	db.AutoMigrate(&model.Testcase{})
	db.AutoMigrate(&model.Testset{})
	db.AutoMigrate(&model.WrkResult{})
}
