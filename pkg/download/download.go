package download

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Download struct {
	Url        string
	Path       string
	SectionNum int
}

func (d Download) Action() error {

	req, err := d.getRequest("HEAD")
	if err != nil{return fmt.Errorf("request can not set: %s", err)}

	res, err := http.DefaultClient.Do(req)
	if err != nil{return fmt.Errorf("client can not set: %s", err)}
	fmt.Printf("got%v\n",res.StatusCode)
	//informational responses (100– 199) Successful responses (200– 299)
	//Redirection messages (300– 399) Client error responses (400– 499)
	if res.StatusCode > 299 {
		return fmt.Errorf("can't process, response is %v", res.StatusCode)
	}
	size, err := strconv.Atoi(res.Header.Get("Content-Length"))
	if err != nil { return err}
	fmt.Println("size of file is: ",math.Round(float64(size)/1000000) , " MB")
	var section = make([][2]int,d.SectionNum)
	var eachSize = size / d.SectionNum
	fmt.Println("size of each section is: ",eachSize/1000 , " KB")
	for i := range section {
		if i == 0 {
			// starting byte of first section
			section[i][0] = 0
		} else {
			// starting byte of other sections
			section[i][0] = section[i-1][1] + 1
		}

		if i < d.SectionNum-1 {
			// ending byte of other sections
			section[i][1] = section[i][0] + eachSize
		} else {
			// ending byte of end section
			section[i][1] = size - 1
		}
	}
	log.Println(section)
	
	wg := sync.WaitGroup{}
	for i,s:= range section {
		wg.Add(1)
		go func(i int,s [2]int){
			err = d.downloadSection(i,s)
			if err != nil {panic(err)}
			defer wg.Done()	
		}(i,s)
	}
	wg.Wait()
	return d.mergeFiles(section)
}

func (d Download) getRequest(method string) (*http.Request,error){
	req, err := http.NewRequest(method,d.Url,nil)
	if err!= nil{return nil,err}
	req.Header.Set("downloadmanager","BKv0.0.01")
	return req,err	
}

func (d Download) downloadSection (i int, s [2]int)error{
	fmt.Printf("start downloading section %v\n",i)
	req, err := d.getRequest("GET")
	if err != nil {return err}
	req.Header.Set("Range",fmt.Sprintf("bytes=%v-%v",s[0],s[1]))
	res, err := http.DefaultClient.Do(req)
	if err != nil {return err}
	if res.StatusCode > 299 {
		return fmt.Errorf("can't process, response is %v", res.StatusCode)
	}
	fmt.Printf("Downloaded %v bytes for section %v\n", res.Header.Get("Content-Length"), i)
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("section-%v.tmp", i), b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (d Download) mergeFiles(sections [][2]int) error {
	f, err := os.OpenFile(d.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {return err}
	defer f.Close()
	for i := range sections {
		tmpFileName := fmt.Sprintf("section-%v.tmp", i)
		b, err := ioutil.ReadFile(tmpFileName)
		if err != nil {
			return err
		}
		_, err = f.Write(b)
		if err != nil {
			return err
		}
		err = os.Remove(tmpFileName)
		if err != nil {
			return err
		}
	}
	return nil
}


func NewDownload ()(*Download){
	var url,path string
	fmt.Println("Add the URL:")
	fmt.Scanln(&url)
	fmt.Println("Add the PATH:")
	fmt.Scanln(&path)
	d := new(Download)
	d.Url = url
	d.Path = path
	d.SectionNum = 10
	return d
}
