package usecases

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CaiqueRibeiro/client-api-ex/client/src/entities"
)

type GetQuotationUseCase struct {
	ServerURL  string
	OutputPath string
}

func NewGetQuotationUseCase() *GetQuotationUseCase {
	return &GetQuotationUseCase{
		ServerURL:  "http://localhost:8080/cotacao",
		OutputPath: "cotacao.txt", // Caminho padrão
	}
}

func (g *GetQuotationUseCase) Execute() (entities.Quotation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.ServerURL, nil)
	if err != nil {
		log.Printf("Erro ao criar requisição: %v", err)
		return entities.Quotation{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao fazer requisição: %v", err)
		return entities.Quotation{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler corpo da resposta: %v", err)
		return entities.Quotation{}, err
	}

	quotation := entities.Quotation{
		Bid: string(body),
	}

	return quotation, nil
}

func (g *GetQuotationUseCase) SaveQuotationToFile(quotation entities.Quotation) error {
	outputPath := g.OutputPath
	if outputPath == "" {
		outputPath = "cotacao.txt" // Usa o padrão se não estiver definido
	}

	file, err := os.Create(outputPath)
	if err != nil {
		log.Printf("Erro ao criar arquivo: %v", err)
		return err
	}
	defer file.Close()

	content := fmt.Sprintf("Dólar: %s", quotation.Bid)
	_, err = file.WriteString(content)
	if err != nil {
		log.Printf("Erro ao escrever no arquivo: %v", err)
		return err
	}

	return nil
}
