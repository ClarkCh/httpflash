package common

import (
	"sync"
)

var (
	VCLIENT uint64
	VCLICK  uint64
)

var (
	URL    string
	METHOD string
	HEADER string
	COOKIE string
	DATA   string
)

var StartMutex sync.RWMutex

var StartChannel = make(chan uint, 65536)
var StopChannel = make(chan uint, 65536)

type result struct {
	ClientNumber       int64
	ClickNumber        int64
	SeccessClickNumber int64
	FalseClickNumber   int64
	SeccessTime        int64
	FalseTime          int64
	FalseConnect       int64
}

var ResultPool result
