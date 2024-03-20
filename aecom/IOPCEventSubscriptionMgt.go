package aecom

import (
	"syscall"
	"unsafe"

	"github.com/huskar-t/opcda/com"
	"golang.org/x/sys/windows"
)

// IID_IOPCEventSubscriptionMgt 65168855-5783-11D1-84A0-00608CB8A7E9
var IID_IOPCEventSubscriptionMgt = windows.GUID{
	Data1: 0x65168855,
	Data2: 0x5783,
	Data3: 0x11D1,
	Data4: [8]byte{0x84, 0xA0, 0x00, 0x60, 0x8C, 0xB8, 0xA7, 0xE9},
}

type IOPCEventSubscriptionMgt struct {
	*com.IUnknown
}
type IOPCEventSubscriptionMgtVtbl struct {
	com.IUnknownVtbl
	SetFilter                uintptr
	GetFilter                uintptr
	SelectReturnedAttributes uintptr
	GetReturnedAttributes    uintptr
	Refresh                  uintptr
	CancelRefresh            uintptr
	GetState                 uintptr
	SetState                 uintptr
}

func (v *IOPCEventSubscriptionMgt) VTable() *IOPCEventSubscriptionMgtVtbl {
	return (*IOPCEventSubscriptionMgtVtbl)(unsafe.Pointer(v.IUnknown.LpVtbl))
}

// virtual HRESULT STDMETHODCALLTYPE SetFilter(
// /* [in] */ DWORD dwEventType,
// /* [in] */ DWORD dwNumCategories,
// /* [size_is][in] */ DWORD *pdwEventCategories,
// /* [in] */ DWORD dwLowSeverity,
// /* [in] */ DWORD dwHighSeverity,
// /* [in] */ DWORD dwNumAreas,
// /* [size_is][in] */ LPWSTR *pszAreaList,
// /* [in] */ DWORD dwNumSources,
// /* [size_is][in] */ LPWSTR *pszSourceList) = 0;
func (v *IOPCEventSubscriptionMgt) SetFilter(eventType uint32, eventCategories []uint32, lowSeverity uint32, highSeverity uint32, areaList []string, sourceList []string) (err error) {
	var ppAreaList, ppSourceList unsafe.Pointer
	if len(areaList) > 0 {
		pAreaList := make([]*uint16, len(areaList))
		for i, a := range areaList {
			pAreaList[i], err = syscall.UTF16PtrFromString(a)
			if err != nil {
				return
			}
		}
		ppAreaList = unsafe.Pointer(&pAreaList[0])
	}
	if len(sourceList) > 0 {
		pSourceList := make([]*uint16, len(sourceList))
		for i, s := range sourceList {
			pSourceList[i], err = syscall.UTF16PtrFromString(s)
			if err != nil {
				return
			}
		}
		ppSourceList = unsafe.Pointer(&pSourceList[0])
	}
	r0, _, _ := syscall.SyscallN(
		v.VTable().SetFilter,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(eventType),
		uintptr(len(eventCategories)),
		uintptr(unsafe.Pointer(&eventCategories[0])),
		uintptr(lowSeverity),
		uintptr(highSeverity),
		uintptr(len(areaList)),
		uintptr(ppAreaList),
		uintptr(len(sourceList)),
		uintptr(ppSourceList))
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE GetFilter(
//
//	/* [out] */ DWORD *pdwEventType,
//	/* [out] */ DWORD *pdwNumCategories,
//	/* [size_is][size_is][out] */ DWORD **ppdwEventCategories,
//	/* [out] */ DWORD *pdwLowSeverity,
//	/* [out] */ DWORD *pdwHighSeverity,
//	/* [out] */ DWORD *pdwNumAreas,
//	/* [size_is][size_is][out] */ LPWSTR **ppszAreaList,
//	/* [out] */ DWORD *pdwNumSources,
//	/* [size_is][size_is][out] */ LPWSTR **ppszSourceList) = 0;
func (v *IOPCEventSubscriptionMgt) GetFilter() (eventType uint32, eventCategories []uint32, lowSeverity uint32, highSeverity uint32, areaList []string, sourceList []string, err error) {
	var pNumCategories, pNumAreas, pNumSources uint32
	var ppEventCategories, ppAreaList, ppSourceList unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		v.VTable().GetFilter,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(&eventType)),
		uintptr(unsafe.Pointer(&pNumCategories)),
		uintptr(unsafe.Pointer(&ppEventCategories)),
		uintptr(unsafe.Pointer(&lowSeverity)),
		uintptr(unsafe.Pointer(&highSeverity)),
		uintptr(unsafe.Pointer(&pNumAreas)),
		uintptr(unsafe.Pointer(&ppAreaList)),
		uintptr(unsafe.Pointer(&pNumSources)),
		uintptr(unsafe.Pointer(&ppSourceList)),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
		return
	}
	eventCategories = make([]uint32, pNumCategories)
	for i := uint32(0); i < pNumCategories; i++ {
		eventCategories[i] = *(*uint32)(unsafe.Pointer(uintptr(ppEventCategories) + uintptr(i)*4))
	}
	areaList = make([]string, pNumAreas)
	for i := uint32(0); i < pNumAreas; i++ {
		pwstr := *(**uint16)(unsafe.Pointer(uintptr(ppAreaList) + uintptr(i)*pointerSize))
		areaList[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
	}
	sourceList = make([]string, pNumSources)
	for i := uint32(0); i < pNumSources; i++ {
		pwstr := *(**uint16)(unsafe.Pointer(uintptr(ppSourceList) + uintptr(i)*pointerSize))
		sourceList[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE SelectReturnedAttributes(
// /* [in] */ DWORD dwEventCategory,
// /* [in] */ DWORD dwCount,
// /* [size_is][in] */ DWORD *dwAttributeIDs) = 0;
func (v *IOPCEventSubscriptionMgt) SelectReturnedAttributes(eventCategory uint32, attributeIDs []uint32) (err error) {
	var pAttributeIDs unsafe.Pointer
	if len(attributeIDs) != 0 {
		pAttributeIDs = unsafe.Pointer(&attributeIDs[0])
	}
	r0, _, _ := syscall.SyscallN(
		v.VTable().SelectReturnedAttributes,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(eventCategory),
		uintptr(len(attributeIDs)),
		uintptr(pAttributeIDs),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE GetReturnedAttributes(
// /* [in] */ DWORD dwEventCategory,
// /* [out] */ DWORD *pdwCount,
// /* [size_is][size_is][out] */ DWORD **ppdwAttributeIDs) = 0;
func (v *IOPCEventSubscriptionMgt) GetReturnedAttributes(eventCategory uint32) (attributeIDs []uint32, err error) {
	var pCount uint32
	var ppAttributeIDs unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		v.VTable().GetReturnedAttributes,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(eventCategory),
		uintptr(unsafe.Pointer(&pCount)),
		uintptr(unsafe.Pointer(&ppAttributeIDs)),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
		return
	}
	attributeIDs = make([]uint32, pCount)
	for i := uint32(0); i < pCount; i++ {
		attributeIDs[i] = *(*uint32)(unsafe.Pointer(uintptr(ppAttributeIDs) + uintptr(i)*4))
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE Refresh(
//
//	/* [in] */ DWORD dwConnection) = 0;
func (v *IOPCEventSubscriptionMgt) Refresh(connection uint32) (err error) {
	r0, _, _ := syscall.SyscallN(
		v.VTable().Refresh,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(connection),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE CancelRefresh(
//
//	/* [in] */ DWORD dwConnection) = 0;
func (v *IOPCEventSubscriptionMgt) CancelRefresh(connection uint32) (err error) {
	r0, _, _ := syscall.SyscallN(
		v.VTable().CancelRefresh,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(connection),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE GetState(
// /* [out] */ BOOL *pbActive,
// /* [out] */ DWORD *pdwBufferTime,
// /* [out] */ DWORD *pdwMaxSize,
// /* [out] */ OPCHANDLE *phClientSubscription) = 0;
func (v *IOPCEventSubscriptionMgt) GetState() (active bool, bufferTime uint32, maxSize uint32, clientSubscription uint32, err error) {
	var pActive int32
	r0, _, _ := syscall.SyscallN(
		v.VTable().GetState,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(&pActive)),
		uintptr(unsafe.Pointer(&bufferTime)),
		uintptr(unsafe.Pointer(&maxSize)),
		uintptr(unsafe.Pointer(&clientSubscription)),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
		return
	}
	active = pActive != 0
	return
}

// virtual HRESULT STDMETHODCALLTYPE SetState(
// /* [in][unique] */ BOOL *pbActive,
// /* [in][unique] */ DWORD *pdwBufferTime,
// /* [in][unique] */ DWORD *pdwMaxSize,
// /* [in] */ OPCHANDLE hClientSubscription,
// /* [out] */ DWORD *pdwRevisedBufferTime,
// /* [out] */ DWORD *pdwRevisedMaxSize) = 0;
func (v *IOPCEventSubscriptionMgt) SetState(active *int32, bufferTime *uint32, maxSize *uint32, clientSubscription uint32) (revisedBufferTime uint32, revisedMaxSize uint32, err error) {
	r0, _, _ := syscall.SyscallN(
		v.VTable().SetState,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(active)),
		uintptr(unsafe.Pointer(bufferTime)),
		uintptr(unsafe.Pointer(maxSize)),
		uintptr(clientSubscription),
		uintptr(unsafe.Pointer(&revisedBufferTime)),
		uintptr(unsafe.Pointer(&revisedMaxSize)),
	)
	if r0 != 0 {
		err = syscall.Errno(r0)
	}
	return
}
