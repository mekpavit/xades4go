package xades4go

const (
	signatureElementTag    = "Signature"
	signedInfoElementTag   = "SignedInfo"
	referenceElementTag    = "Reference"
	transformsElementTag   = "Transformers"
	transformElementTag    = "Transform"
	digestMethodElementTag = "DigestMethod"
	digestValueElementTag  = "DigestValue"

	uriAttributeKey       = "URI"
	algorithmAttributeKey = "Algorithm"
)

type SignatureValidator interface {
	Validate(xmlByte []byte) (ValidationResult, error)
}

type ValidationResult struct {
	ReferenceValidationResults []ReferenceValidationResult
}

type ReferenceValidationResult struct {
	IsValid              bool
	GeneratedDigestValue []byte
	DigestValue          []byte
}
