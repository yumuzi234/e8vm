# IO memory mapping

## Page 0: reserved

## Page 1: interrupt

0-4: flags
4-8: handler SP
8-c: handler PC
c-10: syscall SP
10-14: syscall PC
20-40: interrupt enabling mask
40-60: interrupt pending bits

## Page 2: Basic IO

0: console output byte
1: is console output byte valid
4: console input byte
5: is console input byte valid

8-c: boot arg

10: mouse click input valid
11: mouse click signal
12: mouse click line
13: mouse click col

20: read flags
24: rand number
28: time high
2c: time low

80-84: serial input head
84-88: serial input tail
88-8c: cycles to wait to raise an input interrupt
8c-90: input threshold

90-94: serial output head
94-98: serial output tail
98-9c: cycles to wait
9c-a0: output threshold

c0-e0: serial input ring buffer
e0-100: serial output ring buffer

100: rom command
101: rom name length
102: rom state
103: rom error
104-108: rom offset to read
108-10c: rom address to write
10c-110: rom size of buffer

110-114: rom number of bytes read
114-178: rom file name, max 100 chars

## Page 3: simple network

0: if has incoming packet ready
1: incoming packet flag // reserved
2-4: incoming packet length
8: incoming packet payload

800: if has outgoing packet ready
801: outgoing packet flag // reserved
802-804: outgoing packet length
808: outgoing packet payload

## Page 5: Screen text frame buffer

size 80x24, one byte for each char

## Page 6: Screen color frame buffer

size 80x24, one byte for each char, foreground and background

## Page 7: System information

0-4: number of pages for the physical memory
4-8: number of cores

## Page 8: Start of boot image.
