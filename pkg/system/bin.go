package system

import "os"

const BinEnvVar = "ORION101_BIN"

func SetBinToSelf() {
	if err := os.Setenv(BinEnvVar, Bin()); err != nil {
		panic(err)
	}
}

func Bin() string {
	bin := os.Getenv(BinEnvVar)
	if bin != "" {
		return bin
	}
	return currentBin()
}
