package opcae

import (
	"testing"

	"github.com/huskar-t/opcda/com"
	"github.com/stretchr/testify/assert"
)

const TestProgID = "Matrikon.OPC.Simulation.1"
const TestHost = "localhost"
const KepWareProgID = "Kepware.KEPServerEX_AE.V6"

func TestMain(m *testing.M) {
	com.Initialize()
	defer com.Uninitialize()
	m.Run()
}
func TestConnectEventServer(t *testing.T) {
	eventServer, err := ConnectEventServer(TestProgID, TestHost)
	if err != nil {
		t.Fatalf("connect to opc event server failed: %s\n", err)
	}
	assert.NotNil(t, eventServer)
	eventServer.Disconnect()
}

func TestQueryAvailableFilters(t *testing.T) {
	eventServer, err := ConnectEventServer(TestProgID, TestHost)
	if err != nil {
		t.Fatalf("connect to opc event server failed: %s\n", err)
	}
	assert.NotNil(t, eventServer)
	defer eventServer.Disconnect()
	filters, err := eventServer.QueryAvailableFilters()
	assert.NoError(t, err)
	t.Log(filters)
}

func TestQueryEventCategories(t *testing.T) {
	eventServer, err := ConnectEventServer(TestProgID, TestHost)
	if err != nil {
		t.Fatalf("connect to opc event server failed: %s\n", err)
	}
	assert.NotNil(t, eventServer)
	defer eventServer.Disconnect()
	categories, err := eventServer.QueryEventCategories([]EventCategoryType{OPC_ALL_EVENTS})
	assert.NoError(t, err)
	assert.NotEmpty(t, categories)
	for _, category := range categories {
		t.Log(category)
	}
}

func TestQueryConditionNames(t *testing.T) {
	eventServer, err := ConnectEventServer(TestProgID, TestHost)
	if err != nil {
		t.Fatalf("connect to opc event server failed: %s\n", err)
	}
	assert.NotNil(t, eventServer)
	defer eventServer.Disconnect()
	names, err := eventServer.QueryConditionNames([]EventCategoryType{OPC_ALL_EVENTS})
	assert.NoError(t, err)
	assert.NotEmpty(t, names)
	for _, name := range names {
		t.Log(name)
	}
}

func TestQueryEventAttributes(t *testing.T) {
	eventServer, err := ConnectEventServer(TestProgID, TestHost)
	if err != nil {
		t.Fatalf("connect to opc event server failed: %s\n", err)
	}
	assert.NotNil(t, eventServer)
	defer eventServer.Disconnect()
	categories, err := eventServer.QueryEventCategories([]EventCategoryType{OPC_ALL_EVENTS})
	assert.NoError(t, err)
	assert.NotEmpty(t, categories)
	for _, category := range categories {
		t.Log(category)
	}
	attributes, err := eventServer.QueryEventAttributes(categories[2].ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, attributes)
	for _, attribute := range attributes {
		t.Log(attribute)
	}
}

func TestCreateAreaBrowser(t *testing.T) {
	eventServer, err := ConnectEventServer(TestProgID, TestHost)
	if err != nil {
		t.Fatalf("connect to opc event server failed: %s\n", err)
	}
	assert.NotNil(t, eventServer)
	defer eventServer.Disconnect()
	browser, err := eventServer.CreateAreaBrowser()
	assert.NoError(t, err)
	assert.NotNil(t, browser)
	defer browser.Release()
}
func TestOPCEventServer_CreateEventSubscription(t *testing.T) {
	eventServer, err := ConnectEventServer(TestProgID, TestHost)
	if err != nil {
		t.Fatalf("connect to opc event server failed: %s\n", err)
	}
	assert.NotNil(t, eventServer)
	defer eventServer.Disconnect()
	subscription, revisedBufferTime, revisedMaxSize, err := eventServer.CreateEventSubscription(true, 100, 1000, 100)
	assert.NoError(t, err)
	defer subscription.Release()
	assert.NotNil(t, subscription)
	assert.Equal(t, uint32(100), revisedBufferTime)
	assert.Equal(t, uint32(1000), revisedMaxSize)
}
