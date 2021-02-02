package etreeimpl

import (
	"testing"

	"github.com/beevik/etree"
	"github.com/google/go-cmp/cmp"
	"github.com/mekpavit/xades4go"
)

func TestEtreeEnvelopedSignatureTransformer_Transform(t *testing.T) {
	type args struct {
		nodeSet *etree.Element
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "when element does not contains Signature element, it should do nothing",
			args: args{
				nodeSet: mustCreateElementFromString(`<a><aa><aaa></aaa></aa><ab><aba></aba><abb></abb></ab></a>`),
			},
			want:    []byte(`<a><aa><aaa></aaa></aa><ab><aba></aba><abb></abb></ab></a>`),
			wantErr: false,
		},
		{
			name: "when element contains Signature element, it should remove it",
			args: args{
				nodeSet: mustCreateElementFromString(`<a><aa><aaa></aaa></aa><ab><aba></aba><abb><ds:Signature></ds:Signature></abb></ab></a>`),
			},
			want:    []byte(`<a><aa><aaa></aaa></aa><ab><aba></aba><abb></abb></ab></a>`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := &envelopedSignatureTransformer{}
			gotElement, err := transformer.Transform(xades4go.XML{IsOctetStream: false, NodeSet: tt.args.nodeSet})
			if (err != nil) != tt.wantErr {
				t.Errorf("EnvelopedSignatureTransformer.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotEtreeElement := gotElement.NodeSet.(*etree.Element)
			got, err := completeCanonicalization(gotEtreeElement)
			if err != nil {
				t.Errorf("CompleteCanonicalization returns error: %v", err)
			}
			if diff := cmp.Diff(string(tt.want), string(got)); diff != "" {
				t.Errorf("EnvelopedSignatureTransformer.Transform() result mistmatch (-want+got):\n%s", diff)
			}
		})
	}
}
