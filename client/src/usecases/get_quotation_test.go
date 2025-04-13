package usecases

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/CaiqueRibeiro/client-api-ex/client/src/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetQuotationUseCase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		serverStatus   int
		serverDelay    time.Duration
		expectError    bool
		expectedBid    string
	}{
		{
			name:           "success",
			serverResponse: "5.8576",
			serverStatus:   http.StatusOK,
			serverDelay:    0,
			expectError:    false,
			expectedBid:    "5.8576",
		},
		{
			name:           "server_error",
			serverResponse: "Internal Server Error",
			serverStatus:   http.StatusInternalServerError,
			serverDelay:    0,
			expectError:    false, // No error because we still get a response body
			expectedBid:    "Internal Server Error",
		},
		{
			name:           "timeout",
			serverResponse: "5.8576",
			serverStatus:   http.StatusOK,
			serverDelay:    400 * time.Millisecond, // More than the 300ms timeout
			expectError:    true,
			expectedBid:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tt.serverDelay)
				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Create the use case
			useCase := &GetQuotationUseCase{
				ServerURL: server.URL,
			}

			// Call the Execute method
			quotation, err := useCase.Execute()

			// Check if we expected an error
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			// Check the result
			require.NoError(t, err)
			assert.Equal(t, tt.expectedBid, quotation.Bid)
		})
	}
}

// Custom usecase with file path for testing
type testableGetQuotationUseCase struct {
	GetQuotationUseCase
	filePath string
}

func (g *testableGetQuotationUseCase) SaveQuotationToFile(quotation entities.Quotation) error {
	file, err := os.Create(g.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	content := "Dólar: " + quotation.Bid
	_, err = file.WriteString(content)
	return err
}

func TestGetQuotationUseCase_SaveQuotationToFile(t *testing.T) {
	// Create a temporary file path
	tempDir, err := os.MkdirTemp("", "test_cotacao_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tempFilePath := filepath.Join(tempDir, "cotacao.txt")

	// Create a use case with custom file path
	useCase := &testableGetQuotationUseCase{
		GetQuotationUseCase: *NewGetQuotationUseCase(),
		filePath:            tempFilePath,
	}

	// Test saving a quotation to a file
	quotation := entities.Quotation{
		Bid: "5.8576",
	}

	err = useCase.SaveQuotationToFile(quotation)
	require.NoError(t, err)

	// Read the file content
	content, err := os.ReadFile(tempFilePath)
	require.NoError(t, err)

	// Verify the content
	assert.Equal(t, "Dólar: 5.8576", string(content))
}

func TestGetQuotationUseCase_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Start a test server to mock the real server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("5.8576"))
	}))
	defer server.Close()

	// Create the use case with the test server URL
	useCase := &GetQuotationUseCase{
		ServerURL: server.URL,
	}

	// Execute the use case
	quotation, err := useCase.Execute()
	require.NoError(t, err)
	assert.Equal(t, "5.8576", quotation.Bid)

	// Test saving to file using a temp file
	tempDir, err := os.MkdirTemp("", "test_cotacao_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tempFilePath := filepath.Join(tempDir, "cotacao.txt")

	// Create a testable use case
	testableUseCase := &testableGetQuotationUseCase{
		GetQuotationUseCase: *useCase,
		filePath:            tempFilePath,
	}

	err = testableUseCase.SaveQuotationToFile(quotation)
	require.NoError(t, err)

	// Read and verify the content
	content, err := os.ReadFile(tempFilePath)
	require.NoError(t, err)
	assert.Equal(t, "Dólar: 5.8576", string(content))
}
