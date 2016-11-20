package compare

import "ahlt/model"

type(
	BenchmarkResult struct {
		JobID			uint
		IsError			[]bool
		Url			[]string
		Duration		[]float64
		Thread			[]int
		Connection		[]int
		Latency_Avg		[]float64
		Latency_Stdev		[]float64
		Latency_Max		[]float64
		ReqPerSec_Avg		[]float64
		ReqPerSec_Stdev		[]float64
		ReqPerSec_Max		[]float64
		Requests		[]int
		RequestPerSec		[]float64
		TransferPerSec		[]float64
		TotalTransfer		[]float64
		SocketErrors_Connection	[]int
		SocketErrors_Read	[]int
		SocketErrors_Write	[]int
		SocketErrors_Timeout	[]int
		TotalSocketError	[]int
		Non2xx3xx		[]int
		SuccessRequest		[]int
		TestcaseID		uint
	}
)

func (b *BenchmarkResult) FromWrkResultToJobData(job model.Job){
	b.TestcaseID = job.WrkResult[len(job.WrkResult)-1].TestcaseID
	b.TestcaseID = job.ID
	for _, wrkResult := range job.WrkResult{
		b.IsError = append(b.IsError, wrkResult.IsError)
		b.Url = append(b.Url, wrkResult.Url)
		b.Duration = append(b.Duration, wrkResult.Duration)
		b.Thread = append(b.Thread, wrkResult.Thread)
		b.Connection = append(b.Connection, wrkResult.Connection)
		b.Latency_Avg = append(b.Latency_Avg, wrkResult.Latency_Avg)
		b.Latency_Stdev = append(b.Latency_Stdev, wrkResult.Latency_Stdev)
		b.Latency_Max = append(b.Latency_Max, wrkResult.Latency_Max)
		b.ReqPerSec_Avg = append(b.ReqPerSec_Avg, wrkResult.ReqPerSec_Avg)
		b.ReqPerSec_Stdev = append(b.ReqPerSec_Stdev, wrkResult.ReqPerSec_Stdev)
		b.ReqPerSec_Max = append(b.ReqPerSec_Max, wrkResult.ReqPerSec_Max)
		b.Requests = append(b.Requests, wrkResult.Requests)
		b.RequestPerSec = append(b.RequestPerSec, wrkResult.RequestPerSec)
		b.TransferPerSec = append(b.TransferPerSec, wrkResult.TransferPerSec)
		b.TotalTransfer = append(b.TotalTransfer, wrkResult.TotalTransfer)
		b.SocketErrors_Connection = append(b.SocketErrors_Connection, wrkResult.SocketErrors_Connection)
		b.SocketErrors_Read = append(b.SocketErrors_Read, wrkResult.SocketErrors_Read)
		b.SocketErrors_Write = append(b.SocketErrors_Write, wrkResult.SocketErrors_Write)
		b.SocketErrors_Timeout = append(b.SocketErrors_Timeout, wrkResult.SocketErrors_Timeout)
		b.Non2xx3xx= append(b.Non2xx3xx, wrkResult.Non2xx3xx)
		b.TotalSocketError = append(b.TotalSocketError, wrkResult.SocketErrors_Timeout + wrkResult.SocketErrors_Write + wrkResult.SocketErrors_Read + wrkResult.SocketErrors_Connection)
		b.SuccessRequest = append(b.SuccessRequest, wrkResult.Requests-wrkResult.Non2xx3xx)
	}
}
