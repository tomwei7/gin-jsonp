# gin-jsonp

JSONP Middleware for Gin Framework

### usage

```
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tomwei7/gin-jsonp"
)

func main() {
	r := gin.Default()
	r.Use(jsonp.JsonP())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("127.0.0.1:8080") // listen and server on 0.0.0.0:8080
}
```

**output:**

```
curl -v http://127.0.0.1:8080/ping\?callback\=callback

*   Trying 127.0.0.1...
> GET /ping?callback=callback HTTP/1.1
> Host: 127.0.0.1:8080
> User-Agent: curl/7.43.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< Date: Sun, 28 Aug 2016 05:52:00 GMT
< Content-Length: 28
< 
* Connection #0 to host 127.0.0.1 left intact
callback({"message":"pong"})
```

