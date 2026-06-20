package daemon

import (
	"bytes"
	"testing"
)

func TestBrotliCompression(t *testing.T) {
	original := []byte("Aether is a decentralized peer to peer network architecture.")
	compressed := CompressBrotli(original)
	
	if len(compressed) == 0 {
		t.Fatal("Compression resulted in empty byte array")
	}

	decompressed, err := DecompressBrotli(compressed)
	if err != nil {
		t.Fatalf("Decompression failed: %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		t.Fatalf("Decompressed output does not match original: got %s, want %s", string(decompressed), string(original))
	}
}
