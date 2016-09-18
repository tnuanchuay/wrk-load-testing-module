package model

import "github.com/jinzhu/gorm"

type Job struct{
	gorm.Model
	Name		string
	Url		string
	Method		string
	Testset		uint
	load		string
}

func (r *Job) KeyValueToLoad(keys, values []string){
	keyValue := map[string]string{}
	for i, key := range keys{
		keyValue[key] = values[i]
	}

	for key, value := range keyValue{
		if len(key) > 0 {
			r.load += key + "=" + value + "&"
		}
	}
}