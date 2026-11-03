package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
	"github.com/JGustavoCN/mcp-gdocs-latex/internal/latex"
)

const (
	testDocID   = "1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8" // TCC Teste IA
	targetPdfID = "1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH"            // main.pdf no Google Drive
	localPdf    = "resources/tcc/main.pdf"
)

func main() {
	fmt.Println("====================================================================")
	fmt.Println("🚀 BATERIA DE TESTES LIVE E2E: MCP GOOGLE DOCS & DRIVE")
	fmt.Println("====================================================================")

	ctx := context.Background()

	// Garante credenciais apontando para o arquivo correto
	credPath, _ := filepath.Abs("credentials.json")
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credPath)
	fmt.Printf("[CONFIG] Credenciais ativas: %s\n\n", credPath)

	// ─── ETAPA 1: Paridade Criptográfica do PDF (CheckPDFSync) ──────────────
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println("📌 ETAPA 1: Verificação Criptográfica do PDF (CheckPDFSync)")
	fmt.Println("--------------------------------------------------------------------")
	report1, err := gdocs.CheckPDFSync(ctx, localPdf, targetPdfID)
	if err != nil {
		fmt.Printf("❌ Falha no CheckPDFSync inicial: %v\n", err)
	} else {
		repJSON, _ := json.MarshalIndent(report1, "", "  ")
		fmt.Printf("%s\n", string(repJSON))
	}

	// ─── ETAPA 2: Sincronização Perene do PDF (SyncPDFToDrive) ──────────────
	fmt.Println("\n--------------------------------------------------------------------")
	fmt.Println("📌 ETAPA 2: Upload Idempotente do PDF para o Drive (SyncPDFToDrive)")
	fmt.Println("--------------------------------------------------------------------")
	syncResult, err := gdocs.SyncPDFToDrive(ctx, targetPdfID, localPdf)
	if err != nil {
		fmt.Printf("❌ Falha no SyncPDFToDrive: %v\n", err)
	} else {
		fmt.Printf("✅ %s\n", syncResult)
	}

	// Pós-verificação para comprovar paridade IN_SYNC
	fmt.Println("\n🔍 Rechecagem de paridade pós-sincronização:")
	report2, err := gdocs.CheckPDFSync(ctx, localPdf, targetPdfID)
	if err != nil {
		fmt.Printf("❌ Falha no CheckPDFSync pós-upload: %v\n", err)
	} else {
		repJSON, _ := json.MarshalIndent(report2, "", "  ")
		fmt.Printf("%s\n", string(repJSON))
		if report2.Status == "IN_SYNC" {
			fmt.Println("🎯 SUCESSO: O PDF remoto e o PDF local possuem HASH MD5 100% IDÊNTICO!")
		}
	}

	// ─── ETAPA 3: Configuração de Margens ABNT NBR 14724 ────────────────────
	fmt.Println("\n--------------------------------------------------------------------")
	fmt.Println("📌 ETAPA 3: Aplicação de Margens Oficiais ABNT (ApplyABNTPageMargins)")
	fmt.Println("--------------------------------------------------------------------")
	marginResult, err := gdocs.ApplyABNTPageMargins(ctx, testDocID)
	if err != nil {
		fmt.Printf("❌ Falha ao aplicar margens: %v\n", err)
	} else {
		fmt.Printf("✅ %s\n", marginResult)
	}

	// ─── ETAPA 4: Inspeção Estrutural do Documento (GetDocSkeleton) ─────────
	fmt.Println("\n--------------------------------------------------------------------")
	fmt.Println("📌 ETAPA 4: Inspeção Estrutural e Localização de Índices (GetDocSkeleton)")
	fmt.Println("--------------------------------------------------------------------")
	skeleton, err := gdocs.GetDocSkeleton(ctx, testDocID)
	if err != nil {
		fmt.Printf("❌ Falha ao obter skeleton: %v\n", err)
	} else {
		fmt.Printf("✅ Esqueleto obtido com sucesso (%d bytes de JSON estruturado).\n", len(skeleton))
	}

	// ─── ETAPA 5: Inserção de Tabela Nativa ABNT/IBGE (InsertNativeTable) ────
	fmt.Println("\n--------------------------------------------------------------------")
	fmt.Println("📌 ETAPA 5: Inserção de Tabela Nativa ABNT/IBGE (InsertNativeTable)")
	fmt.Println("--------------------------------------------------------------------")

	// Obter o documento para capturar um índice de inserção seguro no final
	svc, err := gdocs.GetDocsService(ctx)
	if err != nil {
		fmt.Printf("❌ Falha ao obter docs service: %v\n", err)
		return
	}
	doc, err := svc.Documents.Get(testDocID).Do()
	if err != nil {
		fmt.Printf("❌ Falha ao inspecionar documento: %v\n", err)
		return
	}

	var lastContentEnd int64 = 1
	for _, el := range doc.Body.Content {
		if el.EndIndex > lastContentEnd {
			lastContentEnd = el.EndIndex
		}
	}
	// Posição de inserção segura logo antes do último caractere EOF
	insertTarget := lastContentEnd - 1
	if insertTarget < 1 {
		insertTarget = 1
	}
	fmt.Printf("📍 Índice de inserção calculado: %d (Fim do documento)\n", insertTarget)

	latexTable := `\begin{table}[htbp]
\centering
\caption{Matriz do Experimento Fatorial 2x2 de Injeção de Contexto}
\label{tab:fatorial_2x2}
\begin{tabular}{lll}
\toprule
\textbf{Arquitetura} & \textbf{Contexto Institucional (RAG)} & \textbf{Contexto Individual (MCP)} \\
\midrule
Baseline (Sem Contexto) & Ausente & Ausente \\
RAG Isolado & Presente (Regulamentos) & Ausente \\
MCP Isolado & Ausente & Presente (SUAP) \\
Híbrido (RAG + MCP) & Presente (Regulamentos) & Presente (SUAP) \\
\bottomrule
\end{tabular}
\source{Autoria própria (2026).}
\end{table}`

	parsedTable := latex.ParseSingleLaTeXTable(latexTable)
	fmt.Printf("📊 Dimensões da tabela: %d linhas x %d colunas\n", parsedTable.RowsCount, parsedTable.ColsCount)
	fmt.Printf("🏷️ Legenda: %s | Fonte: %s\n", parsedTable.Caption, parsedTable.Source)

	tableRes, err := gdocs.InsertNativeTable(ctx, testDocID, insertTarget, parsedTable)
	if err != nil {
		fmt.Printf("❌ Erro ao inserir tabela nativa: %v\n", err)
	} else {
		resJSON, _ := json.MarshalIndent(tableRes, "", "  ")
		fmt.Printf("✅ %s\n", string(resJSON))
	}

	// ─── ETAPA 6: Inserção de Diagrama Mermaid como Imagem (InsertMermaidDiagram) ───
	fmt.Println("\n--------------------------------------------------------------------")
	fmt.Println("📌 ETAPA 6: Inserção de Diagrama Mermaid / Imagem (InsertMermaidDiagram)")
	fmt.Println("--------------------------------------------------------------------")
	mmdPath := "resources/tcc/figuras/arquitetura.mmd"

	// Recalcula o fim do documento após a inserção da tabela
	docAfterTable, err := svc.Documents.Get(testDocID).Do()
	if err == nil {
		for _, el := range docAfterTable.Body.Content {
			if el.EndIndex > lastContentEnd {
				lastContentEnd = el.EndIndex
			}
		}
		insertTarget = lastContentEnd - 1
	}

	fmt.Printf("📍 Índice de inserção do diagrama: %d\n", insertTarget)
	fmt.Printf("🖼️ Fonte Mermaid: %s\n", mmdPath)

	mermaidRes, err := gdocs.InsertMermaidDiagram(
		ctx,
		testDocID,
		insertTarget,
		mmdPath,
		"Figura 1 – Arquitetura Híbrida RAG-MCP na Validação de Requerimentos",
		"Autoria própria (2026)",
	)
	if err != nil {
		fmt.Printf("❌ Erro ao inserir diagrama Mermaid: %v\n", err)
	} else {
		resJSON, _ := json.MarshalIndent(mermaidRes, "", "  ")
		fmt.Printf("✅ %s\n", string(resJSON))
	}

	// ─── ETAPA 7: Leitura de Verificação Final (ReadDocContent) ──────────────
	fmt.Println("\n--------------------------------------------------------------------")
	fmt.Println("📌 ETAPA 7: Leitura e Auditoria dos Elementos Inseridos (ReadDocContent)")
	fmt.Println("--------------------------------------------------------------------")
	// Lê os últimos 2000 caracteres do documento
	auditStart := insertTarget - 1500
	if auditStart < 1 {
		auditStart = 1
	}
	contentRead, err := gdocs.ReadDocContent(ctx, testDocID, auditStart, lastContentEnd+1000)
	if err != nil {
		fmt.Printf("❌ Falha na leitura de auditoria: %v\n", err)
	} else {
		lines := strings.Split(contentRead, "\n")
		// Mostra as primeiras e últimas linhas do trecho lido
		fmt.Println("📄 Trecho final auditado do Google Doc:")
		for i, l := range lines {
			if i < 40 {
				fmt.Println("  ", l)
			}
		}
	}

	// ─── ETAPA 8: Diagnóstico e Higienização de Anomalias (LintDocumentABNT) ───
	fmt.Println("\n--------------------------------------------------------------------")
	fmt.Println("📌 ETAPA 8: Diagnóstico e Higienização Estrutural (LintDocumentABNT)")
	fmt.Println("--------------------------------------------------------------------")
	lintReport, err := gdocs.LintDocumentABNT(ctx, testDocID, true, false)
	if err != nil {
		fmt.Printf("❌ Erro no linter: %v\n", err)
	} else {
		lintJSON, _ := json.MarshalIndent(lintReport, "", "  ")
		fmt.Printf("✅ %s\n", string(lintJSON))
	}

	fmt.Println("\n====================================================================")
	fmt.Println("🎉 BATERIA DE TESTES LIVE E2E CONCLUÍDA COM SUCESSO!")
	fmt.Println("====================================================================")
}
