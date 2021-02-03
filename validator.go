package xades4go

import "crypto/rsa"

const (
	signatureElementTag              = "Signature"
	signedInfoElementTag             = "SignedInfo"
	referenceElementTag              = "Reference"
	transformsElementTag             = "Transformers"
	transformElementTag              = "Transform"
	digestMethodElementTag           = "DigestMethod"
	digestValueElementTag            = "DigestValue"
	keyInfoElementTag                = "KeyInfo"
	canonicalizationMethodElementTag = "CanonicalizationMethod"
	signatureMethodElementTag        = "SignatureMethod"
	signatureValueElementTag         = "SignatureValue"

	x509DataElementTag        = "X509Data"
	x509CertificateElementTag = "X509Certificate"

	uriAttributeKey       = "URI"
	algorithmAttributeKey = "Algorithm"
)

type SignatureValidator interface {
	Validate(xmlByte []byte) (ValidationResult, error)
}

type ValidationResult struct {
	ReferenceValidationResults []ReferenceValidationResult
	IsSignatureValid           bool
}

type ReferenceValidationResult struct {
	IsValid              bool
	GeneratedDigestValue []byte
	DigestValue          []byte
}

type SignatureValueVerifier interface {
	Verify(signatureAlgorithm string, hashed []byte, signatureValue []byte) error
}

type rsaSignatureValueVerifier struct {
	rsaPublicKey *rsa.PublicKey
}

func (verifier *rsaSignatureValueVerifier) Verify(signatureAlgorithm string, hashed []byte, signatureValue []byte) error {
	hashAlgorithm, err := mapSignatureAlgorithmToCrytoHash(signatureAlgorithm)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(verifier.rsaPublicKey, hashAlgorithm, hashed, signatureValue)
}
