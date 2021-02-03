package xades4go

import (
	"crypto"
	"fmt"
	"errors"
)

const (
	// Canonicalization Algorithm
	CanonicalXML10Algorithm                            = "http://www.w3.org/TR/2001/RECxmlc14n20010315"
	CanonicalXML10WithCommentAlgorithm                 = "http://www.w3.org/TR/2001/RECxmlc14n20010315#WithComments"
	CanonicalXML11Algorithm                            = "http://www.w3.org/2006/12/xmlc14n11"
	CanonicalXML11WithCommentAlgorithm                 = "http://www.w3.org/2006/12/xmlc14n11#WithComments"
	ExclusiveXMLCanonicalization10Algorithm            = "http://www.w3.org/2001/10/xmlexcc14n#"
	ExclusiveXMLCanonicalization10WithCommentAlgorithm = "http://www.w3.org/2001/10/xmlexcc14n#WithComments"

	// Transform Algorithm
	Base64Algorithm                      = "http://www.w3.org/2000/09/xmldsig#base64"
	XPathFilteringAlgorithm              = "http://www.w3.org/TR/1999/REC-xpath-19991116"
	EnvelopedSignatureTransformAlgorithm = "http://www.w3.org/2000/09/xmldsig#enveloped-signature"
	XLSTTransformAlgorithm               = "http://www.w3.org/TR/1999/REC-xslt-19991116"

	// Digest Algorithm
	SHA1MessageDigestAlgorithm   = "http://www.w3.org/2000/09/xmldsig#sha1"
	SHA224MessageDigestAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#sha224"
	SHA256MessageDigestAlgorithm = "http://www.w3.org/2001/04/xmlenc#sha256"
	SHA384MessageDigestAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#sha384"
	SHA512MessageDigestAlgotithm = "http://www.w3.org/2001/04/xmlenc#sha512"

	// Signature Base64Algorithm
	DSASHA1SignatureAlgorithm = "http://www.w3.org/2000/09/xmldsig#dsa-sha1"
	DSASHA256SignatureAlgorithm = "http://www.w3.org/2009/xmldsig11#dsa-sha256"
	RSASHA1SignatureAlgorithm = "http://www.w3.org/2000/09/xmldsig#rsa-sha1"
	RSASHA224SignatureAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha224"
	RSASHA256SignatureAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"
	RSASHA384SignatureAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha384"
	RSASHA512SignatureAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha512"
	ECDSASHA1SignatureAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha1"
	ECDSASHA224SignatureAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha224"
	ECDSASHA256SignatureAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256"
	ECDSASHA384SignatureAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha384"
	ECDSASHA512SignatureAlgorithm = "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha512"
	)

// Transformer is an interface that perform Transform algorithm which follows https://www.w3.org/TR/xmldsig-core1/#sec-TransformAlg.
// interface{} is used as an input and output intentionally to provide freedom to implement Transformer with any 3rd-package. But for convenice, []byte will be used as a octa-stream in XMLDSig document.
type Transformer interface {
	Transform(input XML) (XML, error)
}

// Canonicalizer is an object that canonicalize XML to octet-stream.
type Canonicalizer interface {
	Canonicalize(input XML) ([]byte, error)
}

// Dereferencer is a object that dereference data object from URI (URI attribute of a given Reference element) and XML bytes to XML type.
type Dereferencer interface {
	Dereference(xmlContent []byte, uri string) (XML, error)
}

// XML is an input/output of/from Transformer, Canonicalizer and Dereferencer.
// According to https://www.w3.org/TR/xmldsig-core1, the input and output of Transformer, Canonicalizer and Dereferencer can be either Octet-Stream or NodeSet. The dedicated type for this input/output is needed.
// Since, currently, there is no standard XML library for Go (that support Node Set API); The NodeSet's type here is intentionally left with interface{} to provide the freedom for contributors to implement their own XML Node Set.
type XML struct {
	IsOctetStream bool
	OctetStream   []byte
	NodeSet       interface{}
}

// SignedInfoFactory is an abstact factory that creates Transformer, Canonicalizer and Dereferencer that can work together (have same NodeSet implementation).
// The name, SignedInfoFactory, is come from the fact that this factory only construct objects that related with creating SignedInfo element.
type SignedInfoFactory interface {
	CreateTransformer(algorithm string) (Transformer, error)
	CreateCanonicalizer(canonicalizationAlgorithm string) (Canonicalizer, error)
	CreateDereferencer() Dereferencer
}

// Digester is an object that perform digest algorithm on octet-stream input and return base64-encoded output.
type Digester interface {
	Digest(input []byte) ([]byte, error)
}

type cryptoDigester struct {
	h crypto.Hash
}

func (digester *cryptoDigester) Digest(input []byte) ([]byte, error) {
	hash := digester.h.New()
	_, err := hash.Write(input)
	if err != nil {
		return nil, fmt.Errorf("cannot digest using %s: %w", digester.h.String(), err)
	}
	return hash.Sum(nil), nil
}

func CreateDigester(algorithmName string) (Digester, error) {
	if algorithmName == "" {
		return nil, errors.New("Algorithm must not be empty")
	}
	switch algorithmName {
	case SHA1MessageDigestAlgorithm:
		return &cryptoDigester{h: crypto.SHA1}, nil
	case SHA224MessageDigestAlgorithm:
		return &cryptoDigester{h: crypto.SHA224}, nil
	case SHA256MessageDigestAlgorithm:
		return &cryptoDigester{h: crypto.SHA256}, nil
	case SHA384MessageDigestAlgorithm:
		return &cryptoDigester{h: crypto.SHA384}, nil
	case SHA512MessageDigestAlgotithm:
		return &cryptoDigester{h: crypto.SHA512}, nil
	}
	return nil, fmt.Errorf("this package does not implement %s digest algorithm", algorithmName)
}

func CreateDigesterForSignatureAlgorithm(signatureAlgorithm string) (Digester, error) {
	h, err := mapSignatureAlgorithmToCrytoHash(signatureAlgorithm)
	if err != nil {
		return nil, err
	}
	return &cryptoDigester{h: h}, nil
}

func mapSignatureAlgorithmToCrytoHash(signatureAlgorithm string) (crypto.Hash, error) {
	if signatureAlgorithm == "" {
		return 0, errors.New("Algorithm must not be empty")
	}
	switch signatureAlgorithm {
	case DSASHA1SignatureAlgorithm:
		return crypto.SHA1, nil
	case DSASHA256SignatureAlgorithm:
		return crypto.SHA256, nil
	case ECDSASHA1SignatureAlgorithm:
		return crypto.SHA1, nil
	case RSASHA224SignatureAlgorithm:
		return crypto.SHA224, nil
	case RSASHA256SignatureAlgorithm:
		return crypto.SHA256, nil
	case RSASHA384SignatureAlgorithm:
		return crypto.SHA384, nil
	case RSASHA512SignatureAlgorithm:
		return crypto.SHA512, nil
	case ECDSASHA224SignatureAlgorithm:
		return crypto.SHA224, nil
	case ECDSASHA256SignatureAlgorithm:
		return crypto.SHA256, nil
	case ECDSASHA384SignatureAlgorithm:
		return crypto.SHA384, nil
	case ECDSASHA512SignatureAlgorithm:
		return crypto.SHA512, nil
	}
	return 0, fmt.Errorf("this package does not implement %s signature algorithm", signatureAlgorithm)
}
