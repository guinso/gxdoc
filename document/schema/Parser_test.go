package schema

import (
	"strings"
	"testing"
)

func TestParseDocSchema(t *testing.T) {
	xml := `<documentschema name="invoice" revision="1">
		<validation script="js">return true</validation>
		<char name="jojo" charlen="32"></char>
		<text name="bobo" ></text>
		<integer  name="nini"></integer>
		<decimal name="price" precision="2" optional="false"></decimal>
		<date name="roro"></date>
		<file name="gogo" filter="pdf | zip"></file>
		<reference name="myLink" source="invoice.invNo"></reference>
		<items name="items" optional="true">
			<integer name="no"></integer>
			<char name="code" length="5"></char>
		</items>
	</documentschema>`

	docSchema, err := ParseFromXML(xml)

	if err != nil {
		t.Error("Fail to parse document schema: " + err.Error())
		return
	}

	if docSchema.Revision != 1 {
		t.Errorf("Expect revision is 1, but get %d instead", docSchema.Revision)
		return
	}

	if docSchema.Chars == nil {
		t.Error("Char field should has atleast one element")
		return
	}

	if charField := docSchema.Chars[0]; strings.Compare(charField.Name, "jojo") != 0 {
		t.Errorf("Expect char field name %s, but getting %s instead", "jojo", charField.Name)
		return
	}
}
