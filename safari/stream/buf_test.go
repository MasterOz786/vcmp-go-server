package stream

import "testing"

func TestStreamRoundTrip(t *testing.T) {
	w := NewWriter()
	w.WriteInt(PacketRegister)
	w.WriteString("secret")

	r := NewReader(w.Bytes())
	pkt, err := r.ReadInt()
	if err != nil || pkt != PacketRegister {
		t.Fatalf("ReadInt() = (%d, %v), want (%d, nil)", pkt, err, PacketRegister)
	}
	s, err := r.ReadString()
	if err != nil || s != "secret" {
		t.Fatalf("ReadString() = (%q, %v), want (secret, nil)", s, err)
	}
}

func TestShowRegisterPacket(t *testing.T) {
	w := NewWriter()
	w.WriteInt(PacketShowRegister)
	if len(w.Bytes()) != 4 {
		t.Fatalf("SHOW_REGISTER payload len = %d, want 4", len(w.Bytes()))
	}
}
