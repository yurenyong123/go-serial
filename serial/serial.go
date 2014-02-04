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
	GetAttr( string ) string
	GetAttrs( []string ) []string
	SetAttr( string, string )
	SetAttrs( map[string]string )
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
	Name   string
	Opened bool
}

func (s *SerialBase) IsOpened() bool { return s.Opened }
