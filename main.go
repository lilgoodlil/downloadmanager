package main

import (
	"dmanager/pkg/download"
	"fmt"
	"time"
)

var d download.Download = download.Download{
	Url: "https://sorsore.com/wp-content/uploads/2020/02/kiazhameghashangtare.mp3_83079.mp3",
	Path:"c:/go/amoopoorang.mp3",
	SectionNum: 5,
}

func main() {
	startTime := time.Now()
	err := d.Action()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Download Completed in %v Seconds",time.Now().Sub(startTime).Seconds())

}