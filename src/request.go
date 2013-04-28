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
)

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
func NewRequest(a S3Account, verb, bucket, path string) *Request {
	var url,sign string

	date := time.Now().Format(time.RFC822Z)
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
	fmt.Println(url)

	
	r,_ := http.NewRequest(verb, url, nil)
	c   := &http.Client{}
	req := &Request{verb, url, bucket, *r, *c}
	// http://golang.org/pkg/time/#pkg-constants

	auth := fmt.Sprintf("AWS %s:%s", a.access_key, sign)
	//fmt.Printf("%v\n", sign)
	
	req.req.Header.Add("Authorization", auth)
	req.req.Header.Add("x-amz-date", date)
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

func Sign(access_secret, verb, md5, ctype, date, hdrs, path string) string {
	str := fmt.Sprintf("%s\n%s\n%s\n\nx-amz-date:%s%s\n%s", verb, md5, ctype, date, hdrs, path)
	return base64_hmac(access_secret, str)
}
