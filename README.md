mys3
====

S3 cli client which works with s3cmd config file.

```
(set ~/.s3cfg which works fine with s3cmd before)
$ export GOPATH=`pwd`
$ mkdir src
$ cd src
$ git clone git://github.com/kuenishi/mys3
$ cd mys3
$ go get blah blah blah
$ go build
$ go run main.go ls
```

```
$ ./mys3 ls
$ ./mys3 ls <bucket>
$ ./mys3 info <bucket>
$ ./mys3 put <bucket> <localfile>
$ ./mys3 get <bucket> </path/to/object>
$ ./mys3 mp <bucket> # shows whole mulpart
```
