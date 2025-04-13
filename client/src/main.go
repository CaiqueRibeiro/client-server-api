package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/CaiqueRibeiro/client-api-ex/client/src/usecases"
)

func main() {
	// Analisa os flags da linha de comando
	serverURL := flag.String("server", "http://localhost:8080/cotacao", "URL of the quotation server")
	outputPath := flag.String("output", "cotacao.txt", "Path to save the quotation")
	flag.Parse()

	// Cria um caso de uso personalizado com a URL do servidor e caminho de saída fornecidos
	getQuotationUseCase := &usecases.GetQuotationUseCase{
		ServerURL:  *serverURL,
		OutputPath: *outputPath,
	}

	quotation, err := getQuotationUseCase.Execute()
	if err != nil {
		log.Fatalf("Failed to get quotation: %v", err)
	}

	fmt.Printf("USD-BRL quotation: %s\n", quotation.Bid)

	// Salva a cotação no arquivo especificado
	err = getQuotationUseCase.SaveQuotationToFile(quotation)
	if err != nil {
		log.Fatalf("Failed to save quotation to file: %v", err)
	}

	fmt.Printf("Quotation saved to %s\n", *outputPath)
}
