package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"concernedmate/SurvivorGame/entities"
	"concernedmate/SurvivorGame/physics"

	"github.com/gorilla/websocket"
)

const MAP_BOUNDARY_X = entities.MAP_BOUNDARY_X
const MAP_BOUNDARY_Y = entities.MAP_BOUNDARY_Y

func GameLoop(Game *entities.Game, sync chan bool) {
	startDelta := time.Now()

	for idx, player := range Game.Players {
		/* [0, 0, 0, 0, 0] w a s d space */
		if len(player.Events) == 0 {
			break
		}
		if player.Events[0] == "1" && Game.Players[idx].PosY < MAP_BOUNDARY_Y-float32(player.Size) {
			Game.Players[idx].PosY += float32(player.Speed) * Game.DeltaTime
		}
		if player.Events[1] == "1" && Game.Players[idx].PosX > 0+float32(player.Size) {
			Game.Players[idx].PosX -= float32(player.Speed) * Game.DeltaTime
		}
		if player.Events[2] == "1" && Game.Players[idx].PosY > 0+float32(player.Size) {
			Game.Players[idx].PosY -= float32(player.Speed) * Game.DeltaTime
		}
		if player.Events[3] == "1" && Game.Players[idx].PosX < MAP_BOUNDARY_X-float32(player.Size) {
			Game.Players[idx].PosX += float32(player.Speed) * Game.DeltaTime
		}
		if player.Events[4] == "1" {
			if player.WeaponCD <= 0 {
				Game.CreateProjectileA(&Game.Players[idx])
			}
		}
		if player.WeaponCD > 0 {
			Game.Players[idx].WeaponCD -= int(time.Since(startDelta).Milliseconds()) + int(Game.DeltaTime*1000)
		}
	}

	for idx := range Game.Projectiles {
		if Game.Projectiles[idx].ProjMovement == nil {
			log.Fatal("Projectile", idx, "does not have projectile movement!")
		}
		Game.Projectiles[idx].ProjMovement(&Game.Projectiles[idx], Game.DeltaTime)
	}

	for idx := range Game.Mobs {
		if Game.Mobs[idx].MobMovement == nil {
			log.Fatal("Mob", idx, "does not have mob movement!")
		}
		Game.Mobs[idx].MobMovement(&Game.Mobs[idx], Game.DeltaTime)
	}

	for i := 0; i < len(Game.Projectiles); i++ {
		if Game.Projectiles[i].PosX < float32(-MAP_BOUNDARY_X) ||
			Game.Projectiles[i].PosX > float32(MAP_BOUNDARY_X) ||
			Game.Projectiles[i].PosY < float32(-MAP_BOUNDARY_Y) ||
			Game.Projectiles[i].PosY > float32(MAP_BOUNDARY_Y) {
			Game.Projectiles[i].DeleteFlag = true
		}
	}

	for i := 0; i < len(Game.Mobs); i++ {
		if Game.Mobs[i].PosX < float32(-MAP_BOUNDARY_X) ||
			Game.Mobs[i].PosX > float32(MAP_BOUNDARY_X) ||
			Game.Mobs[i].PosY < float32(-MAP_BOUNDARY_Y) ||
			Game.Mobs[i].PosY > float32(MAP_BOUNDARY_Y) {
			Game.Mobs[i].DeleteFlag = true
		}
	}

	// Projectile Mobs Collision
	for i := 0; i < len(Game.Projectiles); i++ {
		if Game.Projectiles[i].DeleteFlag {
			continue
		}
		for j := 0; j < len(Game.Mobs); j++ {
			if Game.Mobs[j].DeleteFlag {
				continue
			}
			if physics.ProjectileMobCollision(&Game.Projectiles[i], &Game.Mobs[j]) {
				Game.Projectiles[i].DeleteFlag = true
				Game.Mobs[j].Health -= 1
				if Game.Mobs[j].Health <= 0 {
					Game.Mobs[j].DeleteFlag = true
				}
				break
			}
		}
	}

	// Player Mobs Collision
	for i := 0; i < len(Game.Players); i++ {
		for j := 0; j < len(Game.Mobs); j++ {
			if Game.Mobs[j].DeleteFlag {
				continue
			}
			if physics.PlayerMobCollision(&Game.Players[i], &Game.Mobs[j]) {
				Game.Mobs[j].DeleteFlag = true
				Game.Players[i].Health -= 1
				if Game.Players[i].Health <= 0 {
					Game.DestroyPlayer(Game.Players[i].Uid)
				}
				break
			}
		}
	}

	// Cleanup
	Game.DestroyMob()
	Game.DestroyProjectile()

	// CALCULATE DELTATIME
	if time.Since(startDelta).Milliseconds() == 0 {
		time.Sleep(time.Millisecond * 1)
	}
	// fmt.Printf("Game processing for %d ms\n", time.Since(startDelta).Milliseconds())
	Game.DeltaTime = float32(time.Since(startDelta).Milliseconds()) / 1000
	sync <- true
}

func GameSpawners(Game *entities.Game, sync chan bool) {
	for {
		t := <-Game.Ticker.C
		if t.Second()%2 == 0 {
			Game.CreateMobB()
		}
		if t.Second()%3 == 0 {
			Game.CreateMobC()
		}
		Game.CreateMobA()
		<-sync
	}
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
		sync <- true
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		defer func() {
			Game.DestroyPlayer(id)
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
				if players.Uid == id {
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
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	// GAME
	var Game = entities.Game{
		DeltaTime: 0.0,
		Ticker:    *time.NewTicker(time.Millisecond * 500),
	}
	sync := make(chan bool)

	// LOOP
	go Server(&Game, sync)
	go GameSpawners(&Game, sync)
	for {
		if len(Game.Players) > 0 {
			GameLoop(&Game, sync)
		}
	}
}
