package liquid

import (
	"encoding/hex"
	"machine"
	"strconv"
	"time"
)

type LCD struct {
	RS              machine.Pin
	RW              machine.Pin
	Enable          machine.Pin
	Data4           []machine.Pin
	DisplayFunction uint8
	DisplayMode     uint8
	DisplayControl  uint8
	RowOffsets      []uint8
}

func NewLCD4(rs, rw, enable machine.Pin, data []machine.Pin) *LCD {
	return &LCD{
		RS:              rs,
		RW:              rw,
		Enable:          enable,
		DisplayFunction: LCD_4BITMODE | LCD_2LINE | LCD_5x8DOTS,
		DisplayControl:  0,
		DisplayMode:     0,
		Data4:           data[:4],
		RowOffsets:      []uint8{0x00, 0x40, 0x00, 0x00},
	}
}

func (l *LCD) Begin(cols, rows uint8) {
	l.setRowOffsets(0x00, 0x40, 0x00+cols, 0x40+cols)
	l.RS.Configure(machine.PinConfig{Mode: machine.PinOutput})
	l.Enable.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for _, pin := range l.Data4 {
		pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	}
	// SEE PAGE 45/46 FOR INITIALIZATION SPECIFICATION!
	// according to datasheet, we need at least 40ms after power rises above 2.7V
	// before sending commands. Arduino can turn on way before 4.5V so we'll wait 50
	time.Sleep(100 * time.Millisecond)
	// Now we pull both RS and R/W low to begin commands
	l.RS.Low()
	l.Enable.Low()

	// 4bits mode
	// this is according to the hitachi HD44780 datasheet
	// figure 24, pg 46

	// we start in 8bit mode, try to set 4 bit mode
	l.write4(0x03)
	time.Sleep(5 * time.Millisecond) // wait min 4.1ms

	// second try
	l.write4(0x03)
	time.Sleep(5 * time.Millisecond) // wait min 4.1ms

	// third go!
	l.write4(0x03)
	time.Sleep(5 * time.Millisecond)

	// finally, set to 4-bit interface
	l.write4(0x02)

	// Then, set # lines, font size, etc.
	l.command(LCD_FUNCTIONSET | l.DisplayFunction)

	// turn the display on with no cursor or blinking default
	l.DisplayControl = LCD_DISPLAYON | LCD_CURSOROFF | LCD_BLINKOFF
	l.Display()

	// clear it off
	l.Clear()

	// Initialize to default text direction (for romance languages)
	l.DisplayMode = LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT
	// set the entry mode
	l.command(LCD_ENTRYMODESET | l.DisplayMode)
}

func (l *LCD) setRowOffsets(row0, row1, row2, row3 uint8) {
	l.RowOffsets = []uint8{row0, row1, row2, row3}
}

// BusyFlag - when the busy flag is 1, the HD44780U is in the internal operation mode,
// and the next instruction will notbe accepted. When RS = 0 and R/W = 1
// (Table 1), the busy flag is output to DB7. The next instruction must be
// written after ensuring that the busy flag is 0.
func (l *LCD) BusyFlag() bool {
	l.RS.Low()
	l.RW.High()
	return l.Data4[3].Get()
}

func (l *LCD) Clear() {
	l.command(LCD_CLEARDISPLAY)       // clear display, set cursor position to zero
	time.Sleep(20 * time.Millisecond) // this command takes a long time!
}

func (l *LCD) Display() {
	l.DisplayControl |= LCD_DISPLAYON
	l.command(LCD_DISPLAYCONTROL | l.DisplayControl)
}

func (l *LCD) Home() {
	l.command(LCD_RETURNHOME)
	time.Sleep(20 * time.Millisecond)
}

func (l *LCD) SetCursor(col, row uint8) {
	var maxLines uint8 = 2
	if row >= maxLines {
		row = maxLines - 1 // we count rows starting w/0
	}
	l.command(LCD_SETDDRAMADDR | (col + l.RowOffsets[row]))
	time.Sleep(20 * time.Millisecond)
}

func (l *LCD) FirstLine() {
	l.SetCursor(0, 0)
}

func (l *LCD) SecondLine() {
	l.SetCursor(0, 1)
}

func (l *LCD) MoveRight() {
	l.command(LCD_CURSORSHIFT | LCD_CURSORMOVE | LCD_MOVERIGHT)
}

func (l *LCD) write4(value uint8) {
	for i := 0; i < 4; i++ {
		bit := (value >> i) & 0x01
		l.Data4[i].Set(bit != 0x00)
	}
	l.pulseEnable()
}

func (l *LCD) Print(s string) {
	buffer := []uint8(s)
	for _, b := range buffer {
		l.write(b)
	}
}

func (l *LCD) printUint8Array(array []uint8, sep string) {
	length := len(array)
	for i := 0; i < length-1; i++ {
		l.PrintUint8(array[i])
		l.Print(sep)
	}
	l.PrintUint8(array[length-1])
}

func (l *LCD) PrintIP(ip []uint8) {
	for i := 0; i < 3; i++ {
		l.PrintUint8(ip[i])
		l.Print(".")
	}
	l.PrintUint8(ip[3])
}

func (l *LCD) PrintMAC(mac []uint8) {
	length := hex.EncodedLen(len(mac))
	dst := make([]byte, length)
	hex.Encode(dst, mac)

	for i := 0; i < length-1; i += 2 {
		l.Print(string(dst[i : i+2]))
		l.Print(":")
	}
	l.Print(string(dst[length-2:]))
}

func (l *LCD) PrintInt(n int) {
	l.Print(strconv.Itoa(n))
}

func (l *LCD) PrintUint8(u uint8) {
	l.Print(strconv.Itoa(int(u)))
}

func (l *LCD) CursorOff() {
	l.DisplayControl &= 255 - LCD_CURSORON
	l.command(LCD_DISPLAYCONTROL | l.DisplayControl)
}

func (l *LCD) CursorOn() {
	l.DisplayControl |= LCD_CURSORON
	l.command(LCD_DISPLAYCONTROL | l.DisplayControl)
}

func (l *LCD) pulseEnable() {
	l.Enable.Low()
	time.Sleep(10 * time.Microsecond)
	l.Enable.High()
	time.Sleep(10 * time.Microsecond)
	l.Enable.Low()
	time.Sleep(100 * time.Microsecond)
}

func (l *LCD) command(value uint8) {
	mode := l.RS.Get()
	l.send(value, false)
	l.RS.Set(mode)
}

func (l *LCD) write(value uint8) {
	mode := l.RS.Get()
	l.send(value, true)
	l.RS.Set(mode)
}

/************ low level data pushing commands **********/

// write either command or data, with automatic 4/8-bit selection
func (l *LCD) send(value uint8, mode bool) {
	if mode {
		l.RS.High()
	} else {
		l.RS.Low()
	}
	// l.RS.Set(mode)

	// if there is a RW pin indicated, set it low to Write
	if l.RW != 255 {
		l.RW.Low()
	}

	l.write4(value >> 4)
	l.write4(value)
}

func main() {}
