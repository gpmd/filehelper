package filehelper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// sqlEscapeType uses Reflect to detect and handle each different type
// and escape it accordingly
func sqlEscapeType(value reflect.Value) string {
	var escaped string
	switch value.Kind() {
	case reflect.String:
		escaped = escapeString(value.String())
	case reflect.Slice:
		vals := make([]string, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			vals = append(vals, sqlEscapeType(value.Index(i)))
		}
		escaped = strings.Join(vals, ", ")
	case reflect.Int:
		escaped = strconv.FormatInt(value.Int(), 10)
	case reflect.Float32:
		escaped = strconv.FormatFloat(value.Float(), 'f', -1, 32)
	case reflect.Float64:
		escaped = strconv.FormatFloat(value.Float(), 'f', -1, 64)
	default:
		b, err := json.Marshal(value.Interface())
		if err != nil {
			panic(err)
		}
		escaped = sqlEscapeType(reflect.ValueOf(string(b)))
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
