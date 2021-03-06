package mys3

import (
	"code.google.com/p/goconf/conf" //http://code.google.com/p/goconf/
	"net/http"
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"time"
	"io/ioutil"
	"io"
)

func TestSign(a *S3Account, str string) {
	s := base64_hmac(a.secret_key, str)
	fmt.Println(s)
}

func base64_hmac(access_secret, str string) string {
	h := hmac.New(sha1.New, []byte(access_secret))
	h.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

type S3Account struct {
	Host_base, access_key, secret_key string
}
func (a *S3Account) Load( file string ) {
	c,_ := conf.ReadConfigFile(file)
	a.Host_base,_  = c.GetString("default", "host_base")
	a.access_key,_ = c.GetString("default", "access_key")
	a.secret_key,_ = c.GetString("default", "secret_key")
}

type Request struct {
	verb   string
	url    string
	bucket string
 	req    http.Request
	client http.Client
}

func NewRequest(a S3Account, verb,
	bucket, path string, body io.Reader) *Request {
	var url,sign string

	loc,_ := time.LoadLocation("GMT")
	date := time.Now().In(loc).Format(time.RFC1123Z)
	//date := time.Now().Format(time.RFC822Z)
	// http://golang.org/pkg/time/#pkg-constants

	if bucket == "" {
		url = fmt.Sprintf("http://%s%s", a.Host_base, path)
		sign = Sign(a.secret_key, verb, "", "", date, "", path)
	}else if path == "/"{
		url = fmt.Sprintf("http://%s.%s/?delimiter=/", bucket, a.Host_base)
		sign = Sign(a.secret_key, verb, "", "", date, "", "/" + bucket + path)
	}else{
		url = fmt.Sprintf("http://%s.%s%s", bucket, a.Host_base, path)
		sign = Sign(a.secret_key, verb, "", "", date, "", "/" + bucket + path)
	}
	//fmt.Println(url)
	//fmt.Println(sign)
	
	r,_ := http.NewRequest(verb, url, body)
	c   := &http.Client{}
	req := &Request{verb, url, bucket, *r, *c}

	auth := fmt.Sprintf("AWS %s:%s", a.access_key, sign)
	//fmt.Printf("%v\n", sign)
	
	req.req.Header.Add("Authorization", auth)
	req.req.Header.Add("x-amz-date", date)
	req.req.Header.Add("Date", date)
	//req.Header.Add("Content-Length", "0")
	return req
}
func (r *Request) Send() []byte {
	req := &r.req
	resp,_ := r.client.Do(req)
	fmt.Printf("%v -> %v\n", req, resp)
	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)
	return body
}

func (r *Request) SendAndWriteFile(out io.Writer) int64 {
	req := &r.req
	resp,e0 := r.client.Do(req)
	defer resp.Body.Close()
	
	if e0 != nil {
		fmt.Printf("%v\n", req)
		panic("can't fetch data")
	}
	//fmt.Printf("%v -> %v\n", e0, resp.Body)
	n,e := io.Copy(out, resp.Body)
	if e != nil { panic("can't copy to file") }
	return n
}

func Sign(access_secret, verb, md5, ctype, date, hdrs, path string) string {
	str := fmt.Sprintf("%s\n%s\n%s\n\nx-amz-date:%s%s\n%s", verb, md5, ctype, date, hdrs, path)
	return base64_hmac(access_secret, str)
}
