# =============================================
# Makefile — Servidor MCP TCC (Google Docs ↔ LaTeX)
# =============================================
# Uso:
#   make build           → Compila para o SO atual
#   make build-windows   → Compila para Windows (amd64)
#   make build-linux     → Compila para Linux (amd64)
#   make build-mac       → Compila para macOS (amd64)
#   make build-all       → Compila para os 3 sistemas
#   make clean           → Remove binários gerados
#   make run             → Executa localmente via go run
#   make tidy            → Atualiza dependências do go

APP_NAME    := mcp-tcc
CMD_PATH    := ./cmd/mcp-server
BUILD_DIR   := bin
VERSION     := 1.0.0
LDFLAGS     := -s -w

ifeq ($(OS),Windows_NT)
    MKDIR_P           := if not exist $(BUILD_DIR) mkdir $(BUILD_DIR)
    RM_BIN            := if exist $(BUILD_DIR)\$(APP_NAME)* del /f /q $(BUILD_DIR)\$(APP_NAME)*
    BIN_CURRENT       := $(BUILD_DIR)/$(APP_NAME).exe
    BUILD_ENV_CURRENT := set CGO_ENABLED=0&&
    BUILD_ENV_WINDOWS := set CGO_ENABLED=0&& set GOOS=windows&& set GOARCH=amd64&&
    BUILD_ENV_LINUX   := set CGO_ENABLED=0&& set GOOS=linux&& set GOARCH=amd64&&
    BUILD_ENV_MAC     := set CGO_ENABLED=0&& set GOOS=darwin&& set GOARCH=amd64&&
else
    MKDIR_P           := mkdir -p $(BUILD_DIR)
    RM_BIN            := rm -f $(BUILD_DIR)/$(APP_NAME)*
    BIN_CURRENT       := $(BUILD_DIR)/$(APP_NAME)
    BUILD_ENV_CURRENT := CGO_ENABLED=0
    BUILD_ENV_WINDOWS := GOOS=windows GOARCH=amd64 CGO_ENABLED=0
    BUILD_ENV_LINUX   := GOOS=linux GOARCH=amd64 CGO_ENABLED=0
    BUILD_ENV_MAC     := GOOS=darwin GOARCH=amd64 CGO_ENABLED=0
endif

.PHONY: build build-windows build-linux build-mac build-all clean run tidy

# --- Build para o SO atual ---
build:
	@echo Compilando para o SO atual...
	@$(MKDIR_P)
	$(BUILD_ENV_CURRENT) go build -ldflags "$(LDFLAGS)" -o $(BIN_CURRENT) $(CMD_PATH)
	@echo Binario gerado com sucesso: $(BIN_CURRENT)

# --- Cross-compilation ---
build-windows:
	@echo Compilando para Windows (amd64)...
	@$(MKDIR_P)
	$(BUILD_ENV_WINDOWS) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME).exe $(CMD_PATH)
	@echo Binario gerado com sucesso: $(BUILD_DIR)/$(APP_NAME).exe

build-linux:
	@echo Compilando para Linux (amd64)...
	@$(MKDIR_P)
	$(BUILD_ENV_LINUX) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux $(CMD_PATH)
	@echo Binario gerado com sucesso: $(BUILD_DIR)/$(APP_NAME)-linux

build-mac:
	@echo Compilando para macOS (amd64)...
	@$(MKDIR_P)
	$(BUILD_ENV_MAC) go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin $(CMD_PATH)
	@echo Binario gerado com sucesso: $(BUILD_DIR)/$(APP_NAME)-darwin

build-all: build-windows build-linux build-mac
	@echo Todos os binarios foram gerados em $(BUILD_DIR)/

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
