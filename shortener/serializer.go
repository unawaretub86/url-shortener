package shortener

//Decode receives slice of byte and return a tuple of redirect and error
//Encode receives Redirect  and then return a tuple of slice of byte and error
//This is implemented to work with JSON and MessagePack
type RedirectSerializer interface {
	Decode([]byte) (*Redirect, error)
	Encode(*Redirect) ([]byte, error)
}
