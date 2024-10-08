package main

import (
	"flag"
	"log"
	"os"
	"time"

	"code.vegaprotocol.io/oracles-relay/openoracle"
	"github.com/pelletier/go-toml"
	"nebula.exchange/oracle-price-pusher/prices"
	"nebula.exchange/oracle-price-pusher/twelvedata"
)

type Config struct {
	NodeAddr           string        `toml:"node_addr"`
	WalletMnemonic     string        `toml:"wallet_mnemonic"`
	EthereumPrivateKey string        `toml:"ethereum_private_key"`
	Symbol             string        `toml:"symbol"`
	UpdateFrequency    time.Duration `toml:"update_frequency"`
	TwelveDataApiKey   string        `toml:"twelvedata_apikey"`
}

var flags = struct {
	Config string
}{}

func init() {
	flag.StringVar(&flags.Config, "config", "config.toml", "The configuration of the oracle price pusher")
}

func main() {
	flag.Parse()

	// load our configuration
	config, err := loadConfig(flags.Config)
	if err != nil {
		log.Fatalf("unable to read configuration: %v", err)
	}

	if len(config.Symbol) <= 0 {
		log.Fatalf("missing symbol in config")
	}

	if config.UpdateFrequency == 0 || config.UpdateFrequency <= 5*time.Second {
		log.Fatalf("update frequency is too low")
	}

	p, err := prices.New(
		config.WalletMnemonic,
		config.EthereumPrivateKey,
		config.NodeAddr,
	)
	if err != nil {
		log.Fatalf("could not instanties prices: %v", err)
	}

	ticker := time.NewTicker(config.UpdateFrequency)
	_ = p

	for range ticker.C {
		log.Printf("updating prices")
		price, err := twelvedata.Pull(config.Symbol, config.TwelveDataApiKey)
		if err != nil {
			log.Printf("couldn't get prices: %v", err)
			continue
		}

		log.Printf("new price: %v", price)

		ts := uint64(time.Now().Unix())
		err = p.Send(
			ts,
			[]openoracle.OraclePrice{
				{
					Asset:     config.Symbol,
					Price:     price,
					Timestamp: ts,
				},
			},
		)
		if err != nil {
			log.Printf("could not send prices to the network: %v", err)
		}
	}
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
