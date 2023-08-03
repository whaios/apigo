package simple

import (
	"encoding/json"
	"time"
)

type (
	SomeStruct struct {
		Bool    bool                      `json:"bool"`                       // 布尔
		Int     int                       `json:"int"`                        // 整数
		Int64   int64                     `json:"int64"`                      // 大整数
		Float64 float64                   `json:"float64"`                    // 浮点数
		String  string                    `json:"string" validate:"required"` // 字符串
		Bytes   []byte                    `json:"bytes"`
		JSON    json.RawMessage           `json:"json"`
		Time    time.Time                 `json:"time"`
		Slice   []SomeOtherType           `json:"slice"`
		Map     map[string]*SomeOtherType `json:"map"`

		Struct struct {
			X string `json:"x" validate:"required"`
		} `json:"struct"`

		EmptyStruct struct {
			Y string
		} `json:"structWithoutFields"`

		Ptr *SomeOtherType `json:"ptr"`
	}

	SomeOtherType string
)
