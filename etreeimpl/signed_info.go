package etreeimpl

import (
	"errors"
	"fmt"

	"github.com/mekpavit/xades4go"
)

type signedInfoFactory struct{}

func NewSignedInfoFactory() xades4go.SignedInfoFactory {
	return &signedInfoFactory{}
}

func (factory *signedInfoFactory) CreateTransformer(algorithmName string) (xades4go.Transformer, error) {
	if algorithmName == "" {
		return nil, errors.New("Algorithm must not be empty")
	}
	switch algorithmName {
	case xades4go.CanonicalXML10Algorithm:
		return &canonicalXML10Canonicalizer{}, nil
	case xades4go.CanonicalXML10WithCommentAlgorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", algorithmName)
	case xades4go.CanonicalXML11Algorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", algorithmName)
	case xades4go.CanonicalXML11WithCommentAlgorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", algorithmName)
	case xades4go.ExclusiveXMLCanonicalization10Algorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", algorithmName)
	case xades4go.ExclusiveXMLCanonicalization10WithCommentAlgorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", algorithmName)

	case xades4go.Base64Algorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", algorithmName)
	case xades4go.XPathFilteringAlgorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", algorithmName)
	case xades4go.EnvelopedSignatureTransformAlgorithm:
		return &envelopedSignatureTransformer{}, nil
	case xades4go.XLSTTransformAlgorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", algorithmName)
	}
	return nil, fmt.Errorf("%s was not an acceptable Transform algorithm", algorithmName)
}

func (factory *signedInfoFactory) CreateCanonicalizer(canonicalizationAlgorithm string) (xades4go.Canonicalizer, error) {
	if canonicalizationAlgorithm == "" {
		return nil, errors.New("Algorithm must not be empty")
	}
	switch canonicalizationAlgorithm {
	case xades4go.CanonicalXML10Algorithm:
		return &canonicalXML10Canonicalizer{}, nil
	case xades4go.CanonicalXML10WithCommentAlgorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", canonicalizationAlgorithm)
	case xades4go.CanonicalXML11Algorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", canonicalizationAlgorithm)
	case xades4go.CanonicalXML11WithCommentAlgorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", canonicalizationAlgorithm)
	case xades4go.ExclusiveXMLCanonicalization10Algorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", canonicalizationAlgorithm)
	case xades4go.ExclusiveXMLCanonicalization10WithCommentAlgorithm:
		return nil, fmt.Errorf("%s was not implemented by etreeimpl", canonicalizationAlgorithm)
	}
	return nil, fmt.Errorf("%s was not an acceptable Canonicalization algorithm", canonicalizationAlgorithm)
}

func (factory *signedInfoFactory) CreateDereferencer() xades4go.Dereferencer {
	return &dereferencer{}
}
