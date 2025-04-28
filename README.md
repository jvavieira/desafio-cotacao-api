# API de Cotações BRL → USD

Este projeto é parte de um desafio técnico para a vaga de Cloud. Ele consiste em uma aplicação que consulta a cotação do dólar (USD) em relação ao real (BRL), armazena os dados em um banco de dados e disponibiliza endpoints REST para consulta.

## Escolha da Linguagem Go

Apesar de minha experiência sólida com PL/SQL e Java (Spring Boot), sempre tive curiosidade em explorar o Go. Por ser uma linguagem procedural, assim como o PL/SQL, identifiquei similaridades conceituais que poderiam facilitar o aprendizado.

Além disso, o Go é conhecido por sua leveza, simplicidade e excelente desempenho em aplicações REST. Considerando que o projeto envolvia um serviço pequeno e direto, optei por não utilizar Java, que depende da JVM e tende a ser mais pesado para esse tipo de caso.

Outro ponto importante foi o aspecto "Cloud Native" do Go. A linguagem é amplamente utilizada em ferramentas e serviços desenvolvidos pela própria AWS, o que reforça seu alinhamento com práticas modernas de cloud computing.

Essa decisão foi estratégica não apenas por curiosidade técnica, mas por acreditar que o Go seria uma escolha prática, eficiente e adequada para os objetivos e escala deste desafio.

## Arquitetura do Projeto

A arquitetura proposta utiliza os seguintes serviços e ferramentas:

- **Amazon App Runner**: serviço responsável por executar a aplicação principal (API REST) a partir de uma imagem Docker.
- **Amazon ECR (Elastic Container Registry)**: repositório onde são armazenadas as imagens Docker da aplicação e da função Lambda.
- **Amazon DynamoDB**: banco de dados NoSQL utilizado para armazenar as cotações obtidas.
- **AWS Lambda** *(não finalizado)*: função que poderá futuramente buscar e salvar as cotações de forma automatizada.
- **Amazon EventBridge** *(não finalizado)*: utilizado para agendamento de execuções da função Lambda.

A infraestrutura foi provisionada totalmente via **Terraform**, garantindo rastreabilidade e consistência no deploy dos recursos.

## Infraestrutura como Código (IaC)

A infraestrutura do projeto é gerenciada utilizando Terraform e está organizada na pasta `terraform/`. Os principais componentes criados são:

- **DynamoDB**: tabela para armazenar as cotações com `data_hora` como chave primária.
- **ECR**: repositórios separados para a aplicação principal e a função Lambda.
- **App Runner**: serviço que consome a imagem Docker publicada no ECR e executa a API.
- **Secrets Manager**: armazena a chave da API utilizada para buscar as cotações.
- **Lambda e EventBridge** *(não finalizado)*: para execução automatizada e agendada da coleta de dados.

### Execução
Para aplicar a infraestrutura, siga os seguintes passos:

```bash
cd terraform
terraform init
terraform apply -var="fixer_api_key_value=SUA_CHAVE_DO_FIXER"
```
⚠️ A chave da API do ApiLayer é passada via variável no momento do apply, e não está escrita em nenhum código-fonte, mantendo a segurança das credenciais.
⚠️ Gerar token através de longin gratuito no site https://apilayer.com/signup.
⚠️ Certifique-se de que as credenciais da AWS estejam configuradas corretamente.

## API de Cotações

A aplicação em Go expõe dois endpoints REST principais via Amazon App Runner:

### 1. `GET /cotacao/ultima`
Retorna a cotação BRL → USD mais recente consultada via API externa e salva no DynamoDB.

#### Exemplo de resposta:
```json
{
  "moeda_origem": "BRL",
  "moeda_destino": "USD",
  "valor": 5.19,
  "data_hora": "2025-04-21T14:00:00Z"
}
```

### 2. `GET /cotacao/historico?inicio=YYYY-MM-DDTHH:mm&fim=YYYY-MM-DDTHH:mm`
Consulta o histórico de cotações dentro de um intervalo de datas.

#### Parâmetros:
- `inicio`: data/hora inicial (ex: `2025-04-20T00:00`)
- `fim`: data/hora final (ex: `2025-04-22T23:59`)

#### Exemplo de resposta:
```json
[
  {
    "moeda_origem": "BRL",
    "moeda_destino": "USD",
    "valor": 5.20,
    "data_hora": "2025-04-20T12:00:00Z"
  },
  {
    "moeda_origem": "BRL",
    "moeda_destino": "USD",
    "valor": 5.19,
    "data_hora": "2025-04-21T14:00:00Z"
  }
]
```

## Deploy via App Runner

A aplicação é empacotada em uma imagem Docker e enviada ao Amazon Elastic Container Registry (ECR). O serviço App Runner é responsável por executar a imagem e disponibilizar os endpoints públicos.

## Etapas:

1. O Dockerfile na raiz do projeto define a build da aplicação Go.

2. Um workflow no GitHub Actions (lambda-deploy.yml) realiza:

  - Checkout do código
  - Testes unitários
  - Análise de vulnerabilidades com Trivy
  - Build da imagem
  - Push da imagem para o ECR

3. O Terraform provisiona o serviço App Runner apontando para a imagem mais recente do repositório.

## Exemplo de execução do workflow:

 - Vá até a aba Actions no GitHub
 - Escolha o workflow Deploy Lambda Container (nome pode variar)
 - Clique em Run workflow

O serviço App Runner será automaticamente atualizado com a nova imagem da aplicação.

## EventBridge + Lambda *(melhoria não concluída)*

O projeto previa a inclusão de uma função Lambda acoplada ao Amazon EventBridge, com objetivo de automatizar a coleta de cotações em horários agendados. Apesar da estrutura ter sido parcialmente implementada via Terraform, a integração final não foi concluída devido à limitação de tempo, onde o teste realizado manualmente no Console, não retornou o resultado esperado.

## CI/CD da API principal

Além do pipeline da Lambda, foi implementado um pipeline dedicado à API principal.

### Funcionalidades do CI/CD (via GitHub Actions):
- Build da aplicação com Docker
- Testes unitários com Go
- Análise de vulnerabilidades com Trivy
- Deploy automatizado da imagem para o Amazon ECR

### Arquivo do workflow: `.github/workflows/api-deploy.yml`

A execução pode ser feita manualmente na aba **Actions** do GitHub ou ser configurada para rodar em push para a branch `main`.

## Instruções de Teste e Submissão

### Testar a API

1. Acesse o endpoint público fornecido pelo App Runner:

Ex: https://<endpoint>.awsapprunner.com/cotacao/ultima

Ex: https://<endpoint>.awsapprunner.com/cotacao/historico?inicio=2025-04-20T00:00&fim=2025-04-22T23:59

2. Utilize ferramentas como Postman, Insomnia ou curl:
```bash
curl https://<seu-endpoint>/cotacao/ultima
```

### Repositório Público
> 🔗 https://github.com/jvavieira/desafio-cotacao-api

---
Caso a banca avaliadora deseje testar manualmente, todas as instruções estão descritas acima. Obrigado!
