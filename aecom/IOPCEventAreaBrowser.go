package aecom

import (
	"syscall"
	"unsafe"

	"github.com/huskar-t/opcda/com"
	"golang.org/x/sys/windows"
)

// MIDL_DEFINE_GUID(IID, IID_IOPCEventAreaBrowser,0x65168857,0x5783,0x11D1,0x84,0xA0,0x00,0x60,0x8C,0xB8,0xA7,0xE9);
var IID_IOPCEventAreaBrowser = windows.GUID{
	Data1: 0x65168857,
	Data2: 0x5783,
	Data3: 0x11D1,
	Data4: [8]byte{0x84, 0xA0, 0x00, 0x60, 0x8C, 0xB8, 0xA7, 0xE9},
}

type IOPCEventAreaBrowser struct {
	*com.IUnknown
}

type IOPCEventAreaBrowserVtbl struct {
	com.IUnknownVtbl
	ChangeBrowsePosition   uintptr
	BrowseOPCAreas         uintptr
	GetQualifiedAreaName   uintptr
	GetQualifiedSourceName uintptr
}

func (v *IOPCEventAreaBrowser) VTable() *IOPCEventAreaBrowserVtbl {
	return (*IOPCEventAreaBrowserVtbl)(unsafe.Pointer(v.IUnknown.LpVtbl))
}

// virtual HRESULT STDMETHODCALLTYPE ChangeBrowsePosition(
// /* [in] */ OPCAEBROWSEDIRECTION dwBrowseDirection,
// /* [string][in] */ LPCWSTR szString) = 0;
func (v *IOPCEventAreaBrowser) ChangeBrowsePosition(browseDirection uint32, str string) (err error) {
	r0, _, _ := syscall.SyscallN(
		v.VTable().ChangeBrowsePosition,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(browseDirection),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(str))),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE BrowseOPCAreas(
// /* [in] */ OPCAEBROWSETYPE dwBrowseFilterType,
// /* [string][in] */ LPCWSTR szFilterCriteria,
// /* [out] */ LPENUMSTRING *ppIEnumString) = 0;
func (v *IOPCEventAreaBrowser) BrowseOPCAreas(browseFilterType uint32, filterCriteria string) (areas []string, err error) {
	var pString *com.IUnknown
	r0, _, _ := syscall.SyscallN(
		v.VTable().BrowseOPCAreas,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(browseFilterType),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(filterCriteria))),
		uintptr(unsafe.Pointer(&pString)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	ppIEnumString := &com.IEnumString{IUnknown: pString}
	defer ppIEnumString.Release()
	for {
		batch, err := ppIEnumString.Next(100)
		if err != nil {
			return nil, err
		}
		areas = append(areas, batch...)
		if len(batch) < 100 {
			break
		}
	}
	return
}

//virtual HRESULT STDMETHODCALLTYPE GetQualifiedAreaName(
///* [in] */ LPCWSTR szAreaName,
///* [string][out] */ LPWSTR *pszQualifiedAreaName) = 0;
//

func (v *IOPCEventAreaBrowser) GetQualifiedAreaName(areaName string) (qualifiedAreaName string, err error) {
	var pszQualifiedAreaName *uint16
	r0, _, _ := syscall.SyscallN(
		v.VTable().GetQualifiedAreaName,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(areaName))),
		uintptr(unsafe.Pointer(&pszQualifiedAreaName)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	qualifiedAreaName = windows.UTF16PtrToString(pszQualifiedAreaName)
	return
}

// virtual HRESULT STDMETHODCALLTYPE GetQualifiedSourceName(
// /* [in] */ LPCWSTR szSourceName,
// /* [string][out] */ LPWSTR *pszQualifiedSourceName) = 0;
func (v *IOPCEventAreaBrowser) GetQualifiedSourceName(SourceName string) (qualifiedSourceName string, err error) {
	var pszQualifiedSourceName *uint16
	r0, _, _ := syscall.SyscallN(
		v.VTable().GetQualifiedSourceName,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(SourceName))),
		uintptr(unsafe.Pointer(&pszQualifiedSourceName)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	qualifiedSourceName = windows.UTF16PtrToString(pszQualifiedSourceName)
	return
}
