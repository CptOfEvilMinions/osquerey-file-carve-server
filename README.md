# Osquery File Carve Server

## UML diagram
<UML diagram>


## Build project
1. `cd osquery-file-carve-server/`
1. `go mod init github.com/CptOfEvilMinions/osquery-file-carve-server`
1. `go build`

## Assumptions
* All blocks of data sent by Osquery will arrive in order
* All data block sizes for NGINX, Osquery, and Mongo have appropriate settings

## Setup
### Block size configuration
The default setting for all the configs in this repo is to set the data block size at 10MB (10000000 bytes). Osquery has `carver_block_size` set to 10000000 (10MB), NGINX has `client_max_body_size` set to `15MB` for TCP and TLS overhead, and Mongo GridFS blocksize is set to `10MB`.
* Mongo GridFS blockstore has chunk size/block size set to `255K` by default.

### Generate SSL certs

1. `openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout conf/kolide/tls/snakeoil.key -out conf/kolide/tls/snakeoil.crt`
  * KOLIDE WILL ONLY WORK WITH TLS 
1. `openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout conf/nginx/tls/snakeoil.key -out conf/nginx/tls/snakeoil.crt`

### Spin up stack
1. `docker-compose build`
1. `docker-compose run --rm kolide fleet prepare db --config /etc/kolide/kolide.yml`
1. `docker-compose up -d`
1. [Setup Kolide](docs/kolide_osquery.md#Init_Kolide)
1. [Install Osquery](docs/kolide_osquery.md#Install-Osquery-4.4.0-on-Ubuntu-20.04)
1. [Enroll Osquery with Kolide](docs/kolide_osquery.md#Add-Osquery-host-to-Kolide)


## Resource stats
* If the Osquery `carver_block_size` is set to `1000000` (1MB) the osquery-file-carve server will consume roughly 30MB
* If the Osquery `carver_block_size` is set to `10000000` (10MB) the osquery-file-carve server will consume roughly 140MB


## Tested Osquery versions
* `osquery version 4.3.0` 
* `osquery version 4.4.0` 

## To do
* Add the ability to clean up unsucessful uploads from Mongo
  * Files will be deleted at specified clean up interval
* TravisCI build
* UML/network diagram

## Refernces
### NGINX
* [StackOverFlow - nginx not blocking user agents](https://serverfault.com/questions/480492/nginx-not-blocking-user-agents)
* [StackOverFlow - Nginx: location regex for multiple paths](https://serverfault.com/questions/564127/nginx-location-regex-for-multiple-paths)
* [Github - CptOfEvilMinions/BlogProjects](https://github.com/CptOfEvilMinions/BlogProjects/tree/master/kolide-mutual-tls)
* [NGINX Reverse Proxy](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/)
* [Introducing gRPC Support with NGINX 1.13.10](https://www.nginx.com/blog/nginx-1-13-10-grpc/)
* [Kolide Fleet – Breaking out the osquery API & Web UI](https://defensivedepth.com/2020/04/02/kolide-fleet-breaking-out-the-osquery-api-web-ui/)
* [Slack - Osquery - NGINX config](https://osquery.slack.com/archives/C1XCLA5DZ/p1567760131054400?thread_ts=1555590511.044500&cid=C1XCLA5DZ)
* []()

### Kolide
* [Configuring The Fleet Binary](https://github.com/kolide/fleet/blob/master/docs/infrastructure/configuring-the-fleet-binary.md)
* [Kolide Fleet – Breaking out the osquery API & Web UI](https://defensivedepth.com/2020/04/02/kolide-fleet-breaking-out-the-osquery-api-web-ui/)
* []()
* []()
* []()

### Mongo + GridFS
* [GoDocs - package gridfs](https://godoc.org/go.mongodb.org/mongo-driver/mongo/gridfs#Bucket.OpenUploadStream)
* [Golang+MongoDB](https://blog.csdn.net/qq_25490573/article/details/103540311)
* [GoDocs - package mongo](https://godoc.org/go.mongodb.org/mongo-driver/mongo)
* [Github - MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
* [StackOverFlow - Mongodb to Mongodb GridFS](https://stackoverflow.com/questions/30694254/mongodb-to-mongodb-gridfs)
* [golang mongoDB GridFS query storage delete file](https://www.programmersought.com/article/92554631584/)
* [Quick Start: Golang & MongoDB - A Quick Look at GridFS](https://www.mongodb.com/blog/post/quick-start-golang--mongodb--a-quick-look-at-gridfs)
* [Connection String URI Format](https://docs.mongodb.com/manual/reference/connection-string/)
* [Github issue - how to configure TLS/SSL for mongod](https://github.com/docker-library/mongo/issues/250)
* [StackOverFlow - How to use new URL from mongodb 3.6 to connect from golang](https://stackoverflow.com/questions/52052311/how-to-use-new-url-from-mongodb-3-6-to-connect-from-golang)
* [StackOverflow- How to download file in browser from Go server](https://stackoverflow.com/questions/24116147/how-to-download-file-in-browser-from-go-server)
* [primitive.go](https://sourcegraph.com/github.com/mongodb/mongo-go-driver/-/blob/bson/primitive/primitive.go#L74:24)
* [StackOverFlow - Delete all the document older than a date using _id in mongo using mgo](https://stackoverflow.com/questions/34412502/delete-all-the-document-older-than-a-date-using-id-in-mongo-using-mgo)
* []()
* []()


### File I/O
* [StackOverFlow - How to read/write from/to file using Go?](https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file-using-go)
* [StackOverFlow - Getting “bytes.Buffer does not implement io.Writer” error message](https://stackoverflow.com/questions/23454940/getting-bytes-buffer-does-not-implement-io-writer-error-message)
* [GoDocs - Package io](https://golang.org/pkg/io/)
* [Streaming IO in Go](https://medium.com/learning-the-go-programming-language/streaming-io-in-go-d93507931185)

### net/http + HTTP client
* [Format a Go string without printing?](https://stackoverflow.com/questions/11123865/format-a-go-string-without-printing)
* [Using environment variables in Go](https://flaviocopes.com/golang-environment-variables/)
* [Listening multiple ports on golang http servers (using http.Handler)](https://gist.github.com/filewalkwithme/24363472e7424bbe7028)
* [Create a Basic HTTPS Server (using TLS)](https://golangcode.com/basic-https-server-with-certificate/)
* [How can I pass the parameter to a Http handler function](https://groups.google.com/g/golang-nuts/c/SGn1gd290zI)
* [Creating a RESTful API With Golang](https://tutorialedge.net/golang/creating-restful-api-with-golang/)
* [Copy a file in Go](https://shapeshed.com/copy-a-file-in-go/)
* [net/http - DetectContentType](https://golang.org/pkg/net/http/#DetectContentType)
* [Golang transmit files over a net/http server to clients](https://mrwaggel.be/post/golang-transmit-files-over-a-nethttp-server-to-clients/)
* [StackOverFlow - How to download file in browser from Go server](https://stackoverflow.com/questions/24116147/how-to-download-file-in-browser-from-go-server)
* [Making and Using HTTP Middleware](https://www.alexedwards.net/blog/making-and-using-middleware)
* [StackOverFlow - How to set headers in http get request?](https://stackoverflow.com/questions/12864302/how-to-set-headers-in-http-get-request)
* [Making HTTP Requests in Golang](https://medium.com/@masnun/making-http-requests-in-golang-dd123379efe7)
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

### JWT tokens
* [Simple JWT Authentication for Golang (Part 1)](https://dev.to/omnisyle/simple-jwt-authentication-for-golang-part-1-3kfo)
* [Securing Golang API using Json Web Token (JWT)](https://medium.com/@Raulgzm/securing-golang-api-using-json-web-token-jwt-2dc363792a48)
* [Using JWT for Authentication in a Golang Application](https://www.nexmo.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr)
* [https://www.alexedwards.net/blog/making-and-using-middleware](https://www.alexedwards.net/blog/making-and-using-middleware)
* [Github - dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)


### Vault
* [Vault - Token - Lookup a Token (Self)](https://www.vaultproject.io/api-docs/auth/token#lookup-a-token-self)
* [Github - adfinis-sygroup/vault-client - auth_backend.go](https://github.com/adfinis-sygroup/vault-client/blob/master/src/auth_backend.go)
* [Github - adfinis-sygroup/vault-client - auth.go](https://github.com/adfinis-sygroup/vault-client/blob/master/src/auth.go)
* [Vault - Read Health Information](https://www.vaultproject.io/api-docs/system/health)
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
* [Check If a File Exists Before Using It](https://golangcode.com/check-if-a-file-exists/)
* [Golang : Convert Unix timestamp to UTC timestamp](https://socketloop.com/tutorials/golang-convert-unix-timestamp-to-utc-timestamp#:~:text=Solution%20%3A,Coordinated%20Universal%20Time)%20time%20stamp.)
* [How to compare times in Golang?](https://www.geeksforgeeks.org/how-to-compare-times-in-golang/)
* [Subtracting time.Duration from time in Go](https://stackoverflow.com/questions/26285735/subtracting-time-duration-from-time-in-go)
* []()
* []()
* []()
* []()
* []()
* []()