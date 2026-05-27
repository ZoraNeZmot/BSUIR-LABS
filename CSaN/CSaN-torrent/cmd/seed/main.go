package main

import (
	"context"
	"crypto/rand"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"torrent/internal/btwire"
	"torrent/internal/metainfo"
	"torrent/internal/trackerclient"
)

func main() {
	torrentPath := flag.String("torrent", "", "path to .torrent (required)")
	contentPath := flag.String("content", "", "path to file on disk matching torrent size (required)")
	listenPort := flag.Int("port", 6882, "TCP port to listen on (must match tracker announce port)")
	trackerURL := flag.String("tracker", "", "optional announce URL override (e.g. http://127.0.0.1:8080/announce)")
	flag.Parse()

	if *torrentPath == "" || *contentPath == "" {
		log.Fatal("usage: seed -torrent FILE.torrent -content PAYLOAD.bin [-port N] [-tracker URL]")
	}

	t, err := metainfo.Load(*torrentPath)
	if err != nil {
		log.Fatalf("load torrent: %v", err)
	}
	announce := t.Announce
	if *trackerURL != "" {
		announce = *trackerURL
	}

	var peerID [20]byte
	copy(peerID[:8], []byte("-GoSeed1"))
	if _, err := rand.Read(peerID[8:]); err != nil {
		log.Fatalf("peer id: %v", err)
	}
	defer func() {
		_, _, _ = trackerclient.Announce(announce, t.InfoHash, peerID, *listenPort, t.TotalLen, t.TotalLen, 0, "stopped")
	}()

	fi, err := os.Stat(*contentPath)
	if err != nil {
		log.Fatalf("stat content: %v", err)
	}
	if fi.Size() != t.TotalLen {
		log.Fatalf("content size %d does not match torrent total %d", fi.Size(), t.TotalLen)
	}

	f, err := os.Open(*contentPath)
	if err != nil {
		log.Fatalf("open content: %v", err)
	}

	_, _, err = trackerclient.Announce(announce, t.InfoHash, peerID, *listenPort, t.TotalLen, 0, 0, "started")
	if err != nil {
		log.Fatalf("tracker announce: %v", err)
	}

	addr := "0.0.0.0:" + strconv.Itoa(*listenPort)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}
	defer ln.Close()
	log.Printf("seeding %q on %s (press Ctrl+C to stop)", t.Name, addr)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				log.Printf("accept: %v", err)
				continue
			}
		}
		go func(c net.Conn) {
			defer c.Close()
			if err := servePeer(c, t, f, peerID); err != nil && err != io.EOF {
				log.Printf("peer %s: %v", c.RemoteAddr(), err)
			}
		}(conn)
	}
}

func servePeer(conn net.Conn, t *metainfo.Torrent, content *os.File, peerID [20]byte) error {
	if err := btwire.AcceptHandshake(conn, t.InfoHash, peerID); err != nil {
		return err
	}
	if err := btwire.WriteMessage(conn, btwire.MsgBitfield, btwire.FullBitfield(t.NumPieces)); err != nil {
		return err
	}
	// Unchoke immediately so peers that already sent interested still work.
	if err := btwire.WriteMessage(conn, btwire.MsgUnchoke, nil); err != nil {
		return err
	}

	buf := make([]byte, btwire.BlockSize)
	for {
		id, payload, ka, err := btwire.ReadMessage(conn)
		if ka {
			continue
		}
		if err != nil {
			return err
		}
		switch id {
		case btwire.MsgInterested:
			if err := btwire.WriteMessage(conn, btwire.MsgUnchoke, nil); err != nil {
				return err
			}
		case btwire.MsgRequest:
			idx, begin, ln, err := btwire.DecodeRequest(payload)
			if err != nil {
				continue
			}
			pl := pieceLen(t, idx)
			if int64(begin)+int64(ln) > pl || idx < 0 || idx >= t.NumPieces {
				continue
			}
			if ln > len(buf) {
				buf = make([]byte, ln)
			}
			off := int64(idx)*t.PieceLen + int64(begin)
			if _, err := content.ReadAt(buf[:ln], off); err != nil {
				return err
			}
			piece := btwire.EncodePiece(idx, begin, buf[:ln])
			if err := btwire.WriteMessage(conn, btwire.MsgPiece, piece); err != nil {
				return err
			}
		default:
			// ignore
		}
	}
}

func pieceLen(t *metainfo.Torrent, index int) int64 {
	if index == t.NumPieces-1 {
		return t.TotalLen - int64(index)*t.PieceLen
	}
	return t.PieceLen
}
