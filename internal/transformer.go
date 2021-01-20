package internal

import (
	"context"

	"github.com/beevik/etree"
)

type Transformer interface {
	Transform(ctx context.Context, element *etree.Element) (*etree.Element, error)
}
