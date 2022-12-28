package seedcoin

import (
	"encoding/json"
	"io"
	"math"
	"math/big"
	"net/http"
	"sync"
	"time"
)

var gasCalculatorOnceSyncPoint sync.Once

type GasCalculator struct {
	GasCalculationCoef float64
	ticker             *time.Ticker
}

const GasPriceGwei = 200

var singletonCalculator *GasCalculator

const (
	oracleEndpoint    = "https://api.seedcoin.network/oracle/price"
	observingInterval = time.Second * 30
)

func SharedCalculator() *GasCalculator {
	if singletonCalculator == nil {
		gasCalculatorOnceSyncPoint.Do(
			func() {
				singletonCalculator = &GasCalculator{}
			})
	}

	return singletonCalculator
}

func (g *GasCalculator) StartObservingGasCalculationCoef() {
	g.ticker = time.NewTicker(observingInterval)
	for range g.ticker.C {
		updatingError := g.UpdateGasCalculationCoef()
		if updatingError != nil {
			SharedLogger().Log("%s", updatingError)
			println(updatingError)
		}
	}
}

func (g *GasCalculator) UpdateGasCalculationCoef() error {
	resp, reqErr := http.Get(oracleEndpoint)
	if reqErr != nil {
		return reqErr
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return readErr
	}

	var data FeePayload
	parseErr := json.Unmarshal(body, &data)
	if parseErr != nil {
		return parseErr
	}

	g.GasCalculationCoef = data.Price
	SharedLogger().Log("Received price: %f", data.Price)

	return nil
}

func (g *GasCalculator) GasCost(amount *big.Int) uint64 {
	const prec = 512

	x := g.GasCalculationCoef
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
