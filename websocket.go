package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"sync"

	utils "github.com/kasperskaarhoj/toweltornado/utils"

	"github.com/gorilla/websocket"
	log "github.com/s00500/env_logger"
)

type HiScoreEntry struct {
	Name string
	Wind float64
	Time float64 // time in ms
}

type wsToClient struct {
	Ping bool `json:",omitempty"`

	// Sends data for wind...
	DeviceNumber int      `json:",omitempty"`
	WindSpeed    *float64 `json:",omitempty"`

	ShowHiScore   bool   `json:",omitempty"` // If true, shows score board
	ShowGame      bool   `json:",omitempty"` // If true, shows game
	ResetGameView bool   `json:",omitempty"` // If true, resets the game view
	NewGameName   string `json:",omitempty"` // If set, shows game and starts game with name (how fast can you get all hats off...)

	HiScore []*HiScoreEntry `json:",omitempty"`
}

type wsFromClient struct {
	SendHiScore bool `json:",omitempty"` // If true, shall send HiScore

	ShowHiScore   bool   `json:",omitempty"` // If true, shows score board
	ShowGame      bool   `json:",omitempty"` // If true, shows game
	ResetGameView bool   `json:",omitempty"` // If true, resets the game view
	NewGameName   string `json:",omitempty"` // If set, will start a new game

	UpdateHiScoreEntry *HiScoreEntry `json:",omitempty"`
}

type wsServer struct {
	clients []*wsClient
	tt      *TowelTornado
	sync.Mutex
}

type wsClient struct {
	msgToClient chan []byte
	quit        chan bool
}

func (wsserver *wsServer) Push(w *wsClient) {
	wsserver.Lock()
	defer wsserver.Unlock()

	wsserver.clients = append(wsserver.clients, w)
}

func (wsserver *wsServer) Pull(w *wsClient) bool {
	wsserver.Lock()
	defer wsserver.Unlock()
	for i, wsclient := range wsserver.clients {
		if wsclient == w {
			wsserver.clients = append(wsserver.clients[:i], wsserver.clients[i+1:]...)
			//fmt.Println("Removed ", i)
			return true
		}
	}
	return false
}

func (wsserver *wsServer) Iter(routine func(*wsClient)) {
	wsserver.Lock()
	defer wsserver.Unlock()
	for _, wsclient := range wsserver.clients {
		routine(wsclient)
	}
}

func (wsserver *wsServer) createServer(port int, tt *TowelTornado) *http.Server {

	if port <= 0 {
		log.Errorf("Port %d was invalid", port)
		return nil
	}

	wsserver.tt = tt

	// create `ServerMux`
	mux := http.NewServeMux()

	// create a default route handler
	mux.HandleFunc("/", wsserver.homePage)
	mux.HandleFunc("/ws", wsserver.wsEndpoint)

	// create new server
	server := http.Server{
		Addr:    fmt.Sprintf(":%v", port), // :{port}
		Handler: mux,
	}

	// register a cleanup function
	server.RegisterOnShutdown(func() {
		wsserver.Iter(func(w *wsClient) {
			w.quit <- true
		})
	})

	// return new server (pointer)
	return &server
}

func (w *wsClient) Start(wsserver *wsServer, ws *websocket.Conn) {
	w.msgToClient = make(chan []byte, 10) // some buffer size to avoid blocking
	go func() {
		for {
			select {
			case msg := <-w.msgToClient:
				err := ws.WriteMessage(1, msg)
				if err != nil {

					go func() { // This is wrapped in a go-routine since otherwise iteration over things to send would result in a lock on the wsserver mutex inside Push/Pull/Iter. This seems to fix it...:
						log.Println("Removing a Client")
						if !wsserver.Pull(w) {
							log.Should(log.Wrap(err, "on writing to ws client. Tried to remove it, but nothing was removed..."))
						}
						//log.Println("Done Removing a Client")
					}()
				}
			case <-w.quit:
				return
			}
		}
	}()
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (wsserver *wsServer) homePage(w http.ResponseWriter, r *http.Request) {
	// Check if the request path is "/"
	if r.URL.Path == "/" {
		html := string(utils.ReadResourceFile("resources/index.html"))
		fmt.Fprint(w, html)
		return
	}

	// If the path is not "/", read the filename from the request path
	filename := filepath.Join("resources", r.URL.Path)
	// Read the file contents
	fileContent := utils.ReadResourceFile(filename)

	// Set the appropriate content type based on the file extension
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)

	// Write the file contents to the response
	w.Write(fileContent)
}

func (wsserver *wsServer) wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	ww := &wsClient{}
	ww.Start(wsserver, ws)
	log.Println("Adding a Client from ", ws.RemoteAddr())
	wsserver.Push(ww)

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	wsserver.reader(ws, ww)
	ww.quit <- true
	log.Println("Exit")
}

// define a reader which will listen for
// new messages being sent to our WebSocket
// endpoint
func (wsserver *wsServer) reader(conn *websocket.Conn, client *wsClient) {
	for {

		// Read in a message
		_, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Parse that message into a struct:
		wsFromClientObj := &wsFromClient{}
		err = json.Unmarshal(p, wsFromClientObj)
		if log.Should(err) {
			return
		}

		log.Println("Received from websocket: ", log.Indent(wsFromClientObj))

		if wsFromClientObj.SendHiScore {
			msg := wsToClient{}
			msg.HiScore, _ = LoadHiScores("hiscores.json")

			msgAsString, _ := json.Marshal(msg)
			client.msgToClient <- msgAsString
		}

		if wsFromClientObj.ShowHiScore {
			wsserver.BroadcastMessage(&wsToClient{ShowHiScore: true})
		}

		if wsFromClientObj.ResetGameView {
			wsserver.BroadcastMessage(&wsToClient{ResetGameView: true})
		}

		if wsFromClientObj.ShowGame {
			wsserver.BroadcastMessage(&wsToClient{ShowGame: true})
		}

		if wsFromClientObj.NewGameName != "" {
			wsserver.BroadcastMessage(&wsToClient{NewGameName: wsFromClientObj.NewGameName})
		}

		if wsFromClientObj.UpdateHiScoreEntry != nil {
			msg := wsToClient{}
			msg.HiScore = AddToHiscore(wsFromClientObj.UpdateHiScoreEntry)
			wsserver.BroadcastMessage(&msg)
		}
	}
}

func (wsserver *wsServer) BroadcastMessage(msg *wsToClient) {

	msgAsString, err := json.Marshal(msg)
	log.Should(log.Wrap(err, "on marshalling message to client"))

	//	fmt.Println(string(msgAsString))
	wsserver.Iter(func(w *wsClient) { w.msgToClient <- msgAsString })
}
