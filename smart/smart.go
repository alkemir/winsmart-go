package smart

import (
	"errors"
	"fmt"
	"syscall"
)

var ErrNotSupported = errors.New("S.M.A.R.T not supported")

type getVersionInParams struct { //https://docs.microsoft.com/en-us/windows-hardware/drivers/ddi/content/ntdddisk/ns-ntdddisk-_getversioninparams
	bVersion      uint8
	bRevision     uint8
	bReserved     uint8
	bIDEDeviceMap uint8
	fCapabilities uint32
	dwReserved    [4]uint32
}

type sendCmdInParams struct {
	cBufferSize  uint32
	irDriveRegs  ideRegs
	bDriveNumber uint8
	bReserved    [3]uint8
	dwReserved   [4]uint32
	bBuffer      [1]uint8
}

type sendCmdOutParams struct {
	cBufferSize  uint32
	driverStatus driverStatus
	bBuffer      [1]uint8
}

type ideRegs struct {
	bFeaturesReg     uint8
	bSectorCountReg  uint8
	bSectorNumberReg uint8
	bCylLowReg       uint8
	bCylHighReg      uint8
	bDriveHeadReg    uint8
	bCommandReg      uint8
	bReserved        uint8
}

type driverStatus struct {
	bDriverError uint8
	bIDEError    uint8
	bReserved    [2]uint8
	dwReserved   [2]uint32
}

// volume must be something like C:
func Read(driveIndex uint8) error {
	h, err := getFileHandler(fmt.Sprintf("\\\\.\\PHYSICALDRIVE%d", driveIndex))
	if err != nil {
		return err
	}

	if err := checkSupport(h, driveIndex); err != nil {
		syscall.CloseHandle(h)
		return err
	}

	/*{
		bRet = CollectDriveInfo(hDevice, ucDriveIndex)
		bRet = ReadSMARTAttributes(hDevice, ucDriveIndex)
	}*/

	return syscall.CloseHandle(h)
}

func (d *getVersionInParams) supportsSmart() bool {
	return d.fCapabilities&CAP_SMART_CMD == CAP_SMART_CMD
}

func getFileHandler(name string) (fd syscall.Handle, err error) {
	if len(name) == 0 {
		return syscall.InvalidHandle, syscall.ERROR_FILE_NOT_FOUND
	}

	pathp, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return syscall.InvalidHandle, err
	}

	access := uint32(syscall.GENERIC_READ) | uint32(syscall.GENERIC_WRITE)
	sharemode := uint32(syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE)
	createmode := uint32(syscall.OPEN_EXISTING)

	h, e := syscall.CreateFile(pathp, access, sharemode, nil, createmode, syscall.FILE_ATTRIBUTE_SYSTEM, 0)
	return h, e
}
