package pwp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestSerializeMessage_MsgKeepAlive(t *testing.T) {
	msg := MsgKeepAlive{}
	result := SerializeMessage(msg)

	// Expected: 4 bytes length (value 1) + 1 byte ID (255)
	expected := []byte{0, 0, 0, 1, 255}

	if !bytes.Equal(result, expected) {
		t.Errorf("SerializeMessage(MsgKeepAlive) = %v, want %v", result, expected)
	}

	// Verify length field
	length := binary.BigEndian.Uint32(result[0:4])
	if length != 1 {
		t.Errorf("Length field = %d, want 1", length)
	}

	// Verify ID
	if result[4] != 255 {
		t.Errorf("ID = %d, want 255", result[4])
	}
}

func TestSerializeMessage_MsgChoke(t *testing.T) {
	msg := MsgChoke{}
	result := SerializeMessage(msg)

	// Expected: 4 bytes length (value 1) + 1 byte ID (Msg_Choke)
	expected := []byte{0, 0, 0, 1, Msg_Choke}

	if !bytes.Equal(result, expected) {
		t.Errorf("SerializeMessage(MsgChoke) = %v, want %v", result, expected)
	}

	// Verify total length
	if len(result) != 5 {
		t.Errorf("Total length = %d, want 5", len(result))
	}
}

func TestSerializeMessage_MsgUnchoke(t *testing.T) {
	msg := MsgUnchoke{}
	result := SerializeMessage(msg)

	expected := []byte{0, 0, 0, 1, Msg_Unchoke}

	if !bytes.Equal(result, expected) {
		t.Errorf("SerializeMessage(MsgUnchoke) = %v, want %v", result, expected)
	}
}

func TestSerializeMessage_MsgHave(t *testing.T) {
	tests := []struct {
		name  string
		index uint32
	}{
		{"zero index", 0},
		{"small index", 42},
		{"large index", 999999},
		{"max uint32", 0xFFFFFFFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgHave{Index: tt.index}
			result := SerializeMessage(msg)

			// Expected: 4 bytes length (value 5) + 1 byte ID + 4 bytes index
			if len(result) != 9 {
				t.Errorf("Total length = %d, want 9", len(result))
			}

			// Check length field
			length := binary.BigEndian.Uint32(result[0:4])
			if length != 5 {
				t.Errorf("Length field = %d, want 5", length)
			}

			// Check ID
			if result[4] != Msg_Have {
				t.Errorf("ID = %d, want %d", result[4], Msg_Have)
			}

			// Check index
			index := binary.BigEndian.Uint32(result[5:9])
			if index != tt.index {
				t.Errorf("Index = %d, want %d", index, tt.index)
			}
		})
	}
}

func TestSerializeMessage_MsgRequest(t *testing.T) {
	tests := []struct {
		name   string
		index  uint32
		begin  uint32
		length uint32
	}{
		{"all zeros", 0, 0, 0},
		{"typical request", 5, 16384, 16384},
		{"max values", 0xFFFFFFFF, 0xFFFFFFFF, 0xFFFFFFFF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgRequest{
				Index:  tt.index,
				Begin:  tt.begin,
				Length: tt.length,
			}
			result := SerializeMessage(msg)
			
			// Expected: 4 bytes length (value 13) + 1 byte ID + 12 bytes payload
			if len(result) != 17 {
				t.Errorf("Total length = %d, want 17", len(result))
			}

			// Check length field
			length := binary.BigEndian.Uint32(result[0:4])
			if length != 13 {
				t.Errorf("Length field = %d, want 13", length)
			}

			// Check ID
			if result[4] != Msg_Request {
				t.Errorf("ID = %d, want %d", result[4], Msg_Request)
			}

			// Check index
			index := binary.BigEndian.Uint32(result[5:9])
			if index != tt.index {
				t.Errorf("Index = %d, want %d", index, tt.index)
			}

			// Check begin
			begin := binary.BigEndian.Uint32(result[9:13])
			if begin != tt.begin {
				t.Errorf("Begin = %d, want %d", begin, tt.begin)
			}

			// Check length
			reqLength := binary.BigEndian.Uint32(result[13:17])
			if reqLength != tt.length {
				t.Errorf("Length = %d, want %d", reqLength, tt.length)
			}
		})
	}
}

func TestXx(t *testing.T) {
	// msg := MsgBitfield{
	// 	bits: []byte{255, 2, 3, 92, 122, 122,122,122,255, 2, 3, 92, 122,112,136,132,45,122},
	// }
	msg := MsgPiece{
		Index: 69,
		Begin: 999999999,
		Block: []byte{0x12, 0x34, 0x30, 0x39} ,
	}
	fmt.Println()
	fmt.Println("payload :", msg.Block)
	fmt.Println()

	result := SerializeMessage(msg)
	fmt.Println("res - ", result)
	fmt.Printf("%s", hex.Dump(result))

}

// Helper: encode and append a message
// func writeMessageStream(tmpfile *os.File, buf *[]byte, m Message) *os.File {
// 	payload := m.Payload()
// 	msg := append([]byte{m.ID()}, payload...)
// 	length := uint32(len(msg))

// 	lb := make([]byte, 4)
// 	binary.BigEndian.PutUint32(lb, length)

// 	*buf = append(*buf, lb...)
// 	*buf = append(*buf, msg...)

// 	tmpfile.Write(*buf)
// 	tmpfile.Seek(0, 0)
// 	return tmpfile
// }

// func TestSimulate(t *testing.T) {
// 	// 	// simulate input stream — writing a MsgHave message into a buffer
// 	buf := make([]byte, 0)
// 	tmpfile, _ := os.CreateTemp("", "msgbuf")
// 	writer := writeMessageStream(tmpfile, &buf, MsgHave{Index: 123})
// 	writer = writeMessageStream(tmpfile, &buf, MsgChoke{})
// 	writer = writeMessageStream(tmpfile, &buf, MsgBitfield{
// 		Bitfield: []byte{255, 2, 3, 92, 122, 122,122,122,255, 2, 3, 92, 122,112,136,132,45,122},
// 	})
// 	writer = writeMessageStream(tmpfile, &buf, MsgPiece{
// 		Index: 69,
// 		Begin: 999999999,
// 		Block: []byte{0x12, 0x34, 0x30, 0x39} ,
// 	})
	
// 	reader := bufio.NewReader(writer)
// 	for {
// 		msg, err := ReadMessage(reader)
// 		if err != nil {
// 			log.Fatal(err)
// 		}	
// 		fmt.Println("Received message")
// 		// fmt.Printf("type: %T\n", msg)
// 		fmt.Println("ID: ", msg.ID(), "Payload: ", msg.Payload())
// 	} 
// 	// switch m := msg.(type) {
// 	// case MsgHave:
// 	// 	log.Printf("Have: %d", m.Index)
// 	// default:
// 	// 	log.Printf("Unhandled message type")
// 	// }
// }

