package domain

// Type describes the beer type.
type Type int

const (
	// Unknown beer type
	Unknown Type = iota + 1 // Unknown beer type
	// Ale is a beer type, see  https://en.wikipedia.org/wiki/Ale
	Ale
	// Bitter is a beer type, see https://en.wikipedia.org/wiki/Bitter_(beer)
	Bitter
	// Lager is a beer type, see  https://en.wikipedia.org/wiki/Lager
	Lager
	// IndiaPaleAle (IPA) is a beer type, see https://en.wikipedia.org/wiki/India_pale_ale
	IndiaPaleAle
	// Stout is a beer type, see https://en.wikipedia.org/wiki/Stout
	Stout
	// Pilsner is a beer type, see https://en.wikipedia.org/wiki/Pilsner
	Pilsner
	// Porter is a beer type, see https://en.wikipedia.org/wiki/Porter_(beer)
	Porter
	// PaleAle is a beer type, see https://en.wikipedia.org/wiki/Pale_ale
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
