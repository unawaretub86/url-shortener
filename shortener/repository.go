package shortener

//This connect our business logic with our repository
type RedirectRepository interface {
	Find(string) (*Redirect, error)
	Store(*Redirect) error
}
