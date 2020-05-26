package filehelper

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// sqlEscapeType uses Reflect to detect and handle each different type
// and escape it accordingly
func sqlEscapeType(value reflect.Value) string {
	sqlTypes := map[reflect.Kind]string {
		reflect.String: "VARCHAR",
		reflect.Int: "BIGINT",
		reflect.Float32: "FLOAT",
		reflect.Float64: "FLOAT",
	}

	var escaped string
	switch value.Kind() {
	case reflect.String:
		escaped = escapeString(value.String())
	case reflect.Slice:
		vals := make([]string, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			vals = append(vals, sqlEscapeType(value.Index(i)))
		}
		elemtype := "ARRAY"
		if sqlType, ok := sqlTypes[value.Type().Elem().Kind()]; ok {
			elemtype = sqlType
		}
		escaped = fmt.Sprintf("%s[%s]", elemtype, strings.Join(vals, ", "))
	case reflect.Int:
		escaped = strconv.FormatInt(value.Int(), 10)
	case reflect.Float32:
		escaped = strconv.FormatFloat(value.Float(), 'f', -1, 32)
	case reflect.Float64:
		escaped = strconv.FormatFloat(value.Float(), 'f', -1, 64)
	case reflect.Struct:
		vals := make([]string, 0, value.NumField())
		for i := 0; i < value.NumField(); i++ {
			vals = append(vals, sqlEscapeType(value.Field(i)))
		}
		escaped = strings.Join(vals, ", ")
	default:
		fmt.Println("Unsupported type")
	}
	return escaped
}

// escapeString, escapes unwanted characters from strings
// taken from https://gist.github.com/siddontang/8875771
func escapeString(source string) string {
	dest := make([]byte, 0, 2*len(source))
	var escape byte
	for i := 0; i < len(source); i++ {
		c := source[i]
		escape = 0
		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
			break
		case '\n': /* Must be escaped for logs */
			escape = 'n'
			break
		case '\r':
			escape = 'r'
			break
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		case '"': /* Better safe than sorry */
			escape = '"'
			break
		case '\032': /* This gives problems on Win32 */
			escape = 'Z'
		}
		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}
	return fmt.Sprintf("'%s'", dest)
}