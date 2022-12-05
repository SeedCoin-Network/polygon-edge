package main

import (
	_ "embed"

	"github.com/0xPolygon/polygon-edge/command/root"
	"github.com/0xPolygon/polygon-edge/licenses"
	"github.com/0xPolygon/polygon-edge/seedcoin"
)

var (
	//go:embed LICENSE
	license string
)

func main() {
	licenses.SetLicense(license)
	go seedcoin.SharedCalculator().StartObservingGasCalculationCoef()
	root.NewRootCommand().Execute()
}
