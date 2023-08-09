package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Récupérer l'origine de la demande
		origin := r.Header.Get("Origin")

		// Remplacez "http://localhost:5173" par votre URL autorisée
		allowedOrigin := "http://localhost:5173"

		// Vérifier si l'origine de la demande correspond à l'URL autorisée
		return origin == allowedOrigin
	},
}

type Client struct {
	conn *websocket.Conn
}

type Server struct {
	Router    *Router
	clients   map[*Client]bool
	Socket    *websocket.Upgrader
	broadcast chan []byte
}

func NewServer() *Server {
	fmt.Println("Création du server")
	return &Server{
		Router:    NewRouter(),
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
	}
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Erreur lors de la mise à niveau de la connexion WebSocket: %v", err)
		return
	}
	defer conn.Close()

	client := &Client{
		conn: conn,
	}
	s.clients[client] = true

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Erreur lors de la lecture du message: %v", err)
			delete(s.clients, client)
			break
		}

		// Traitez le message reçu (p) comme JSON
		var data map[string]interface{}
		if err := json.Unmarshal(p, &data); err != nil {
			log.Printf("Erreur lors de la conversion JSON: %v", err)
			continue
		}

		// Vous pouvez maintenant manipuler les données JSON comme une carte (map)
		// data contient les données envoyées par le client en JSON

		// Répondre au client (facultatif)

		s.broadcast <- p
	}
}

func (s *Server) broadcastMessage() {
	for {
		message := <-s.broadcast
		for client := range s.clients {
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Erreur lors de l'envoi du message: %v", err)
				client.conn.Close()
				delete(s.clients, client)
			}
		}
	}
}

func (s *Server) Run() {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "authorization"},
		ExposedHeaders:   []string{"authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(s.Router.Router)
	s.Router.Router.HandleFunc("/ws", s.handleWebSocket)
	go s.broadcastMessage()
	if err := http.ListenAndServe(":3000", handler); err != nil {
		log.Fatalf("Erreur pendant le lancement du serveur%v", err)
	}

	fmt.Println("Server running")
}
