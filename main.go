package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/vmware/vmw-guestinfo/rpcvmx"

	"github.com/vmware/vmw-guestinfo/vmcheck"
	ovf "github.com/vmware/vmw-ovflib"
)

func isBackdoorAvailable() bool {
	res, err := vmcheck.IsVirtualWorld()
	return res && (err == nil)
}

func readConfig(key string) (string, error) {
	data, err := rpcvmx.NewConfig().String(key, "")
	if err == nil {
		log.Printf("Read from %q: %q\n", key, data)
	} else {
		log.Printf("Failed to read from %q: %v\n", key, err)
	}
	return data, err
}

func main() {
	ovfEnvPtr := flag.String("ovfenv", defaultOvfFilePath, "path to ovf environment file")

	flag.Parse()

	var ovfEnv []byte

	if _, err := os.Stat(*ovfEnvPtr); os.IsNotExist(err) {
		data, err := readConfig("ovfenv")
		if err != nil {
			ovfEnv = make([]byte, 0)
		} else {
			ovfEnv = []byte(data)
		}
	} else {
		if !isBackdoorAvailable() {
			os.Exit(1)
		}

		ovfEnv, err = ioutil.ReadFile(*ovfEnvPtr)
		if err != nil {
			ovfEnv = make([]byte, 0)
		}
	}

	env, err := ovf.ReadEnvironment(ovfEnv)
	if err != nil {
		os.Exit(1)
	}

	for _, arg := range flag.Args() {
		val, ok := env.Properties[arg]
		if !ok {
			os.Exit(1)
		}
		fmt.Println(val)
	}
}
