package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type receivedDecimalReport struct {
	DataId    [32]byte
	Timestamp uint32
	Answer    *big.Int
}

func encodeReceivedDecimalReport(in []receivedDecimalReport) ([]byte, error) {
	arrayOfTuplesType, err := abi.NewType(
		"tuple[]", "",
		[]abi.ArgumentMarshaling{
			{Name: "dataId", Type: "bytes32"},
			{Name: "timestamp", Type: "uint32"},
			{Name: "answer", Type: "uint224"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple[] type: %w", err)
	}
	args := abi.Arguments{
		{Name: "reports", Type: arrayOfTuplesType},
	}
	return args.Pack(in)
}
