package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	chatPort = "8081"
	votePort = "8082"
	filePort = "8083"
	drawPort = "8084"
)

func main() {
	chatHub := newHub()
	voteHub := newHub()
	fileHub := newHub()

	go chatHub.run()
	go voteHub.run()
	go fileHub.run()

	var localIP, ok = getLocalIP()
	if !ok {
		log.Fatal("Couldn't get IP")
		return
	}

	http.HandleFunc("/", serveHome)
	go runService(localIP+":"+chatPort, "chat", chatHub)
	go runService(localIP+":"+votePort, "vote", voteHub)
	go runService(localIP+":"+filePort, "file", fileHub)
	// go runService(localIP+drawPort, "draw", drawHub)

	waitForQuit()
}

func getLocalIP() (ip string, ok bool) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "Unknown IP", false
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), true
			}
		}
	}
	return "Unknown IP", false
}

func waitForQuit() {
	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Press 'q' to quit.")
	for scanner.Scan() {
		if scanner.Text() == "q" {
			log.Println("Exiting...")
			os.Exit(0)
		}
	}
}
