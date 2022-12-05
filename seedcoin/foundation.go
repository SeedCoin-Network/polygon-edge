package seedcoin

import (
	"fmt"

	"github.com/0xPolygon/polygon-edge/types"
)

type Foundation struct {
	Id      uint64
	Address string
	Name    string
}

func (f Foundation) Description() string {
	desc := fmt.Sprintf(
		"-----\nDescription of Foundation:\n\tid: %d,\n\taddress: %s,\n\tname: %s",
		f.Id,
		f.Address,
		f.Name,
	)
	return desc
}

type Foundations []Foundation

var DefaultFoundations Foundations = Foundations{
	{
		Id:      2000,
		Address: "0xc425C9564906f4B8f8a849C5C1F2b4272534C6D7",
		Name:    "Seedcoin Foundation",
	},
	{
		Id:      2001,
		Address: "0x6965DF0B63B73E0060ca9ad1EB563561675219BF",
		Name:    "AARP Foundation",
	},
	{
		Id:      2002,
		Address: "0x4dD3CBD243A5982B2dF9889B283b122855194a30",
		Name:    "Adelson Foundation",
	},
}

func (f Foundations) FoundationsIDs() []uint64 {
	keys := make([]uint64, 0, len(f))

	for _, foundation := range f {
		keys = append(keys, foundation.Id)
	}

	return keys
}

func (f Foundations) SearchFoundationByID(id uint64) *Foundation {
	for _, foundation := range f {
		if foundation.Id == id {
			return &foundation
		}
	}
	return nil
}

func (f Foundations) ContainsFoundationWithID(id uint64) bool {
	for _, element := range f {
		if element.Id == id {
			return true
		}
	}

	return false
}

func (f Foundations) ContainsAddress(addr types.Address) bool {
	for _, element := range f {
		if element.AddressObject() == addr {
			return true
		}
	}

	return false
}

func (f Foundation) AddressObject() types.Address {
	return types.StringToAddress(f.Address)
}

func SeedcoinFoundation() *Foundation {
	return DefaultFoundations.SearchFoundationByID(2000)
}
