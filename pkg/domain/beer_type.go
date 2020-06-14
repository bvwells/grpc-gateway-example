package domain

// BeerType describes the beer type.
type BeerType int

const (
	// Unspecified beer type
	Unspecified BeerType = iota + 1 // Unspecified beer type
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
	"Unspecified",
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
func (t BeerType) String() string {
	if Unspecified <= t && t <= PaleAle {
		return types[t-1]
	}
	return "Unspecified"
}
