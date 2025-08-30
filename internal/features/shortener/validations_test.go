package shortener

import (
	"testing"
)

func TestShortModel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		model   ShortModel
		wantErr bool
	}{
		{
			name:    "Valid model",
			model:   ShortModel{UserID: 1, OriginalUrl: "https://example.com", ShortUrl: "abc123"},
			wantErr: false,
		},
		{
			name:    "Missing OriginalUrl",
			model:   ShortModel{UserID: 1, OriginalUrl: "", ShortUrl: "abc123"},
			wantErr: true,
		},
		{
			name:    "Missing ShortUrl",
			model:   ShortModel{UserID: 1, OriginalUrl: "https://example.com", ShortUrl: ""},
			wantErr: true,
		},
		{
			name:    "Invalid UserID",
			model:   ShortModel{UserID: 0, OriginalUrl: "https://example.com", ShortUrl: "abc123"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got nil. Value: %v", tt.model)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got %v error for input: %v", err, tt.name)
			}

		})
	}
}

func TestShortenRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     ShortenRequest
		wantErr bool
	}{
		{
			name:    "Valid request",
			req:     ShortenRequest{UserID: 1, Url: "https://example.com"},
			wantErr: false,
		},
		{
			name:    "Empty URL",
			req:     ShortenRequest{UserID: 1, Url: ""},
			wantErr: true,
		},
		{
			name:    "Invalid URL",
			req:     ShortenRequest{UserID: 1, Url: "htp:/invalid"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got nil. Value: %v", tt.req)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error but got %v error for input: %v", err, tt.name)
			}

		})
	}
}
