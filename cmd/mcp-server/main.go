// Servidor MCP — TCC Google Docs ↔ LaTeX
//
// Este é o entrypoint do servidor. Ele inicializa o servidor MCP,
// registra todas as ferramentas e inicia o transporte via stdio.
//
// IMPORTANTE: Todos os logs vão para stderr para não conflitar
// com o protocolo MCP que utiliza stdout.
package main

import (
	"log"
	"os"

	mcptools "github.com/JGustavoCN/mcp-gdocs-latex/internal/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "TCC_DocsLaTeX_Bridge"
	serverVersion = "1.0.0"
)

func main() {
	// Configura logging para stderr (o protocolo MCP usa stdout para JSON-RPC)
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ldate | log.Ltime)

	log.Println("============================================")
	log.Println("  Servidor MCP — TCC Google Docs ↔ LaTeX")
	log.Printf("  Versão: %s", serverVersion)
	log.Println("============================================")

	// Log de configuração
	if allowedIDs := os.Getenv("ALLOWED_DOC_IDS"); allowedIDs != "" {
		log.Printf("[INFO] ALLOWED_DOC_IDS configurado — acesso restrito ativado")
	} else {
		log.Println("[INFO] ALLOWED_DOC_IDS não definido — acesso livre a todos os documentos compartilhados")
	}

	// Cria o servidor MCP
	s := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithToolCapabilities(false),
	)

	// Registra todas as ferramentas
	mcptools.RegisterTools(s)

	log.Println("[INFO] Servidor MCP iniciado. Aguardando conexões via stdio...")

	// Inicia o transporte via stdin/stdout
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("[FATAL] Erro no servidor MCP: %v", err)
	}
}
