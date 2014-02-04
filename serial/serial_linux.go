package serial

import ( "syscall"; "os"; "unsafe"; "strings"; "strconv" )

// #include <unistd.h>
// #include <sys/termios.h>
// #include <sys/ioctl.h>
import "C"


type Serial struct {
	SerialBase
	f *os.File
}


var portsPrefixes = []string{"/dev/ttyS","/dev/ttyACM","/dev/ttyUSB","/dev/rfcomm"}
const portsMaxNum = 32
var ports []string
func ScanPorts() {
	for _,pref := range portsPrefixes {
		for suf:=0; suf<=portsMaxNum; suf++ {
			fn := pref + strconv.Itoa( suf )
			_, err := os.Stat( fn )
			if err == nil {
				ports = append( ports, fn )
			}
		}
	}
}
func GetPorts() []string { return ports }


func getTermios(f *os.File, dst *C.struct_termios) {
    syscall.Syscall( syscall.SYS_IOCTL, uintptr(f.Fd()), uintptr(C.TCGETS), uintptr(unsafe.Pointer(dst)))
}
func setTermios(f *os.File, src *C.struct_termios) {
    syscall.Syscall( syscall.SYS_IOCTL, uintptr(f.Fd()), uintptr(C.TCSETS), uintptr(unsafe.Pointer(src)))
}


func (s *Serial) Open( name string ) bool {
	if ( !s.Opened ) {
		var err error
		s.f, err = os.OpenFile( name, os.O_RDWR | syscall.O_NOCTTY | syscall.O_NDELAY | syscall.O_NONBLOCK, 0666 )
		
		if (err == nil) {
			s.Name   = name
			s.Opened = true

		    var tios C.struct_termios
			getTermios(s.f, &tios)
			
		    tios.c_cflag |= C.CLOCAL | C.CREAD
		    tios.c_lflag &= C.tcflag_t( ^uint16(C.ICANON | C.ECHO | C.ECHOE | C.ECHOK | C.ECHONL | C.ISIG | C.IEXTEN))
		    tios.c_iflag &= C.tcflag_t( ^uint16(C.BRKINT | C.ICRNL | C.INPCK | C.ISTRIP | C.IXON))
		    tios.c_oflag &= C.tcflag_t( ^uint16(C.OPOST))
		    tios.c_cc[C.VMIN] = 1;
		    tios.c_cc[C.VTIME] = 0;
		
		    setTermios(s.f, &tios);
		}
	}
	return s.Opened
}
func (s *Serial) Close() {
	if s.Opened {
		s.f.Close()
		s.Opened = false
	}
}


func (s *Serial) Available() uint32 {
	if (s.Opened) {
		var nb uint32
		_,_,err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(C.TIOCINQ), uintptr(unsafe.Pointer(&nb)))
		if err != 0 { return 0 }
		return nb
	}
	return 0
}
func (s *Serial) Flush() {
	if (s.Opened) {
//		syscall.Syscall( syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(TCFLSH), uintptr(C.TCIOFLUSH) )
		C.tcflush( C.int(s.f.Fd()), C.TCIOFLUSH )
	}
}


func (s *Serial) Read(n uint32) []byte {
	buf := []byte{}
	if (s.Opened) {
		size := s.Available()
		if (n > 0) && (n < size) { size = n }
		if (size > 0) {
			buf1 := make( []byte, size, size )
			_, err := s.f.Read( buf1 )
			if ( err == nil ) { buf = buf1 }
		}
	}
	return buf
}
func (s *Serial) Write(buf []byte) {
	if s.Opened && (len(buf) > 0) {
		s.f.Write( buf )
	}
}


func (s *Serial) Baud() uint32 {
	var baud uint32 = 0
    if s.Opened {
	    var tios C.struct_termios
		getTermios(s.f, &tios)
		switch C.cfgetispeed( &tios ) {
			case C.B0:       baud = 0
			case C.B50:      baud = 50
			case C.B75:      baud = 75
			case C.B110:     baud = 110
			case C.B134:     baud = 134
			case C.B150:     baud = 150
			case C.B200:     baud = 200
			case C.B300:     baud = 300
			case C.B600:     baud = 600
			case C.B1200:    baud = 1200
			case C.B1800:    baud = 1800
			case C.B2400:    baud = 2400
			case C.B4800:    baud = 4800
			case C.B9600:    baud = 9600
			case C.B19200:   baud = 19200
			case C.B38400:   baud = 38400
			case C.B57600:   baud = 57600
			case C.B115200:  baud = 115200
			case C.B230400:  baud = 230400
			case C.B460800:  baud = 460800
			case C.B500000:  baud = 500000
			case C.B576000:  baud = 576000
			case C.B921600:  baud = 921600
			case C.B1000000: baud = 1000000
			case C.B1152000: baud = 1152000
			case C.B1500000: baud = 1500000
			case C.B2000000: baud = 2000000
			case C.B2500000: baud = 2500000
			case C.B3000000: baud = 3000000
			case C.B3500000: baud = 3500000
			case C.B4000000: baud = 4000000
		}
	}
	return baud
}
func (s *Serial) SetBaud(b uint32) {
    if s.Opened {
	    var tios C.struct_termios
		getTermios(s.f, &tios)
		
		var baud C.speed_t
		switch b {
			case 0:       baud = C.B0
			case 50:      baud = C.B50
			case 75:      baud = C.B75
			case 110:     baud = C.B110
			case 134:     baud = C.B134
			case 150:     baud = C.B150
			case 200:     baud = C.B200
			case 300:     baud = C.B300
			case 600:     baud = C.B600
			case 1200:    baud = C.B1200
			case 1800:    baud = C.B1800
			case 2400:    baud = C.B2400
			case 4800:    baud = C.B4800
			case 9600:    baud = C.B9600
			case 19200:   baud = C.B19200
			case 38400:   baud = C.B38400
			case 57600:   baud = C.B57600
			case 115200:  baud = C.B115200
			case 230400:  baud = C.B230400
			case 460800:  baud = C.B460800
			case 500000:  baud = C.B500000
			case 576000:  baud = C.B576000
			case 921600:  baud = C.B921600
			case 1000000: baud = C.B1000000
			case 1152000: baud = C.B1152000
			case 1500000: baud = C.B1500000
			case 2000000: baud = C.B2000000
			case 2500000: baud = C.B2500000
			case 3000000: baud = C.B3000000
			case 3500000: baud = C.B3500000
			case 4000000: baud = C.B4000000
			default: baud = 10000000
		}
		if (baud != 10000000) {
			C.cfsetispeed( &tios, baud )
			C.cfsetospeed( &tios, baud )
			setTermios(s.f, &tios)
		}
	}
}


func (s *Serial) Bits() uint8 {
	var bits uint8 = 0
    if s.Opened {
	    var tios C.struct_termios
		getTermios(s.f, &tios)
		
		switch (tios.c_cflag & C.CSIZE) {
			case C.CS8: bits = 8
			case C.CS7: bits = 7
			case C.CS6: bits = 6
			case C.CS5: bits = 5
		}
    }
    return bits
}
func (s *Serial) SetBits(b uint8) {
    if s.Opened {
    	switch b {
    		case 8: b = C.CS8
    		case 7: b = C.CS7
    		case 6: b = C.CS6
    		case 5: b = C.CS5
    		default: b = 0
    	}
    	if (b != 0) {
		    var tios C.struct_termios
			getTermios(s.f, &tios)
			tios.c_cflag &= C.tcflag_t( ^ uint16(C.CSIZE) )
			tios.c_cflag |= C.tcflag_t( b )
			setTermios(s.f, &tios)
    	}
    }
}

func (s *Serial) Stops() uint8 {
    var stops uint8 = 0
    if s.Opened {
	    var tios C.struct_termios
		getTermios(s.f, &tios)
    	
    	if (tios.c_cflag & C.CSTOPB) != 0 {
    		stops = 2
    	} else { stops = 1 }
    }
    return stops
}
func (s *Serial) SetStops(stops uint8) {
    if s.Opened {
	    var tios C.struct_termios
		getTermios(s.f, &tios)
		if stops == 1 {
			tios.c_cflag &= C.tcflag_t( ^uint16(C.CSTOPB))
		} else if stops == 2 { tios.c_cflag |= C.CSTOPB }
		setTermios(s.f, &tios)
    }
}

func (s *Serial) Parity() string {
    parity := "n"
    if s.Opened {
	    var tios C.struct_termios
		getTermios(s.f, &tios)
		
		if (tios.c_cflag & C.PARENB) != 0 {
			if (tios.c_cflag & C.PARODD) != 0 {
				parity = "o"
			} else {
				parity = "e"
			}
		}
	}
	return parity
}
func (s *Serial) SetParity(p string) {
    if s.Opened {
	    var tios C.struct_termios
		getTermios(s.f, &tios)
		
		switch p {
			case "n":
				tios.c_cflag &= C.tcflag_t( ^uint16(C.PARENB))
			case "o":
				tios.c_cflag |= C.PARENB | C.PARODD
			case "e":
				tios.c_cflag |= C.PARENB
				tios.c_cflag &= C.tcflag_t( ^uint16(C.PARODD))
		}

		setTermios(s.f, &tios)
	}
}


func (s *Serial) DTR() uint8 {
	var dtr uint8 = 0
    if s.Opened {
    	var status uint
    	_,_,err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(C.TIOCMGET), uintptr(unsafe.Pointer(&status)))
    	if (err == 0) {
    		if (status & C.TIOCM_DTR) != 0 { dtr = 1 }
    	}
    }
    return dtr
}
func (s *Serial) SetDTR(state uint8) {
    if s.Opened {
    	var status uint
    	_,_,err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(C.TIOCMGET), uintptr(unsafe.Pointer(&status)))
    	if (err == 0) {
    		if (state == 0) {
    			status &= ^uint(C.TIOCM_DTR)
    		} else {
    			status |= C.TIOCM_DTR
    		}
    		syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(C.TIOCMSET), uintptr(unsafe.Pointer(&status)))
    	}
    }
}

func (s *Serial) RTS() uint8 {
	var rts uint8 = 0
    if s.Opened {
    	var status uint
    	_,_,err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(C.TIOCMGET), uintptr(unsafe.Pointer(&status)))
    	if (err == 0) {
    		if (status & C.TIOCM_RTS) != 0 { rts = 1 }
    	}
    }
    return rts
}
func (s *Serial) SetRTS(state uint8) {
    if s.Opened {
    	var status uint
    	_,_,err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(C.TIOCMGET), uintptr(unsafe.Pointer(&status)))
    	if (err == 0) {
    		if (state == 0) {
    			status &= ^uint(C.TIOCM_RTS)
    		} else {
    			status |= C.TIOCM_RTS
    		}
    		syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(C.TIOCMSET), uintptr(unsafe.Pointer(&status)))
    	}
    }
}

func (s *Serial) DSR() uint8 {
	var dsr uint8 = 0
    if s.Opened {
    	var status uint
    	_,_,err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(C.TIOCMGET), uintptr(unsafe.Pointer(&status)))
    	if (err == 0) {
    		if (status & C.TIOCM_DSR) != 0 { dsr = 1 }
    	}
    }
    return dsr
}
func (s *Serial) DCD() uint8 {
	var dcd uint8 = 0
    if s.Opened {
    	var status uint
    	_,_,err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(C.TIOCMGET), uintptr(unsafe.Pointer(&status)))
    	if (err == 0) {
    		if (status & C.TIOCM_CD) != 0 { dcd = 1 }
    	}
    }
    return dcd
}
func (s *Serial) CTS() uint8 {
	var cts uint8 = 0
    if s.Opened {
    	var status uint
    	_,_,err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(C.TIOCMGET), uintptr(unsafe.Pointer(&status)))
    	if (err == 0) {
    		if (status & C.TIOCM_CTS) != 0 { cts = 1 }
    	}
    }
    return cts
}
func (s *Serial) RI() uint8 {
	var ri uint8 = 0
    if s.Opened {
    	var status uint
    	_,_,err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(s.f.Fd()), uintptr(C.TIOCMGET), uintptr(unsafe.Pointer(&status)))
    	if (err == 0) {
    		if (status & C.TIOCM_RI) != 0 { ri = 1 }
    	}
    }
    return ri
}


func (s *Serial) GetKeys() []string {
	return []string{"name","opened","baud","bits","stops","parity","dtr","dsr","rts","cts","dcd","ri"}
}
func (s *Serial) GetAttr( key string ) string {
	val := ""
	switch strings.ToLower(key) {
		case "name":   val = s.Name
		case "opened": val = strconv.FormatBool(s.Opened)
		
		case "baud":   val = strconv.Itoa(int(s.Baud()))
		case "bits":   val = strconv.Itoa(int(s.Bits()))
		case "stops":  val = strconv.Itoa(int(s.Stops()))
		case "parity": val = s.Parity()
		
		case "dtr":    val = strconv.Itoa(int(s.DTR()))
		case "dsr":    val = strconv.Itoa(int(s.DSR()))
		case "rts":    val = strconv.Itoa(int(s.RTS()))
		case "cts":    val = strconv.Itoa(int(s.CTS()))
		case "dcd":    val = strconv.Itoa(int(s.DCD()))
		case "ri":     val = strconv.Itoa(int(s.RI()))
	}
	return val
}
func (s *Serial) SetAttr( key string, val string ) {
	switch strings.ToLower(key) {
		case "baud":
			val,err := strconv.ParseInt(val, 0, 32)
			if (err == nil) { s.SetBaud( uint32( val )) }
		case "bits":
			val,err := strconv.ParseInt(val, 0, 8)
			if (err == nil) { s.SetBits( uint8( val )) }
		case "stops":
			val,err := strconv.ParseInt(val, 0, 8)
			if (err == nil) { s.SetStops( uint8( val )) }
		case "parity":
			if len(val) > 0 {
				s.SetParity( (strings.ToLower(val))[0:1] )
			}
		case "dtr":
			if strings.ToLower(val) == "false" { val = "0"
			} else if strings.ToLower(val) == "true" { val = "1" }
			val,err := strconv.ParseInt(val, 0, 8)
			if (err == nil) { s.SetDTR( uint8( val )) }
		case "rts":
			if strings.ToLower(val) == "false" { val = "0"
			} else if strings.ToLower(val) == "true" { val = "1" }
			val,err := strconv.ParseInt(val, 0, 8)
			if (err == nil) { s.SetRTS( uint8( val )) }
	}
}
func (s *Serial) GetAttrs( keys []string ) []string {
	var vals []string
	for _,key := range keys {
		vals = append(vals, s.GetAttr(key))
	}
	return vals
}
func (s *Serial) SetAttrs( attrs map[string]string ) {
	for k,v := range attrs {
		s.SetAttr( k, v )
	}
}
