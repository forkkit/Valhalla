package entity

import (
	"fmt"

	"github.com/Hucaru/Valhalla/mnet"
	"github.com/Hucaru/Valhalla/mpacket"
)

type instance struct {
	id      int
	fieldID int32
	npcs    []npc
	conns   []mnet.Client
	players *Players
}

func (inst *instance) delete() error {
	return nil
}

func (inst *instance) addPlayer(conn mnet.Client) error {
	for _, npc := range inst.npcs {
		conn.Send(PacketNpcShow(npc))
		conn.Send(PacketNpcSetController(npc.spawnID, true))
	}

	connPlayer, _ := inst.players.GetFromConn(conn)
	for _, other := range inst.conns {
		otherPlayer, _ := inst.players.GetFromConn(other)
		other.Send(PacketMapPlayerEnter(connPlayer.char))
		conn.Send(PacketMapPlayerEnter(otherPlayer.char))
	}

	// show all monsters on field

	// show all the rooms

	inst.conns = append(inst.conns, conn)
	return nil
}

func (inst *instance) removePlayer(conn mnet.Client) error {
	index := -1

	for i, v := range inst.conns {
		if v == conn {
			index = i
			break
		} else {

		}
	}
	if index == -1 {
		return fmt.Errorf("player does not exist in instance")
	}

	inst.conns = append(inst.conns[:index], inst.conns[index+1:]...)

	// if in room, remove
	player, _ := inst.players.GetFromConn(conn)
	for _, v := range inst.conns {
		v.Send(PacketMapPlayerLeft(player.char.id))
	}

	return nil
}

func (inst instance) send(p mpacket.Packet) error {
	for _, v := range inst.conns {
		v.Send(p)
	}

	return nil
}

func (inst instance) sendExcept(p mpacket.Packet, exception mnet.Client) error {
	for _, v := range inst.conns {
		if v == exception {
			continue
		}

		v.Send(p)
	}

	return nil
}