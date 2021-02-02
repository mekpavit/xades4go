package xades4go

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
