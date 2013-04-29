
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

	flag.Usage = usage
	flag.Parse()
	subcmd := flag.Arg(0)
	
	switch subcmd {
	case "get":
		mys3.GetObject(*a, flag.Arg(1), flag.Arg(2))

	case "put":
		mys3.PutObject(*a, flag.Arg(1), flag.Arg(2))

	case "info":
		if len(flag.Args()) < 2 {
			flag.Usage()
		}
		mys3.ShowPolicy(*a, flag.Arg(1))
		mys3.ShowLocation(*a, flag.Arg(1))

	case "ls":
		if len(flag.Args()) == 1 {
			mys3.ListAllMyBuckets(*a)
		}else if len(flag.Args()) == 2 {
			mys3.ListBucket(*a, flag.Arg(1))
		}

	case "mp":
		mys3.ListMultipartUploads(*a, flag.Arg(1))
		
	default:
		flag.Usage()
	}
}
