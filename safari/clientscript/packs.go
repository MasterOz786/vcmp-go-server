package clientscript

import "github.com/masteroz/vcmp-go-server/safari/stream"

func showPacket(packetID int32) []byte {
	w := stream.NewWriter()
	w.WriteInt(packetID)
	return w.Bytes()
}

func ShowRegister() []byte  { return showPacket(stream.PacketShowRegister) }
func HideRegister() []byte  { return showPacket(stream.PacketHideRegister) }
func HidePacks() []byte     { return showPacket(stream.PacketHidePacks) }

func ShowPacks(team, currentPack int) []byte {
	w := stream.NewWriter()
	w.WriteInt(stream.PacketShowPacks)
	w.WriteInt(int32(team))
	w.WriteInt(int32(currentPack))
	return w.Bytes()
}

func PackFeedback(message string) []byte {
	w := stream.NewWriter()
	w.WriteInt(stream.PacketPackFeedback)
	w.WriteString(message)
	return w.Bytes()
}

func HydraCam(mode, vehicleID int) []byte {
	w := stream.NewWriter()
	w.WriteInt(stream.PacketHydraCam)
	w.WriteInt(int32(mode))
	w.WriteInt(int32(vehicleID))
	return w.Bytes()
}
