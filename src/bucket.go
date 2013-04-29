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
func (b Bucket) pp() {
	fmt.Printf("%s\t%s\n", b.CreationDate, b.Name)
}

type Buckets struct {
	Bucket []Bucket
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Owner Owner
	Buckets Buckets
}

func ListAllMyBuckets(a S3Account) {
	req := NewRequest(a, "GET", "", "/", nil)
	body := req.Send()

	lamb := ListAllMyBucketsResult{}
	xml.Unmarshal(body, &lamb)
	//fmt.Printf("%v\n", lamb)
	for _,b := range lamb.Buckets.Bucket {
		b.pp()
	}
}

// maybe CommonPrefixes and Contents can be
// aggregated with same Interface function 'PP'
// in the far future.
type CommonPrefixes struct {
	Prefix string
}
func (cp CommonPrefixes) pp () {
	fmt.Printf("\t\t\t\tDIR\t%s\n", cp.Prefix)
}
type Contents struct {
	Key, LastModified, ETag, StorageClass string
	Size int
	Owner Owner
}
func (c Contents) pp () {
	fmt.Printf("%s\t%d\t%s\n",
		c.LastModified, c.Size, c.Key)
}
type ListBucketResult struct {
	XMLName xml.Name `xml:"ListBucketResult"`
	Name, Prefix, Marker string
	MaxKeys int
	Delimiter string
	IsTruncated bool
	CommonPrefixes []CommonPrefixes
	Contents []Contents
}

// todo; add listing subdirectories and Marker
//    and recursive?
func ListBucket(a S3Account, bucket string) {
	req := NewRequest(a, "GET", bucket, "/", nil)
	body := req.Send()
	lbr := ListBucketResult{}
	xml.Unmarshal(body, &lbr)
	for _,cp := range lbr.CommonPrefixes {
		cp.pp()
	}
	for _,c := range lbr.Contents {
		c.pp()
	}
}

func ListMultipartUploads(a S3Account, bucket string) {
	req := NewRequest(a, "GET", bucket, "/?uploads", nil)
	body := req.Send()
	fmt.Printf("%v\n", string(body))
}

func ShowPolicy(a S3Account, bucket string) {
	req := NewRequest(a, "GET", bucket, "/?policy", nil)
	body := req.Send()
	fmt.Printf("%v\n", string(body))
}

func ShowLocation(a S3Account, bucket string) {
	req := NewRequest(a, "GET", bucket, "/?location", nil)
	body := req.Send()
	fmt.Printf("%v\n", string(body))

}