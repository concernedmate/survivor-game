package engines

import (
	"concernedmate/SurvivorGame/entities"
	"concernedmate/SurvivorGame/physics"
	"fmt"
	"log"
	"time"
)

const MAP_BOUNDARY_X = entities.MAP_BOUNDARY_X
const MAP_BOUNDARY_Y = entities.MAP_BOUNDARY_Y

func GameRoutine(Game *entities.Game, id_room string) {
	timer := time.Now()
	for {
		if len(Game.Players) > 0 {
			GameLoop(Game)
			timer = time.Now()
		} else {
			if time.Since(timer).Seconds() >= 3 {
				fmt.Printf("closing room %s\n", id_room)
				delete(Rooms, id_room)
				break
			}
		}
	}
}

func GameLoop(Game *entities.Game) {
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
					Game.AddScorePlayer(Game.Projectiles[i].OwnerId, Game.Mobs[j].ScoreVal)
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

	Game.DeltaTime = float32(time.Since(startDelta).Milliseconds()) / 1000
	Game.Sync <- true
}

func GameSpawners(Game *entities.Game) {
	for Game != nil {
		t := <-Game.Ticker.C
		if t.Second()%2 == 0 {
			Game.CreateMobB()
		}
		if t.Second()%3 == 0 {
			Game.CreateMobC()
		}
		Game.CreateMobA()
		<-Game.Sync
	}
}
