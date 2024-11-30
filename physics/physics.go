package physics

import "concernedmate/SurvivorGame/entities"

func PlayerMobCollision(player *entities.Player, mob *entities.Mob) bool {
	if player.PosX > mob.PosX-float32(mob.Size/2) && player.PosX < mob.PosX+float32(mob.Size/2) &&
		player.PosY > mob.PosY-float32(mob.Size/2) && player.PosY < mob.PosY+float32(mob.Size/2) {
		return true
	}
	if mob.PosX > player.PosX-float32(player.Size/2) && mob.PosX < player.PosX+float32(player.Size/2) &&
		mob.PosY > player.PosY-float32(player.Size/2) && mob.PosY < player.PosY+float32(player.Size/2) {
		return true
	}
	return false
}

func ProjectileMobCollision(proj *entities.Projectile, mob *entities.Mob) bool {
	if proj.PosX > mob.PosX-float32(mob.Size/2) && proj.PosX < mob.PosX+float32(mob.Size/2) &&
		proj.PosY > mob.PosY-float32(mob.Size/2) && proj.PosY < mob.PosY+float32(mob.Size/2) {
		return true
	}
	if mob.PosX > proj.PosX-float32(proj.Size/2) && mob.PosX < proj.PosX+float32(proj.Size/2) &&
		mob.PosY > proj.PosY-float32(proj.Size/2) && mob.PosY < proj.PosY+float32(proj.Size/2) {
		return true
	}
	return false
}
