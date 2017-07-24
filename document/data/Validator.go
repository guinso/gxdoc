package data

import (
	"fmt"
	"gxdoc/document/schema"
	"regexp"
	"time"

	"github.com/beevik/etree"
)

//ValidateXML is to validate input XML document is compliance with predefined document schema
func ValidateXML(rawXML string, schema *schema.DocSchema) error {
	//TODO: check input XML is well form or not...

	//load XML into an instance
	xmlDoc := etree.NewDocument()
	if readErr := xmlDoc.ReadFromString(rawXML); readErr != nil {
		return readErr
	}

	root := xmlDoc.FindElement(schema.Name)
	if root == nil {
		return fmt.Errorf("XML root is not match. Expected root is %s", schema.Name)
	}

	//validate chars
	if len(schema.Chars) > 0 {
		for _, charSchema := range schema.Chars {
			if charErr := validateChar(root, &charSchema); charErr != nil {
				return charErr
			}
		}
	}

	//validate texts
	if len(schema.Texts) > 0 {
		for _, textSchema := range schema.Texts {
			if textErr := validateText(root, &textSchema); textErr != nil {
				return textErr
			}
		}
	}

	//validate integers
	if len(schema.Texts) > 0 {
		for _, intSchema := range schema.Integers {
			if intErr := validateInteger(root, &intSchema); intErr != nil {
				return intErr
			}
		}
	}

	//validate decimal
	if len(schema.Decimals) > 0 {
		for _, decimalSchema := range schema.Decimals {
			if decimalErr := validateDecimal(root, &decimalSchema); decimalErr != nil {
				return decimalErr
			}
		}
	}

	//validate date
	if len(schema.Dates) > 0 {
		for _, dateSchema := range schema.Dates {
			if dateErr := validateDate(root, &dateSchema); dateErr != nil {
				return dateErr
			}
		}
	}

	//validate file
	if len(schema.Files) > 0 {
		for _, fileSchema := range schema.Files {
			if fileErr := validateFile(root, &fileSchema); fileErr != nil {
				return fileErr
			}
		}
	}

	//validate items
	if len(schema.Items) > 0 {
		for _, itemSchema := range schema.Items {
			if itemsErr := validateItem(root, &itemSchema); itemsErr != nil {
				return itemsErr
			}
		}
	}

	return nil
}

func validateItem(root *etree.Element, schema *schema.ItemsSchema) error {
	elements := root.FindElements(schema.Name)
	eleCount := len(elements)

	if eleCount == 0 && schema.Optional == false {
		return fmt.Errorf("<%s> cannot be empty", schema.Name)
	}

	//check each item
	for _, element := range elements {
		if itemErr := validateSubItem(element, schema); itemErr != nil {
			return itemErr
		}
	}

	return nil
}

func validateSubItem(element *etree.Element, schema *schema.ItemsSchema) error {
	//validate chars
	if len(schema.Chars) > 0 {
		for _, charSchema := range schema.Chars {
			if charErr := validateChar(element, &charSchema); charErr != nil {
				return charErr
			}
		}
	}

	//validate texts
	if len(schema.Texts) > 0 {
		for _, textSchema := range schema.Texts {
			if textErr := validateText(element, &textSchema); textErr != nil {
				return textErr
			}
		}
	}

	//validate integers
	if len(schema.Texts) > 0 {
		for _, intSchema := range schema.Integers {
			if intErr := validateInteger(element, &intSchema); intErr != nil {
				return intErr
			}
		}
	}

	//validate decimal
	if len(schema.Decimals) > 0 {
		for _, decimalSchema := range schema.Decimals {
			if decimalErr := validateDecimal(element, &decimalSchema); decimalErr != nil {
				return decimalErr
			}
		}
	}

	//validate date
	if len(schema.Dates) > 0 {
		for _, dateSchema := range schema.Dates {
			if dateErr := validateDate(element, &dateSchema); dateErr != nil {
				return dateErr
			}
		}
	}

	//validate file
	if len(schema.Files) > 0 {
		for _, fileSchema := range schema.Files {
			if fileErr := validateFile(element, &fileSchema); fileErr != nil {
				return fileErr
			}
		}
	}

	return nil
}

func validateFile(root *etree.Element, schema *schema.FileSchema) error {
	//validate element count
	if countErr := validateElementCount(root, schema.Name, schema.Optional); countErr != nil {
		return countErr
	}

	file := root.FindElement(schema.Name)
	if file != nil && len(file.Text()) > 0 {
		//inner text only store file name
		//wheres physical file path is parsed from multi-part form

		len := len(file.ChildElements())
		if len > 0 {
			return fmt.Errorf("<%s> File type not allow sub element; found %d",
				schema.Name, len)
		}
	}

	return nil
}

func validateDate(root *etree.Element, schema *schema.DateSchema) error {
	//validate element count
	if countErr := validateElementCount(root, schema.Name, schema.Optional); countErr != nil {
		return countErr
	}

	date := root.FindElement(schema.Name)
	if date != nil && len(date.Text()) > 0 {
		_, parseErr := time.Parse("2006-01-02", date.Text())
		if parseErr != nil {
			return parseErr
		}
	}

	//validate no sub elements
	if date != nil {
		if len(date.ChildElements()) > 0 {
			return fmt.Errorf("<%s> Date type is not allowed to have sub element; content: %s",
				schema.Name, date.Text())
		}
	}

	return nil
}

func validateDecimal(root *etree.Element, schema *schema.DecimalSchema) error {
	//validate element count
	if countErr := validateElementCount(root, schema.Name, schema.Optional); countErr != nil {
		return countErr
	}

	decimal := root.FindElement(schema.Name)
	if decimal != nil && len(decimal.Text()) > 0 {
		//validate text pattern is in correct decimal format
		pattern := fmt.Sprintf("^(([-+]?[1-9][0-9]*|0))(\\.[0-9]{1,%d})?$", schema.Precision)
		if isMatch, _ := regexp.MatchString(pattern, decimal.Text()); !isMatch {
			return fmt.Errorf("Element <%s> content is not decimal, found: %s",
				schema.Name, decimal.Text())
		}
	}

	//validate no sub elements
	if decimal != nil {
		if len(decimal.ChildElements()) > 0 {
			return fmt.Errorf("<%s> Decimal type is not allowed to have sub element; content: %s",
				schema.Name, decimal.Text())
		}
	}

	return nil
}

func validateInteger(root *etree.Element, schema *schema.IntegerSchema) error {
	//validate element count
	if countErr := validateElementCount(root, schema.Name, schema.Optional); countErr != nil {
		return countErr
	}

	//validate content is integer or not
	integer := root.FindElement(schema.Name)
	if integer != nil && len(integer.Text()) > 0 {
		if isMatch, _ := regexp.MatchString("^(([-+]?[1-9][0-9]*|0)$)", integer.Text()); !isMatch {
			return fmt.Errorf("Element <%s> content is not integer, found: %s",
				schema.Name, integer.Text())
		}
	}

	//validate no sub elements
	if integer != nil {
		if len(integer.ChildElements()) > 0 {
			return fmt.Errorf("<%s> Integer type is not allowed to have sub element; content: %s",
				schema.Name, integer.Text())
		}
	}

	return nil
}

func validateChar(root *etree.Element, schema *schema.CharSchema) error {
	//validate element count
	if countErr := validateElementCount(root, schema.Name, schema.Optional); countErr != nil {
		return countErr
	}

	//validate char length
	if schema.Length > 0 {
		char := root.FindElement(schema.Name)
		if char != nil && !isContentEmptyAllowable(char, schema.Optional) {
			content := char.Text()
			length := len(content)
			if length != schema.Length {
				return fmt.Errorf("Character length is not match with <%s> setting; "+
					"expected: %d, actual: %d, content: %s",
					schema.Name, schema.Length, len(content), content)
			}
		}
	}

	//apply regular expression validation
	if len(schema.Regex) > 0 {
		tag := root.FindElement(schema.Name)
		if tag != nil && !isContentEmptyAllowable(tag, schema.Optional) {
			isMatch, _ := regexp.MatchString(schema.Regex, tag.Text())
			if !isMatch {
				return fmt.Errorf("Regular expression not match with <%s> setting; "+
					"Content: %s, Pattern: %s",
					schema.Name, tag.Text(), schema.Regex)
			}
		}
	}

	//validate no sub elements
	char := root.FindElement(schema.Name)
	if char != nil {
		if len(char.ChildElements()) > 0 {
			return fmt.Errorf("<%s> Char type is not allowed to have sub element; content: %s",
				schema.Name, char.Text())
		}
	}

	return nil
}

func validateText(root *etree.Element, schema *schema.TextSchema) error {
	//validate element count
	if countErr := validateElementCount(root, schema.Name, schema.Optional); countErr != nil {
		return countErr
	}

	//apply regular expression validation
	if len(schema.Regex) > 0 {
		tag := root.FindElement(schema.Name)
		if tag != nil && !isContentEmptyAllowable(tag, schema.Optional) {
			isMatch, _ := regexp.MatchString(schema.Regex, tag.Text())
			if !isMatch {
				return fmt.Errorf("Regular expression not match with <%s> setting; "+
					"Content: %s, Pattern: %s",
					schema.Name, tag.Text(), schema.Regex)
			}
		}
	}

	text := root.FindElement(schema.Name)
	//validate no sub elements
	if text != nil {
		if len(text.ChildElements()) > 0 {
			return fmt.Errorf("<%s> Text type is not allowed to have sub element; content: %s",
				schema.Name, text.Text())
		}
	}

	return nil
}

func isContentEmptyAllowable(element *etree.Element, optional bool) bool {
	return len(element.Text()) == 0 && optional
}

func validateElementCount(root *etree.Element, tagName string, optional bool) error {
	chars := root.FindElements(tagName)

	//validate quantity
	count := len(chars)
	qtyOk := (count == 1) || (count == 0 && optional == true)
	if !qtyOk {
		if count == 0 && optional == false {
			return fmt.Errorf("XML element %s not found", tagName)
		}

		return fmt.Errorf("XML element <%s> is allow to occur one time but found %d",
			tagName, count)
	}

	return nil
}

//ValidateJSON to validate input Json document is compliance with predefined document schema
func ValidateJSON(rawJSON string, schema *schema.DocSchema) error {
	return nil
}
