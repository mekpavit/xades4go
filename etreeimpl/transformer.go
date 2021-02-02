package etreeimpl

import (
	"errors"
	"fmt"

	"github.com/beevik/etree"
	"github.com/mekpavit/xades4go"
)

type signedInfoFactory struct{}

func (factory *signedInfoFactory) CreateTransformer(algorithmName string) (xades4go.Transformer, error) {
	if algorithmName == "" {
		return nil, errors.New("Algorithm must not be empty")
	}
	switch algorithmName {
	case "http://www.w3.org/TR/2001/RECxmlc14n20010315":
		return &canonicalXML10Canonicalizer{}, nil
	case "http://www.w3.org/TR/2001/RECxmlc14n20010315#WithComments":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/2006/12/xmlc14n11":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/2006/12/xmlc14n11#WithComments":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/2001/10/xmlexcc14n#":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/2001/10/xmlexcc14n#WithComments":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)

	case "http://www.w3.org/2000/09/xmldsig#base64":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/TR/1999/RECxpath19991116":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/2000/09/xmldsig#envelopedsignature":
		return &envelopedSignatureTransformer{}, nil
	case "http://www.w3.org/TR/1999/RECxslt19991116":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	}
	return nil, fmt.Errorf("%s was not an acceptable Transform algorithm", algorithmName)
}

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
	resultNodeSet := nodeSet.Copy()
	signatureElement := resultNodeSet.FindElement("//Signature")
	if signatureElement == nil {
		return resultNodeSet, nil
	}
	parentOfSignatureElement := signatureElement.Parent()
	if parentOfSignatureElement == nil {
		return resultNodeSet, nil
	}
	parentOfSignatureElement.RemoveChild(signatureElement)
	return resultNodeSet, nil
}
