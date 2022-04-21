package main

import (
	"fmt"
	"time"
)

func main() {
	s := "2022-04-20T16:36:24Z"
	//s := time.Now().Format("2006-01-02 15:04:05")
	//fmt.Println(t, reflect.TypeOf(t))

	t, err := time.ParseInLocation("2006-01-02T15:04:05+8:00", s, time.Local)
	if err != nil {
		t, err = time.ParseInLocation("2006-01-02T15:04:05Z", s, time.Local)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(t)
	fmt.Println(t.Format("2006-01-02 15:04:05"))
}
