package etreeimpl

import (
	"errors"

	"github.com/beevik/etree"
	"github.com/mekpavit/xades4go"
)

type envelopedSignatureTransformer struct{}

func (transformer *envelopedSignatureTransformer) Transform(input xades4go.XML) (xades4go.XML, error) {
	var inputNodeSet *etree.Element
	if input.IsOctetStream {
		var err error
		inputNodeSet, err = createNodeSetFromBytes(input.OctetStream)
		if err != nil {
			return xades4go.XML{}, err
		}
	} else {
		var ok bool
		inputNodeSet, ok = input.NodeSet.(*etree.Element)
		if !ok {
			return xades4go.XML{}, errors.New("input must be []byte or *etree.Element")
		}
	}
	outputNodeSet, err := transformer.transform(inputNodeSet)
	if err != nil {
		return xades4go.XML{}, err
	}
	return xades4go.XML{IsOctetStream: false, NodeSet: outputNodeSet}, nil
}

func (transformer *envelopedSignatureTransformer) transform(nodeSet *etree.Element) (*etree.Element, error) {
	signatureElement := nodeSet.FindElement("//Signature")
	if signatureElement == nil {
		return nodeSet, nil
	}
	parentOfSignatureElement := signatureElement.Parent()
	if parentOfSignatureElement == nil {
		return nodeSet, nil
	}
	parentOfSignatureElement.RemoveChild(signatureElement)
	return nodeSet, nil
}
