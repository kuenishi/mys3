package mys3

import (
	//"encoding/xml"
	"fmt"
	"strings"
	"os"
)

func GetObject (a S3Account, bucket, object string) {
	req := NewRequest(a, "GET", bucket, object)
	filename := basename(object)
	f,e := os.Create(filename)
	if e != nil {
		panic("can't open file")
	}
	fmt.Printf("%v %v\n", f, e)
	n := req.SendAndWriteFile(f)
	f.Close()
	fmt.Printf("%d bytes written to %s\n", n, filename)
}

func basename (object string) string {
	a := strings.Split(object, "/")
	return a[len(a)-1]
}