package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"sort"
	"time"
	"uniswapliquidity/nftposmgr"
	"uniswapliquidity/v3factory"
	"uniswapliquidity/v3pool"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var base = big.NewInt(2)
var power128 = big.NewInt(128)
var power32 = big.NewInt(32)
var power96 = big.NewInt(96)
var Q128 = new(big.Int).Exp(base, power128, nil)
var Q32 = new(big.Int).Exp(base, power32, nil)
var Q96 = new(big.Int).Exp(base, power96, nil)
var MIN_TICK = big.NewInt(-887272)
var MAX_TICK = big.NewInt(887272)
var zero = big.NewInt(0)
var one = big.NewInt(1)
var maxUInt256, _ = new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

type MaskAndMulByPair struct {
	mask  *big.Int
	mulBy *big.Int
}

type AmountResult struct {
	TokenId *big.Int `json:"tokenId"`
	Amount0 *big.Int `json:"amount0"`
	Amount1 *big.Int `json:"amount1"`
}

func mulShift(val, mulBy *big.Int) *big.Int {
	mul := new(big.Int).Mul(val, mulBy)
	return new(big.Int).Div(mul, Q128)
}

func notZero(val *big.Int, mask *big.Int) bool {
	result := new(big.Int).And(val, mask)
	return result.Cmp(zero) != 0
}

func HexBigInt(value string) (*big.Int, bool) {
	intValue, ok := new(big.Int).SetString(value, 16)
	return intValue, ok
}

func GetSqrtRatioAtTick(tick *big.Int) (*big.Int, error) {
	valid := tick.Cmp(MIN_TICK) >= 0 && tick.Cmp(MAX_TICK) <= 0

	if !valid {
		return nil, errors.New("invariant tick")
	}

	absTick := new(big.Int).Abs(tick)
	ratio0, _ := HexBigInt("fffcb933bd6fad37aa2d162d1a594001")
	ratio1, _ := HexBigInt("100000000000000000000000000000000")
	mask1, _ := HexBigInt("1")

	mask2, _ := HexBigInt("2")
	mask4, _ := HexBigInt("4")
	mask8, _ := HexBigInt("8")
	mask10, _ := HexBigInt("10")
	mask20, _ := HexBigInt("20")
	mask40, _ := HexBigInt("40")
	mask80, _ := HexBigInt("80")
	mask100, _ := HexBigInt("100")
	mask200, _ := HexBigInt("200")
	mask400, _ := HexBigInt("400")
	mask800, _ := HexBigInt("800")
	mask1000, _ := HexBigInt("1000")
	mask2000, _ := HexBigInt("2000")
	mask4000, _ := HexBigInt("4000")
	mask8000, _ := HexBigInt("8000")
	mask10000, _ := HexBigInt("10000")
	mask20000, _ := HexBigInt("20000")
	mask40000, _ := HexBigInt("40000")
	mask80000, _ := HexBigInt("80000")

	mulByMask2, _ := HexBigInt("fff97272373d413259a46990580e213a")
	mulByMask4, _ := HexBigInt("fff2e50f5f656932ef12357cf3c7fdcc")
	mulByMask8, _ := HexBigInt("ffe5caca7e10e4e61c3624eaa0941cd0")
	mulByMask10, _ := HexBigInt("ffcb9843d60f6159c9db58835c926644")
	mulByMask20, _ := HexBigInt("ff973b41fa98c081472e6896dfb254c0")
	mulByMask40, _ := HexBigInt("ff2ea16466c96a3843ec78b326b52861")
	mulByMask80, _ := HexBigInt("fe5dee046a99a2a811c461f1969c3053")
	mulByMask100, _ := HexBigInt("fcbe86c7900a88aedcffc83b479aa3a4")
	mulByMask200, _ := HexBigInt("f987a7253ac413176f2b074cf7815e54")
	mulByMask400, _ := HexBigInt("f3392b0822b70005940c7a398e4b70f3")
	mulByMask800, _ := HexBigInt("e7159475a2c29b7443b29c7fa6e889d9")
	mulByMask1000, _ := HexBigInt("d097f3bdfd2022b8845ad8f792aa5825")
	mulByMask2000, _ := HexBigInt("a9f746462d870fdf8a65dc1f90e061e5")
	mulByMask4000, _ := HexBigInt("70d869a156d2a1b890bb3df62baf32f7")
	mulByMask8000, _ := HexBigInt("31be135f97d08fd981231505542fcfa6")
	mulByMask10000, _ := HexBigInt("9aa508b5b7a84e1c677de54f3e99bc9")
	mulByMask20000, _ := HexBigInt("5d6af8dedb81196699c329225ee604")
	mulByMask40000, _ := HexBigInt("2216e584f5fa1ea926041bedfe98")
	mulByMask80000, _ := HexBigInt("48a170391f7dc42444e8fa2")

	var ratio *big.Int
	if notZero(absTick, mask1) {
		ratio = ratio0
	} else {
		ratio = ratio1
	}

	pairs := []MaskAndMulByPair{
		{mask2, mulByMask2},
		{mask4, mulByMask4},
		{mask8, mulByMask8},
		{mask10, mulByMask10},
		{mask20, mulByMask20},
		{mask40, mulByMask40},
		{mask80, mulByMask80},
		{mask100, mulByMask100},
		{mask200, mulByMask200},
		{mask400, mulByMask400},
		{mask800, mulByMask800},
		{mask1000, mulByMask1000},
		{mask2000, mulByMask2000},
		{mask4000, mulByMask4000},
		{mask8000, mulByMask8000},
		{mask10000, mulByMask10000},
		{mask20000, mulByMask20000},
		{mask40000, mulByMask40000},
		{mask80000, mulByMask80000},
	}

	for _, v := range pairs {
		if notZero(absTick, v.mask) {
			ratio = mulShift(ratio, v.mulBy)
		}
	}

	if tick.Cmp(zero) == 1 {
		ratio = new(big.Int).Div(maxUInt256, ratio)
	}

	result := new(big.Int).Div(ratio, Q32)

	if new(big.Int).Mod(ratio, Q32).Cmp(zero) == 1 {
		result = new(big.Int).Add(result, one)
	}

	return result, nil
}

func GetAmount0Delta(sqrtRatioAX96 *big.Int, sqrtRatioBX96 *big.Int, liquidity *big.Int) *big.Int {
	// fmt.Println("GetAmount0Delta", sqrtRatioAX96, sqrtRatioBX96, liquidity)
	var ax96 *big.Int = sqrtRatioAX96
	var bx96 *big.Int = sqrtRatioBX96
	if ax96.Cmp(bx96) == 1 {
		ax96, bx96 = bx96, ax96
	}
	numerator1 := new(big.Int).Lsh(liquidity, 96)
	numerator2 := new(big.Int).Sub(bx96, ax96)

	mulNumerators := new(big.Int).Mul(numerator1, numerator2)
	divBx96 := new(big.Int).Div(mulNumerators, bx96)
	return new(big.Int).Div(divBx96, ax96)
}

func GetAmount1Delta(sqrtRatioAX96 *big.Int, sqrtRatioBX96 *big.Int, liquidity *big.Int) *big.Int {
	var ax96 *big.Int = sqrtRatioAX96
	var bx96 *big.Int = sqrtRatioBX96
	if ax96.Cmp(bx96) == 1 {
		ax96, bx96 = bx96, ax96
	}

	ratioDiff := new(big.Int).Sub(bx96, ax96)
	mulLiqudity := new(big.Int).Mul(liquidity, ratioDiff)
	return new(big.Int).Div(mulLiqudity, Q96)
}

func GetTokenIdsForAddress(nftPosManager *nftposmgr.Nftposmgr, address common.Address, count *big.Int) ([]*big.Int, error) {
	tokenIds := []*big.Int{}
	for index := range make([]int, count.Int64()) {
		tokenId, err := nftPosManager.TokenOfOwnerByIndex(&bind.CallOpts{}, address, big.NewInt(int64(index)))
		if err != nil {
			return nil, err
		}
		time.Sleep(500 * time.Millisecond)
		tokenIds = append(tokenIds, tokenId)
	}
	return tokenIds, nil
}

func GetAmount0(tickCurrent, sqrtPriceX96, tickLower, tickUpper, liquidity *big.Int) (*big.Int, error) {
	if tickCurrent.Cmp(tickLower) == -1 {
		sqrtRatioAX96, err := GetSqrtRatioAtTick(tickLower)
		if err != nil {
			return nil, err
		}

		sqrtRatioBX96, err := GetSqrtRatioAtTick(tickUpper)
		if err != nil {
			return nil, err
		}

		return GetAmount0Delta(
			sqrtRatioAX96,
			sqrtRatioBX96,
			liquidity,
		), nil
	} else if tickCurrent.Cmp(tickUpper) == -1 {
		sqrtRatioBX96, err := GetSqrtRatioAtTick(tickUpper)
		if err != nil {
			return nil, err
		}

		return GetAmount0Delta(
			sqrtPriceX96,
			sqrtRatioBX96,
			liquidity,
		), nil
	} else {
		return big.NewInt(0), nil
	}
}

func GetAmount1(tickCurrent, sqrtPriceX96, tickLower, tickUpper, liquidity *big.Int) (*big.Int, error) {
	if tickCurrent.Cmp(tickLower) == -1 {
		return big.NewInt(0), nil
	} else if tickCurrent.Cmp(tickUpper) == -1 {
		sqrtRatioAX96, err := GetSqrtRatioAtTick(tickLower)
		if err != nil {
			return nil, err
		}

		return GetAmount1Delta(
			sqrtRatioAX96,
			sqrtPriceX96,
			liquidity,
		), nil
	} else {
		sqrtRatioAX96, err := GetSqrtRatioAtTick(tickLower)
		if err != nil {
			return nil, err
		}

		sqrtRatioBX96, err := GetSqrtRatioAtTick(tickUpper)
		if err != nil {
			return nil, err
		}

		return GetAmount1Delta(
			sqrtRatioAX96,
			sqrtRatioBX96,
			liquidity,
		), nil
	}
}

func GetAmountOfNft(client *ethclient.Client, uniswapV3Factory *v3factory.V3factory, nftPosManager *nftposmgr.Nftposmgr, tokenId *big.Int) (*AmountResult, error) {
	response, err := nftPosManager.Positions(&bind.CallOpts{}, tokenId)

	if err != nil {
		return nil, err
	}

	token0 := response.Token0
	token1 := response.Token1
	tickLower := response.TickLower
	tickUpper := response.TickUpper
	liquidity := response.Liquidity
	fee := response.Fee

	tokens := []common.Address{token0, token1}
	sort.SliceStable(tokens, func(l, r int) bool {
		return bytes.Compare(tokens[l].Bytes(), tokens[r].Bytes()) == -1
	})

	token0, token1 = tokens[0], tokens[1]
	poolAddress, err := uniswapV3Factory.GetPool(&bind.CallOpts{}, token0, token1, fee)

	if err != nil {
		return nil, err
	}

	v3Pool, err := v3pool.NewV3pool(poolAddress, client)

	if err != nil {
		return nil, err
	}

	slotResponse, err := v3Pool.Slot0(&bind.CallOpts{})

	if err != nil {
		return nil, err
	}

	tickCurrent := slotResponse.Tick
	sqrtPriceX96 := slotResponse.SqrtPriceX96

	amount0, err := GetAmount0(tickCurrent, sqrtPriceX96, tickLower, tickUpper, liquidity)
	if err != nil {
		return nil, err
	}

	amount1, err := GetAmount1(tickCurrent, sqrtPriceX96, tickLower, tickUpper, liquidity)
	if err != nil {
		return nil, err
	}

	return &AmountResult{
		tokenId,
		amount0,
		amount1,
	}, nil
}

func main() {
	url := flag.String("url", "", "url to node network (for example, infura url)")
	wallet := flag.String("wallet", "", "wallet hex address")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Usage()
	flag.Parse()

	if *url == "" {
		panic("Please set url using --url")
	}

	if *wallet == "" {
		panic("Please set wallet address using --wallet")
	}

	client, err := ethclient.Dial(*url)
	if err != nil {
		log.Fatal(err)
	}

	nftPositionManagerAddress := common.HexToAddress("0xC36442b4a4522E871399CD717aBDD847Ab11FE88")
	uniswapV3FactoryAddress := common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")
	walletAddress := common.HexToAddress(*wallet)

	nftPosMgr, err := nftposmgr.NewNftposmgr(nftPositionManagerAddress, client)

	if err != nil {
		log.Fatal(err)
	}

	uniswapV3Factory, err := v3factory.NewV3factory(uniswapV3FactoryAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	countOfNfts, err := nftPosMgr.BalanceOf(&bind.CallOpts{}, walletAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("NFTS: ", countOfNfts)
	tokenIds, err := GetTokenIdsForAddress(nftPosMgr, walletAddress, countOfNfts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Token ids: ", tokenIds)
	for _, tokenId := range tokenIds {
		result, err := GetAmountOfNft(client, uniswapV3Factory, nftPosMgr, tokenId)
		if err != nil {
			log.Fatal(err)
		}

		bytes, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(bytes))
		time.Sleep(1000 * time.Millisecond)
	}
}
