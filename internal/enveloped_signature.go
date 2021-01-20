package internal

import (
	"context"

	"github.com/beevik/etree"
)

type EnvelopedSignatureTransformer struct{}

func (transformer *EnvelopedSignatureTransformer) Transform(ctx context.Context, nodeSet *etree.Element) (*etree.Element, error) {
	resultNodeSet := nodeSet.Copy()
	signatureElement := resultNodeSet.FindElement("//Signature")
	if signatureElement == nil {
		return resultNodeSet, nil
	}
	parentOfSignatureElement := signatureElement.Parent()
	if parentOfSignatureElement == nil {
		return resultNodeSet, nil
	}
	parentOfSignatureElement.RemoveChild(signatureElement)
	return resultNodeSet, nil
}
