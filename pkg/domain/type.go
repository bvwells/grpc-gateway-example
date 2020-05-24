package domain

// Type describes the beer type.
type Type int

const (
	Unknown Type = iota + 1
	Ale
	Bitter
	Lager
	IndiaPaleAle
	Stout
	Pilsner
	Porter
	PaleAle
)

var types = [...]string{
	"Unknown",
	"Ale",
	"Bitter",
	"Lager",
	"IndiaPaleAle",
	"Stout",
	"Pilsner",
	"Porter",
	"PaleAle",
}

// String returns the string representation of a beer type.
func (t Type) String() string {
	if Unknown <= t && t <= PaleAle {
		return types[t-1]
	}
	return "Unknown"
}
