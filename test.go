package main

import "encoding/json"

type Data struct {
	Name string
}

type Wrapper struct {
	Data  any
	Error error
}

func main() {

	wrapper := Wrapper{
		Data:  Data{Name: "test"},
		Error: nil,
	}

	w, _ := json.Marshal(wrapper)

	println(string(w))

	var zver Wrapper
	json.Unmarshal(w, &zver)

	res := zver.Data.([]byte)
	println(res)
}
