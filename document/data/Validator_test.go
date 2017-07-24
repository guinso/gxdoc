package data

import (
	"testing"

	"github.com/guinso/gxdoc/document/schema"
)

func TestValidateXML(t *testing.T) {

	schemaDoc := schema.DocSchema{
		Name:     "invoice",
		Revision: 1,
		Chars: []schema.CharSchema{
			schema.CharSchema{
				Name:     "docNo",
				Optional: false,
				Length:   6,
				Regex:    "^inv[0-9]{3}$"}},
		Texts: []schema.TextSchema{
			schema.TextSchema{
				Name:     "remark",
				Optional: false,
				Regex:    "^(ko){2} [0-9]+$"}},
		Integers: []schema.IntegerSchema{
			schema.IntegerSchema{
				Name:     "qty",
				Optional: false}},
		Decimals: []schema.DecimalSchema{
			schema.DecimalSchema{
				Name:      "price",
				Optional:  false,
				Precision: 2}},
		Dates: []schema.DateSchema{
			schema.DateSchema{
				Name:     "date",
				Optional: false}},
		Files: []schema.FileSchema{
			schema.FileSchema{
				Name:     "attachment",
				Optional: false,
				Filter:   "*.pdf"}},
		Items: []schema.ItemsSchema{
			schema.ItemsSchema{
				Name:     "order",
				Optional: false,
				Integers: []schema.IntegerSchema{
					schema.IntegerSchema{
						Name:     "no",
						Optional: false},
					schema.IntegerSchema{
						Name:     "quantity",
						Optional: false}},
				Texts: []schema.TextSchema{
					schema.TextSchema{
						Name:     "description",
						Optional: false}},
				Decimals: []schema.DecimalSchema{
					schema.DecimalSchema{
						Name:      "cost",
						Optional:  false,
						Precision: 2}}}}}

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, &schemaDoc, t)

	positiveTest(&schemaDoc, t)

	negativeTest(&schemaDoc, t)
}

func positiveTest(schemaDoc *schema.DocSchema, t *testing.T) {
	//==== docNo ====
	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv123</docNo>
		<remark>koko 06653</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv999</docNo>
		<remark>koko 06653</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv000</docNo>
		<remark>koko 06653</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	//==== qty ====
	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>+123</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>0</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>12</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	//==== price ====
	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>43</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>+7.92</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>-5.7</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>0.20</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>-9.3</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	//==== remark ====
	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 06653</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 0</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 193757362</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	//==== orders ====
	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 193757362</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)

	runXMLValidationPositive(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 193757362</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`, schemaDoc, t)
}

func negativeTest(schemaDoc *schema.DocSchema, t *testing.T) {

	//==== docNo ====
	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv0001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <docNo> fail with 'inv0001'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv 001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <docNo> fail with 'inv 001'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>PO004</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <docNo> fail with 'PO004'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001A</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <docNo> fail with 'inv001A'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>
			<fofo>123</fofo>
			<koko>asd</koko>
		</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <docNo> fail with '<fofo>123</fofo><koko>asd</koko>'")

	//==== remark ====
	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>kokoko 123</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <remark> fail with 'kokoko 123'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>kokoA 123</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <remark> fail with 'kokoA 123'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko123</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <remark> fail with 'koko123'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark> koko 123</remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <remark> fail with ' koko 123'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123 </remark>
		<qty>-8900</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <remark> fail with 'koko 123 '")

	//==== qty ====
	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>012</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <qty> fail with '012'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty> 12</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <qty> fail with ' 12'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>20 </qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <qty> fail with '20 '")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>3.6</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <qty> fail with '3.6'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>2.5e+45</qty>
		<price>95.30</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <qty> fail with '2.5e+45'")

	//==== price ====
	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>05.60</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <price> fail with '05.60'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>3.5e-4</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <price> fail with '3.5e-4'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>-04.20</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <price> fail with '-04.20'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>1.45678</price>
		<date>2017-04-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <price> fail with '1.45678'")

	//==== date ====
	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>1.45</price>
		<date>2017-04-32</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <date> fail with '2017-04-32'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>1.45</price>
		<date>2017-23-02</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <date> fail with '2017-23-02'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>1.45</price>
		<date>23-04-2015</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <date> fail with '23-04-2015'")

	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>1.45</price>
		<date>2017-4-14</date>
		<attachment>info.pdf</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <date> fail with '2017-4-14'")

	//==== file ====
	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>1.45</price>
		<date>2017-04-14</date>
		<attachment>
			<koko>qwe</koko>
			<gogo>123</gogo>
		</attachment>
		<order>
			<no>1</no>
			<description>torch light</description>
			<cost>6.70</cost>
			<quantity>1</quantity>
		</order>
		<order>
			<no>2</no>
			<description>pen</description>
			<cost>3.20</cost>
			<quantity>4</quantity>
		</order>
		<order>
			<no>3</no>
			<description>A4 paper</description>
			<cost>20.00</cost>
			<quantity>1</quantity>
		</order>
	</invoice>`,
		schemaDoc, t, "Expect <file> fail with '<koko>qwe</koko><gogo>123</gogo>'")

	//==== order ====
	runXMLValidationNegative(`<?xml version="1.0"?>
	<invoice>
		<docNo>inv001</docNo>
		<remark>koko 123</remark>
		<qty>-8900</qty>
		<price>1.45</price>
		<date>2017-04-14</date>
		<attachment>jerr.pdf</attachment>
	</invoice>`,
		schemaDoc, t, "Expect <order> fail with empty element")
}

func runXMLValidationPositive(xml string, schema *schema.DocSchema, t *testing.T) {
	err := ValidateXML(xml, schema)

	if err != nil {
		t.Error(err.Error())
	}
}

func runXMLValidationNegative(xml string, schema *schema.DocSchema, t *testing.T, errMsg string) {
	err := ValidateXML(xml, schema)

	if err == nil {
		t.Error(errMsg)
	}
}
