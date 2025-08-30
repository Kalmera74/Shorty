package validation

import (
	"testing"
)

func TestValidateUrl(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid http url",
			input:   "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid https url",
			input:   "https://example.com/path?query=1",
			wantErr: false,
		},
		{
			name:    "missing scheme",
			input:   "example.com",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid url with spaces",
			input:   "http://exa mple.com",
			wantErr: true,
		},
		{
			name:    "just garbage text",
			input:   "%%%^^^",
			wantErr: true,
		}, {
			name:    "invalid http url",
			input:   "htp:/invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUrl(tt.input)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got nil. Value: %v", tt.input)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got %v error for input: %v", err, tt.name)
			}

		})
	}
}
