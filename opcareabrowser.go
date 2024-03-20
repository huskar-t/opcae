package opcae

import (
	"github.com/huskar-t/opcae/aecom"

	"github.com/huskar-t/opcda/com"
)

type OPCAreaBrowser struct {
	browser *aecom.IOPCEventAreaBrowser
}

func NewOPCAreaBrowser(unknown *com.IUnknown) *OPCAreaBrowser {
	return &OPCAreaBrowser{
		browser: &aecom.IOPCEventAreaBrowser{IUnknown: unknown},
	}
}

func (b *OPCAreaBrowser) MoveToRoot() error {
	return b.browser.ChangeBrowsePosition(OPCAE_BROWSE_TO, "")
}

func (b *OPCAreaBrowser) MoveUP() error {
	return b.browser.ChangeBrowsePosition(OPCAE_BROWSE_UP, "")
}

func (b *OPCAreaBrowser) MoveDown(area string) error {
	return b.browser.ChangeBrowsePosition(OPCAE_BROWSE_DOWN, area)
}

func (b *OPCAreaBrowser) BrowseOPCAreas(browseFilterType BrowseType, filterCriteria string) (areas []string, err error) {
	return b.browser.BrowseOPCAreas(uint32(browseFilterType), filterCriteria)
}

func (b *OPCAreaBrowser) GetQualifiedAreaName(areaName string) (qualifiedAreaName string, err error) {
	return b.browser.GetQualifiedAreaName(areaName)
}

func (b *OPCAreaBrowser) GetQualifiedSourceName(sourceName string) (qualifiedSourceName string, err error) {
	return b.browser.GetQualifiedSourceName(sourceName)
}

func (b *OPCAreaBrowser) Release() error {
	b.browser.Release()
	return nil
}
