package opcae

import (
	"github.com/huskar-t/opcae/aecom"
	"unsafe"

	"github.com/huskar-t/opcda/com"
)

type OPCEventSubscription struct {
	cookie               uint32
	receiver             chan *EventSinkOnEventData
	container            *com.IConnectionPointContainer
	point                *com.IConnectionPoint
	event                *IOPCEventSink
	eventSubscriptionMgt *aecom.IOPCEventSubscriptionMgt
	common               *com.IOPCCommon
	clientHandle         uint32
}

func NewOPCEventSubscription(unknown *com.IUnknown, common *com.IOPCCommon, clientHandle, receiverBufSize uint32) (*OPCEventSubscription, error) {
	var iUnknownContainer *com.IUnknown
	err := unknown.QueryInterface(&com.IID_IConnectionPointContainer, unsafe.Pointer(&iUnknownContainer))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			iUnknownContainer.Release()
		}
	}()
	container := &com.IConnectionPointContainer{IUnknown: iUnknownContainer}
	point, err := container.FindConnectionPoint(&IID_IOPCEventSink)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			point.Release()
		}
	}()
	receiver := make(chan *EventSinkOnEventData, receiverBufSize)
	event := NewEventSink(receiver)
	cookie, err := point.Advise((*com.IUnknown)(unsafe.Pointer(event)))
	if err != nil {
		return nil, err
	}
	return &OPCEventSubscription{
		eventSubscriptionMgt: &aecom.IOPCEventSubscriptionMgt{IUnknown: unknown},
		common:               common,
		clientHandle:         clientHandle,
		container:            container,
		point:                point,
		event:                event,
		cookie:               cookie,
		receiver:             receiver,
	}, nil
}

func (es *OPCEventSubscription) GetClientHandle() uint32 {
	return es.clientHandle
}

func (es *OPCEventSubscription) GetState() (active bool, bufferTime uint32, maxSize uint32, clientSubscription uint32, err error) {
	return es.eventSubscriptionMgt.GetState()
}

func (es *OPCEventSubscription) SetActive(active bool) error {
	comBool := com.BoolToComBOOL(active)
	_, _, err := es.eventSubscriptionMgt.SetState(&comBool, nil, nil, es.clientHandle)
	return err
}

func (es *OPCEventSubscription) SetBufferTime(bufferTime uint32) (uint32, error) {
	revisedBufferTime, _, err := es.eventSubscriptionMgt.SetState(nil, &bufferTime, nil, es.clientHandle)
	return revisedBufferTime, err
}

func (es *OPCEventSubscription) SetMaxSize(maxSize uint32) (uint32, error) {
	_, revisedMaxSize, err := es.eventSubscriptionMgt.SetState(nil, nil, &maxSize, es.clientHandle)
	return revisedMaxSize, err
}

func (es *OPCEventSubscription) SetFilter(events []EventCategoryType, eventCategories []uint32, lowSeverity uint32, highSeverity uint32, areaList []string, sourceList []string) error {
	return es.eventSubscriptionMgt.SetFilter(MarshalEventCategoryType(events), eventCategories, lowSeverity, highSeverity, areaList, sourceList)
}

func (es *OPCEventSubscription) GetFilter() (events []EventCategoryType, eventCategories []uint32, lowSeverity uint32, highSeverity uint32, areaList []string, sourceList []string, err error) {
	cEvents, eventCategories, lowSeverity, highSeverity, areaList, sourceList, err := es.eventSubscriptionMgt.GetFilter()
	return UnmarshalEventCategoryType(cEvents), eventCategories, lowSeverity, highSeverity, areaList, sourceList, err
}

func (es *OPCEventSubscription) SelectReturnedAttributes(eventCategory uint32, attributeIDs []uint32) (err error) {
	return es.eventSubscriptionMgt.SelectReturnedAttributes(eventCategory, attributeIDs)
}

func (es *OPCEventSubscription) GetReturnedAttributes(eventCategory uint32) (attributeIDs []uint32, err error) {
	return es.eventSubscriptionMgt.GetReturnedAttributes(eventCategory)
}

func (es *OPCEventSubscription) Refresh() error {
	return es.eventSubscriptionMgt.Refresh(es.cookie)
}

func (es *OPCEventSubscription) CancelRefresh() error {
	return es.eventSubscriptionMgt.CancelRefresh(es.cookie)
}

func (es *OPCEventSubscription) GetReceiver() <-chan *EventSinkOnEventData {
	return es.receiver
}

func (es *OPCEventSubscription) Release() error {
	err := es.point.Unadvise(es.cookie)
	es.point.Release()
	es.container.Release()
	return err
}
