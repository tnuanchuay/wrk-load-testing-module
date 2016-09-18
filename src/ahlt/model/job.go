package model

import (
	"github.com/jinzhu/gorm"
	"os/exec"
	"fmt"
	"bufio"
	"strings"
	"io/ioutil"
	"ahlt/static"
)

type Job struct{
	gorm.Model
	Name          string
	RequestUrl    string
	RequestMethod string
	Label         string
	Testset       uint
	Load          string
}

func (r *Job) KeyValueToLoad(keys, values []string){
	keyValue := map[string]string{}
	for i, key := range keys{
		keyValue[key] = values[i]
	}

	for key, value := range keyValue{
		if len(key) > 0 {
			r.Load += key + "=" + value + "&"
		}
	}
}



func (j *Job) RunWrk(ts Testcase, label string){
	t := ts.Thread
	c := ts.Connection
	d := ts.Duration

	url := j.RequestUrl
	var command *exec.Cmd

	command = exec.Command("wrk", "-t"+t, "-c"+c, "-d"+d, "-s", fmt.Sprintf("lua/%s.lua", j.Name),url)

	fmt.Println("label", label)
	if label == "time" {
		j.Label = j.Label + "," + d
	}else{
		j.Label = j.Label + "," +c
	}

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

	//wrkResult := WrkResult{}
	//wrkResult.SetData(url, out, time)
	//
	//mongoChan <- wrkResult
}

func (j *Job) GenerateScript(filename string){
	script := ""
	script += fmt.Sprintf(static.LUA_METHOD, j.RequestMethod)
	if len(j.Load) > 0{
		script += fmt.Sprintf(static.LUA_LOAD, j.Load)
		script += fmt.Sprintf(static.LUA_CONTENTTYPE, "application/x-www-form-urlencoded")
	}
	ioutil.WriteFile(fmt.Sprintf("lua/%s.lua", filename), []byte(script), 0644)
}
