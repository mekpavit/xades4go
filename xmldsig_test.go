package xades4go_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mekpavit/xades4go"
	"github.com/mekpavit/xades4go/etreeimpl"
)

func testXMLDSigSignatureValidator(t *testing.T, name string, xmldsigSignatureValidator xades4go.SignatureValidator) {
	type args struct {
		xmlBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    xades4go.ValidationResult
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(name+": "+tt.name, func(t *testing.T) {
			got, err := xmldsigSignatureValidator.Validate(tt.args.xmlBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Validate() result mismatch (-want+got):\n%s", diff)
			}
		})
	}
}

func Test_XMLDSigSignatureValidator(t *testing.T) {
	testXMLDSigSignatureValidator(t, "etreeimpl", xades4go.NewXMLDSigValidator(etreeimpl.NewSignedInfoFactory()))
}
