package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	envStorageType = "PROMHSD_STORAGE"
	envStorageArgs = "PROMHSD_%s_ARGS"
)

func getStorage() string {
	storageVal := os.Getenv(envStorageType)
	return storageVal
}

func getStorageArgs(storage string) string {
	storageArgs := os.Getenv(fmt.Sprintf(envStorageArgs, strings.ToUpper(storage)))
	return storageArgs
}
