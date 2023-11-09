# Web service with Gin

To create a web server in Golang it could be with the Gin framework, which help us to make the things easier to accept request and make the operations through the web.

First is need it to create a folder to hold the server project

```bash
mkdir web-service-gin
```

This is going to be paste in the path where we are going to create the folder

Then it is necessary to manage the dependencies with the go.mod

```bash
go mod init example.com/web-service-gin
```


After this create the main.go to create the main package

```go
package main  
  
// album represents data about a record album  
type album struct {  
    ID     string  `json:"id"`  
    Title  string  `json:"title"`  
    Artist string  `json:"artist"`  
    Price  float64 `json:"price"`  
}  
  
// `json:"id"` it's a tag to say the lenguage whic will be the  
// key in the JSON  
  
// albums slice to seed record album data.var albums = []album{  
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},  
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},  
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},  
}
```


## Write a handler to return all the items

Now it is necessary to create items to handle the client's requests

- Logic to prepare the response
- Code to map the request path to your logic

The function to return all the album items is

```go
// getAlbums responds with the list of all albums as JSON
func getAlbums(c *gin.Context) {  
    c.IndentedJSON(http.StatusOK, albums)  
}
```

**IMPORTANT**: The most important part is gin.Context; it carries the request details, validates and serializes JSON and more. (This is different from Gin's built-in **context** package)

The Context.IndentedJSON serialize the struct into JSON

You can use Context.JSON instead of Context.IndentedJSON, but the second option is easier to debug and it is smaller.


```go
func main() {  
    router := gin.Default()  
    router.GET("/albums", getAlbums)  
  
    router.Run("localhost:8080")  
}
```

Now it the main function the GET method and the "/albums" url will be accepting request and returning the response of "getAlbums" .

After all of this the code will be like 

```go
package main  
  
import (  
    "net/http"  
  
    "github.com/gin-gonic/gin")  
  
// album represents data about a record album  
type album struct {  
    ID     string  `json:"id"`  
    Title  string  `json:"title"`  
    Artist string  `json:"artist"`  
    Price  float64 `json:"price"`  
}  
  
// `json:"id"` it's a tag to say the lenguage whic will be the  
// key in the JSON  
  
// albums slice to seed record album data.var albums = []album{  
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},  
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},  
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},  
}  
  
func main() {  
    router := gin.Default()  
    router.GET("/albums", getAlbums)  
  
    router.Run("localhost:8080")  
}  
  
// getAlbums responds with the list of all albums as JSONfunc getAlbums(c *gin.Context) {  
    c.IndentedJSON(http.StatusOK, albums)  
}
```

To run this code is with 

```shell
go get .
```

With this the dependencies in the import will be downloaded to the system, to satisfy the requirements.

Now it is necessary to run the server and make the request, is with the next commands

```bash
go run .
```

```bash
curl http://localhost:8080/albums
```


## Add items with new handler

Now it's going to be created the POST method to add new items, the logic is

- Logic to add the new albums to the existing list
- A bit of code to route the POST request to your logic

```go
// postAlbums adds an album from JSON received in the request bodyfunc postAlbums(c *gin.Context) {  
    var newAlbum album  
  
    // Call BindJson to bind the received JSON to  
    // newAlbum    if err := c.BindJSON(&newAlbum); err != nil {  
       return  
    }  
  
    // Add the new album to the slice.  
    albums = append(albums, newAlbum)  
    c.IndentedJSON(http.StatusCreated, newAlbum)  
}
```

With this function can add some new elements to the albums

```shell
```
curl http://localhost:8080/albums \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'
```
```

With this call can create a new item in the album


## Write a handler to return an specific item

To get an specific item is with: GET /albums/[id]. To do this it is necessary:

- Add logic to retrieve the requested album
- Map the path to logic

The function to add is the next:

```go
// getAlbumByID locates the album whose ID value matches the id// parameter sent by the client, then returns that album as a responsefunc getAlbumByID(c *gin.Context) {  
    id := c.Param("id")  
  
    // Loop over the list of albums, looking for  
    // an album whose ID value matches the parameter    for _, a := range albums {  
       if a.ID == id {  
          c.IndentedJSON(http.StatusOK, a)  
          return  
       }  
    }  
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})  
}
```

And now change in the main function to accept the new url

```go
func main() {  
    router := gin.Default()  
    router.GET("/albums", getAlbums)  
    router.POST("/albums", postAlbums)  
    router.GET("/albums/:id", getAlbumByID)  
  
    router.Run("localhost:8080")  
}
```


And now to request for an specific albums is with the next way:

```shell
curl http://localhost:8080/albums/2
```

The last number could be replaced with the ID that we want to request.


## Summary code

The complete code is:

```go
package main  
  
import (  
    "net/http"  
  
    "github.com/gin-gonic/gin")  
  
// album represents data about a record album  
type album struct {  
    ID     string  `json:"id"`  
    Title  string  `json:"title"`  
    Artist string  `json:"artist"`  
    Price  float64 `json:"price"`  
}  
  
// `json:"id"` it's a tag to say the lenguage whic will be the  
// key in the JSON  
  
// albums slice to seed record album data.var albums = []album{  
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},  
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},  
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},  
}  
  
func main() {  
    router := gin.Default()  
    router.GET("/albums", getAlbums)  
    router.POST("/albums", postAlbums)  
    router.GET("/albums/:id", getAlbumByID)  
  
    router.Run("localhost:8080")  
}  
  
// getAlbums responds with the list of all albums as JSONfunc getAlbums(c *gin.Context) {  
    c.IndentedJSON(http.StatusOK, albums)  
}  
  
// postAlbums adds an album from JSON received in the request bodyfunc postAlbums(c *gin.Context) {  
    var newAlbum album  
  
    // Call BindJson to bind the received JSON to  
    // newAlbum    if err := c.BindJSON(&newAlbum); err != nil {  
       return  
    }  
  
    // Add the new album to the slice.  
    albums = append(albums, newAlbum)  
    c.IndentedJSON(http.StatusCreated, newAlbum)  
}  
  
// getAlbumByID locates the album whose ID value matches the id// parameter sent by the client, then returns that album as a responsefunc getAlbumByID(c *gin.Context) {  
    id := c.Param("id")  
  
    // Loop over the list of albums, looking for  
    // an album whose ID value matches the parameter    for _, a := range albums {  
       if a.ID == id {  
          c.IndentedJSON(http.StatusOK, a)  
          return  
       }  
    }  
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})  
}
```