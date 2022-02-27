package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"loggerhead/loggerhead"
	"net/http"
)

var (
	workDir     = flag.String("f", ".", "Working directory to be used during runtime")
	program     = flag.String("b", "", "Program to start in host OS binary format for instance python")
	programArgs = flag.String("a", "", "space separated arguement to be pass to the program for instance python filename.py")
	stdout      = flag.String("o", "", "file path to show stdout and stderror result. default to stdout.in home directory")
)

func main() {
	flag.Parse()

	app := gin.Default()
	app.LoadHTMLGlob("fe/*.html")

	ssHandler := &loggerhead.SocketHandler{}
	app.GET("/socket.io/", gin.WrapH(ssHandler))
	app.POST("/socket.io/", gin.WrapH(ssHandler))

	app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	go loggerhead.Watch(
		*program,
		*workDir,
		*programArgs,
		*stdout,
	)

	http.ListenAndServe(":5050", app)
}
