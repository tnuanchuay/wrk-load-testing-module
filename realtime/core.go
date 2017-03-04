package realtime

import (
	"github.com/tspn/wrk-load-testing-module/model"
	"sync"
	"runtime"
	"fmt"
	"bufio"
	"strings"
	"os/exec"
	"strconv"
	"encoding/json"
	"gopkg.in/kataras/iris.v6/adaptors/websocket"
)

type WrkEngine struct{
	status bool
	Concurrency	int
	SamplingTime	int
	Url		string
	ResultJob	[]model.WrkResult
	wg		sync.WaitGroup
}

func (j *WrkEngine) GetState() bool{
	return j.status
}

func (j *WrkEngine) SetSamplingTime(time int){
	if time == 0{
		j.SamplingTime = 1
	}else{
		j.SamplingTime = time
	}
}

func (j *WrkEngine) SetConcurrency(concurrent int){
	if concurrent < runtime.NumCPU(){
		j.Concurrency = runtime.NumCPU()
	}else{
		j.Concurrency = concurrent
	}
}

func (j *WrkEngine) SetUrl(url string){
	j.Url = url
}

func (*WrkEngine)New()*WrkEngine{
	return &WrkEngine{}
}

func (j *WrkEngine) Stop(){
	j.status = false;
	j.SamplingTime = 0
	j.Concurrency = runtime.NumCPU()
	j.Url = ""
}

func (j *WrkEngine) Start(so websocket.Connection){
	j.status = true
	go func() {
		for j.status {
			j.wg.Add(1)
			go func() {
				var testcase model.Testcase
				testcase.Thread = strconv.Itoa(runtime.NumCPU())
				testcase.Connection = strconv.Itoa(j.Concurrency)
				testcase.Duration = strconv.Itoa(j.SamplingTime)
				result := j.RunForResult(testcase, j.Url)
				j.ResultJob = append(j.ResultJob, *result)
				j.wg.Done()
			}()
			j.wg.Wait()
			var result = j.ResultJob[len(j.ResultJob)-1]
			if result.IsError {
				j.status = false
				(so).Emit("err", "err")
			}else{
				data := map[string]interface{}{
					"status" : j.status,
					"url" : j.Url,
					"concurrecy" : j.Concurrency,
					"sampling" : j.SamplingTime,
					"rps" : result.RequestPerSec,
					"errratio" : (float64(result.Non2xx3xx) / float64(result.Requests)),
				}
				jsonData, _ := json.Marshal(data)
				(so).Emit("data", jsonData)
			}

		}
	}()
}

func (*WrkEngine) RunForResult(ts model.Testcase, url string) *model.WrkResult{
	t := ts.Thread
	c := ts.Connection
	d := ts.Duration

	var command *exec.Cmd

	command = exec.Command("wrk", "-t"+t, "-c"+c, "-d"+d, url)

	fmt.Println(command.Args)
	cmdReader, _ := command.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	var out string
	go func() {
		for scanner.Scan() {
			out = fmt.Sprintf("%s\n%s", out, scanner.Text())
			if strings.Contains(out, "Transfer"){
				break;
			}
		}
	}()
	command.Start()
	command.Wait()
	fmt.Println(out)

	wrkResult := model.WrkResult{}
	wrkResult.SetData(url, out)

	return &wrkResult
}
