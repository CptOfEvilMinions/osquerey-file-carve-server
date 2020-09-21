## Oquery-file-carve server

## File storeage backend
This GoLang file server supports writing file uploads to the file system or to Mongo GridFS. 

### File system
#### Pros
* Simplest configuration
* Doesn't require management of a backend database

#### Cons
* Doesn't allow for scalability horizontally

### Mongo GridFS
#### Pros
* Ability to scale horizontally 

#### Cons
* Requires Mongo GridFS
* More resources
* Requires management of a backend database

## References
### File streaming
* [StackOverFlow - How to read/write from/to file using Go?](https://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file-using-go)
* [StackOverFlow - Getting “bytes.Buffer does not implement io.Writer” error message](https://stackoverflow.com/questions/23454940/getting-bytes-buffer-does-not-implement-io-writer-error-message)
* [GoDocs - Package io](https://golang.org/pkg/io/)
* [Streaming IO in Go](https://medium.com/learning-the-go-programming-language/streaming-io-in-go-d93507931185)

### Mongo + GridFS streaming
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