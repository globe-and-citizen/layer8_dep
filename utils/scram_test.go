package utils

import "testing"

func TestXOR(t *testing.T) {
	tc := []struct {
		name string
		a    []byte
		b    []byte
		want []byte
	}{
		{
			name: "equal length",
			a:    []byte{0x01, 0x02, 0x03},
			b:    []byte{0x04, 0x05, 0x06},
			want: []byte{0x05, 0x07, 0x05},
		},
		{
			name: "unequal length",
			a:    []byte{0x01, 0x02, 0x03},
			b:    []byte{0x04, 0x05},
			want: nil,
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			got, err := XOR(c.a, c.b)
			if err != nil {
				if c.want != nil {
					t.Errorf("got %v, want %v", err, c.want)
				}
			} else if string(got) != string(c.want) {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestConcatenateScramAttributes(t *testing.T) {
	tc := []struct {
		name       string
		attributes map[string]string
		want       string
	}{
		{
			name:       "empty attributes",
			attributes: map[string]string{},
			want:       "",
		},
		{
			name:       "one attribute",
			attributes: map[string]string{"a": "b"},
			want:       "a=b",
		},
		{
			name:       "two attributes",
			attributes: map[string]string{"a": "b", "c": "d"},
			want:       "a=b,c=d",
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			got := ConcatenateScramAttributes(c.attributes)
			if got != c.want {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}

func TestParseScramAttributes(t *testing.T) {
	tc := []struct {
		name             string
		attributeMessage string
		want             map[string]string
	}{
		{
			name:             "empty attributes",
			attributeMessage: "",
			want:             map[string]string{},
		},
		{
			name:             "one attribute",
			attributeMessage: "a=b",
			want:             map[string]string{"a": "b"},
		},
		{
			name:             "two attributes",
			attributeMessage: "a=b,c=d",
			want:             map[string]string{"a": "b", "c": "d"},
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			got := ParseScramAttributes(c.attributeMessage)
			if len(got) != len(c.want) {
				t.Errorf("got %v, want %v", got, c.want)
			}
			for k, v := range got {
				if c.want[k] != v {
					t.Errorf("got %v, want %v", got, c.want)
				}
			}
		})
	}
}
