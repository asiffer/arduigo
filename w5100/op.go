package w5100

//
//
// Wiznet W5100 memory map
//
// 0x0000 	+---------------------+
// 			| Common Registers    |
// 0x0030	+---------------------+
// 			| Reserved            |
// 0x0400	+---------------------+
// 			| Socket Registers    |
// 			|                     |
// 0x0800	+---------------------+
// 			| Reserved            |
// 			|                     |
// 0x4000	+---------------------+
// 			| Tx Memory           |
// 			|                     |
// 			|                     |
// 			|                     |
// 			|                     |
// 			|                     |
// 0x6000	+---------------------+
// 			| Rx Memory           |
// 			|                     |
// 			|                     |
// 			|                     |
// 			|                     |
// 			|                     |
// 0x8000	+---------------------+

// MaxSockNum is the W5100 maximum number of sockets
const MaxSockNum = 4

const (
	// TxBufBase is the Tx memory base address
	TxBufBase uint16 = 0x4000
	// RxBufBase is the Rx memory base address
	RxBufBase uint16 = 0x6000
	// TxRxMaxBufSize is the Tx/Rx maximum buffer size at the socket level (4 sockets)
	TxRxMaxBufSize uint16 = 0x2000 / 4
)

const (

	// SMask is the Tx buffer MASK
	SMask uint16 = 0x07FF
	// RMask is the Rx buffer MASK
	RMask uint16 = 0x07FF
	// SSIZE is the Max Tx buffer size
	SSIZE uint16 = 2048 // Max Tx buffer size
	// RSIZE is the Max Rx buffer size
	RSIZE   uint16 = 2048
	CH_BASE uint16 = 0x0400
	CH_SIZE uint16 = 0x0100
)

// opcodes
const (
	// READ opcode
	READ uint8 = 0x0F
	// WRITE opcode
	WRITE uint8 = 0xF0
)

// W5100 registers
