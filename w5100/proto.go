package w5100

// IPProto is the socket protocol
var IPProto = struct {
	IP   uint8
	ICMP uint8
	IGMP uint8
	GGP  uint8
	TCP  uint8
	PUP  uint8
	UDP  uint8
	IDP  uint8
	ND   uint8
	RAW  uint8
}{
	0,   // IP
	1,   // ICMP
	2,   // IGMP
	3,   // GGP
	6,   // TCP
	12,  // PUP
	17,  // UDP
	22,  // IDP
	77,  // ND
	255, // RAW
}
