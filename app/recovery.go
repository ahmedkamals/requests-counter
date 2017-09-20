package app

import (
	"encoding/json"
)

type (
	Recovery struct {
		Index     uint8      `json:"index"`
		TimeStamp uint64     `json:"time_stamp"`
		Data      [60]uint64 `json:"data"`
	}
)

func NewRecovery(index uint8, timestamp uint64, data [60]uint64) *Recovery {
	return &Recovery{
		Index:     index,
		TimeStamp: timestamp,
		Data:      data,
	}
}

func NewEmptyRecovery() *Recovery {
	return new(Recovery)
}

func (r *Recovery) Unmarshal(data []byte) error {
	savedData := &Recovery{}

	err := json.Unmarshal(data, savedData)

	if nil != err {
		return err
	}

	r.Index = savedData.Index
	r.TimeStamp = savedData.TimeStamp
	r.Data = savedData.Data

	return nil
}

func (r *Recovery) Marshall() ([]byte, error) {
	return json.Marshal(r)
}
