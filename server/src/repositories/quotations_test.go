package repositories

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/CaiqueRibeiro/client-api-ex/server/src/gateways"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	suite.Suite
	db         *sql.DB
	repository *QuotationsRepository
}

func (suite *RepositoryTestSuite) SetupTest() {
	// Create an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(suite.T(), err)

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

	suite.db = db
	suite.repository = NewQuotationsRepository(db)
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *RepositoryTestSuite) TestCreate() {
	// Create a test quotation
	createDate, _ := time.Parse("2006-01-02 15:04:05", "2023-11-29 17:55:42")
	quotation := gateways.Quotation{
		USDBRL: gateways.USDBRL{
			Code:       "USD",
			Codein:     "BRL",
			Name:       "Dólar Americano/Real Brasileiro",
			High:       "5.8688",
			Low:        "5.8213",
			VarBid:     "0.0313",
			PctChange:  "0.54",
			Bid:        "5.8576",
			Ask:        "5.8582",
			Timestamp:  "1701278942",
			CreateDate: createDate,
		},
	}

	// Test the Create method
	err := suite.repository.Create(quotation)
	assert.NoError(suite.T(), err)

	// Verify the quotation was saved
	var count int
	err = suite.db.QueryRow("SELECT COUNT(*) FROM quotations").Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	// Verify the data
	var savedBid string
	err = suite.db.QueryRow("SELECT bid FROM quotations").Scan(&savedBid)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "5.8576", savedBid)
}

func (suite *RepositoryTestSuite) TestCreateWithContextTimeout() {
	// Create a test quotation
	createDate, _ := time.Parse("2006-01-02 15:04:05", "2023-11-29 17:55:42")
	quotation := gateways.Quotation{
		USDBRL: gateways.USDBRL{
			Code:       "USD",
			Codein:     "BRL",
			Name:       "Dólar Americano/Real Brasileiro",
			High:       "5.8688",
			Low:        "5.8213",
			VarBid:     "0.0313",
			PctChange:  "0.54",
			Bid:        "5.8576",
			Ask:        "5.8582",
			Timestamp:  "1701278942",
			CreateDate: createDate,
		},
	}

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Let the context expire
	time.Sleep(2 * time.Nanosecond)

	// Test the CreateWithContext method with an expired context
	err := suite.repository.CreateWithContext(ctx, quotation)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "context deadline exceeded")
}

// Run the test suite
func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
