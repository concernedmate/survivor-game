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

const APP_PORT = 3000

var Rooms map[string]*entities.Game

func initRooms() {
	rooms_map := make(map[string]*entities.Game)
	Rooms = rooms_map
}

func getRoom(id_room string) *entities.Game {
	return Rooms[id_room]
}

func OpenRoom(id_room string) {
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
	go GameRoutine(&Game, id_room)
}

func Server() {
	initRooms()

	var upgrader = websocket.Upgrader{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			id_room := r.FormValue("id_room")
			if id_room == "" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			if getRoom(id_room) == nil {
				fmt.Printf("starting room %s ...\n", id_room)
				go OpenRoom(id_room)
			}

			for getRoom(id_room) == nil {
			}

			http.Redirect(w, r, fmt.Sprintf("/game?id_room=%s", id_room), http.StatusSeeOther)
			return
		}

		if r.Method == "GET" {
			http.ServeFile(w, r, "./engines/client/index.html")
		}
	})
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./engines/client/client.html")
	})
	http.HandleFunc("/client.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./engines/client/client.js")
	})

	http.HandleFunc("/ws_client", func(w http.ResponseWriter, r *http.Request) {
		id_room := r.URL.Query().Get("id_room")
		room := getRoom(id_room)
		if room == nil {
			fmt.Println("id_room:", id_room, "not found")
			return
		}

		Game := room

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
					select {
					case Game.Sync <- true:
					default:
					}
				}
			}
		}
	})
	http.HandleFunc("/ws_server", func(w http.ResponseWriter, r *http.Request) {
		id_room := r.URL.Query().Get("id_room")
		room := getRoom(id_room)
		if room == nil {
			fmt.Println("id_room:", id_room, "not found")
			return
		}

		Game := room

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
			if dur > 32 {
				mobsData := map[int][]int{}
				for _, mob := range Game.Mobs {
					mobsData[mob.Size] = append(mobsData[mob.Size], []int{int(mob.PosX), int(mob.PosY)}...)
				}

				projsData := map[string][]int{}
				for _, proj := range Game.Projectiles {
					sizex := fmt.Sprintf("%d|%d", proj.Size, int(proj.PosX))
					projsData[sizex] = append(projsData[sizex], int(proj.PosY))
				}

				playersData := []any{}
				for _, player := range Game.Players {
					data := []any{
						player.Uid, player.Score,
						player.Health, player.Mana,
						int(player.PosX), int(player.PosY),
						player.Size,
					}
					playersData = append(playersData, data)
				}

				err := ws.WriteJSON(map[int]any{
					0: mobsData,
					1: projsData,
					2: playersData,
				})
				if err != nil {
					log.Println("errs:", err)
					break
				}
				timer = time.Now()
			}
		}
	})

	fmt.Printf("Server is started on port %d\n", APP_PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", APP_PORT), nil))
}
