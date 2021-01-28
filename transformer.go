package xades4go

import (
	"context"
	"fmt"
	"errors"

	"github.com/beevik/etree"
)

type Transformer interface {
	Transform(ctx context.Context, input interface{}) (interface{}, error)
}

type EnvelopedSignatureTransformer struct{}

func (transformer *EnvelopedSignatureTransformer) Transform(ctx context.Context, input interface{}) (interface{}, error) {
	inputNodeSet, ok := input.(*etree.Element)
	if !ok {
		return nil, errors.New("input must be []byte or *etree.Element")
	}
	return transformer.transform(ctx, inputNodeSet)
}

func (transformer *EnvelopedSignatureTransformer) transform(ctx context.Context, nodeSet *etree.Element) (*etree.Element, error) {
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

func CreateTransformerFromAlgorithm(algorithmName string) (Transformer, error) {
	if algorithmName == "" {
		return nil, errors.New("Algorithm must not be empty")
	}
	switch algorithmName {
	case "http://www.w3.org/TR/2001/RECxmlc14n20010315":
		return &C14NTransformer{}, nil
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
		return &EnvelopedSignatureTransformer{}, nil
	case "http://www.w3.org/TR/1999/RECxslt19991116":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	}
	return nil, fmt.Errorf("%s was not an acceptable Transform algorithm", algorithmName)
}
