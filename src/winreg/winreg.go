package winreg

import (
	"dbgutil"
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"io"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
	"winpriv"
)

var st_RootKeyMap map[string]registry.Key
var st_AccessMap map[string]uint32
var (
	modadvapi32 = windows.NewLazySystemDLL("advapi32.dll")
	modkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procRegLoadKeyW   = modadvapi32.NewProc("RegLoadKeyW")
	procRegUnLoadKeyW = modadvapi32.NewProc("RegUnLoadKeyW")
	procRegSaveKeyW   = modadvapi32.NewProc("RegSaveKeyW")
)

const (
	HKLM_ROOT     = "HKLM"
	HKCU_ROOT     = "HKCU"
	HKU_ROOT      = "HKU"
	HKCC_ROOT     = "HKCC"
	HKCR_ROOT     = "HKCR"
	USERS_ROOT    = "USERS"
	REG_ALL       = "ALL"
	REG_EXECUTE   = "EXECUTE"
	REG_QUERY     = "QUERY_VALUE"
	REG_READ      = "READ"
	REG_SET       = "SET_VALUE"
	REG_WRITE     = "WRITE"
	REG_ENUMERATE = "ENUMERATE_SUB_KEYS"
)

func init() {
	st_RootKeyMap = make(map[string]registry.Key)
	st_RootKeyMap[HKLM_ROOT] = registry.LOCAL_MACHINE
	st_RootKeyMap[HKCU_ROOT] = registry.CURRENT_USER
	st_RootKeyMap[HKCR_ROOT] = registry.CLASSES_ROOT
	st_RootKeyMap[HKU_ROOT] = registry.USERS
	st_RootKeyMap[HKCC_ROOT] = registry.CURRENT_CONFIG
	st_RootKeyMap[USERS_ROOT] = registry.USERS
	st_AccessMap = make(map[string]uint32)
	st_AccessMap[REG_ALL] = registry.ALL_ACCESS
	st_AccessMap[REG_EXECUTE] = registry.EXECUTE
	st_AccessMap[REG_QUERY] = registry.QUERY_VALUE
	st_AccessMap[REG_READ] = registry.READ
	st_AccessMap[REG_SET] = registry.SET_VALUE
	st_AccessMap[REG_WRITE] = registry.WRITE
	st_AccessMap[REG_ENUMERATE] = registry.ENUMERATE_SUB_KEYS
}

func getRootKey(root string) (key registry.Key, err error) {
	key, ok := st_RootKeyMap[root]
	if ok {
		err = nil
		return
	}
	err = dbgutil.FormatError("can not find %s", root)
	return
}

func getRegAccess(accesstype string) (access uint32, err error) {
	access = uint32(0)
	strs := strings.Split(accesstype, "|")
	for _, s := range strs {
		val, ok := st_AccessMap[s]
		if !ok {
			err = fmt.Errorf("can not find %s as type", s)
			return
		}
		access |= val
	}
	err = nil
	return
}

func ReadRegString(root, path, key string) (value string, typestr string, err error) {
	var rbuf []byte
	var rbufsize int
	var val32 uint32
	var val64 uint64
	var valtype uint32
	var retn int
	var startn, curn, curidx, cntn int
	var u16ptr []uint16
	rk, err := getRootKey(root)
	if err != nil {
		return
	}
	k, err := registry.OpenKey(rk, path, registry.QUERY_VALUE)
	if err != nil {
		err = dbgutil.FormatError("open [%s].[%s] error[%s]", root, path, err.Error())
		return
	}
	defer k.Close()

	rbuf = nil
	rbufsize = 32
	for {
		rbuf = make([]byte, rbufsize, rbufsize)
		retn, valtype, err = k.GetValue(key, rbuf)
		if err != nil {
			errtype := reflect.TypeOf(err)
			if strings.Compare(errtype.Name(), "Errno") == 0 {
				if err.(syscall.Errno) == syscall.ERROR_MORE_DATA {
					rbufsize <<= 1
					continue
				}
			}
			err = dbgutil.FormatError("read [%s].[%s].[%s] error [%s]", root, path, key, err.Error())
			return
		}
		if rbufsize > retn {
			break
		}
		rbufsize <<= 1
	}
	typestr = ""
	value = ""
	err = nil
	if valtype == registry.NONE {
		typestr = "NONE"
		value = ""
	} else if valtype == registry.SZ || valtype == registry.EXPAND_SZ {
		typestr = "SZ"
		if valtype == registry.EXPAND_SZ {
			typestr = "EXPAND_SZ"
		}
		if (retn % 2) != 0 {
			retn++
		}
		u16ptr = (*[1 << 29]uint16)(unsafe.Pointer(&rbuf[0]))[:retn]
		retn /= 2
		cntn = 0
		for cntn < retn {
			if u16ptr[cntn] == uint16(0) {
				break
			}
			cntn++
		}
		value = string(utf16.Decode(u16ptr[:cntn]))
	} else if valtype == registry.BINARY || valtype == registry.FULL_RESOURCE_DESCRIPTOR ||
		valtype == registry.RESOURCE_LIST {
		typestr = "BINARY"
		if valtype == registry.FULL_RESOURCE_DESCRIPTOR {
			typestr = "FULL_RESOURCE_DESCRIPTOR"
		}
		if valtype == registry.RESOURCE_LIST {
			typestr = "RESOURCE_LIST"
		}
		value = ""
		for i, c := range rbuf[:retn] {
			if i > 0 {
				value += ","
			}
			value += fmt.Sprintf("0x%02x", c)
		}
	} else if valtype == registry.DWORD {
		typestr = "DWORD"
		val32 = 0
		for i, c := range rbuf[:retn] {
			val32 += (uint32(c) << (uint32(i) * 8))
		}
		value = fmt.Sprintf("%d", val32)
	} else if valtype == registry.DWORD_BIG_ENDIAN {
		typestr = "DWORD_BIG_ENDIAN"
		val32 = 0
		for _, c := range rbuf[:retn] {
			val32 <<= 8
			val32 += uint32(c)
		}
		value = fmt.Sprintf("%d", val32)
	} else if valtype == registry.LINK {
		typestr = "LINK"
		value = string(rbuf[:retn])
	} else if valtype == registry.MULTI_SZ {
		typestr = "MULTI_SZ"
		startn = 0
		curidx = 0
		/*we add 1 last one*/
		if (retn % 2) != 0 {
			retn++
		}
		retn = retn / 2
		u16ptr = (*[1 << 29]uint16)(unsafe.Pointer(&rbuf[0]))[:retn]
		for startn < retn {
			curn = startn
			for curn < retn {
				if u16ptr[curn] == uint16(0x0) {
					break
				}
				curn++
			}
			if curn == startn {
				break
			}
			if curidx > 0 {
				value += "\\0"
			}
			value += string(utf16.Decode(u16ptr[startn:curn]))
			curidx++
			/*for the next one*/
			startn = curn + 1
		}
	} else if valtype == registry.QWORD {
		typestr = "QWORD"
		val64 = 0
		for i, c := range rbuf[:retn] {
			val64 += (uint64(c) << (uint32(i) * 8))
		}
		value = fmt.Sprintf("%v", val64)
	} else {
		err = dbgutil.FormatError("can not find type %d (%s\\%s)", valtype, path, key)
	}
	return
}

func parserXNumber(value string) (ival64 int64, err error) {
	if strings.HasPrefix(value, "0x") ||
		strings.HasPrefix(value, "0X") {
		ival64, err = strconv.ParseInt(value[2:], 16, 64)
	} else if strings.HasPrefix(value, "x") ||
		strings.HasPrefix(value, "X") {
		ival64, err = strconv.ParseInt(value[1:], 16, 64)
	} else {
		ival64, err = strconv.ParseInt(value, 10, 64)
	}

	return
}

func WriteRegString(root, path, key, value, typestr string) error {
	var tmpstr string
	var rbuf []byte
	var cbyte byte
	var ival64 int64
	var strs []string
	var startn, curn int
	rk, err := getRootKey(root)
	if err != nil {
		return err
	}
	k, err := registry.OpenKey(rk, path, registry.WRITE)
	if err != nil {
		return err
	}
	defer k.Close()
	if strings.Compare(typestr, "SZ") == 0 || strings.Compare(typestr, "EXPAND_SZ") == 0 {
		err = k.SetStringValue(key, value)
	} else if strings.Compare(typestr, "MULTI_SZ") == 0 {
		startn = 0
		strs = make([]string, 0, 0)
		for startn < len(value) {
			curn = startn
			for curn < len(value) {
				if value[curn] == '\\' &&
					curn < (len(value)-1) && value[(curn+1)] == '0' {
					curn += 2
					break
				}
				curn++
			}
			if curn < len(value) {
				strs = append(strs, value[startn:(curn-2)])
			} else {
				strs = append(strs, value[startn:curn])
			}
			startn = curn
		}
		err = k.SetStringsValue(key, strs)
	} else if strings.Compare(typestr, "BINARY") == 0 {
		startn = 0
		rbuf = make([]byte, 0, 0)
		for startn < len(value) {
			curn = startn
			for curn < len(value) {
				if curn == ',' {
					curn++
					break
				}
				curn++
			}
			if curn < len(value) {
				tmpstr = value[startn:(curn - 1)]
			} else {
				tmpstr = value[startn:curn]
			}

			ival64, err = parserXNumber(tmpstr)

			if err != nil {
				return err
			}

			cbyte = byte(ival64)
			rbuf = append(rbuf, cbyte)
			startn = curn
		}

		err = k.SetBinaryValue(key, rbuf)
	} else if strings.Compare(typestr, "DWORD") == 0 || strings.Compare(typestr, "DWORD_BIG_ENDIAN") == 0 {
		ival64, err = parserXNumber(value)
		if err != nil {
			return err
		}
		err = k.SetDWordValue(key, uint32(ival64))
	} else if strings.Compare(typestr, "QWORD") == 0 {
		ival64, err = parserXNumber(value)
		if err != nil {
			return err
		}
		err = k.SetQWordValue(key, uint64(ival64))
	} else {
		err = fmt.Errorf("unknown type (%s)", typestr)
	}
	return err
}

func CreateRegKey(root, path string, accesstype string, existok bool) error {
	rk, err := getRootKey(root)
	if err != nil {
		return err
	}
	access, err := getRegAccess(accesstype)
	if err != nil {
		return err
	}
	newk, existed, err := registry.CreateKey(rk, path, access)
	if err != nil {
		return err
	}
	defer newk.Close()
	if existed && !existok {
		err = fmt.Errorf("[%s] %s existed", root, path)
		return err
	}
	return nil
}

func DeleteRegKey(root, path string) error {
	rk, err := getRootKey(root)
	if err != nil {
		return err
	}

	err = registry.DeleteKey(rk, path)
	if err != nil {
		errtype := reflect.TypeOf(err)
		if strings.Compare(errtype.Name(), "Errno") == 0 {
			if err.(syscall.Errno) == syscall.ERROR_FILE_NOT_FOUND {
				/*that means not find ,so ok*/
				return nil
			}
			err = dbgutil.FormatError("errno (%d)", (int)(err.(syscall.Errno)))
		}
		return err
	}
	return nil
}

func DeleteRegValue(root, path, value string) error {
	var curk registry.Key
	rk, err := getRootKey(root)
	if err != nil {
		return err
	}

	curk, err = registry.OpenKey(rk, path, registry.WRITE|registry.EXECUTE)
	if err != nil {
		if reflect.TypeOf(err).Name() == "Errno" {
			if err.(syscall.Errno) == syscall.ERROR_FILE_NOT_FOUND {
				/*that means not find ,so ok*/
				return nil
			}
		}
		return err
	}
	defer curk.Close()

	err = curk.DeleteValue(value)
	if err != nil {
		if reflect.TypeOf(err).Name() == "Errno" {
			if err.(syscall.Errno) == syscall.ERROR_FILE_NOT_FOUND {
				/*that means not find ,so ok*/
				return nil
			}
		}

		return err
	}
	return nil
}

func EnumerateRegKeys(root, path string) (retkeys []string, err error) {
	var maxnum int = 16
	rk, err := getRootKey(root)
	if err != nil {
		return
	}

	k, err := registry.OpenKey(rk, path, registry.EXECUTE)
	if err != nil {
		return
	}
	defer k.Close()

	for {
		retkeys, err = k.ReadSubKeyNames(maxnum)
		if err != nil {
			if err != io.EOF {
				err = dbgutil.FormatError("subkeys error(%s)", err.Error())
			} else {
				err = nil
			}
		}

		if err != nil {
			return
		}

		if len(retkeys) < maxnum {
			return
		}
		maxnum <<= 1
	}

	return
}

func EnumerateRegValueKeys(root, path string) (valkeys []string, err error) {
	var maxnum int = 4
	rk, err := getRootKey(root)
	if err != nil {
		return
	}

	k, err := registry.OpenKey(rk, path, registry.EXECUTE)
	if err != nil {
		return
	}
	defer k.Close()
	for {
		valkeys = []string{}
		valkeys, err = k.ReadValueNames(maxnum)
		if err != nil {
			if err != io.EOF {
				err = dbgutil.FormatError("valuekeys error(%s)", err.Error())
			} else {
				err = nil
			}
		}

		if err != nil {
			return
		}

		if len(valkeys) < maxnum {
			return
		}
		maxnum <<= 1
	}

	return
}

func LoadHive(file, root, subkey string) (err error) {
	var hkey registry.Key
	var ufile *uint16
	var usubkey *uint16
	var r0 uintptr
	ufile, err = syscall.UTF16PtrFromString(file)
	if err != nil {
		err = dbgutil.FormatError("UTF16 from [%s] error[%s]", file, err.Error())
		return
	}
	usubkey, err = syscall.UTF16PtrFromString(subkey)
	if err != nil {
		err = dbgutil.FormatError("UTF16 from [%s] error[%s]", subkey, err.Error())
		return
	}

	hkey, err = getRootKey(root)
	if err != nil {
		return
	}

	err = winpriv.SetPrivLedge(winpriv.SE_RESTORE_NAME, true)
	if err != nil {
		return
	}
	defer winpriv.SetPrivLedge(winpriv.SE_RESTORE_NAME, false)
	err = winpriv.SetPrivLedge(winpriv.SE_BACKUP_NAME, true)
	if err != nil {
		return
	}
	defer winpriv.SetPrivLedge(winpriv.SE_BACKUP_NAME, false)

	r0, _, _ = syscall.Syscall6(procRegLoadKeyW.Addr(), 3, uintptr(syscall.Handle(hkey)), uintptr(unsafe.Pointer(usubkey)), uintptr(unsafe.Pointer(ufile)), 0, 0, 0)
	if r0 != 0 {
		err = dbgutil.FormatError("load [%s] => [%s].[%s] error[%d]", file, root, subkey, r0)
		return
	}
	err = nil
	return
}

func UnLoadHive(root, subkey string) (err error) {
	var hkey registry.Key
	var usubkey *uint16
	var r0 uintptr
	usubkey, err = syscall.UTF16PtrFromString(subkey)
	if err != nil {
		err = dbgutil.FormatError("UTF16 from [%s] error[%s]", subkey, err.Error())
		return
	}

	hkey, err = getRootKey(root)
	if err != nil {
		return
	}

	err = winpriv.SetPrivLedge(winpriv.SE_RESTORE_NAME, true)
	if err != nil {
		return
	}
	defer winpriv.SetPrivLedge(winpriv.SE_RESTORE_NAME, false)
	err = winpriv.SetPrivLedge(winpriv.SE_BACKUP_NAME, true)
	if err != nil {
		return
	}
	defer winpriv.SetPrivLedge(winpriv.SE_BACKUP_NAME, false)

	r0, _, _ = syscall.Syscall6(procRegUnLoadKeyW.Addr(), 2, uintptr(syscall.Handle(hkey)), uintptr(unsafe.Pointer(usubkey)), 0, 0, 0, 0)
	if r0 != 0 {
		err = dbgutil.FormatError("unload [%s].[%s] error[%d]", root, subkey, r0)
		return
	}
	err = nil
	return
}

func SaveHive(file, root, subkey string) (err error) {
	var hkey registry.Key
	var ufile *uint16
	var r0 uintptr
	var k registry.Key
	ufile, err = syscall.UTF16PtrFromString(file)
	if err != nil {
		err = dbgutil.FormatError("UTF16 from [%s] error[%s]", file, err.Error())
		return
	}

	hkey, err = getRootKey(root)
	if err != nil {
		return
	}

	k, err = registry.OpenKey(hkey, subkey, registry.READ)
	if err != nil {
		err = dbgutil.FormatError("open [%s].[%s] error[%s]", root, subkey, err.Error())
		return
	}
	defer k.Close()

	err = winpriv.SetPrivLedge(winpriv.SE_BACKUP_NAME, true)
	if err != nil {
		return
	}
	defer winpriv.SetPrivLedge(winpriv.SE_BACKUP_NAME, false)

	r0, _, _ = syscall.Syscall6(procRegSaveKeyW.Addr(), 3, uintptr(syscall.Handle(k)), uintptr(unsafe.Pointer(ufile)), 0, 0, 0, 0)
	if r0 != 0 {
		err = dbgutil.FormatError("save [%s].[%s] => [%s] error[%d]", root, subkey, file, r0)
		return
	}
	err = nil
	return
}
