package tracing

import (
	"strings"

	"google.golang.org/grpc/metadata"
)

// MDReaderWriter 读取元数据
type MDReaderWriter struct {
	metadata.MD
}

// ForeachKey 遍历所有的key并调用handler
func (c MDReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vs := range c.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// Set 设置key和val 到元数据中
func (c MDReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	c.MD[key] = append(c.MD[key], val)
}
