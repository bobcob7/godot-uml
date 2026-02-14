// Package encoding implements PlantUML URL encoding/decoding (DEFLATE + custom base64).
//
// PlantUML uses a custom base64 alphabet for URL-safe encoding:
//
//	0-9 → 0-9, A-Z → 10-35, a-z → 36-61, - → 62, _ → 63
package encoding

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"strings"
)

// PlantUML's custom base64 alphabet.
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"

// Encode compresses PlantUML text using DEFLATE and encodes with custom base64.
func Encode(text string) (string, error) {
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		return "", fmt.Errorf("creating deflate writer: %w", err)
	}
	if _, err := w.Write([]byte(text)); err != nil {
		return "", fmt.Errorf("writing deflate data: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("closing deflate writer: %w", err)
	}
	return encode64(buf.Bytes()), nil
}

// Decode decodes a PlantUML-encoded string back to plain text.
func Decode(encoded string) (string, error) {
	data, err := decode64(encoded)
	if err != nil {
		return "", fmt.Errorf("decoding base64: %w", err)
	}
	r := flate.NewReader(bytes.NewReader(data))
	defer func() { _ = r.Close() }()
	result, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("decompressing: %w", err)
	}
	return string(result), nil
}

func encode64(data []byte) string {
	var sb strings.Builder
	for i := 0; i < len(data); i += 3 {
		switch {
		case i+2 < len(data):
			append4(&sb, data[i], data[i+1], data[i+2])
		case i+1 < len(data):
			append4(&sb, data[i], data[i+1], 0)
		default:
			append4(&sb, data[i], 0, 0)
		}
	}
	return sb.String()
}

func append4(sb *strings.Builder, b1, b2, b3 byte) {
	c1 := b1 >> 2
	c2 := ((b1 & 0x3) << 4) | (b2 >> 4)
	c3 := ((b2 & 0xF) << 2) | (b3 >> 6)
	c4 := b3 & 0x3F
	sb.WriteByte(alphabet[c1])
	sb.WriteByte(alphabet[c2])
	sb.WriteByte(alphabet[c3])
	sb.WriteByte(alphabet[c4])
}

func decode64(s string) ([]byte, error) {
	var result []byte
	for i := 0; i < len(s); i += 4 {
		end := min(i+4, len(s))
		chunk := s[i:end]
		vals := make([]byte, 4)
		for j := range 4 {
			if j < len(chunk) {
				idx := strings.IndexByte(alphabet, chunk[j])
				if idx < 0 {
					return nil, fmt.Errorf("invalid character %q at position %d", chunk[j], i+j)
				}
				vals[j] = byte(idx)
			}
		}
		result = append(result, (vals[0]<<2)|(vals[1]>>4))
		result = append(result, (vals[1]<<4)|(vals[2]>>2))
		result = append(result, (vals[2]<<6)|vals[3])
	}
	return result, nil
}
