package entities

import (
	"math/rand"
	"time"
)

type Player struct {
	Uid    string
	PosX   float32
	PosY   float32
	Size   int
	Speed  int
	Events []string

	Health int
	Mana   int

	Score    int
	WeaponCD int
}

type Mob struct {
	PosX  float32
	PosY  float32
	Size  int
	Speed int

	Health int
	Range  int

	AttackRange      int
	AttackAoe        int
	AttackProjectile string

	MobMovement    func(mob *Mob, deltaTime float32)
	MobMovModifier int
}

type Projectile struct {
	OwnerId  string
	Friendly bool
	PosX     float32
	PosY     float32
	Size     int
	Speed    int

	Damage          int
	DamageEvolution int

	Health       int
	Lifetime     int
	ProjMovement func(proj *Projectile, deltaTime float32)
}

type Obstacle struct {
	PosX float32
	PosY float32
	Size int
}

type Game struct {
	Players     []Player
	Mobs        []Mob
	Obstacles   []Obstacle
	Projectiles []Projectile

	DeltaTime float32
	Ticker    time.Ticker
}

func (Game *Game) DestroyProjectile(idx int) {
	if idx > len(Game.Projectiles)-1 || idx < 0 {
		panic("DestroyProjectile: projectile index invalid!")
	}
	Game.Projectiles = append(Game.Projectiles[:idx], Game.Projectiles[idx+1:]...)
}

func (Game *Game) DestroyMob(idx int) {
	if idx > len(Game.Mobs)-1 || idx < 0 {
		panic("DestroyMob: mob index invalid!")
	}
	Game.Mobs = append(Game.Mobs[:idx], Game.Mobs[idx+1:]...)
}

func (Game *Game) DestroyPlayer(uid string) {
	for idx, players := range Game.Players {
		if players.Uid == uid {
			Game.Players = append(Game.Players[:idx], Game.Players[idx+1:]...)
		}
	}
}

func (Game *Game) CreatePlayer(id string) {
	Game.Players = append(Game.Players, Player{
		Uid:   id,
		PosX:  float32(rand.Intn(250)) + float32(rand.Intn(250)),
		PosY:  950,
		Size:  30,
		Speed: 400,

		Health: 10,
		Mana:   10,

		Score: 0,
	})
}

func (Game *Game) CreateProjectileA(Player *Player) {
	Game.Projectiles = append(Game.Projectiles, Projectile{
		OwnerId: Player.Uid,
		PosX:    Player.PosX,
		PosY:    Player.PosY,
		Size:    10,
		Speed:   600,

		Damage: 0,
		ProjMovement: func(proj *Projectile, deltaTime float32) {
			proj.PosY -= float32(proj.Speed) * deltaTime
		},
	})
	Player.WeaponCD = 100
}

func (Game *Game) CreateMobA() {
	newmob := Mob{
		PosX:  20 + float32(rand.Intn(950)),
		PosY:  10,
		Size:  20,
		Speed: 300,

		Health: 1,
		Range:  0,

		AttackRange: 0,
		AttackAoe:   0,
		MobMovement: func(mob *Mob, deltaTime float32) {
			mob.PosY += float32(mob.Speed) * deltaTime
		},
	}
	Game.Mobs = append(Game.Mobs, newmob)
}

func (Game *Game) CreateMobB() {
	startX := 225 + float32(rand.Intn(575))
	newmob := Mob{
		PosX:  startX,
		PosY:  10,
		Size:  25,
		Speed: 400,

		Health: 1,
		Range:  0,

		AttackRange:    0,
		AttackAoe:      0,
		MobMovModifier: 1,
		MobMovement: func(mob *Mob, deltaTime float32) {
			mob.PosY += float32(mob.Speed) * deltaTime
			mob.PosX += float32(mob.Speed) * deltaTime * float32(mob.MobMovModifier)
			if mob.PosX > 900 || mob.PosX < 100 {
				mob.MobMovModifier *= -1
			}
		},
	}
	Game.Mobs = append(Game.Mobs, newmob)
}

func (Game *Game) CreateMobC() {
	startX := 125 + float32(rand.Intn(750))
	newmob1 := Mob{
		PosX:  startX,
		PosY:  10,
		Size:  35,
		Speed: 500,

		Health: 1,
		Range:  0,

		AttackRange: 0,
		AttackAoe:   0,
		MobMovement: func(mob *Mob, deltaTime float32) {
			mob.PosY += float32(mob.Speed) * deltaTime
		},
	}
	newmob2 := Mob{
		PosX:  startX,
		PosY:  10,
		Size:  25,
		Speed: 500,

		Health: 1,
		Range:  0,

		AttackRange:    0,
		AttackAoe:      0,
		MobMovModifier: 1,
		MobMovement: func(mob *Mob, deltaTime float32) {
			mob.PosY += float32(mob.Speed) * deltaTime
			mob.PosX += float32(mob.Speed) * deltaTime * float32(mob.MobMovModifier)
			if mob.PosX > startX+100 || mob.PosX < startX-100 {
				mob.MobMovModifier *= -1
			}
		},
	}
	newmob3 := Mob{
		PosX:  startX,
		PosY:  10,
		Size:  25,
		Speed: 500,

		Health: 1,
		Range:  0,

		AttackRange:    0,
		AttackAoe:      0,
		MobMovModifier: -1,
		MobMovement: func(mob *Mob, deltaTime float32) {
			mob.PosY += float32(mob.Speed) * deltaTime
			mob.PosX += float32(mob.Speed) * deltaTime * float32(mob.MobMovModifier)
			if mob.PosX > startX+100 || mob.PosX < startX-100 {
				mob.MobMovModifier *= -1
			}
		},
	}
	Game.Mobs = append(Game.Mobs, newmob1)
	Game.Mobs = append(Game.Mobs, newmob2)
	Game.Mobs = append(Game.Mobs, newmob3)
}
