package SQLBuilder

import (
	"strings"
	"testing"

	"github.com/guinso/gxschema"
)

func TestGenerateSQLTable(t *testing.T) {
	schema := gxschema.DxDoc{
		Name:     "invoice",
		Revision: 1,
		ID:       "733bee1b-f79a-4cb7-b675-842317b994b5",
		Items: []gxschema.DxItem{
			gxschema.DxInt{
				Name: "qty",
			},
			gxschema.DxStr{
				Name:           "inv no",
				EnableLenLimit: true,
				LenLimit:       6,
			},
			gxschema.DxFile{
				Name:    "attachment",
				IsArray: false,
			},
			gxschema.DxBool{
				Name: "isMandatory",
			},
			gxschema.DxDecimal{
				Name:      "total price",
				Precision: 2,
			},
			gxschema.DxSection{
				Name:    "items",
				IsArray: true,
				Items: []gxschema.DxItem{
					gxschema.DxStr{Name: "description"},
					gxschema.DxInt{Name: "qty"},
					gxschema.DxDecimal{Name: "unit price"},
				},
			},
		},
	}
	sqlStr, err := GenerateSQLTable(&schema)
	if err != nil {
		t.Error(err)
		return
	}

	expectedSQL := "CREATE TABLE `data_733bee1b-f79a-4cb7-b675-842317b994b5_r1`(\n" +
		"`id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
		"`qty` int(11) NOT NULL,\n" +
		"`inv no` char(6) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
		"`isMandatory` tinyint(1) NOT NULL,\n" +
		"`total price` decimal(11,2) NOT NULL,\n" +
		"PRIMARY KEY(`id`)\n" +
		") ENGINE=innodb DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;\n\n" +

		"CREATE TABLE `data_733bee1b-f79a-4cb7-b675-842317b994b5_r1_items`(\n" +
		"`id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
		"`parent_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
		"`description` text COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
		"`qty` int(11) NOT NULL,\n" +
		"`unit price` decimal(11,0) NOT NULL,\n" +
		"PRIMARY KEY(`id`),\n" +
		"CONSTRAINT `data_733bee1b-f79a-4cb7-b675-842317b994b5_r1_items_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `data_733bee1b-f79a-4cb7-b675-842317b994b5_r1` (`id`)\n" +
		") ENGINE=innodb DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;\n\n" +

		"CREATE TABLE `data_733bee1b-f79a-4cb7-b675-842317b994b5_r1_attachment`(\n" +
		"`id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
		"`parent_id` char(36) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
		"`filename` char(200) COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
		"`filepath` text COLLATE utf8mb4_unicode_ci NOT NULL,\n" +
		"PRIMARY KEY(`id`),\n" +
		"UNIQUE KEY `parent_id` (`parent_id`),\n" +
		"CONSTRAINT `data_733bee1b-f79a-4cb7-b675-842317b994b5_r1_attachment_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `data_733bee1b-f79a-4cb7-b675-842317b994b5_r1` (`id`)\n" +
		") ENGINE=innodb DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;\n\n"

	if strings.Compare(expectedSQL, sqlStr) != 0 {
		t.Errorf("output SQL not same as expected:\nExpected:\n%s\n\nOutput:\n%s\n====*", expectedSQL, sqlStr)
	}
	//t.Log(sqlStr)
	//t.Errorf("saja fail")
}
