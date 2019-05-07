package varmint

type Trackable interface {
	Track() <-chan Status
}
