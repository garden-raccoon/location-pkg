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
func (loc Location) Proto() *proto.Location {
	return &proto.Location{
		LocationUuid: loc.LocationUuid.Bytes(),
		Name:         loc.Name,
		Country:      loc.Country,
		Lang:         loc.Lang,
		Currency:     loc.Currency,
		City:         loc.City,
	}

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

func LocationsFromProto(pb *proto.AllLocations) []*Location {
	var locations []*Location
	for _, loc := range pb.Locations {
		locations = append(locations, LocationFromProto(loc))
	}
	return locations
}
