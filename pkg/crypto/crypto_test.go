package crypto

import "testing"

func TestEncodeDecode(t *testing.T) {
	crypt := NewCrypto("testfuck", "testf")

	encoded := crypt.Encode([]byte("test"))
	t.Log(string(encoded))
	encoded = crypt.Encode([]byte("test"))
	t.Log(string(encoded))
	t.Log(string(crypt.Decode(encoded)))

}
