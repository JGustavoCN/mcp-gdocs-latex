//go:build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
	"google.golang.org/api/docs/v1"
)

func main() {
	ctx := context.Background()
	docID := "1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8"

	svc, err := gdocs.GetDocsService(ctx)
	if err != nil {
		log.Fatalf("Erro svc: %v", err)
	}

	doc, err := svc.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalf("Erro doc: %v", err)
	}
	_ = doc

	// 1. Dividir parágrafos longos inserindo quebra de parágrafo dupla (\n\n) no ponto exato
	replacements := []struct {
		target string
		repl   string
	}{
		// P1: MCP introdução
		{
			target: "inconsistências de dados (ANTHROPIC, 2024; RAY, 2025). Para endereçar essa fragmentação",
			repl:   "inconsistências de dados (ANTHROPIC, 2024; RAY, 2025).\n\nPara endereçar essa fragmentação",
		},
		// P2: MCP entidades e transportes
		{
			target: "sinalizações formais de erro (JSON-RPC WORKING GROUP, 2010; RAY, 2025). Para o transporte dessas mensagens",
			repl:   "sinalizações formais de erro (JSON-RPC WORKING GROUP, 2010; RAY, 2025).\n\nPara o transporte dessas mensagens",
		},
		// P3: MCP protótipo node e transportes
		{
			target: "persistência em disco (ZULLO et al., 2025). Na integração interativa com o Claude Desktop",
			repl:   "persistência em disco (ZULLO et al., 2025).\n\nNa integração interativa com o Claude Desktop",
		},
		// P4: Yuri no IFES
		{
			target: "treze testes de conversação tabulados. Essa avaliação confrontou as saídas",
			repl:   "treze testes de conversação tabulados.\n\nEssa avaliação confrontou as saídas",
		},
		// Enumeração dos objetivos específicos (conforme LaTeX \begin{enumerate} e ABNT NBR 6024)
		{
			target: "• Realizar consultas exploratórias",
			repl:   "1. Realizar consultas exploratórias",
		},
		{
			target: "• Implementar um protótipo conversacional",
			repl:   "2. Implementar um protótipo conversacional",
		},
		{
			target: "• Validar tecnicamente a comunicação",
			repl:   "3. Validar tecnicamente a comunicação",
		},
		{
			target: "• Modelar a orquestração das fontes",
			repl:   "4. Modelar a orquestração das fontes",
		},
		{
			target: "• Desenvolver um módulo de geração documental",
			repl:   "5. Desenvolver um módulo de geração documental",
		},
		{
			target: "• Construir um conjunto de dados de avaliação",
			repl:   "6. Construir um conjunto de dados de avaliação",
		},
		{
			target: "• Mensurar e comparar a Taxa de Conformidade",
			repl:   "7. Mensurar e comparar a Taxa de Conformidade",
		},
	}

	var reqs []*docs.Request
	for _, r := range replacements {
		reqs = append(reqs, &docs.Request{
			ReplaceAllText: &docs.ReplaceAllTextRequest{
				ContainsText: &docs.SubstringMatchCriteria{
					MatchCase: true,
					Text:      r.target,
				},
				ReplaceText: r.repl,
			},
		})
	}

	_, err = svc.Documents.BatchUpdate(docID, &docs.BatchUpdateDocumentRequest{Requests: reqs}).Do()
	if err != nil {
		log.Fatalf("Erro ao aplicar correções textuais: %v", err)
	}
	fmt.Println("Substituições de divisão e alíneas aplicadas com sucesso.")

	// 2. Formata com regras ABNT
	fmt.Println("Reformatando com regras ABNT...")
	_, err = gdocs.FormatDocumentABNTExtended(ctx, docID, -1, -1, true, true, gdocs.ABNTDefaultFontFamily, true, true, false)
	if err != nil {
		log.Fatalf("Erro format: %v", err)
	}

	// 3. Executa Lint
	opts := gdocs.LintOptions{
		AutoFix:              false,
		ClearHighlights:      false,
		ExpectedFontFamily:   gdocs.ABNTDefaultFontFamily,
		CheckStructuralOrder: true,
		CheckLatexLeaks:      true,
		CheckSpacingAndText:  true,
	}
	report, err := gdocs.LintDocumentABNT(ctx, docID, false, false, opts)
	if err != nil {
		log.Fatalf("Erro lint: %v", err)
	}

	repJSON, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println("\n=== RELATÓRIO DO LINT FINAL ===")
	fmt.Println(string(repJSON))
}
