package seedcoin

import (
	"math"
	"math/big"
	"sync"

	"github.com/0xPolygon/polygon-edge/types"
)

var gasCalculatorOnceSyncPoint sync.Once

type GasCalculator struct{}

const (
	GasPriceGwei = 200
	prec         = 512
)

var singletonCalculator *GasCalculator

func SharedCalculator() *GasCalculator {
	if singletonCalculator == nil {
		gasCalculatorOnceSyncPoint.Do(
			func() {
				singletonCalculator = &GasCalculator{}
			})
	}

	return singletonCalculator
}

func (g *GasCalculator) GasCost(amount *big.Int, isExecutionCalculation bool, header *types.Header) uint64 {
	lastPrice, err := LastPrice()
	var x float64
	if header != nil {
		priceFromBlock := ExtractPriceFromBlockValue(header.CoinPrice)
		x = priceFromBlock
	} else {
		if err != nil {
			SharedLogger().Log(
				"Couldn't load last price from file%s",
				"FAIL",
			)
			x = 1
		} else {
			x = lastPrice
		}
	}
	// λ=0.01+0.98/(1+(x+1)^{24})
	value := (1.0 + math.Pow((x+1.0), 24))
	lambda := 0.01 + 0.98/value

	bigLambda := new(big.Float).SetPrec(prec).SetFloat64(lambda)
	bigFloatAmount := new(big.Float).SetPrec(prec).SetInt(amount)
	bigNormalizer := new(big.Float).SetPrec(prec).SetFloat64(1e-9)
	bigNormalizedAmount := new(big.Float).SetPrec(prec).Mul(bigFloatAmount, bigNormalizer)
	bigResult := new(big.Float).SetPrec(prec).Mul(bigNormalizedAmount, bigLambda)
	bigGasPrice := new(big.Float).SetPrec(prec).SetUint64(GasPriceGwei)
	bigTotalAmount := new(big.Float).SetPrec(prec).Quo(bigResult, bigGasPrice)

	plainGasPrice, _ := bigTotalAmount.Uint64()

	return plainGasPrice
}

func (g *GasCalculator) BaseComission(amount *big.Int) *big.Int {
	devider := new(big.Int).SetUint64(100)

	baseComission := new(big.Int).Quo(amount, devider)

	return baseComission
}
