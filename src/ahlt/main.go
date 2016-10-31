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
	"runtime"
	"ahlt/unit/si"
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

	db.Find(&model.Job{}, "complete = ?", false).Update("exit_interrupt", true)

	var jobProgress map[uint]float64 = make(map[uint]float64)
	var wrkChannel = make(chan *model.Job, 100)

	iris.Config.IsDevelopment = true

	iris.Static("/assets", "./static/assets", 1)

	iris.Get("/", func(ctx *iris.Context){
		ctx.Redirect("/run")
	})

	iris.Get("/run", func(ctx *iris.Context){
		home := view_controller.Home{}.GetViewControl(db)
		ctx.Render("home.html", home)
	})

	iris.Get("/result", func(ctx *iris.Context){
		resultViewControl := view_controller.Result{}.GetViewControl(db)
		ctx.Render("job.html", resultViewControl)
	})

	iris.Get("/result/:id", func(ctx *iris.Context){
		stringId := ctx.Param("id")
		intId, _ := strconv.Atoi(stringId)
		resultJobViewControl := view_controller.JobResult{}.GetJobViewControl(db, uint(intId))
		if resultJobViewControl == nil{
			ctx.Redirect("/result")
		}else{
			ctx.Render("result.html", resultJobViewControl)
		}
	})

	iris.Get("/result/:id/del", func(ctx *iris.Context){
		id := ctx.Param("id")
		db.Delete(&model.Job{}, "id = ?", id)
		ctx.Redirect("/result")
	})

	iris.Get("/testset", func(ctx *iris.Context){
		testsetPage := view_controller.TestsetPage{}.GetPageViewControl(db)
		ctx.Render("testset.html", testsetPage)
	})

	iris.Get("/testset/:id", func(ctx *iris.Context){
		stringId := ctx.Param("id")

		if stringId == "new"{
			ctx.Render("testset-new.html", nil)
			return;
		}

		id, err := strconv.Atoi(stringId)
		if err != nil{
			ctx.Redirect("/testset")
			return;
		}

		var testset model.Testset
		db.Find(&testset, "id = ?", id).Related(&testset.Testcase)

		var vc = view_controller.EditTestsetPage{Testset:testset}

		ctx.Render("testset-edit.html", vc)
	})

	iris.Post("/testset", func(ctx *iris.Context){
		name := string(ctx.FormValue("name"))
		cpu := runtime.NumCPU()
		duration := string(ctx.FormValue("duration"))
		connections := ctx.FormValues("connection")

		var testset model.Testset
		testset.Name = name

		for _, connection := range connections{
			var testcase model.Testcase
			testcase.Duration = duration
			testcase.Connection = connection
			floatConnection, _ := si.SIToFloat(connection)
			if int(floatConnection) <= cpu{
				testcase.Thread = strconv.Itoa(int(floatConnection))
			}else{
				testcase.Thread = strconv.Itoa(cpu)
			}

			testset.Testcase = append(testset.Testcase, testcase)
		}

		db.Create(&testset)
		ctx.Redirect("/testset")
	})

	iris.Get("/testset/:id/del", func(ctx *iris.Context){
		id := string(ctx.Param("id"))
		intId, _ := strconv.Atoi(id)
		db.Delete(&model.Testset{}, "id = ?", uint(intId)).Related(&model.Testcase{})
		ctx.Redirect("/testset")
	})

	iris.Post("/testset/:id/edit", func(ctx *iris.Context){
		name := string(ctx.FormValue("name"))
		cpu := runtime.NumCPU()
		duration := string(ctx.FormValue("duration"))
		connections := ctx.FormValues("connection")
		id := string(ctx.FormValue("id"))

		var testset model.Testset
		intId, _ := strconv.Atoi(id)
		uIntId := uint(intId)
		db.Find(&testset, "id = ?", uIntId)
		testset.Name = name

		db.Delete(&model.Testcase{}, "testset_id = ?", uIntId)

		for _, connection := range connections{
			var testcase model.Testcase
			testcase.Duration = duration
			testcase.Connection = connection
			floatConnection, _ := si.SIToFloat(connection)
			if int(floatConnection) <= cpu{
				testcase.Thread = strconv.Itoa(int(floatConnection))
			}else{
				testcase.Thread = strconv.Itoa(cpu)
			}

			testset.Testcase = append(testset.Testcase, testcase)
		}

		db.Save(&testset)
		ctx.Redirect("/testset")
	})

	iris.Get("/rerun/:id", func(ctx *iris.Context){
		sId := string(ctx.Param("id"))
		var job model.Job
		db.Find(&job, "id = ?", sId)

		job.ID = 0
		job.ExitInterrupt = false
		job.Complete = false

		db.Create(&job)
		
		jobProgress[job.ID] = 1;

		wrkChannel <- &job
		ctx.Redirect("/result")
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

		job.ExitInterrupt = false
		job.Complete = false

		db.Create(&job)
		jobProgress[job.ID] = 1;
		wrkChannel <- &job
		ctx.Redirect("/")
	})

	iris.Get("/realtime", func(ctx *iris.Context){
		ctx.Render("realtime.html", nil)
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
				server.BroadcastTo("real-time", "_" + strconv.Itoa(int(job.ID)),
					map[string]interface{}{
						"rx":100,
						"ok":true,
					})
			}else {
				server.BroadcastTo("real-time", "_" + strconv.Itoa(int(job.ID)),
					map[string]interface{}{
						"rx":fmt.Sprintf("%.2f",jobProgress[job.ID]),
						"ok":true,
					})
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
						ok := !job.RunWrk(testcase, "con", scriptFile, db)
						jobProgress[job.ID] = float64(i+1) / float64(len(testset.Testcase)) *100.0

						server.BroadcastTo("real-time", "_" + strconv.Itoa(int(job.ID)),
							map[string]interface{}{
								"rx":fmt.Sprintf("%.2f",jobProgress[job.ID]),
								"ok":ok,
							})


						db.Save(&job)
						time.Sleep(10 * time.Second)
					}
					job.Complete = true
					job.ExitInterrupt = false
					db.Save(&job)
					wg.Done()
				}()
			}

			wg.Wait()
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
