package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/beevik/etree"
)

type Transformer interface {
	Transform(ctx context.Context, element *etree.Element) (*etree.Element, error)
}

func CreateTransformerFromAlgorithm(algorithmName string) (Transformer, error) {
	if algorithmName == "" {
		return nil, errors.New("Algorithm must not be empty")
	}
	switch algorithmName {
	case "http://www.w3.org/TR/2001/REC-xml-c14n-20010315":
		return &C14NTransformer{}, nil
	case "http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/2006/12/xml-c14n11":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/2006/12/xml-c14n11#WithComments":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/2001/10/xml-exc-c14n#":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/2001/10/xml-exc-c14n#WithComments":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)

	case "http://www.w3.org/2000/09/xmldsig#base64":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/TR/1999/REC-xpath-19991116":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	case "http://www.w3.org/2000/09/xmldsig#enveloped-signature":
		return &EnvelopedSignatureTransformer{}, nil
	case "http://www.w3.org/TR/1999/REC-xslt-19991116":
		return nil, fmt.Errorf("%s was not implemented by this package", algorithmName)
	}
	return nil, fmt.Errorf("%s was not an acceptable Transform algorithm", algorithmName)
}
