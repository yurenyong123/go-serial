Golang serial-port OOP interface.
Examples:

```go

import ( "serial" )

//open
p := serial.Serial{}
p.Open( "/dev/ttyUSB0" )
if !p.IsOpened() { os.Exit(0) }

//parameters config
p.SetBaud(9600)
p.SetBits(8)
p.SetStops(2)
p.SetParity("n")

//line control
ri := p.RI()
p.SetRTS(0)
p.SetDTR(1)

//set/get attribute
fmt.Println( "Attributes: ", p.GetKeys() )
p.SetAttrs([]string{"baud","bits","stops","parity"}, []string{"38400","8","2","n"})
fmt.Println( p.GetAttrs( []string{"baud","bits","stops","parity"} ))

//close
p.Close()

```

