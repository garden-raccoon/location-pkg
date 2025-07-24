package models

import (
	proto "github.com/garden-raccoon/location-pkg/protocols/location-pkg"
	"github.com/gofrs/uuid"
)

type Location struct {
	LocationUuid uuid.UUID `json:"location_uuid"`
	Name         string    `json:"name"`
	Country      string    `json:"country"`
	Lang         []string  `json:"lang"`
	Currency     string    `json:"currency"`
	City         string    `json:"city"`
}

// Proto is
func Proto(loc *Location) (*proto.Location, error) {
	d := &proto.Location{
		LocationUuid: loc.LocationUuid.Bytes(),
		Name:         loc.Name,
		Country:      loc.Country,
		Lang:         loc.Lang,
		Currency:     loc.Currency,
		City:         loc.City,
	}
	return d, nil
}

func LocationFromProto(pb *proto.Location) *Location {
	return &Location{
		LocationUuid: uuid.FromBytesOrNil(pb.LocationUuid),
		Name:         pb.Name,
		Country:      pb.Country,
		Lang:         pb.Lang,
		Currency:     pb.Currency,
		City:         pb.City,
	}
}

/*
func MealsFromProto(pb *proto.Meals) []*Meal {
	var meals []*Meal
	for _, b := range pb.Meals {
		meal := MealFromProto(b)
		meals = append(meals, meal)
	}
	return meals
}

*/

func LocationsFromProto(pb *proto.AllLocations) []*Location {
	var locations []*Location
	for _, loc := range pb.Locations {
		locations = append(locations, LocationFromProto(loc))
	}
	return locations
}
