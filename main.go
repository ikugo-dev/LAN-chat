package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	chatPort = ":8080"
	votePort = ":8081"
	filePort = ":8082"
	drawPort = ":8083"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

func main() {
	flag.Parse()

	chatHub := newHub()
	voteHub := newHub()
	fileHub := newHub()

	go chatHub.run()
	go voteHub.run()
	go fileHub.run()

	var localIP = getLocalIP()

	http.HandleFunc("/", serveHome)
	go runService(localIP+chatPort, "chat", chatHub)
	go runService(localIP+votePort, "vote", voteHub)
	go runService(localIP+filePort, "file", fileHub)
	// go runService(localIP+drawPort, "draw", drawHub)

	waitForQuit()
}

func runService(adress, endpoint string, hub *Hub) {
	log.Println("Starting " + endpoint + " on " + adress)

	http.HandleFunc("/ws/"+endpoint, func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err := http.ListenAndServe(adress, nil)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "Unknown IP"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "Unknown IP"
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
