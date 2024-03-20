package opcae

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOPCEventSubscription_GetFilter(t *testing.T) {
	eventServer, err := ConnectEventServer(TestProgID, TestHost)
	if err != nil {
		t.Fatalf("connect to opc event server failed: %s\n", err)
	}
	assert.NotNil(t, eventServer)
	defer eventServer.Disconnect()
	subscription, revisedBufferTime, revisedMaxSize, err := eventServer.CreateEventSubscription(true, 0, 0, 100)
	assert.NoError(t, err)
	defer subscription.Release()
	assert.NotNil(t, subscription)
	assert.Equal(t, uint32(0), revisedBufferTime)
	assert.Equal(t, uint32(0), revisedMaxSize)
	events, categories, lowSeverity, highSeverity, areaList, sourceList, err := subscription.GetFilter()
	assert.NoError(t, err)
	assert.Equal(t, []EventCategoryType{OPC_SIMPLE_EVENT, OPC_TRACKING_EVENT, OPC_CONDITION_EVENT, OPC_ALL_EVENTS}, events)
	assert.Empty(t, categories)
	assert.Equal(t, uint32(1), lowSeverity)
	assert.Equal(t, uint32(1000), highSeverity)
	assert.Empty(t, areaList)
	assert.Empty(t, sourceList)
}

func TestSubscription(t *testing.T) {
	eventServer, err := ConnectEventServer(TestProgID, TestHost)
	if err != nil {
		t.Fatalf("connect to opc event server failed: %s\n", err)
	}
	assert.NotNil(t, eventServer)
	defer eventServer.Disconnect()
	subscription, revisedBufferTime, revisedMaxSize, err := eventServer.CreateEventSubscription(true, 0, 0, 100)
	assert.NoError(t, err)
	defer subscription.Release()
	assert.NotNil(t, subscription)
	assert.Equal(t, uint32(0), revisedBufferTime)
	assert.Equal(t, uint32(0), revisedMaxSize)
	receiver := subscription.GetReceiver()
	data := <-receiver
	for _, event := range data.Events {
		t.Logf("%#v", event)
	}
}
