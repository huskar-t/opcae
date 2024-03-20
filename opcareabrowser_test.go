package opcae

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAreaBrowser(t *testing.T) {
	eventServer, err := ConnectEventServer(KepWareProgID, TestHost)
	if err != nil {
		t.Fatalf("connect to opc event server failed: %s\n", err)
	}
	assert.NotNil(t, eventServer)
	defer eventServer.Disconnect()
	browser, err := eventServer.CreateAreaBrowser()
	assert.NoError(t, err)
	assert.NotNil(t, browser)
	defer browser.Release()
	err = browser.MoveToRoot()
	assert.NoError(t, err)
	areas, err := browser.BrowseOPCAreas(OPC_AREA, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, areas)
	expectAreas := []string{"_CustomAlarms", "_System"}
	assert.Equal(t, expectAreas, areas)
	t.Log(areas)
	for _, area := range areas {
		qualifiedArea, err := browser.GetQualifiedAreaName(area)
		assert.NoError(t, err)
		assert.Equal(t, area, qualifiedArea)
	}
	err = browser.MoveDown("_System")
	assert.NoError(t, err)
	areas, err = browser.BrowseOPCAreas(OPC_AREA, "")
	assert.NoError(t, err)
	assert.Empty(t, areas)
	err = browser.MoveUP()
	assert.NoError(t, err)
	areas, err = browser.BrowseOPCAreas(OPC_AREA, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, areas)
	assert.Equal(t, expectAreas, areas)
}
