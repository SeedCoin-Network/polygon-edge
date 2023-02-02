package seedcoin

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"math"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/0xPolygon/polygon-edge/types"
)

var gasCalculatorOnceSyncPoint sync.Once

type GasCalculator struct {
	GasCalculationCoef float64
	ticker             *time.Ticker
	priceFromBlock     float64
	mode               WorkingMode
}

type WorkingMode uint64

var (
	InlineMode  WorkingMode = 0
	SyncingMode WorkingMode = 1
)

const GasPriceGwei = 200

var singletonCalculator *GasCalculator

const (
	oracleEndpoint    = "https://api.seedcoin.network/oracle/price"
	observingInterval = time.Second * 120
)

func SharedCalculator() *GasCalculator {
	if singletonCalculator == nil {
		gasCalculatorOnceSyncPoint.Do(
			func() {
				singletonCalculator = &GasCalculator{
					mode: InlineMode,
				}
			})
	}

	return singletonCalculator
}

func (g *GasCalculator) SetMode(mode WorkingMode) {
	g.mode = mode
	switch mode {
	case InlineMode:
		SharedLogger().Log("Mode changed to %s", "Inline")
	case SyncingMode:
		SharedLogger().Log("Mode changed to %s", "Syncing")
	}
}

func (g *GasCalculator) GetMode() WorkingMode {
	return g.mode
}

func (g *GasCalculator) ApplyPriceFromBlockHeader(header *types.Header) {
	if header == nil {
		SharedLogger().Log("[SYNCING] Trying to apply price from block failed, block header is nil %s", "")
		return
	}

	if header.CoinPrice == nil {
		value := 1.0
		SharedLogger().Log("[SYNCING] coin_price in header missing, using value %f", value)
		g.priceFromBlock = value
	} else {
		bitsPrice := binary.BigEndian.Uint64(header.CoinPrice)
		priceFromBlock := math.Float64frombits(bitsPrice)
		SharedLogger().Log("[SYNCING] coin_price in header found, using value %f", priceFromBlock)
		g.priceFromBlock = priceFromBlock
	}
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

func (g *GasCalculator) GasCost(amount *big.Int, isExecutionCalculation bool) uint64 {
	const prec = 512

	var x float64
	if isExecutionCalculation {
		switch g.mode {
		case SyncingMode:
			x = g.priceFromBlock
		case InlineMode:
			x = g.GasCalculationCoef
		default:
			x = g.GasCalculationCoef
		}
	} else {
		x = g.GasCalculationCoef
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
