package liquid

const (
	LCD_CLEARDISPLAY   uint8 = 0x01
	LCD_RETURNHOME     uint8 = 0x02
	LCD_ENTRYMODESET   uint8 = 0x04
	LCD_DISPLAYCONTROL uint8 = 0x08
	LCD_CURSORSHIFT    uint8 = 0x10
	LCD_FUNCTIONSET    uint8 = 0x20
	LCD_SETCGRAMADDR   uint8 = 0x40
	LCD_SETDDRAMADDR   uint8 = 0x80

	// flags for display entry mode
	LCD_ENTRYRIGHT          uint8 = 0x00
	LCD_ENTRYLEFT           uint8 = 0x02
	LCD_ENTRYSHIFTINCREMENT uint8 = 0x01
	LCD_ENTRYSHIFTDECREMENT uint8 = 0x00

	// flags for display on/off control
	LCD_DISPLAYON  uint8 = 0x04
	LCD_DISPLAYOFF uint8 = 0x00
	LCD_CURSORON   uint8 = 0x02
	LCD_CURSOROFF  uint8 = 0x00
	LCD_BLINKON    uint8 = 0x01
	LCD_BLINKOFF   uint8 = 0x00

	// flags for display/cursor shift
	LCD_DISPLAYMOVE uint8 = 0x08
	LCD_CURSORMOVE  uint8 = 0x00
	LCD_MOVERIGHT   uint8 = 0x04
	LCD_MOVELEFT    uint8 = 0x00

	// flags for function set
	LCD_8BITMODE uint8 = 0x10
	LCD_4BITMODE uint8 = 0x00
	LCD_2LINE    uint8 = 0x08
	LCD_1LINE    uint8 = 0x00
	LCD_5x10DOTS uint8 = 0x04
	LCD_5x8DOTS  uint8 = 0x00
)
