package mys3

import (
	//"encoding/xml"
	"fmt"
	"strings"
	"os"
)

func GetObject (a S3Account, bucket, object string) {
	req := NewRequest(a, "GET", bucket, object, nil)
	filename := basename(object)
	f,e := os.Create(filename)
	if e != nil { panic("can't create file to output") }
	n := req.SendAndWriteFile(f)
	f.Close()
	fmt.Printf("%d bytes written to %s\n", n, filename)
}

func PutObject (a S3Account, bucket, path string) {
	filename := basename(path)
	f,e := os.Open(path)
	fi,_ := f.Stat()
	size := fi.Size()

	req := NewRequest(a, "PUT", bucket, "/" + filename, f)
	req.req.Header.Add("Content-Length", string(size))
	if e != nil { panic("can't open file to PUT") }
	req.Send()
	f.Close()
	fmt.Printf("%d bytes read from %s\n", size, path)
}


func basename (object string) string {
	a := strings.Split(object, "/")
	return a[len(a)-1]
}