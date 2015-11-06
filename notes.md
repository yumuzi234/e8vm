## Read-only disk

```
struct {
	cmd uint8 // 0 is idle, 1 is waiting for read, 2 is read complete
	err uint8 // 0 for no error, 1 for file not found, 2 for offset overflow
	namelen uint8 // name length, max 100 bytes
	_ uint8
	
	offset uint32 // read offset
	addr uint32 // read into physical addr
	size uint32 // read number of bytes
	nread uint32 // number of bytes read

	filename [100]uint8

	_ uint32
	_ uint32
}
```

## Constants

- true and false should be consts.
- there should be also typed consts.
