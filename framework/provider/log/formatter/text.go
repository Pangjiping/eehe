package formatter

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Pangjiping/eehe/framework/contract"
)

func TextFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	separator := "\t"

	prefix := Prefix(level)

	buffer.WriteString(prefix)
	buffer.WriteString(separator)

	ts := t.Format(time.RFC3339)
	buffer.WriteString(ts)
	buffer.WriteString(separator)

	buffer.WriteString("\"")
	buffer.WriteString(msg)
	buffer.WriteString("\"")
	buffer.WriteString(separator)

	buffer.WriteString(fmt.Sprint(fields))
	return buffer.Bytes(), nil
}
