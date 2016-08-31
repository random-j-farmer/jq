package jq

import (
	"fmt"
	"testing"
)

var simpleJSON = `
{"simple": {
	"string": "yes",
	"numStr": "2.75",
	"number": 3.14,
	"bool": true
},
"collection": [
	"yes",
	"2.75",
	3.14,
	true
]
}`

var table = []struct {
	keys     []interface{}
	expStr   string
	errStr   bool
	expInt   int64
	errInt   bool
	expFloat float64
	errFloat bool
	errSlice bool
	errMap   bool
}{
	{[]interface{}{"simple", "string"}, "yes", false, 0, true, 0.0, true, true, true},
	{[]interface{}{"simple", "numStr"}, "2.75", false, 2, false, 2.75, false, true, true},
	{[]interface{}{"simple", "number"}, "3.14", false, 3, false, 3.14, false, true, true},
	{[]interface{}{"simple", "bool"}, "true", false, 1, false, 0.0, true, true, true},
	{[]interface{}{"collection", 0}, "yes", false, 0, true, 0.0, true, true, true},
	{[]interface{}{"collection", 1}, "2.75", false, 2, false, 2.75, false, true, true},
	{[]interface{}{"collection", 2}, "3.14", false, 3, false, 3.14, false, true, true},
	{[]interface{}{"collection", 3}, "true", false, 1, false, 0.0, true, true, true},
	{[]interface{}{"doesNotExist"}, "", false, 0, false, 0.0, false, false, false},
	{[]interface{}{"collection", 666}, "", false, 0, false, 0.0, false, false, false},
	{[]interface{}{"collection"}, "", true, 0, true, 0.0, true, false, true},
	{[]interface{}{"simple"}, "", true, 0, true, 0.0, true, true, false},
}

func Test_String(t *testing.T) {
	q, err := Unmarshal([]byte(simpleJSON))
	if err != nil {
		t.Errorf("error parsing json: %v", err)
	}

	for _, inOut := range table {
		v, err := q.StringError(inOut.keys...)
		if inOut.errStr {
			if err == nil {
				t.Errorf("StringError(%v): would have expected error!", inOut.keys)
			}
		} else {
			if err != nil {
				t.Errorf("StringError(%v): %v", inOut.keys, err)
			}
		}
		if v != inOut.expStr {
			t.Errorf("StringError(%v) exp<>act:\n%v\n%v", inOut.keys, inOut.expStr, v)
		}
	}
}

func ExampleQuery_String() {
	q, _ := Unmarshal([]byte(`{"name": {"first": "Random", "last": "Farmer"}}`))
	fmt.Printf("%s %s", q.String("name", "first"), q.String("name", "last"))
	// Output: Random Farmer
}

func Test_Int(t *testing.T) {
	q, err := Unmarshal([]byte(simpleJSON))
	if err != nil {
		t.Errorf("error parsing json: %v", err)
	}

	for _, inOut := range table {
		v, err := q.Int64Error(inOut.keys...)
		if inOut.errInt {
			if err == nil {
				t.Errorf("IntError(%v): would have expected error!", inOut.keys)
			}
		} else {
			if err != nil {
				t.Errorf("IntError(%v): %v", inOut.keys, err)
			}
		}
		if v != inOut.expInt {
			t.Errorf("IntError(%v) exp<>act: %v<>%v", inOut.keys, inOut.expInt, v)
		}
	}

}

func Test_Float(t *testing.T) {
	q, err := Unmarshal([]byte(simpleJSON))
	if err != nil {
		t.Errorf("error parsing json: %v", err)
	}

	for _, inOut := range table {
		v, err := q.FloatError(inOut.keys...)
		if inOut.errFloat {
			if err == nil {
				t.Errorf("FloatError(%v): would have expected error!", inOut.keys)
			}
		} else {
			if err != nil {
				t.Errorf("FloatError(%v): %v", inOut.keys, err)
			}
		}
		if v != inOut.expFloat {
			t.Errorf("FloatError(%v) exp<>act: %v<>%v", inOut.keys, inOut.expFloat, v)
		}
	}
}

func ExampleQuery_Float() {
	q, _ := Unmarshal([]byte(`{"x": {"a": 17.5, "b": 2}}`))
	a, b, c := q.Float("x", "a"), q.Float("x", "b"), q.Float("x", "c")
	fmt.Printf("%g*%g*%g = %g\n", a, b, c, a*b*c)
	// Output: 17.5*2*0 = 0
}

func Test_Slice(t *testing.T) {
	q, err := Unmarshal([]byte(simpleJSON))
	if err != nil {
		t.Errorf("error parsing json: %v", err)
	}

	for _, inOut := range table {
		_, err := q.SliceError(inOut.keys...)
		if inOut.errSlice {
			if err == nil {
				t.Errorf("SliceError(%v): would have expected error!", inOut.keys)
			}
		} else {
			if err != nil {
				t.Errorf("SliceError(%v): %v", inOut.keys, err)
			}
		}
	}

	coll := q.Slice("collection")
	if len(coll) != 4 {
		t.Errorf("expected collection length 4: %d", len(coll))
	}
}

func Test_Map(t *testing.T) {
	q, err := Unmarshal([]byte(simpleJSON))
	if err != nil {
		t.Errorf("error parsing json: %v", err)
	}

	for _, inOut := range table {
		_, err := q.MapError(inOut.keys...)
		if inOut.errMap {
			if err == nil {
				t.Errorf("MapError(%v): would have expected error!", inOut.keys)
			}
		} else {
			if err != nil {
				t.Errorf("MapError(%v): %v", inOut.keys, err)
			}
		}
	}
}
