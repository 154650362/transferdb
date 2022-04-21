package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	err := os.Setenv("LD_LIBRARY_PATH", "/Users/liuyu/soft/instantclient_19_8")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(os.Getenv("LD_LIBRARY_PATH"))

}
