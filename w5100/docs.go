// Package w5100 aims to manage the wiznet W5100 ethernet shield
//
// Examples
//
// First you have to initialize the shield. Then you can define
// the mac and ip addresses
//  w := w5100.Init()
//  w.SetMACAddress([]uint8{0x00, 0x08, 0xDC, 0xAF, 0xEE, 0x00})
//  w.SetIPAddress([]uint8{192, 168, 1, 15})
//
// You can also open a socket
//  socketID := uint8(0)
//  port := uint16(30000)
//  flags := uint8(0)
//  sock, err := w.Socket(socketID, w5100.Mode.TCP, port, flags)
package w5100
