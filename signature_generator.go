package xades4go

type SignatureGenerator interface {
	SignXMLBytes(xmlBytes []byte, dataObjectReferences []ReferenceGenerationDetail) ([]byte, error)
}

type ReferenceGenerationDetail struct {
	URIOfDataObjectBeingSigned string
	TransformAlgorithms        []string
	DigestAlgorithm            string
}
