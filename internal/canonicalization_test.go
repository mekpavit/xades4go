package internal

import (
	"context"
	"testing"

	"github.com/beevik/etree"
	"github.com/google/go-cmp/cmp"
)

func TestC14NTransformer_Transform(t *testing.T) {
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
			name: "when document subset is given, it should correctly propagate namespaces from ancestor",
			args: args{
				nodeSet: mustCreateElementFromString(`<?xml version="1.0" encoding="UTF-8"?>
<ietf:c14n11Xmllang xmlns:ietf="http://www.ietf.org"
xmlns:w3c="http://www.w3.org">
   <ietf:e1 xml:lang="EN">
      <ietf:e11>
         <ietf:e111 />
      </ietf:e11>
      <ietf:e12 at="2">
         <ietf:e121 />
      </ietf:e12>
   </ietf:e1>
   <ietf:e2 >
      <ietf:e21 />
   </ietf:e2>
</ietf:c14n11Xmllang>`).FindElement("e1"),
			},
			want: []byte(`<ietf:e1 xmlns:ietf="http://www.ietf.org" xmlns:w3c="http://www.w3.org" xml:lang="EN">
      <ietf:e11>
         <ietf:e111></ietf:e111>
      </ietf:e11>
      <ietf:e12 at="2">
         <ietf:e121></ietf:e121>
      </ietf:e12>
   </ietf:e1>`),
			wantErr: false,
		},
		{
			name: "when utf-8 character references are given, it should correctly handle them",
			args: args{
				nodeSet: mustCreateElementFromString(`<?xml version="1.0" encoding="ISO-8859-1"?>
<doc>&#169;</doc>`),
			},
			want:    []byte(`<doc>©</doc>`),
			wantErr: false,
		},
		{
			name: "when xml character referrences are given, it should be hanndled",
			args: args{
				nodeSet: mustCreateElementFromString(`<doc>
   <text>First line&#x0d;&#10;Second line</text>
   <value>&#x32;</value>
   <compute><![CDATA[value>"0" && value<"10" ?"valid":"error"]]></compute>
   <compute expr='value>"0" &amp;&amp; value&lt;"10" ?"valid":"error"'>valid</compute>
   <norm attr=' &apos;   &#x20;&#13;&#xa;&#9;   &apos; '/>
   <normNames attr='   A   &#x20;&#13;&#xa;&#9;   B   '/>
   <normId id=' &apos;   &#x20;&#13;&#xa;&#9;   &apos; '/>
</doc>`),
			},
			want: []byte(`<doc>
   <text>First line&#xD;
Second line</text>
   <value>2</value>
   <compute>value&gt;"0" &amp;&amp; value&lt;"10" ?"valid":"error"</compute>
   <compute expr="value>&quot;0&quot; &amp;&amp; value&lt;&quot;10&quot; ?&quot;valid&quot;:&quot;error&quot;">valid</compute>
   <norm attr=" '    &#xD;&#xA;&#x9;   ' "></norm>
   <normNames attr="   A    &#xD;&#xA;&#x9;   B   "></normNames>
   <normId id=" '    &#xD;&#xA;&#x9;   ' "></normId>
</doc>`),
			wantErr: false,
		},
		{
			name: "when unnecessary spaces and superflous namespaces are given, it should remove them",
			args: args{
				nodeSet: mustCreateElementFromString(`<doc>
   <e1   />
   <e2   ></e2>
   <e3   name = "elem3"   id="elem3"   />
   <e4   name="elem4"   id="elem4"   ></e4>
   <e5 a:attr="out" b:attr="sorted" attr2="all" attr="I'm"
      xmlns:b="http://www.ietf.org"
      xmlns:a="http://www.w3.org"
      xmlns="http://example.org"/>
   <e6 xmlns:a="http://www.w3.org">
      <e7 xmlns="http://www.ietf.org">
         <e8 xmlns:a="http://www.w3.org">
            <e9 xmlns:a="http://www.ietf.org" attr="default"/>
         </e8>
      </e7>
   </e6>
</doc>`),
			},
			want: []byte(`<doc>
   <e1></e1>
   <e2></e2>
   <e3 id="elem3" name="elem3"></e3>
   <e4 id="elem4" name="elem4"></e4>
   <e5 xmlns="http://example.org" xmlns:a="http://www.w3.org" xmlns:b="http://www.ietf.org" attr="I'm" attr2="all" b:attr="sorted" a:attr="out"></e5>
   <e6 xmlns:a="http://www.w3.org">
      <e7 xmlns="http://www.ietf.org">
         <e8>
            <e9 xmlns:a="http://www.ietf.org" attr="default"></e9>
         </e8>
      </e7>
   </e6>
</doc>`),
			wantErr: false,
		},
		{
			name: "when clean node set is given, it should done nothing",
			args: args{
				nodeSet: mustCreateElementFromString(`<doc>
   <clean>   </clean>
   <dirty>   A   B   </dirty>
   <mixed>
      A
      <clean>   </clean>
      B
      <dirty>   A   B   </dirty>
      C
   </mixed>
</doc>`),
			},
			want: []byte(`<doc>
   <clean>   </clean>
   <dirty>   A   B   </dirty>
   <mixed>
      A
      <clean>   </clean>
      B
      <dirty>   A   B   </dirty>
      C
   </mixed>
</doc>`),
			wantErr: false,
		},
		{
			name: "when node set contains xml declaration and comment, it should correctly remove them",
			args: args{
				nodeSet: mustCreateElementFromString(`<?xml version="1.0"?>


<doc>Hello, world!<!-- Comment 1 --></doc>


<!-- Comment 2 -->

<!-- Comment 3 -->`),
			},
			want:    []byte(`<doc>Hello, world!</doc>`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := &C14NTransformer{}
			gotElement, err := transformer.Transform(context.Background(), tt.args.nodeSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("C14N11Transformer.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := CompleteCanonicalization(gotElement)
			if err != nil {
				t.Errorf("CompleteCanonicalization returns error: %v", err)
				return
			}
			if diff := cmp.Diff(string(tt.want), string(got)); diff != "" {
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
