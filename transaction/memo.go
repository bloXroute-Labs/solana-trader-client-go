package transaction

import (
	"bytes"
	"encoding/base64"
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

func AddMemo(instructions []solana.Instruction, memoContent string, blockhash solana.Hash, owner solana.PublicKey, privateKeys map[solana.PublicKey]solana.PrivateKey) (string, error) {
	memo := CreateTraderAPIMemoInstruction(memoContent)

	instructions = append(instructions, memo)

	txnBytes, err := buildFullySignedTxn(blockhash, owner, instructions, privateKeys)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(txnBytes), nil
}

// AddMemoToSerializedTxn adds memo instruction to a serialized transaction, it's primarily used if the user
// doesn't want to interact with Trader-API directly
func AddMemoToSerializedTxn(txBase64, memoContent string,
	owner solana.PublicKey, privateKeys map[solana.PublicKey]solana.PrivateKey) (string, error) {
	signedTxBytes, err := solanarpc.DataBytesOrJSONFromBase64(txBase64)
	if err != nil {
		return "", err
	}
	unsignedTx := solanarpc.TransactionWithMeta{Transaction: signedTxBytes}
	solanaTx, err := unsignedTx.GetTransaction()
	if err != nil {
		return "", err
	}

	var instructions []solana.Instruction
	for _, cmpInst := range solanaTx.Message.Instructions {
		accounts, err := cmpInst.ResolveInstructionAccounts(&solanaTx.Message)
		if err != nil {
			return "", err
		}
		instProgID, err := solanaTx.Message.ResolveProgramIDIndex(cmpInst.ProgramIDIndex)
		if err != nil {
			return "", err
		}
		instructions = append(instructions, &solana.GenericInstruction{
			AccountValues: accounts,
			ProgID:        instProgID,
			DataBytes:     cmpInst.Data,
		})
	}

	return AddMemo(instructions, memoContent, solanaTx.Message.RecentBlockhash, owner, privateKeys)
}

func buildFullySignedTxn(recentBlockHash solana.Hash, owner solana.PublicKey, instructions []solana.Instruction, privateKeys map[solana.PublicKey]solana.PrivateKey) ([]byte, error) {

	privateKeysGetter := func(key solana.PublicKey) *solana.PrivateKey {
		if pbKey, ok := privateKeys[key]; ok {
			return &pbKey
		}
		return nil
	}

	txBuilder := solana.NewTransactionBuilder()
	for _, inst := range instructions {
		txBuilder.AddInstruction(inst)
	}

	txBuilder.SetRecentBlockHash(recentBlockHash)
	txBuilder.SetFeePayer(owner)

	tx, err := txBuilder.Build()
	if err != nil {
		return nil, err
	}
	_, err = tx.Sign(privateKeysGetter)
	if err != nil {
		return nil, err
	}
	return tx.MarshalBinary()
}
