package main

import (
	"go/build"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"path/filepath"
)

var (
	upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
	}
	wsConnA *websocket.Conn
	wsConnB *websocket.Conn
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Serving file.\n")
	home := filepath.Join(defaultAssetPath(), "home.html")
	http.ServeFile(w, r, home)
}

func foobar(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I hate %s!", r.URL.Path[1:])
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("upgrading to ws.\n")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	go processWs(conn)
}

func processWs(conn *websocket.Conn) {
	// echo service
	for {
    messageType, p, err := conn.ReadMessage()
		fmt.Printf("received message %s.\n", p)
    if err != nil {
			return
    }
    if err = conn.WriteMessage(messageType, []byte("I HATE YOU")); err != nil {
			return
    }
	}
}

func wsHandlerA(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("upgrading to ws.\n")

	var err error
	wsConnA, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	go processWsA()
}

func processWsA() {
	// echo service
	for {
    messageType, p, err := wsConnA.ReadMessage()
		fmt.Printf("received message %s.\n", p)
    if err != nil {
			return
    }
    if err = wsConnB.WriteMessage(messageType, p); err != nil {
			return
    }
	}
}

func wsHandlerB(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("upgrading to ws.\n")
	var err error
	wsConnB, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	go processWsB()
}

func processWsB() {
	// echo service
	for {
    messageType, p, err := wsConnB.ReadMessage()
		fmt.Printf("received message %s.\n", p)
    if err != nil {
			return
    }
    if err = wsConnA.WriteMessage(messageType, p); err != nil {
			return
    }
	}
}


func chatAHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Serving file.\n")
	home := filepath.Join(defaultAssetPath(), "chata.html")
	http.ServeFile(w, r, home)
}

func chatBHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Serving file.\n")
	home := filepath.Join(defaultAssetPath(), "chatb.html")
	http.ServeFile(w, r, home)
}

func defaultAssetPath() string {
	p, err := build.Default.Import("github.com/mclazarus/operator", "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}



func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/foobar", foobar)
	http.HandleFunc("/chata", chatAHandler)
	http.HandleFunc("/chatb", chatBHandler)
	http.HandleFunc("/ws/chata", wsHandlerA)
	http.HandleFunc("/ws/chatb", wsHandlerB)
	http.HandleFunc("/ws", wsHandler)
	http.ListenAndServe(":8080", nil)
}
