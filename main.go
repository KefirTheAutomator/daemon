package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"flag"
	"os"

	"github.com/sevlyar/go-daemon"
)

var (
	pidFile string
	logFile string
)

func init() {
	flag.StringVar(&pidFile, "pidFile", "", "Where store the pid file?")
	flag.StringVar(&logFile, "logFile", "", "Where store the log file?")
	flag.Parse()
	checkFlag(pidFile, "pidFile is necessary")
	checkFlag(logFile, "logFile is necessary")
}

func main() {
	context := &daemon.Context {
		PidFileName: pidFile,
		PidFilePerm: 0640,
		LogFileName: logFile,
		LogFilePerm: 0640,
		Umask: 027,
	}

	daemon_process, err := context.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if daemon_process != nil {
		return
	}
	defer context.Release()

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")

	serveHTTP()
}

func serveHTTP() {
	http.HandleFunc("/", httpHandler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}

func httpHandler(writer http.ResponseWriter, request *http.Request) {
	log.Printf("request from %s: %s %q",
		request.RemoteAddr,
		request.Method,
		request.URL)

	fmt.Fprintf(writer,
		"go-daemon: %q",
		html.EscapeString(request.URL.Path))
}

func checkFlag(flagToCheck string, msgToPrint string) {
	if flagToCheck == "" {
		fmt.Println(msgToPrint)
		os.Exit(1)
	}
}
