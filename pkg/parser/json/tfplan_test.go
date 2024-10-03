package json

import (
	"testing"

	"github.com/Checkmarx/kics/v2/pkg/model"
	"github.com/stretchr/testify/require"
)

func TestJson_parseTFPlan(t *testing.T) {
	type args struct {
		doc model.Document
	}
	dummyValues := map[string]interface{}{
		"name": "Production DB",
		"size": 256,
	}

	dummyExpectedValues := map[string]interface{}{
		"name": "Production DB",
		"size": float64(256),
	}

	tests := []struct {
		name    string
		args    args
		want    model.Document
		wantErr bool
	}{
		{
			name: "test - parse as tfplan",
			args: args{
				doc: model.Document{
					"format_version":    "0.2",
					"terraform_version": "1.0.5",
					"variables":         map[string]interface{}{},
					"planned_values": map[string]interface{}{
						"root_module": map[string]interface{}{
							"resources": []map[string]interface{}{
								{
									"address":          "fakewebservices_database.prod_db",
									"mode":             "managed",
									"type":             "fakewebservices_database",
									"name":             "prod_db",
									"provider_name":    "registry.terraform.io/hashicorp/fakewebservices",
									"schema_version":   0,
									"values":           dummyValues,
									"sensitive_values": map[string]interface{}{},
								},
							},
						},
					},
					"resource_changes": []map[string]interface{}{},
					"configuration":    map[string]interface{}{},
				},
			},
			want: model.Document{
				"resource": map[string]interface{}{
					"fakewebservices_database": map[string]interface{}{
						"fakewebservices_database.prod_db": dummyExpectedValues,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "test - should not parse tfplan",
			args: args{
				doc: model.Document{
					"resource": map[string]interface{}{
						"name": "martin",
					},
				},
			},
			want:    model.Document{},
			wantErr: true,
		},
		{
			name: "test - child module is parsed",
			args: args{
				doc: model.Document{
					"format_version":    "0.2",
					"terraform_version": "1.0.5",
					"variables":         map[string]interface{}{},
					"planned_values": map[string]interface{}{
						"root_module": map[string]interface{}{
							"resources": []map[string]interface{}{
								{
									"address":       "fakewebservices_database.prod_db",
									"mode":          "managed",
									"type":          "fakewebservices_database",
									"name":          "prod_db",
									"provider_name": "registry.terraform.io/hashicorp/fakewebservices",
									"values":        dummyValues,
								},
							},
							"child_modules": []map[string]interface{}{
								{
									"resources": []map[string]interface{}{
										{
											"address":       "module.ilovetests.fakewebservices_database.prod_db",
											"mode":          "managed",
											"type":          "fakewebservices_database",
											"name":          "prod_db",
											"provider_name": "registry.terraform.io/hashicorp/fakewebservices",
											"values":        dummyValues,
										},
									},
									"child_modules": []map[string]interface{}{
										{
											"resources": []map[string]interface{}{
												{
													"address":       "module.ilovetests.module.ilovego.fakewebservices_database.prod_db",
													"mode":          "managed",
													"type":          "fakewebservices_database",
													"name":          "prod_db",
													"provider_name": "registry.terraform.io/hashicorp/fakewebservices",
													"values":        dummyValues,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: model.Document{
				"resource": map[string]interface{}{
					"fakewebservices_database": map[string]interface{}{
						"fakewebservices_database.prod_db":                                  dummyExpectedValues,
						"module.ilovetests.fakewebservices_database.prod_db":                dummyExpectedValues,
						"module.ilovetests.module.ilovego.fakewebservices_database.prod_db": dummyExpectedValues,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTFPlan(tt.args.doc)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
