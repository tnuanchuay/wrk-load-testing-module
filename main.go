package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kataras/iris"
	"github.com/tspn/wrk-load-testing-module/model"
	"github.com/tspn/wrk-load-testing-module/realtime"
	"github.com/tspn/wrk-load-testing-module/unit/si"
	"github.com/tspn/wrk-load-testing-module/view-controller"
	"github.com/tspn/wrk-load-testing-module/wrk"
	"github.com/tspn/wrk-load-testing-module/ws"
	"gopkg.in/kataras/iris.v6/adaptors/websocket"
)

const (
	AFTERBURN_WAIT_TIME = 70
)

func main() {
	var realtimeWrkEngine realtime.WrkEngine
	var sockets ws.GroupSocket
	var realtimeSockets ws.GroupSocket
	var ecSockets ws.GroupSocket
	realtimeInUsed := false
	leakyBucket := make(chan int)

	app := iris.New()

	go func() {
		leakyBucket <- 1
	}()

	db, err := gorm.Open("sqlite3", "database.db")
	if err != nil {
		panic("Cannot open database")
		os.Exit(-1)
	}
	defer db.Close()

	initializeModel(db)
	initializeTestset(db)

	db.Find(&model.Job{}, "complete = ?", false).Update("exit_interrupt", true)
	db.Find(&model.ECJob{}, "is_done = ?", 0).Update("is_done", -1)

	var jobProgress map[uint]float64 = make(map[uint]float64)
	var wrkChannel = make(chan *model.Job, 100)

	//app.Config.IsDevelopment = true

	app.StaticWeb("/assets", "./static/assets")

	app.Get("/", func(ctx *iris.Context) {
		ctx.Redirect("/hello")
	})

	app.Get("/hello", func(ctx *iris.Context) {
		ctx.Render("hello.html", nil)
	})

	app.Get("/run", func(ctx *iris.Context) {
		home := view_controller.Home{}.GetViewControl(db)
		ctx.Render("home.html", home)
	})

	app.Get("/result", func(ctx *iris.Context) {
		resultViewControl := view_controller.Result{}.GetViewControl(db)
		ctx.Render("job.html", struct {
			Result *view_controller.Result
			Host   string
		}{Result: resultViewControl, Host: ctx.Host()})
	})

	app.Get("/result/:id", func(ctx *iris.Context) {
		stringId := ctx.Param("id")
		intId, _ := strconv.Atoi(stringId)
		resultJobViewControl := view_controller.JobResult{}.GetJobViewControl(db, uint(intId))
		if resultJobViewControl == nil {
			ctx.Redirect("/result")
		} else {
			ctx.Render("result.html", struct {
				Result *view_controller.JobResult
				Host   string
			}{Result: resultJobViewControl, Host: ctx.Host()})
		}
	})

	app.Get("/result/:id/diff", func(ctx *iris.Context) {
		stringId := ctx.Param("id")
		intId, _ := strconv.Atoi(stringId)
		jobDiffSelectViewControl := view_controller.DiffSelect{}.GetViewControl(db, uint(intId))
		ctx.Render("diff-select.html", jobDiffSelectViewControl)
	})

	app.Get("/result/:id/diff/:id2", func(ctx *iris.Context) {
		stringId1 := ctx.Param("id")
		stringId2 := ctx.Param("id2")
		intId1, _ := strconv.Atoi(stringId1)
		intId2, _ := strconv.Atoi(stringId2)
		jobDiffViewControl := view_controller.DiffPage{}.GetViewControl(db, uint(intId1), uint(intId2))
		ctx.Render("diff.html", jobDiffViewControl)
	})

	app.Get("/result/:id/del", func(ctx *iris.Context) {
		id := ctx.Param("id")
		db.Delete(&model.Job{}, "id = ?", id)
		ctx.Redirect("/result")
	})

	app.Get("/testset", func(ctx *iris.Context) {
		testsetPage := view_controller.TestsetPage{}.GetPageViewControl(db)
		ctx.Render("testset.html", testsetPage)
	})

	app.Get("/testset/:id", func(ctx *iris.Context) {
		stringId := ctx.Param("id")

		if stringId == "new" {
			ctx.Render("testset-new.html", nil)
			return
		}

		id, err := strconv.Atoi(stringId)
		if err != nil {
			ctx.Redirect("/testset")
			return
		}

		var testset model.Testset
		db.Find(&testset, "id = ?", id).Related(&testset.Testcase)

		var vc = view_controller.EditTestsetPage{Testset: testset}

		ctx.Render("testset-edit.html", vc)
	})

	app.Post("/testset", func(ctx *iris.Context) {
		name := string(ctx.FormValue("name"))
		cpu := runtime.NumCPU()
		duration := string(ctx.FormValue("duration"))
		connections := ctx.FormValues()["connection"]

		var testset model.Testset
		testset.Name = name

		for _, connection := range connections {
			var testcase model.Testcase
			testcase.Duration = duration
			testcase.Connection = connection
			floatConnection, _ := si.SIToFloat(connection)
			if int(floatConnection) <= cpu {
				testcase.Thread = strconv.Itoa(int(floatConnection))
			} else {
				testcase.Thread = strconv.Itoa(cpu)
			}

			testset.Testcase = append(testset.Testcase, testcase)
		}

		db.Create(&testset)
		ctx.Redirect("/testset")
	})

	app.Get("/testset/:id/del", func(ctx *iris.Context) {
		id := string(ctx.Param("id"))
		intId, _ := strconv.Atoi(id)
		db.Delete(&model.Testset{}, "id = ?", uint(intId)).Related(&model.Testcase{})
		ctx.Redirect("/testset")
	})

	app.Post("/testset/:id/edit", func(ctx *iris.Context) {
		name := string(ctx.FormValue("name"))
		cpu := runtime.NumCPU()
		duration := string(ctx.FormValue("duration"))
		connections := ctx.FormValues()["connection"]
		id := string(ctx.FormValue("id"))

		var testset model.Testset
		intId, _ := strconv.Atoi(id)
		uIntId := uint(intId)
		db.Find(&testset, "id = ?", uIntId)
		testset.Name = name

		db.Delete(&model.Testcase{}, "testset_id = ?", uIntId)

		for _, connection := range connections {
			var testcase model.Testcase
			testcase.Duration = duration
			testcase.Connection = connection
			floatConnection, _ := si.SIToFloat(connection)
			if int(floatConnection) <= cpu {
				testcase.Thread = strconv.Itoa(int(floatConnection))
			} else {
				testcase.Thread = strconv.Itoa(cpu)
			}

			testset.Testcase = append(testset.Testcase, testcase)
		}

		db.Save(&testset)
		ctx.Redirect("/testset")
	})

	app.Get("/rerun/:id", func(ctx *iris.Context) {
		sId := string(ctx.Param("id"))
		var job model.Job
		db.Find(&job, "id = ?", sId)

		job.ID = 0
		job.ExitInterrupt = false
		job.Complete = false

		db.Create(&job)

		jobProgress[job.ID] = 1

		wrkChannel <- &job
		ctx.Redirect("/result")
	})

	app.Post("/wrk", func(ctx *iris.Context) {
		name := string(ctx.FormValue("name"))
		url := string(ctx.FormValue("url"))
		method := string(ctx.FormValue("method"))
		testset := string(ctx.FormValue("testset"))
		keys := ctx.FormValues()["key"]
		values := ctx.FormValues()["value"]

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
		jobProgress[job.ID] = 1
		wrkChannel <- &job
		ctx.Redirect("/result")
	})

	app.Get("/realtime/reset", func(ctx *iris.Context) {
		realtimeWrkEngine.Stop()
		realtimeInUsed = false
		go func() {
			leakyBucket <- 1
		}()
		realtimeSockets.BroadCast("exit", map[string]interface{}{
			"exit": "exit",
		})
		ctx.Redirect("/realtime")
	})

	app.Get("/realtime", func(ctx *iris.Context) {
		if !realtimeInUsed {
			ctx.Render("realtime.html", map[string]interface{}{"Host": ctx.Host()})
		} else {
			ctx.Render("realtimebusy.html", nil)
		}

		fmt.Println("realtimeInUsed =", realtimeInUsed)
	})

	app.Get("/ec", func(ctx *iris.Context) {
		vc := view_controller.ECDashPage{}.GetPageViewControl(db)
		ctx.Render("ec.html", map[string]interface{}{
			"Result": vc,
			"Host":   ctx.Host(),
		})
	})

	app.Post("/ec/test", func(ctx *iris.Context) {
		url := string(ctx.FormValue("url"))
		ecResult := model.ECJob{}
		ecResult.Url = url
		ecResult.IsDone = 0
		db.Save(&ecResult)
		ctx.Redirect("/ec", 302)
		wg := sync.WaitGroup{}
		go func() {
			<-leakyBucket
			minCon := 0
			maxCon := 50000
			socketErrorNum := 0
			result := model.WrkResult{}
			getAnswer := false
			errNum := 0
			for !getAnswer {
				wg.Add(1)
				go func() {
					defer wg.Done()
					currentTarget := (minCon + maxCon) / 2
					result = wrk.Run(url,
						strconv.Itoa(runtime.NumCPU()),
						strconv.Itoa(currentTarget), "10s")
					time.Sleep(AFTERBURN_WAIT_TIME * time.Second)
					if result.IsError {
						fmt.Println("error", url, currentTarget)
						errNum++
						return
					}

					errNum = 0

					socketErrorNum = result.SocketErrors_Connection

					successRequestRatio := float64(result.Non2xx3xx) / float64(result.Requests) * 100.0

					fmt.Println("result", minCon, maxCon, successRequestRatio)

					if socketErrorNum != 0 {
						maxCon = currentTarget - int(float64(socketErrorNum)*0.50)
						if minCon > maxCon {
							temp := minCon
							minCon = maxCon
							maxCon = temp
						}
						return
					}

					if ((1 < successRequestRatio) && (successRequestRatio < 5)) || (currentTarget == maxCon) {
						getAnswer = true
					} else if successRequestRatio < 1 {
						minCon = currentTarget
					} else if 5 < successRequestRatio {
						maxCon = currentTarget
					}

					if float64(minCon)/float64(maxCon) > 0.99 {
						getAnswer = true
					}
				}()
				wg.Wait()
				if errNum == 3 {
					ecResult.IsDone = -1
					db.Save(&ecResult)
					return
				}
			}
			capacity := (minCon + maxCon) / 2
			fmt.Println("capacity of ", url, "=", capacity, "and can work at delivery rate", result.RequestPerSec)
			ecResult.RequestPerSec = result.RequestPerSec
			ecResult.LowNumber = minCon
			ecResult.HighNumber = maxCon
			ecResult.Estimate = capacity
			ecResult.IsDone = 1
			ecResult.TimeoutError = result.SocketErrors_Timeout
			ecResult.WriteError = result.SocketErrors_Write
			ecResult.ReadError = result.SocketErrors_Read
			db.Save(&ecResult)
			fmt.Println(len(ecSockets.Sockets))
			ecSockets.BroadCast("refresh", map[string]interface{}{"command": "need to refresh"})
			leakyBucket <- 1
		}()
	})

	wss := websocket.New(websocket.Config{
		ReadBufferSize:10000,
		WriteBufferSize:10000,
		Endpoint:"/end_point",
	})

	wss.OnConnection(func(c websocket.Connection) {

		c.On("get-progress", func(msg string) {
			i, _ := strconv.Atoi(msg)
			progress := jobProgress[uint(i)]

			var job model.Job
			db.Find(&job, "id = ?", i)
			var interfaceValue map[string]interface{}
			if progress == 0 {
				interfaceValue = map[string]interface{}{
					"progress": 100,
					"ok":       true,
				}
			} else {
				interfaceValue = map[string]interface{}{
					"progress": fmt.Sprintf("%.2f", jobProgress[job.ID]),
					"ok":       true,
				}
			}

			c.Emit("ROOM"+strconv.Itoa(int(job.ID)), interfaceValue)
			time.Sleep(10 * time.Millisecond)
		})

		c.On("regis", func(msg string) {
			switch msg {
			case "/realtime":
				realtimeSockets.Sockets = append(realtimeSockets.Sockets, &c)
			case "/result":
				sockets.Sockets = append(sockets.Sockets, &c)
			case "/ec":
				ecSockets.Sockets = append(ecSockets.Sockets, &c)
				fmt.Println(len(ecSockets.Sockets))
			}
			fmt.Println(msg)
		})

		c.On("realtime", func(msg string) {
			var request realtime.Request
			fmt.Println(msg)
			request.Parse(msg)
			if (realtimeWrkEngine.GetState() != request.EngineStatus) && (request.EngineStatus == true) {
				<-leakyBucket
				realtimeWrkEngine.SetConcurrency(request.Concurrency)
				realtimeWrkEngine.SetSamplingTime(request.SamplingTime)
				realtimeWrkEngine.SetUrl(request.Url)
				realtimeWrkEngine.Start(c)
				realtimeInUsed = true
				realtimeSockets.BroadcastAllExcept("exit", map[string]interface{}{
					"exit": "exit",
				}, c)
			} else if (realtimeWrkEngine.GetState() == request.EngineStatus) && (request.EngineStatus == true) {
				realtimeWrkEngine.SetConcurrency(request.Concurrency)
				realtimeWrkEngine.SetSamplingTime(request.SamplingTime)
			} else if request.EngineStatus == false {
				realtimeWrkEngine.Stop()
				go func() {
					leakyBucket <- 1
				}()
				realtimeInUsed = false
			}

			fmt.Println("realtimeInUsed =", realtimeInUsed)
		})
	})

	app.Adapt(wss)

	go func() {
		wg := sync.WaitGroup{}
		for {
			select {
			case job := <-wrkChannel:
				wg.Add(1)
				<-leakyBucket
				go func() {
					var testset model.Testset
					db.Find(&testset, "id = ?", job.Testset).Related(&testset.Testcase)
					scriptFile := job.GenerateScript(job.Name)
					for i, testcase := range testset.Testcase {
						ok := !job.RunWrk(testcase, "con", scriptFile, db)
						jobProgress[job.ID] = float64(i+1) / float64(len(testset.Testcase)) * 100.0
						sockets.BroadCast("ROOM"+strconv.Itoa(int(job.ID)), map[string]interface{}{
							"progress": fmt.Sprintf("%.2f", jobProgress[job.ID]),
							"ok":       ok,
						})
						db.Save(&job)
						time.Sleep(AFTERBURN_WAIT_TIME * time.Second)
					}
					job.Complete = true
					job.ExitInterrupt = false
					db.Save(&job)
					wg.Done()
					leakyBucket <- 1
				}()
			}

			wg.Wait()
		}
	}()

	app.Listen(":2559")
}

func initializeTestset(db *gorm.DB) {
	var t1 model.Testset

	t1.Name = "simple testset"
	db.First(&t1, "name = ?", t1.Name).Related(&t1.Testcase)

	if t1.ID == 1 {
		return
	}

	t1.Testcase = append(t1.Testcase,
		model.Testcase{Thread: "1",
			Connection: "1",
			Duration:   "30s",
		})
	t1.Testcase = append(t1.Testcase,
		model.Testcase{Thread: "4",
			Connection: "10",
			Duration:   "30s",
		})
	t1.Testcase = append(t1.Testcase,
		model.Testcase{Thread: "4",
			Connection: "100",
			Duration:   "30s",
		})
	t1.Testcase = append(t1.Testcase,
		model.Testcase{Thread: "4",
			Connection: "1k",
			Duration:   "30s",
		})
	t1.Testcase = append(t1.Testcase,
		model.Testcase{Thread: "4",
			Connection: "10k",
			Duration:   "30s",
		})

	db.Save(&t1)
}

func initializeModel(db *gorm.DB) {
	db.AutoMigrate(&model.Job{})
	db.AutoMigrate(&model.Testcase{})
	db.AutoMigrate(&model.Testset{})
	db.AutoMigrate(&model.WrkResult{})
	db.AutoMigrate(&model.ECJob{})
}
