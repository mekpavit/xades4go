package etreeimpl

import (
	"errors"
	"fmt"
	"sort"

	"github.com/beevik/etree"
	"github.com/mekpavit/xades4go"
)

// C14NTransformer is a XML canonicalizer that follows https://www.w3.org/TR/2001/REC-xml-c14n-20010315 processing model.
// It's still not support some type of XML which are
// 1. XML that contains <!ENTITY ...>
// 2. XML that contains <!ATTLIST ...>
// 3. XML that contains Processing Instruction nodes
// 4. XML that contains empty default namespace (xmlns="")
type canonicalXML10Canonicalizer struct {
}

func (transformer *canonicalXML10Canonicalizer) Transform(input xades4go.XML) (xades4go.XML, error) {
	canonicalizedXML, err := transformer.Canonicalize(input)
	if err != nil {
		return xades4go.XML{}, err
	}
	return xades4go.XML{IsOctetStream: true, OctetStream: canonicalizedXML}, nil
}

func createNodeSetFromBytes(xmlContent []byte) (*etree.Element, error) {
	etreeDoc := etree.NewDocument()
	err := etreeDoc.ReadFromBytes(xmlContent)
	if err != nil {
		return nil, fmt.Errorf("error while parsing XML bytes to node set: %w", err)
	}
	return etreeDoc.Root(), nil
}

func (transformer *canonicalXML10Canonicalizer) Canonicalize(input xades4go.XML) ([]byte, error) {
	var inputNodeSet *etree.Element
	if input.IsOctetStream {
		var err error
		inputNodeSet, err = createNodeSetFromBytes(input.OctetStream)
		if err != nil {
			return nil, err
		}
	} else {
		var ok bool
		inputNodeSet, ok = input.NodeSet.(*etree.Element)
		if !ok {
			return nil, errors.New("input must be []byte or *etree.Element")
		}
	}
	resultNodeSet := inputNodeSet.Copy()
	directAncestorNamespaces := collectAncestorNamespaces(inputNodeSet)
	propagateAncestorNamespacesTo(resultNodeSet, directAncestorNamespaces)
	traverseAndTransformElement(resultNodeSet, map[string]string{})
	return completeCanonicalization(resultNodeSet)
}

func traverseAndTransformElement(element *etree.Element, directAncestorNamespaces map[string]string) {
	copyOfDirectAncestorNamespaces := make(map[string]string)
	for namespaceFullKey, namespaceURI := range directAncestorNamespaces {
		copyOfDirectAncestorNamespaces[namespaceFullKey] = namespaceURI
	}

	for _, attr := range element.Attr {
		if !isNamespace(attr) {
			continue
		}
		declaredNamespaceURI, isRedeclaredNamespace := directAncestorNamespaces[attr.FullKey()]
		if !isRedeclaredNamespace || (isRedeclaredNamespace && attr.Value != declaredNamespaceURI) {
			copyOfDirectAncestorNamespaces[attr.FullKey()] = attr.Value
			continue
		}
		element.RemoveAttr(attr.FullKey())
	}

	sort.Sort(attributesByLexicographicalOrder(element.Attr))

	for _, child := range element.Child {
		if childElement, isElementNode := child.(*etree.Element); isElementNode {
			traverseAndTransformElement(childElement, copyOfDirectAncestorNamespaces)
		}
		if _, isCommentNode := child.(*etree.Comment); isCommentNode {
			element.RemoveChild(child)
		}
	}
}

func collectAncestorNamespaces(element *etree.Element) map[string]etree.Attr {
	result := make(map[string]etree.Attr)
	for parent := element.Parent(); parent != nil; parent = parent.Parent() {
		for _, attr := range parent.Attr {
			if !isNamespace(attr) {
				continue
			}
			if _, isAlreadyCollected := result[attr.FullKey()]; isAlreadyCollected {
				continue
			}
			result[attr.FullKey()] = attr
		}
	}
	return result
}

func propagateAncestorNamespacesTo(element *etree.Element, ancestorNamespaces map[string]etree.Attr) {
	for namespaceFullKey, namespaceAttr := range ancestorNamespaces {
		if redeclaredNamespace := element.SelectAttr(namespaceFullKey); redeclaredNamespace != nil {
			continue
		}
		element.Attr = append(element.Attr, namespaceAttr)
	}
}

type attributesByLexicographicalOrder []etree.Attr

func (attributes attributesByLexicographicalOrder) Len() int {
	return len(attributes)
}

func (attributes attributesByLexicographicalOrder) Less(i int, j int) bool {
	x, y := attributes[i], attributes[j]
	if isNamespace(x) && isNamespace(y) {
		return x.FullKey() < y.FullKey()
	}
	if !isNamespace(x) && !isNamespace(y) {
		if isUnqualifiedAttribute(x) && isUnqualifiedAttribute(y) {
			return x.Key < y.Key
		}
		if !isUnqualifiedAttribute(x) && !isUnqualifiedAttribute(y) {
			if x.NamespaceURI() != y.NamespaceURI() {
				return x.NamespaceURI() < y.NamespaceURI()
			}
			return x.Key < y.Key
		}
		if isUnqualifiedAttribute(x) && !isUnqualifiedAttribute(y) {
			return true
		}
		return false
	}
	if isNamespace(x) && !isNamespace(y) {
		return true
	}
	return false
}

func (attributes attributesByLexicographicalOrder) Swap(i int, j int) {
	attributes[i], attributes[j] = attributes[j], attributes[i]
}

func isNamespace(attribute etree.Attr) bool {
	return attribute.Space == "xmlns" || attribute.FullKey() == "xmlns"
}

func isUnqualifiedAttribute(attribute etree.Attr) bool {
	return !isNamespace(attribute) && attribute.Space == ""
}

// completeCanonicalization performs the canonicalization that does not include in the Canonicalization Transformer.
// If any of Canonicalization is used, this function MUST be called to ensure the correctness of canonicalization.
//
// The canonicalization handled by this function are:
// 1. Canonicalize end tags from <aaa    /> to <aaa></aaa>
// 2. Canonicalize text nodes
// 3. Canonicalize attribute nodes
func completeCanonicalization(element *etree.Element) ([]byte, error) {
	doc := etree.NewDocument()
	doc.SetRoot(element)
	doc.WriteSettings = etree.WriteSettings{
		CanonicalEndTags: true,
		CanonicalText:    true,
		CanonicalAttrVal: true,
		UseCRLF:          false,
	}
	result, err := doc.WriteToBytes()
	if err != nil {
		return nil, fmt.Errorf("error while canonicalize element contents and convert to bytes: %w", err)
	}
	return result, nil
}
