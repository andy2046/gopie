package tlv

import (
	"bytes"
	"testing"
)

func TestTLV(t *testing.T) {
	v := "hola, tlv!"
	typ := uint(8)
	buf := new(bytes.Buffer)
	codec := &Codec{TypeBytes: Bytes1, LenBytes: Bytes2}
	writer := NewWriter(buf, codec)

	record := &Record{
		Payload: []byte(v),
		Type:    typ,
	}
	writer.Write(record)

	reader := bytes.NewReader(buf.Bytes())
	tlvReader := NewReader(reader, codec)
	next, _ := tlvReader.Next()

	if next.Type != typ {
		t.Errorf("expected %d got %d", typ, next.Type)
	}

	if r := string(next.Payload); r != v {
		t.Errorf("expected %s got %s", v, r)
	}

	t.Logf("type: %d\n", next.Type)
	t.Logf("payload: %s\n", string(next.Payload))
}
