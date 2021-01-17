package w5100

import "errors"

// Socket is a w5100 socket
type Socket struct {
	uint8         // internal socket id
	wiznet *W5100 // pointer to the parent ethernet board
	sBase  uint16
	rBase  uint16
}

// ID returns the internal socket id (from 0 to MaxSockNum)
func (sock *Socket) ID() uint8 {
	return sock.uint8
}

func (sock *Socket) computeBases() {
	sock.sBase = TxBufBase + SSIZE*uint16(sock.uint8)
	sock.rBase = RxBufBase + RSIZE*uint16(sock.uint8)
}

// exec sends a command to the socket (CR register)
func (sock *Socket) exec(cmd uint8) {
	// Send command to socket
	sock.write(SocketRegister.CR, cmd)
	// Wait for command to complete
	for sock.read(SocketRegister.CR) != 0 {
	}
}

// Close does the job
func (sock *Socket) Close() {
	sock.exec(Command.CLOSE)
	sock.write(SocketRegister.IR, 0xFF)
}

// read returns the value stored at the given address (socket register)
func (sock *Socket) read(addr uint16) uint8 {
	return sock.wiznet.read(CH_BASE + uint16(sock.uint8)*CH_SIZE + addr)
}

// write sets the value at a given address (socket register)
func (sock *Socket) write(addr uint16, data uint8) {
	sock.wiznet.write(CH_BASE+uint16(sock.uint8)*CH_SIZE+addr, data)
}

// write16 does thes same thing as write but twice (two bytes instead of one)
func (sock *Socket) write16(addr uint16, data uint16) {
	sock.write(addr, uint8(data>>8))
	sock.write(addr+1, uint8(data&0xFF))
}

// read16 reads 2 bytes
func (sock *Socket) read16(addr uint16) uint16 {
	res := uint16(sock.read(addr))
	res2 := uint16(sock.read(addr + 1))
	res = res << 8
	res2 = res2 & 0xFF
	res = res | res2
	return res
}

// readBuffer generalizes read by reading several bytes
func (sock *Socket) readBuffer(addr uint16, size uint16) []uint8 {
	return sock.wiznet.readBuffer(CH_BASE+uint16(sock.uint8)*CH_SIZE+addr, size)
}

// writeBuffer generalizes write by writing several bytes
func (sock *Socket) writeBuffer(addr uint16, buffer []uint8) {
	sock.wiznet.writeBuffer(CH_BASE+uint16(sock.uint8)*CH_SIZE+addr, buffer)
}

// Listen does the job. It establisheds the connection for the channel
// in passive (server) mode. This function waits for the request from the peer.
func (sock *Socket) Listen() error {
	if sock.read(SocketRegister.SR) != Status.INIT {
		return errors.New("The socket is not in INIT mode")
	}
	sock.exec(Command.LISTEN)
	return nil
}

// Connect establishes the connection for the channel in Active (client) mode.
// This function waits for the untill the connection is established.
func (sock *Socket) Connect(addr []uint8, port uint16) error {
	if port == 0 {
		return errors.New("Port to connect is set to zero")
	}

	// check address
	var null, full uint8

	for _, b := range addr[:4] {
		if b == 0x00 {
			null++
		} else if b == 0xFF {
			full++
		}
	}
	if full == 4 || null == 4 {
		return errors.New("The IP address to connect is not valid")
	}

	// set destination IP
	sock.writeBuffer(SocketRegister.DIPR, addr[:4])
	// set destination port
	sock.write16(SocketRegister.DPORT, port)
	// connect
	sock.exec(Command.CONNECT)

	return nil
}

// Disconnect does its job
func (sock *Socket) Disconnect() {
	sock.exec(Command.DISCON)
}

func (sock *Socket) getTXFreeSize() uint16 {
	var val uint16
	var val1 uint16 = 1

	for val != val1 {
		val1 = sock.read16(SocketRegister.TxFSR)
		if val1 != 0 {
			val = sock.read16(SocketRegister.TxFSR)
		}
	}
	return val
}

func (sock *Socket) getRXReceivedSize() uint16 {
	var val uint16
	var val1 uint16 = 1

	for val != val1 {
		val1 = sock.read16(SocketRegister.RxRSR)
		if val1 != 0 {
			val = sock.read16(SocketRegister.RxRSR)
		}
	}
	return val
}

func (sock *Socket) sendDataProcessingOffset(dataOffset uint16, data []uint8) {
	ptr := sock.read16(SocketRegister.TxWR) //readSnTX_WR(s);
	ptr += dataOffset
	offset := ptr & SMask
	dstAddr := offset + sock.sBase
	size := uint16(len(data))

	if offset+size > SSIZE {
		// Wrap around circular buffer
		// size = SSIZE - offset
		sock.wiznet.writeBuffer(dstAddr, data[:SSIZE-offset])
		sock.wiznet.writeBuffer(sock.sBase, data[SSIZE-offset:])
	} else {
		sock.wiznet.writeBuffer(dstAddr, data)
	}

	ptr += size
	sock.write16(SocketRegister.TxWR, ptr)
}

// Send sends some bytes
func (sock *Socket) Send(buf []uint8) uint16 {
	var status uint8
	var ret, freesize uint16
	size := len(buf)

	if size == 0 {
		return 0
	}

	if size > int(SSIZE) {
		ret = SSIZE // check size not to exceed MAX size.
	} else {
		ret = uint16(size)
	}

	// if freebuf is available, start.
	for freesize < ret {
		freesize = sock.getTXFreeSize()
		status = sock.read(SocketRegister.SR)
		if (status != Status.ESTABLISHED) && (status != Status.CLOSE_WAIT) {
			ret = 0
			break
		}
	}

	// copy data
	sock.sendDataProcessingOffset(0, buf)
	sock.exec(Command.SEND)

	for (sock.read(SocketRegister.IR) & Interrupt.SEND_OK) != Interrupt.SEND_OK {
		if sock.read(SocketRegister.SR) == Status.CLOSED {
			sock.Close()
			return 0
		}
	}

	sock.write(SocketRegister.IR, Interrupt.SEND_OK)
	return ret
}

func (sock *Socket) recvDataProcessing(size uint16) []uint8 {
	ptr := sock.read16(SocketRegister.RxRD)
	data := sock.readData(ptr, size)
	// no peek case
	ptr += size
	sock.write16(SocketRegister.RxRD, ptr)
	return data
}

func (sock *Socket) readData(src uint16, size uint16) []uint8 {
	srcMask := src & RMask
	srcPtr := sock.rBase + srcMask

	if (srcMask + size) > RSIZE {
		newSize := RSIZE - srcMask
		data0 := sock.wiznet.readBuffer(srcPtr, newSize)
		data1 := sock.wiznet.readBuffer(sock.rBase, size-newSize)
		return append(data0, data1...)
	}
	return sock.wiznet.readBuffer(srcPtr, size)
}

// Recv is an application I/F function which is used to receive the data in TCP mode.
// It continues to wait for data as much as the application wants to receive.
func (sock *Socket) Recv(size uint16) []uint8 {
	// Check how much data is available
	ret := sock.getRXReceivedSize()
	if ret == 0 {
		// No data available.
		status := sock.read(SocketRegister.SR)
		if status == Status.LISTEN || status == Status.CLOSED || status == Status.CLOSE_WAIT {
			// The remote end has closed its side of the connection, so this is the eof state
			return []uint8{}
		}
		// The connection is still up, but there's no data waiting to be read
		return nil

	} else if ret > size {
		ret = size
	}

	data := sock.recvDataProcessing(ret)
	sock.exec(Command.RECV)
	return data
}
