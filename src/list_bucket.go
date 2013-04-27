package mys3

import (
	"encoding/xml"
	"fmt"
)
//<Owner><ID>145700b3de19e3d3211592a03adafc7e928f4a50f1d591594b48882a1afc0281</ID>
// <DisplayName>nJaiZ00q</DisplayName></Owner>
type Owner struct {
	ID, DisplayName  string
}
type Bucket struct {
//<Bucket><Name>resource.basho.co.jp</Name><CreationDate>2013-03-15T12:49:04.000Z</CreationDate></Bucket>
	Name, CreationDate string
	
}
type Buckets struct {
	Bucket []Bucket
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner Owner
	Buckets Buckets
}

func DoListAllMyBucketsResult(s []byte) {
	lamb := ListAllMyBucketsResult{}
	xml.Unmarshal(s, &lamb)
	fmt.Printf("%v\n", lamb)
	fmt.Printf("%v\n", lamb.Buckets.Bucket)
}
