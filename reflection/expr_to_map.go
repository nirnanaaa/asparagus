package reflection

import (
	"regexp"
	"strings"
)

var exprRex = regexp.MustCompile(`(\w+)=("(.*?)"|\w+)`)

// ExprToMap turns an expression (key=value) into a map of strings.
//
func ExprToMap(expr string) map[string]interface{} {
	data := exprRex.FindAllStringSubmatch(expr, -1)

	res := make((map[string]interface{}))
	for _, kv := range data {
		k := kv[1]
		v := kv[2]
		// TODO: maybe adjust regex to skip quotes upon extraction?
		res[k] = strings.Replace(v, `"`, "", -1)
	}
	return res
}
