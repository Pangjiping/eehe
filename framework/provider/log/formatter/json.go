package formatter

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/Pangjiping/eehe/framework/contract"
	"github.com/pkg/errors"
)

func JsonFormatter(level contract.LogLevel, t time.Time, msg string, fields map[string]interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	if fields == nil {
		fields = map[string]interface{}{}
	}

	fields["msg"] = msg
	fields["level"] = level
	fields["timestamp"] = t.Format(time.RFC3339)
	content, err := json.Marshal(fields)
	if err != nil {
		return buffer.Bytes(), errors.Wrap(err, "json format error")
	}

	buffer.Write(content)
	return buffer.Bytes(), nil
}
