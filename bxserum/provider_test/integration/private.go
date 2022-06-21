package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/bloXroute-Labs/serum-client-go/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	orderSubmissionTimeout = 30 * time.Second
	txConfirmationAttempts = 3
	txConfirmationTimeout  = 5 * time.Second
	submitGoodOrders       = false
)

var publicKey string

func init() {
	var ok bool

	publicKey, ok = os.LookupEnv("PUBLIC_KEY")
	if !ok {
		fmt.Printf("env variable `PUBLIC_KEY` must be set")
		os.Exit(0)
	}
}

func testSubmitOrder(
	t *testing.T,
	submitOrderFn func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string,
	submitOrderErrFn func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string,
) {
	// submit transaction with client order ID
	if submitGoodOrders {
		ctx, cancel := context.WithTimeout(context.Background(), orderSubmissionTimeout)
		defer cancel()
		// not enabled by default due to a lack of funds without the cancel option
		txHash := submitOrderFn(ctx, publicKey, publicKey, "SOL/USDC", pb.Side_S_ASK, 0.1, 10_000, provider.PostOrderOpts{
			ClientOrderID: 5000,
		})

		verifyTx(t, txHash)
	}

	// payer mismatch: tried to use SOL public key
	ctx, cancel := context.WithTimeout(context.Background(), orderSubmissionTimeout)
	defer cancel()
	err := submitOrderErrFn(ctx, publicKey, publicKey, "SOL/USDC", pb.Side_S_BID, 0.1, 10_000, provider.PostOrderOpts{})
	require.Equal(t, "invalid payer specified: owner cannot match payer unless selling SOL", err)

	// quantity too low
	ctx, cancel = context.WithTimeout(context.Background(), orderSubmissionTimeout)
	defer cancel()
	err = submitOrderErrFn(ctx, publicKey, publicKey, "SOL/USDC", pb.Side_S_ASK, 0.0000001, 10_000, provider.PostOrderOpts{})
	require.Equal(t, "Transaction simulation failed: Error processing Instruction 2: invalid program argument", err)

	// bad open orders address
	privateKey, _ := solana.NewRandomPrivateKey()
	ctx, cancel = context.WithTimeout(context.Background(), orderSubmissionTimeout)
	defer cancel()
	err = submitOrderErrFn(ctx, publicKey, publicKey, "SOL/USDC", pb.Side_S_ASK, 0.0000001, 10_000, provider.PostOrderOpts{OpenOrdersAddress: privateKey.PublicKey().String()})
	require.Equal(t, "Transaction simulation failed: Error processing Instruction 2: invalid program argument", err)
}

type txConfirmation struct {
	TxHash string `json:"txHash"`
}

func verifyTx(t *testing.T, txHash string) {
	ok := false
	for i := 0; i < txConfirmationAttempts; i++ {
		foundHash, err := checkSOLScan(txHash)
		if err == nil {
			ok = true
			assert.Equal(t, txHash, foundHash)
		}
		time.Sleep(txConfirmationTimeout)
	}
	assert.True(t, ok)
}

func checkSOLScan(txHash string) (string, error) {
	url := fmt.Sprintf("https://public-api.solscan.io/transaction/%s", txHash)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var confirmation txConfirmation
	err = json.Unmarshal(b, &confirmation)
	if err != nil {
		return "", err
	}

	return confirmation.TxHash, nil
}
