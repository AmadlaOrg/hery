package object

import (
	"reflect"
	"testing"
)

func TestMergeMultilevel(t *testing.T) {
	tests := []struct {
		name     string
		dTo      map[string]interface{}
		dFrom    map[string]interface{}
		inplace  bool
		expected map[string]interface{}
	}{
		{
			name: "simple merge",
			dTo: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{
					"subkey1": "subvalue1",
				},
			},
			dFrom: map[string]interface{}{
				"key2": map[string]interface{}{
					"subkey2": "subvalue2",
				},
				"key3": "value3",
			},
			inplace: false,
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{
					"subkey1": "subvalue1",
					"subkey2": "subvalue2",
				},
				"key3": "value3",
			},
		},
		{
			name: "inplace merge",
			dTo: map[string]interface{}{
				"key1": "value1",
			},
			dFrom: map[string]interface{}{
				"key2": "value2",
			},
			inplace: true,
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "overwrite value",
			dTo: map[string]interface{}{
				"key1": "value1",
			},
			dFrom: map[string]interface{}{
				"key1": "newvalue1",
			},
			inplace: false,
			expected: map[string]interface{}{
				"key1": "newvalue1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeMultilevel(tt.dTo, tt.dFrom, tt.inplace)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDeepCopy(t *testing.T) {
	original := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"subkey1": "subvalue1",
		},
	}

	copied := DeepCopy(original)

	if !reflect.DeepEqual(original, copied) {
		t.Errorf("expected %v, got %v", original, copied)
	}

	// Modify the copy and ensure the original is not affected
	copied["key1"] = "newvalue1"
	copied["key2"].(map[string]interface{})["subkey1"] = "newsubvalue1"

	if reflect.DeepEqual(original, copied) {
		t.Errorf("original map was modified when the copy was changed")
	}
}
