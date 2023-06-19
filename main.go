package main

import (
	"dmanager/pkg/download"
	"fmt"
	"time"
)

func main() {

	d := download.NewDownload()
	startTime := time.Now()
	err := d.Action()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Download Completed in %v Seconds",time.Now().Sub(startTime).Seconds())

}