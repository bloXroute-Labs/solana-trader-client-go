package transaction

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
)

const BxMemoMarkerMsg = "Powered by bloXroute Trader Api"

var TraderAPIMemoProgram = solana.MustPublicKeyFromBase58("HQ2UUt18uJqKaQFJhgV9zaTdQxUZjNrsKFgoEDquBkcx")

// CreateTraderAPIMemoInstruction generates a transaction instruction that places a memo in the transaction log
// Having a memo instruction with signals Trader-API usage is required
func CreateTraderAPIMemoInstruction(msg string) solana.Instruction {
	if msg == "" {
		msg = BxMemoMarkerMsg
	}
	buf := new(bytes.Buffer)
	buf.Write([]byte(msg))

	instruction := &solana.GenericInstruction{
		AccountValues: nil,
		ProgID:        TraderAPIMemoProgram,
		DataBytes:     buf.Bytes(),
	}

	return instruction
}

func addMemo(tx *solana.Transaction) error {
	memoInstruction := CreateTraderAPIMemoInstruction("")
	memoData, err := memoInstruction.Data()
	if err != nil {
		return err
	}

	cutoff := uint16(len(tx.Message.AccountKeys))
	for _, instruction := range tx.Message.Instructions {
		for i, accountIdx := range instruction.Accounts {
			if accountIdx >= cutoff {
				instruction.Accounts[i] = accountIdx + 1
			}
		}
	}

	tx.Message.AccountKeys = append(tx.Message.AccountKeys, memoInstruction.ProgramID())
	tx.Message.Instructions = append(tx.Message.Instructions, solana.CompiledInstruction{
		ProgramIDIndex: cutoff,
		Accounts:       nil,
		Data:           memoData,
	})

	return nil
}

// AddMemoAndSign adds memo instruction to a serialized transaction, it's primarily used if the user
// doesn't want to interact with Trader-API directly
func AddMemoAndSign(txBase64 string, privateKey solana.PrivateKey) (string, error) {
	signedTxBytes, err := solanarpc.DataBytesOrJSONFromBase64(txBase64)
	if err != nil {
		return "", err
	}
	unsignedTx := solanarpc.TransactionWithMeta{Transaction: signedTxBytes}
	solanaTx, err := unsignedTx.GetTransaction()
	if err != nil {
		return "", err
	}

	if len(solanaTx.Message.AccountKeys) >= 32 {
		return "", fmt.Errorf("transaction has too many account keys")
	}

	for _, key := range solanaTx.Message.AccountKeys {
		if key == TraderAPIMemoProgram {
			return "", fmt.Errorf("transaction already has bloXroute memo instruction")
		}
	}

	err = addMemo(solanaTx)
	if err != nil {
		return "", err
	}

	err = signTx(solanaTx, privateKey)
	if err != nil {
		return "", err
	}

	txnBytes, err := solanaTx.MarshalBinary()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(txnBytes), nil

}
