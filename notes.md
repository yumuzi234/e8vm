- the vm spawns a small terminal of 80x40 chars
- each char has a foreground color and a background color
3byte for foreground, 3-byte for background
1 byte for the content, and 1 byte for flags
this will require 25k of memory for the frame buffer