package console

import (
	"fmt"
	"encoding/json"
)

type ConsolePipline struct{

}

func NewConsolePipline() *ConsolePipline{
	return &ConsolePipline{}
}

func (c *ConsolePipline)ProcessData(v interface{},taskName string,processName string){
	bytes,_ := json.Marshal(v)
	fmt.Println("Pipline :",string(bytes))
}