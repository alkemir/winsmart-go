package smart

// GetLogicalDrives returns a slice of physical drive indexes found
func GetLogicalDrives() []uint8 {
	drives := make([]uint8, 0)
	b, _, _ := pGetLogicalDrives.Call()

	for idx := uint8(0); b != 0; idx++ {
		if b&1 == 1 {
			drives = append(drives, idx)
		}
		b = b >> 1
	}

	return drives
}
