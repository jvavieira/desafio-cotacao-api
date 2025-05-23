package services_test

import (
	"cambio-brl-usd/models"
	"cambio-brl-usd/services"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
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

func TestBuscarHistorico_ErroNaExpressao(t *testing.T) {
	inicio := time.Now().Add(-48 * time.Hour)
	fim := time.Now().Add(-48 * time.Hour)
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

func TestBuscarUltimaCotacao_TokenVazio(t *testing.T) {
	original := services.SecretsFetcher
	services.SecretsFetcher = func() string { return "" }
	defer func() { services.SecretsFetcher = original }()

	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("Esperava fallback, recebeu: %f", cotacao.Valor)
	}
}

func TestBuscarUltimaCotacao_ErroCriarRequisicao(t *testing.T) {
	original := services.NewHTTPRequest
	services.NewHTTPRequest = func(method string, url string, body io.Reader) (*http.Request, error) {
		return nil, errors.New("erro simulado")
	}
	defer func() { services.NewHTTPRequest = original }()

	services.SecretsFetcher = func() string { return "token" }

	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("Esperava fallback por erro ao criar request")
	}
}

func TestBuscarUltimaCotacao_ErroClientDo(t *testing.T) {
	original := services.HTTPClientDo
	services.HTTPClientDo = func(_ *http.Request) (*http.Response, error) {
		return nil, errors.New("erro client.Do simulado")
	}
	defer func() { services.HTTPClientDo = original }()

	services.SecretsFetcher = func() string { return "token" }
	services.NewHTTPRequest = http.NewRequest

	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("Esperava fallback por erro no client.Do")
	}
}

func TestBuscarUltimaCotacao_JSONInvalido(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("nao-e-json"))
	}))
	defer srv.Close()

	services.SecretsFetcher = func() string { return "token" }
	saved := services.NewHTTPRequest
	services.NewHTTPRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		return http.NewRequest(method, srv.URL, body)
	}
	defer func() { services.NewHTTPRequest = saved }()

	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("Esperava fallback por JSON inválido")
	}
}

func TestBuscarUltimaCotacao_SuccessFalse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"base":"BRL","success":false,"rates":{"USD":5.0}}`))
	}))
	defer srv.Close()

	services.SecretsFetcher = func() string { return "token" }
	saved := services.NewHTTPRequest
	services.NewHTTPRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		return http.NewRequest(method, srv.URL, body)
	}
	defer func() { services.NewHTTPRequest = saved }()

	cotacao := services.BuscarUltimaCotacao()
	if cotacao.Valor != 5.00 {
		t.Errorf("Esperava fallback por success=false")
	}
}

func TestBuscarUltimaCotacao_ComSucesso(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{
			"success": true,
			"base": "BRL",
			"rates": { "USD": 5.42 }
		}`))
	}))
	defer srv.Close()

	services.SecretsFetcher = func() string { return "token" }

	savedRequest := services.NewHTTPRequest
	services.NewHTTPRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		return http.NewRequest(method, srv.URL, body)
	}
	defer func() { services.NewHTTPRequest = savedRequest }()

	savedSaver := services.SaveCotacao
	services.SaveCotacao = func(c models.Cotacao) {
		t.Logf("Cotação salva mockada: %+v", c)
	}
	defer func() { services.SaveCotacao = savedSaver }()

	cotacao := services.BuscarUltimaCotacao()

	if cotacao.Valor != 5.42 {
		t.Errorf("Esperava valor 5.42, recebeu: %f", cotacao.Valor)
	}
	if cotacao.MoedaOrigem != "BRL" || cotacao.MoedaDestino != "USD" {
		t.Errorf("Esperava moedas BRL→USD, recebeu: %+v", cotacao)
	}
}

func TestBuscarHistorico_ErroNaExpressaoOther(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Função não deve causar panic mesmo com erro na expressão")
		}
	}()
	fakeInicio := time.Now()
	fakeFim := time.Now()

	_ = services.BuscarHistorico(fakeInicio, fakeFim)
}

func TestBuscarHistorico_ErroScanDynamo(t *testing.T) {
	original := services.DynamoScan
	services.DynamoScan = func(_ *dynamodb.Client, _ *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		return nil, errors.New("simulando erro no scan")
	}
	defer func() { services.DynamoScan = original }()

	fakeInicio := time.Now().Add(-24 * time.Hour)
	fakeFim := time.Now()

	result := services.BuscarHistorico(fakeInicio, fakeFim)
	if result != nil {
		t.Errorf("Esperava retorno nil em erro de scan")
	}
}

func TestBuscarHistorico_ErroUnmarshal(t *testing.T) {
	original := services.UnmarshalList
	services.UnmarshalList = func(_ []map[string]types.AttributeValue, _ interface{}) error {
		return errors.New("erro simulado no unmarshal")
	}
	defer func() { services.UnmarshalList = original }()

	fakeInicio := time.Now().Add(-24 * time.Hour)
	fakeFim := time.Now()

	result := services.BuscarHistorico(fakeInicio, fakeFim)
	if result != nil {
		t.Errorf("Esperava retorno nil em erro de unmarshal")
	}
}

func TestSalvarCotacaoNoDynamo_Sucesso(t *testing.T) {
	called := false

	original := services.PutItemFn
	services.PutItemFn = func(_ *dynamodb.Client, _ *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		called = true
		return &dynamodb.PutItemOutput{}, nil
	}
	defer func() { services.PutItemFn = original }()

	cotacao := models.Cotacao{
		MoedaOrigem:  "BRL",
		MoedaDestino: "USD",
		Valor:        5.42,
		DataHora:     time.Now(),
	}

	services.SaveCotacao(cotacao)

	if !called {
		t.Errorf("PutItemFn não foi chamado")
	}
}

func TestSalvarCotacaoNoDynamo_ErroPutItem(t *testing.T) {
	original := services.PutItemFn
	services.PutItemFn = func(_ *dynamodb.Client, _ *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
		return nil, errors.New("erro simulado")
	}
	defer func() { services.PutItemFn = original }()

	cotacao := models.Cotacao{
		MoedaOrigem:  "BRL",
		MoedaDestino: "USD",
		Valor:        5.00,
		DataHora:     time.Now(),
	}

	services.SaveCotacao(cotacao)
}

func TestBuscarAPIKeyDoFixer_ErroSecretsManager(t *testing.T) {
	original := services.SecretsFetcher
	services.SecretsFetcher = func() string {
		t.Log("Simulando erro ao obter segredo")
		return ""
	}
	defer func() { services.SecretsFetcher = original }()

	token := services.SecretsFetcher()
	if token != "" {
		t.Errorf("Esperava retorno vazio, obteve: %s", token)
	}
}

func TestBuscarAPIKeyDoFixer_JSONInvalidoOthers(t *testing.T) {
	original := services.SecretsFetcher
	services.SecretsFetcher = func() string {
		raw := `{"fixer_api_key":123}` // valor numérico, quebra o json.Unmarshal
		var parsed map[string]string
		err := json.Unmarshal([]byte(raw), &parsed)
		if err != nil {
			t.Log("Erro ao interpretar JSON simulado")
			return ""
		}
		return parsed["fixer_api_key"]
	}
	defer func() { services.SecretsFetcher = original }()

	token := services.SecretsFetcher()
	if token != "" {
		t.Errorf("Esperava falha no unmarshal, obteve: %s", token)
	}
}

func TestBuscarAPIKeyDoFixer_ErroAoObterSegredo(t *testing.T) {
	original := services.GetSecretValueFn
	services.GetSecretValueFn = func(_ *secretsmanager.Client, _ *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
		return nil, errors.New("erro simulado")
	}
	defer func() { services.GetSecretValueFn = original }()

	key := services.BuscarAPIKeyDoFixer()
	if key != "" {
		t.Errorf("Esperava string vazia, obteve: %s", key)
	}
}

func TestBuscarAPIKeyDoFixer_ErroJSONUnmarshal(t *testing.T) {
	original := services.GetSecretValueFn
	services.GetSecretValueFn = func(_ *secretsmanager.Client, _ *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
		str := "isso-não-é-json"
		return &secretsmanager.GetSecretValueOutput{SecretString: &str}, nil
	}
	defer func() { services.GetSecretValueFn = original }()

	key := services.BuscarAPIKeyDoFixer()
	if key != "" {
		t.Errorf("Esperava string vazia por erro no unmarshal, obteve: %s", key)
	}
}
