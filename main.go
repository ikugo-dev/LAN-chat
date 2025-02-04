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
		http.Error(w, "Use port numbers for navigation, not path", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed, you can only GET", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./pages/home.html")
}

func main() {
	flag.Parse()

	chatHub := newHub()
	voteHub := newHub()
	fileHub := newHub()

	go chatHub.run()
	go voteHub.run()
	go fileHub.run()

	var localIP, ok = getLocalIP()
	if !ok {
		return
	}

	http.HandleFunc("/", serveHome)
	go runService(localIP+chatPort, "chat", chatHub, handleTextMessage)
	go runService(localIP+votePort, "vote", voteHub, handleVoteMessage)
	go runService(localIP+filePort, "file", fileHub, handleFileMessage)
	// go runService(localIP+drawPort, "draw", drawHub)

	waitForQuit()
}

func runService(address, endpoint string, hub *Hub, handler func(msg Message, c *Client)) {
	log.Println("Starting " + endpoint + " on " + address)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/"+endpoint+".html")
	})
	mux.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./pages"))))
	mux.HandleFunc("/ws/"+endpoint, func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r, handler)
	})

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
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

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request, handler func(msg Message, c *Client)) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller
	go client.writePump()
	go client.readPump(handler)
}
