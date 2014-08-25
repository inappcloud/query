package query

import (
	"bytes"
	"fmt"
	"strings"
)

func replacePlaceholders(sql string) string {
	buf := &bytes.Buffer{}
	for i := 1; ; i++ {
		p := strings.Index(sql, "?")
		if p == -1 {
			break
		}

		buf.WriteString(sql[:p])
		fmt.Fprintf(buf, "$%d", i)
		sql = sql[p+1:]
	}

	buf.WriteString(sql)
	return buf.String()
}
