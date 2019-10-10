package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tealeg/xlsx"
)

type etherscanFuelBalanceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func main() {
	addrPtr := flag.String("file", "example.txt", `Filename you want to scan through. Must be delineated by lines. EG:
		0x743009bf25ff1bd56305b51f395dd1d0784a427f
		0x743009bf25ff1bd56305b51f395dd1d0784a427f
		0x743009bf25ff1bd56305b51f395dd1d0784a427f
		`)

	tokenPtr := flag.String("contractAddress", "0xea38eaa3c86c8f9b751533ba2e562deb9acded40", `Token address you want to scan through. Must be a valid ERC20 address, eg: 0xea38eaa3c86c8f9b751533ba2e562deb9acded40
			`)

	file, err := os.Open(*addrPtr)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	defer file.Close()

	var addresses []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		addresses = append(addresses, scanner.Text())
	}

	excel := xlsx.NewFile()

	sheet, err := excel.AddSheet("ERC20 Current Balance List")
	row := sheet.AddRow()
	vals := []string{"Address", "Current Balance"}
	for _, val := range vals {
		cell := row.AddCell()
		cell.Value = val
	}
	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	for _, addr := range addresses {
		// avoid etherscan rate limit
		time.Sleep(300 * time.Millisecond)
		handle(sheet, addr, *tokenPtr)
	}

	if err := excel.Save("ERC20 Current Balance List.xlsx"); err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
}

func handle(sheet *xlsx.Sheet, addr, tokenAddr string) {
	row := sheet.AddRow()

	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=tokenbalance&contractaddress=%s&address=%s&tag=latest", tokenAddr, addr)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("couldnt generate request to etherscan")
		os.Exit(1)
	}

	var fb etherscanFuelBalanceResponse
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&fb)

	if fb.Status != "1" || fb.Message != "OK" {
		fmt.Printf("unsuccessful response from etherscan for address %s, returning early", addr)
		return
	}

	vals := []string{addr, fb.Result}

	for _, val := range vals {
		cell := row.AddCell()
		cell.Value = val
	}
}
