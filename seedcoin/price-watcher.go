package seedcoin

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	PriceFilename               string = "/tmp/seedcoin_price.dat"
	priceBlockWritingMultiplier uint64 = 1e+9
)

var (
	oracleEndpoint    string
	observingInterval time.Duration
	ticker            *time.Ticker
)

// MARK: - Public functions&methods

func Prepare() {
	err := godotenv.Load(".env")
	check(err)

	oracleEndpoint = os.Getenv("ORACLE_URL")
	if len(oracleEndpoint) == 0 {
		check(errors.New("ORACLE_URL is missing"))
	}

	SharedLogger().Log(
		"We will use this oracle %s",
		oracleEndpoint,
	)

	intervalFromEnv := os.Getenv("PRICE_OBSERVING_INTERVAL")
	if len(intervalFromEnv) == 0 {
		check(errors.New("PRICE_OBSERVING_INTERVAL is missing"))
	}

	interval, err := strconv.Atoi(intervalFromEnv)
	check(err)

	observingInterval = time.Second * time.Duration(interval)

	SharedLogger().Log(
		"Price observing duration - %02.f seconds",
		observingInterval.Seconds(),
	)

	bytesPrice, err := download()
	check(err)

	SharedLogger().Log(
		"First time price downloading completed%s",
		"OK",
	)

	priceModel, err := parsePrice(bytesPrice)
	check(err)

	SharedLogger().Log(
		"Received price model: %s",
		priceModel.Description(),
	)

	err = writePriceToFile(bytesPrice)
	check(err)

	SharedLogger().Log(
		"First time price was written to file%s",
		"OK",
	)
}

func WatchForPrice() {
	ticker = time.NewTicker(observingInterval)
	for range ticker.C {
		UpdatePrice()
	}
}

func UpdatePrice() {
	bytesPrice, err := download()
	if err != nil {
		SharedLogger().Log(
			"Price updating error: %s",
			err.Error(),
		)
		return
	}
	priceModel, err := parsePrice(bytesPrice)
	if err != nil {
		SharedLogger().Log(
			"Price model validation error: %s",
			err.Error(),
		)
		return
	}

	SharedLogger().Log(
		"Received price model: %s",
		priceModel.Description(),
	)

	err = writePriceToFile(bytesPrice)
	if err != nil {
		SharedLogger().Log(
			"Price was written to file%s",
			"OK",
		)
	}
	SharedLogger().Log(
		"Price was written to file%s",
		"OK",
	)
}

func LastPrice() (float64, error) {
	lastPriceFromFile, err := readPriceFromFile()
	if err != nil {
		return 0, err
	}

	lastPrice := math.Round(lastPriceFromFile*100) / 100

	return lastPrice, nil
}

func PreparePriceForWritingToBlock(price float64) uint64 {
	bigPrice := new(big.Float).SetPrec(prec).SetFloat64(price)
	bigNormalizer := new(big.Float).SetPrec(prec).SetUint64(priceBlockWritingMultiplier)
	bigResult := new(big.Float).SetPrec(prec).Mul(bigPrice, bigNormalizer)
	result, _ := bigResult.Uint64()
	return result
}

func ExtractPriceFromBlockValue(value uint64) float64 {
	bigValue := new(big.Float).SetPrec(prec).SetUint64(value)
	bigNormalizer := new(big.Float).SetPrec(prec).SetUint64(priceBlockWritingMultiplier)
	bigResult := new(big.Float).SetPrec(prec).Quo(bigValue, bigNormalizer)
	result, _ := bigResult.Float64()
	return result
}

// MARK: - Private functions&methods

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func download() ([]byte, error) {
	resp, err := http.Get(oracleEndpoint)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func parsePrice(bytesArray []byte) (FeePayload, error) {
	var data FeePayload
	if err := json.Unmarshal(bytesArray, &data); err != nil {
		return FeePayload{}, err
	}
	return data, nil
}

func writePriceToFile(bytesPrice []byte) error {
	f, err := os.Create(PriceFilename)
	if err != nil {
		return err
	}

	defer f.Close()

	bytesCount, err := f.Write(bytesPrice)
	if err != nil {
		return err
	}
	if bytesCount == 0 {
		return errors.New("\t0 bytes was written, something went wrong")
	}

	return nil
}

func readPriceFromFile() (float64, error) {
	_, err := os.Stat(PriceFilename)
	if errors.Is(err, os.ErrNotExist) {
		return 0, err
	}

	priceBytes, err := os.ReadFile(PriceFilename)
	if err != nil {
		return 0, err
	}

	feePriceModel, err := parsePrice(priceBytes)
	if err != nil {
		return 0, err
	}

	return feePriceModel.Price, nil
}
