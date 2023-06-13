package masker

import (
	"errors"
	"io"

	"github.com/goccy/go-json"
)

type MaskWriter struct {
	writer     io.Writer
	keysToMask map[string]struct{}
	maskValue  string
}

func NewMaskWriter(w io.Writer, keys []string, mask string) MaskWriter {
	keysToMask := make(map[string]struct{})
	for _, key := range keys {
		keysToMask[key] = struct{}{}
	}
	return MaskWriter{
		writer:     w,
		keysToMask: keysToMask,
		maskValue:  mask,
	}
}

func (mw *MaskWriter) Write(p []byte) (n int, err error) {
	var d interface{}
	if err := json.Unmarshal(p, &d); err != nil {
		return 0, err
	}
	m, ok := d.(map[string]interface{})
	if !ok {
		return 0, errors.New("data must be a map[string]interface{}")
	}
	m = mw.applyMaskRecursively(m)
	b, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}
	return mw.writer.Write(b)
}

func (mw *MaskWriter) applyMaskRecursively(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		switch v := v.(type) {
		case map[string]interface{}:
			m[k] = mw.applyMaskRecursively(v)
		case []interface{}:
			nv := make([]interface{}, len(v))
			for i, e := range v {
				if em, ok := e.(map[string]interface{}); ok {
					nv[i] = mw.applyMaskRecursively(em)
				} else {
					nv[i] = e
				}
			}
			m[k] = nv
		default:
			if _, ok := mw.keysToMask[k]; ok {
				m[k] = mw.maskValue
			}
		}
	}
	return m
}
