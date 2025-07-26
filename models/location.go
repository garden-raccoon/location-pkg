package models

import (
	proto "github.com/garden-raccoon/location-pkg/protocols/location-pkg"
	"github.com/gofrs/uuid"
)

type Location struct {
	LocationUuid uuid.UUID   `json:"location_uuid"`
	Name         string      `json:"name"`
	Country      string      `json:"country"`
	Lang         []string    `json:"lang"`
	Currency     string      `json:"currency"`
	City         string      `json:"city"`
	Activeness   *Activeness `json:"activeness"`
}
type MealLocation struct {
	City     string    `json:"location_uuid"`
	MealUUID uuid.UUID `json:"meal_uuid"`
	Price    float64   `json:"price"`
}

type Activeness struct {
	IsActive bool `json:"is_active"`
}

// Proto is
func (loc Location) Proto() *proto.Location {
	l := &proto.Location{
		LocationUuid: loc.LocationUuid.Bytes(),
		Name:         loc.Name,
		Country:      loc.Country,
		Lang:         loc.Lang,
		Currency:     loc.Currency,
		City:         loc.City,
	}
	if loc.Activeness != nil {
		l.Activeness = &proto.Activeness{
			IsActive: loc.Activeness.IsActive,
		}
	}
	return l

}

func LocationFromProto(pb *proto.Location) *Location {
	l := &Location{
		LocationUuid: uuid.FromBytesOrNil(pb.LocationUuid),
		Name:         pb.Name,
		Country:      pb.Country,
		Lang:         pb.Lang,
		Currency:     pb.Currency,
		City:         pb.City,
	}
	if pb.Activeness != nil {
		l.Activeness = &Activeness{
			IsActive: pb.Activeness.IsActive,
		}
	}
	return l
}

func LocationsFromProto(pb *proto.AllLocations) []*Location {
	var locations []*Location
	for _, loc := range pb.Locations {
		locations = append(locations, LocationFromProto(loc))
	}
	return locations
}
