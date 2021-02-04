package xades4go

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/beevik/etree"
)

type XMLDSigSignatureValidator struct {
	signedInfoFactory                SignedInfoFactory
	defaultCanonicalizationAlgorithm string
}

func NewXMLDSigSignatureValidator(signedInfoFactory SignedInfoFactory) SignatureValidator {
	return &XMLDSigSignatureValidator{
		signedInfoFactory:                signedInfoFactory,
		defaultCanonicalizationAlgorithm: CanonicalXML10Algorithm,
	}
}

func (validator *XMLDSigSignatureValidator) Validate(xmlBytes []byte) (ValidationResult, error) {
	rootElement, err := createEtreeElementFromXMLBytes(xmlBytes)
	if err != nil {
		return ValidationResult{}, err
	}
	signatureElement, err := mustFoundOnlyOneElement(rootElement, signatureElementTag)
	if err != nil {
		return ValidationResult{}, err
	}
	signedInfoElement, err := mustFoundOnlyOneChildElement(signatureElement, signedInfoElementTag)
	if err != nil {
		return ValidationResult{}, err
	}
	references, err := mustFoundAtLeastOneChildElement(signedInfoElement, referenceElementTag)
	if err != nil {
		return ValidationResult{}, err
	}
	result := ValidationResult{}
	for referenceIndex, reference := range references {
		uriAttribute, err := mustFoundAttribute(reference, uriAttributeKey)
		if err != nil {
			return ValidationResult{}, errors.New("this validator does not support anonymous referecing (no URI attribute)")
		}
		xmlInput, err := validator.signedInfoFactory.CreateDereferencer().DereferenceByURI(xmlBytes, uriAttribute.Value)
		if err != nil {
			return ValidationResult{}, fmt.Errorf("at Reference#%d element, cannot dereference the given URI: %w", referenceIndex, err)
		}
		transformsElement, err := mustFoundOnlyOneIfFound(reference, transformsElementTag)
		if err != nil {
			return ValidationResult{}, fmt.Errorf("at Reference#%d element: %w", referenceIndex, err)
		}
		if transformsElement != nil {
			transformElements, err := mustFoundAtLeastOneChildElement(transformsElement, transformElementTag)
			if err != nil {
				return ValidationResult{}, fmt.Errorf("at Reference#%d element: %w", referenceIndex, err)
			}
			for transformIndex, transformElement := range transformElements {
				algorithmAttribute, err := mustFoundAttribute(transformElement, algorithmAttributeKey)
				if err != nil {
					return ValidationResult{}, fmt.Errorf("at Transform#%d element of Reference#%d element: %w", transformIndex, referenceIndex, err)
				}
				transformer, err := validator.signedInfoFactory.CreateTransformer(algorithmAttribute.Value)
				if err != nil {
					return ValidationResult{}, fmt.Errorf("error while creating Transformer at Transform#%d element of Reference#%d element: %w", transformIndex, referenceIndex, err)
				}
				xmlInput, err = transformer.Transform(xmlInput)
				if err != nil {
					return ValidationResult{}, fmt.Errorf("error while transforming at Transform#%d element of Reference#%d element: %w", transformIndex, referenceIndex, err)
				}
			}
		}
		transformedDataObjectToBeDigested := xmlInput.OctetStream
		if !xmlInput.IsOctetStream {
			canonicalizer, err := validator.signedInfoFactory.CreateCanonicalizer(validator.defaultCanonicalizationAlgorithm)
			if err != nil {
				return ValidationResult{}, fmt.Errorf("error while creating canonicalizer at Reference#%d: %w", referenceIndex, err)
			}
			transformedDataObjectToBeDigested, err = canonicalizer.Canonicalize(xmlInput)
			if err != nil {
				return ValidationResult{}, fmt.Errorf("error while canonicalizing at Reference#%d: %w", referenceIndex, err)
			}
		}
		digestMethodElement, err := mustFoundOnlyOneChildElement(reference, digestMethodElementTag)
		if err != nil {
			return ValidationResult{}, fmt.Errorf("at Reference#%d element: %w", referenceIndex, err)
		}
		algorithmAttribute, err := mustFoundAttribute(digestMethodElement, algorithmAttributeKey)
		if err != nil {
			return ValidationResult{}, fmt.Errorf("at Reference#%d element: %w", referenceIndex, err)
		}
		digester, err := CreateDigester(algorithmAttribute.Value)
		if err != nil {
			return ValidationResult{}, fmt.Errorf("error while creating Digester at Reference#%d: %w", referenceIndex, err)
		}
		generatedDigestValue, err := digester.Digest(transformedDataObjectToBeDigested)
		if err != nil {
			return ValidationResult{}, fmt.Errorf("error while digesting at Reference#%d: %w", referenceIndex, err)
		}
		digestValueElement, err := mustFoundOnlyOneChildElement(reference, digestValueElementTag)
		if err != nil {
			return ValidationResult{}, fmt.Errorf("at Reference#%d element: %w", referenceIndex, err)
		}
		digestValue := []byte(digestValueElement.Text())
		referenceValidationResult := ReferenceValidationResult{
			IsValid:              false,
			GeneratedDigestValue: string(generatedDigestValue),
			DigestValue:          string(digestValue),
		}
		if bytes.Compare(generatedDigestValue, digestValue) == 0 {
			referenceValidationResult.IsValid = true
		}
		result.ReferenceValidationResults = append(result.ReferenceValidationResults, referenceValidationResult)
	}
	canonicalizationMethodElement, err := mustFoundOnlyOneChildElement(signedInfoElement, canonicalizationMethodElementTag)
	if err != nil {
		return ValidationResult{}, err
	}
	algorithmAttribute, err := mustFoundAttribute(canonicalizationMethodElement, algorithmAttributeKey)
	if err != nil {
		return ValidationResult{}, err
	}
	canonicalizationAlgorithm := algorithmAttribute.Value
	signatureMethodElement, err := mustFoundOnlyOneChildElement(signedInfoElement, signatureMethodElementTag)
	if err != nil {
		return ValidationResult{}, err
	}
	algorithmAttribute, err = mustFoundAttribute(signatureMethodElement, algorithmAttributeKey)
	if err != nil {
		return ValidationResult{}, err
	}
	signedInfoInput, err := validator.signedInfoFactory.CreateDereferencer().DereferenceByPath(xmlBytes, "//"+signatureElementTag+"/"+signedInfoElementTag)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("cannot derefernce SignedInfo element: %w", err)
	}
	canonicalizer, err := validator.signedInfoFactory.CreateCanonicalizer(canonicalizationAlgorithm)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("cannot create canonicalizer from CanonicalizationMethod element: %w", err)
	}
	canonicalizedSignedInfo, err := canonicalizer.Canonicalize(signedInfoInput)
	if err != nil {
		return ValidationResult{}, fmt.Errorf("error while canonicalizing SignedInfo element: %w", err)
	}
	signatureMethodAlgorithm := algorithmAttribute.Value
	keyInfoElement, err := mustFoundOnlyOneIfFound(signatureElement, keyInfoElementTag)
	if err != nil {
		return ValidationResult{}, err
	}
	signatureValueElement, err := mustFoundOnlyOneChildElement(signatureElement, signatureValueElementTag)
	if err != nil {
		return ValidationResult{}, err
	}
	signatureValue := signatureValueElement.Text()
	possibleSignatureVerifiers, err := createPossibleSignatureVerifiersFromKeyInfoElement(keyInfoElement, signatureMethodAlgorithm)
	if err != nil {
		return ValidationResult{}, err
	}
	isSignatureValid := false
	for _, signatureVerifier := range possibleSignatureVerifiers {
		err := signatureVerifier.Verify(signatureMethodAlgorithm, canonicalizedSignedInfo, []byte(signatureValue))
		if err == nil {
			isSignatureValid = true
			break
		}
	}
	result.IsSignatureValid = isSignatureValid
	return result, nil
}

func createEtreeElementFromXMLBytes(xmlBytes []byte) (*etree.Element, error) {
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(xmlBytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse xmlBytes to etree's element: %w", err)
	}
	return doc.Root(), nil
}

func mustFoundOnlyOneElement(root *etree.Element, tag string) (*etree.Element, error) {
	foundElements := root.FindElements("//" + tag)
	if len(foundElements) == 0 {
		return nil, fmt.Errorf("%s element not found", tag)
	}
	if len(foundElements) > 1 {
		return nil, fmt.Errorf("found more than one %s element", tag)
	}
	return foundElements[0], nil
}

func mustFoundOnlyOneChildElement(parent *etree.Element, childTag string) (*etree.Element, error) {
	foundElements := parent.SelectElements(childTag)
	if len(foundElements) == 0 {
		return nil, fmt.Errorf("%s element was not found on %s element", childTag, parent.FullTag())
	}
	if len(foundElements) > 1 {
		return nil, fmt.Errorf("found more than one %s element on %s element", childTag, parent.FullTag())
	}
	return foundElements[0], nil
}

func mustFoundAtLeastOneChildElement(parent *etree.Element, childTag string) ([]*etree.Element, error) {
	foundElements := parent.SelectElements(childTag)
	if len(foundElements) == 0 {
		return nil, fmt.Errorf("%s element was not found on %s element", childTag, parent.FullTag())
	}
	return foundElements, nil
}

func mustFoundAttribute(element *etree.Element, attributeKey string) (etree.Attr, error) {
	attribute := element.SelectAttr(attributeKey)
	if attribute == nil {
		return etree.Attr{}, fmt.Errorf("attribute %s is not found on %s element", attributeKey, element.FullTag())
	}
	return *attribute, nil
}

func mustFoundOnlyOneIfFound(parent *etree.Element, childTag string) (*etree.Element, error) {
	foundElements := parent.SelectElements(childTag)
	if len(foundElements) > 1 {
		return nil, fmt.Errorf("found more than one %s elemnt on %s element", childTag, parent.FullTag())
	}
	if len(foundElements) == 1 {
		return foundElements[0], nil
	}
	return nil, nil
}

func createPossibleSignatureVerifiersFromKeyInfoElement(keyInfoElement *etree.Element, signatureAlgorithm string) ([]SignatureValueVerifier, error) {
	result := make([]SignatureValueVerifier, 0)
	x509DataElements := keyInfoElement.SelectElements(x509DataElementTag)
	if len(x509DataElements) > 0 {
		for _, x509Element := range x509DataElements {
			x509CertificateElements := x509Element.SelectElements(x509CertificateElementTag)
			if len(x509CertificateElements) > 0 {
				for _, x509CertificateElement := range x509CertificateElements {
					asn1Certificate, err := base64.StdEncoding.DecodeString(x509CertificateElement.Text())
					if err != nil {
						return nil, errors.New("cannot base64-decode attached certificate: " + err.Error())
					}
					certificate, err := x509.ParseCertificate(asn1Certificate)
					if err != nil {
						return nil, errors.New("cannot parse attached certificate: " + err.Error())
					}
					switch pub := certificate.PublicKey.(type) {
					case *rsa.PublicKey:
						result = append(result, &rsaSignatureValueVerifier{rsaPublicKey: pub})
					}
				}
			}
		}
	}
	return result, nil
}
