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
	router.Use(ginglog.Logger(3 * time.Second))

	router.GET("/restart", func(c *gin.Context) {
		providedAuth := c.Param("auth")

		if providedAuth == auth {
			whichServer := c.Param("server")

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
			c.String(http.StatusUnauthorized, fmt.Sprintf("Incorrect authorization: %s", providedAuth, auth))
		}
	})

	return router
}

func executeServerRestart(serverNumber int) {
	command := fmt.Sprintf("systemctl restart game-server-%d", serverNumber)
	executeCommand(command)
}

func executeCommand(command string) {
	cmd := exec.Command(command)

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}
