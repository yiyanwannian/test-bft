package main

import (
	"flag"
	"fmt"
	"math/big"
	"github.com/syndtr/goleveldb/leveldb"
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
	db, err := leveldb.OpenFile(".", nil)
	if err != nil {
		panic(err)
	}
	balanceKey := fmt.Sprintf("B/%s", "accountAddr")
	underlayDBBalanceKey := append(append([]byte("SYSTEM_CONTRACT_DPOS_ERC20"), '#'), balanceKey...)

	// 2. get the balance from db
	val, err := db.Get(underlayDBBalanceKey, nil)
	if err != nil {
		panic(err)
	}
	balance, ok := big.NewInt(0).SetString(string(val), 10)
	if !ok {
		panic(fmt.Sprintf("covert balance bytes[%s] to big.Int failed", balance))
	}
	fmt.Printf("before update balance: %s", balance.String())

	// 3. update the balance in the db
	if err = db.Put(underlayDBBalanceKey, []byte("1000"), nil); err != nil {
		panic(err)
	}

	// 4. check update balance
	if val, err = db.Get(underlayDBBalanceKey, nil); err != nil {
		panic(err)
	}
	if balance, ok = big.NewInt(0).SetString(string(val), 10); !ok {
		panic(fmt.Sprintf("covert balance bytes[%s] to big.Int failed", balance))
	}
	fmt.Printf("after update balance: %s", balance.String())
	if balance.Int64() != 1000 {
		panic("update balance failed")
	}
}
