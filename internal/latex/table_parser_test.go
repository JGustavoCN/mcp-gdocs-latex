package latex

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAllTablesInLatex(t *testing.T) {
	sampleLatex := `
\begin{table}[ht]
\centering
\caption{Delimitação do escopo experimental em relação ao protótipo funcional}
\label{tab:escopo_recorte}
\begin{tabular}{lll}
\toprule
\textbf{Dimensão} & \textbf{Protótipo Funcional} & \textbf{Bancada Experimental} \\
\midrule
Entrada de dados & Chat livre & Cenários simulados \\
Validação & ROD completo & Trancamento e Proficiência \\
Saída & Folha A4 e JSON & Avaliação da TCD \\
\bottomrule
\end{tabular}
\fonte{Autor (2026).}
\end{table}
`

	tbls := ParseAllTablesInLatex(sampleLatex)
	if len(tbls) != 1 {
		t.Fatalf("Esperava 1 tabela, obteve %d", len(tbls))
	}

	tbl := tbls[0]
	if tbl.Label != "tab:escopo_recorte" {
		t.Errorf("Label incorreto: %s", tbl.Label)
	}

	if tbl.ColsCount != 3 {
		t.Errorf("ColsCount esperado: 3, obteve %d", tbl.ColsCount)
	}

	if tbl.RowsCount != 4 { // 1 cabeçalho + 3 dados
		t.Errorf("RowsCount esperado: 4, obteve %d", tbl.RowsCount)
	}

	if len(tbl.Headers) != 3 {
		t.Errorf("Headers esperado: 3, obteve %d", len(tbl.Headers))
	}

	if tbl.Headers[0] != "Dimensão" {
		t.Errorf("Header[0] esperado: 'Dimensão', obteve %q", tbl.Headers[0])
	}
}

func TestParseRealTCCCaptitulosTables(t *testing.T) {
	capDir := filepath.Join("..", "..", "resources", "tcc", "capitulos")
	if _, err := os.Stat(capDir); os.IsNotExist(err) {
		t.Skip("resources/tcc/capitulos não encontrado")
	}

	files, err := filepath.Glob(filepath.Join(capDir, "*.tex"))
	if err != nil {
		t.Fatalf("Erro ao listar arquivos .tex: %v", err)
	}

	totalTables := 0
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("Erro ao ler arquivo %s: %v", f, err)
		}
		tbls := ParseAllTablesInLatex(string(content))
		totalTables += len(tbls)
	}

	t.Logf("Total de tabelas e quadros extraídos dos capítulos do TCC: %d", totalTables)
	if totalTables == 0 {
		t.Errorf("Nenhuma tabela foi encontrada nos arquivos do TCC")
	}
}

