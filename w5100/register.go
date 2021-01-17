package w5100

const (
	// MR is the Mode Register
	// This register is used for S/W reset,
	// memory test mode, ping block mode,
	// PPPoE mode and indirect bus I/F.
	MR uint16 = 0x0000
	// RMSR is the RX Memory Size Register
	RMSR uint16 = 0x001A
	// TMSR is the TX Memory Size Register
	TMSR uint16 = 0x001B
	// GWR is the Gateway IP Address Register
	GWR uint16 = 0x0001
	// SUBR is the Subnet Mask Register
	SUBR uint16 = 0x0005
	// SHAR is the Source Hardware Address Register
	SHAR uint16 = 0x0009
	// SIPR is the Source IP Address Register
	SIPR uint16 = 0x000F
)

// MR: Mode register (8 bits)
//    7		6	  5		4 	  3 	2 	  1 	0
// +-----+-----+-----+-----+-----+-----+-----+-----+
// | RST | 	   |	 | PB  |PPPoE|     | AI  | IND |
// +-----+-----+-----+-----+-----+-----+-----+-----+
const (
	// RST bit.
	// S/W Reset If this bit is ‘1’, internal register
	// will be initialized. It will be automatically
	// cleared after reset.
	RST uint8 = 0x07
	// PB is the Ping Mode Bloc bit
	// 0: disable Ping block
	// 1: enable Ping block
	// If the bit is set as 1, there is no response to
	// the ping request.
	PB uint8 = 0x04
	// PPPoE mode
	// 0: disable PPPoE mode
	// 1: enable PPPoE mode
	// If you use ADSL without router or etc, you should
	// set the bit as 1 to connect to ADSL Server
	PPPoE uint8 = 0x03
	// AI is the Address Auto-Increment in Indirect Bus I/F bit
	// 0: disable auto-increment
	// 1: enable auto-increment
	// At the Indirect Bus I/F mode, if this bit
	// is set as 1, the address will be automatically
	// increased by 1 whenever read and write are performed
	AI uint8 = 0x01
	// IND is the Indirect Bus I/F mode bit
	IND uint8 = 0x00
)
