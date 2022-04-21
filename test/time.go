package main

import (
	"fmt"
	"time"
)

func main() {
	s := "2022-04-20T16:36:24+08:00"
	//s := time.Now().Format("2006-01-02 15:04:05")
	//fmt.Println(t, reflect.TypeOf(t))

	t, err := time.ParseInLocation("2006-01-02T15:04:05+08:00", s, time.Local)
	if err != nil {
		fmt.Println(t, err)
	}

	fmt.Println(t.Format("2006-01-02 15:04:05"))

}
