package block

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
	"time"
)

const PumpFunProgramID = "6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P"

//const PumpFunMintAuthorityProgramID = "TSLvdd1pWpHVjahSpsvCXUbgwsL3JAcvokwaKt1eokM"

func connect3PWs(header http.Header, rpcHost string) (*websocket.Conn, error) {
	requestBody := `{"jsonrpc": "2.0","id": "1","method": "blockSubscribe","params": ["all",  {"maxSupportedTransactionVersion":0, "commitment": "finalized", "encoding": "base64", "showRewards": true, "transactionDetails": "full"}]}`
	if strings.Contains(rpcHost, "helius") {
		rpcHost = fmt.Sprintf("wss://%s", rpcHost)
		requestBody = `{"jsonrpc":"2.0","id":420,"method":"transactionSubscribe","params":[{"vote":false,"failed":false,"accountInclude":["6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P"],"accountRequired":["TSLvdd1pWpHVjahSpsvCXUbgwsL3JAcvokwaKt1eokM"]},{"commitment":"processed","encoding":"base64","transaction_details":"full","showRewards":true,"maxSupportedTransactionVersion":0}]}`
	} else if strings.Contains(rpcHost, "160.202.128.215") {
		rpcHost = fmt.Sprintf("ws://%s/ws", rpcHost)
	} else {
		rpcHost = fmt.Sprintf("wss://%s/ws", rpcHost)
	}
	logger.Log().Infow("connecting to third party", "rpcHost", rpcHost, "requestBody", requestBody)
	ws, _, err := websocket.DefaultDialer.Dial(rpcHost, header)
	if err != nil {
		logger.Log().Errorw("dial error ", "err", err)
		return nil, err
	}

	err = ws.WriteMessage(websocket.TextMessage, []byte(requestBody))
	if err != nil {
		logger.Log().Errorw("failed to write message", "err", err)
		return nil, err
	}

	return ws, nil
}

func TransactionFromBase64(txBase64 string) (*solana.Transaction, error) {
	txBytes, err := solanarpc.DataBytesOrJSONFromBase64(txBase64)
	if err != nil {
		return nil, err
	}

	tx, err := solanarpc.TransactionWithMeta{Transaction: txBytes}.GetTransaction()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func StartBenchmarking(ctx context.Context, pumpTxMap *utils.LockedMap[string, benchmark.PumpTxInfo], header http.Header, rpcHost string) error {
	isHelius := false
	if strings.Contains(rpcHost, "helius") {
		isHelius = true
	}
	ws, err := connect3PWs(header, rpcHost)
	if err != nil {
		return err
	}
	ch := make(chan []byte, 10)
	go func() {
		for {
			_, response, err := ws.ReadMessage()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				time.Sleep(time.Second)
				logger.Log().Errorw("ReadMessage", "error", err)
				err := ws.Close()
				if err != nil {
					logger.Log().Errorw("ws.Close()", "error", err)
				}
				ws, err = connect3PWs(header, rpcHost)
				if err != nil {
					logger.Log().Errorw("connect3PWs", "error", err)
				} else {
					logger.Log().Infow("reconnected to ws successfully")
				}
				continue
			}

			ch <- response
		}
	}()

	for {

		select {
		case response := <-ch:
			if isHelius {
				processHelius(pumpTxMap, response)
			} else {
				process(pumpTxMap, response)
			}

		case <-ctx.Done():
			logger.Log().Infow("end of third party processing")
			err := ws.Close()
			if err != nil {
				logger.Log().Errorw("ws.Close()", "error", err)
			}
			return nil
		}
	}
}

func processHelius(pumpTxMap *utils.LockedMap[string, benchmark.PumpTxInfo], response []byte) {
	var tx HeliusTx
	err := json.Unmarshal(response, &tx)
	if err != nil {
		logger.Log().Debugw("Unmarshal : ", "error", err)
		return
	}
	if tx.Params.Result.Transaction.Meta.Err != nil {
		return
	}
	for _, fulltx := range tx.Params.Result.Transaction.Transaction {
		if fulltx == "base64" {
			continue
		}
		txParsed, err := TransactionFromBase64(fulltx)
		if err != nil {
			logger.Log().Errorw("TransactionFromBase64 ", "error", err)
			continue
		}

		foundPump := false

		for _, key := range txParsed.Message.AccountKeys {
			if key.String() == PumpFunProgramID {
				logger.Log().Infow("helius found pump tx")
				foundPump = true
				break
			}
		}

		if !foundPump {
			continue
		}

		for _, sig := range txParsed.Signatures {
			sigStr := sig.String()
			//logger.Log().Infow("helius signature incoming", "sig", sigStr)
			pumpTxMap.Set(sigStr, benchmark.PumpTxInfo{
				TimeSeen: time.Now(),
			})

		}

	}
}

func process(pumpTxMap *utils.LockedMap[string, benchmark.PumpTxInfo], response []byte) {
	var block FullBlock
	err := json.Unmarshal(response, &block)
	if err != nil {
		logger.Log().Errorw("Unmarshal : ", "error", err)
		return
	}

	for _, fulltx := range block.Params.Result.Value.Block.Txs {
		for _, tx := range fulltx.Transaction {
			if fulltx.Meta.Err != nil {
				continue
			}
			if tx == "base64" {
				continue
			}
			txParsed, err := TransactionFromBase64(tx)
			if err != nil {
				logger.Log().Errorw("TransactionFromBase64 ", "error", err)
				continue
			}

			foundPump := false
			now := time.Now()
			for _, key := range txParsed.Message.AccountKeys {
				if key.String() == PumpFunProgramID {
					logger.Log().Infow("3p found pump tx")
					foundPump = true
					break
				}
			}
			if !foundPump {
				continue
			}

			for _, sig := range txParsed.Signatures {
				sigStr := sig.String()
				logger.Log().Infow("rpcNode signature incoming", "sig", sigStr)
				pumpTxMap.Set(sigStr, benchmark.PumpTxInfo{
					TimeSeen: now,
				})

			}

		}
	}
}
