package winpriv

import (
	"bytes"
	"dbgutil"
	"encoding/binary"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

const (
	NEW_SE_PRIVILEGE_ENABLED                  uint32 = 0x00000002
	SE_ASSIGNPRIMARYTOKEN_NAME                string = "SeAssignPrimaryTokenPrivilege"
	SE_AUDIT_NAME                             string = "SeAuditPrivilege"
	SE_BACKUP_NAME                            string = "SeBackupPrivilege"
	SE_CHANGE_NOTIFY_NAME                     string = "SeChangeNotifyPrivilege"
	SE_CREATE_GLOBAL_NAME                     string = "SeCreateGlobalPrivilege"
	SE_CREATE_PAGEFILE_NAME                   string = "SeCreatePagefilePrivilege"
	SE_CREATE_PERMANENT_NAME                  string = "SeCreatePermanentPrivilege"
	SE_CREATE_SYMBOLIC_LINK_NAME              string = "SeCreateSymbolicLinkPrivilege"
	SE_CREATE_TOKEN_NAME                      string = "SeCreateTokenPrivilege"
	SE_DEBUG_NAME                             string = "SeDebugPrivilege"
	SE_DELEGATE_SESSION_USER_IMPERSONATE_NAME string = "SeDelegateSessionUserImpersonatePrivilege"
	SE_ENABLE_DELEGATION_NAME                 string = "SeEnableDelegationPrivilege"
	SE_IMPERSONATE_NAME                       string = "SeImpersonatePrivilege"
	SE_INC_BASE_PRIORITY_NAME                 string = "SeIncreaseBasePriorityPrivilege"
	SE_INCREASE_QUOTA_NAME                    string = "SeIncreaseQuotaPrivilege"
	SE_INC_WORKING_SET_NAME                   string = "SeIncreaseWorkingSetPrivilege"
	SE_LOAD_DRIVER_NAME                       string = "SeLoadDriverPrivilege"
	SE_LOCK_MEMORY_NAME                       string = "SeLockMemoryPrivilege"
	SE_MACHINE_ACCOUNT_NAME                   string = "SeMachineAccountPrivilege"
	SE_MANAGE_VOLUME_NAME                     string = "SeManageVolumePrivilege"
	SE_PROF_SINGLE_PROCESS_NAME               string = "SeProfileSingleProcessPrivilege"
	SE_RELABEL_NAME                           string = "SeRelabelPrivilege"
	SE_REMOTE_SHUTDOWN_NAME                   string = "SeRemoteShutdownPrivilege"
	SE_RESTORE_NAME                           string = "SeRestorePrivilege"
	SE_SECURITY_NAME                          string = "SeSecurityPrivilege"
	SE_SHUTDOWN_NAME                          string = "SeShutdownPrivilege"
	SE_SYNC_AGENT_NAME                        string = "SeSyncAgentPrivilege"
	SE_SYSTEM_ENVIRONMENT_NAME                string = "SeSystemEnvironmentPrivilege"
	SE_SYSTEM_PROFILE_NAME                    string = "SeSystemProfilePrivilege"
	SE_SYSTEMTIME_NAME                        string = "SeSystemtimePrivilege"
	SE_TAKE_OWNERSHIP_NAME                    string = "SeTakeOwnershipPrivilege"
	SE_TCB_NAME                               string = "SeTcbPrivilege"
	SE_TIME_ZONE_NAME                         string = "SeTimeZonePrivilege"
	SE_TRUSTED_CREDMAN_ACCESS_NAME            string = "SeTrustedCredManAccessPrivilege"
	SE_UNDOCK_NAME                            string = "SeUndockPrivilege"
	SE_UNSOLICITED_INPUT_NAME                 string = "SeUnsolicitedInputPrivilege"
)

var (
	advapi32                  = syscall.MustLoadDLL("Advapi32.dll")
	procLookupPrivilegeValueW = advapi32.MustFindProc("LookupPrivilegeValueW")
	procAdjustTokenPrivileges = advapi32.MustFindProc("AdjustTokenPrivileges")
)

func lookupPrivilegeValue_func(systemName string, name string, luid *int64) error {
	var _p0 *uint16
	var _p1 *uint16
	var err error
	var _r1 uintptr
	var _e1 syscall.Errno
	_p0, err = syscall.UTF16PtrFromString(systemName)
	if err != nil {
		return err
	}
	_p1, err = syscall.UTF16PtrFromString(name)
	if err != nil {
		return err
	}

	_r1, _, _e1 = syscall.Syscall(procLookupPrivilegeValueW.Addr(), 3, uintptr(unsafe.Pointer(_p0)), uintptr(unsafe.Pointer(_p1)), uintptr(unsafe.Pointer(luid)))
	if _r1 == 0 {
		if _e1 != 0 {
			err = error(_e1)
		} else {
			err = syscall.EINVAL
		}
		return err
	}
	return nil
}

func adjustTokenPrivilege_func(token windows.Token, input *byte, outputSize uint32, output *byte, requiredsize *uint32) error {
	var _r0 uintptr
	var _e1 syscall.Errno
	var err error = nil
	_r0, _, _e1 = syscall.Syscall6(procAdjustTokenPrivileges.Addr(), 6, uintptr(token), uintptr(0), uintptr(unsafe.Pointer(input)), uintptr(outputSize), uintptr(unsafe.Pointer(output)), uintptr(unsafe.Pointer(requiredsize)))
	if _r0 == 0 {
		if _e1 != 0 {
			err = error(_e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

func setTokenPrivilege_func(token windows.Token, sysname string, name string, enabled bool) error {
	var err error
	var val int64
	var b bytes.Buffer
	err = lookupPrivilegeValue_func(sysname, name, &val)
	if err != nil {
		return err
	}
	binary.Write(&b, binary.LittleEndian, uint32(1))
	binary.Write(&b, binary.LittleEndian, val)
	if enabled {
		binary.Write(&b, binary.LittleEndian, uint32(NEW_SE_PRIVILEGE_ENABLED))
	} else {
		binary.Write(&b, binary.LittleEndian, uint32(0))
	}

	err = adjustTokenPrivilege_func(token, &b.Bytes()[0], uint32(b.Len()), nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func SetPrivLedge(name string, enabled bool) error {
	var token windows.Token
	var h windows.Handle
	var err error
	h, err = windows.GetCurrentProcess()
	if err != nil {
		err = dbgutil.FormatError("GetCurrentProcess error[%s]", err.Error())
		return err
	}

	err = windows.OpenProcessToken(h, syscall.TOKEN_QUERY|syscall.TOKEN_ADJUST_PRIVILEGES, &token)
	if err != nil {
		err = dbgutil.FormatError("OpenProcessToken error[%s]", err.Error())
		return err
	}
	defer token.Close()
	err = setTokenPrivilege_func(token, "", name, enabled)
	if err != nil {
		if enabled {
			err = dbgutil.FormatError("enable [%s] error [%s]", name, err.Error())
		} else {
			err = dbgutil.FormatError("disable [%s] error [%s]", name, err.Error())
		}
		return err
	}
	return nil
}
