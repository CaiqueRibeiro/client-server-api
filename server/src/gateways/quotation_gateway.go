package gateways

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

type USDBRL struct {
	Code       string    `json:"code"`
	Codein     string    `json:"codein"`
	Name       string    `json:"name"`
	High       string    `json:"high"`
	Low        string    `json:"low"`
	VarBid     string    `json:"varBid"`
	PctChange  string    `json:"pctChange"`
	Bid        string    `json:"bid"`
	Ask        string    `json:"ask"`
	Timestamp  string    `json:"timestamp"`
	CreateDate time.Time `json:"create_date"`
}

type Quotation struct {
	USDBRL `json:"USDBRL"`
}

type QuotationGateway struct {
	URL string
}

func NewQuotationGateway() *QuotationGateway {
	return &QuotationGateway{
		URL: "https://economia.awesomeapi.com.br/json/last/USD-BRL",
	}
}

func (g *QuotationGateway) GetQuotation() (Quotation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*200))
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.URL, nil)
	if err != nil {
		log.Printf("Erro ao criar requisição: %v", err)
		return Quotation{}, err
	}

	c := &http.Client{}

	resp, err := c.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Tempo excedido ao chamar API externa: %v", err)
		} else {
			log.Printf("Erro ao chamar API externa: %v", err)
		}
		return Quotation{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler corpo da resposta: %v", err)
		return Quotation{}, err
	}

	var rawQuotation map[string]map[string]interface{}
	err = json.Unmarshal(body, &rawQuotation)
	if err != nil {
		log.Printf("Erro ao desserializar JSON: %v", err)
		return Quotation{}, err
	}

	usdbrl := rawQuotation["USDBRL"]
	createDate, err := time.Parse("2006-01-02 15:04:05", usdbrl["create_date"].(string))
	if err != nil {
		log.Printf("Erro ao analisar create_date: %v", err)
		return Quotation{}, err
	}

	quotation := Quotation{
		USDBRL: USDBRL{
			Code:       usdbrl["code"].(string),
			Codein:     usdbrl["codein"].(string),
			Name:       usdbrl["name"].(string),
			High:       usdbrl["high"].(string),
			Low:        usdbrl["low"].(string),
			VarBid:     usdbrl["varBid"].(string),
			PctChange:  usdbrl["pctChange"].(string),
			Bid:        usdbrl["bid"].(string),
			Ask:        usdbrl["ask"].(string),
			Timestamp:  usdbrl["timestamp"].(string),
			CreateDate: createDate,
		},
	}

	return quotation, nil
}
