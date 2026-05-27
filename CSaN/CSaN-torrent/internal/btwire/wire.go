package btwire

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	MsgChoke = iota
	MsgUnchoke
	MsgInterested
	MsgNotInterested
	MsgHave
	MsgBitfield
	MsgRequest
	MsgPiece
	MsgCancel
)

// BlockSize is the default block size for BitTorrent requests.
const BlockSize = 16384

// Handshake performs the BitTorrent handshake and verifies info_hash.
func Handshake(conn io.ReadWriter, infoHash, peerID [20]byte) error {
	buf := make([]byte, 68)
	buf[0] = 19
	copy(buf[1:20], "BitTorrent protocol")
	// 8 reserved bytes already zero
	copy(buf[28:48], infoHash[:])
	copy(buf[48:68], peerID[:])
	if _, err := conn.Write(buf); err != nil {
		return err
	}
	if _, err := io.ReadFull(conn, buf); err != nil {
		return err
	}
	if buf[0] != 19 || string(buf[1:20]) != "BitTorrent protocol" {
		return fmt.Errorf("invalid handshake response")
	}
	var remoteIH [20]byte
	copy(remoteIH[:], buf[28:48])
	if remoteIH != infoHash {
		return fmt.Errorf("info_hash mismatch in handshake")
	}
	return nil
}

// AcceptHandshake reads the initiator handshake, verifies info_hash, then sends ours.
func AcceptHandshake(conn io.ReadWriter, infoHash, peerID [20]byte) error {
	buf := make([]byte, 68)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return err
	}
	if buf[0] != 19 || string(buf[1:20]) != "BitTorrent protocol" {
		return fmt.Errorf("invalid handshake")
	}
	var remoteIH [20]byte
	copy(remoteIH[:], buf[28:48])
	if remoteIH != infoHash {
		return fmt.Errorf("info_hash mismatch in handshake")
	}
	out := make([]byte, 68)
	out[0] = 19
	copy(out[1:20], "BitTorrent protocol")
	copy(out[28:48], infoHash[:])
	copy(out[48:68], peerID[:])
	if _, err := conn.Write(out); err != nil {
		return err
	}
	return nil
}

// ReadMessage reads one peer wire message. keepalive is true when length prefix is 0.
func ReadMessage(r io.Reader) (id byte, payload []byte, keepalive bool, err error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return 0, nil, false, err
	}
	n := binary.BigEndian.Uint32(lenBuf[:])
	if n == 0 {
		return 0, nil, true, nil
	}
	if n > 1<<24 {
		return 0, nil, false, fmt.Errorf("message too large: %d", n)
	}
	payload = make([]byte, n)
	if _, err := io.ReadFull(r, payload); err != nil {
		return 0, nil, false, err
	}
	return payload[0], payload[1:], false, nil
}

// WriteMessage writes a peer wire message (id + payload).
func WriteMessage(w io.Writer, id byte, payload []byte) error {
	n := uint32(1 + len(payload))
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], n)
	if _, err := w.Write(lenBuf[:]); err != nil {
		return err
	}
	if _, err := w.Write([]byte{id}); err != nil {
		return err
	}
	if len(payload) > 0 {
		if _, err := w.Write(payload); err != nil {
			return err
		}
	}
	return nil
}

// EncodeRequest builds payload for MsgRequest.
func EncodeRequest(pieceIndex int, begin int, length int) []byte {
	out := make([]byte, 12)
	binary.BigEndian.PutUint32(out[0:4], uint32(pieceIndex))
	binary.BigEndian.PutUint32(out[4:8], uint32(begin))
	binary.BigEndian.PutUint32(out[8:12], uint32(length))
	return out
}

// DecodeRequest parses MsgRequest payload.
func DecodeRequest(payload []byte) (index int, begin int, length int, err error) {
	if len(payload) != 12 {
		return 0, 0, 0, fmt.Errorf("invalid request length")
	}
	index = int(binary.BigEndian.Uint32(payload[0:4]))
	begin = int(binary.BigEndian.Uint32(payload[4:8]))
	length = int(binary.BigEndian.Uint32(payload[8:12]))
	if length < 1 || length > 1<<17 {
		return 0, 0, 0, fmt.Errorf("invalid request block length")
	}
	return index, begin, length, nil
}

// FullBitfield returns a bitfield marking all numPieces as available.
func FullBitfield(numPieces int) []byte {
	if numPieces < 1 {
		return nil
	}
	nbytes := (numPieces + 7) / 8
	bf := make([]byte, nbytes)
	for p := 0; p < numPieces; p++ {
		bf[p/8] |= 1 << (7 - (p % 8))
	}
	return bf
}

// EncodePiece builds MsgPiece payload (index, begin, block).
func EncodePiece(index, begin int, block []byte) []byte {
	out := make([]byte, 8+len(block))
	binary.BigEndian.PutUint32(out[0:4], uint32(index))
	binary.BigEndian.PutUint32(out[4:8], uint32(begin))
	copy(out[8:], block)
	return out
}

// DecodePiece parses MsgPiece payload: index, begin, block.
func DecodePiece(payload []byte) (index int, begin int, block []byte, err error) {
	if len(payload) < 8 {
		return 0, 0, nil, fmt.Errorf("piece message too short")
	}
	index = int(binary.BigEndian.Uint32(payload[0:4]))
	begin = int(binary.BigEndian.Uint32(payload[4:8]))
	block = append([]byte(nil), payload[8:]...)
	return index, begin, block, nil
}
