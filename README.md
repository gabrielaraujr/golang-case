# Proposal

* [Overview](#overview)
    - [Domínio](#domain)
* [Instalação](#installation)
    - [Repositório](#repository)
    - [Configuração](#configure)
* [Roadmap](#roadmap)
    - [Verificando o ambiente](#checking)
    - [Executando o caso de uso](#run_use_case)
    - [Consultando status](#consult_use_case)
* [Regras de Análise](#analyze_rules)
    - [Documentos](#document)
    - [Crédito](#credit)
    - [Fraude](#fraud)
* [Testando cenários](#testing)
* [Fluxo de Estados](#states)
* [Monitoramento](#monitoring)
* [Estrutura](#structure)

<!--suppress HtmlDeprecatedAttribute -->

<div id="overview"></div>

## Overview

Case de sistema de captura e análise de propostas para abertura de contas.

<div id='domain'></div>

### Domínio

O sistema é composto por dois microsserviços:

- **Account Service**: API REST que gerencia propostas e coordena o fluxo através de eventos.
- **Risk Analysis Service**: Consumer de eventos que executa análises de documentos, crédito e fraude.

Para decisões arquiteturais detalhadas, veja:
- [Decisões Arquiteturais](docs/decisoes-arquiteturais.md)
- [Decisões Técnicas](docs/decisoes-tecnicas.md)

<div id='installation'></div>

## Instalação

<div id='repository'></div>

### Repositório

Clone o repositório usando a linha de comando:

```bash
git clone https://github.com/gabrielaraujr/golang-case.git
```

<div id='configure'></div>

### Configuração

Verifique se algum processo usa as portas: **4566**, **5432**, **8001**. Se alguma das portas estiver em uso, vai precisar libera-los.

Para instalar e configurar o projeto, execute na raiz do projeto:

```bash
make configure
```

Aguarde até todos os containers estarem prontos (~30 segundos).

<div id='roadmap'></div>

## Roadmap

<div id='checking'></div>

### Verificando o ambiente

Para executar o caso de uso, basta estar com o ambiente docker inicializado.

<div id='run_use_case'></div>

### Executando o caso de uso

Pode executar com a linha de comando:

```bash
make create-proposal
```

Ou se quiser fazer manualmente:

```bash
curl -X POST http://localhost:8001/proposals \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Gabriel Silva",
    "cpf": "12345678900",
    "salary": 5000.00,
    "email": "gabriel@email.com",
    "phone": "11999999999",
    "birthdate": "15-01-1990",
    "address": {
      "street": "Rua Teste 123",
      "city": "São Paulo",
      "state": "SP",
      "zip_code": "01234567"
    }
  }'
```

Guarde o `id` retornado na resposta.

<div id='consult_use_case'></div>

### Consultar status da proposta

```bash
curl http://localhost:8001/proposals/{id}
```

Aguarde 5-10 segundos para o processamento completo.

<div id='analyze_rules'></div>

## Regras de Análise

<div id='document'></div>

### Documentos
- **Aprovado**: CPF com 11 dígitos e nome com 3+ caracteres
- **Rejeitado**: CPF inválido ou nome muito curto

<div id='credit'></div>

### Crédito
- **Aprovado**: Salário > R$ 3.000,00
- **Rejeitado**: Salário ≤ R$ 3.000,00

<div id='fraud'></div>

### Fraude
- **Aprovado**: Último dígito do CPF é par
- **Rejeitado**: Último dígito do CPF é ímpar

<div id='testing'></div>

## Testando cenários

### Proposta aprovada

```bash
# CPF com último dígito par (fraude OK) + salário > 3k (crédito OK)
curl -X POST http://localhost:8001/proposals \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Maria Santos",
    "cpf": "98765432100",
    "salary": 5000.00,
    "email": "maria@email.com",
    "phone": "11988888888",
    "birthdate": "20-03-1985",
    "address": {
      "street": "Av Paulista 1000",
      "city": "São Paulo",
      "state": "SP",
      "zip_code": "01310100"
    }
  }'
```

**Status final:** `approved`

### Proposta rejeitada (fraude)

```bash
# CPF com último dígito ímpar = fraude detectada
curl -X POST http://localhost:8001/proposals \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "João Lima",
    "cpf": "11122233341",
    "salary": 5000.00,
    "email": "joao@email.com",
    "phone": "11977777777",
    "birthdate": "10-05-1990",
    "address": {
      "street": "Rua ABC 789",
      "city": "Rio de Janeiro",
      "state": "RJ",
      "zip_code": "20000000"
    }
  }'
```

**Status final:** `rejected`

### Proposta rejeitada (crédito)

```bash
# Salário baixo (≤ 3000)
curl -X POST http://localhost:8001/proposals \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Pedro Costa",
    "cpf": "55566677788",
    "salary": 2000.00,
    "email": "pedro@email.com",
    "phone": "11966666666",
    "birthdate": "01-12-1992",
    "address": {
      "street": "Rua XYZ 456",
      "city": "Curitiba",
      "state": "PR",
      "zip_code": "80000000"
    }
  }'
```

**Status final:** `rejected`

### Proposta rejeitada (documentos)

```bash
# CPF com menos de 11 dígitos
curl -X POST http://localhost:8001/proposals \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Ana Silva",
    "cpf": "123",
    "salary": 5000.00,
    "email": "ana@email.com",
    "phone": "11955555555",
    "birthdate": "25-08-1988",
    "address": {
      "street": "Rua Teste 999",
      "city": "Belo Horizonte",
      "state": "MG",
      "zip_code": "30000000"
    }
  }'
```

**Status final:** `rejected`

<div id='states'></div>

## Fluxo de Estados

```
pending → analyzing → approved/rejected
```

- **pending**: *Proposta criada* -> aguardando para análise
- **analyzing**: Análises em andamento
- **approved**: Todas as análises aprovadas
- **rejected**: Alguma análise reprovou

<div id='monitoring'></div>

## Monitoramento

```bash
# Ver logs
make logs

# Verificar filas
make check-queue        # Fila de propostas
make check-results      # Fila de análise de risco

# Rodar testes
make test

# Remover containers e limpar cache
make clean
```

<div id='structure'></div>

## Estrutura

```
account/          # Serviço de gestão de propostas (API REST)
risk-analysis/    # Serviço de análise de risco (Consumer)
docs/             # Documentação arquitetural
```
