package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"

	"BubbleMathics/db"
	"BubbleMathics/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	router := mux.NewRouter()
	db.Connect()

	router.HandleFunc("/api/questions", getQuestions).Methods("GET")
	router.HandleFunc("/ws", handleConnections)

	// Serve static files from the React app
	staticDir := filepath.Join(".", "build")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(staticDir)))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var questions []models.Question
	collection := db.GetCollection("questions")
	cursor, err := collection.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	for cursor.Next(r.Context()) {
		var question models.Question
		if err := cursor.Decode(&question); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		questions = append(questions, question)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(questions)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer ws.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go pingWebSocket(ctx, ws)

	for {
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		fmt.Printf("Message received: %+v\n", msg)
		err = ws.WriteJSON(msg)
		if err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}

func pingWebSocket(ctx context.Context, ws *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("WebSocket ping error:", err)
				return
			}
		}
	}
}
