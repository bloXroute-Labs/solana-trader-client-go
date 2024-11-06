package main

import (
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/utils"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/pumpfun_newtoken_compare/block"
	utils2 "github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

import (
	"context"
)

const updateInterval = 10 * time.Second

var (
	DurationFlag = &cli.DurationFlag{
		Name:     "run-time",
		Usage:    "amount of time to run script for (seconds)",
		Required: false,
		//Value:    time.Second * 3600, // 1 HOUR
		Value: time.Minute * 15,
		//Value: time.Second * 15,
	}
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	utils.OutputFileFlag = &cli.StringFlag{
		Name:     "output",
		Usage:    "file to output CSV results to",
		Required: false,
		Value:    "pump_fun_trader_api_comparison.csv",
	}

	app := &cli.App{
		Name:  "benchmark-traderapi-pumpfun-newtokens",
		Usage: "Compares Solana Trader API pumpfun new token stream",
		Flags: []cli.Flag{
			DurationFlag,
			utils.OutputFileFlag,
		},
		Action: run,
	}

	err = app.Run(os.Args)
	defer func() {
		if logger.Log() != nil {
			_ = logger.Log().Sync()
		}
	}()

	if err != nil {
		panic(err)
	}

}

func run(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pumpTxMap := utils2.NewLockedMap[string, benchmark.PumpTxInfo]()

	startTime := time.Now()
	duration := c.Duration(DurationFlag.Name)
	runCtx, runCancel := context.WithTimeout(ctx, duration)
	defer runCancel()
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-sigc:
				logger.Log().Info("shutdown requested")
				runCancel()
				return
			case <-ctx.Done():
				runCancel()
				return
			}
		}
	}()
	thirdPartyEndpoint, ok := os.LookupEnv("THIRD_PARTY_ENDPOINT")
	if !ok {
		log.Infof("THIRD_PARTY_ENDPOINT environment variable not set: requests will be slower")
	}
	skip3rdParty := false
	messageChan := make(chan *benchmark.NewTokenResult, 100)

	authHeader, ok := os.LookupEnv("AUTH_HEADER")
	if !ok {
		return errors.New("AUTH_HEADER not set in environment")
	}
	getBlockEndpoint, ok := os.LookupEnv("GET_BLOCK_ENDPOINT")
	if !ok {
		log.Infof("GET_BLOCK_ENDPOINT environment variable not set: requests will be slower")
	}
	firstPartyEndpoint, ok := os.LookupEnv("FIRST_PARTY_ENDPOINT")
	if !ok {
		return errors.New("FIRST_PARTY_ENDPOINT not set in environment")
	}

	if strings.Contains(thirdPartyEndpoint, ":1809") {
		// the third party is another trader-api in this case
		skip3rdParty = true
		err := startTraderAPIStream(false, runCtx, messageChan, authHeader, thirdPartyEndpoint, pumpTxMap)
		if err != nil {
			panic(err)
		}
	}

	if !skip3rdParty {
		go func() {
			err := block.StartThirdParty(
				runCtx,
				pumpTxMap,
				http.Header{},
				thirdPartyEndpoint,
				messageChan,
			)
			if err != nil {
				logger.Log().Errorw("startDetecting", "error", err)
			}
		}()
	}

	err := startTraderAPIStream(true, runCtx, messageChan, authHeader, firstPartyEndpoint, pumpTxMap)
	if err != nil {
		panic(err)
	}
	ticker := time.NewTicker(updateInterval)
	var tradeUpdates []*benchmark.NewTokenResult
	solanaRpc := rpc.New(fmt.Sprintf("https://%s", getBlockEndpoint))

Loop:
	for {
		select {
		case msg, ok := <-messageChan:
			if ok {
				populateSlotInfos(msg, solanaRpc)
				tradeUpdates = append(tradeUpdates, msg)
			}
		case <-ticker.C:
			elapsedTime := time.Now().Sub(startTime).Round(time.Second)
			logger.Log().Infof("waited %v out of %v...", elapsedTime, duration)
			if elapsedTime >= duration {
				break Loop
			}
		case <-runCtx.Done():
			break Loop

		}
	}
	time.Sleep(time.Second)

	logger.Log().Infow("finished collecting data points", "tradercount", len(tradeUpdates))

	PrintSummary(duration, tradeUpdates)

	return nil
}

func startTraderAPIStream(isFirstParty bool, runCtx context.Context, messageChan chan *benchmark.NewTokenResult, authHeader, firstPartyEndpoint string, pumpTxMap *utils2.LockedMap[string, benchmark.PumpTxInfo]) error {
	traderOS, err := stream.NewTraderWSPPumpFunNewToken(isFirstParty, messageChan, pumpTxMap, firstPartyEndpoint, authHeader)
	if err != nil {
		return err
	}

	go func() {
		var err error

		_, err = traderOS.Run(runCtx)
		if err != nil {
			panic(err)
			return
		}
	}()

	return nil
}

func PrintSummary(runtime time.Duration, datapoints []*benchmark.NewTokenResult) {
	traderFaster := 0
	tpFaster := 0
	var sumDiff, total int64
	fmt.Println("BlockTime         TraderAPIEventTime     ThirdPartyEventTime       Diff(thirdParty)       Diff(Blocktime)")
	for _, vs := range datapoints {

		total++
		fmt.Print(fmt.Sprintf("%d", vs.BlockTime.UnixMilli()))
		fmt.Print(fmt.Sprintf("     %d", vs.TraderAPIEventTime.UnixMilli()))
		fmt.Print(fmt.Sprintf("          %d", vs.ThirdPartyEventTime.UnixMilli()))
		fmt.Print(fmt.Sprintf("            %f sec", vs.Diff.Seconds()))
		fmt.Println(fmt.Sprintf("        %f sec", vs.TraderAPIEventTime.Sub(vs.BlockTime).Seconds()))
		diffMillis := vs.TraderAPIEventTime.UnixMilli() - vs.ThirdPartyEventTime.UnixMilli()
		sumDiff += diffMillis
		if vs.TraderAPIEventTime.Before(vs.ThirdPartyEventTime) {
			traderFaster++
		} else if vs.ThirdPartyEventTime.Before(vs.TraderAPIEventTime) {
			tpFaster++
		}
	}

	fmt.Println("Run time: ", runtime)
	fmt.Println()

	fmt.Println("Total events: ", total)
	fmt.Println()

	fmt.Println("Faster counts: ")
	fmt.Println(fmt.Sprintf(" traderAPIFaster   %d", traderFaster))
	fmt.Println(fmt.Sprintf(" thirdPartyFaster  %d", tpFaster))
	if total != 0 {
		fmt.Println(fmt.Sprintf(" Avg time Diff in millis  %f", float64(sumDiff/total)))
	}
}

func populateSlotInfos(msg *benchmark.NewTokenResult, solanaRpc *rpc.Client) {
	tryCount := 0
	for {
		if tryCount > 3 {
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(time.Second)
		}
		tryCount++
		if tryCount >= 10 {
			logger.Log().Infow("failed to find info for tx",
				"sig", msg.TxHash)
			break
		}

		mstv := uint64(0)
		slotInfo, err := solanaRpc.GetBlockWithOpts(
			context.Background(),
			uint64(msg.Slot),
			&rpc.GetBlockOpts{
				TransactionDetails:             rpc.TransactionDetailsNone,
				Commitment:                     rpc.CommitmentConfirmed,
				MaxSupportedTransactionVersion: &mstv,
			},
		)
		if err != nil {
			logger.Log().Errorw("error occurred when getting slot info",
				"tryCount", tryCount, "err", err)
			continue
		}
		msg.BlockTime = slotInfo.BlockTime.Time()
		return
	}

}
