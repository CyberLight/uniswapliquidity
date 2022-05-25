## Install abigen
* `go mod init uniswapliquidity`
* `go get github.com/ethereum/go-ethereum`
* `go mod download --json`
    * find `go-etherium` and grab Dir information like this `/go/pkg/mod/github.com/ethereum/go-ethereum@v1.10.17/cmd/abigen`
* `cd /go/pkg/mod/github.com/ethereum/go-ethereum@v1.10.17/cmd/abigen`
* `go build`
    * ignore error at the end of build process about permission denied
* `go install`
* check availability from command line `abigen`

## Generate contract from abi
* `mkdir nftposmgr`
* `mkdir v3factory`
* `mkdir v3pool`
* `abigen --abi=./abis/NFTPositionManager.json --pkg=nftposmgr --out=nftposmgr/main.go`
* `abigen --abi=./abis/UniswapV3Factory.json --pkg=v3factory --out=v3factory/main.go`
* `abigen --abi=./abis/UniswapV3Pool.json --pkg=v3pool --out=v3pool/main.go`

## Run program
* `go build`
* `./uniswapliquidity -h`
```bash
Usage of ./uniswapliquidity:
  -url string
        url to node network (for example, infura url)
  -wallet string
        wallet hex address
```
* `./uniswapliquidity --url=https://kovan.infura.io/v3/xxxxxxxxxxxxxxxxxxxxx --wallet=0xXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX`
```
NFTS:  7
Token ids:  [XXXXX XXXXX XXXXX XXXXX XXXXX XXXXX XXXXX]
{"tokenId":XXXXX,"amount0":0,"amount1":0}
{"tokenId":XXXXX,"amount0":1335185372,"amount1":124340604754067}
{"tokenId":XXXXX,"amount0":0,"amount1":0}
{"tokenId":XXXXX,"amount0":9999999999999999996,"amount1":0}
{"tokenId":XXXXX,"amount0":0,"amount1":99551921146955144}
{"tokenId":XXXXX,"amount0":169614025463176389361,"amount1":101162888585634929975}
{"tokenId":XXXXX,"amount0":99999999999999999,"amount1":970078212752337457449}
```
