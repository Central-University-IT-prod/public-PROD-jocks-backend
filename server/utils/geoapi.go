package utils

import (
	"fmt"

	"github.com/yuriizinets/go-nominatim"
)

func ErrNotFound(target string) error {
	return fmt.Errorf("%v: not found", target)
}

type Coords struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func GetCoords(location string) (Coords, error) {
	n := nominatim.Nominatim{}
	res, err := n.Search(nominatim.SearchParameters{
		Query: location,
	})

	if err != nil {
		return Coords{}, err
	}

	if len(res) == 0 {
		return Coords{}, ErrNotFound(location)
	}

	return Coords{Lat: res[0].Lat, Lng: res[0].Lng}, nil
}
