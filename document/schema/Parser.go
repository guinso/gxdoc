package schema

import "encoding/xml"

type DocSchema struct {
	XMLName    xml.Name          `xml:"documentschema"`
	Name       string            `xml:"name,attr"`
	Revision   int               `xml:"revision,attr"`
	Items      []ItemsSchema     `xml:"items"`
	Validation *ValidationSchema `xml:"validation"`
	Chars      []CharSchema      `xml:"char"`
	Texts      []TextSchema      `xml:"text"`
	Integers   []IntegerSchema   `xml:"integer"`
	Decimals   []DecimalSchema   `xml:"decimal"`
	Dates      []DateSchema      `xml:"date"`
	Files      []FileSchema      `xml:"file"`
	References []ReferenceSchema `xml:"reference"`
}

type ItemsSchema struct {
	XMLName    xml.Name          `xml:"items"`
	Name       string            `xml:"name,attr"`
	Optional   bool              `xml:"optional,attr"`
	Chars      []CharSchema      `xml:"char"`
	Texts      []TextSchema      `xml:"text"`
	Integers   []IntegerSchema   `xml:"integer"`
	Decimals   []DecimalSchema   `xml:"decimal"`
	Dates      []DateSchema      `xml:"date"`
	Files      []FileSchema      `xml:"file"`
	References []ReferenceSchema `xml:"reference"`
}

type ValidationSchema struct {
	XMLName    xml.Name `xml:"validation"`
	ScriptType string   `xml:"script,attr"`
	Script     string   `xml:"innerxml"`
}

type CharSchema struct {
	XMLName  xml.Name `xml:"char"`
	Length   int      `xml:"charlen,attr"`
	Name     string   `xml:"name,attr"`
	Optional bool     `xml:"optional,attr"`
	Regex    string   `xml:"regex,attr"`
}

type TextSchema struct {
	XMLName  xml.Name `xml:"text"`
	Name     string   `xml:"name,attr"`
	Optional bool     `xml:"optional,attr"`
	Regex    string   `xml:"regex,attr"`
}

type IntegerSchema struct {
	XMLName  xml.Name `xml:"integer"`
	Name     string   `xml:"name,attr"`
	Optional bool     `xml:"optional,attr"`
}

type DecimalSchema struct {
	XMLName   xml.Name `xml:"decimal"`
	Name      string   `xml:"name.attr"`
	Precision int      `xml:"precision,attr"`
	Optional  bool     `xml:"optional,attr"`
}

type DateSchema struct {
	XMLName  xml.Name `xml:"date"`
	Name     string   `xml:"name,attr"`
	Optional bool     `xml:"optional,attr"`
}

type FileSchema struct {
	XMLName  xml.Name `xml:"file"`
	Name     string   `xml:"name,attr"`
	Filter   string   `xml:"filter,attr"`
	Optional bool     `xml:"optional,attr"`
}

//TODO: need to study further in programming use case(s)
type ReferenceSchema struct {
	XMLName  xml.Name `xml:"reference"`
	Name     string   `xml:"name,attr"`
	Source   string   `xml:"source,attr"`
	Optional bool     `xml:"optional,attr"`
}

// ParseFromXML is to convert xml to DocSchema data object
func ParseFromXML(xmlString string) (*DocSchema, error) {
	result := DocSchema{}

	parseErr := xml.Unmarshal([]byte(xmlString), &result)

	if parseErr != nil {
		return nil, parseErr
	}

	return &result, nil
}

//UnmarshalXML parse Char xml tag in custom method
func (ri *CharSchema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type cField CharSchema //prevent recursion
	item := cField{Optional: false}

	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}

	*ri = (CharSchema)(item)
	return nil
}

//UnmarshalXML parse Text xml tag in custom method
func (ri *TextSchema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type cField TextSchema //prevent recursion
	item := cField{Optional: false}

	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}

	*ri = (TextSchema)(item)
	return nil
}

//UnmarshalXML parse Integer xml tag in custom method
func (ri *IntegerSchema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type cField IntegerSchema //prevent recursion
	item := cField{Optional: false}

	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}

	*ri = (IntegerSchema)(item)
	return nil
}

//UnmarshalXML parse Decimal xml tag in custom method
func (ri *DecimalSchema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type cField DecimalSchema //prevent recursion
	item := cField{Optional: false}

	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}

	*ri = (DecimalSchema)(item)
	return nil
}

//UnmarshalXML parse Date xml tag in custom method
func (ri *DateSchema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type cField DateSchema //prevent recursion
	item := cField{Optional: false}

	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}

	*ri = (DateSchema)(item)
	return nil
}

//UnmarshalXML parse File xml tag in custom method
func (ri *FileSchema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type cField FileSchema //prevent recursion
	item := cField{Optional: false}

	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}

	*ri = (FileSchema)(item)
	return nil
}

//UnmarshalXML parse Reference xml tag in custom method
func (ri *ReferenceSchema) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type cField ReferenceSchema //prevent recursion
	item := cField{Optional: false}

	if err := d.DecodeElement(&item, &start); err != nil {
		return err
	}

	*ri = (ReferenceSchema)(item)
	return nil
}
