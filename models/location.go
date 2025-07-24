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
}

// Proto is
func Proto(loc *Location) (*proto.Location, error) {
	d := &proto.Location{
		LocationUuid: loc.LocationUuid.Bytes(),
		Name:         loc.Name,
		Country:      loc.Country,
		Lang:         loc.Lang,
		Currency:     loc.Currency,
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
	}
}
