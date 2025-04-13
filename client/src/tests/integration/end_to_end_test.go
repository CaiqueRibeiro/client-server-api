package integration

import (
	"database/sql"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EndToEndTestSuite struct {
	suite.Suite
	serverCmd   *exec.Cmd
	serverURL   string
	tempDir     string
	cotacaoPath string
	dbPath      string
}

func (suite *EndToEndTestSuite) SetupSuite() {
	// Skip if we're in short mode
	if testing.Short() {
		suite.T().Skip("Skipping end-to-end test in short mode")
	}

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "client_server_test_*")
	require.NoError(suite.T(), err)
	suite.tempDir = tempDir

	// Set paths
	suite.cotacaoPath = filepath.Join(suite.tempDir, "cotacao.txt")
	suite.dbPath = filepath.Join(suite.tempDir, "test_quotations.db")
	suite.serverURL = "http://localhost:8081"

	// Start the server with a different port and DB path
	suite.serverCmd = exec.Command("go", "run", "../../../server/src/main.go",
		"-port", "8081",
		"-db", suite.dbPath)

	// Set environment variables if needed
	suite.serverCmd.Env = append(os.Environ(), "TEST_MODE=true")

	// Start the server
	err = suite.serverCmd.Start()
	require.NoError(suite.T(), err)

	// Wait for the server to start
	time.Sleep(1 * time.Second)
}

func (suite *EndToEndTestSuite) TearDownSuite() {
	// Kill the server
	if suite.serverCmd != nil && suite.serverCmd.Process != nil {
		suite.serverCmd.Process.Kill()
	}

	// Clean up the temporary directory
	os.RemoveAll(suite.tempDir)
}

func (suite *EndToEndTestSuite) TestClientServerIntegration() {
	// Skip if needed parts are missing
	_, err := os.Stat("../../../server/src/main.go")
	if os.IsNotExist(err) {
		suite.T().Skip("Server main.go not found, skipping end-to-end test")
	}

	// Run the client with our test server URL
	clientCmd := exec.Command("go", "run", "../../main.go",
		"-server", suite.serverURL+"/cotacao",
		"-output", suite.cotacaoPath)

	// Run the client
	output, err := clientCmd.CombinedOutput()
	require.NoError(suite.T(), err, "Client failed to run: %s", output)

	_, err = os.Stat(suite.cotacaoPath)
	assert.NoError(suite.T(), err, "cotacao.txt was not created")

	content, err := os.ReadFile(suite.cotacaoPath)
	require.NoError(suite.T(), err)
	assert.Contains(suite.T(), string(content), "Dólar: ")

	db, err := sql.Open("sqlite3", suite.dbPath)
	if err == nil {
		defer db.Close()

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM quotations").Scan(&count)
		if err == nil {
			assert.GreaterOrEqual(suite.T(), count, 1, "No records were inserted in the database")
		}
	}
}

func TestEndToEndSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping end-to-end test in short mode")
	}
	suite.Run(t, new(EndToEndTestSuite))
}
