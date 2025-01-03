package engines

import (
	"concernedmate/SurvivorGame/entities"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"math/rand"

	"github.com/gorilla/websocket"
)

var Chan bool
var Rooms map[string]*entities.Game

func Server() {
	var upgrader = websocket.Upgrader{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./engines/client/client.html")
	})
	http.HandleFunc("/client.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./engines/client/client.js")
	})

	http.HandleFunc("/ws_client", func(w http.ResponseWriter, r *http.Request) {
		id_room := r.Header.Get("id_room")
		Game := Rooms[id_room]

		if len(Game.Players) > 3 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Room is full",
			})
			return
		}
		id := "player" + strconv.Itoa(rand.Intn(1000))
		Game.CreatePlayer(id)
		Game.Sync <- true
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer func() {
			Game.DestroyPlayer(id)
			Game.Sync <- true
			ws.Close()
		}()

		ws.WriteJSON(map[string]string{"id": id})
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("errs:", err)
				break
			}
			playerEvent := string(message)
			for idx, players := range Game.Players {
				if players.Uid == id {
					Game.Players[idx].Events = strings.Split(playerEvent, ",")
					Game.Sync <- true
				}
			}
		}
	})
	http.HandleFunc("/ws_server", func(w http.ResponseWriter, r *http.Request) {
		id_room := r.Header.Get("id_room")
		Game := Rooms[id_room]

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer ws.Close()
		timer := time.Now()
		for {
			<-Game.Sync
			dur := time.Since(timer).Milliseconds()
			if dur > 30 {
				mobsData := map[int][]int{}
				for _, mob := range Game.Mobs {
					mobsData[mob.Size] = append(mobsData[mob.Size], []int{int(mob.PosX), int(mob.PosY)}...)
				}

				projsData := map[int][]int{}
				for _, proj := range Game.Projectiles {
					projsData[proj.Size] = append(projsData[proj.Size], []int{int(proj.PosX), int(proj.PosY)}...)
				}

				playersData := []map[string]any{}
				for _, player := range Game.Players {
					playersData = append(playersData, map[string]any{
						"Uid":    player.Uid,
						"Health": player.Health,
						"Mana":   player.Mana,
						"Score":  player.Score,
						"PosX":   int(player.PosX),
						"PosY":   int(player.PosY),
						"Size":   player.Size,
					})
				}

				err := ws.WriteJSON(map[string]any{
					"mobs":        mobsData,
					"projectiles": projsData,
					"players":     playersData,
				})
				if err != nil {
					log.Println("errs:", err)
					break
				}
				timer = time.Now()
			}
		}
	})

	fmt.Println("Server is started on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func OpenRoom(id_room string) {
	if Rooms == nil {
		Rooms = make(map[string]*entities.Game)
	}
	if Rooms[id_room] != nil {
		return
	}
	// GAME
	var Game = entities.Game{
		DeltaTime: 0.0,
		Ticker:    *time.NewTicker(time.Millisecond * 500),
		Sync:      make(chan bool),
	}
	Rooms[id_room] = &Game

	// LOOP
	go GameSpawners(&Game)
	for {
		if len(Game.Players) > 0 {
			GameLoop(&Game)
		}
	}
}
