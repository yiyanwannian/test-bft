package main

import (
	"flag"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"math/big"
)

func main() {
	var(
		dbPath string
		accAddr string
	)
	flag.StringVar(&dbPath, "db_path", "", "db of the path")
	flag.StringVar(&accAddr, "user_addr", "", "acc of the addr")
	flag.Parse()
	if len(dbPath) == 0 || len(accAddr) == 0{
		panic(fmt.Sprintf("db_path[%s] or user_addr[%s] is null", dbPath, accAddr))
	}

	// 1. open db and init key
	for i := 0; i <= 9; i++{
		db, err := leveldb.OpenFile(fmt.Sprintf("%s/%d", dbPath, i), nil)
		if err != nil {
			continue
			//panic(err)
		}
		balanceKey := fmt.Sprintf("B/%s", accAddr)
		underlayDBBalanceKey := append(append([]byte("SYSTEM_CONTRACT_DPOS_ERC20"), '#'), balanceKey...)
		//fmt.Printf("key: %s, hash: %x\n", string(underlayDBBalanceKey), sha256.Sum256(underlayDBBalanceKey))

		// 2. get the balance from db
		val, err := db.Get(underlayDBBalanceKey, nil)
		if err != nil {
			continue
			panic(fmt.Sprintf("get balance failed, error: %s", err))
		}
		balance, ok := big.NewInt(0).SetString(string(val), 10)
		if !ok {
			continue
			panic(fmt.Sprintf("covert balance bytes[%s] to big.Int failed", balance))
		}
		fmt.Printf("before update balance: %s\n", balance.String())

		// 3. update the balance in the db
		if err = db.Put(underlayDBBalanceKey, []byte("0"), nil); err != nil {
			continue
			panic(fmt.Sprintf("update balance failed, error: %s", err))
		}

		// 4. check update balance
		if val, err = db.Get(underlayDBBalanceKey, nil); err != nil {
			continue
			panic(fmt.Sprintf("get balance failed, error: %s", err))
		}
		if balance, ok = big.NewInt(0).SetString(string(val), 10); !ok {
			continue
			panic(fmt.Sprintf("covert balance bytes[%s] to big.Int failed", balance))
		}
		fmt.Printf("after update balance: %s\n", balance.String())
		if balance.Int64() != 0 {
			continue
			panic("update balance failed")
		}
	}

}
