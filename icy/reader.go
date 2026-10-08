package icy

import (
	"io"
	"strings"
)

type Reader struct {
	r          io.Reader
	metaint    int
	bytesUntil int
	onMeta     func(map[string]string)
}

func NewReader(r io.Reader, metaint int, onMeta func(map[string]string)) *Reader {
	if onMeta == nil {
		return &Reader{
			r:          r,
			metaint:    metaint,
			bytesUntil: metaint,
			onMeta:     func(map[string]string) {},
		}
	}
	return &Reader{
		r:          r,
		metaint:    metaint,
		bytesUntil: metaint,
		onMeta:     onMeta,
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.metaint <= 0 {
		return r.r.Read(p)
	}
	if r.bytesUntil == 0 {
		if err := r.readMetadata(); err != nil {
			return 0, err
		}
		r.bytesUntil = r.metaint
	}
	toRead := min(len(p), r.bytesUntil)
	n, err = r.r.Read(p[:toRead])
	r.bytesUntil -= n
	return n, err
}

func (r *Reader) readMetadata() error {
	lenByte := make([]byte, 1)
	_, err := io.ReadFull(r.r, lenByte)
	if err != nil {
		return err
	}
	metadataLength := int(lenByte[0]) * 16
	if metadataLength == 0 {
		return nil
	}
	metadata := make([]byte, metadataLength)
	_, err = io.ReadFull(r.r, metadata)
	if err != nil {
		return err
	}

	metaStr := strings.TrimRight(string(metadata), "\x00")
	if metaStr != "" {
		metaMap := make(map[string]string)
		for p := range strings.SplitSeq(metaStr, ";") {
			kv := strings.SplitN(p, "=", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				val := strings.Trim(strings.TrimSpace(kv[1]), "'")
				metaMap[key] = val
			}
		}
		r.onMeta(metaMap)
	}

	return nil
}
