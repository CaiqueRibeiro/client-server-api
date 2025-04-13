package integration

import (
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/CaiqueRibeiro/client-api-ex/server/src/gateways"
	"github.com/CaiqueRibeiro/client-api-ex/server/src/handlers"
	"github.com/CaiqueRibeiro/client-api-ex/server/src/repositories"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ServerIntegrationTestSuite struct {
	suite.Suite
	server   *httptest.Server
	db       *sql.DB
	tempFile string
}

func (suite *ServerIntegrationTestSuite) SetupSuite() {
	// Create a temporary SQLite database file
	tempFile, err := os.CreateTemp("", "test_quotations_*.db")
	require.NoError(suite.T(), err)
	suite.tempFile = tempFile.Name()
	tempFile.Close()

	// Open the database
	db, err := sql.Open("sqlite3", suite.tempFile)
	require.NoError(suite.T(), err)
	suite.db = db

	// Create the table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS quotations (
		id TEXT PRIMARY KEY,
		code TEXT,
		codein TEXT,
		name TEXT,
		high TEXT,
		low TEXT,
		varBid TEXT,
		pctChange TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		create_date TEXT
	)`)
	require.NoError(suite.T(), err)

	// Create dependencies
	repository := repositories.NewQuotationsRepository(db)
	gateway := gateways.NewQuotationGateway()
	handler := handlers.NewQuotationHandler(gateway, repository)

	// Create a test server
	mux := http.NewServeMux()
	mux.HandleFunc("GET /cotacao", handler.HandleGetQuotation)
	server := httptest.NewServer(mux)
	suite.server = server
}

func (suite *ServerIntegrationTestSuite) TearDownSuite() {
	suite.server.Close()
	suite.db.Close()
	os.Remove(suite.tempFile)
}

func (suite *ServerIntegrationTestSuite) TestQuotationEndpoint() {
	// Make a request to the server
	resp, err := http.Get(suite.server.URL + "/cotacao")
	require.NoError(suite.T(), err)
	defer resp.Body.Close()

	// Check the status code
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(suite.T(), err)

	// The response should not be empty (exact bid value will vary)
	assert.NotEmpty(suite.T(), string(body))

	// Verify that the quotation was saved in the database
	var count int
	err = suite.db.QueryRow("SELECT COUNT(*) FROM quotations").Scan(&count)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	// Verify that the bid was saved
	var bid string
	err = suite.db.QueryRow("SELECT bid FROM quotations").Scan(&bid)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), string(body), bid)
}

// Run the test suite
func TestServerIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	suite.Run(t, new(ServerIntegrationTestSuite))
}
