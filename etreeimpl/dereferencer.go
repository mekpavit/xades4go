package etreeimpl

import (
	"fmt"

	"github.com/mekpavit/xades4go"
)

type dereferencer struct{}

func (d *dereferencer) Dereference(xmlContent []byte, uri string) (xades4go.XML, error) {
	nodeSet, err := createNodeSetFromBytes(xmlContent)
	if err != nil {
		return xades4go.XML{}, err
	}
	dereferencedNodeSet := nodeSet.FindElement("")
	if dereferencedNodeSet == nil {
		return xades4go.XML{}, fmt.Errorf("cannot find any node set from %s", uri)
	}
	return xades4go.XML{IsOctetStream: false, NodeSet: dereferencedNodeSet}, nil
}
