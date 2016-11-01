package realtime


import(
	"encoding/json"
)
type Request struct{
	RequestToStart	bool	`json:"e"`
	SamplingTime	int	`json:"d"`
	Concurrency	int	`json:"c"`
}

func (r * Request) FromJSON(j string){
	json.Unmarshal([]byte(j), r)
}


