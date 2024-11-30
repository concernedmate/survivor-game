package entities

import (
	"math/rand"
	"time"
)

const MAP_BOUNDARY_X = 1280
const MAP_BOUNDARY_Y = 720

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

	DeleteFlag bool
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

	DeleteFlag bool
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

func (Game *Game) DestroyProjectile() {
	var newArr []Projectile
	for _, proj := range Game.Projectiles {
		if !proj.DeleteFlag {
			newArr = append(newArr, proj)
		}
	}
	Game.Projectiles = newArr
}

func (Game *Game) DestroyMob() {
	var newArr []Mob
	for _, mob := range Game.Mobs {
		if !mob.DeleteFlag {
			newArr = append(newArr, mob)
		}
	}
	Game.Mobs = newArr
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
		PosX:  float32(rand.Intn(MAP_BOUNDARY_X/4)) + float32(rand.Intn(MAP_BOUNDARY_X/4)),
		PosY:  MAP_BOUNDARY_Y - 50,
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
		DeleteFlag: false,
	})
	Player.WeaponCD = 150
}

func (Game *Game) CreateMobA() {
	newmob := Mob{
		PosX:  float32(rand.Intn(MAP_BOUNDARY_X)),
		PosY:  10,
		Size:  35,
		Speed: 300,

		Health: 1,
		Range:  0,

		AttackRange: 0,
		AttackAoe:   0,
		MobMovement: func(mob *Mob, deltaTime float32) {
			mob.PosY += float32(mob.Speed) * deltaTime
		},
		DeleteFlag: false,
	}
	Game.Mobs = append(Game.Mobs, newmob)
}

func (Game *Game) CreateMobB() {
	startX := MAP_BOUNDARY_X/4 + float32(rand.Intn(MAP_BOUNDARY_X/2))
	newmob := Mob{
		PosX:  startX,
		PosY:  10,
		Size:  40,
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
		DeleteFlag: false,
	}
	Game.Mobs = append(Game.Mobs, newmob)
}

func (Game *Game) CreateMobC() {
	startX := MAP_BOUNDARY_X/4 + float32(rand.Intn(MAP_BOUNDARY_X/4*3))
	newmob1 := Mob{
		PosX:  startX,
		PosY:  10,
		Size:  50,
		Speed: 500,

		Health: 1,
		Range:  0,

		AttackRange: 0,
		AttackAoe:   0,
		MobMovement: func(mob *Mob, deltaTime float32) {
			mob.PosY += float32(mob.Speed) * deltaTime
		},
		DeleteFlag: false,
	}
	newmob2 := Mob{
		PosX:  startX,
		PosY:  10,
		Size:  40,
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
		DeleteFlag: false,
	}
	newmob3 := Mob{
		PosX:  startX,
		PosY:  10,
		Size:  40,
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
		DeleteFlag: false,
	}
	Game.Mobs = append(Game.Mobs, newmob1)
	Game.Mobs = append(Game.Mobs, newmob2)
	Game.Mobs = append(Game.Mobs, newmob3)
}
