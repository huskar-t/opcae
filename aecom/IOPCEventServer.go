package aecom

import (
	"syscall"
	"time"
	"unsafe"

	"github.com/huskar-t/opcda/com"
	"golang.org/x/sys/windows"
)

// 65168851-5783-11D1-84A0-00608CB8A7E9
var IID_IOPCEventServer = windows.GUID{
	Data1: 0x65168851,
	Data2: 0x5783,
	Data3: 0x11D1,
	Data4: [8]byte{0x84, 0xA0, 0x00, 0x60, 0x8C, 0xB8, 0xA7, 0xE9},
}

type IOPCEventServer struct {
	*com.IUnknown
}

type IOPCEventServerVtbl struct {
	com.IUnknownVtbl
	GetStatus                uintptr
	CreateEventSubscription  uintptr
	QueryAvailableFilters    uintptr
	QueryEventCategories     uintptr
	QueryConditionNames      uintptr
	QuerySubConditionNames   uintptr
	QuerySourceConditions    uintptr
	QueryEventAttributes     uintptr
	TranslateToItemIDs       uintptr
	GetConditionState        uintptr
	EnableConditionByArea    uintptr
	EnableConditionBySource  uintptr
	DisableConditionByArea   uintptr
	DisableConditionBySource uintptr
	AckCondition             uintptr
	CreateAreaBrowser        uintptr
}

func (v *IOPCEventServer) Vtbl() *IOPCEventServerVtbl {
	return (*IOPCEventServerVtbl)(unsafe.Pointer(v.IUnknown.LpVtbl))
}

//typedef /* [public][public] */ struct __MIDL___MIDL_itf_opc_ae_0262_0005
//{
//FILETIME ftStartTime;
//FILETIME ftCurrentTime;
//FILETIME ftLastUpdateTime;
//OPCEVENTSERVERSTATE dwServerState;
//WORD wMajorVersion;
//WORD wMinorVersion;
//WORD wBuildNumber;
//WORD wReserved;
///* [string] */ LPWSTR szVendorInfo;
//} 	OPCEVENTSERVERSTATUS;

type OPCEVENTSERVERSTATUS struct {
	FtStartTime      windows.Filetime
	FtCurrentTime    windows.Filetime
	FtLastUpdateTime windows.Filetime
	DwServerState    int32
	WMajorVersion    uint16
	WMinorVersion    uint16
	WBuildNumber     uint16
	WReserved        uint16
	SzVendorInfo     *uint16
}

type EventServerStatus struct {
	StartTime      time.Time
	CurrentTime    time.Time
	LastUpdateTime time.Time
	ServerState    int32
	MajorVersion   uint16
	MinorVersion   uint16
	BuildNumber    uint16
	Reserved       uint16
	VendorInfo     string
}

func (v *IOPCEventServer) GetStatus() (*EventServerStatus, error) {
	var pStatus *OPCEVENTSERVERSTATUS
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().GetStatus,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(&pStatus)),
	)
	if int32(r0) < 0 {
		return nil, syscall.Errno(r0)
	}
	defer func() {
		if pStatus != nil {
			if pStatus.SzVendorInfo != nil {
				com.CoTaskMemFree(unsafe.Pointer(pStatus.SzVendorInfo))
			}
		}
	}()
	status := &EventServerStatus{
		StartTime:      time.Unix(0, pStatus.FtStartTime.Nanoseconds()),
		CurrentTime:    time.Unix(0, pStatus.FtCurrentTime.Nanoseconds()),
		LastUpdateTime: time.Unix(0, pStatus.FtLastUpdateTime.Nanoseconds()),
		ServerState:    pStatus.DwServerState,
		MajorVersion:   pStatus.WMajorVersion,
		MinorVersion:   pStatus.WMinorVersion,
		BuildNumber:    pStatus.WBuildNumber,
		Reserved:       pStatus.WReserved,
		VendorInfo:     windows.UTF16PtrToString(pStatus.SzVendorInfo),
	}
	return status, nil
}

// HRESULT ( STDMETHODCALLTYPE *CreateEventSubscription )(
// IOPCEventServer2 * This,
// /* [in] */ BOOL bActive,
// /* [in] */ DWORD dwBufferTime,
// /* [in] */ DWORD dwMaxSize,
// /* [in] */ OPCHANDLE hClientSubscription,
// /* [in] */ REFIID riid,
// /* [iid_is][out] */ LPUNKNOWN *ppUnk,
// /* [out] */ DWORD *pdwRevisedBufferTime,
// /* [out] */ DWORD *pdwRevisedMaxSize);
func (v *IOPCEventServer) CreateEventSubscription(active bool, bufferTime, maxSize uint32, clientSubscriptionHandle uint32, riid *windows.GUID) (unknown *com.IUnknown, revisedBufferTime, revisedMaxSize uint32, err error) {
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().CreateEventSubscription,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(com.BoolToComBOOL(active)),
		uintptr(bufferTime),
		uintptr(maxSize),
		uintptr(clientSubscriptionHandle),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(&unknown)),
		uintptr(unsafe.Pointer(&revisedBufferTime)),
		uintptr(unsafe.Pointer(&revisedMaxSize)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

// HRESULT ( STDMETHODCALLTYPE *QueryAvailableFilters )(
// IOPCEventServer * This,
// /* [out] */ DWORD *pdwFilterMask);

func (v *IOPCEventServer) QueryAvailableFilters() (filterMask uint32, err error) {
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().QueryAvailableFilters,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(&filterMask)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	return
}

// HRESULT ( STDMETHODCALLTYPE *QueryEventCategories )(
// IOPCEventServer * This,
// /* [in] */ DWORD dwEventType,
// /* [out] */ DWORD *pdwCount,
// /* [size_is][size_is][out] */ DWORD **ppdwEventCategories,
// /* [size_is][size_is][out] */ LPWSTR **ppszEventCategoryDescs);
var pointerSize uintptr = unsafe.Sizeof(uintptr(0))

func (v *IOPCEventServer) QueryEventCategories(eventType uint32) (eventCategories []uint32, eventCategoryDescs []string, err error) {
	var pdwCount uint32
	var pEventCategories unsafe.Pointer
	var pEventCategoryDescs unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().QueryEventCategories,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(eventType),
		uintptr(unsafe.Pointer(&pdwCount)),
		uintptr(unsafe.Pointer(&pEventCategories)),
		uintptr(unsafe.Pointer(&pEventCategoryDescs)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer func() {
		com.CoTaskMemFree(pEventCategories)
		com.CoTaskMemFree(pEventCategoryDescs)
	}()
	eventCategories = make([]uint32, pdwCount)
	eventCategoryDescs = make([]string, pdwCount)
	for i := 0; i < int(pdwCount); i++ {
		eventCategories[i] = *(*uint32)(unsafe.Pointer(uintptr(pEventCategories) + uintptr(i)*4))
		pwstr := *(**uint16)(unsafe.Pointer(uintptr(pEventCategoryDescs) + uintptr(i)*pointerSize))
		eventCategoryDescs[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
	}
	return
}

//HRESULT ( STDMETHODCALLTYPE *QueryConditionNames )(
//IOPCEventServer * This,
///* [in] */ DWORD dwEventCategory,
///* [out] */ DWORD *pdwCount,
///* [size_is][size_is][out] */ LPWSTR **ppszConditionNames);

func (v *IOPCEventServer) QueryConditionNames(eventCategory uint32) (conditionNames []string, err error) {
	var pdwCount uint32
	var pConditionNames unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().QueryConditionNames,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(eventCategory),
		uintptr(unsafe.Pointer(&pdwCount)),
		uintptr(unsafe.Pointer(&pConditionNames)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer com.CoTaskMemFree(pConditionNames)
	conditionNames = make([]string, pdwCount)
	for i := 0; i < int(pdwCount); i++ {
		pwstr := *(**uint16)(unsafe.Pointer(uintptr(pConditionNames) + uintptr(i)*pointerSize))
		conditionNames[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
	}
	return
}

//virtual HRESULT STDMETHODCALLTYPE QuerySubConditionNames(
///* [in] */ LPWSTR szConditionName,
///* [out] */ DWORD *pdwCount,
///* [size_is][size_is][out] */ LPWSTR **ppszSubConditionNames) = 0;

func (v *IOPCEventServer) QuerySubConditionNames(conditionName string) (subConditionNames []string, err error) {
	var pdwCount uint32
	var pSubConditionNames unsafe.Pointer
	pConditionName, err := syscall.UTF16PtrFromString(conditionName)
	if err != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().QuerySubConditionNames,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(pConditionName)),
		uintptr(unsafe.Pointer(&pdwCount)),
		uintptr(unsafe.Pointer(&pSubConditionNames)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer com.CoTaskMemFree(pSubConditionNames)
	subConditionNames = make([]string, pdwCount)
	for i := 0; i < int(pdwCount); i++ {
		pwstr := *(**uint16)(unsafe.Pointer(uintptr(pSubConditionNames) + uintptr(i)*pointerSize))
		subConditionNames[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
	}
	return
}

//virtual HRESULT STDMETHODCALLTYPE QuerySourceConditions(
///* [in] */ LPWSTR szSource,
///* [out] */ DWORD *pdwCount,
///* [size_is][size_is][out] */ LPWSTR **ppszConditionNames) = 0;

func (v *IOPCEventServer) QuerySourceConditions(source string) (conditionNames []string, err error) {
	var pdwCount uint32
	var pConditionNames unsafe.Pointer
	pSource, err := syscall.UTF16PtrFromString(source)
	if err != nil {
		return
	}
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().QuerySourceConditions,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(pSource)),
		uintptr(unsafe.Pointer(&pdwCount)),
		uintptr(unsafe.Pointer(&pConditionNames)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer com.CoTaskMemFree(pConditionNames)
	conditionNames = make([]string, pdwCount)
	for i := 0; i < int(pdwCount); i++ {
		pwstr := *(**uint16)(unsafe.Pointer(uintptr(pConditionNames) + uintptr(i)*pointerSize))
		conditionNames[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE QueryEventAttributes(
// /* [in] */ DWORD dwEventCategory,
// /* [out] */ DWORD *pdwCount,
// /* [size_is][size_is][out] */ DWORD **ppdwAttrIDs,
// /* [size_is][size_is][out] */ LPWSTR **ppszAttrDescs,
// /* [size_is][size_is][out] */ VARTYPE **ppvtAttrTypes) = 0;
func (v *IOPCEventServer) QueryEventAttributes(eventCategory uint32) (attrIDs []uint32, attrDescs []string, attrTypes []uint16, err error) {
	var pdwCount uint32
	var pAttrIDs unsafe.Pointer
	var pAttrDescs unsafe.Pointer
	var pvtAttrTypes unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().QueryEventAttributes,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(eventCategory),
		uintptr(unsafe.Pointer(&pdwCount)),
		uintptr(unsafe.Pointer(&pAttrIDs)),
		uintptr(unsafe.Pointer(&pAttrDescs)),
		uintptr(unsafe.Pointer(&pvtAttrTypes)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer func() {
		com.CoTaskMemFree(pAttrIDs)
		com.CoTaskMemFree(pAttrDescs)
		com.CoTaskMemFree(pvtAttrTypes)
	}()
	attrIDs = make([]uint32, pdwCount)
	attrDescs = make([]string, pdwCount)
	attrTypes = make([]uint16, pdwCount)
	for i := 0; i < int(pdwCount); i++ {
		attrIDs[i] = *(*uint32)(unsafe.Pointer(uintptr(pAttrIDs) + uintptr(i)*4))
		pwstr := *(**uint16)(unsafe.Pointer(uintptr(pAttrDescs) + uintptr(i)*pointerSize))
		attrDescs[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
		attrTypes[i] = *(*uint16)(unsafe.Pointer(uintptr(pvtAttrTypes) + uintptr(i)*2))
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE TranslateToItemIDs(
// /* [in] */ LPWSTR szSource,
// /* [in] */ DWORD dwEventCategory,
// /* [in] */ LPWSTR szConditionName,
// /* [in] */ LPWSTR szSubconditionName,
// /* [in] */ DWORD dwCount,
// /* [size_is][in] */ DWORD *pdwAssocAttrIDs,
// /* [size_is][size_is][out] */ LPWSTR **ppszAttrItemIDs,
// /* [size_is][size_is][out] */ LPWSTR **ppszNodeNames,
// /* [size_is][size_is][out] */ CLSID **ppCLSIDs) = 0;
func (v *IOPCEventServer) TranslateToItemIDs(source string, eventCategory uint32, conditionName, subconditionName string, assocAttrIDs []uint32) (attrItemIDs, nodeNames []string, clsIDs []windows.GUID, err error) {
	var pSource, pConditionName, pSubconditionName *uint16
	pSource, err = syscall.UTF16PtrFromString(source)
	if err != nil {
		return
	}
	pConditionName, err = syscall.UTF16PtrFromString(conditionName)
	if err != nil {
		return
	}
	pSubconditionName, err = syscall.UTF16PtrFromString(subconditionName)
	if err != nil {
		return
	}
	var pdwAssocAttrIDs unsafe.Pointer
	if len(assocAttrIDs) > 0 {
		pdwAssocAttrIDs = unsafe.Pointer(&assocAttrIDs[0])
	}
	var pdwCount uint32
	var pAttrItemIDs unsafe.Pointer
	var pNodeNames unsafe.Pointer
	var pCLSIDs unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().TranslateToItemIDs,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(pSource)),
		uintptr(eventCategory),
		uintptr(unsafe.Pointer(pConditionName)),
		uintptr(unsafe.Pointer(pSubconditionName)),
		uintptr(len(assocAttrIDs)),
		uintptr(pdwAssocAttrIDs),
		uintptr(unsafe.Pointer(&pAttrItemIDs)),
		uintptr(unsafe.Pointer(&pNodeNames)),
		uintptr(unsafe.Pointer(&pCLSIDs)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer func() {
		com.CoTaskMemFree(pAttrItemIDs)
		com.CoTaskMemFree(pNodeNames)
		com.CoTaskMemFree(pCLSIDs)
	}()
	attrItemIDs = make([]string, pdwCount)
	nodeNames = make([]string, pdwCount)
	clsIDs = make([]windows.GUID, pdwCount)
	for i := 0; i < int(pdwCount); i++ {
		pwstr := *(**uint16)(unsafe.Pointer(uintptr(pAttrItemIDs) + uintptr(i)*pointerSize))
		attrItemIDs[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
		pwstr = *(**uint16)(unsafe.Pointer(uintptr(pNodeNames) + uintptr(i)*pointerSize))
		nodeNames[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
		clsIDs[i] = *(*windows.GUID)(unsafe.Pointer(uintptr(pCLSIDs) + uintptr(i)*unsafe.Sizeof(windows.GUID{})))
	}
	return
}

// typedef /* [public][public] */ struct __MIDL___MIDL_itf_opc_ae_0262_0006
//
//	{
//	WORD wState;
//	WORD wReserved1;
//	LPWSTR szActiveSubCondition;
//	LPWSTR szASCDefinition;
//	DWORD dwASCSeverity;
//	LPWSTR szASCDescription;
//	WORD wQuality;
//	WORD wReserved2;
//	FILETIME ftLastAckTime;
//	FILETIME ftSubCondLastActive;
//	FILETIME ftCondLastActive;
//	FILETIME ftCondLastInactive;
//	LPWSTR szAcknowledgerID;
//	LPWSTR szComment;
//	DWORD dwNumSCs;
//	/* [size_is] */ LPWSTR *pszSCNames;
//	/* [size_is] */ LPWSTR *pszSCDefinitions;
//	/* [size_is] */ DWORD *pdwSCSeverities;
//	/* [size_is] */ LPWSTR *pszSCDescriptions;
//	DWORD dwNumEventAttrs;
//	/* [size_is] */ VARIANT *pEventAttributes;
//	/* [size_is] */ HRESULT *pErrors;
//	} 	OPCCONDITIONSTATE;
type OPCCONDITIONSTATE struct {
	WState               uint16
	WReserved1           uint16
	SzActiveSubCondition *uint16
	SzASCDefinition      *uint16
	DwASCSeverity        uint32
	SzASCDescription     *uint16
	WQuality             uint16
	WReserved2           uint16
	FtLastAckTime        windows.Filetime
	FtSubCondLastActive  windows.Filetime
	FtCondLastActive     windows.Filetime
	FtCondLastInactive   windows.Filetime
	SzAcknowledgerID     *uint16
	SzComment            *uint16
	DwNumSCs             uint32
	PszSCNames           unsafe.Pointer
	PszSCDefinitions     unsafe.Pointer
	PdwSCSeverities      unsafe.Pointer
	PszSCDescriptions    unsafe.Pointer
	DwNumEventAttrs      uint32
	PEventAttributes     unsafe.Pointer
	PErrors              unsafe.Pointer
}

type ConditionState struct {
	State              uint16
	Reserved1          uint16
	ActiveSubCondition string
	ASCDefinition      string
	ASCSeverity        uint32
	ASCDescription     string
	Quality            uint16
	Reserved2          uint16
	LastAckTime        time.Time
	SubCondLastActive  time.Time
	CondLastActive     time.Time
	CondLastInactive   time.Time
	AcknowledgerID     string
	Comment            string
	NumSCs             uint32
	SCNames            []string
	SCDefinitions      []string
	SCSeverities       []uint32
	SCDescriptions     []string
	NumEventAttrs      uint32
	EventAttributes    []com.VARIANT
	Errors             []int32
}

// virtual HRESULT STDMETHODCALLTYPE GetConditionState(
// /* [in] */ LPWSTR szSource,
// /* [in] */ LPWSTR szConditionName,
// /* [in] */ DWORD dwNumEventAttrs,
// /* [size_is][in] */ DWORD *pdwAttributeIDs,
// /* [out] */ OPCCONDITIONSTATE **ppConditionState) = 0;
func (v *IOPCEventServer) GetConditionState(source, conditionName string, attributeIDs []uint32) (conditionState *ConditionState, err error) {
	var pSource, pConditionName *uint16
	pSource, err = syscall.UTF16PtrFromString(source)
	if err != nil {
		return
	}
	pConditionName, err = syscall.UTF16PtrFromString(conditionName)
	if err != nil {
		return
	}
	var pdwAttributeIDs unsafe.Pointer
	if len(attributeIDs) > 0 {
		pdwAttributeIDs = unsafe.Pointer(&attributeIDs[0])
	}
	var pConditionState *OPCCONDITIONSTATE
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().GetConditionState,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(pSource)),
		uintptr(unsafe.Pointer(pConditionName)),
		uintptr(len(attributeIDs)),
		uintptr(pdwAttributeIDs),
		uintptr(unsafe.Pointer(&pConditionState)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer com.CoTaskMemFree(unsafe.Pointer(pConditionState))
	conditionState = &ConditionState{
		State:              pConditionState.WState,
		Reserved1:          pConditionState.WReserved1,
		ActiveSubCondition: windows.UTF16PtrToString(pConditionState.SzActiveSubCondition),
		ASCDefinition:      windows.UTF16PtrToString(pConditionState.SzASCDefinition),
		ASCSeverity:        pConditionState.DwASCSeverity,
		ASCDescription:     windows.UTF16PtrToString(pConditionState.SzASCDescription),
		Quality:            pConditionState.WQuality,
		Reserved2:          pConditionState.WReserved2,
		LastAckTime:        time.Unix(0, pConditionState.FtLastAckTime.Nanoseconds()),
		SubCondLastActive:  time.Unix(0, pConditionState.FtSubCondLastActive.Nanoseconds()),
		CondLastActive:     time.Unix(0, pConditionState.FtCondLastActive.Nanoseconds()),
		CondLastInactive:   time.Unix(0, pConditionState.FtCondLastInactive.Nanoseconds()),
		AcknowledgerID:     windows.UTF16PtrToString(pConditionState.SzAcknowledgerID),
		Comment:            windows.UTF16PtrToString(pConditionState.SzComment),
		NumSCs:             pConditionState.DwNumSCs,
		NumEventAttrs:      pConditionState.DwNumEventAttrs,
	}
	conditionState.SCNames = make([]string, conditionState.NumSCs)
	conditionState.SCDefinitions = make([]string, conditionState.NumSCs)
	conditionState.SCSeverities = make([]uint32, conditionState.NumSCs)
	conditionState.SCDescriptions = make([]string, conditionState.NumSCs)
	for i := 0; i < int(conditionState.NumSCs); i++ {
		pwstr := *(**uint16)(unsafe.Pointer(uintptr(pConditionState.PszSCNames) + uintptr(i)*pointerSize))
		conditionState.SCNames[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
		pwstr = *(**uint16)(unsafe.Pointer(uintptr(pConditionState.PszSCDefinitions) + uintptr(i)*pointerSize))
		conditionState.SCDefinitions[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
		conditionState.SCSeverities[i] = *(*uint32)(unsafe.Pointer(uintptr(pConditionState.PdwSCSeverities) + uintptr(i)*4))
		pwstr = *(**uint16)(unsafe.Pointer(uintptr(pConditionState.PszSCDescriptions) + uintptr(i)*pointerSize))
		conditionState.SCDescriptions[i] = windows.UTF16PtrToString(pwstr)
		com.CoTaskMemFree(unsafe.Pointer(pwstr))
	}
	conditionState.EventAttributes = make([]com.VARIANT, conditionState.NumEventAttrs)
	conditionState.Errors = make([]int32, conditionState.NumEventAttrs)
	for i := 0; i < int(conditionState.NumEventAttrs); i++ {
		variant := *(*com.VARIANT)(unsafe.Pointer(uintptr(pConditionState.PEventAttributes) + uintptr(i)*unsafe.Sizeof(com.VARIANT{})))
		errNo := *(*int32)(unsafe.Pointer(uintptr(pConditionState.PErrors) + uintptr(i)*4))
		conditionState.EventAttributes[i] = variant
		conditionState.Errors[i] = int32(errNo)
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE EnableConditionByArea(
// /* [in] */ DWORD dwNumAreas,
// /* [size_is][in] */ LPWSTR *pszAreas) = 0;
func (v *IOPCEventServer) EnableConditionByArea(areas []string) (err error) {
	var pdwNumAreas uint32
	var ppszAreas unsafe.Pointer
	if len(areas) > 0 {
		pdwNumAreas = uint32(len(areas))
		ppszAreas = unsafe.Pointer(&areas[0])
	}
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().EnableConditionByArea,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(pdwNumAreas),
		uintptr(ppszAreas),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE EnableConditionBySource(
// /* [in] */ DWORD dwNumSources,
// /* [size_is][in] */ LPWSTR *pszSources) = 0;
func (v *IOPCEventServer) EnableConditionBySource(sources []string) (err error) {
	var pdwNumSources uint32
	var ppszSources unsafe.Pointer
	if len(sources) > 0 {
		pdwNumSources = uint32(len(sources))
		ppszSources = unsafe.Pointer(&sources[0])
	}
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().EnableConditionBySource,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(pdwNumSources),
		uintptr(ppszSources),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE DisableConditionByArea(
// /* [in] */ DWORD dwNumAreas,
// /* [size_is][in] */ LPWSTR *pszAreas) = 0;
func (v *IOPCEventServer) DisableConditionByArea(areas []string) (err error) {
	var pdwNumAreas uint32
	var ppszAreas unsafe.Pointer
	if len(areas) > 0 {
		pdwNumAreas = uint32(len(areas))
		ppszAreas = unsafe.Pointer(&areas[0])
	}
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().DisableConditionByArea,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(pdwNumAreas),
		uintptr(ppszAreas),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE DisableConditionBySource(
// /* [in] */ DWORD dwNumSources,
// /* [size_is][in] */ LPWSTR *pszSources) = 0;
func (v *IOPCEventServer) DisableConditionBySource(sources []string) (err error) {
	var pdwNumSources uint32
	var ppszSources unsafe.Pointer
	if len(sources) > 0 {
		pdwNumSources = uint32(len(sources))
		ppszSources = unsafe.Pointer(&sources[0])
	}
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().DisableConditionBySource,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(pdwNumSources),
		uintptr(ppszSources),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE AckCondition(
// /* [in] */ DWORD dwCount,
// /* [string][in] */ LPWSTR szAcknowledgerID,
// /* [string][in] */ LPWSTR szComment,
// /* [size_is][in] */ LPWSTR *pszSource,
// /* [size_is][in] */ LPWSTR *pszConditionName,
// /* [size_is][in] */ FILETIME *pftActiveTime,
// /* [size_is][in] */ DWORD *pdwCookie,
// /* [size_is][size_is][out] */ HRESULT **ppErrors) = 0;
func (v *IOPCEventServer) AckCondition(acknowledgerID, comment string, sources, conditionNames []string, activeTimes []time.Time, cookies []uint32) (errors []int32, err error) {
	var pdwCount uint32
	var pszSource, pszConditionName, pftActiveTime, pdwCookie unsafe.Pointer
	if len(sources) > 0 {
		pdwCount = uint32(len(sources))
		pszSource = unsafe.Pointer(&sources[0])
		pszConditionName = unsafe.Pointer(&conditionNames[0])
		pftActiveTime = unsafe.Pointer(&activeTimes[0])
		pdwCookie = unsafe.Pointer(&cookies[0])
	}
	var ppErrors unsafe.Pointer
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().AckCondition,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(pdwCount),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(acknowledgerID))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(comment))),
		uintptr(pszSource),
		uintptr(pszConditionName),
		uintptr(pftActiveTime),
		uintptr(pdwCookie),
		uintptr(unsafe.Pointer(&ppErrors)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
		return
	}
	defer com.CoTaskMemFree(ppErrors)
	errors = make([]int32, pdwCount)
	for i := 0; i < int(pdwCount); i++ {
		errNo := *(*int32)(unsafe.Pointer(uintptr(ppErrors) + uintptr(i)*4))
		if errNo < 0 {
			errors[i] = int32(errNo)
		}
	}
	return
}

// virtual HRESULT STDMETHODCALLTYPE CreateAreaBrowser(
// /* [in] */ REFIID riid,
// /* [iid_is][out] */ LPUNKNOWN *ppUnk) = 0;
func (v *IOPCEventServer) CreateAreaBrowser(riid *windows.GUID) (unknown *com.IUnknown, err error) {
	r0, _, _ := syscall.SyscallN(
		v.Vtbl().CreateAreaBrowser,
		uintptr(unsafe.Pointer(v.IUnknown)),
		uintptr(unsafe.Pointer(riid)),
		uintptr(unsafe.Pointer(&unknown)),
	)
	if int32(r0) < 0 {
		err = syscall.Errno(r0)
	}
	return
}
