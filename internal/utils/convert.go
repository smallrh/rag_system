package utils

import (
	"encoding/json"
	"fmt"
	"io"
)

// Convert 将 JSON 请求体转换为指定的 Go 实体
func Convert[T any](body io.Reader) (T, error) {
	var entity T

	if body == nil {
		return entity, fmt.Errorf("request body is empty")
	}

	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&entity); err != nil {
		return entity, fmt.Errorf("failed to decode request body: %w", err)
	}

	return entity, nil
}
