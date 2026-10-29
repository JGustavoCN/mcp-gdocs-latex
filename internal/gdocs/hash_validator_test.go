package gdocs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCalculateLocalPDFMD5_RealTCC(t *testing.T) {
	// Caminho relativo ao diretório raiz do projeto
	tccPDFPath := filepath.Join("..", "..", "resources", "tcc", "main.pdf")

	if _, err := os.Stat(tccPDFPath); os.IsNotExist(err) {
		t.Skip("resources/tcc/main.pdf não encontrado para teste de integração local")
	}

	md5Hash, size, err := CalculateLocalPDFMD5(tccPDFPath)
	if err != nil {
		t.Fatalf("Erro ao calcular hash do main.pdf: %v", err)
	}

	if size != 1098114 {
		t.Errorf("Tamanho do main.pdf esperado: 1098114 bytes; obteve %d", size)
	}

	if len(md5Hash) != 32 {
		t.Errorf("Hash MD5 inválido: %s (comprimento %d)", md5Hash, len(md5Hash))
	}
}

func TestCalculateLocalPDFMD5_CorruptedNoEOF(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_corrupt_*.pdf")
	if err != nil {
		t.Fatalf("Erro ao criar temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Escreve dados sem o marcador %%EOF
	tmpFile.WriteString("%PDF-1.4\nConteúdo sem marcador final...")
	tmpFile.Close()

	_, _, err = CalculateLocalPDFMD5(tmpFile.Name())
	if err == nil {
		t.Fatalf("Esperava erro de integridade ISO 32000-1 para PDF sem %%%%EOF")
	}

	if !strings.Contains(err.Error(), "%%EOF") {
		t.Errorf("Mensagem de erro esperada mencionando %%%%EOF, obteve: %v", err)
	}
}

