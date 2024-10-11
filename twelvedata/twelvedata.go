package twelvedata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"code.vegaprotocol.io/vega/libs/num"
)

var (
	urlQuery = func(symbol, apiKey, micCode string) string {
		return fmt.Sprintf("https://api.twelvedata.com/price?symbol=%s&apikey=%s&mic_code=%v", symbol, apiKey, micCode)
	}
)

type response struct {
	Price   string `json:"price"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Pull(symbol, apiKey, micCode string) (string, error) {
	req, err := http.NewRequest("GET", urlQuery(symbol, apiKey, micCode), nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	price := response{}
	err = json.Unmarshal(body, &price)
	if err != nil {
		return "", err
	}

	if price.Status == "error" {
		return "", fmt.Errorf("api error: %v", price.Message)
	}

	priceD, err := num.DecimalFromString(price.Price)
	if err != nil {
		return "", err
	}
	priceD = priceD.Mul(num.DecimalFromFloat(10).Pow(num.DecimalFromInt64(8)))

	return priceD.Truncate(0).String(), nil
}
