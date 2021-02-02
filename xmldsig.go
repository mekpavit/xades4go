package xades4go

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/beevik/etree"
)

type XMLDSigValidator struct {
	signedInfoFactory                SignedInfoFactory
	defaultCanonicalizationAlgorithm string
}

func (validator *XMLDSigValidator) Validate(xmlBytes []byte) (ValidationResult, error) {
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
		xmlInput, err := validator.signedInfoFactory.CreateDereferencer().Dereference(xmlBytes, uriAttribute.Value)
		if err != nil {
			return ValidationResult{}, fmt.Errorf("at Reference#%d element, cannot dereference the given URI: %w", referenceIndex, err)
		}
		transformsElement := reference.SelectElement(transformsElementTag)
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
		digestMethodElement, err := mustFoundOnlyOneElement(reference, digestMethodElementTag)
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
		digestValueElement, err := mustFoundOnlyOneElement(reference, digestValueElementTag)
		if err != nil {
			return ValidationResult{}, fmt.Errorf("at Reference#%d element: %w", referenceIndex, err)
		}
		digestValue := []byte(digestValueElement.Text())
		referenceValidationResult := ReferenceValidationResult{
			IsValid:              false,
			GeneratedDigestValue: generatedDigestValue,
			DigestValue:          digestValue,
		}
		if bytes.Compare(generatedDigestValue, digestValue) == 0 {
			referenceValidationResult.IsValid = true
		}
		result.ReferenceValidationResults = append(result.ReferenceValidationResults, referenceValidationResult)
	}
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
