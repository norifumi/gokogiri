package test

import (
	"gokogiri"
    "gokogiri/help"
	"testing"
	"strings"
)


func TestCDataNodeSetContent(t *testing.T) {
	doc := gokogiri.HtmlParseString("<html></html>")
	htmlNode := doc.RootElement()
	scriptNode := htmlNode.NewChild("script", "")
	scriptNode.SetCDataContent("//<![CDATA[\nalert('boo')\n//]]>")
	if !strings.Contains(doc.DumpHTML(), "<script>//<![CDATA[") {
		t.Error("Should have actually made a CDATA tag")
	}
    doc.Free()
    help.XmlCleanUpParser()
    if help.XmlMemoryAllocation() != 0 {
        t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
        help.XmlMemoryLeakReport()
    }

}
