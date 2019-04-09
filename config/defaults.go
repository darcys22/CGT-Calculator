package config

import (
	"os"
	"os/user"
	"path/filepath"
	"time"
)

var DefaultConfig = Config{}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}

	DefaultConfig.GainsDatabase = filepath.Join(home,".config","cgtcalc","gainsdb")
	DefaultConfig.PriceDatabase = filepath.Join(home,".config","cgtcalc","pricedb")

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	//currentLocation := now.Location()
	currentLocation := time.UTC

	if currentMonth <= 6 {
		DefaultConfig.EndFY = time.Date(currentYear-1, 6, 30, 23, 59, 59, 999999999, currentLocation)
		DefaultConfig.StartFY = time.Date(currentYear-2, 7, 1, 0, 0, 0, 0, currentLocation)
	} else {
		DefaultConfig.EndFY = time.Date(currentYear, 6, 30, 23, 59, 59, 999999999, currentLocation)
		DefaultConfig.StartFY = time.Date(currentYear-1, 7, 1, 0, 0, 0, 0, currentLocation)
	}
	DefaultConfig.Files = make(map[string]string)


}

func defaultConfiglocation() string {
	wd, _ := os.Getwd()
	file := filepath.Join(wd,"cgt.conf")
	return file
}
