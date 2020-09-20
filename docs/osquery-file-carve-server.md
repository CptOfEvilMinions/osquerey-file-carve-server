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