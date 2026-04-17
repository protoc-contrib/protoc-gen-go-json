package testpb_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/protoc-contrib/protoc-gen-go-json/internal/generator/testpb"
)

// TestRoundTrip exercises the generated MarshalJSON and UnmarshalJSON methods
// through encoding/json to verify they honor json.Marshaler / json.Unmarshaler
// contracts, including oneofs, maps, proto3 optional, and nested messages.
func TestRoundTrip(t *testing.T) {
	present := "present"
	empty := ""

	cases := []struct {
		name  string
		value proto.Message
	}{
		{
			name: "basic with oneof int",
			value: &testpb.Basic{
				A: "hello",
				B: &testpb.Basic_Int{Int: 42},
			},
		},
		{
			name: "nested",
			value: &testpb.Nested_Message{
				Basic: &testpb.Basic{
					A: "hello",
					B: &testpb.Basic_Str{Str: "world"},
				},
			},
		},
		{
			name: "optional present",
			value: &testpb.Basic{
				A: "hello",
				O: &present,
			},
		},
		{
			name: "optional empty",
			value: &testpb.Basic{
				A: "hello",
				O: &empty,
			},
		},
		{
			name: "with map",
			value: &testpb.Basic{
				A:   "hello",
				Map: map[string]string{"k": "v"},
			},
		},
		{
			name:  "with enum",
			value: &testpb.WithEnum{Kind: testpb.WithEnum_ONE},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			if len(bs) == 0 {
				t.Fatal("marshal produced empty output")
			}

			got := reflect.New(reflect.ValueOf(tc.value).Elem().Type()).Interface().(proto.Message)
			if err := json.Unmarshal(bs, got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if !proto.Equal(got, tc.value) {
				t.Fatalf("roundtrip mismatch:\n got:  %+v\n want: %+v", got, tc.value)
			}
		})
	}
}

// TestRoundTripEmbedded covers the embedded-by-value case — a Go struct that
// embeds a proto message. proto.Equal does not apply, so we rely on
// reflect.DeepEqual. Kept separate from TestRoundTrip so the proto-native
// cases use proto.Equal.
func TestRoundTripEmbedded(t *testing.T) {
	type basicWrapper struct{ testpb.Basic }

	want := &basicWrapper{
		Basic: testpb.Basic{
			A: "hello",
			B: &testpb.Basic_Int{Int: 42},
		},
	}

	bs, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := &basicWrapper{}
	if err := json.Unmarshal(bs, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !proto.Equal(&got.Basic, &want.Basic) {
		t.Fatalf("roundtrip mismatch:\n got:  %+v\n want: %+v", got, want)
	}
}
