package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"
	"github.com/bvwells/grpc-gateway-example/pkg/infrastructure"
	"github.com/bvwells/grpc-gateway-example/pkg/usecases"

	"github.com/google/uuid"
)

// Beer is a beer obviously!
// The fields field contains the following parameters
// - name
// - country
// - cat_name
// - style_name
// - name_breweries
type Beer struct {
	DatasetID string                 `json:"datasetid"`
	RecordID  string                 `json:"recordid"`
	Fields    map[string]interface{} `json:"fields"`
}

func main() {
	f, err := os.Open("./open-beer-database.json")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	var beers []*Beer
	err = json.Unmarshal(b, &beers)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	settings := &infrastructure.PostgresSettings{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "ilovebeer",
		DBName:   "beers",
	}

	generateID := func() string {
		return uuid.New().String()
	}
	repo, err := infrastructure.NewPostgresBeerRepository(settings, generateID)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	interactor := usecases.NewBeerInteractor(repo)

	for _, beer := range beers {
		name := getField(beer, "name")
		if name == "" {
			continue
		}

		_, err := interactor.CreateBeer(context.Background(), &domain.CreateBeerParams{
			Name:    name,
			Type:    getType(getField(beer, "style_name")),
			Brewer:  getField(beer, "name_breweries"),
			Country: getField(beer, "country"),
		})
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

func getField(beer *Beer, fieldName string) string {
	field, ok := beer.Fields[fieldName]
	if !ok {
		return ""
	}
	val, ok := field.(string)
	if !ok {
		return ""
	}
	return val
}

func getType(in string) domain.Type {
	switch in {
	case "American-Style Brown Ale":
		return domain.Ale
	case "English-Style Pale Mild Ale":
		return domain.PaleAle
	case "Belgian-Style Quadrupel":
		return domain.Lager
	case "Foreign (Export)-Style Stout":
		return domain.Stout
	case "Other Belgian-Style Ales":
		return domain.Ale
	case "Belgian-Style Tripel":
		return domain.Ale
	case "Winter Warmer":
		return domain.Ale
	case "American-Style India Pale Ale":
		return domain.IndiaPaleAle
	case "Scottish-Style Light Ale":
		return domain.Ale
	case "Belgian-Style Pale Ale":
		return domain.PaleAle
	case "English-Style Dark Mild Ale":
		return domain.Ale
	case "European Low-Alcohol Lager":
		return domain.Lager
	case "Imperial or Double Red Ale":
		return domain.Ale
	case "American-Style Lager":
		return domain.Lager
	case "Out of Category":
		return domain.Unknown
	case "American-Style Imperial Stout":
		return domain.Stout
	case "Herb and Spice Beer":
		return domain.Ale
	case "Specialty Beer":
		return domain.Ale
	case "Kellerbier - Ale":
		return domain.Ale
	case "Fruit Beer":
		return domain.Ale
	case "South German-Style Hefeweizen":
		return domain.Ale
	case "Extra Special Bitter":
		return domain.Ale
	case "Irish-Style Red Ale":
		return domain.Ale
	case "Ordinary Bitter":
		return domain.Bitter
	case "Vienna-Style Lager":
		return domain.Lager
	case "Smoke Beer":
		return domain.Ale
	case "Traditional German-Style Bock":
		return domain.Ale
	case "Sweet Stout":
		return domain.Stout
	case "French & Belgian-Style Saison":
		return domain.Ale
	case "Imperial or Double India Pale Ale":
		return domain.IndiaPaleAle
	case "German-Style Schwarzbier":
		return domain.Ale
	case "American-Style Pale Ale":
		return domain.PaleAle
	case "Old Ale":
		return domain.Ale
	case "American-Style Barley Wine Ale":
		return domain.Ale
	case "American-Style Dark Lager":
		return domain.Lager
	case "Belgian-Style Dubbel":
		return domain.Ale
	case "German-Style Brown Ale/Altbier":
		return domain.Ale
	case "Oatmeal Stout":
		return domain.Stout
	case "Belgian-Style Dark Strong Ale":
		return domain.Ale
	case "American-Style Stout":
		return domain.Stout
	case "Special Bitter or Best Bitter":
		return domain.Bitter
	case "American-Style Light Lager":
		return domain.Lager
	case "German-Style Oktoberfest":
		return domain.Ale
	case "Classic English-Style Pale Ale":
		return domain.PaleAle
	case "American Rye Ale or Lager":
		return domain.Lager
	case "Specialty Honey Lager or Ale":
		return domain.Lager
	case "German-Style Doppelbock":
		return domain.Ale
	case "South German-Style Weizenbock":
		return domain.Ale
	case "German-Style Pilsener":
		return domain.Pilsner
	case "American-Style Strong Pale Ale":
		return domain.PaleAle
	case "case Baltic-Style Porter":
		return domain.Porter
	case "Belgian-Style Pale Strong Ale":
		return domain.PaleAle
	case "Belgian-Style Fruit Lambic":
		return domain.Ale
	case "Scotch Ale":
		return domain.Ale
	case "Belgian-Style White":
		return domain.Ale
	case "Dark American-Belgo-Style Ale":
		return domain.Ale
	case "Pumpkin Beer":
		return domain.Ale
	case "Classic Irish-Style Dry Stout":
		return domain.Stout
	case "American-Style India Black Ale":
		return domain.IndiaPaleAle
	case "Porter":
		return domain.Porter
	case "English-Style India Pale Ale":
		return domain.IndiaPaleAle
	case "Bamberg-Style Bock Rauchbier":
		return domain.Ale
	case "American-Style Amber/Red Ale":
		return domain.Ale
	case "German-Style Heller Bock/Maibock":
		return domain.Ale
	case "Golden or Blonde Ale":
		return domain.Ale
	case "American-Style Cream Ale or Lager":
		return domain.Lager
	case "Strong Ale":
		return domain.Ale
	default:
		return domain.Unknown
	}
}
