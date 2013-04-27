
package main

import (
	"flag"
	"fmt"
	"os"
	"mys3/src"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: mys3 (ls|...) [bucketname] [path]\n")
    flag.PrintDefaults()
    os.Exit(2)
}

func main() {
	a := &mys3.S3Account{}
	a.Load(os.Getenv("HOME") + "/.s3cfg")
	
	//fmt.Printf("%s\n", a.host_base)
	//fmt.Printf("%s\n", a.access_key)
	//fmt.Printf("%s\n", secret_key)

	flag.Usage = usage
	flag.Parse()
	//args := flag.Args()
	subcmd := flag.Arg(0)
	//fmt.Printf("subcmd: %s\n", subcmd)

	//fmt.Printf("%v %d\n", flag.Args(), len(flag.Args()))
	
	switch subcmd {
	case "info":
	case "ls":

		if len(flag.Args()) == 1 {
			//url := fmt.Sprintf("http://%s/", a.Host_base)
			req := mys3.NewRequest(*a, "GET", "", "/")
			body := req.Send()
			mys3.DoListAllMyBucketsResult(body)
		}else if len(flag.Args()) == 2 {
			//url := fmt.Sprintf("http://%s.%s/", flag.Arg(1), a.Host_base)
			///fmt.Println(url)
			req := mys3.NewRequest(*a, "GET", flag.Arg(1), "/")
			body := req.Send()
			fmt.Printf("%v\n", string(body))
		}
		//fmt.Printf("%v\n -> %v '%v'", req, resp, string(body))
	default:
		flag.Usage()
	}
	//str := "GET\n\n\n\nx-amz-date:Sat, 27 Apr 2013 14:16:05 +0000\n/"
	//fmt.Printf("%s\n%s => %v\n", secret_key, str, base64_hmac(secret_key, str))
}
