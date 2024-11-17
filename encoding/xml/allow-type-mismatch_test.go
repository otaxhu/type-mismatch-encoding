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

package xml

import (
	"slices"
	"strings"
	"testing"
)

func TestAllowTypeMismatchDecode(t *testing.T) {
	type T struct {
		XMLName      struct{}  `xml:"t"`
		String       string    `xml:"string"`
		Int          int       `xml:"int"`
		Float64      float64   `xml:"float64"`
		SliceString  []string  `xml:"sliceString>i"`
		SliceInt     []int     `xml:"sliceInt>i"`
		SliceFloat64 []float64 `xml:"sliceFloat64>i"`

		AttrString  string  `xml:"attrString,attr"`
		AttrInt     int     `xml:"attrInt,attr"`
		AttrFloat64 float64 `xml:"attrFloat64,attr"`
	}
	equalsT := func(a, b T) bool {
		return a.String == b.String &&
			a.Int == b.Int &&
			a.Float64 == b.Float64 &&
			slices.Equal(a.SliceString, b.SliceString) &&
			slices.Equal(a.SliceInt, b.SliceInt) &&
			slices.Equal(a.SliceFloat64, b.SliceFloat64)
	}
	baseT := T{
		String:  "test",
		Int:     123,
		Float64: 123.123,
		SliceString: []string{
			"test1",
			"test2",
			"test3",
		},
		SliceInt: []int{
			123,
			456,
			789,
		},
		SliceFloat64: []float64{
			123.123,
			456.456,
			789.789,
		},
		AttrString:  "testAttr",
		AttrInt:     123,
		AttrFloat64: 123.123,
	}
	testCases := []struct {
		name      string
		input     string
		expectedT func() T
	}{
		{
			name: "MismatchedType_String_GotDocument",
			input: Header + `
				<t attrString="testAttr" attrInt="123" attrFloat64="123.123">
					<string><MISMATCHED_TYPE></MISMATCHED_TYPE></string>
					<int>123</int>
					<float64>123.123</float64>
					<sliceString>
						<i>test1</i>
						<i>test2</i>
						<i>test3</i>
					</sliceString>
					<sliceInt>
						<i>123</i>
						<i>456</i>
						<i>789</i>
					</sliceInt>
					<sliceFloat64>
						<i>123.123</i>
						<i>456.456</i>
						<i>789.789</i>
					</sliceFloat64>
				</t>
			`,
			expectedT: func() T {
				ret := baseT
				ret.String = ""
				return ret
			},
		},
		{
			name: "MismatchedType_Int_GotString",
			input: Header + `
				<t attrString="testAttr" attrInt="123" attrFloat64="123.123">
					<string>test</string>
					<int>MISMATCHED_TYPE</int>
					<float64>123.123</float64>
					<sliceString>
						<i>test1</i>
						<i>test2</i>
						<i>test3</i>
					</sliceString>
					<sliceInt>
						<i>123</i>
						<i>456</i>
						<i>789</i>
					</sliceInt>
					<sliceFloat64>
						<i>123.123</i>
						<i>456.456</i>
						<i>789.789</i>
					</sliceFloat64>
				</t>
			`,
			expectedT: func() T {
				ret := baseT
				ret.Int = 0
				return ret
			},
		},
		{
			name: "MismatchedType_Int_GotFloat64",
			input: Header + `
				<t attrString="testAttr" attrInt="123" attrFloat64="123.123">
					<string>test</string>
					<int>123.123</int>
					<float64>123.123</float64>
					<sliceString>
						<i>test1</i>
						<i>test2</i>
						<i>test3</i>
					</sliceString>
					<sliceInt>
						<i>123</i>
						<i>456</i>
						<i>789</i>
					</sliceInt>
					<sliceFloat64>
						<i>123.123</i>
						<i>456.456</i>
						<i>789.789</i>
					</sliceFloat64>
				</t>
			`,
			expectedT: func() T {
				ret := baseT
				ret.Int = 0
				return ret
			},
		},
		{
			name: "MismatchedType_Float64_GotString",
			input: Header + `
				<t attrString="testAttr" attrInt="123" attrFloat64="123.123">
					<string>test</string>
					<int>123</int>
					<float64>MISMATCHED_TYPE</float64>
					<sliceString>
						<i>test1</i>
						<i>test2</i>
						<i>test3</i>
					</sliceString>
					<sliceInt>
						<i>123</i>
						<i>456</i>
						<i>789</i>
					</sliceInt>
					<sliceFloat64>
						<i>123.123</i>
						<i>456.456</i>
						<i>789.789</i>
					</sliceFloat64>
				</t>
			`,
			expectedT: func() T {
				ret := baseT
				ret.Float64 = 0
				return ret
			},
		},
		{
			name: "MismatchedType_SliceString_GotString",
			input: Header + `
				<t attrString="testAttr" attrInt="123" attrFloat64="123.123">
					<string>test</string>
					<int>123</int>
					<float64>123.123</float64>
					<sliceString>MISMATCHED_TYPE</sliceString>
					<sliceInt>
						<i>123</i>
						<i>456</i>
						<i>789</i>
					</sliceInt>
					<sliceFloat64>
						<i>123.123</i>
						<i>456.456</i>
						<i>789.789</i>
					</sliceFloat64>
				</t>
			`,
			expectedT: func() T {
				ret := baseT
				ret.SliceString = []string{}
				return ret
			},
		},
		{
			name: "MismatchedType_SliceInt_InsideItems_GotString",
			input: Header + `
				<t attrString="testAttr" attrInt="123" attrFloat64="123.123">
					<string>test</string>
					<int>123</int>
					<float64>123.123</float64>
					<sliceString>
						<i>test1</i>
						<i>test2</i>
						<i>test3</i>
					</sliceString>
					<sliceInt>
						<i>MISMATCHED_TYPE</i>
						<i>MISMATCHED_TYPE</i>
						<i>MISMATCHED_TYPE</i>
					</sliceInt>
					<sliceFloat64>
						<i>123.123</i>
						<i>456.456</i>
						<i>789.789</i>
					</sliceFloat64>
				</t>
			`,
			expectedT: func() T {
				ret := baseT
				ret.SliceInt = make([]int, 3)
				return ret
			},
		},
		{
			name: "MismatchedType_SliceInt_InsideItems_GotFloat64",
			input: Header + `
				<t attrString="testAttr" attrInt="123" attrFloat64="123.123">
					<string>test</string>
					<int>123</int>
					<float64>123.123</float64>
					<sliceString>
						<i>test1</i>
						<i>test2</i>
						<i>test3</i>
					</sliceString>
					<sliceInt>
						<i>123.123</i>
						<i>456.456</i>
						<i>789.789</i>
					</sliceInt>
					<sliceFloat64>
						<i>123.123</i>
						<i>456.456</i>
						<i>789.789</i>
					</sliceFloat64>
				</t>
			`,
			expectedT: func() T {
				ret := baseT
				ret.SliceInt = make([]int, 3)
				return ret
			},
		},
		{
			name: "MismatchedType_SliceFloat64_InsideItems_GotString",
			input: Header + `
				<t attrString="testAttr" attrInt="123" attrFloat64="123.123">
					<string>test</string>
					<int>123</int>
					<float64>123.123</float64>
					<sliceString>
						<i>test1</i>
						<i>test2</i>
						<i>test3</i>
					</sliceString>
					<sliceInt>
						<i>123</i>
						<i>456</i>
						<i>789</i>
					</sliceInt>
					<sliceFloat64>
						<i>MISMATCHED_TYPE</i>
						<i>MISMATCHED_TYPE</i>
						<i>MISMATCHED_TYPE</i>
					</sliceFloat64>
				</t>
			`,
			expectedT: func() T {
				ret := baseT
				ret.SliceFloat64 = make([]float64, 3)
				return ret
			},
		},
		{
			name: "MismatchedType_AttrInt_GotString",
			input: Header + `
				<t attrString="testAttr" attrInt="MISMATCHED_TYPE" attrFloat64="123.123">
					<string>test</string>
					<int>123</int>
					<float64>123.123</float64>
					<sliceString>
						<i>test1</i>
						<i>test2</i>
						<i>test3</i>
					</sliceString>
					<sliceInt>
						<i>123</i>
						<i>456</i>
						<i>789</i>
					</sliceInt>
					<sliceFloat64>
						<i>123.123</i>
						<i>456.456</i>
						<i>789.789</i>
					</sliceFloat64>
				</t>
			`,
			expectedT: func() T {
				ret := baseT
				ret.AttrInt = 0
				return ret
			},
		},
		{
			name: "MismatchedType_AttrInt_GotFloat64",
			input: Header + `
				<t attrString="testAttr" attrInt="123.123" attrFloat64="123.123">
					<string>test</string>
					<int>123</int>
					<float64>123.123</float64>
					<sliceString>
						<i>test1</i>
						<i>test2</i>
						<i>test3</i>
					</sliceString>
					<sliceInt>
						<i>123</i>
						<i>456</i>
						<i>789</i>
					</sliceInt>
					<sliceFloat64>
						<i>123.123</i>
						<i>456.456</i>
						<i>789.789</i>
					</sliceFloat64>
				</t>
			`,
			expectedT: func() T {
				ret := baseT
				ret.AttrInt = 0
				return ret
			},
		},
		{
			name: "MismatchedType_AttrFloat64_GotString",
			input: Header + `
				<t attrString="testAttr" attrInt="123" attrFloat64="MISMATCHED_TYPE">
					<string>test</string>
					<int>123</int>
					<float64>123.123</float64>
					<sliceString>
						<i>test1</i>
						<i>test2</i>
						<i>test3</i>
					</sliceString>
					<sliceInt>
						<i>123</i>
						<i>456</i>
						<i>789</i>
					</sliceInt>
					<sliceFloat64>
						<i>123.123</i>
						<i>456.456</i>
						<i>789.789</i>
					</sliceFloat64>
				</t>
			`,
			expectedT: func() T {
				ret := baseT
				ret.AttrFloat64 = 0
				return ret
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dec := NewDecoder(strings.NewReader(tc.input))
			dec.AllowTypeMismatch = true
			var gotT T
			err := dec.Decode(&gotT)
			if err != nil {
				t.Fatal(err)
			}
			expectedT := tc.expectedT()
			if !equalsT(gotT, expectedT) {
				t.Fatalf("expected:\n\t%v\ngot:\n\t%v", expectedT, gotT)
			}
		})
	}
}
