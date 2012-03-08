package html

//#include "helper.h"
import "C"
import (
	"unsafe"
	"os"
	"bytes"
	"gokogiri/xml"
)

var fragmentWrapperStart = []byte("<html><body>")
var bodySigBytes = []byte("<body")

var ErrFailParseFragment = os.NewError("failed to parse html fragment")

const initChildrenNumber = 4
func ParseFragment(document xml.Document, content, encoding, url []byte, options int) (fragment *xml.DocumentFragment, err os.Error) {
	//deal with trivial cases
	if len(content) == 0 { return }

	//if a document is not provided, we should create an empty Html document
	//a fragment must reside in a document
	if document == nil {
		document = CreateEmptyDocument(encoding)
	} 
	
	containBody := (bytes.Index(content, bodySigBytes) >= 0)
	
	//wrap the content
	content = append(fragmentWrapperStart, content...)

	//set up pointers before calling the C function
	var contentPtr, urlPtr unsafe.Pointer
	contentPtr   = unsafe.Pointer(&content[0])
	contentLen   := len(content)
	if len(url) > 0  { urlPtr = unsafe.Pointer(&url[0]) }
	
	htmlPtr := C.htmlParseFragment(document.DocPtr(), contentPtr, C.int(contentLen), urlPtr, C.int(options), nil, 0)
	

	//Note we've parsed the fragment within the given document 
	//the root is not the root of the document; rather it's the root of the subtree from the fragment
	html := xml.NewNode(unsafe.Pointer(htmlPtr), document)

	if html == nil {
		err = ErrFailParseFragment
		return
	}
	root := html
	if ! containBody {
		root = html.FirstChild()
		html.Remove() //remove html otherwise it's leaked
	}

	fragment = &xml.DocumentFragment{}
	fragment.Node = root

	nodes := make([]xml.Node, 0, initChildrenNumber)
	child := root.FirstChild()
	for ; child != nil; child = child.NextSibling() {
		nodes = append(nodes, child)
	}
	fragment.Children = xml.NewNodeSet(document, nodes)
	document.BookkeepFragment(fragment)
	return
}