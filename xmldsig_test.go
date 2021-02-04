package xades4go_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mekpavit/xades4go"
	"github.com/mekpavit/xades4go/etreeimpl"
)

func Test_XMLDSigSignatureValidator(t *testing.T) {
	runTestOfXMLDSigSignatureValidator(t, "etreeimpl", xades4go.NewXMLDSigSignatureValidator(etreeimpl.NewSignedInfoFactory()))
}

func runTestOfXMLDSigSignatureValidator(t *testing.T, name string, xmldsigSignatureValidator xades4go.SignatureValidator) {
	type args struct {
		xmlBytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    xades4go.ValidationResult
		wantErr bool
	}{
		{
			name: "When KeyInfo is X509Data, transform algorithms are enveloped signature and canonical XML 1.0, it should pass the valition",
			args: args{
				xmlBytes: []byte(`<rsm:TaxInvoice_CrossIndustryInvoice xmlns:rsm="urn:etda:uncefact:data:standard:TaxInvoice_CrossIndustryInvoice:2" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:etda:uncefact:data:standard:TaxInvoice_CrossIndustryInvoice:2">
    <rsm:ExchangedDocumentContext xmlns:ram="urn:etda:uncefact:data:standard:TaxInvoice_ReusableAggregateBusinessInformationEntity:2">
        <ram:GuidelineSpecifiedDocumentContextParameter>
            <ram:ID schemeAgencyID="ETDA" schemeVersionID="v2.0">ER3-2560</ram:ID>
        </ram:GuidelineSpecifiedDocumentContextParameter>
    </rsm:ExchangedDocumentContext>
    <rsm:ExchangedDocument xmlns:ram="urn:etda:uncefact:data:standard:TaxInvoice_ReusableAggregateBusinessInformationEntity:2">
        <ram:ID>INV01</ram:ID>
        <ram:Name>ใบกำกับภาษี</ram:Name>
        <ram:TypeCode>388</ram:TypeCode>
        <ram:IssueDateTime>2017-12-19T00:00:00.000</ram:IssueDateTime>


        <ram:CreationDateTime>2017-12-19T00:00:00.000</ram:CreationDateTime>

    </rsm:ExchangedDocument>
    <rsm:SupplyChainTradeTransaction xmlns:ram="urn:etda:uncefact:data:standard:TaxInvoice_ReusableAggregateBusinessInformationEntity:2">
        <ram:ApplicableHeaderTradeAgreement>
            <ram:SellerTradeParty>
                <ram:Name>บริษัท ขยันหมั่นเพียร จำกัด</ram:Name>
                <ram:SpecifiedTaxRegistration>
                    <ram:ID schemeID="TXID">123456789012300000</ram:ID>
                </ram:SpecifiedTaxRegistration>
                <ram:DefinedTradeContact>

                    <ram:EmailURIUniversalCommunication>
                        <ram:URIID>natueal@example.com</ram:URIID>
                    </ram:EmailURIUniversalCommunication>


                    <ram:TelephoneUniversalCommunication>
                        <ram:CompleteNumber>+66-81234567</ram:CompleteNumber>
                    </ram:TelephoneUniversalCommunication>


                </ram:DefinedTradeContact>


                <ram:PostalTradeAddress>
                    <ram:PostcodeCode>71180</ram:PostcodeCode>
                    <ram:LineOne>สำนักงานอุทยานแห่งชาติ </ram:LineOne>

                    <ram:CityName>7107</ram:CityName>
                    <ram:CitySubDivisionName>710705</ram:CitySubDivisionName>
                    <ram:CountryID schemeID="3166-1 alpha-2">TH</ram:CountryID>
                    <ram:CountrySubDivisionID>71</ram:CountrySubDivisionID>
                    <ram:BuildingNumber>777/777 </ram:BuildingNumber>
                </ram:PostalTradeAddress>
            </ram:SellerTradeParty>
            <ram:BuyerTradeParty>
                <ram:Name>บริษัททำดีจำกัด</ram:Name>
                <ram:SpecifiedTaxRegistration>
                    <ram:ID schemeID="TXID">222222222222200000</ram:ID>
                </ram:SpecifiedTaxRegistration>
                <ram:DefinedTradeContact>

                    <ram:EmailURIUniversalCommunication>
                        <ram:URIID>patiw@example.com</ram:URIID>
                    </ram:EmailURIUniversalCommunication>


                    <ram:TelephoneUniversalCommunication>
                        <ram:CompleteNumber>+66-97778889</ram:CompleteNumber>
                    </ram:TelephoneUniversalCommunication>


                </ram:DefinedTradeContact>


                <ram:PostalTradeAddress>
                    <ram:PostcodeCode>11344</ram:PostcodeCode>
                    <ram:LineOne>หาดทุ่งวัวแล่น</ram:LineOne>

                    <ram:CityName>8603</ram:CityName>
                    <ram:CitySubDivisionName>860303</ram:CitySubDivisionName>
                    <ram:CountryID schemeID="3166-1 alpha-2">TH</ram:CountryID>
                    <ram:CountrySubDivisionID>86</ram:CountrySubDivisionID>
                    <ram:BuildingNumber>77/79</ram:BuildingNumber>
                </ram:PostalTradeAddress>
            </ram:BuyerTradeParty>

        </ram:ApplicableHeaderTradeAgreement>
        <ram:ApplicableHeaderTradeDelivery>
            <ram:ShipToTradeParty>
                <ram:DefinedTradeContact>
                    <ram:PersonName>สมพร ใจงาม</ram:PersonName>
                </ram:DefinedTradeContact>
            </ram:ShipToTradeParty>
        </ram:ApplicableHeaderTradeDelivery>
        <ram:ApplicableHeaderTradeSettlement>
            <ram:InvoiceCurrencyCode listID="ISO 4217 3A">THB</ram:InvoiceCurrencyCode>
            <ram:ApplicableTradeTax>
                <ram:TypeCode>VAT</ram:TypeCode>
                <ram:CalculatedRate>7</ram:CalculatedRate>
                <ram:BasisAmount>9999</ram:BasisAmount>
                <ram:CalculatedAmount>699.93</ram:CalculatedAmount>
            </ram:ApplicableTradeTax>


            <ram:SpecifiedTradeSettlementHeaderMonetarySummation>
                <ram:LineTotalAmount>9999</ram:LineTotalAmount>
                <ram:TaxBasisTotalAmount>9999</ram:TaxBasisTotalAmount>
                <ram:TaxTotalAmount>699.93</ram:TaxTotalAmount>
                <ram:GrandTotalAmount>10698.93</ram:GrandTotalAmount>
            </ram:SpecifiedTradeSettlementHeaderMonetarySummation>
        </ram:ApplicableHeaderTradeSettlement>
        <ram:IncludedSupplyChainTradeLineItem>
            <ram:AssociatedDocumentLineDocument>
                <ram:LineID>1</ram:LineID>
            </ram:AssociatedDocumentLineDocument>
            <ram:SpecifiedTradeProduct>


                <ram:Name>สินค้าทดสอบ</ram:Name>
            </ram:SpecifiedTradeProduct>
            <ram:SpecifiedLineTradeAgreement>
                <ram:GrossPriceProductTradePrice>
                    <ram:ChargeAmount>9999</ram:ChargeAmount>
                </ram:GrossPriceProductTradePrice>
            </ram:SpecifiedLineTradeAgreement>
            <ram:SpecifiedLineTradeDelivery>
                <ram:BilledQuantity unitCode="AS">1</ram:BilledQuantity>
            </ram:SpecifiedLineTradeDelivery>
            <ram:SpecifiedLineTradeSettlement>

                <ram:SpecifiedTradeSettlementLineMonetarySummation>
                    <ram:NetLineTotalAmount>9999</ram:NetLineTotalAmount>
                    <ram:NetIncludingTaxesLineTotalAmount>10698.93</ram:NetIncludingTaxesLineTotalAmount>
                </ram:SpecifiedTradeSettlementLineMonetarySummation>
            </ram:SpecifiedLineTradeSettlement>
        </ram:IncludedSupplyChainTradeLineItem>
    </rsm:SupplyChainTradeTransaction>
<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#" Id="xmldsig-5b38fead-4352-464f-b3b3-3f6cd5c9fbf9"><ds:SignedInfo><ds:CanonicalizationMethod Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/><ds:SignatureMethod Algorithm="http://www.w3.org/2001/04/xmldsig-more#rsa-sha512"/><ds:Reference Id="xmldsig-5b38fead-4352-464f-b3b3-3f6cd5c9fbf9-ref0" URI=""><ds:Transforms><ds:Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/><ds:Transform Algorithm="http://www.w3.org/TR/2001/REC-xml-c14n-20010315"/></ds:Transforms><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha512"/><ds:DigestValue>y2/Zx52P9Ck3r1/Rb8Xn516CcuT8i4I57hPKWk++6rv8kmk0Azd+intm2yNgtVyKdHaRt/qAL4YWmgHu91Z7tQ==</ds:DigestValue></ds:Reference><ds:Reference Type="http://uri.etsi.org/01903#SignedProperties" URI="#xmldsig-5b38fead-4352-464f-b3b3-3f6cd5c9fbf9-signedprops"><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha512"/><ds:DigestValue>u/ejCCgofcQ7jpaZuyc6RAkd4CuEugPVFx31aFJ3iIEoRh4ZxDkryGHmmPvrQXAp/nEMp4GkcedrQLHJT7kZEA==</ds:DigestValue></ds:Reference></ds:SignedInfo><ds:SignatureValue Id="xmldsig-5b38fead-4352-464f-b3b3-3f6cd5c9fbf9-sigvalue">BOxF2QGxUzpNcP5YcJ6IIMLLWWmrqEscAPE7a+yr/x3kgJnwMPWm4D3ae0F5zifA+OeZGhQjPZ9ctIZYssVUZtJNIBrHuTpFndZvL9H07//HWjJUOi7Gv8qeCRw1FsuSdlTew9TrucNH6zfCm5KQIGCkH2nplV3oNhvqewA6cjYW1nJmnAyfaaDcD0x7xl6dmgXQ9xv573eCFKP72GRzhwxr36llqLvaoJMHUGkq59wCYc7oMgj5d+vG5fSA6BXsfXPWxmp1gyeoa7UIPT62Pvy79RjH9WtmiQcUcm4/iA+fX0jCOXhKeJEOz15fhTCHQnWp4CYNjVBOzoh2xOhlEg==</ds:SignatureValue><ds:KeyInfo><ds:X509Data><ds:X509Certificate>MIIFrjCCA5agAwIBAgIIew6XEhSld4EwDQYJKoZIhvcNAQELBQAwgbUxCzAJBgNVBAYTAnRoMT0wOwYDVQQKDDRNaW5pc3RyeSBvZiBJbmZvcm1hdGlvbiBhbmQgQ29tbXVuaWNhdGlvbiBUZWNobm9sb2d5MUkwRwYDVQQLDEBFbGVjdHJvbmljIFRyYW5zYWN0aW9ucyBEZXZlbG9wbWVudCBBZ2VuY3kgKFB1YmxpYyBPcmdhbml6YXRpb24pMRwwGgYDVQQDDBNUZURBIENBIGZvciBUZXN0aW5nMB4XDTE5MDcwNDA1MDY1MFoXDTIyMDcwNDA1MDY1MFowZTELMAkGA1UEBhMCVEgxMzAxBgNVBAoMKkVsZWN0cm9uaWMgVHJhbnNhY3Rpb25zIERldmVsb3BtZW50IEFnZW5jeTEhMB8GA1UEAwwYQ29kZSBTaWduaW5nIENlcnRpZmljYXRlMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoVWH2z88nV0+GuxpyHKZjNWxTv7syQqwL+FiF2KErwdI9rRSKFz5GyOhB5N6Nzjh0AflcfUGrOa+AbjNi+5MGUC3uL2ugo7jMXx2Rwp90aHhGU8jhE/Dx6pdUiQSd6ZLyCTYzIrb4okgDRzsJTe3pFfnM0ScspiU2GMRCMmQgPYWob6BFPgzqIcYK99f82CENg3PFlm4bUWTgvVgF0TTevjLvH8Dx8LvnORW05Hk+jHWbPQuHtarmubUowFv9N1EboBEbFAPxhOp67vctzuoDfOtLaC0unfkXBxzSpnpKg7ZDSEbT4C6hDMqXaAGZxRWtE6ggVQfjzICOs8ZhV+zkwIDAQABo4IBDzCCAQswVQYIKwYBBQUHAQEESTBHMEUGCCsGAQUFBzAChjlodHRwOi8vcmVwby10ZXN0LnRlZGEudGgvY2VydC9UZURBQ0Fmb3JUZXN0aW5nLmNhY2VydC5jcnQwHQYDVR0OBBYEFBnP8q2USQRJ1oTlTds9irSSivl7MAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgwFoAUw7M9c+QxW38JHlzBhP71WNfouZQwQgYDVR0fBDswOTA3oDWgM4YxaHR0cDovL3JlcG8tdGVzdC50ZWRhLnRoL2NybC9UZURBQ0Fmb3JUZXN0aW5nLmNybDALBgNVHQ8EBAMCBPAwEwYDVR0lBAwwCgYIKwYBBQUHAwMwDQYJKoZIhvcNAQELBQADggIBACrFj0Ee4paSEBzmskqyLatVvbnDfUUfDMMkQrSGcD2l2lNaorAtcBZeVJTRMt+doJTNPwpAFbW3rbbAAX+PKAn5M8F2dcj0W/Q6dIw1pQyuRJIBgJ7BwXq7fwbEV3C1AUV3EXTGND4hz7LYRqCIuLi6ODdT3/HBQlEBQtNhKLBBciE81mWKvaQ1g/hAbPZOSDW7WBEw8Kjj1vbPS0lviar8TurRwbwDlYMk6NzpSGPJYUrxjYw54ZJx/1QngKGK6wsZiV0sj5JbbfxjTwWOhEl2LdulQJ8KNZv+ajQMZqtEeAreAHLyGSG6xgOpPV9aHP9LDTR/d5qi3JB5fwMOvEsWWvzoKzvilR6WO3hYL8qQi/Y4C7oYMkjxVBAALXi2PH4cZSA26SkR2gHQ8FMO1o+StqkBBjkrtdyhvr+PxijFSh25T3rLlAPBDCALSUPRdLg848k07CleGBzDDETNsFnUhiZXCzD6TWEKqdMVItzXCuCe+bCX8/wvsVC48chMdxjVHLR3P8csyK+tPS+Te9ipsI3ZgIoDWilNJhKMyaQbmI+zHFzBVVE9cVMkFsGOWh0lKscQA3k1CnhvfIsppyz/ZK6sj7/7Q5+is/4ay5vGhgSPXVhN0kCW5u+esGouVPJLMfqvwhh4V1a+9sQtsRDugqPhzk00/DrI0byUjl65</ds:X509Certificate><ds:X509Certificate>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</ds:X509Certificate></ds:X509Data></ds:KeyInfo><ds:Object><xades:QualifyingProperties xmlns:xades="http://uri.etsi.org/01903/v1.3.2#" Target="#xmldsig-5b38fead-4352-464f-b3b3-3f6cd5c9fbf9"><xades:SignedProperties Id="xmldsig-5b38fead-4352-464f-b3b3-3f6cd5c9fbf9-signedprops"><xades:SignedSignatureProperties><xades:SigningTime>2020-08-17T18:27:35+07:00</xades:SigningTime><xades:SigningCertificate><xades:Cert><xades:CertDigest><ds:DigestMethod Algorithm="http://www.w3.org/2001/04/xmlenc#sha512"/><ds:DigestValue>1r8x/T+ReH+ehYxWyrRkILA2U0whmsLYLrewFjuhiTGSLQX7RchTPF0eZJvhPihlaIGa7xBMzS58iILPkh71MA==</ds:DigestValue></xades:CertDigest><xades:IssuerSerial><ds:X509IssuerName>C=TH,O=Ministry of Information and Communication Technology,CN=TeDA CA for Testing</ds:X509IssuerName><ds:X509SerialNumber>8867190820250679169</ds:X509SerialNumber></xades:IssuerSerial></xades:Cert></xades:SigningCertificate></xades:SignedSignatureProperties></xades:SignedProperties></xades:QualifyingProperties></ds:Object></ds:Signature></rsm:TaxInvoice_CrossIndustryInvoice>`),
			},
			want: xades4go.ValidationResult{
				ReferenceValidationResults: []xades4go.ReferenceValidationResult{
					{
						IsValid:              true,
						GeneratedDigestValue: `y2/Zx52P9Ck3r1/Rb8Xn516CcuT8i4I57hPKWk++6rv8kmk0Azd+intm2yNgtVyKdHaRt/qAL4YWmgHu91Z7tQ==`,
						DigestValue:          `y2/Zx52P9Ck3r1/Rb8Xn516CcuT8i4I57hPKWk++6rv8kmk0Azd+intm2yNgtVyKdHaRt/qAL4YWmgHu91Z7tQ==`,
					},
					{
						IsValid:              true,
						GeneratedDigestValue: `u/ejCCgofcQ7jpaZuyc6RAkd4CuEugPVFx31aFJ3iIEoRh4ZxDkryGHmmPvrQXAp/nEMp4GkcedrQLHJT7kZEA==`,
						DigestValue:          `u/ejCCgofcQ7jpaZuyc6RAkd4CuEugPVFx31aFJ3iIEoRh4ZxDkryGHmmPvrQXAp/nEMp4GkcedrQLHJT7kZEA==`,
					},
				},
				IsSignatureValid: true,
			},
			wantErr: false,
		},
	}
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
