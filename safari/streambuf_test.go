package safari

import "testing"

func TestStreamRoundTrip(t *testing.T) {
	w := NewStreamWriter()
	w.WriteInt(PacketRegister)
	w.WriteString("test-pass")

	r := NewStreamReader(w.Bytes())
	pkt, err := r.ReadInt()
	if err != nil || pkt != PacketRegister {
		t.Fatalf("ReadInt() = (%d, %v), want (%d, nil)", pkt, err, PacketRegister)
	}
	pass, err := r.ReadString()
	if err != nil || pass != "test-pass" {
		t.Fatalf("ReadString() = (%q, %v), want (test-pass, nil)", pass, err)
	}
}

func TestStreamShowRegisterPacket(t *testing.T) {
	w := NewStreamWriter()
	w.WriteInt(PacketShowRegister)
	if len(w.Bytes()) != 4 {
		t.Fatalf("SHOW_REGISTER payload len = %d, want 4", len(w.Bytes()))
	}
}
