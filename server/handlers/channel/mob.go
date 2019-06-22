package channel

import (
	"math/rand"
	"time"

	skills "github.com/Hucaru/Valhalla/constant/skill"
	"github.com/Hucaru/Valhalla/game"
	"github.com/Hucaru/Valhalla/mnet"
	"github.com/Hucaru/Valhalla/mpacket"
	"github.com/Hucaru/Valhalla/nx"
)

func mobControl(conn mnet.Client, reader mpacket.Reader) {
	mobSpawnID := reader.ReadInt32()
	moveID := reader.ReadInt16()
	bits := reader.ReadByte()
	action := reader.ReadInt8()
	skillData := reader.ReadUint32()

	actualAction := int(byte(action >> 1))

	if action < 0 {
		actualAction = -1
	}

	skillPossible := (bits & 0x0F) != 0

	player, ok := game.Players[conn]

	if !ok {
		return
	}

	char := player.Char()
	game.Maps[char.MapID].GetMobFromSpawnID(mobSpawnID, player.InstanceID)
	mob, err := game.Maps[char.MapID].GetMobFromSpawnID(mobSpawnID, player.InstanceID)

	if err != nil {
		return
	}

	if mob.Controller != conn { // prevents hijack and reassigns controller to anyone except hijacker
		if newController := game.Maps[char.MapID].FindControllerExcept(conn, player.InstanceID); newController != nil {
			mob.ChangeController(newController)
		}

		return
	}

	// Update mob position information
	moveData, finalData := parseMovement(reader)

	if !validateMobMovement(*mob, moveData) {
		return
	}

	mob.X = finalData.X
	mob.Y = finalData.Y
	mob.Foothold = finalData.Foothold
	mob.Stance = finalData.Stance

	moveBytes := generateMovementBytes(moveData)

	// Perform the action received
	if actualAction >= 21 && actualAction <= 25 {
		performSkill(mob, int16(skillData>>16), byte(skillData>>8), byte(skillData))
	} else if actualAction > 12 && actualAction < 20 {
		attackID := byte(actualAction - 12)

		// check mob can use attack
		if level, valid := mob.Skills[attackID]; valid {
			levels, err := nx.GetMobSkill(attackID)

			if err != nil {
				return
			}

			if int(level) < len(levels) {
				skill := levels[level]
				mob.MP = mob.MP - skill.MpCon
				if mob.MP < 0 {
					mob.MP = 0
				}
			}

		}

		mob.LastAttackTime = time.Now().Unix()
	}

	// Calculate the next action
	mob.CanUseSkill = skillPossible

	if !mob.CanUseSkill || (mob.StatBuff&skills.MobStat.SealSkill > 0) || (time.Now().Unix()-mob.LastSkillUseTime) < 3 {
		// there are more reasons as to why a mob cannot use a skill
		mob.SkillID = 0
	} else {
		mob.SkillID, mob.SkillLevel = chooseNextSkill(mob)
	}

	mob.Acknowledge(moveID, skillPossible, mob.SkillID, mob.SkillLevel)
	game.Maps[char.MapID].SendExcept(game.PacketMobMove(mobSpawnID, skillPossible, byte(action), skillData, moveBytes), conn, player.InstanceID)
}

func chooseNextSkill(mob *game.Mob) (byte, byte) {
	var skillID, skillLevel byte

	skillsToChooseFrom := []byte{}

	for id, level := range mob.Skills {
		levels, err := nx.GetMobSkill(level)

		if err != nil {
			continue
		}

		if int(skillLevel) >= len(levels) {
			continue
		}

		skillData := levels[skillLevel]

		// Skill HP check
		if (mob.HP * 100 / mob.MaxHP) < skillData.Hp {
			continue
		}

		// Skill cooldown check
		if val, ok := mob.SkillTimes[id]; ok {
			if (val + skillData.Interval) > time.Now().Unix() { // Is cooldown in seconds?
				continue
			}
		}

		// Check summon limit
		// if skillData.Limit {

		// }

		// Determine if stats can be buffed
		if mob.StatBuff > 0 {
			alreadySet := false

			switch id {
			case skills.Mob.WeaponAttackUp:
				fallthrough
			case skills.Mob.WeaponAttackUpAoe:
				alreadySet = mob.StatBuff&skills.MobStat.PowerUp > 0

			case skills.Mob.MagicAttackUp:
				fallthrough
			case skills.Mob.MagicAttackUpAoe:
				alreadySet = mob.StatBuff&skills.MobStat.MagicUp > 0

			case skills.Mob.WeaponDefenceUp:
				fallthrough
			case skills.Mob.WeaponDefenceUpAoe:
				alreadySet = mob.StatBuff&skills.MobStat.PowerGuardUp > 0

			case skills.Mob.MagicDefenceUp:
				fallthrough
			case skills.Mob.MagicDefenceUpAoe:
				alreadySet = mob.StatBuff&skills.MobStat.MagicGuardUp > 0

			case skills.Mob.WeaponImmunity:
				alreadySet = mob.StatBuff&skills.MobStat.PhysicalImmune > 0

			case skills.Mob.MagicImmunity:
				alreadySet = mob.StatBuff&skills.MobStat.MagicImmune > 0

			// case skills.Mob.WeaponDamageReflect:

			// case skills.Mob.MagicDamageReflect:

			case skills.Mob.McSpeedUp:
				alreadySet = mob.StatBuff&skills.MobStat.Speed > 0

			default:
			}

			if alreadySet {
				continue
			}

		}

		skillsToChooseFrom = append(skillsToChooseFrom, id)
	}

	if len(skillsToChooseFrom) > 0 {
		nextID := skillsToChooseFrom[rand.Intn(len(skillsToChooseFrom))]

		skillID = nextID

		for id, level := range mob.Skills {
			if id == nextID {
				skillLevel = level
			}
		}
	}

	if skillLevel == 0 {
		skillID = 0
	}

	return skillID, skillLevel
}

func performSkill(mob *game.Mob, delay int16, skillLevel, skillID byte) {
	if skillID != mob.SkillID || (mob.StatBuff&skills.MobStat.SealSkill > 0) {
		skillID = 0
		return
	}

	levels, err := nx.GetMobSkill(skillID)

	if err != nil {
		mob.SkillID = 0
		return
	}

	var skillData nx.MobSkill
	for i, v := range levels {
		if i == int(skillLevel) {
			skillData = v
		}
	}

	mob.MP = mob.MP - skillData.MpCon
	if mob.MP < 0 {
		mob.MP = 0
	}

	currentTime := time.Now().Unix()

	mob.SkillTimes[skillID] = currentTime
	mob.LastSkillUseTime = currentTime

	// Handle all the different skills!
	switch skillID {

	}
}