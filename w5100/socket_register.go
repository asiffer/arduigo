package w5100

// SocketRegister are the socket-wise registers
var SocketRegister = struct {
	// MR is the Socket Mode Register
	// This register sets up socket option or protocol
	// type for each socket
	MR    uint16
	CR    uint16 // Command
	IR    uint16 // Interrupt
	SR    uint16 // Status
	PORT  uint16 // Source Port
	DHAR  uint16 // Destination Hardw Addr (length = 6)
	DIPR  uint16 // Destination IP Addr (length = 4)
	DPORT uint16 // Destination Port
	MSSR  uint16 // Max Segment Size
	PROTO uint16 // Protocol in IP RAW Mode
	TOS   uint16 // IP TOS
	TTL   uint16 // IP TTL
	TxFSR uint16 // TX Free Size
	TxRD  uint16 // TX Read Pointer
	TxWR  uint16 // TX Write Pointer
	RxRSR uint16 // RX Free Size
	RxRD  uint16 // RX Read Pointer
	RxWR  uint16 // RX Write Pointer (supported?)
}{
	0x0000, // Mode
	0x0001, // Command
	0x0002, // Interrupt
	0x0003, // Status
	0x0004, // Source Port
	0x0006, // Destination Hardw Addr (length = 6,
	0x000C, // Destination IP Addr (length = 4,
	0x0010, // Destination Port
	0x0012, // Max Segment Size
	0x0014, // Protocol in IP RAW Mode
	0x0015, // IP TOS
	0x0016, // IP TTL
	0x0020, // TX Free Size
	0x0022, // TX Read Pointer
	0x0024, // TX Write Pointer
	0x0026, // RX Free Size
	0x0028, // RX Read Pointer
	0x002A, // RX Write Pointer (supported?)
}

// Interrupt is the SnIR register
var Interrupt = struct {
	SEND_OK uint8
	TIMEOUT uint8
	RECV    uint8
	DISCON  uint8
	CON     uint8
}{
	0x10, // SEND_OK
	0x08, // TIMEOUT
	0x04, // RECV
	0x02, // DISCON
	0x01, // CON
}

// Status is the SnSR register
var Status = struct {
	CLOSED      uint8
	INIT        uint8
	LISTEN      uint8
	SYNSENT     uint8
	SYNRECV     uint8
	ESTABLISHED uint8
	FIN_WAIT    uint8
	CLOSING     uint8
	TIME_WAIT   uint8
	CLOSE_WAIT  uint8
	LAST_ACK    uint8
	UDP         uint8
	IPRAW       uint8
	MACRAW      uint8
	PPPOE       uint8
}{
	0x00, // CLOSED
	0x13, // INIT
	0x14, // LISTEN
	0x15, // SYNSENT
	0x16, // SYNRECV
	0x17, // ESTABLISHED
	0x18, // FIN_WAIT
	0x1A, // CLOSING
	0x1B, // TIME_WAIT
	0x1C, // CLOSE_WAIT
	0x1D, // LAST_ACK
	0x22, // UDP
	0x32, // IPRAW
	0x42, // MACRAW
	0x5F, // PPPOE
}

// Mode is the SnMR register
var Mode = struct {
	CLOSE  uint8
	TCP    uint8
	UDP    uint8
	IPRAW  uint8
	MACRAW uint8
	PPPOE  uint8
	ND     uint8
	MULTI  uint8
}{
	0x00, // CLOSE
	0x01, // TCP
	0x02, // UDP
	0x03, // IPRAW
	0x04, // MACRAW
	0x05, // PPPOE
	0x20, // ND
	0x80, // MULTI
}

// Command are the socket commands (SnCR)
var Command = struct {
	OPEN      uint8
	LISTEN    uint8
	CONNECT   uint8
	DISCON    uint8
	CLOSE     uint8
	SEND      uint8
	SEND_MAC  uint8
	SEND_KEEP uint8
	RECV      uint8
}{
	0x01, // OPEN
	0x02, // LISTEN
	0x04, // CONNECT
	0x08, // DISCON
	0x10, // CLOSE
	0x20, // SEND
	0x21, // SEND_MAC
	0x22, // SEND_KEEP
	0x40, // RECV
}
