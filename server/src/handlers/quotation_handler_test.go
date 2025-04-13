package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CaiqueRibeiro/client-api-ex/server/src/gateways"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock gateway
type MockQuotationGateway struct {
	mock.Mock
}

func (m *MockQuotationGateway) GetQuotation() (gateways.Quotation, error) {
	args := m.Called()
	return args.Get(0).(gateways.Quotation), args.Error(1)
}

// Mock repository
type MockQuotationsRepository struct {
	mock.Mock
}

func (m *MockQuotationsRepository) Create(quotation gateways.Quotation) error {
	args := m.Called(quotation)
	return args.Error(0)
}

func (m *MockQuotationsRepository) CreateWithContext(ctx context.Context, quotation gateways.Quotation) error {
	args := m.Called(ctx, quotation)
	return args.Error(0)
}

func TestHandleGetQuotation(t *testing.T) {
	// Test cases
	tests := []struct {
		name                string
		gatewayError        error
		repositoryError     error
		expectedStatus      int
		expectedResponse    string
		quotationToReturn   gateways.Quotation
		isContextDeadlineEx bool
	}{
		{
			name:             "success",
			gatewayError:     nil,
			repositoryError:  nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: "5.8576",
			quotationToReturn: gateways.Quotation{
				USDBRL: gateways.USDBRL{
					Bid: "5.8576",
				},
			},
			isContextDeadlineEx: false,
		},
		{
			name:             "gateway error",
			gatewayError:     errors.New("gateway error"),
			repositoryError:  nil,
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "gateway error\n",
			quotationToReturn: gateways.Quotation{
				USDBRL: gateways.USDBRL{},
			},
			isContextDeadlineEx: false,
		},
		{
			name:             "repository error",
			gatewayError:     nil,
			repositoryError:  errors.New("repository error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "repository error\n",
			quotationToReturn: gateways.Quotation{
				USDBRL: gateways.USDBRL{
					Bid: "5.8576",
				},
			},
			isContextDeadlineEx: false,
		},
		{
			name:             "context deadline exceeded",
			gatewayError:     nil,
			repositoryError:  context.DeadlineExceeded,
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "context deadline exceeded\n",
			quotationToReturn: gateways.Quotation{
				USDBRL: gateways.USDBRL{
					Bid: "5.8576",
				},
			},
			isContextDeadlineEx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockGateway := new(MockQuotationGateway)
			mockRepository := new(MockQuotationsRepository)

			// Setup expectations
			mockGateway.On("GetQuotation").Return(tt.quotationToReturn, tt.gatewayError)

			// We only mock the repository call if the gateway call succeeds
			if tt.gatewayError == nil {
				mockRepository.On("CreateWithContext", mock.Anything, tt.quotationToReturn).Return(tt.repositoryError)
			}

			// Create the handler with the mocks
			handler := NewQuotationHandler(mockGateway, mockRepository)

			// Create a request and recorder for the test
			req := httptest.NewRequest(http.MethodGet, "/cotacao", nil)
			recorder := httptest.NewRecorder()

			// Call the handler
			handler.HandleGetQuotation(recorder, req)

			// Assert the results
			assert.Equal(t, tt.expectedStatus, recorder.Code)
			assert.Equal(t, tt.expectedResponse, recorder.Body.String())

			// Verify the expectations
			mockGateway.AssertExpectations(t)
			if tt.gatewayError == nil {
				mockRepository.AssertExpectations(t)
			}
		})
	}
}

func TestHandleGetQuotationIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create real instances (this is an integration test)
	gateway := gateways.NewQuotationGateway()

	// For the repository, we'll still use a mock to avoid DB dependencies
	mockRepository := new(MockQuotationsRepository)
	mockRepository.On("CreateWithContext", mock.Anything, mock.Anything).Return(nil)

	// Create the handler
	handler := NewQuotationHandler(gateway, mockRepository)

	// Create a request and recorder for the test
	req := httptest.NewRequest(http.MethodGet, "/cotacao", nil)
	recorder := httptest.NewRecorder()

	// Call the handler
	handler.HandleGetQuotation(recorder, req)

	// Assert the results
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.NotEmpty(t, recorder.Body.String())

	// Verify the expectations
	mockRepository.AssertExpectations(t)
}
