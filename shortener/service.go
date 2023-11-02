package shortener

type RedirectService interface {
	Find(string) (*Redirect, error)
	Store(*Redirect) error
}
