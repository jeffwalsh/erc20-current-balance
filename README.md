## Example
To run the example, you can use `./bin/erc20-current-balance` and it will create a spreadsheet of current balances for a single user, defined in example.txt, for the FUEL token.

## To configure and run:

Make a file in the same directory, `yourfile.txt.`
Make it line-delineated list of Ethereum addresses.
Get the contract address of the ERC20 token you want to check.
`./bin/erc20-current-balance -contractAddress=ERC20_ADDRESS -file=YOURFILE.TXT`
