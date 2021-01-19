package internal

import (
	"context"
	"reflect"
	"testing"

	"github.com/beevik/etree"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestX(t *testing.T) {
	doc := etree.NewDocument()
	doc.ReadFromString(`<a xmlns:a="sssss"></a>`)
	t.Error(doc.Root().SelectAttr("xmlns:a"))
	doc.WriteSettings.CanonicalAttrVal = true
	doc.WriteSettings.CanonicalEndTags = true
	doc.WriteSettings.CanonicalText = true
	r, err := doc.WriteToString()
	require.NoError(t, err)
	t.Error(r)
}

func TestC14N11Transformer_Transform(t *testing.T) {
	type args struct {
		ctx     context.Context
		nodeSet *etree.Element
	}
	tests := []struct {
		name        string
		transformer *C14N11Transformer
		args        args
		want        []byte
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := &C14N11Transformer{}
			got, err := transformer.Transform(tt.args.ctx, tt.args.nodeSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("C14N11Transformer.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("C14N11Transformer.Transform() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lexicographicalSortingTransformer_Transform(t *testing.T) {
	type args struct {
		element *etree.Element
	}
	tests := []struct {
		name      string
		args      args
		wantBytes []byte
	}{
		{
			name: "when element contains unsorted attributes, it should sort the attributes in lexicographical order",
			args: args{
				element: mustCreateElementFromString(`<e5 a:attr="out" b:attr="sorted" attr2="all" attr="I'm"
      xmlns:b="http://www.ietf.org"
      xmlns:a="http://www.w3.org"
      xmlns="http://example.org"/>
`),
			},
			wantBytes: []byte(`<e5 xmlns="http://example.org" xmlns:a="http://www.w3.org" xmlns:b="http://www.ietf.org" attr="I'm" attr2="all" b:attr="sorted" a:attr="out"></e5>`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := &lexicographicalSortingTransformer{}
			transformer.Transform(tt.args.element)
			got, err := canonicalizeElementContentsAndConvertToBytes(tt.args.element)
			require.NoError(t, err)
			if diff := cmp.Diff(string(tt.wantBytes), string(got)); diff != "" {
				t.Errorf("Transform() result mismatch (-want+got):\n%s", diff)
			}
		})
	}
}

func Test_superfluosNamespaceRemovingTransformer_Transform(t *testing.T) {
	type args struct {
		element *etree.Element
	}
	tests := []struct {
		name      string
		args      args
		wantBytes []byte
	}{
		{
			name: "when the element contains both superfluos namespaces and non-superflous namespaces, it should remove only superflous namespaces",
			args: args{
				element: mustCreateElementFromString(`<a xmlns:a="https://example.com/a" xmlns:b="https://example.com/b">
<b xmlns:a="https://b.example.com/a" xmlns:b="https://example.com/b" xmlns:c="https://b.example.com/c"></b>
</a>`).FindElement("b"),
			},
			wantBytes: []byte(`<b xmlns:a="https://b.example.com/a" xmlns:c="https://b.example.com/c"></b>`),
		},
		{
			name: "when the element does not contains superfluos namespaces, it do nothing",
			args: args{
				element: mustCreateElementFromString(`<a xmlns:a="https://example.com/a" xmlns:b="https://example.com/b">
<b xmlns:a="https://b.example.com/a" xmlns:c="https://b.example.com/c"></b>
</a>`).FindElement("b"),
			},
			wantBytes: []byte(`<b xmlns:a="https://b.example.com/a" xmlns:c="https://b.example.com/c"></b>`),
		},
		{
			name: "when the element contains superfluos namespaces, it should remove them",
			args: args{
				element: mustCreateElementFromString(`<a xmlns:a="https://example.com/a" xmlns:b="https://example.com/b">
<b xmlns:a="https://example.com/a" xmlns:b="https://example.com/b"></b>
</a>`).FindElement("b"),
			},
			wantBytes: []byte(`<b></b>`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := &superfluosNamespaceRemovingTransformer{}
			transformer.Transform(tt.args.element)
			got, err := canonicalizeElementContentsAndConvertToBytes(tt.args.element)
			require.NoError(t, err)
			if diff := cmp.Diff(string(tt.wantBytes), string(got)); diff != "" {
				t.Errorf("Transform() result mismatch (-want+got):\n%s", diff)
			}
		})
	}
}

func Test_commentRemovingTransformer_Transform(t *testing.T) {
	type args struct {
		element *etree.Element
	}
	tests := []struct {
		name      string
		args      args
		wantBytes []byte
	}{
		{
			name: "when the element does not contain comment node in its children, it should do nothing",
			args: args{
				element: mustCreateElementFromString(`<a>
<b></b>

<c></c>
</a>`),
			},
			wantBytes: []byte(`<a>
<b></b>

<c></c>
</a>`),
		},
		{
			name: "when the element contains comment node in its children, it should remove that comment node",
			args: args{
				element: mustCreateElementFromString(`<a>
<b></b>
<!-- some comment
-->
<c></c>
</a>`),
			},
			wantBytes: []byte(`<a>
<b></b>

<c></c>
</a>`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := &commentRemovingTransformer{}
			transformer.Transform(tt.args.element)
			got, err := canonicalizeElementContentsAndConvertToBytes(tt.args.element)
			require.NoError(t, err)
			if diff := cmp.Diff(string(tt.wantBytes), string(got)); diff != "" {
				t.Errorf("Transform() result mismatch (-want+got):\n%s", diff)
			}
		})
	}
}

func mustCreateElementFromString(xmlContent string) *etree.Element {
	doc := etree.NewDocument()
	doc.ReadFromString(xmlContent)
	return doc.Root()
}
