package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type receivedBundledReport struct {
	DataId    [32]byte
	Timestamp uint32
	Bundle    []byte
}

func encodeReceivedBundledReports(in []receivedBundledReport) ([]byte, error) {
	arrayOfTuplesType, err := abi.NewType(
		"tuple[]", "",
		[]abi.ArgumentMarshaling{
			{Name: "dataId", Type: "bytes32"},
			{Name: "timestamp", Type: "uint32"},
			{Name: "bundle", Type: "bytes"},
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
