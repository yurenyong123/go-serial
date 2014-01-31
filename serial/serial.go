package serial

import ( )

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
}
