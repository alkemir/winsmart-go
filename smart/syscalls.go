package smart

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	SMART_GET_VERSION        uint32 = (0x0007 << 16) | (0x0001 << 14) | (0x0020 << 2) | (0) // (IOCTL_DISK_BASE) << 16) | ((FILE_READ_ACCESS) << 14) | ((0x0020) << 2) | (METHOD_BUFFERED) https://msdn.microsoft.com/en-us/library/windows/hardware/ff566202(v=vs.85).aspx
	SMART_SEND_DRIVE_COMMAND uint32 = (0x0007 << 16) | (0x0003 << 14) | (0x0021 << 2) | (0) // IOCTL_DISK_BASE << 16 | (FILE_READ_ACCESS | FILE_WRITE_ACCESS) << 14 | 33 << 2 | METHOD_BUFFERED
	CAP_SMART_CMD                   = uint32(4)
	ENABLE_SMART                    = uint8(0xD8)
	SMART_CYL_LOW                   = uint8(0x4F)
	SMART_CYL_HI                    = uint8(0xC2)
	SMART_CMD                       = uint8(0xB0)
	DRIVE_HEAD_REG                  = uint8(0xA0)
)

var (
	modkernel         = syscall.MustLoadDLL("kernel32.dll")
	pDeviceIoControl  = modkernel.MustFindProc("DeviceIoControl")  // https://docs.microsoft.com/en-us/windows/desktop/api/ioapiset/nf-ioapiset-deviceiocontrol
	pGetLogicalDrives = modkernel.MustFindProc("GetLogicalDrives") // https://docs.microsoft.com/en-us/windows/desktop/api/fileapi/nf-fileapi-getlogicaldrives
)

func checkSupport(h syscall.Handle, driveIndex uint8) error {
	driveInfo := &getVersionInParams{}
	lOut := uint32(0)

	r1, _, err := pDeviceIoControl.Call(uintptr(h), uintptr(SMART_GET_VERSION), uintptr(0), uintptr(0), uintptr(unsafe.Pointer(driveInfo)), unsafe.Sizeof(getVersionInParams{}), uintptr(unsafe.Pointer(&lOut)), uintptr(0))
	if r1 == 0 {
		return err
	}

	fmt.Println(1)

	if !driveInfo.supportsSmart() {
		return ErrNotSupported
	}

	fmt.Println(2)

	stCOP := &sendCmdOutParams{}
	stCIP := &sendCmdInParams{
		bDriveNumber: driveIndex,
		irDriveRegs: ideRegs{
			bFeaturesReg:     ENABLE_SMART,
			bSectorCountReg:  1,
			bSectorNumberReg: 1,
			bCylLowReg:       SMART_CYL_LOW,
			bCylHighReg:      SMART_CYL_HI,
			bDriveHeadReg:    DRIVE_HEAD_REG,
			bCommandReg:      SMART_CMD}}

	r1, _, err = pDeviceIoControl.Call(uintptr(h), uintptr(SMART_SEND_DRIVE_COMMAND), uintptr(unsafe.Pointer(stCIP)), uintptr(unsafe.Sizeof(sendCmdInParams{})), uintptr(unsafe.Pointer(stCOP)), unsafe.Sizeof(sendCmdOutParams{}), uintptr(unsafe.Pointer(&lOut)), uintptr(0))
	if r1 == 0 {
		return err
	}

	return nil
}
