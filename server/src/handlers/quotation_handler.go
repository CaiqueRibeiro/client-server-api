package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/CaiqueRibeiro/client-api-ex/server/src/gateways"
)

// Interfaces para dependências
type QuotationGateway interface {
	GetQuotation() (gateways.Quotation, error)
}

type QuotationRepository interface {
	Create(quotation gateways.Quotation) error
	CreateWithContext(ctx context.Context, quotation gateways.Quotation) error
}

type QuotationHandler struct {
	gateway    QuotationGateway
	repository QuotationRepository
}

func NewQuotationHandler(gateway QuotationGateway, repository QuotationRepository) *QuotationHandler {
	return &QuotationHandler{
		gateway:    gateway,
		repository: repository,
	}
}

func (h *QuotationHandler) HandleGetQuotation(w http.ResponseWriter, r *http.Request) {
	quotation, err := h.gateway.GetQuotation()
	if err != nil {
		log.Printf("Erro ao obter cotação da API: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*10))
	defer cancel()

	err = h.repository.CreateWithContext(ctx, quotation)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Tempo excedido ao persistir cotação no banco de dados: %v", err)
		} else {
			log.Printf("Erro ao persistir cotação no banco de dados: %v", err)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bid := quotation.Bid

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(bid))
}
