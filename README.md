# Stress Test CLI - Desafio Técnico FullCycle

[![Go Version](https://img.shields.io/badge/Go-1.25-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-ready-brightgreen.svg)](https://www.docker.com/)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](https://github.com)

Sistema CLI em Go para realizar testes de carga em serviços web. Permite simular múltiplas requisições HTTP simultâneas e gerar relatórios detalhados sobre o desempenho do serviço testado.

## 📋 Funcionalidades

- ✅ Testes de carga HTTP com concorrência configurável
- ✅ Relatório detalhado com métricas de performance
- ✅ Distribuição de códigos de status HTTP
- ✅ Medição de tempo de execução
- ✅ Execução via Docker
- ✅ Testes unitários completos

## 🚀 Tecnologias

- Go 1.25.5
- Docker
- HTTP Client nativo

## 📦 Instalação

### Via Go

```bash
# Clone o repositório
git clone https://github.com/adalbertofjr/desafio-tecnico-stress-test.git
cd desafio-tecnico-stress-test

# Execute diretamente
go run ./cmd/stresstest/main.go --url=<URL> --requests=<N> --concurrency=<C>
```

### Via Docker

```bash
# Build da imagem
docker build -t stresstest .

# Execute com Docker
docker run --rm stresstest --url=<URL> --requests=<N> --concurrency=<C>
```

## 🔧 Uso

### Parâmetros CLI

| Parâmetro | Descrição | Obrigatório |
|-----------|-----------|-------------|
| `--url` | URL do serviço a ser testado | Sim |
| `--requests` | Número total de requisições | Sim |
| `--concurrency` | Número de requisições simultâneas | Sim |

### Exemplos

**Teste básico:**
```bash
go run ./cmd/stresstest/main.go \
  --url=http://localhost:8080/health \
  --requests=100 \
  --concurrency=10
```

**Teste de alta carga:**
```bash
go run ./cmd/stresstest/main.go \
  --url=https://api.exemplo.com/endpoint \
  --requests=10000 \
  --concurrency=100
```

**Via Docker (acessando localhost da máquina host):**
```bash
docker run --rm stresstest \
  --url=http://host.docker.internal:8080/health \
  --requests=1000 \
  --concurrency=50
```

**Via Docker (acessando serviço externo):**
```bash
docker run --rm stresstest \
  --url=https://google.com \
  --requests=500 \
  --concurrency=25
```

## 📊 Relatório Gerado

O sistema gera um relatório detalhado ao final da execução:

```
=========================================
Relatório de Stress Test
=========================================
URL Testada: http://localhost:8080/health
Total de Requisições: 1000
Requisições Executadas: 1000
Duração Total: 5.23s
=========================================
Detalhes das Respostas:
  Código (200): 950 Sucesso
  Código (404): 30 Falha
  Código (500): 20 Falha
=========================================
```

### Métricas incluídas:
- **URL Testada**: Endpoint que recebeu as requisições
- **Total de Requisições**: Número configurado de requests
- **Requisições Executadas**: Confirmação de execução
- **Duração Total**: Tempo total em segundos
- **Distribuição de Status**: Contagem por código HTTP

## 🧪 Testes

Execute os testes unitários:

```bash
# Executar todos os testes
go test ./stresstest/

# Com detalhes verbose
go test -v ./stresstest/

# Com cobertura
go test -cover ./stresstest/

# Teste específico
go test -v -run TestExecute ./stresstest/
```

### Cobertura de Testes

- ✅ Testes de concorrência (race conditions)
- ✅ Contagem de status codes
- ✅ Incremento thread-safe de contadores
- ✅ Execução de requisições HTTP
- ✅ Geração de relatórios
- ✅ Validação de múltiplos códigos de erro

## 🐳 Docker

### Build

```bash
docker build -t stresstest .
```

### Executar

```bash
# Com remoção automática do container
docker run --rm stresstest \
  --url=http://host.docker.internal:8080/api \
  --requests=1000 \
  --concurrency=10
```

### Notas sobre Docker e Networking

- Para acessar `localhost` da máquina host use `host.docker.internal` (macOS/Windows)
- No Linux, use `--network=host` ou o IP da máquina
- Para serviços externos, use a URL normalmente

## 🏗️ Arquitetura

### Estrutura de Diretórios

```
.
├── cmd/
│   └── stresstest/
│       └── main.go           # Entry point da aplicação
├── stresstest/
│   ├── stresstest.go         # Lógica principal de testes
│   ├── stresstest_test.go    # Testes da lógica principal
│   ├── report.go             # Geração de relatórios
│   └── report_test.go        # Testes dos relatórios
├── Dockerfile                # Containerização
├── .dockerignore
├── go.mod
└── README.md
```

### Componentes Principais

#### 1. **CLI (main.go)**
- Processa argumentos via `flag`
- Valida parâmetros obrigatórios
- Executa o teste de carga

#### 2. **Executor (stresstest.go)**
- Gerencia concorrência com **channel semáforo**
- Controla número de goroutines simultâneas
- Coleta métricas de execução
- Thread-safe com `sync.Mutex`

#### 3. **Gerador de Relatórios (report.go)**
- Formata métricas coletadas
- Exibe distribuição de status HTTP
- Calcula tempo total de execução

### Fluxo de Funcionamento

```
┌─────────────┐
│  Parse CLI  │  --url, --requests, --concurrency
└──────┬──────┘
       │
       ▼
┌─────────────────────┐
│  Validate Params    │  Verifica obrigatoriedade
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│  Execute()          │
│  - Inicia timer     │
│  - Cria semáforo    │
│  - Spawna goroutines│
└──────┬──────────────┘
       │
       ├─► Goroutine 1 ──► HTTP GET ──► Incrementa counter
       ├─► Goroutine 2 ──► HTTP GET ──► Registra status
       ├─► Goroutine N ──► HTTP GET ──► Thread-safe
       │
       ▼
┌─────────────────────┐
│  WaitGroup.Wait()   │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│  PrintReport()      │  Exibe métricas
└─────────────────────┘
```

### Concorrência Implementada

O sistema usa um **semáforo baseado em channel** para controlar a concorrência de forma eficiente:

```go
sem := make(chan struct{}, concurrency)
for range requests {
    wg.Add(1)
    sem <- struct{}{}  // Bloqueia se atingir limite de concorrência
    
    go func() {
        defer func() { <-sem }()  // Libera slot ao finalizar
        makeRequest(...)           // Executa requisição HTTP
    }()
}
wg.Wait()
```

**Vantagens:**
- ✅ Limita goroutines ativas simultaneamente
- ✅ Evita sobrecarga de memória
- ✅ Backpressure natural (bloqueia se cheio)
- ✅ Simplicidade (sem libraries externas)

### Thread Safety

- **`sync.Mutex`**: Protege contadores compartilhados
- **`sync.WaitGroup`**: Sincroniza término das goroutines
- **Incremento atômico**: Lock/unlock em todas as escritas

## 🎯 Requisitos do Desafio

- [x] Sistema CLI em Go
- [x] Parâmetros via CLI (--url, --requests, --concurrency)
- [x] Execução de requests HTTP
- [x] Distribuição com concorrência
- [x] Relatório com tempo total
- [x] Relatório com total de requests
- [x] Relatório com status HTTP 200
- [x] Relatório com distribuição de outros códigos
- [x] Execução via Docker

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'feat: Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

### Commits Semânticos

Seguimos [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` Nova funcionalidade
- `fix:` Correção de bug
- `docs:` Mudanças na documentação
- `test:` Adição ou correção de testes
- `refactor:` Refatoração de código
- `perf:` Melhoria de performance
- `chore:` Tarefas de manutenção

## 📄 Licença

Este projeto foi desenvolvido como desafio técnico para o curso de Pós-Graduação em Golang da FullCycle.

## 👨‍💻 Autor

Desenvolvido por **Adalberto Fernandes Jr.**

## 🙏 Agradecimentos

- FullCycle - Pós-Graduação Go Expert
- Comunidade Go Brasil

---

**Nota:** Este é um projeto educacional desenvolvido para demonstração de conceitos de testes de carga, concorrência e CLI em Go.
