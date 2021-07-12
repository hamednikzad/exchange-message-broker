package broker

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Up and running"))
	//http.ServeFile(w, r, "home.html")
}

func Run(addr string) *Hub {
	hub := NewHub()
	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(w, r)
	})
	go func() {
		logrus.Infof("Listening on %s", addr)
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			logrus.Error(err)
		}
	}()

	return hub
}
