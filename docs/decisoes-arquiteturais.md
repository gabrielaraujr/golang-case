# Decisões Arquiteturais
## Ordem natural

1. **Cliente envia proposta**
2. **Análise de Documentos**
3. **Análise de Crédito**
4. **Análise de Fraude**
5. **Decisão Final**
---
## Anotações sobre processos

- Reprovar em qualquer análise **encerrará a proposta**, simplificando o fluxo de decisões
- Análise de Documentos é essencial processar antes de prosseguir para outras análises
- Análise de Crédito e Fraude podem rodar **em paralelo**
	- Atrasar abertura de conta é piorar experiência do cliente
- Apenas um serviço incluindo todas as análises, pois todos participam de um domínio de risco, já que todos apresentam risco à empresa
- **Serviço da Proposta** não executa análises mas envia notificações e aprovação final com base no resultado delas
- **Serviço de Análise** não notifica o cliente diretamente, apenas retorna o estado para o **Serviço da Proposta**

---
## Serviço de Proposta responsável por notificar

Extrair um microsserviço só para notificação nesse case seria overengineering. Levo em consideração:

- Notificação não vai ser um domínio tão complexo
- Não existe tantos serviços
- Não existe diversos times separados

### Benefícios

- Menor custo operacional
- Menos serviços para manter
- Menor latência
- Fluxo mais simples
- Deploys mais rápidos
- Observabilidade centralizada

### Trade-offs

Poderia avaliar decisão de separar Notificação em um serviço isolado só para casos de existir diversos times em projetos separados que necessitam disso. 

Separar poderia trazer mais autonomia na evolução dos produtos, quando cada um quer **configuração própria**; múltiplos canais (e-mail, SMS, push); ou requisitos não funcionais próprios (DQL, retries).

---
## Serviço de Análise de Risco (Todas as análises)

Levo em consideração que todas as análises realizadas durante a proposta tem um objetivo em comum: **evitar riscos à instituição**. Dado o escopo do case, prazo pequeno e sendo o único dono, juntar validação de documentos, análise de crédito e análise de fraude em um único serviço é a **arquitetura adequada ao contexto**, evitando overengineering.

Optei por concentrar **todas as análises** em um único serviço, pois:

- Documentos, crédito e fraude fazem parte do mesmo domínio de **risco**
- Não há necessidade de evolução ou escalabilidade independente entre essas análises
- O fluxo possui uma ordem natural interna
- Crédito e fraude podem executar em paralelo após identidade validada
- Reduz custo operacional, número de serviços e coordenação entre sistemas

### Benefícios

- Menor custo operacional
- Menos serviços para manter
- Menor complexidade arquitetural
- Menos eventos e filas
- Execução paralela interna das análises de crédito e fraude
- Decisão de risco centralizada
- Observabilidade mais simples

### Trade-offs

A separação das análises em serviços isolados pode ser considerada no futuro caso exista:

- Escala significativamente diferente entre validação de documentos e análises de risco
- Times distintos com donos separados
- Alta evolução independente

---
## Análise de Crédito e Fraude em paralelo

Crédito e fraude executam em paralelo **dentro do mesmo serviço**, após validação dos documentos do cliente.

Caso qualquer uma das análises resulte em reprovação, o Serviço de Análise de Risco publica imediatamente um evento de reprovação.

Ao consumir esse evento, o Serviço de Proposta encerra o fluxo da proposta, atualizando seu estado para reprovado e notificando o cliente, sem a necessidade de aguardar a conclusão das demais análises.

Seguimos optando por:

- Simplicidade
- Menos acoplamento entre serviços
- Menos estado compartilhado
- Menos coordenação distribuída

---
## Microsserviços e responsabilidades

1. **Serviço de Proposta** (Central)
	
	Responsável por **coordenar todo o fluxo com base nos estados**.
	
	- Receber proposta via API HTTP
	- Persistir estado da proposta
	- Publicar evento de criação de proposta
	- Gerenciar ciclo de vida da proposta (state machine)
	- Receber resultado das análises externas
	- Decisão final da proposta
	- Notificar cliente a cada mudança de estado
	
	**Observações:**
	- Não executa análise
	- Atua como gerenciar do estado da proposta
	- Responsável por notificar acompanhamento do processo ao cliente
	
2. **Serviço de Análise de Risco** (Documentos, Crédito e Fraude)
	
	Responsável por **evitar riscos financeiros e fraude**, tratando validação de identidade, crédito e de fraude no mesmo domínio.
	
	- Consumir evento de criação da proposta
	- Realizar análise de documentos
	- Realizar análise de crédito
	- Realizar análise de fraude em paralelo com a análise de crédito
	- Organizar os resultados das análises
	- Publicar os resultados das análises
	
	**Observações:**
	- Validação de documentos é obrigatória e bloqueante
	- Crédito e fraude executam em paralelo após documentos aprovados
	- Reprovação em qualquer etapa encerra o fluxo da proposta
	- Não notifica o cliente, apenas o serviço de proposta

---
## Resumo arquitetural

- Arquitetura em microsserviços
- Serviço de Proposta atua como coordenador por estado
- Notificações embutidas no Serviço de Proposta
- Comunicação assíncrona via eventos para garantir desacoplamento
- Análise de documentos é etapa obrigatória e bloqueante
- Análises de crédito e fraude executam em paralelo
- Reprovação em qualquer análise encerra o fluxo
- Sem retentativas nem fluxos alternativos
- Trade-offs de custo computacional são aceitos para reduzir acoplamento e complexidade
- Arquitetura feita para ser simples, evolutiva e adequada ao contexto de escopo do case