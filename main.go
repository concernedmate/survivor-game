package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"concernedmate/SurvivorGame/entities"

	"github.com/gorilla/websocket"
)

const MAP_BOUNDARY_X = 1000
const MAP_BOUNDARY_Y = 1000

func GameLoop(Game *entities.Game, sync chan bool) {
	startDelta := time.Now()
	sync <- true

	// TODO
	for idx, player := range Game.Players {
		/* [0, 0, 0, 0, 0] w a s d space */
		if len(player.Events) == 0 {
			break
		}
		if player.Events[0] == "1" {
			Game.Players[idx].PosY += float32(player.Speed) * Game.DeltaTime
		}
		if player.Events[1] == "1" {
			Game.Players[idx].PosX -= float32(player.Speed) * Game.DeltaTime
		}
		if player.Events[2] == "1" {
			Game.Players[idx].PosY -= float32(player.Speed) * Game.DeltaTime
		}
		if player.Events[3] == "1" {
			Game.Players[idx].PosX += float32(player.Speed) * Game.DeltaTime
		}
		if player.Events[4] == "1" {
			if player.WeaponCD <= 0 {
				Game.Projectiles = append(Game.Projectiles, entities.Projectile{
					OwnerId: player.Uid,
					PosX:    player.PosX + float32(rand.Intn(5)) - float32(rand.Intn(10)),
					PosY:    player.PosY + float32(rand.Intn(5)) - float32(rand.Intn(10)),
					Size:    10,
					Speed:   100,

					Damage: 0,
				})
				Game.Players[idx].WeaponCD = 100
			}
		}
		if player.WeaponCD > 0 {
			Game.Players[idx].WeaponCD -= int(time.Since(startDelta).Milliseconds()) + int(Game.DeltaTime*1000)
		}
	}
	sync <- true

	for idx, proj := range Game.Projectiles {
		if proj.PosX < 0 {
			Game.Projectiles[idx].PosX -= float32(proj.Speed) * float32(Game.DeltaTime)
			Game.Projectiles[idx].PosY -= float32(proj.Speed) * float32(Game.DeltaTime)
		} else {
			Game.Projectiles[idx].PosX += float32(proj.Speed) * float32(Game.DeltaTime)
			Game.Projectiles[idx].PosY += float32(proj.Speed) * float32(Game.DeltaTime)
		}
	}

	for i := 0; i < len(Game.Projectiles); i++ {
		if i < 0 {
			break
		}
		if Game.Projectiles[i].PosX < float32(-MAP_BOUNDARY_X) ||
			Game.Projectiles[i].PosX > float32(MAP_BOUNDARY_X) ||
			Game.Projectiles[i].PosY < float32(-MAP_BOUNDARY_Y) ||
			Game.Projectiles[i].PosY > float32(MAP_BOUNDARY_Y) {
			Game.DestroyProjectile(i)
			fmt.Printf("Destroyed projectile: %d\n", i)
			i--
		}
	}

	// CALCULATE DELTATIME
	if time.Since(startDelta).Milliseconds() == 0 {
		time.Sleep(time.Millisecond * 1)
	}
	fmt.Printf("Game processing for %d ms\n", time.Since(startDelta).Milliseconds())
	Game.DeltaTime = float32(time.Since(startDelta).Milliseconds()) / 1000
	sync <- true
}

func Server(Game *entities.Game, sync chan bool) {
	var upgrader = websocket.Upgrader{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/client.html")
	})
	http.HandleFunc("/client.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./client/client.js")
	})
	http.HandleFunc("/ws_client", func(w http.ResponseWriter, r *http.Request) {
		id := strconv.Itoa(rand.Intn(1000))
		Game.Players = append(Game.Players, entities.Player{
			Uid:   "player" + id,
			PosX:  0,
			PosY:  0,
			Size:  25,
			Speed: 100,

			Health: 10,
			Mana:   10,

			Score: 0,
		})
		sync <- true
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer func() {
			Game.DestroyPlayer("player" + id)
			sync <- true
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
				if players.Uid == "player"+id {
					Game.Players[idx].Events = strings.Split(playerEvent, ",")
					sync <- true
				}
			}
		}
	})
	http.HandleFunc("/ws_server", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer ws.Close()
		timer := time.Now()
		for {
			<-sync
			dur := time.Since(timer).Milliseconds()
			if dur > 10 {
				err := ws.WriteJSON(map[string]any{
					"mobs":        Game.Mobs,
					"projectiles": Game.Projectiles,
					"players":     Game.Players,
					"obstacles":   Game.Obstacles,
				})
				if err != nil {
					log.Println("errs:", err)
					break
				}
				timer = time.Now()
			}
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	// GAME
	var Game = entities.Game{
		DeltaTime: 0.0,
	}
	sync := make(chan bool, 1)

	// LOOP
	go Server(&Game, sync)
	for {
		if len(Game.Players) > 0 {
			GameLoop(&Game, sync)
		}
	}
}
