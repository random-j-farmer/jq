package jq

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Query wraps a parsed json object
type Query struct {
	json interface{}
}

// New creates a query object for an unmarshalled json object.
func New(json interface{}) *Query {
	return &Query{json: json}
}

// Unmarshal unmarshals from a json blob and creates a query struct.
func Unmarshal(b []byte) (*Query, error) {
	var parsed interface{}
	err := json.Unmarshal(b, &parsed)
	if err != nil {
		return nil, err
	}

	return &Query{json: parsed}, nil
}

// String extracts a string from the json document.
func (q *Query) String(args ...interface{}) string {
	s, _ := q.StringError(args...)
	return s
}

// Int is a convenience wrapper for Int64Error
func (q *Query) Int(args ...interface{}) int {
	i, _ := q.Int64Error(args...)
	return int(i)
}

// Int64 extracts an integer from the json document.
func (q *Query) Int64(args ...interface{}) int64 {
	i, _ := q.Int64Error(args...)
	return i
}

// UInt64 is a convenience wrapper for Int64Error
func (q *Query) UInt64(args ...interface{}) uint64 {
	i, _ := q.Int64Error(args...)
	return uint64(i)
}

// Float extracts a float from the json document.
func (q *Query) Float(args ...interface{}) float64 {
	f, _ := q.FloatError(args...)
	return f
}

// Slice extracts a collection from the json document.
func (q *Query) Slice(args ...interface{}) []interface{} {
	slice, _ := q.SliceError(args...)
	return slice
}

// Map extracs a map from the json document.
func (q *Query) Map(args ...interface{}) map[string]interface{} {
	m, _ := q.MapError(args...)
	return m
}

// StringError extracts a string from the json document.
func (q *Query) StringError(args ...interface{}) (string, error) {
	value, err := findChild(q.json, args)
	if err != nil {
		return "", err
	}

	switch value := value.(type) {
	case string:
		return value, nil
	case float64:
		return fmt.Sprintf("%v", value), nil
	case bool:
		return fmt.Sprintf("%v", value), nil
	case nil:
		return "", nil
	default:
		return "", fmt.Errorf("can not convert to string: %T %v", value, value)
	}

}

// Int64Error extracts an integer from the json document.
func (q *Query) Int64Error(args ...interface{}) (int64, error) {
	value, err := findChild(q.json, args)
	if err != nil {
		return 0, err
	}

	switch value := value.(type) {
	case string:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return 0, err
			}
			return int64(f), nil
		}
		return i, nil
	case float64:
		return int64(value), nil
	case bool:
		if value {
			return 1, nil
		}
		return 0, nil
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("can not convert to int: %T %v", value, value)
	}

}

// FloatError extracts a float from the json document.
func (q *Query) FloatError(args ...interface{}) (float64, error) {
	value, err := findChild(q.json, args)
	if err != nil {
		return 0.0, err
	}

	switch value := value.(type) {
	case float64:
		return value, nil
	case string:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0.0, err
		}
		return f, nil
	case nil:
		return 0.0, nil
	default:
		return 0.0, fmt.Errorf("can not convert to float: %T %v", value, value)
	}

}

// SliceError extracts a collection.
func (q *Query) SliceError(args ...interface{}) ([]interface{}, error) {
	value, err := findChild(q.json, args)
	if err != nil {
		return nil, err
	}

	switch value := value.(type) {
	case nil:
		return nil, nil
	case []interface{}:
		return value, nil
	default:
		return nil, fmt.Errorf("can not convert to []interface{}: %T %v", value, value)
	}
}

// MapError extracs a map
func (q *Query) MapError(args ...interface{}) (map[string]interface{}, error) {
	value, err := findChild(q.json, args)
	if err != nil {
		return nil, err
	}

	switch value := value.(type) {
	case nil:
		return nil, nil
	case map[string]interface{}:
		return value, nil
	default:
		return nil, fmt.Errorf("can not convert to map[string]interface{}: %T %v", value, value)
	}
}

func findChild(node interface{}, args []interface{}) (interface{}, error) {
	if node == nil || len(args) == 0 {
		return node, nil
	}
	arg := args[0]

	var child interface{}
	var err error

	switch arg := arg.(type) {
	case string:
		child, err = byString(node, arg)
	case int:
		child, err = byInt(node, arg)
	default:
		return nil, fmt.Errorf("can not index by: %T %v", arg, arg)
	}
	if err != nil {
		return nil, err
	}
	return findChild(child, args[1:])
}

func byString(node interface{}, arg string) (interface{}, error) {
	switch node := node.(type) {
	case map[string]interface{}:
		if child, ok := node[arg]; ok {
			return child, nil
		}
		// XXX: error?
		return nil, nil

	default:
		return nil, fmt.Errorf("unhandled json node: %T %v", node, node)
	}
}

func byInt(node interface{}, arg int) (interface{}, error) {
	switch node := node.(type) {
	case []interface{}:
		if arg >= len(node) {
			// XXX: or error?
			return nil, nil
		}
		return node[arg], nil

	default:
		return nil, fmt.Errorf("unhandled json node: %T %v", node, node)
	}
}
