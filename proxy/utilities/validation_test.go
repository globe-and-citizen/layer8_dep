package utilities

import "testing"

func TestValidateRequiredFields(t *testing.T) {
	type tp struct {
		Name string `validate:"nonzero"`
		Age  int    `validate:"nonzero"`
	}

	tc := []struct {
		name      string
		input     tp
		expectErr bool
	}{
		{
			name: "valid",
			input: tp{
				Name: "test",
				Age:  10,
			},
			expectErr: false,
		},
		{
			name: "invalid",
			input: tp{
				Age: 10,
			},
			expectErr: true,
		},
		{
			name: "invalid",
			input: tp{
				Name: "test",
			},
			expectErr: true,
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateRequiredFields(c.input)
			if err != nil && !c.expectErr {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && c.expectErr {
				t.Errorf("expected error but got nil")
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tc := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{
			name:      "valid",
			input:     "ab@mail.com",
			expectErr: false,
		},
		{
			name:      "invalid",
			input:     "abmail.com",
			expectErr: true,
		},
		{
			name:      "invalid",
			input:     "ab@mail",
			expectErr: true,
		},
		{
			name:      "invalid",
			input:     "abmail",
			expectErr: true,
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateEmail(c.input)
			if err != nil && !c.expectErr {
				t.Errorf("unexpected error: %v", err)
			}
			if err == nil && c.expectErr {
				t.Errorf("expected error but got nil")
			}
		})
	}
}
