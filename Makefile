# ====================================================
# Makefile — Servidor MCP TCC (Google Docs ↔ LaTeX)
# ====================================================
# Uso:
#   make / make build    → Compila para o SO atual (bin/mcp-tcc.exe no Windows)
#   make build-windows   → Compila explicitamente para Windows (amd64)
#   make build-linux     → Compila para Linux (amd64)
#   make build-mac       → Compila para macOS (amd64)
#   make build-all       → Compila para os 3 sistemas operacionais
#   make test            → Executa suíte de testes unitários
#   make clean           → Remove binários gerados em bin/
#   make run             → Executa localmente via go run
#   make inspect         → Abre o MCP Inspector no navegador para testar ferramentas via UI
#   make inspect-dev     → Abre o MCP Inspector executando diretamente via go run
#   make fmt             → Formata o código-fonte Go com go fmt
#   make tidy            → Atualiza e organiza dependências go.mod/go.sum

.DEFAULT_GOAL := build

APP_NAME    := mcp-tcc
CMD_PATH    := ./cmd/mcp-server
BUILD_DIR   := bin
VERSION     := 1.0.0
LDFLAGS     := -s -w -X main.serverVersion=$(VERSION)

ifeq ($(OS),Windows_NT)
    MKDIR_P           := if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
    RM_BIN            := if exist $(BUILD_DIR)\$(APP_NAME)* del /f /q $(BUILD_DIR)\$(APP_NAME)* && if exist $(BUILD_DIR)\*.exe~ del /f /q $(BUILD_DIR)\*.exe~
    BIN_CURRENT       := $(BUILD_DIR)/$(APP_NAME).exe
    BUILD_ENV_CURRENT := set CGO_ENABLED=0&&
    BUILD_ENV_WINDOWS := set CGO_ENABLED=0&& set GOOS=windows&& set GOARCH=amd64&&
    BUILD_ENV_LINUX   := set CGO_ENABLED=0&& set GOOS=linux&& set GOARCH=amd64&&
    BUILD_ENV_MAC     := set CGO_ENABLED=0&& set GOOS=darwin&& set GOARCH=amd64&&
else
    MKDIR_P           := mkdir -p $(BUILD_DIR)
    RM_BIN            := rm -f $(BUILD_DIR)/$(APP_NAME)* $(BUILD_DIR)/*.exe~
    BIN_CURRENT       := $(BUILD_DIR)/$(APP_NAME)
    BUILD_ENV_CURRENT := CGO_ENABLED=0
    BUILD_ENV_WINDOWS := GOOS=windows GOARCH=amd64 CGO_ENABLED=0
    BUILD_ENV_LINUX   := GOOS=linux GOARCH=amd64 CGO_ENABLED=0
    BUILD_ENV_MAC     := GOOS=darwin GOARCH=amd64 CGO_ENABLED=0
endif

.PHONY: build build-windows build-linux build-mac build-all test clean run tidy inspect inspect-dev fmt

# --- Build para o SO atual ---
build:
	@echo Compilando para o SO atual ($(APP_NAME) v$(VERSION))...
	@$(MKDIR_P)
	$(BUILD_ENV_CURRENT) go build -ldflags "$(LDFLAGS)" -o $(BIN_CURRENT) $(CMD_PATH)
	@echo Binario gerado com sucesso: $(BIN_CURRENT)

# --- Cross-compilation ---
build-windows:
	@echo Compilando para Windows amd64 ($(APP_NAME).exe v$(VERSION))...
	@$(MKDIR_P)
	$(BUILD_ENV_WINDOWS) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME).exe $(CMD_PATH)
	@echo Binario gerado com sucesso: $(BUILD_DIR)/$(APP_NAME).exe

build-linux:
	@echo Compilando para Linux amd64 ($(APP_NAME)-linux v$(VERSION))...
	@$(MKDIR_P)
	$(BUILD_ENV_LINUX) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux $(CMD_PATH)
	@echo Binario gerado com sucesso: $(BUILD_DIR)/$(APP_NAME)-linux

build-mac:
	@echo Compilando para macOS amd64 ($(APP_NAME)-darwin v$(VERSION))...
	@$(MKDIR_P)
	$(BUILD_ENV_MAC) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin $(CMD_PATH)
	@echo Binario gerado com sucesso: $(BUILD_DIR)/$(APP_NAME)-darwin

build-all: build-windows build-linux build-mac
	@echo Todos os binarios foram gerados em $(BUILD_DIR)/

# --- Testes unitários ---
test:
	@echo Executando testes unitarios...
	go test ./internal/...

# --- Limpeza ---
clean:
	@echo Removendo binarios compilados...
	@$(RM_BIN)
	@echo Limpeza concluida.

# --- Execucao local ---
run:
	@echo Executando servidor MCP localmente...
	go run $(CMD_PATH)

# --- Dependencias ---
tidy:
	@echo Resolvendo dependencias...
	go mod tidy
	@echo Dependencias atualizadas.

# --- Formatacao do codigo Go ---
fmt:
	@echo Formatando codigo-fonte Go...
	go fmt ./...

# --- MCP Inspector (Interface web interativa para testar e clicar nas ferramentas) ---
inspect: build
	@echo Abrindo o MCP Inspector no navegador para interagir com o binario...
	npx @modelcontextprotocol/inspector $(BIN_CURRENT)

inspect-dev:
	@echo Abrindo o MCP Inspector no navegador com go run...
	npx @modelcontextprotocol/inspector go run $(CMD_PATH)
