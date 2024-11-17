// Copyright 2024 Oscar Pernia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package json

import (
	"maps"
	"slices"
	"strings"
	"testing"
)

func TestAllowTypeMismatchDecode(t *testing.T) {
	type T struct {
		String  string         `json:"string"`
		Int     int            `json:"int"`
		Float64 float64        `json:"float64"`
		Object  map[string]any `json:"object"`
		Slice   []any          `json:"slice"`
	}

	equalsT := func(a, b T) bool {
		return a.String == b.String &&
			a.Int == b.Int &&
			a.Float64 == b.Float64 &&
			maps.Equal(a.Object, b.Object) &&
			slices.Equal(a.Slice, b.Slice)
	}

	baseT := T{
		String:  "test",
		Int:     123,
		Float64: 123.123,
		Object: map[string]any{
			"foo": "bar",
		},
		Slice: []any{
			1.0, 2.0, 3.0, // .0 appended because they get unmarshaled as float64
		},
	}

	testCases := []struct {
		CaseName

		input          string
		expectedT      func() T
		expectingError bool
	}{
		{
			CaseName: Name("MismatchedType_WholeStruct"),
			input:    `"test"`, // JSON string
			expectedT: func() T {
				return T{} // zero-value
			},
			expectingError: false,
		},
		{
			CaseName: Name("MismatchedType_String"),
			input: `
				{
					"string": 123,
					"int": 123,
					"float64": 123.123,
					"object": {"foo": "bar"},
					"slice": [
						1,
						2,
						3
					]
				}
			`,
			expectedT: func() T {
				ret := baseT
				ret.String = "" // zero-value
				return ret
			},
			expectingError: false,
		},
		{
			CaseName: Name("MismatchedType_Int"),
			input: `
				{
					"string": "test",
					"int": "MISMATCHED_TYPE",
					"float64": 123.123,
					"object": {"foo": "bar"},
					"slice": [
						1,
						2,
						3
					]
				}
			`,
			expectedT: func() T {
				ret := baseT
				ret.Int = 0
				return ret
			},
			expectingError: false,
		},
		{
			CaseName: Name("MismatchedType_Int_GotFloat64"),
			input: `
				{
					"string": "test",
					"int": 123.123,
					"float64": 123.123,
					"object": {"foo": "bar"},
					"slice": [
						1,
						2,
						3
					]
				}
			`,
			expectedT: func() T {
				ret := baseT
				ret.Int = 0
				return ret
			},
		},
		{
			CaseName: Name("MismatchedType_Float64"),
			input: `
				{
					"string": "test",
					"int": 123,
					"float64": "MISMATCHED_TYPE",
					"object": {"foo": "bar"},
					"slice": [
						1,
						2,
						3
					]
				}
			`,
			expectedT: func() T {
				ret := baseT
				ret.Float64 = 0
				return ret
			},
			expectingError: false,
		},
		{
			CaseName: Name("MismatchedType_Object"),
			input: `
				{
					"string": "test",
					"int": 123,
					"float64": 123.123,
					"object": "MISMATCHED_TYPE",
					"slice": [
						1,
						2,
						3
					]
				}
			`,
			expectedT: func() T {
				ret := baseT
				ret.Object = nil
				return ret
			},
			expectingError: false,
		},
		{
			CaseName: Name("MismatchedType_Slice"),
			input: `
				{
					"string": "test",
					"int": 123,
					"float64": 123.123,
					"object": {"foo": "bar"},
					"slice": "MISMATCHED_TYPE"
				}
			`,
			expectedT: func() T {
				ret := baseT
				ret.Slice = nil
				return ret
			},
			expectingError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			dec := NewDecoder(strings.NewReader(tc.input))

			dec.AllowTypeMismatch() // Function to test

			var gotT T

			err := dec.Decode(&gotT)
			if err != nil && tc.expectingError {
				return
			} else if err != nil && !tc.expectingError {
				t.Fatalf("expected (Decoder).Decode() to not return an error, got:\n\t%v", err)
			} else if err == nil && tc.expectingError {
				t.Fatalf("expected (Decoder).Decode() to return an error, got nil")
			}

			expectedT := tc.expectedT()

			if !equalsT(gotT, expectedT) {
				t.Fatalf("expected:\n\t%v\n\tgot:\n\t%v", expectedT, gotT)
			}
		})
	}

}
