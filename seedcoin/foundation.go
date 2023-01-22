package seedcoin

import (
	"fmt"

	"github.com/0xPolygon/polygon-edge/types"
)

const seedcoinID = 4815162342

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
		Id:      seedcoinID,
		Address: "0x020bB66EC46751CB672a2Ea00b244DEe75970e76",
		Name:    "Seedcoin Foundation",
	},
	{
		Id:      200,
		Address: "0xe0cdBF5cf15D853f2Ef51331f8FD47b95F867085",
		Name:    "Woodland",
	},
	{
		Id:      199,
		Address: "0x340F36Ca8c06AaB6c1EB23684392971b56F8A80D",
		Name:    "Sealand",
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
	return DefaultFoundations.SearchFoundationByID(seedcoinID)
}
