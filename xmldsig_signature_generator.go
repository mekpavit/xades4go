package xades4go

type XMLDSigSignatureGenerator struct {
}

func (generator *XMLDSigSignatureGenerator) SignXMLBytes(xmlBytes []byte, dataObjectReferences []ReferenceGenerationDetail) ([]byte, error) {
	_, err := createEtreeElementFromXMLBytes(xmlBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
