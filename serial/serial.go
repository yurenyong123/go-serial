package serial

import (  )

type ISerialBase interface {
	Open( string ) bool
	Close()

	Available() uint32
	Flush()
	
	Read(uint32) []uint8
	Write([]uint8)
}

type ISerialKV interface {
	GetKeys() []string
	GetStr( key string ) string
	GetStrs( keys []string ) []string
	SetStr( key string, val string )
	SetStrs( keys []string, vals []string )
}

type ISerial interface {
	Open( string ) bool
	Close()

	Available() uint32
	Flush()
	
	Read(n uint32) []uint8
	Write([]uint8)

	Baud() uint32
	SetBaud(b uint32)
	Bits() uint8
	SetBits(b uint8)
	Stops() uint8
	SetStops(stops uint8)
	Parity() string
	SetParity(p string)

	DTR() uint8
	SetDTR(state uint8)
	RTS() uint8
	SetRTS(state uint8)
	DSR() uint8
	DCD() uint8
	CTS() uint8
	RI() uint8
}

type SerialBase struct {
	Name string
	Opened bool

	Open func( string ) bool
	Close func()

	Available func() uint32
	Flush func()
	
	Read func(n uint32) []uint8
	Write func([]uint8)

	Baud func() uint32
	SetBaud func(b uint32)
	Bits func() uint8
	SetBits func(b uint8)
	Stops func() uint8
	SetStops func(stops uint8)
	Parity func() string
	SetParity func(p string)

	DTR func() uint8
	SetDTR func(state uint8)
	RTS func() uint8
	SetRTS func(state uint8)
	DSR func() uint8
	DCD func() uint8
	CTS func() uint8
	RI func() uint8
}

func (s *SerialBase) IsOpened() bool { return s.Opened }
