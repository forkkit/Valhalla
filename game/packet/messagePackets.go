package packet

import (
	"github.com/Hucaru/Valhalla/consts/opcodes"
	"github.com/Hucaru/Valhalla/mpacket"
)

func MessageRedText(msg string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelInfoMessage)
	p.WriteByte(9)
	p.WriteString(msg)

	return p
}

func MessageGuildPointsChange(ammount int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelInfoMessage)
	p.WriteByte(6)
	p.WriteInt32(ammount)

	return p
}

func MessageFameChange(ammount int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelInfoMessage)
	p.WriteByte(4)
	p.WriteInt32(ammount)

	return p
}

// sends the [item name] has passed its expeiration date and will be removed from your inventory
func MessageItemExpired(itemID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelInfoMessage)
	p.WriteByte(2)
	p.WriteInt32(itemID)
	return p
}

func MessageItemExpired2(itemID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelInfoMessage)
	p.WriteByte(8)
	p.WriteByte(1)
	p.WriteInt32(itemID)
	return p
}

func MessageMesosChangeChat(ammount int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelInfoMessage)
	p.WriteByte(5)
	p.WriteInt32(ammount)

	return p
}

func MessageUnableToPickUp(itemNotAvailable bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelInfoMessage)
	p.WriteByte(0)
	if itemNotAvailable {
		p.WriteByte(0xFE)
	} else {
		p.WriteByte(0xFF)
	}

	return p
}

func MessageDropPickUp(isMesos bool, itemID, ammount int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelInfoMessage)
	p.WriteByte(0)

	if isMesos {
		p.WriteInt32(ammount)
		p.WriteInt32(0)
	} else {
		p.WriteInt32(itemID)
		p.WriteInt32(ammount)
	}

	return p
}

func MessageExpGained(whiteText, appearInChat bool, ammount int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelInfoMessage)
	p.WriteByte(3)
	p.WriteBool(whiteText)
	p.WriteInt32(ammount)
	p.WriteBool(appearInChat)

	return p
}

func MessageNotice(msg string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelBroadcastMessage)
	p.WriteByte(0)
	p.WriteString(msg)

	return p
}

func MessageDialogueBox(msg string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelBroadcastMessage)
	p.WriteByte(1)
	p.WriteString(msg)

	return p
}

func MessageWhiteBar(msg string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelBroadcastMessage)
	p.WriteByte(2)
	p.WriteString(msg) // not sure how string is formated

	return p
}

// Need to figure out how to display the username and  atm it bastardises it.
func MessageBroadcastChannel(senderName string, msg string, channel byte, ear bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelBroadcastMessage)
	p.WriteByte(3)
	p.WriteString(senderName)
	p.WriteByte(channel)
	if ear {
		p.WriteByte(0x01)
	} else {
		p.WriteByte(0x00)
	}

	return p
}

func MessageScrollingHeader(msg string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelBroadcastMessage)
	p.WriteByte(4)
	p.WriteBool(bool(len(msg) > 0))
	p.WriteString(msg)

	return p
}

func MessageBubblessChat(msgType byte, sender string, msg string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelBubblessChat)
	p.WriteByte(msgType) // 0x00 buddy chat, 0x01 - party, 0x02 - guild
	p.WriteString(sender)
	p.WriteString(msg)

	return p
}

func MessageWhisper(sender string, message string, channel byte) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelWhisper)
	p.WriteByte(0x12)
	p.WriteString(sender)
	p.WriteByte(channel)
	p.WriteString(message)

	return p
}

func MessageFindResult(character string, isAdmin, inCashShop, sameChannel bool, mapID int32) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelWhisper)

	if isAdmin {
		p.WriteByte(0x05)
		p.WriteString("User not found")

	} else if mapID > 0 {
		p.WriteByte(0x9)
		p.WriteString(character)

		if inCashShop {
			p.WriteByte(0x02)
			p.WriteInt32(0) // ?
		} else if sameChannel {
			p.WriteByte(0x01)
			p.WriteInt32(mapID)
			p.WriteInt32(0) // ?
		} else {
			p.WriteByte(0x01)
			p.WriteInt32(mapID)
		}

		p.WriteInt32(0) // ?
	} else {
		p.WriteByte(0x0A)
		p.WriteString(character)
		p.WriteByte(0) // ?
	}

	return p
}

func MessageAllChat(senderID int32, isAdmin bool, msg string) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelAllChatMsg)
	p.WriteInt32(senderID)
	p.WriteBool(isAdmin)
	p.WriteString(msg)

	return p
}

// Implement logic for these
func MessageGmBan(good bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelEmployee)
	if good {
		p.WriteByte(4)
		p.WriteByte(0)
	} else {
		p.WriteByte(6)
		p.WriteByte(1)
	}

	return p
}

func MessageGmRemoveFromRanks() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelEmployee)
	p.WriteByte(6)
	p.WriteByte(0)

	return p
}

func MessageGmWarning(good bool) mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelEmployee)
	p.WriteByte(14)
	if good {
		p.WriteByte(1)
	} else {
		p.WriteByte(0)
	}

	return p
}

func MessageGmBlockedAccess() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelEmployee)
	p.WriteByte(4)
	p.WriteByte(0)

	return p
}

func MessageGmUnblock() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelEmployee)
	p.WriteByte(5)
	p.WriteByte(0)

	return p
}

// Don't know what this is used for
func MessageGmWrongNpc() mpacket.Packet {
	p := mpacket.CreateWithOpcode(opcodes.SendChannelEmployee)
	p.WriteByte(8)
	p.WriteInt16(0)

	return p
}