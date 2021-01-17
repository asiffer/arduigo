package w5100

import (
	"device/avr"
	"errors"
)

var localPort uint16 = 2000

// W5100 basic structure to manage the ethernet card
type W5100 struct {
	spi     *SPI
	sockets [MaxSockNum]Socket
}

// Init inits the W5100 chip
func Init() *W5100 {
	w := &W5100{spi: DefaultSPI()}
	w.ConfigureSPI()
	w.write(RMSR, 0x55)
	w.write(TMSR, 0x55)
	w.write(MR, 1<<RST)
	return w
}

// ConfigureSPI sets the SPI in "Begin" mode
//
// It starts a transaction with the parameters of
// the W5100 ethernet board
//
//  SPIE (Enable Interrupt SPI bit) = 0 : Inactivation de l'interruption SPI
// 	SPE (Enable SPI) = 1 : Active le module SPI
// 	DORD = 0 : Envoie les données en commençant par le bit de poids fort
// 	MSTR = 1 : Configure l'Arduino en mode MAÎTRE
// 	CPOL = 0 : Configure l'impulsion d'horloge inactive au niveau BAS
// 	CPHA = 0 : Valide les données sur le front sur le front montant
// 	SPR1 = 0 et SPR0 = 0: Configure la vitesse de communication à 00 = le plus rapide (Fosc/4= 4MHz)
//
func (w *W5100) ConfigureSPI() {
	// SPI Enable bit on SPCR register (SPI Control Register)
	w.spi.SPCR.SetBits(avr.SPCR_SPE)
	// Master/Slave select bit on SPCR register
	w.spi.SPCR.SetBits(avr.SPCR_MSTR)
	// SPI Mode bit on SPCR register
	// SPI_MODE0 : CPOL -> 0, CPHA -> 0
	w.spi.SPCR.ClearBits(avr.SPCR_CPOL)
	w.spi.SPCR.ClearBits(avr.SPCR_CPHA)
	// SPI data rate bit on SPCR register andSPSR register (SPI State Register)
	// SPR1 and SPR2 (SPI Clock Rate Select) bits
	// Les bits SPR configurent la fréquence du signal d'horloge. Quand l'esclave lit l’horloge d’une broche d’entrée, les bits SPR n’ont aucun effet sur l’esclave. La fréquence de l'horloge SPI est liée à la fréquence de l'oscillateur AVR. Plus le signal d'horloge SPI est rapide, plus le transfert de données sera rapide. vous devez respecter la fréquence d'horloge maximale spécifiée par l'esclave. le Le tableau suivant résume la relation entre la fréquence SCK et les bits SPR:
	// SPR1 	SPR0 	SCK frequency
	// 0 		0 		fosc/4
	// 0 		1 		fosc/16
	// 1 		0 		fosc/64
	// 1 		1 		fosc/128
	w.spi.SPCR.ClearBits(avr.SPCR_SPR0)
	w.spi.SPCR.ClearBits(avr.SPCR_SPR1)
	// select MOST Signifiant Bit first
	w.spi.SPCR.ClearBits(avr.SPCR_DORD)
	// SPI state (SPSR) to zero
	w.spi.SPSR.Set(0)
}

// In SPI Mode, W5100 operates in "unit of 32-bit stream".
// The unit of 32-bit stream  is composed of
// 	- 1 byte OP-Code Field,
// 	- 2 bytes Address Field,
// 	- 1 byte data Field.
func (w *W5100) write(addr uint16, data uint8) {
	w.spi.SS.Low()
	w.spi.Transfer(WRITE)
	w.spi.Transfer(uint8(addr >> 8))
	w.spi.Transfer(uint8(addr & 0xFF))
	w.spi.Transfer(data)
	w.spi.SS.High()
}

// writeBuffer write several bytes
func (w *W5100) writeBuffer(addr uint16, buffer []uint8) {
	for _, data := range buffer {
		w.write(addr, data)
		addr++
	}
}

// In SPI Mode, W5100 operates in "unit of 32-bit stream".
// The unit of 32-bit stream  is composed of
// 	- 1 byte OP-Code Field,
// 	- 2 bytes Address Field,
// 	- 1 byte data Field.
func (w *W5100) read(addr uint16) uint8 {
	w.spi.SS.Low()
	w.spi.Transfer(READ)
	w.spi.Transfer(uint8(addr >> 8))
	w.spi.Transfer(uint8(addr & 0xFF))
	data := w.spi.Transfer(0)
	w.spi.SS.High()
	return data
}

// readBuffer reads a register
func (w *W5100) readBuffer(addr uint16, size uint16) []uint8 {
	buffer := make([]uint8, size)
	var i uint16
	for i = 0; i < size; i++ {
		buffer[i] = w.read(addr)
		addr++
	}
	return buffer
}

// GetIPAddress returns the internal IP address
func (w *W5100) GetIPAddress() []uint8 {
	return w.readBuffer(SIPR, 4)
}

// SetIPAddress defines the internal IP address
func (w *W5100) SetIPAddress(ip []byte) {
	w.writeBuffer(SIPR, ip[:4])
}

// SetMACAddress defines the internal MAC address
func (w *W5100) SetMACAddress(mac []uint8) {
	w.writeBuffer(SHAR, mac)
}

// GetMACAddress returns the internal MAC address
func (w *W5100) GetMACAddress() []uint8 {
	return w.readBuffer(SHAR, 6)
}

// GetGatewayIP returns the IP address of the gateway
func (w *W5100) GetGatewayIP() []uint8 {
	return w.readBuffer(GWR, 4)
}

// SetGatewayIP sets the IP address of the gateway
func (w *W5100) SetGatewayIP(gw []uint8) {
	w.writeBuffer(GWR, gw[:4])
}

// GetSubnetMask returns the subnet mask
func (w *W5100) GetSubnetMask() uint8 {
	return w.read(SUBR)
}

// SetSubnetMask defines the subnet mask
func (w *W5100) SetSubnetMask(m uint8) {
	w.write(SUBR, m)
}

// initSocket prepares a socket given its id
func (w *W5100) initSocket(id uint8) *Socket {
	if id >= MaxSockNum {
		return nil
	}
	s := &Socket{uint8: id,
		wiznet: w,
		sBase:  TxBufBase + SSIZE*uint16(id),
		rBase:  RxBufBase + RSIZE*uint16(id)}
	s.Close()
	return s
}

// Socket creates a new socket
func (w *W5100) Socket(slot uint8, proto uint8, port uint16, flag uint8) (*Socket, error) {
	if slot >= MaxSockNum {
		return nil, errors.New("Socket number is greater than the maximum number of sockets")
	}

	if proto != Mode.TCP && proto != Mode.UDP && proto != Mode.IPRAW && proto != Mode.MACRAW && proto != Mode.PPPOE {
		return nil, errors.New("Bad protocol")
	}

	// first close the socket
	socket := w.initSocket(slot)
	socket.write(SocketRegister.MR, proto|flag)

	if port == 0 {
		port = localPort
		localPort++
	}

	// write the port
	socket.write16(SocketRegister.PORT, port)
	// now open the socket
	socket.exec(Command.OPEN)
	return socket, nil

}

func main() {}
