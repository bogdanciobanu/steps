package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestUnmarshalMap(t *testing.T) {
	os.Setenv("SP_DEBUG", "true")
	tests := []struct {
		name          string
		input         string
		ok            bool
		output        [][]string
		outputHeaders []string
		inputHeaders  []string
		showIndexes   bool
		byOrder       bool
	}{
		{
			name:          "String JSON",
			input:         `[{"key1":"val1", "key2":"val2"}, {"key1":"val3", "key2":"val4"}]`,
			ok:            true,
			output:        [][]string{{"val1", "val2"}, {"val3", "val4"}},
			outputHeaders: []string{"key1", "key2"},
			inputHeaders:  []string{"by_keys"},
		},
		{
			name:          "String JSON given headers",
			input:         `[{"key1":"val1", "key2":"val2"}, {"key1":"val3", "key2":"val4"}]`,
			ok:            true,
			output:        [][]string{{"val1", "val2"}, {"val3", "val4"}},
			outputHeaders: []string{"key1", "key2"},
			inputHeaders:  []string{"key1", "key2"},
			byOrder:       true,
		},
		{
			name:          "String JSON one header",
			input:         `[{"key1":"val1", "key2":"val2"}, {"key1":"val3", "key2":"val4"}]`,
			ok:            true,
			output:        [][]string{{"val1"}, {"val3"}},
			outputHeaders: []string{"key1"},
			inputHeaders:  []string{"key1"},
			byOrder:       true,
		},
		{
			name:          "String JSON show indexes",
			input:         `[{"key1":"val1", "key2":"val2"}, {"key1":"val3", "key2":"val4"}]`,
			ok:            true,
			output:        [][]string{{"0", "val1", "val2"}, {"1", "val3", "val4"}},
			outputHeaders: []string{"IDX", "key1", "key2"},
			inputHeaders:  []string{"by_keys"},
			showIndexes:   true,
		},
		{
			name:          "JSON show indexes",
			input:         `[{"integer":1, "string":"str", "bool":true, "float":3.1}, {"integer":2, "string":"str2", "bool":false, "float":4.1}]`,
			ok:            true,
			output:        [][]string{{"0", "1", "str", "true", "3.1"}, {"2", "str2", "false", "4.1"}},
			outputHeaders: []string{"IDX", "integer", "string", "bool", "float"},
			inputHeaders:  []string{"by_keys"},
			showIndexes:   true,
		},
		{
			name:  "invalid JSON",
			input: `[`,
			ok:    false,
		},
		{
			name:  "Not map",
			input: `[1,2,3]`,
			ok:    false,
		},
	}
	for _, tt := range tests {
		tabulate := &Tabulate{}
		t.Run(tt.name, func(t *testing.T) {
			unmarshaledMap, ok := tabulate.unmarshalMap(tt.input)
			require.Equal(t, tt.ok, ok)
			if !tt.ok {
				return
			}
			tabulate.handleMapArray(Args{
				Headers:     tt.inputHeaders,
				ShowIndexes: tt.showIndexes,
			}, unmarshaledMap)

			if tt.byOrder {
				assert.Equal(t, tt.output, tabulate.input)
			} else {
				require.Len(t, tabulate.input, len(tt.output))
				for i, expectedOutputArr := range tt.output {
					for _, val := range expectedOutputArr {
						found := false
						for _, actualVal := range tabulate.input[i] {
							if actualVal == val {
								found = true
								break
							}
						}
						require.True(t, found, "value %q not found in actual input at index %d: %#v", val, i, tabulate.input)
					}
				}
			}

			if tt.byOrder {
				assert.Equal(t, tt.outputHeaders, tabulate.headers)
			} else {
				require.Len(t, tabulate.headers, len(tt.outputHeaders))
				for _, header := range tt.outputHeaders {
					found := false
					for _, resHeader := range tabulate.headers {
						if resHeader == header {
							found = true
							break
						}
					}
					require.True(t, found, "header %q not found in output headers (%#v)", header, tt.outputHeaders)
				}
			}
		})
	}
}

func TestUnmarshalMapArrays(t *testing.T) {
	os.Setenv("SP_DEBUG", "true")
	tests := []struct {
		name          string
		input         string
		ok            bool
		output        [][]string
		outputHeaders []string
		inputHeaders  []string
		showIndexes   bool
		byOrder       bool
		groupKey      string
	}{
		{
			name:          "String JSON",
			input:         `{"one":[{"key1":"val1", "key2":"val2"}, {"key1":"val3", "key2":"val4"}], "two":[{"key1":"val5", "key2":"val6"}]}`,
			ok:            true,
			output:        [][]string{{"one", "val1", "val2"}, {"one", "val3", "val4"}, {"two", "val5", "val6"}},
			outputHeaders: []string{"group", "key1", "key2"},
			inputHeaders:  []string{"by_keys"},
			groupKey:      "group",
		},
		{
			name:          "String JSON given header",
			input:         `{"one":[{"key1":"val1", "key2":"val2"}, {"key1":"val3", "key2":"val4"}], "two":[{"key1":"val5", "key2":"val6"}]}`,
			ok:            true,
			output:        [][]string{{"one", "val1"}, {"one", "val3"}, {"two", "val5"}},
			outputHeaders: []string{"group", "key1"},
			inputHeaders:  []string{"key1"},
			byOrder:       true,
			groupKey:      "group",
		},
		{
			name:          "String JSON given header and index",
			input:         `{"one":[{"key1":"val1", "key2":"val2"}, {"key1":"val3", "key2":"val4"}], "two":[{"key1":"val5", "key2":"val6"}]}`,
			ok:            true,
			output:        [][]string{{"0", "one", "val1"}, {"1", "one", "val3"}, {"0", "two", "val5"}},
			outputHeaders: []string{"IDX", "group", "key1"},
			inputHeaders:  []string{"key1"},
			byOrder:       true,
			showIndexes:   true,
			groupKey:      "group",
		},
		{
			name:          "non string JSON",
			input:         `{"one":[{"integer":1, "string":"str", "bool":true, "float":3.1}, {"integer":2, "string":"str2", "bool":false, "float":4.1}], "two":[{"integer":3, "string":"str3", "bool":false, "float":5.1}]}`,
			ok:            true,
			output:        [][]string{{"one", "1", "str", "true", "3.1"}, {"one", "2", "str2", "false", "4.1"}, {"two", "3", "str3", "false", "5.1"}},
			outputHeaders: []string{"group", "integer", "string", "bool", "float"},
			inputHeaders:  []string{"by_keys"},
			groupKey:      "group",
		},
	}
	for _, tt := range tests {
		tabulate := &Tabulate{}
		t.Run(tt.name, func(t *testing.T) {
			unmarshaledMap, ok := tabulate.unmarshalMapArrays(tt.input)
			require.Equal(t, tt.ok, ok)
			if !tt.ok {
				return
			}
			tabulate.handleMapOfArrays(Args{
				Headers:         tt.inputHeaders,
				ShowIndexes:     tt.showIndexes,
				GroupHeaderName: tt.groupKey,
			}, unmarshaledMap)

			if tt.byOrder {
				assert.Equal(t, tt.output, tabulate.input)
			} else {
				require.Len(t, tabulate.input, len(tt.output))
				for i, expectedOutputArr := range tt.output {
					require.Len(t, tabulate.input[i], len(expectedOutputArr), "expected: %#v", expectedOutputArr)
					for _, val := range expectedOutputArr {
						found := false
						for _, actualVal := range tabulate.input[i] {
							if actualVal == val {
								found = true
								break
							}
						}
						require.True(t, found, "value %q not found in actual input at index %d: %#v", val, i, tabulate.input)
					}
				}
			}

			if tt.byOrder {
				assert.Equal(t, tt.outputHeaders, tabulate.headers)
			} else {
				require.Len(t, tabulate.headers, len(tt.outputHeaders))
				for _, header := range tt.outputHeaders {
					found := false
					for _, resHeader := range tabulate.headers {
						if resHeader == header {
							found = true
							break
						}
					}
					require.True(t, found, "header %q not found in output headers (%#v)", header, tabulate.headers)
				}
			}
		})
	}
}

func TestUnmarshalArray(t *testing.T) {
	os.Setenv("SP_DEBUG", "true")
	tests := []struct {
		name          string
		input         string
		ok            bool
		output        [][]string
		outputHeaders []string
		inputHeaders  []string
		showIndexes   bool
		err           error
	}{
		{
			name:          "String JSON",
			input:         `[["val1", "val2"],["val3", "val4"]]`,
			ok:            true,
			output:        [][]string{{"val1", "val2"}, {"val3", "val4"}},
			outputHeaders: []string{"key1", "key2"},
			inputHeaders:  []string{"key1", "key2"},
		},
		{
			name:          "String JSON not specifying headers",
			input:         `[["val1", "val2"],["val3", "val4"]]`,
			ok:            true,
			inputHeaders:  []string{"by_keys"},
			err: fmt.Errorf("can't build headers when arrays list provided, please specify headers explicitly"),
		},
		{
			name:          "String JSON show indexes",
			input:         `[["val1", "val2"],["val3", "val4"]]`,
			ok:            true,
			output:        [][]string{{"0", "val1", "val2"}, {"1", "val3", "val4"}},
			outputHeaders: []string{"IDX", "key1", "key2"},
			inputHeaders:  []string{"key1", "key2"},
			showIndexes:   true,
		},
		{
			name:          "JSON show indexes",
			input:         `[[1, "str", true, 3.1], [2, "str2", false, 4.1]]`,
			ok:            true,
			output:        [][]string{{"0", "1", "str", "true", "3.1"}, {"1", "2", "str2", "false", "4.1"}},
			outputHeaders: []string{"IDX", "integer", "string", "bool", "float"},
			inputHeaders:  []string{"integer", "string", "bool", "float"},
			showIndexes:   true,
		},
		{
			name:  "invalid JSON",
			input: `[`,
			ok:    false,
		},
		{
			name:  "Not array of arrays",
			input: `[1,2,3]`,
			ok:    false,
		},
	}
	for _, tt := range tests {
		tabulate := &Tabulate{}
		t.Run(tt.name, func(t *testing.T) {
			unmarshaledMap, ok := tabulate.unmarshalArray(tt.input)
			require.Equal(t, tt.ok, ok)
			if !tt.ok {
				return
			}
			err := tabulate.handleStringArrays(Args{
				Headers:     tt.inputHeaders,
				ShowIndexes: tt.showIndexes,
			}, unmarshaledMap)
			if tt.err != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.output, tabulate.input)
			assert.Equal(t, tt.outputHeaders, tabulate.headers)
		})
	}
}
