package etreeimpl

import (
	"fmt"
	"strings"

	"github.com/mekpavit/xades4go"
)

type dereferencer struct{}

func (d *dereferencer) DereferenceByURI(xmlContent []byte, uri string) (xades4go.XML, error) {
	nodeSet, err := createNodeSetFromBytes(xmlContent)
	if err != nil {
		return xades4go.XML{}, err
	}
	if uri == "" {
		return xades4go.XML{IsOctetStream: false, NodeSet: nodeSet}, nil
	}
	idOfDataObject := strings.TrimPrefix(uri, "#")
	dereferencedNodeSet := nodeSet.FindElement(fmt.Sprintf("//[@Id='%s']", idOfDataObject))
	if dereferencedNodeSet == nil {
		return xades4go.XML{}, fmt.Errorf("cannot find any node set from uri -> %s", uri)
	}
	return xades4go.XML{IsOctetStream: false, NodeSet: dereferencedNodeSet}, nil
}

func (d *dereferencer) DereferenceByPath(xmlContent []byte, path string) (xades4go.XML, error) {
	nodeSet, err := createNodeSetFromBytes(xmlContent)
	if err != nil {
		return xades4go.XML{}, err
	}
	dereferencedNodeSet := nodeSet.FindElement(path)
	if dereferencedNodeSet == nil {
		return xades4go.XML{}, fmt.Errorf("cannot find any node set from path -> %s", path)
	}
	return xades4go.XML{IsOctetStream: false, NodeSet: dereferencedNodeSet}, nil
}
