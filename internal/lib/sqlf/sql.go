package sqlf

import "fmt"

func SQLFormat(query string, args ...interface{}) string {
	argsStr := make([]any, 0, len(args))
	for _, a := range args {
		if a == "" {
			a = nil
		}
		switch a.(type) {
		case nil:
			argsStr = append(argsStr, "NULL")
		case string:
			argsStr = append(argsStr, fmt.Sprintf("'%s'", a))
		default:
			argsStr = append(argsStr, fmt.Sprint(a))
		}
	}
	return fmt.Sprintf(query, argsStr...)
}
