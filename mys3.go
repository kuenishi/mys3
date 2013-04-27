
package main

import (
	"flag"
	"fmt"
	//"runtime"
	"time"
	"log"
	"os"
	"code.google.com/p/goconf/conf" //http://code.google.com/p/goconf/
	"net/http"
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha1"
	"io/ioutil"
)

func die(ch chan int) {
	defer func() {
		var err = recover()
		ch <- -127
		log.Println("work failed:", err)
		close(ch)
	}()
	panic("oops!\n")
	close(ch)
	ch <- 23
}

func base64_hmac(access_secret, str string) string {
	h := hmac.New(sha1.New, []byte(access_secret))
	h.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func sign(access_secret, verb, md5, ctype, date, hdrs, path string) string {
	str := fmt.Sprintf("%s\n%s\n%s\n\nx-amz-date:%s%s\n%s", verb, md5, ctype, date, hdrs, path)
	fmt.Println(str)
	return base64_hmac(access_secret, str)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
    flag.PrintDefaults()
    os.Exit(2)
}


func main() {
/*
	fmt.Printf("pa -> %d\n", (pa(add, 1))(2))

	runtime.GOMAXPROCS(5)
	var ch = make(chan int, 10)
	var i int
	i = 2
	go hello(1,i,3, ch)
	//go die(ch)
	fmt.Printf("hello, world\n")
	i = <-ch
	fmt.Printf("recvd: %d\n", i)
*/

	c,_ := conf.ReadConfigFile("/Users/kuenishi/.s3cfg")
	host_base,_ := c.GetString("default", "host_base")
	access_key,_ := c.GetString("default", "access_key")
	secret_key,_ := c.GetString("default", "secret_key")
	
	fmt.Printf("%s\n", host_base)
	fmt.Printf("%s\n", access_key)
	fmt.Printf("%s\n", secret_key)

	flag.Usage = usage
	flag.Parse()
	//args := flag.Args()
	subcmd := flag.Arg(0)
	fmt.Printf("subcmd: %s\n", subcmd)

	fmt.Printf("%v\n", flag.Args())
	
	switch subcmd {
	case "info":
		fmt.Printf("info")
		client := &http.Client{}
		//bucket := flag.Arg(1)
		url := fmt.Sprintf("http://%s", host_base)
		req,_ := http.NewRequest("GET", url, nil)
		resp,_ := client.Do(req)
		fmt.Printf("%v\n -> %v", req, resp)
	case "ls":
		fmt.Println("ls")
		client := &http.Client{}
		url := fmt.Sprintf("http://%s/", host_base)
		req,_ := http.NewRequest("GET", url, nil)
		date := time.Now().Format(time.RFC822Z)
		// http://golang.org/pkg/time/#pkg-constants

		fmt.Printf("%v\n", date)
		sign := sign(secret_key, "GET", "", "", date, "", "/")
		auth := fmt.Sprintf("AWS %s:%s", access_key, sign)
			
		req.Header.Add("Authorization", auth)
		req.Header.Add("x-amz-date", date)
		//req.Header.Add("Content-Length", "0")
		resp,_ := client.Do(req)
		body,_ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%v\n -> %v '%v'", req, resp, string(body))

	}
	//str := "GET\n\n\n\nx-amz-date:Sat, 27 Apr 2013 14:16:05 +0000\n/"
	//fmt.Printf("%s\n%s => %v\n", secret_key, str, base64_hmac(secret_key, str))
}
