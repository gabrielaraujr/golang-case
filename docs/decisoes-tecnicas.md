# Decisões Técnicas

<span style="color: gray">Com base nas</span> [Decisões Arquiteturais](decisoes-arquiteturais.md)
## Linguagem

- **Golang**
	Motivos:
    - Linguagem já definida pelo contexto da vaga
    - Boa performance e baixo consumo de recursos
    - Forte aderência a microsserviços e aplicações cloud-native
    - Ecossistema maduro para APIs HTTP e processamento assíncrono

---
## Comunicação assíncrona e eventos

- **Apenas SQS**
    Motivos:
    - Consumo desacoplado por serviço
    - Evita dependência entre serviços
    - Facilita paralelismo e resiliência
    - Simples e suficiente para o escopo do case

---
## Persistência de dados

- **Banco relacional (MySQL ou PostgreSQL)**
- Utilizado apenas pelo **Serviço de Proposta**
	Motivos:
    - Facilidade de modelagem e entendimento
    - Garantia de consistência no estado da proposta
    - Fluxo não possui tantas requisições nem é altamente concorrente
    - Saldos e transações financeiras não fazem parte do escopo

---
## Segurança

- **OAuth 2.0 + JWT (conceitual)**
	Motivos:
    - Autenticação desacoplada das APIs
    - Uso de token JWT para chamadas HTTP
    - Serviços validam token
    - Adequado para arquitetura em microsserviços

---

## Princípios 12-Factor

- Serviços sem estado (Estado salvo em banco)
- Configurações via variáveis de ambiente
- Logs (observabilidade centralizada)
- Serviços externos tratados de maneira controlada (Tratar erros em caso de timetout ou indisponibilidade)

---
## Reação a eventos (Publica / Consome)

### Serviço de Proposta

**Consome:**
- `DocumentsApproved`
- `DocumentsRejected`
- `CreditApproved`
- `CreditRejected`
- `FraudApproved`
- `FraudRejected`
- `RiskAnalysisCompleted`

**Publica:**
- `ProposalCreated`
- `ProposalStatusChanged`
- `ProposalApproved`
- `ProposalRejected`

---
### Serviço de Análise de Risco

**Consome:**
- `ProposalCreated`

**Publica:**
- `DocumentsApproved`
- `DocumentsRejected`
- `CreditApproved`
- `CreditRejected`
- `FraudApproved`
- `FraudRejected`
- `RiskAnalysisCompleted`

---
## Resumo técnico

- Comunicação assíncrona orientada a eventos
- Dois microsserviços apenas
- Serviço de Proposta coordena estado e comunicação
- Serviço de Análise de Risco executa todas as análises
- Simplicidade priorizada sobre granularidade excessiva
- Trade-offs explícitos e conscientes
- Desenvolvimento local usando LocalStack e Docker