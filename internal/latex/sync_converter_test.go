package latex

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanTCCStructure(t *testing.T) {
	tccDir := filepath.Join("..", "..", "resources", "tcc")
	if _, err := os.Stat(tccDir); os.IsNotExist(err) {
		t.Skip("Diretório resources/tcc não encontrado, pulando teste")
	}

	elements, err := ScanTCCStructure(tccDir)
	if err != nil {
		t.Fatalf("Erro ao varrer estrutura do TCC: %v", err)
	}

	if len(elements) == 0 {
		t.Fatalf("Esperava encontrar elementos estruturais no TCC, obteve 0")
	}

	chapters := 0
	sections := 0
	tables := 0
	figures := 0

	for _, el := range elements {
		switch el.Type {
		case ElementChapter:
			chapters++
		case ElementSection:
			sections++
		case ElementTable:
			tables++
		case ElementFigure:
			figures++
		}
	}

	t.Logf("TCC Structure parsed: %d elementos totais | %d capítulos | %d seções | %d tabelas | %d figuras",
		len(elements), chapters, sections, tables, figures)

	if chapters < 4 {
		t.Errorf("Esperava pelo menos 4 capítulos, obteve %d", chapters)
	}
	if tables < 15 {
		t.Errorf("Esperava pelo menos 15 tabelas mapeadas, obteve %d", tables)
	}
	if figures < 3 {
		t.Errorf("Esperava pelo menos 3 figuras mapeadas, obteve %d", figures)
	}
}

func TestConvertGDocsTextToLatex(t *testing.T) {
	input := `Este é um teste com **negrito** e *itálico*. Citação [@feitoza2021] e "aspas normais". Taxa de 50% de sucesso.`
	expectedSubstrings := []string{
		`\textbf{negrito}`,
		`\textit{itálico}`,
		`\cite{feitoza2021}`,
		"``aspas normais''",
		`50\%`,
	}

	output := ConvertGDocsTextToLatex(input)

	for _, sub := range expectedSubstrings {
		if !strings.Contains(output, sub) {
			t.Errorf("Esperava que o output contivesse %q, mas obteve:\n%s", sub, output)
		}
	}
}

func TestUpdateLatexSectionInFile(t *testing.T) {
	tmpDir := t.TempDir()
	texFile := filepath.Join(tmpDir, "capitulo_teste.tex")

	initialContent := `\chapter{Introdução}
\label{cap:intro}

\section{Contexto}
\label{sec:contexto}
Texto original do contexto que precisa ser alterado.

\section{Justificativa}
\label{sec:justificativa}
Texto da justificativa intacto.
`
	if err := os.WriteFile(texFile, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Falha ao criar arquivo temporário: %v", err)
	}

	newText := "Novo texto atualizado a partir do Google Docs com **destaque** e 100% de precisão."
	msg, err := UpdateLatexSectionInFile(texFile, "Contexto", newText)
	if err != nil {
		t.Fatalf("Erro inesperado em UpdateLatexSectionInFile: %v", err)
	}

	t.Logf("Resultado: %s", msg)

	updatedBytes, err := os.ReadFile(texFile)
	if err != nil {
		t.Fatalf("Erro ao ler arquivo atualizado: %v", err)
	}
	updatedStr := string(updatedBytes)

	if !strings.Contains(updatedStr, `\textbf{destaque}`) {
		t.Errorf("Esperava encontrar \\textbf{destaque} no arquivo atualizado, obteve:\n%s", updatedStr)
	}
	if !strings.Contains(updatedStr, `100\%`) {
		t.Errorf("Esperava encontrar 100\\%% no arquivo atualizado, obteve:\n%s", updatedStr)
	}
	if !strings.Contains(updatedStr, `Texto da justificativa intacto.`) {
		t.Errorf("Esperava que a seção posterior (Justificativa) permanecesse intacta, obteve:\n%s", updatedStr)
	}

	// Verifica se o arquivo de backup foi criado
	backupFile := texFile + ".bak"
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		t.Errorf("Arquivo de backup .bak não foi criado em %s", backupFile)
	}
}

func TestGetNumberedTCCStructure(t *testing.T) {
	tccDir := filepath.Join("..", "..", "resources", "tcc")
	known, err := GetNumberedTCCStructure(tccDir)
	if err != nil {
		t.Fatalf("Erro: %v", err)
	}
	t.Logf("Total de seções numeradas: %d", len(known))
	for _, k := range known {
		t.Logf("Level %d | %s | %s", k.Level, k.TargetHeading, k.FullABNTTitle)
	}
}
