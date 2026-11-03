package latex

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseTCCMetadata_RealFile(t *testing.T) {
	metaPath := filepath.Join("..", "..", "resources", "tcc", "metadados.tex")
	meta, err := ParseTCCMetadata(metaPath)
	if err != nil {
		t.Fatalf("Erro inesperado ao fazer parse de %s: %v", metaPath, err)
	}

	if !strings.Contains(meta.Titulo, "Avaliação de uma Arquitetura Híbrida RAG-MCP") {
		t.Errorf("Título incorreto, obteve: %q", meta.Titulo)
	}

	if meta.Autor != "José Gustavo Correia Nascimento" {
		t.Errorf("Autor incorreto, obteve: %q", meta.Autor)
	}

	if !strings.Contains(meta.Instituicao, "Instituto Federal de Educação, Ciência e Tecnologia de Sergipe") {
		t.Errorf("Instituição incorreta, obteve: %q", meta.Instituicao)
	}

	if !strings.Contains(meta.Orientador, "Diego Armando") {
		t.Errorf("Orientador incorreto, obteve: %q", meta.Orientador)
	}

	if meta.Coorientador != "" {
		t.Errorf("Coorientador deveria ser vazio pois está comentado no metadados.tex, obteve: %q", meta.Coorientador)
	}

	if !strings.Contains(meta.Preambulo, "Campus Lagarto") || !strings.Contains(meta.Preambulo, "obtenção do grau") {
		t.Errorf("Preâmbulo truncado ou incompleto, obteve: %q", meta.Preambulo)
	}

	lines := meta.InstituicaoLines()
	if len(lines) < 3 {
		t.Errorf("Esperava pelo menos 3 linhas na instituição, obteve %d: %v", len(lines), lines)
	}
}
