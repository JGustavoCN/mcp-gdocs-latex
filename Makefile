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

APP_NAME    := mcp-tcc
CMD_PATH    := ./cmd/mcp-server
BUILD_DIR   := bin
VERSION     := 1.0.0
LDFLAGS     := -s -w

.PHONY: build build-windows build-linux build-mac build-all clean run tidy

# --- Build para o SO atual ---
build:
	@echo "🔨 Compilando para o SO atual..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_PATH)
	@echo "✅ Binário gerado: $(BUILD_DIR)/$(APP_NAME)"

# --- Cross-compilation ---
build-windows:
	@echo "🪟 Compilando para Windows (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME).exe $(CMD_PATH)
	@echo "✅ Binário gerado: $(BUILD_DIR)/$(APP_NAME).exe"

build-linux:
	@echo "🐧 Compilando para Linux (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux $(CMD_PATH)
	@echo "✅ Binário gerado: $(BUILD_DIR)/$(APP_NAME)-linux"

build-mac:
	@echo "🍎 Compilando para macOS (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin $(CMD_PATH)
	@echo "✅ Binário gerado: $(BUILD_DIR)/$(APP_NAME)-darwin"

build-all: build-windows build-linux build-mac
	@echo ""
	@echo "🎉 Todos os binários foram gerados em $(BUILD_DIR)/"

# --- Limpeza ---
clean:
	@echo "🧹 Removendo binários..."
	rm -rf $(BUILD_DIR)
	@echo "✅ Limpo."

# --- Execução local ---
run:
	@echo "🚀 Executando servidor MCP localmente..."
	go run $(CMD_PATH)

# --- Dependências ---
tidy:
	@echo "📦 Resolvendo dependências..."
	go mod tidy
	@echo "✅ Dependências atualizadas."
