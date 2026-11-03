package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
)

func main() {
	docID := "1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8"
	if len(os.Args) > 1 {
		docID = os.Args[1]
	}

	contentPath := `C:\Users\joseg\.gemini\antigravity\brain\0c400b2f-aaca-4279-9a2f-8ec020311d37\scratch\cap4_limpo.md`
	bytes, err := os.ReadFile(contentPath)
	if err != nil {
		log.Fatalf("Erro ao ler %s: %v", contentPath, err)
	}

	content := string(bytes)
	log.Printf("[INFO] Iniciando injeção do Capítulo 4 (%d caracteres) no documento %s...", len(content), docID)

	ctx := context.Background()
	msg, err := gdocs.SyncSectionToDoc(ctx, docID, "4 DESENVOLVIMENTO DA SOLUÇÃO", content)
	if err != nil {
		log.Fatalf("[ERRO] Falha em SyncSectionToDoc: %v", err)
	}

	fmt.Println(msg)
}
