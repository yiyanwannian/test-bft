package main

import (
	"flag"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"math/big"
	"os"
)

func main() {
	var (
		step    int
		dbPath  string
		accAddr string
		contrnm string
	)
	flag.IntVar(&step, "step", 0, "0: modify amount; 1: recovery amount")
	flag.StringVar(&dbPath, "db_path", "", "db of the path")
	flag.StringVar(&accAddr, "user_addr", "", "acc of the addr")
	flag.StringVar(&contrnm, "contract_name", "", "name of the contract")
	flag.Parse()
	if len(dbPath) == 0 || len(accAddr) == 0 {
		panic(fmt.Sprintf("db_path[%s] or user_addr[%s] is null", dbPath, accAddr))
	}
	if info, err := os.Stat(dbPath); err != nil {
		panic(err)
	} else if !info.IsDir() {
		panic("should pass dir path")
	}

	var oncePrint bool
	// 1. open db and init key
	for i := 0; i <= 9; i++ {
		db, err := leveldb.OpenFile(fmt.Sprintf("%s/%d", dbPath, i), nil)
		if err != nil {
			fmt.Println("open db err: ", err)
			continue
		}
		//balanceKey := fmt.Sprintf("B/%s", accAddr)
		balanceKey := fmt.Sprintf("%s#balance"", accAddr) //address+"#balance"
		underlayDBBalanceKey := append(append([]byte(contrnm), '#'), balanceKey...)

		// 2. get the balance from db
		val, err := db.Get(underlayDBBalanceKey, nil)
		if err != nil && err != leveldb.ErrNotFound {
			fmt.Printf("query acc[%s] balance from db failed, reason: %s", accAddr, err)
			continue
		}
		amount := "0"
		if len(val) != 0 {
			amount = string(val)
		}
		balance, ok := big.NewInt(0).SetString(amount, 10)
		if !ok {
			panic(fmt.Sprintf("covert balance bytes[%s] to big.Int failed", balance))
		}
		if !oncePrint {
			fmt.Printf("before update balance: %s\n", balance.String())
		}

		// 3. update the balance in the db
		if step == 0 {
			if err = db.Delete(underlayDBBalanceKey, nil); err != nil {
				continue
			}
		} else if step == 1 {
			if err = db.Put(underlayDBBalanceKey, []byte("100"), nil); err != nil {
				continue
			}
		}

		// 4. check update balance
		if step == 0 {
			if val, err = db.Get(underlayDBBalanceKey, nil); err != leveldb.ErrNotFound {
				panic(fmt.Sprintf("get balance should get not found error, but get another error: %s", err))
			}
			val = []byte("0")
		} else if step == 1 {
			if val, err = db.Get(underlayDBBalanceKey, nil); err != nil {
				panic(fmt.Sprintf("get balance failed, error: %s", err))
			}
		}

		if balance, ok = big.NewInt(0).SetString(string(val), 10); !ok {
			panic(fmt.Sprintf("covert balance bytes[%s] to big.Int failed", balance))
		}
		if !oncePrint {
			fmt.Printf("after update balance: %s\n", balance.String())
		}
		oncePrint = true

		if step == 0 {
			if balance.Int64() != 0 {
				panic(fmt.Sprintf("amount mismatch, expect: %d, actual: %d", 0, balance.Int64()))
			}
		} else if step == 1 {
			if balance.Int64() != 100 {
				panic(fmt.Sprintf("amount mismatch, expect: %d, actual: %d", 100, balance.Int64()))
			}
		}
	}
}
