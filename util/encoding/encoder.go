package encoding

import "io"

type Encoder interface {
	Encode(data interface{}) (err error)
}

type EncoderFactory interface {
	MakeEncoder(w io.Writer) Encoder
}

type Decoder interface {
	Decode(data interface{}) (err error)
}

type DecoderFactory interface {
	MakeDecoder(r io.Reader) Decoder
}

type EncoderFactoryFunc func(w io.Writer) Encoder

func (f EncoderFactoryFunc) MakeEncoder(w io.Writer) Encoder {
	return f(w)
}

type DecoderFactoryFunc func(r io.Reader) Decoder

func (f DecoderFactoryFunc) MakeDecoder(r io.Reader) Decoder {
	return f(r)
}
