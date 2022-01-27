package BLC

import "fmt"

func (cli *CLI) addressLists() {
	fmt.Println("打印所有的钱包地址。")
	wallets := NewWallets()
	for address, _ := range wallets.WalletsMap {
		fmt.Println("address: ", address)
	}
}
