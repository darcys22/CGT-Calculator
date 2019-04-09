package config

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"time"

	"github.com/naoina/toml"
	log "github.com/Sirupsen/logrus"
)


type Config struct {
	Name					string								`toml:"name"`
	GainsDatabase	string								`toml:"gains_database"`
	PriceDatabase	string								`toml:"price_database"`
	StartFY				time.Time							`toml:"start_f_y"`
	EndFY					time.Time							`toml:"end_f_y"`
	Files					map[string]string			`toml:"files"`
}

func Parse(data []byte) (*Config, error) {
	c := DefaultConfig
	err := toml.Unmarshal(data, &c)
	return &c, err
}

func (c *Config) Encode() ([]byte, error) {
	var buf bytes.Buffer
	e := toml.NewEncoder(&buf)
	err := e.Encode(c)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func LoadConfig(file string) (*Config, error) {
	if file == "" {
		file = defaultConfiglocation()
	}

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	c, err := Parse(dat)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (cfg *Config) Save(file string) (error) {
	if file == "" {
		file = defaultConfiglocation()
	}

	f, err := os.Create(file)
	if err != nil {
		log.Fatal("Could not create config file ", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()
	err = toml.NewEncoder(writer).Encode(cfg)
	if err != nil {
		log.Fatal("Could not write config file ", err)
	}

	return err
}

func NewConfig() (*Config) {
	return &DefaultConfig
}
