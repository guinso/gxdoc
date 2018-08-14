package SQLBuilder

import (
	"fmt"
	"strings"

	"github.com/guinso/rdbmstool"

	"github.com/guinso/gxschema"
)

const dxID = "id"
const dxParentID = "parent_id"
const dxFileName = "filename"
const dxFilePath = "filepath"

//GenerateSQLTable generate SQL to create a set of datatables based on document schema
//Returns SQL string
func GenerateSQLTable(item *gxschema.DxDoc) (string, error) {
	//generate table
	subBuilders := []rdbmstool.TableBuilder{}

	builder := rdbmstool.NewTableBuilder()
	builder.TableName(fmt.Sprintf("data_%s_r%d", item.ID, item.Revision))
	builder.AddColumnChar(dxID, 36, false) //primary key
	builder.AddPrimaryKey(dxID)

	for _, subItem := range item.Items {

		if intItem, ok := subItem.(gxschema.DxInt); ok {
			arr, err := getIntBuilder(builder, &intItem, "", builder.GetTableName())
			if err != nil {
				return "", err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else if strItem, ok := subItem.(gxschema.DxStr); ok {
			arr, err := getStrBuilder(builder, &strItem, "", builder.GetTableName())
			if err != nil {
				return "", err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else if boolItem, ok := subItem.(gxschema.DxBool); ok {
			arr, err := getBoolBuilder(builder, &boolItem, "", builder.GetTableName())
			if err != nil {
				return "", err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else if decimalItem, ok := subItem.(gxschema.DxDecimal); ok {
			arr, err := getDecimalBuilder(builder, &decimalItem, "", builder.GetTableName())
			if err != nil {
				return "", err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else if fileItem, ok := subItem.(gxschema.DxFile); ok {
			arr, err := getFileBuilder(builder, &fileItem, "", builder.GetTableName())
			if err != nil {
				return "", err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else if sectionItem, ok := subItem.(gxschema.DxSection); ok {
			arr, err := getSectionBuilder(builder, &sectionItem, "", builder.GetTableName())
			if err != nil {
				return "", err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else {
			return "", fmt.Errorf(
				"unrecognize DxItem to build data table statement (SQL): %s",
				subItem.GetName())
		}
	}

	var tmpSQL, SQLStr string
	var err error
	tmpSQL, err = builder.SQL()
	if err != nil {
		return "", fmt.Errorf("failed to generate SQL statement: %s", err.Error())
	}
	SQLStr = tmpSQL + "\n\n"

	for i := len(subBuilders) - 1; i >= 0; i-- {
		tmpSQL, err = subBuilders[i].SQL()
		if err != nil {
			return "", fmt.Errorf("failed to generate SQL statement: %s", err.Error())
		}

		SQLStr += tmpSQL + "\n\n"
	}

	return SQLStr, nil
}

func getSubTableName(baseTableName string, path string, name string) string {
	if strings.Compare(path, "") == 0 {
		return fmt.Sprintf("%s_%s",
			baseTableName,
			strings.Replace(name, " ", "-", -1))
	}

	return fmt.Sprintf("%s_%s_%s",
		baseTableName,
		path,
		strings.Replace(name, " ", "-", -1))
}

//getSanatizeColName sanatize item's name by replacing white space to underscore
// func getSanatizeColName(name string) string {
// 	return strings.Replace(name, " ", "_", -1)
// }

func getIntBuilder(builder *rdbmstool.TableBuilder, item *gxschema.DxInt,
	path string, baseTableName string) ([]rdbmstool.TableBuilder, error) {
	if item.IsArray {
		//add sub table
		intBuilder := rdbmstool.NewTableBuilder()
		intBuilder.TableName(getSubTableName(baseTableName, path, item.Name))
		intBuilder.AddColumnInt(item.Name, 11, item.IsOptional)
		intBuilder.AddColumnChar(dxID, 36, false)
		intBuilder.AddColumnChar(dxParentID, 36, false)
		intBuilder.AddPrimaryKey(dxID)
		intBuilder.AddForeignKey(dxParentID, builder.GetTableName(), dxID)

		return []rdbmstool.TableBuilder{*intBuilder}, nil
	}

	builder.AddColumnInt(
		item.Name,
		11,
		item.IsOptional)

	return nil, nil
}

func getStrBuilder(builder *rdbmstool.TableBuilder, item *gxschema.DxStr,
	path string, baseTableName string) ([]rdbmstool.TableBuilder, error) {
	if item.IsArray {
		strBuilder := rdbmstool.NewTableBuilder()
		strBuilder.TableName(getSubTableName(baseTableName, path, item.Name))
		strBuilder.AddColumnChar(dxID, 36, false)
		strBuilder.AddColumnChar(dxParentID, 36, false)
		if item.EnableLenLimit {
			strBuilder.AddColumnChar(
				item.Name,
				item.LenLimit,
				item.IsOptional)
		} else {
			strBuilder.AddColumnText(
				item.Name,
				item.IsOptional)
		}
		strBuilder.AddPrimaryKey(dxID)
		strBuilder.AddForeignKey(dxParentID, builder.GetTableName(), dxID)

		return []rdbmstool.TableBuilder{*strBuilder}, nil
	}

	if item.EnableLenLimit {
		builder.AddColumnChar(
			item.Name,
			item.LenLimit,
			item.IsOptional)
	} else {
		builder.AddColumnText(
			item.Name,
			item.IsOptional)
	}

	return nil, nil
}

func getBoolBuilder(builder *rdbmstool.TableBuilder, item *gxschema.DxBool,
	path string, baseTableName string) ([]rdbmstool.TableBuilder, error) {
	if item.IsArray {
		boolBuilder := rdbmstool.NewTableBuilder()
		boolBuilder.TableName(getSubTableName(baseTableName, path, item.Name))
		boolBuilder.AddColumnChar(dxID, 36, false)
		boolBuilder.AddColumnChar(dxParentID, 36, false)
		boolBuilder.AddColumnBoolean(item.Name, item.IsOptional)
		boolBuilder.AddPrimaryKey(dxID)
		boolBuilder.AddForeignKey(dxParentID, builder.GetTableName(), dxID)

		return []rdbmstool.TableBuilder{*boolBuilder}, nil
	}

	builder.AddColumnBoolean(item.Name, item.IsOptional)

	return nil, nil
}

func getDecimalBuilder(builder *rdbmstool.TableBuilder, item *gxschema.DxDecimal,
	path string, baseTableName string) ([]rdbmstool.TableBuilder, error) {
	if item.IsArray {
		decimalBuilder := rdbmstool.NewTableBuilder()
		decimalBuilder.TableName(getSubTableName(baseTableName, path, item.Name))
		decimalBuilder.AddColumnChar(dxID, 36, false)
		decimalBuilder.AddColumnChar(dxParentID, 36, false)
		decimalBuilder.AddColumnDecimal(
			item.Name,
			11,
			item.Precision,
			item.IsOptional)
		decimalBuilder.AddPrimaryKey(dxID)
		decimalBuilder.AddForeignKey(dxParentID, builder.GetTableName(), dxID)

		return []rdbmstool.TableBuilder{*decimalBuilder}, nil
	}

	builder.AddColumnDecimal(
		item.Name,
		11,
		item.Precision,
		item.IsOptional)

	return nil, nil
}

func getFileBuilder(builder *rdbmstool.TableBuilder, item *gxschema.DxFile,
	path string, baseTableName string) ([]rdbmstool.TableBuilder, error) {

	fileBuilder := rdbmstool.NewTableBuilder()
	fileBuilder.TableName(getSubTableName(baseTableName, path, item.Name))
	fileBuilder.AddColumnChar(dxID, 36, false)
	fileBuilder.AddColumnChar(dxParentID, 36, false)
	fileBuilder.AddColumnChar(dxFileName, 200, false)
	fileBuilder.AddColumnText(dxFilePath, false)
	fileBuilder.AddPrimaryKey(dxID)
	fileBuilder.AddForeignKey(dxParentID, builder.GetTableName(), dxID)

	if !item.IsArray {
		fileBuilder.AddUniqueKey(dxParentID)
	}

	return []rdbmstool.TableBuilder{*fileBuilder}, nil
}

func getSectionBuilder(builder *rdbmstool.TableBuilder, item *gxschema.DxSection,
	path string, baseTableName string) ([]rdbmstool.TableBuilder, error) {

	subBuilders := []rdbmstool.TableBuilder{}

	fileBuilder := rdbmstool.NewTableBuilder()
	fileBuilder.TableName(getSubTableName(baseTableName, path, item.Name))
	fileBuilder.AddColumnChar(dxID, 36, false)
	fileBuilder.AddColumnChar(dxParentID, 36, false)
	fileBuilder.AddPrimaryKey(dxID)
	fileBuilder.AddForeignKey(dxParentID, builder.GetTableName(), dxID)

	//loop all items
	if strings.Compare(path, "") == 0 {
		path = strings.Replace(item.Name, " ", "-", -1)
	} else {
		path = path + "_" + strings.Replace(item.Name, " ", "-", -1)
	}
	for _, subItem := range item.Items {
		if intItem, ok := subItem.(gxschema.DxInt); ok {
			arr, err := getIntBuilder(fileBuilder, &intItem, path, fileBuilder.GetTableName())
			if err != nil {
				return nil, err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else if strItem, ok := subItem.(gxschema.DxStr); ok {
			arr, err := getStrBuilder(fileBuilder, &strItem, path, fileBuilder.GetTableName())
			if err != nil {
				return nil, err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else if boolItem, ok := subItem.(gxschema.DxBool); ok {
			arr, err := getBoolBuilder(fileBuilder, &boolItem, path, fileBuilder.GetTableName())
			if err != nil {
				return nil, err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else if decimalItem, ok := subItem.(gxschema.DxDecimal); ok {
			arr, err := getDecimalBuilder(fileBuilder, &decimalItem, path, fileBuilder.GetTableName())
			if err != nil {
				return nil, err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else if fileItem, ok := subItem.(gxschema.DxFile); ok {
			arr, err := getFileBuilder(fileBuilder, &fileItem, path, fileBuilder.GetTableName())
			if err != nil {
				return nil, err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else if sectionItem, ok := subItem.(gxschema.DxSection); ok {
			arr, err := getSectionBuilder(fileBuilder, &sectionItem, path, fileBuilder.GetTableName())
			if err != nil {
				return nil, err
			}

			for _, tmp := range arr {
				subBuilders = append(subBuilders, tmp)
			}
		} else {
			return nil, fmt.Errorf(
				"unrecognize DxItem to build data table statement (SQL): %s_%s",
				path, subItem.GetName())
		}
	}

	if !item.IsArray {
		fileBuilder.AddUniqueKey(dxParentID)
	}

	return append(subBuilders, *fileBuilder), nil
}
