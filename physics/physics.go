package physics

import "concernedmate/SurvivorGame/entities"

func PlayerMobCollision(player *entities.Player, mob *entities.Mob) bool {
	// TODO
	if player.PosX > mob.PosX-float32(mob.Size/2) && player.PosX < mob.PosX+float32(mob.Size/2) &&
		player.PosY > mob.PosY-float32(mob.Size/2) && player.PosY < mob.PosY+float32(mob.Size/2) {
		return true
	}
	return false
}

func ProjectileMobCollision(proj *entities.Projectile, mob *entities.Mob) bool {
	// TODO
	if proj.PosX > mob.PosX-float32(mob.Size/2) && proj.PosX < mob.PosX+float32(mob.Size/2) &&
		proj.PosY > mob.PosY-float32(mob.Size/2) && proj.PosY < mob.PosY+float32(mob.Size/2) {
		return true
	}
	return false
}
