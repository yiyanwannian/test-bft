package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterh/liner"
	tikvcfg "github.com/yiyanwannian/client-go/config"
	"github.com/yiyanwannian/client-go/rawkv"
)

// ./test-bft --endpoint="127.0.0.1:2379"

var (
	history_fn = filepath.Join(os.TempDir(), ".liner_history")
)

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

	line := liner.NewLiner()
	line.SetCtrlCAborts(true)

	if f, err := os.Open(history_fn); err == nil {
		_, _ = line.ReadHistory(f)
		_ = f.Close()
	}

	db, err := rawkv.NewClient(ctx, addrs, tikvcfg.Default())
	if err != nil {
		panic(fmt.Sprintf("Error opening %s by tikvdbprovider: %v", endpoint, err))
	}

	defer exitFunc(line, db)

	var keyStr string
	for {
		keyStr = ""
		fmt.Println("------- please input ------")
		if keyStr, err = line.Prompt("key> "); err == nil {
			line.AppendHistory(keyStr)
		} else if err == liner.ErrPromptAborted {
			fmt.Println("exiting")
			exitFunc(line, db)
			break
		} else {
			fmt.Println(fmt.Sprintf("Error reading line: ", err))
		}

		val, err1 := db.Get(ctx, []byte(keyStr))
		if err1 != nil && !strings.Contains(err1.Error(), "not found") {
			fmt.Println(fmt.Sprintf("can not get key: %s, err: %v", keyStr, err))
			continue
		}
		fmt.Println(fmt.Sprintf("raw value: %s", val))
		fmt.Println(fmt.Sprintf("hex encoded value: %s", hex.EncodeToString(val)))
		fmt.Println("-------------------------")
		fmt.Println()
	}
}

func exitFunc(line *liner.State, db *rawkv.Client) {
	if f, err := os.Create(history_fn); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}

	db.Close()
	line.Close()
}