package opcae

import (
	"github.com/huskar-t/opcae/aecom"
	"sync/atomic"
	"unsafe"

	"github.com/huskar-t/opcda/com"
	"golang.org/x/sys/windows"
)

type OPCEventServer struct {
	iServer                  *aecom.IOPCEventServer
	iCommon                  *com.IOPCCommon
	Name                     string
	Node                     string
	location                 com.CLSCTX
	clientSubscriptionHandle uint32
	eventSubscriptions       []*OPCEventSubscription
	browsers                 []*OPCAreaBrowser
}

func ConnectEventServer(progID, node string) (eventServer *OPCEventServer, err error) {
	location := com.CLSCTX_LOCAL_SERVER
	if !com.IsLocal(node) {
		location = com.CLSCTX_REMOTE_SERVER
	}
	clsid, err := windows.GUIDFromString(progID)
	if err != nil {
		return nil, err
	}
	iUnknownServer, err := com.MakeCOMObjectEx(node, location, &clsid, &aecom.IID_IOPCEventServer)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			iUnknownServer.Release()
		}
	}()
	var iUnknownCommon *com.IUnknown
	err = iUnknownServer.QueryInterface(&com.IID_IOPCCommon, unsafe.Pointer(&iUnknownCommon))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			iUnknownCommon.Release()
		}
	}()
	server := &aecom.IOPCEventServer{IUnknown: iUnknownServer}
	common := &com.IOPCCommon{IUnknown: iUnknownCommon}
	eventServer = &OPCEventServer{
		iServer:  server,
		iCommon:  common,
		Name:     progID,
		Node:     node,
		location: location,
	}
	return eventServer, nil
}

func (v *OPCEventServer) GetStatus() (*aecom.EventServerStatus, error) {
	return v.iServer.GetStatus()
}

// CreateEventSubscription
// active: FALSE if the Event Subscription is to be created inactive and TRUE if it is to be created as active.
// bufferTime: The requested buffer time. The buffer time is in milliseconds and tells the server how often to send event notifications. A value of 0 for dwBufferTime means that the server should send event notifications as soon as it gets them.
// maxSize: The requested maximum number of events that will be sent in a single IOPCEventSink::OnEvent callback. A value of 0 means that there is no limit to the number of events that will be sent in a single callback
func (v *OPCEventServer) CreateEventSubscription(active bool, bufferTime, maxSize, receiverBufSize uint32) (*OPCEventSubscription, uint32, uint32, error) {
	clientSubscriptionHandle := atomic.AddUint32(&v.clientSubscriptionHandle, 1)
	unknown, revisedBufferTime, revisedMaxSize, err := v.iServer.CreateEventSubscription(active, bufferTime, maxSize, clientSubscriptionHandle, &aecom.IID_IOPCEventSubscriptionMgt)
	if err != nil {
		return nil, 0, 0, err
	}
	sub, err := NewOPCEventSubscription(unknown, v.iCommon, clientSubscriptionHandle, receiverBufSize)
	if err != nil {
		return nil, 0, 0, err
	}
	return sub, revisedBufferTime, revisedMaxSize, nil
}

func (v *OPCEventServer) QueryAvailableFilters() ([]Filter, error) {
	filterMask, err := v.iServer.QueryAvailableFilters()
	if err != nil {
		return nil, err
	}
	return ParseFilter(filterMask), nil
}

type EventCategory struct {
	ID          uint32
	Description string
}

func (v *OPCEventServer) QueryEventCategories(categories []EventCategoryType) ([]*EventCategory, error) {
	category := MarshalEventCategoryType(categories)
	ids, descs, err := v.iServer.QueryEventCategories(category)
	if err != nil {
		return nil, err
	}
	result := make([]*EventCategory, len(ids))
	for i := range ids {
		result[i] = &EventCategory{
			ID:          ids[i],
			Description: descs[i],
		}
	}
	return result, nil
}

func (v *OPCEventServer) QueryConditionNames(categories []EventCategoryType) ([]string, error) {
	category := MarshalEventCategoryType(categories)
	return v.iServer.QueryConditionNames(category)
}

func (v *OPCEventServer) QuerySourceConditions(source string) ([]string, error) {
	return v.iServer.QuerySourceConditions(source)
}

func (v *OPCEventServer) QuerySubConditionNames(conditionName string) ([]string, error) {
	return v.iServer.QuerySubConditionNames(conditionName)
}

type EventAttribute struct {
	ID          uint32
	Description string
	Type        uint16
}

func (v *OPCEventServer) QueryEventAttributes(eventCategoryID uint32) ([]*EventAttribute, error) {
	ids, descs, types, err := v.iServer.QueryEventAttributes(eventCategoryID)
	if err != nil {
		return nil, err
	}
	result := make([]*EventAttribute, len(ids))
	for i := range ids {
		result[i] = &EventAttribute{
			ID:          ids[i],
			Description: descs[i],
			Type:        types[i],
		}
	}
	return result, nil
}

type ItemID struct {
	ID    string
	Name  string
	CLSID windows.GUID
}

func (v *OPCEventServer) TranslateToItemIDs(source string, eventCategoryID uint32, conditionName string, subConditionName string, assocAttrIDs []uint32) ([]*ItemID, error) {
	ids, names, clsIDs, err := v.iServer.TranslateToItemIDs(source, eventCategoryID, conditionName, subConditionName, assocAttrIDs)
	if err != nil {
		return nil, err
	}
	result := make([]*ItemID, len(ids))
	for i := range ids {
		result[i] = &ItemID{
			ID:    ids[i],
			Name:  names[i],
			CLSID: clsIDs[i],
		}
	}
	return result, nil
}

func (v *OPCEventServer) CreateAreaBrowser() (*OPCAreaBrowser, error) {
	unknown, err := v.iServer.CreateAreaBrowser(&aecom.IID_IOPCEventAreaBrowser)
	if err != nil {
		return nil, err
	}
	return NewOPCAreaBrowser(unknown), nil
}

func (v *OPCEventServer) Disconnect() error {
	for _, subscription := range v.eventSubscriptions {
		subscription.Release()
	}
	v.iServer.Release()
	return nil
}
