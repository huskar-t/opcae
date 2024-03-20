package opcae

import (
	"syscall"
	"time"
	"unsafe"

	"github.com/huskar-t/opcda/com"
	"golang.org/x/sys/windows"
)

var IID_IOPCEventSink = windows.GUID{
	Data1: 0x6516885F,
	Data2: 0x5783,
	Data3: 0x11D1,
	Data4: [8]byte{0x84, 0xA0, 0x00, 0x60, 0x8C, 0xB8, 0xA7, 0xE9},
}

type IOPCEventSink struct {
	lpVtbl   *IOPCEventSinkVtbl
	ref      int32
	clsid    *windows.GUID
	receiver chan *EventSinkOnEventData
}

type IOPCEventSinkVtbl struct {
	pQueryInterface uintptr
	pAddRef         uintptr
	pRelease        uintptr
	pOnEvent        uintptr
}

func NewEventSink(
	receiver chan *EventSinkOnEventData,
) *IOPCEventSink {
	return &IOPCEventSink{
		lpVtbl: &IOPCEventSinkVtbl{
			pQueryInterface: syscall.NewCallback(EventSinkQueryInterface),
			pAddRef:         syscall.NewCallback(EventSinkAddRef),
			pRelease:        syscall.NewCallback(EventSinkRelease),
			pOnEvent:        syscall.NewCallback(EventSinkOnEvent),
		},
		ref:      0,
		clsid:    &IID_IOPCEventSink,
		receiver: receiver,
	}
}

func EventSinkQueryInterface(this unsafe.Pointer, iid *windows.GUID, punk *unsafe.Pointer) uintptr {
	er := (*IOPCEventSink)(this)
	*punk = nil
	if com.IsEqualGUID(iid, er.clsid) || com.IsEqualGUID(iid, com.IID_IUnknown) {
		EventSinkAddRef(this)
		*punk = this
		return com.S_OK
	}
	return com.E_POINTER
}

func EventSinkAddRef(this unsafe.Pointer) uintptr {
	er := (*IOPCEventSink)(this)
	er.ref++
	return uintptr(er.ref)
}

func EventSinkRelease(this unsafe.Pointer) uintptr {
	er := (*IOPCEventSink)(this)
	er.ref--
	return uintptr(er.ref)
}

// typedef /* [public][public] */ struct __MIDL___MIDL_itf_opc_ae_0262_0004
//
//	{
//	WORD wChangeMask;
//	WORD wNewState;
//	/* [string] */ LPWSTR szSource;
//	FILETIME ftTime;
//	/* [string] */ LPWSTR szMessage;
//	DWORD dwEventType;
//	DWORD dwEventCategory;
//	DWORD dwSeverity;
//	/* [string] */ LPWSTR szConditionName;
//	/* [string] */ LPWSTR szSubconditionName;
//	WORD wQuality;
//	WORD wReserved;
//	BOOL bAckRequired;
//	FILETIME ftActiveTime;
//	DWORD dwCookie;
//	DWORD dwNumEventAttrs;
//	/* [size_is] */ VARIANT *pEventAttributes;
//	/* [string] */ LPWSTR szActorID;
//	} 	ONEVENTSTRUCT;
type ONEVENTSTRUCT struct {
	WChangeMask        uint16
	WNewState          uint16
	SzSource           *uint16
	FtTime             windows.Filetime
	SzMessage          *uint16
	DwEventType        uint32
	DwEventCategory    uint32
	DwSeverity         uint32
	SzConditionName    *uint16
	SzSubconditionName *uint16
	WQuality           uint16
	WReserved          uint16
	BAckRequired       int32
	FtActiveTime       windows.Filetime
	DwCookie           uint32
	DwNumEventAttrs    uint32
	PEventAttributes   unsafe.Pointer
	SzActorID          *uint16
}

type EventSinkOnEventData struct {
	ClientHandle uint32
	Refresh      bool
	LastRefresh  bool
	Events       []*OnEventStruct
}
type OnEventStruct struct {
	ChangeMask []ChangeMask
	NewState   State
	Source     string
	Time       time.Time
	Message    string
	EventType  uint32
	Category   uint32
	Severity   uint32
	Condition  string
	Subcond    string
	Quality    uint16
	Reserved   uint16
	AckReq     bool
	ActiveTime time.Time
	Cookie     uint32
	NumAttrs   uint32
	Attributes []interface{}
	ActorID    string
}

const VariantSize = unsafe.Sizeof(com.VARIANT{})

// virtual HRESULT STDMETHODCALLTYPE OnEvent(
// /* [in] */ OPCHANDLE hClientSubscription,
// /* [in] */ BOOL bRefresh,
// /* [in] */ BOOL bLastRefresh,
// /* [in] */ DWORD dwCount,
// /* [size_is][in] */ ONEVENTSTRUCT *pEvents) = 0;
func EventSinkOnEvent(this *com.IUnknown, clientHandle uint32, refresh int32, lastRefresh int32, count uint32, events unsafe.Pointer) uintptr {
	er := (*IOPCEventSink)(unsafe.Pointer(this))
	evt := &EventSinkOnEventData{
		ClientHandle: clientHandle,
		Refresh:      refresh != 0,
		LastRefresh:  lastRefresh != 0,
	}
	evt.Events = make([]*OnEventStruct, count)
	for i := uint32(0); i < count; i++ {
		e := (*ONEVENTSTRUCT)(unsafe.Pointer(uintptr(events) + uintptr(i)*unsafe.Sizeof(ONEVENTSTRUCT{})))
		evt.Events[i] = &OnEventStruct{
			ChangeMask: ParseChangeMask(e.WChangeMask),
			NewState:   State(e.WNewState),
			Source:     windows.UTF16PtrToString(e.SzSource),
			Time:       time.Unix(0, e.FtTime.Nanoseconds()),
			Message:    windows.UTF16PtrToString(e.SzMessage),
			EventType:  e.DwEventType,
			Category:   e.DwEventCategory,
			Severity:   e.DwSeverity,
			Condition:  windows.UTF16PtrToString(e.SzConditionName),
			Subcond:    windows.UTF16PtrToString(e.SzSubconditionName),
			Quality:    e.WQuality,
			Reserved:   e.WReserved,
			AckReq:     e.BAckRequired != 0,
			ActiveTime: time.Unix(0, e.FtActiveTime.Nanoseconds()),
			Cookie:     e.DwCookie,
			NumAttrs:   e.DwNumEventAttrs,
			ActorID:    windows.UTF16PtrToString(e.SzActorID),
		}
		evt.Events[i].Attributes = make([]interface{}, e.DwNumEventAttrs)
		for j := uint32(0); j < e.DwNumEventAttrs; j++ {
			variant := (*com.VARIANT)(unsafe.Pointer(uintptr(e.PEventAttributes) + uintptr(j)*VariantSize))
			evt.Events[i].Attributes[j] = variant.Value()
		}
	}
	er.receiver <- evt
	return uintptr(com.S_OK)
}
