package game

import (
	"fmt"

	"github.com/Hucaru/Valhalla/game/def"
	"github.com/Hucaru/Valhalla/game/packet"
	"github.com/Hucaru/Valhalla/mnet"
	"github.com/Hucaru/Valhalla/mpacket"
)

type roomContainer map[int32]*Room

var Rooms = make(roomContainer)

const (
	OmokRoom         = 0x01
	MemoryRoom       = 0x02
	TradeRoom        = 0x03
	PersonalShop     = 0x04
	omokMaxPlayers   = 4
	memoryMaxPlayers = 4
	tradeMaxPlayers  = 2
)

type Room struct {
	ID       int32
	RoomType byte
	players  []mnet.MConnChannel

	Name, Password string
	inProgress     bool
	board          [15][15]byte
	cards          []byte
	BoardType      byte
	leaveAfterGame [2]bool

	accepted int
	items    [2][9]Item
	mesos    [2]int32

	maxPlayers int
}

var roomCounter = int32(0)

func (rc *roomContainer) getNewRoomID() int32 {
	roomCounter++

	if roomCounter == 0 {
		roomCounter = 1
	}

	return roomCounter
}

func (rc *roomContainer) CreateMemoryRoom(name, password string, boardType byte) int32 {
	id := rc.getNewRoomID()
	r := &Room{ID: id, RoomType: MemoryRoom, Name: name, Password: password, BoardType: boardType, maxPlayers: memoryMaxPlayers}
	Rooms[id] = r
	return id
}

func (rc *roomContainer) CreateOmokRoom(name, password string, boardType byte) int32 {
	id := rc.getNewRoomID()
	r := &Room{ID: id, RoomType: OmokRoom, Name: name, Password: password, BoardType: boardType, maxPlayers: omokMaxPlayers}
	Rooms[id] = r
	return id
}

func (rc *roomContainer) CreateTradeRoom() int32 {
	id := rc.getNewRoomID()
	r := &Room{ID: id, RoomType: TradeRoom, maxPlayers: tradeMaxPlayers}
	Rooms[id] = r
	return id
}

func (r *Room) IsOwner(conn mnet.MConnChannel) bool {
	if len(r.players) > 0 && r.players[0] == conn {
		return true
	}

	return false
}

func (r *Room) Broadcast(p mpacket.Packet) {
	for _, v := range r.players {
		v.Send(p)
	}
}

func (r *Room) AddPlayer(conn mnet.MConnChannel) {
	if len(r.players) == r.maxPlayers {
		conn.Send(packet.RoomFull())
	}

	r.players = append(r.players, conn)

	player := Players[conn]
	roomPos := byte(len(r.players)) - 1

	if roomPos == 0 {
		Maps[player.Char().MapID].Send(packet.MapShowGameBox(player.Char().ID, r.ID, r.RoomType, r.BoardType, r.Name, bool(len(r.Password) > 0), r.inProgress, 0x01), player.InstanceID)
	}

	r.Broadcast(packet.RoomJoin(r.RoomType, roomPos, player.Char()))

	displayInfo := []def.Character{}

	for _, v := range r.players {
		displayInfo = append(displayInfo, Players[v].Char())
	}

	if len(displayInfo) > 0 {
		conn.Send(packet.RoomShowWindow(r.RoomType, r.BoardType, byte(r.maxPlayers), roomPos, r.Name, displayInfo))
		// update box on map
	}
}

func (r *Room) closeRoom() {

}

func (r *Room) RemovePlayer(conn mnet.MConnChannel, msgCode byte) bool {
	closeRoom := false
	roomSlot := -1

	for i, v := range r.players {
		if v == conn {
			roomSlot = i
			break
		}
	}

	if roomSlot < 0 {
		return false
	}

	r.players = append(r.players[:roomSlot], r.players[roomSlot+1:]...)

	switch r.RoomType {
	case TradeRoom:
		if r.accepted > 0 {
			r.Broadcast(packet.RoomLeave(byte(roomSlot), 7))
		} else {
			r.Broadcast(packet.RoomLeave(byte(roomSlot), 2))
		}

		closeRoom = true
	case MemoryRoom:
		fallthrough
	case OmokRoom:
		player := Players[conn]

		if roomSlot == 0 {
			closeRoom = true

			Maps[player.Char().MapID].Send(packet.MapRemoveGameBox(player.Char().ID), player.InstanceID)

			for i, v := range r.players {
				v.Send(packet.RoomLeave(byte(i), 0))
			}
		} else {
			conn.Send(packet.RoomLeave(byte(roomSlot), msgCode))
			r.Broadcast(packet.RoomLeave(byte(roomSlot), msgCode))

			if msgCode == 5 {

				r.Broadcast(packet.RoomYellowChat(0, player.Char().Name))
			}

			// Update player positions from index roomSlot onwards (not + 1 as we have removed the gone player)
		}
	default:
		fmt.Println("have not implemented remove player for room type", r.RoomType)
	}

	return closeRoom
}