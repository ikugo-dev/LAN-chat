package main

import (
	"log"
	"net/http"
)

func runService(address, endpoint string, hub *Hub) {
	log.Println("Starting " + endpoint + " on " + address)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pages/"+endpoint+".html")
	})
	mux.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./pages"))))
	mux.HandleFunc("/ws/"+endpoint, func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
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

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller
	go client.writePump()
	go client.readPump()
}

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
