# Osquery File Carve Server

## Init project
1. `cd osquery-file-carve-server/`
1. `go mod init github.com/CptOfEvilMinions/osquery-file-carve-server`

## Spin up Kolide stack
1. `openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout conf/nginx/tls/snakeoil.key -out conf/nginx/tls/snakeoil.crt`
1. ``
1. ``


## Tested Osquery versions
* `osquery version 4.3.0` 

## To do
* Add an API to download files
* Add LDAP auth to API
* Add MongoDB GridFS support for file storage
* Add logging
* TravisCI build
* UML/network diagram

## Refernces

### NGINX
* [StackOverFlow - nginx not blocking user agents](https://serverfault.com/questions/480492/nginx-not-blocking-user-agents)
* [StackOverFlow - Nginx: location regex for multiple paths](https://serverfault.com/questions/564127/nginx-location-regex-for-multiple-paths)
* [Github - CptOfEvilMinions/BlogProjects](https://github.com/CptOfEvilMinions/BlogProjects/tree/master/kolide-mutual-tls)

### Mongo + GriFS
* [GoDocs - package gridfs](https://godoc.org/go.mongodb.org/mongo-driver/mongo/gridfs#Bucket.OpenUploadStream)
* [Golang+MongoDB](https://blog.csdn.net/qq_25490573/article/details/103540311)
* [GoDocs - package mongo](https://godoc.org/go.mongodb.org/mongo-driver/mongo)
* [Github - MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
* [StackOverFlow - Mongodb to Mongodb GridFS](https://stackoverflow.com/questions/30694254/mongodb-to-mongodb-gridfs)
* [golang mongoDB GridFS query storage delete file](https://www.programmersought.com/article/92554631584/)
* [Quick Start: Golang & MongoDB - A Quick Look at GridFS](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--a-quick-look-at-gridfs)
* []()
* []()
* []()
* []()

### net/http
* [Format a Go string without printing?](https://stackoverflow.com/questions/11123865/format-a-go-string-without-printing)
* [Using environment variables in Go](https://flaviocopes.com/golang-environment-variables/)
* [Listening multiple ports on golang http servers (using http.Handler)](https://gist.github.com/filewalkwithme/24363472e7424bbe7028)
* [Create a Basic HTTPS Server (using TLS)](https://golangcode.com/basic-https-server-with-certificate/)
* [How can I pass the parameter to a Http handler function](https://groups.google.com/g/golang-nuts/c/SGn1gd290zI)
* [Creating a RESTful API With Golang](https://tutorialedge.net/golang/creating-restful-api-with-golang/)
* []()
* []()
* []()
* []()
* []()

### Docker
* [Build Docker with Go app: cannot find package](https://stackoverflow.com/questions/47837149/build-docker-with-go-app-cannot-find-package)
* [Create Docker microcontainers for golang using alpine image](https://medium.com/@utranand/create-docker-microcontainers-for-golang-using-alpine-image-68559b688e7d)
* [DockerHub - Golang](https://hub.docker.com/_/golang?tab=tags&page=1)
* [DockerHub - Mongo](https://hub.docker.com/_/mongo?tab=tags&page=1)
* [DockerHub - NGINX](https://hub.docker.com/_/nginx?tab=tags)
* []()
* []()
* []()
* []()
* []()

### GoLang
* [package uuid ](https://pkg.go.dev/github.com/google/uuid?tab=doc)
* [UUID1 vs UUID4 ](https://www.sohamkamani.com/blog/2016/10/05/uuid1-vs-uuid4/)
* [Home / Go Cookbook / Generating good unique ids in Go edit](https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html)
* [An example HTTP JSON response package with Golang](http://www.inanzzz.com/index.php/post/rqu6/an-example-http-json-response-package-with-golang)
* [Golang Response Snippets: JSON, XML and more](https://www.alexedwards.net/blog/golang-response-snippets)
* [Go by Example: String Formatting](https://gobyexample.com/string-formatting)
* [How to sort a slice of ints in Golang?](https://www.geeksforgeeks.org/how-to-sort-a-slice-of-ints-in-golang/)
* [Getting a slice of keys from a map](https://stackoverflow.com/questions/21362950/getting-a-slice-of-keys-from-a-map)
* [How to find length of Map in Go?](https://www.golangprograms.com/how-to-find-length-of-map-in-go.html)
* [Go by Example: Tickers](https://gobyexample.com/tickers)
* [How to Check if a Key Exists in a Map in Go](https://goruncode.com/how-to-check-if-a-key-exists-in-a-map-in-go/)
* [[go-nuts] Initialize map for struct](https://groups.google.com/forum/#!topic/golang-nuts/cSPpHPGf_a8)
* [Golang Read Write Create and Delete text file](https://www.golangprograms.com/golang-read-write-create-and-delete-text-file.html)
* [Go Mutex Tutorial](https://tutorialedge.net/golang/go-mutex-tutorial/)
* [How to get hours difference between two dates](https://stackoverflow.com/questions/23838857/how-to-get-hours-difference-between-two-dates)
* [Go by Example: Mutexes](https://gobyexample.com/mutexes)
* [How to use global var across files in a package?](https://stackoverflow.com/questions/34195360/how-to-use-global-var-across-files-in-a-package)
* [Encode and Decode Strings using Base 64](https://golangcode.com/base-64-encode-decode/)
* [Converting Int data type to Float in Go](https://www.golangprograms.com/convert-int-to-float-in-golang.html)
* [time.Since() Function in Golang With Examples](https://www.geeksforgeeks.org/time-since-function-in-golang-with-examples/)
* []()
* []()
* []()
* []()
* []()