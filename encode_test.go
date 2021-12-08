package nbt

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

func TestEncoder_Encode_intArray(t *testing.T) {
	// Test marshal pure Int array
	v := [3]int32{0, -10, 3}
	out := []byte{tagInt32Array, 0x00, 0x00, 0, 0, 0, 3,
		0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0xff, 0xf6,
		0x00, 0x00, 0x00, 0x03,
	}
	if data, err := Marshal(v); err != nil {
		t.Error(err)
	} else if !bytes.Equal(data, out) {
		t.Errorf("output binary not right: get % 02x, want % 02x ", data, out)
	}

	// Test marshal in a struct
	v2 := struct {
		Ary [3]int32 `nbt:"ary"`
	}{[3]int32{0, -10, 3}}
	out = []byte{tagStruct, 0x00, 0x00,
		tagInt32Array, 0x00, 0x03, 'a', 'r', 'y', 0, 0, 0, 3,
		0x00, 0x00, 0x00, 0x00, // 0
		0xff, 0xff, 0xff, 0xf6, // -10
		0x00, 0x00, 0x00, 0x03, // 3
		tagEnd,
	}
	if data, err := Marshal(v2); err != nil {
		t.Error(err)
	} else if !bytes.Equal(data, out) {
		t.Errorf("output binary not right: get % 02x, want % 02x ", data, out)
	}
}

func TestEncoder_Encode_floatArray(t *testing.T) {
	// Test marshal pure Int array
	v := []float32{0.3, -100, float32(math.NaN())}
	out := []byte{tagSlice, 0x00, 0x00, tagFloat32, 0, 0, 0, 3,
		0x3e, 0x99, 0x99, 0x9a, // 0.3
		0xc2, 0xc8, 0x00, 0x00, // -100
		0x7f, 0xc0, 0x00, 0x00, // NaN
	}
	if data, err := Marshal(v); err != nil {
		t.Error(err)
	} else if !bytes.Equal(data, out) {
		t.Errorf("output binary not right: get % 02x, want % 02x ", data, out)
	}
}

func TestEncoder_Encode_string(t *testing.T) {
	v := "Test"
	out := []byte{tagString, 0x00, 0x00, 0, 4,
		'T', 'e', 's', 't'}

	if data, err := Marshal(v); err != nil {
		t.Error(err)
	} else if !bytes.Equal(data, out) {
		t.Errorf("output binary not right: got % 02x, want % 02x ", data, out)
	}
}

func TestEncoder_Encode_interfaceArray(t *testing.T) {
	type Struct1 struct {
		Val int32
	}

	type Struct2 struct {
		Val float32
	}

	tests := []struct {
		name string
		args []interface{}
		want []byte
	}{
		{
			name: "Two element interface array",
			args: []interface{}{Struct1{3}, Struct2{0.3}},
			want: []byte{
				tagSlice, 0x00, 0x00 /*no name*/, tagStruct, 0, 0, 0, 2,
				// 1st element
				tagInt32, 0x00, 0x03, 'V', 'a', 'l', 0x00, 0x00, 0x00, 0x03, // 3
				tagEnd,
				// 2nd element
				tagFloat32, 0x00, 0x03, 'V', 'a', 'l', 0x3e, 0x99, 0x99, 0x9a, // 0.3
				tagEnd,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.args)
			if err != nil {
				t.Error(err)
			} else if !bytes.Equal(data, tt.want) {
				t.Errorf("Marshal([]interface{}) got = % 02x, want % 02x", data, tt.want)
				return
			}
		})
	}
}

func TestEncoder_Encode_structArray(t *testing.T) {
	type Struct1 struct {
		Val int32
	}

	type Struct2 struct {
		T   int32
		Ele Struct1
	}

	type StructCont struct {
		V []Struct2
	}

	tests := []struct {
		name string
		args StructCont
		want []byte
	}{
		{
			name: "One element struct array",
			args: StructCont{[]Struct2{{3, Struct1{3}}, {-10, Struct1{-10}}}},
			want: []byte{
				tagStruct, 0x00, 0x00,
				tagSlice, 0x00, 0x01, 'V', tagStruct, 0, 0, 0, 2,
				// Struct2
				tagInt32, 0x00, 0x01, 'T', 0x00, 0x00, 0x00, 0x03,
				tagStruct, 0x00, 0x03, 'E', 'l', 'e',
				tagInt32, 0x00, 0x03, 'V', 'a', 'l', 0x00, 0x00, 0x00, 0x03, // 3
				tagEnd,
				tagEnd,
				// 2nd element
				tagInt32, 0x00, 0x01, 'T', 0xff, 0xff, 0xff, 0xf6,
				tagStruct, 0x00, 0x03, 'E', 'l', 'e',
				tagInt32, 0x00, 0x03, 'V', 'a', 'l', 0xff, 0xff, 0xff, 0xf6, // -10
				tagEnd,
				tagEnd,
				tagEnd,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.args)
			if err != nil {
				t.Error(err)
			} else if !bytes.Equal(data, tt.want) {
				t.Errorf("Marshal([]struct{}) got = % 02x, want % 02x", data, tt.want)
				return
			}
		})
	}
}

func TestEncoder_Encode_map(t *testing.T) {
	v := map[string][]int32{
		"Aaaaa":    {1, 2, 3, 4, 5},
		"Xi_Xi_Mi": {0, 0, 4, 7, 2},
	}

	b, err := Marshal(v)
	if err != nil {
		t.Fatal(err)
	}

	var data struct {
		Aaaaa []int32
		XXM   []int32 `nbt:"Xi_Xi_Mi"`
	}

	NewDecoder(bytes.NewReader(b)).Decode(&data)
	if !reflect.DeepEqual(data.Aaaaa, v["Aaaaa"]) {
		t.Fatalf("Marshal map error: got: %q, want %q", data.Aaaaa, v["Aaaaa"])
	}
	if !reflect.DeepEqual(data.XXM, v["Xi_Xi_Mi"]) {
		t.Fatalf("Marshal map error: got: %#v, want %#v", data.XXM, v["Xi_Xi_Mi"])
	}
}

func TestEncoder_Encode_interface(t *testing.T) {
	data := map[string]interface{}{
		"Key":   int32(12),
		"Value": "aaaaa",
	}
	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(data); err != nil {
		t.Fatalf("Encode error: %v", err)
	}

	var container struct {
		Key   int32
		Value string
	}
	if err := NewDecoder(&buf).Decode(&container); err != nil {
		t.Fatalf("Decode error: %v", err)
	}

	if container.Key != 12 || container.Value != "aaaaa" {
		t.Fatalf("want: (%v, %v), but got (%v, %v)", 12, "aaaaa", container.Key, container.Value)
	}
}
