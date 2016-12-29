package wrk

import (
	"github.com/tspn/wrk-load-testing-module/model"
	"os/exec"
	"fmt"
	"bufio"
	"strings"
)

func Run(url, t, c, d string) model.WrkResult{
	var command *exec.Cmd

	command = exec.Command("wrk", "-t"+t, "-c"+c, "-d"+d ,url)

	//fmt.Println(command.Args)
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
	//fmt.Println(out)

	wrk := model.WrkResult{}
	wrk.SetData(url, out)
	return wrk
}
