package w5100

import (
	"device/avr"
	"machine"
	"runtime/volatile"
)

// SPI structure to manage SPI com
type SPI struct {
	MOSI machine.Pin
	MISO machine.Pin
	SCK  machine.Pin
	SS   machine.Pin
	SPCR *volatile.Register8 // SPI Control Register
	SPSR *volatile.Register8 // SPI State Register
	SPDR *volatile.Register8 // SPI Data Register
}

// DefaultSPI creates a default AVR SPI interface and init the pins
func DefaultSPI() *SPI {
	spi := &SPI{
		MOSI: machine.D11,
		MISO: machine.D12,
		SCK:  machine.D13,
		SS:   machine.D10,
		SPCR: avr.SPCR,
		SPSR: avr.SPSR,
		SPDR: avr.SPDR,
	}

	// init the board
	spi.Init()
	return spi
}

// Init sets the pins of the board
func (spi *SPI) Init() {
	output := machine.PinConfig{Mode: machine.PinOutput}
	input := machine.PinConfig{Mode: machine.PinInput}
	// To Initialize the SPI bus by setting SCK, MOSI, and SS to outputs,
	// pulling SCK and MOSI low, and SS high.
	spi.SCK.Configure(output)
	spi.MOSI.Configure(output)
	spi.SS.Configure(output)
	spi.MISO.Configure(input)

	spi.SCK.Low()
	spi.MOSI.Low()
	spi.SS.High() // disable device (~RESET)
}

// End seems to do nothing
func (*SPI) End() {
	// does nothing?
}

// Transfer sends a single byte to the SPI bus
func (spi *SPI) Transfer(data uint8) uint8 {
	// Put data into register
	spi.SPDR.Set(data)
	// Wait until the SPIF bit is set (data is transfered)
	for !spi.SPSR.HasBits(avr.SPSR_SPIF) {
	}
	// return the value of the data register
	return spi.SPDR.Get()
}
