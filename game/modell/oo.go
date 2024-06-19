package main

import (
	"fmt"
	"time"
)

func main() {
	time.UTC, _ = time.LoadLocation("Asia/Shanghai")

	testTime, _ := time.Parse("2006-01-02 15:04:05", "2024-06-13 00:00:00")
	fmt.Println(testTime.Unix(), time.Now().Unix())
}
