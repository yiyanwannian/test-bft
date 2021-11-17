package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	tikvcfg "github.com/yiyanwannian/client-go/config"
	"github.com/yiyanwannian/client-go/rawkv"
)

// ./test-bft --endpoint="127.0.0.1:2379"

func main() {
	var (
		endpoint  string
	)

	flag.StringVar(&endpoint, "endpoint", "", "endpoint of the tikv pds")
	flag.Parse()

	ctx := context.Background()
	addrs := make([]string, 0, 5)
	endpoints := strings.Split(endpoint, ",")
	for _, ep := range endpoints {
		addrs = append(addrs, ep)
	}
	if len(addrs) == 0 {
		panic(fmt.Sprintf("endpoint is invalidate: %s", endpoint))
	}
	db, err := rawkv.NewClient(ctx, addrs, tikvcfg.Default())
	if err != nil {
		panic(fmt.Sprintf("Error opening %s by tikvdbprovider: %v", endpoint, err))
	}
	defer db.Close()

	for {
		fmt.Println("------- please input ------")
		fmt.Print("key: ")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(fmt.Sprintf("read input err: %v", err))
			continue
		}
		if len(text) < 2 {
			fmt.Println("empty key")
			continue
		}

		text = text[:len(text) -1]
		val, err1 := db.Get(ctx, []byte(text))
		if err1 != nil && !strings.Contains(err1.Error(), "not found") {
			fmt.Println( fmt.Sprintf("can not get key: %s, err: %v", text, err))
			continue
		}
		fmt.Println(fmt.Sprintf("raw value: %s", string(val)))
		fmt.Println(fmt.Sprintf("stringed value: %s", string(val)))
		fmt.Println("-------------------------")
		fmt.Println()
	}
}
