package opcae

type Filter uint32

const (
	OPC_FILTER_BY_EVENT    Filter = 0x1
	OPC_FILTER_BY_CATEGORY Filter = 0x2
	OPC_FILTER_BY_SEVERITY Filter = 0x4
	OPC_FILTER_BY_AREA     Filter = 0x8
	OPC_FILTER_BY_SOURCE   Filter = 0x10
)

var filterList = []Filter{
	OPC_FILTER_BY_EVENT,
	OPC_FILTER_BY_CATEGORY,
	OPC_FILTER_BY_SEVERITY,
	OPC_FILTER_BY_AREA,
	OPC_FILTER_BY_SOURCE,
}

func ParseFilter(filter uint32) (filters []Filter) {
	for _, f := range filterList {
		if filter&uint32(f) != 0 {
			filters = append(filters, f)
		}
	}
	return
}

func MarshalFilter(filters []Filter) (filter uint32) {
	for _, f := range filters {
		filter |= uint32(f)
	}
	return
}

type BrowseType uint32

const (
	OPC_AREA   BrowseType = 0x1
	OPC_SOURCE BrowseType = 0x2
)

type EventCategoryType uint32

const (
	OPC_SIMPLE_EVENT    EventCategoryType = 0x1
	OPC_TRACKING_EVENT  EventCategoryType = 0x2
	OPC_CONDITION_EVENT EventCategoryType = 0x4
	OPC_ALL_EVENTS      EventCategoryType = 0x7
)

func MarshalEventCategoryType(categories []EventCategoryType) (category uint32) {
	for _, c := range categories {
		category |= uint32(c)
	}
	return
}

func UnmarshalEventCategoryType(category uint32) (categories []EventCategoryType) {
	for _, c := range []EventCategoryType{
		OPC_SIMPLE_EVENT,
		OPC_TRACKING_EVENT,
		OPC_CONDITION_EVENT,
		OPC_ALL_EVENTS,
	} {
		if category&uint32(c) != 0 {
			categories = append(categories, c)
		}
	}
	return
}

type ChangeMask uint16

const (
	OPC_CHANGE_ACTIVE_STATE ChangeMask = 0x1
	OPC_CHANGE_ACK_STATE    ChangeMask = 0x2
	OPC_CHANGE_ENABLE_STATE ChangeMask = 0x4
	OPC_CHANGE_QUALITY      ChangeMask = 0x8
	OPC_CHANGE_SEVERITY     ChangeMask = 0x10
	OPC_CHANGE_SUBCONDITION ChangeMask = 0x20
	OPC_CHANGE_MESSAGE      ChangeMask = 0x40
	OPC_CHANGE_ATTRIBUTE    ChangeMask = 0x80
)

var changeMaskList = []ChangeMask{
	OPC_CHANGE_ACTIVE_STATE,
	OPC_CHANGE_ACK_STATE,
	OPC_CHANGE_ENABLE_STATE,
	OPC_CHANGE_QUALITY,
	OPC_CHANGE_SEVERITY,
	OPC_CHANGE_SUBCONDITION,
	OPC_CHANGE_MESSAGE,
	OPC_CHANGE_ATTRIBUTE,
}

func ParseChangeMask(mask uint16) (masks []ChangeMask) {
	for _, m := range changeMaskList {
		if mask&uint16(m) != 0 {
			masks = append(masks, m)
		}
	}
	return
}

type State uint16

const (
	OPC_CONDITION_ENABLED State = 0x1
	OPC_CONDITION_ACTIVE  State = 0x2
	OPC_CONDITION_ACKED   State = 0x4
)
