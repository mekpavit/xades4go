package xades4go

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
)

const (
	signatureElementTag              = "Signature"
	signedInfoElementTag             = "SignedInfo"
	referenceElementTag              = "Reference"
	transformsElementTag             = "Transforms"
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
	Validate(xmlBytes []byte) (ValidationResult, error)
}

type ValidationResult struct {
	ReferenceValidationResults []ReferenceValidationResult
	IsSignatureValid           bool
}

type ReferenceValidationResult struct {
	IsValid              bool
	GeneratedDigestValue string
	DigestValue          string
}

// SignatureValueVerifier is an object that verify Base64-encoded signature value (in SignatureValue element) against canonicalized SignedInfo element using the given signature algorithm.
type SignatureValueVerifier interface {
	Verify(signatureAlgorithm string, canonicalizedSignedInfo []byte, base64SignatureValue []byte) error
}

type rsaSignatureValueVerifier struct {
	rsaPublicKey *rsa.PublicKey
}

func (verifier *rsaSignatureValueVerifier) Verify(signatureAlgorithm string, canonicalizedSignedInfo []byte, base64SignatureValue []byte) error {
	hashAlgorithm, err := mapSignatureAlgorithmToCrytoHash(signatureAlgorithm)
	if err != nil {
		return err
	}
	h := hashAlgorithm.New()
	_, err = h.Write(canonicalizedSignedInfo)
	if err != nil {
		return fmt.Errorf("cannot hash SignedInfo using %s: %w", hashAlgorithm.String(), err)
	}
	signatureValue, err := base64.StdEncoding.DecodeString(string(base64SignatureValue))
	if err != nil {
		return fmt.Errorf("SignatureValue is not base64-encoded: %w", err)
	}
	return rsa.VerifyPKCS1v15(verifier.rsaPublicKey, hashAlgorithm, h.Sum(nil), []byte(signatureValue))
}
