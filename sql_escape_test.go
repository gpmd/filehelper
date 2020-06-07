package filehelper

import "testing"

type testSQLTemplateStruct struct {
	Template string
	Values   interface{}
	Result   string
}

func TestSQLEscape(t *testing.T) {
	tests := map[string]testSQLTemplateStruct{
		"string": testSQLTemplateStruct{
			Template: `select * from t where name = {{sql .}}`,
			Values:   "sur\"na'me",
			Result:   `select * from t where name = 'sur"na\'me'`,
		},
		"int": testSQLTemplateStruct{
			Template: `select * from t where answer = {{sql .}}`,
			Values:   42,
			Result:   `select * from t where answer = 42`,
		},
		"string list": testSQLTemplateStruct{
			Template: `select * from t where name in ({{sql .}})`,
			Values:   []string{"name\n", "sur\"na'me", "3"},
			Result:   `select * from t where name in ('name\n', 'sur"na\'me', '3')`,
		},
		"int list": testSQLTemplateStruct{
			Template: `select * from t where name in ({{sql .}})`,
			Values:   []int{1, 2, 3},
			Result:   `select * from t where name in (1, 2, 3)`,
		},
		"float list": testSQLTemplateStruct{
			Template: `select * from t where name in ({{sql .}})`,
			Values:   []float64{1.000001, 2.1, 3.0},
			Result:   `select * from t where name in (1.000001, 2.1, 3)`,
		},
		"json": testSQLTemplateStruct{
			Template: `select * from t where name = {{sql .}}`,
			Values:   map[string]float64{"a": 1.000001, "b": 2.1, "c": 3.0},
			Result:   `select * from t where name = '{"a":1.000001,"b":2.1,"c":3}'`,
		},
	}

	for name, test := range tests {
		res, err := Template(test.Template, test.Values)
		if err != nil {
			panic(err)
		}
		if res != test.Result {
			t.Errorf("%s: '%s' != '%s'", name, res, test.Result)
		}
	}

}
