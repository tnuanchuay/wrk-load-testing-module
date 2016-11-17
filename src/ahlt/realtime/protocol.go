package realtime


import(
	"encoding/json"
)
type Request struct{
	EngineStatus	bool	`json:"e"`
	SamplingTime	int	`json:"d"`
	Concurrency	int	`json:"c"`
	Url		string	`json:"url"`
}

func (r * Request) Parse(j string){
	json.Unmarshal([]byte(j), r)
}


