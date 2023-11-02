package msgpack

import (
	errs "github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"

	"github.com/unawaretub86/url-shortener/shortener"
)

// This is for implement the redirect serializer interface on this struct
type Redirect struct{}

// This gets a slice of byte an then it generates a redirect
// and then we use json unmarshal to the input and put it inside of the redirects
func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, errs.Wrap(err, "serializer.Redirect.Decoder")
	}

	return redirect, nil
}

// This takes the input which is a redirect struct type and then this will pass back
// a slice of bytes and an error
//
// Json Marshal will give us the raw bytes
func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errs.Wrap(err, "serializer.Redirect.Decoder")
	}

	return rawMsg, nil
}
