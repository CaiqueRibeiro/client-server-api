package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/CaiqueRibeiro/client-api-ex/server/src/gateways"
	"github.com/google/uuid"
)

type QuotationsRepository struct {
	Db *sql.DB
}

func NewQuotationsRepository(db *sql.DB) *QuotationsRepository {
	return &QuotationsRepository{Db: db}
}

func (r *QuotationsRepository) Create(quotation gateways.Quotation) error {
	// Usa o método com contexto de background para compatibilidade retroativa
	return r.CreateWithContext(context.Background(), quotation)
}

func (r *QuotationsRepository) CreateWithContext(ctx context.Context, quotation gateways.Quotation) error {
	query := `
		INSERT INTO quotations (
		id,
		code,
		codein,
		name,
		high,
		low,
		varBid,
		pctChange,
		bid,
		ask,
		timestamp,
		create_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	id := uuid.New().String()

	_, err := r.Db.ExecContext(
		ctx,
		query,
		id,
		quotation.Code,
		quotation.Codein,
		quotation.Name,
		quotation.High,
		quotation.Low,
		quotation.VarBid,
		quotation.PctChange,
		quotation.Bid,
		quotation.Ask,
		quotation.Timestamp,
		quotation.CreateDate,
	)
	if err != nil {
		return fmt.Errorf("falha ao inserir cotação: %w", err)
	}

	return nil
}
