package main

import (
	"log"
	"net/http"
)

func basePage(address, service string, specificFunctionaliy func(mux *http.ServeMux)) {
	log.Println("Starting " + service + " on " + address)
	mux := http.NewServeMux()

	specificFunctionaliy(mux) // Extracted the logic of different services

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}
	err := server.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func runService(address, endpoint string, hub *Hub) {
	specificFunctionality := func(mux *http.ServeMux) {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./pages/"+endpoint+".html")
		})
		mux.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./pages"))))
		mux.HandleFunc("/ws/"+endpoint, func(w http.ResponseWriter, r *http.Request) {
			serveWs(hub, w, r)
		})
	}
	basePage(address, endpoint, specificFunctionality)

}

func runLandingPage(address, localIP string) {
	specificFunctionality := func(mux *http.ServeMux) {
		fileServer := http.FileServer(http.Dir("./pages"))
		mux.Handle("/", fileServer)
		mux.HandleFunc("/ip", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(localIP))
		})
	}
	basePage(address, "landing page", specificFunctionality)
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
