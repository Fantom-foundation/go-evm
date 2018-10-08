package state

import (
	"bytes"
	"encoding/json"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
)

type TxError struct {
	Tx     ethTypes.Transaction   `json:"tx"`
	Error  string                 `json:"error"`
}

func (te *TxError) Marshal() ([]byte, error) {
	bf := bytes.NewBuffer([]byte{})
	enc := json.NewEncoder(bf)
	if err := enc.Encode(te); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (te *TxError) Unmarshal(data []byte) error {
	bf := bytes.NewBuffer(data)
	dec := json.NewDecoder(bf)
	if err := dec.Decode(te); err != nil {
		return err
	}
	return nil
}
