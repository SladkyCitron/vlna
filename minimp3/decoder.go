package minimp3

import (
	"encoding/binary"
	"io"

	"github.com/SladkyCitron/resona/afmt"
	"github.com/SladkyCitron/resona/aio"
	"github.com/SladkyCitron/resona/encoding/pcm"
	"github.com/tosone/minimp3"
)

type Decoder struct {
	pcmdec aio.SampleReader
}

func NewDecoder(r io.Reader) (*Decoder, error) {
	mp3dec, err := minimp3.NewDecoder(r)
	if err != nil {
		return nil, err
	}
	pcmdec := pcm.NewDecoder(mp3dec, afmt.SampleFormat{
		BitDepth: 16,
		Encoding: afmt.SampleEncodingInt,
		Endian:   binary.LittleEndian,
	})

	return &Decoder{pcmdec: pcmdec}, nil
}

func (d *Decoder) ReadSamples(samples []float32) (int, error) {
	return d.pcmdec.ReadSamples(samples)
}
