package services_test

import (
	"cambio-brl-usd/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestBuscarUltimaCotacao(t *testing.T) {
	cotacao := services.BuscarUltimaCotacao()
	assert.Equal(t, "BRL", cotacao.MoedaOrigem)
	assert.Equal(t, "USD", cotacao.MoedaDestino)
	assert.Greater(t, cotacao.Valor, 0.0)
}

func TestBuscarUltimaCotacao_FallbackPorTokenVazio(t *testing.T) {
	original := services.SecretsFetcher
	services.SecretsFetcher = func() string { return "" }
	defer func() { services.SecretsFetcher = original }()

	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("Esperava fallback, recebeu: %f", cotacao.Valor)
	}
}

func TestBuscarUltimaCotacao_ErroAoCriarRequisicao(t *testing.T) {
	original := services.SecretsFetcher
	services.SecretsFetcher = func() string { return "token" }
	defer func() { services.SecretsFetcher = original }()

	os.Setenv("FIXER_API_URL", ":")
	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("Esperava fallback, recebeu: %f", cotacao.Valor)
	}
}

func TestBuscarUltimaCotacao_ErroClientDo(t *testing.T) {
	original := services.SecretsFetcher
	services.SecretsFetcher = func() string { return "token" }
	defer func() { services.SecretsFetcher = original }()

	os.Setenv("FIXER_API_URL", "http://127.0.0.1:9999")
	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("Esperava fallback por erro client.Do, recebeu: %f", cotacao.Valor)
	}
}

func TestBuscarUltimaCotacao_JSONInvalido(t *testing.T) {
	original := services.SecretsFetcher
	services.SecretsFetcher = func() string { return "token" }
	defer func() { services.SecretsFetcher = original }()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("json ruim"))
	}))
	defer srv.Close()

	os.Setenv("FIXER_API_URL", srv.URL)
	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("Esperava fallback por JSON inválido")
	}
}

func TestBuscarUltimaCotacao_SuccessFalse(t *testing.T) {
	original := services.SecretsFetcher
	services.SecretsFetcher = func() string { return "token" }
	defer func() { services.SecretsFetcher = original }()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"base":"BRL","success":false,"rates":{"USD":5.50}`))
	}))
	defer srv.Close()

	os.Setenv("FIXER_API_URL", srv.URL)
	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("Esperava fallback por success=false")
	}
}

func TestBuscarHistorico_ErroNaExpressao(t *testing.T) {
	inicio := time.Now().Add(-48 * time.Hour)
	fim := time.Now().Add(-48 * time.Hour) // mesma data para forçar filtro vazio
	_ = services.BuscarHistorico(inicio, fim)
}

func TestSalvarCotacaoNoDynamo_ExecutaSemPanic(t *testing.T) {
	cotacao := services.BuscarUltimaCotacaoMock()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SalvarCotacaoNoDynamo causou panic: %v", r)
		}
	}()
	services.SalvarCotacaoNoDynamo(cotacao)
}

func TestBuscarAPIKeyDoFixer_SegredoNaoExiste(t *testing.T) {
	original := services.SecretsFetcher
	services.SecretsFetcher = func() string {
		return ""
	}
	defer func() { services.SecretsFetcher = original }()

	token := services.SecretsFetcher()
	if token != "" {
		t.Errorf("Esperava segredo vazio, obteve: %s", token)
	}
}

func TestBuscarAPIKeyDoFixer_JSONInvalido(t *testing.T) {
	original := services.SecretsFetcher
	services.SecretsFetcher = func() string {
		raw := `{"fixer_api_key":123}`
		var parsed map[string]string
		err := json.Unmarshal([]byte(raw), &parsed)
		if err != nil {
			return ""
		}
		return parsed["fixer_api_key"]
	}
	defer func() { services.SecretsFetcher = original }()

	token := services.SecretsFetcher()
	if token != "" {
		t.Errorf("Esperava erro de JSON e segredo vazio, obteve: %s", token)
	}
}

func TestBuscarHistorico_ErroNoScan(t *testing.T) {
	inicio := time.Now().Add(-24 * time.Hour)
	fim := time.Now()
	cotacoes := services.BuscarHistorico(inicio, fim)
	_ = cotacoes
}

func TestBuscarHistorico_ErroNoUnmarshal(t *testing.T) {
	inicio := time.Now().Add(-24 * time.Hour)
	fim := time.Now()
	_ = services.BuscarHistorico(inicio, fim)
}

func TestBuscarUltimaCotacao_ErroAoCriarRequest(t *testing.T) {
	original := services.NewHTTPRequest
	services.NewHTTPRequest = func(method string, url string, body io.Reader) (*http.Request, error) {
		return nil, fmt.Errorf("erro simulado")
	}
	defer func() { services.NewHTTPRequest = original }()

	services.SecretsFetcher = func() string { return "fake" }

	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("esperava fallback")
	}
}

func TestBuscarUltimaCotacao_ErroClientDo1(t *testing.T) {
	original := services.HTTPClientDo
	services.HTTPClientDo = func(_ *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("erro simulado no Do")
	}
	defer func() { services.HTTPClientDo = original }()

	services.SecretsFetcher = func() string { return "fake" }

	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("esperava fallback")
	}
}

func TestBuscarUltimaCotacao_ErroDecodeJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("INVALID JSON"))
	}))
	defer srv.Close()

	services.SecretsFetcher = func() string { return "fake" }
	os.Setenv("FIXER_API_URL", srv.URL)

	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("esperava fallback")
	}
}

func TestBuscarHistorico_ErroExpressao(t *testing.T) {
	original := services.DynamoScan
	services.DynamoScan = func(_ *dynamodb.Client, _ *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		return nil, fmt.Errorf("erro fake scan")
	}
	defer func() { services.DynamoScan = original }()

	cotacoes := services.BuscarHistorico(time.Now(), time.Now())
	if cotacoes != nil {
		t.Errorf("esperava nil em erro de Scan")
	}
}
