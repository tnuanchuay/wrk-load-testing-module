package realtime

import (
	"ahlt/model"
	"sync"
	"runtime"
	"fmt"
	"bufio"
	"strings"
	"os/exec"
	"strconv"
	"github.com/kataras/iris"
	"encoding/json"
)

type WrkEngine struct{
	status bool
	Concurrency	int
	samplingTime	int
	url		string
	resultJob	[]model.WrkResult
	wg		sync.WaitGroup
}

func (j *WrkEngine) GetState() bool{
	return j.status
}

func (j *WrkEngine) SetSamplingTime(time int){
	if time == 0{
		j.samplingTime = 1
	}else{
		j.samplingTime = time
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
	j.url = url
}

func (*WrkEngine)New()*WrkEngine{
	return &WrkEngine{}
}

func (j *WrkEngine) Stop(){
	j.status = false;
	j.samplingTime = 0
	j.Concurrency = runtime.NumCPU()
	j.url = ""
}

func (j *WrkEngine) Start(so iris.WebsocketConnection){
	j.status = true
	go func() {
		for j.status {
			j.wg.Add(1)
			go func() {
				var testcase model.Testcase
				testcase.Thread = strconv.Itoa(runtime.NumCPU())
				testcase.Connection = strconv.Itoa(j.Concurrency)
				testcase.Duration = strconv.Itoa(j.samplingTime)
				result := j.RunForResult(testcase, j.url)
				j.resultJob = append(j.resultJob, *result)
				j.wg.Done()
			}()
			j.wg.Wait()
			var result = j.resultJob[len(j.resultJob)-1]
			if result.IsError {
				j.status = false
				(so).Emit("err", "err")
			}else{
			data := map[string]interface{}{
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
