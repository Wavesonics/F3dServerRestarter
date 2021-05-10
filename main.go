package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	ginglog "github.com/szuecs/gin-glog"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	var ipAddr string
	var portNum int

	nullAuth := "null"
	var auth string

	flag.StringVar(&auth, "auth", nullAuth, "Authentication password")

	flag.StringVar(&ipAddr, "a", "0.0.0.0", "IP address for repository  to listen on")
	flag.IntVar(&portNum, "p", 8080, "TCP port for repository to listen on")
	flag.Parse()

	if auth == nullAuth {
		glog.Error("auth not provided")
		return
	}

	serveAddr := net.JoinHostPort(ipAddr, strconv.Itoa(portNum))
	router := initApp(auth)
	http.ListenAndServe(serveAddr, router)

	glog.Info("hello world")
}

func initApp(auth string) http.Handler {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(ginglog.Logger(3 * time.Second))

	router.GET("/restart", func(c *gin.Context) {
		providedAuth := c.Query("auth")

		if providedAuth == auth {
			whichServer := c.Query("server")

			switch whichServer {
			case "official1":
				{
					executeServerRestart(1)
				}
			case "official2":
				{
					executeServerRestart(2)
				}
			}

			c.String(http.StatusOK, "Server[s] restarted")
		} else {
			c.String(http.StatusUnauthorized, "Incorrect authorization")
		}
	})

	return router
}

func executeServerRestart(serverNumber int) {
	scriptName := fmt.Sprintf("restart-server-%d.sh", serverNumber)
	executeShellScript(scriptName)
}

func executeShellScript(shellScren string) {
	cmd := exec.Command("/bin/sh", shellScren)

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}
