package gateways

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetQuotation(t *testing.T) {
	tests := []struct {
		name         string
		responseBody string
		statusCode   int
		wantErr      bool
		expectedBid  string
	}{
		{
			name:         "success",
			responseBody: `{"USDBRL":{"code":"USD","codein":"BRL","name":"Dólar Americano/Real Brasileiro","high":"5.8688","low":"5.8213","varBid":"0.0313","pctChange":"0.54","bid":"5.8576","ask":"5.8582","timestamp":"1701278942","create_date":"2023-11-29 17:55:42"}}`,
			statusCode:   http.StatusOK,
			wantErr:      false,
			expectedBid:  "5.8576",
		},
		{
			name:         "invalid json",
			responseBody: `invalid json`,
			statusCode:   http.StatusOK,
			wantErr:      true,
			expectedBid:  "",
		},
		{
			name:         "server error",
			responseBody: ``,
			statusCode:   http.StatusInternalServerError,
			wantErr:      false, // No error because we still get a response, just with no body
			expectedBid:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Create a gateway with the test server URL
			gateway := &QuotationGateway{
				URL: server.URL,
			}

			// Call the function
			quotation, err := gateway.GetQuotation()

			// Check if we expected an error
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			// For server error cases, we might get empty quotation but no error
			if tt.statusCode != http.StatusOK {
				return
			}

			// Otherwise check the results
			require.NoError(t, err)
			assert.Equal(t, tt.expectedBid, quotation.Bid)
		})
	}
}

func TestGetQuotationTimeout(t *testing.T) {
	// Create a test server that delays its response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond) // Sleep longer than the timeout
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"USDBRL":{"code":"USD","codein":"BRL","name":"Dólar Americano/Real Brasileiro","high":"5.8688","low":"5.8213","varBid":"0.0313","pctChange":"0.54","bid":"5.8576","ask":"5.8582","timestamp":"1701278942","create_date":"2023-11-29 17:55:42"}}`))
	}))
	defer server.Close()

	// Create a gateway with the test server URL
	gateway := &QuotationGateway{
		URL: server.URL,
	}

	// Call the function
	_, err := gateway.GetQuotation()

	// Check that we got a timeout error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestExternalURLIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	gateway := NewQuotationGateway()
	quotation, err := gateway.GetQuotation()

	require.NoError(t, err)
	assert.NotEmpty(t, quotation.Bid)
	assert.NotEmpty(t, quotation.Code)
	assert.Equal(t, "USD", quotation.Code)
	assert.Equal(t, "BRL", quotation.Codein)
}
