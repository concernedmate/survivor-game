package entities

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

	Health   int
	Lifetime int
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
