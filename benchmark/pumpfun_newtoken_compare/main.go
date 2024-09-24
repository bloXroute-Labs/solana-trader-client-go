package main

import (
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/utils"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/pumpfun_newtoken_compare/block"
	utils2 "github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"os/signal"
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
	utils.APIWSEndpoint = &cli.StringFlag{
		Name:  "solana-trader-ws-endpoint",
		Usage: "Solana Trader API API websocket connection endpoint",
		Value: "wss://pump-ny.solana.dex.blxrbdn.com/ws",
	}
	app := &cli.App{
		Name:  "benchmark-traderapi-pumpfun-newtokens",
		Usage: "Compares Solana Trader API pumpfun new token stream",
		Flags: []cli.Flag{
			utils.APIWSEndpoint,
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
	rpcHost, ok := os.LookupEnv("RPC_HOST")
	if !ok {
		log.Infof("RPC_HOST environment variable not set: requests will be slower")
	}

	go func() {
		err := block.StartBenchmarking(
			runCtx,
			pumpTxMap,
			http.Header{},
			rpcHost,
		)
		if err != nil {
			logger.Log().Errorw("startDetecting", "error", err)
		}
	}()

	traderAPIEndpoint := c.String(utils.APIWSEndpoint.Name)

	authHeader, ok := os.LookupEnv("AUTH_HEADER")
	if !ok {
		return errors.New("AUTH_HEADER not set in environment")
	}
	messageChan := make(chan *benchmark.NewTokenResult, 100)
	traderOS, err := stream.NewTraderWSPPumpFunNewToken(messageChan, pumpTxMap, traderAPIEndpoint, authHeader, rpcHost)
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

	ticker := time.NewTicker(updateInterval)
	var tradeUpdates []*benchmark.NewTokenResult
Loop:
	for {
		select {
		case msg, ok := <-messageChan:
			if ok {
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

func PrintSummary(runtime time.Duration, datapoints []*benchmark.NewTokenResult) {
	traderFaster := 0
	tpFaster := 0
	total := 0
	fmt.Println("BlockTime         TraderAPIEventTime     ThirdPartyEventTime       Diff(thirdParty)       Diff(Blocktime)")
	for _, vs := range datapoints {
		total++
		fmt.Print(fmt.Sprintf("%d", vs.BlockTime.UnixMilli()))
		fmt.Print(fmt.Sprintf("     %d", vs.TraderAPIEventTime.UnixMilli()))
		fmt.Print(fmt.Sprintf("          %d", vs.ThirdPartyEventTime.UnixMilli()))
		fmt.Print(fmt.Sprintf("            %f sec", vs.Diff.Seconds()))
		fmt.Println(fmt.Sprintf("        %f sec", vs.TraderAPIEventTime.Sub(vs.BlockTime).Seconds()))

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
}
