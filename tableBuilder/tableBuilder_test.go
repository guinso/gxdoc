package tableBuilder

import (
	"testing"

	"github.com/guinso/gxschema"
)

func TestGenerateSQLBuilder(t *testing.T) {
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
	sqlStr, err := GenerateSQLBuilder(&schema)
	if err != nil {
		t.Error(err)
	}

	t.Log(sqlStr)
	t.Errorf("saja fail")
}

/*
func TestGenerateSQLBuilder(t *testing.T) {
	type args struct {
		item *gxschema.DxDoc
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSQLBuilder(tt.args.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSQLBuilder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateSQLBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
