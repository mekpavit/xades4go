package internal

import (
	"context"
	"fmt"
	"sort"

	"github.com/beevik/etree"
)

// C14N11Transformer is a XML canonicalizer that follows https://www.w3.org/TR/xml-c14n11/ processing model.
// It's still not support some type of XML which are
// 1. XML that contains nodes with "xml:base" attribute
// 2. XML that contains <!ENTITY ...>
// 3. XML that contains <!ATTLIST ...>
// 4. XML that contains Processing Instruction nodes
type C14N11Transformer struct{}

func (transformer *C14N11Transformer) Transform(ctx context.Context, nodeSet *etree.Element) ([]byte, error) {
	resultNodeSet := nodeSet.Copy()
	propagateParentNamespacesTo(resultNodeSet)
	traverseAndTransformElement(
		resultNodeSet,
		new(superfluosNamespaceRemovingTransformer),
		new(lexicographicalSortingTransformer),
		new(commentRemovingTransformer),
	)
	return canonicalizeElementContentsAndConvertToBytes(resultNodeSet)
}

func propagateParentNamespacesTo(element *etree.Element) {
	allAncestorNamespaces := collectNamespacesFromAllAncestorsOf(element)
	for namespaceFullKey, namespace := range allAncestorNamespaces {
		if element.SelectAttr(namespaceFullKey) != nil {
			continue
		}
		element.Attr = append(element.Attr, namespace)
	}
}

func collectNamespacesFromAllAncestorsOf(element *etree.Element) map[string]etree.Attr {
	result := make(map[string]etree.Attr)
	for parent := element.Parent(); parent != nil; parent = parent.Parent() {
		for _, attr := range parent.Attr {
			if attr.Space != "xmlns" && attr.FullKey() != "xmlns" {
				continue
			}
			result[attr.FullKey()] = attr
		}
	}
	return result
}

func getOnlyNamespaceAttrFrom(element *etree.Element) map[string]etree.Attr {
	result := make(map[string]etree.Attr)
	for _, attr := range element.Attr {
		if attr.Space != "xmlns" && attr.FullKey() != "xmlns" {
			continue
		}
		result[attr.FullKey()] = attr
	}
	return result
}

func traverseAndTransformElement(element *etree.Element, transformers ...elementTransformer) {
	for _, transformer := range transformers {
		transformer.Transform(element)
	}
	for _, childElement := range element.ChildElements() {
		traverseAndTransformElement(childElement, transformers...)
	}
}

type elementTransformer interface {
	Transform(element *etree.Element)
}

type lexicographicalSortingTransformer struct{}

func (transformer *lexicographicalSortingTransformer) Transform(element *etree.Element) {
	sort.Sort(attributesByLexicographicalOrder(element.Attr))
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

type superfluosNamespaceRemovingTransformer struct{}

func (transformer *superfluosNamespaceRemovingTransformer) Transform(element *etree.Element) {
	parent := element.Parent()
	if parent == nil {
		return
	}
	namespacesOfParent := getOnlyNamespaceAttrFrom(parent)
	for _, attr := range element.Attr {
		if namespaceAttr, isRedeclared := namespacesOfParent[attr.FullKey()]; isRedeclared && namespaceAttr.Value == attr.Value {
			element.RemoveAttr(attr.FullKey())
		}
	}
}

type commentRemovingTransformer struct{}

func (transformer *commentRemovingTransformer) Transform(element *etree.Element) {
	for _, token := range element.Child {
		if _, isCommentNode := token.(*etree.Comment); isCommentNode {
			element.RemoveChild(token)
		}
	}
}

func canonicalizeElementContentsAndConvertToBytes(element *etree.Element) ([]byte, error) {
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
