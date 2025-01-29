package main

import (
	"encoding/json"

	callback "github.com/epistax1s/gomer/internal/utils"
)

func main() {

	j1 := `
	{"type":"x","page":10}
	`
	j2 := `
	{"type":"y","page12":10}
	`
	
	var data1 callback.CallbackType
	json.Unmarshal([]byte(j1), &data1)

	var data2 callback.CallbackType
	json.Unmarshal([]byte(j2), &data2)

	

/* 	x := &callback.AdminGroupNext{
		CallbackType: callback.CallbackType{
			Type: "x",
		},
		Page: 10,
	}

	jsonData1, err := json.Marshal(x)
	log.Println(err)
	if err == nil {
		log.Println(string(jsonData1))
	} */
}
