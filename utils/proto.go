package utils

import (
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
)

func ConvertProtoAddressLookupTable(addressLookupTableProto map[string]*pb.PublicKeys) (map[solana.PublicKey]solana.PublicKeySlice, error) {
	addressLookupTable := make(map[solana.PublicKey]solana.PublicKeySlice)

	for pk, accounts := range addressLookupTableProto {
		solanaPk, err := solana.PublicKeyFromBase58(pk)
		if err != nil {
			return nil, err
		}

		var solanaPkSlice solana.PublicKeySlice

		for _, acc := range accounts.Pks {
			accPk, err := solana.PublicKeyFromBase58(acc)
			if err != nil {
				return nil, err
			}

			solanaPkSlice = append(solanaPkSlice, accPk)
		}

		addressLookupTable[solanaPk] = solanaPkSlice
	}

	return addressLookupTable, nil
}

func ConvertJupiterInstructions(instructions []*pb.InstructionJupiter) ([]solana.Instruction, error) {
	var solanaInstructions []solana.Instruction

	for _, inst := range instructions {
		programID, err := solana.PublicKeyFromBase58(inst.ProgramID)
		if err != nil {
			return nil, err
		}

		var accountMetaSlice solana.AccountMetaSlice

		for _, acc := range inst.Accounts {
			programID, err := solana.PublicKeyFromBase58(acc.ProgramID)
			if err != nil {
				return nil, err
			}
			accountMetaSlice = append(accountMetaSlice, solana.NewAccountMeta(
				programID, acc.IsWritable, acc.IsSigner))
		}

		solanaInstructions = append(solanaInstructions, solana.NewInstruction(programID, accountMetaSlice, inst.Data))
	}

	return solanaInstructions, nil
}

func ConvertRaydiumInstructions(instructions []*pb.InstructionRaydium) ([]solana.Instruction, error) {
	var solanaInstructions []solana.Instruction

	for _, inst := range instructions {
		programID, err := solana.PublicKeyFromBase58(inst.ProgramID)
		if err != nil {
			return nil, err
		}

		var accountMetaSlice solana.AccountMetaSlice

		for _, acc := range inst.Accounts {
			programID, err := solana.PublicKeyFromBase58(acc.ProgramID)
			if err != nil {
				return nil, err
			}
			accountMetaSlice = append(accountMetaSlice, solana.NewAccountMeta(
				programID, acc.IsWritable, acc.IsSigner))
		}

		solanaInstructions = append(solanaInstructions, solana.NewInstruction(programID, accountMetaSlice, inst.Data))
	}

	return solanaInstructions, nil
}
