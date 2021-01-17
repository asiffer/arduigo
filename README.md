# arduigo
Go code to manage arduino devices.

These libraries can be included in your `Go` projects
that target arduino board (with the [tinygo](https://tinygo.org/) compiler)

**Note**: The [tinygo project](https://tinygo.org/) already proposed generic [drivers](https://tinygo.org/microcontrollers/drivers/)
that can target arduino devices. However, I failed to make them work on my Uno so I have merely
reimplemented the specific legacy `C/C++` libraries into `Go`.

**Current libraries**
- [`w5100`](w5100/) to manage wiznet W5100-based ethernet shields (like the [`Ethernet`](https://github.com/Wiznet/WIZ_Ethernet_Library) library)
- [`liquid`](liquid/) to manage LCD device (like the [`LiquidCrystal`](https://github.com/arduino-libraries/LiquidCrystal) library)