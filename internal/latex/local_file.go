// Package latex fornece funções para manipulação segura
// de arquivos .tex no disco local.
package latex

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// UpdateLocalLatex reescreve o conteúdo de um arquivo .tex local.
// Antes de sobrescrever, cria automaticamente um backup (.bak) do original.
// Por segurança, aceita APENAS arquivos com extensão .tex.
func UpdateLocalLatex(filePath, newContent string) (string, error) {
	// Validação de segurança: apenas arquivos .tex
	if !strings.HasSuffix(strings.ToLower(filePath), ".tex") {
		return "", fmt.Errorf(
			"por segurança, esta ferramenta só permite reescrever arquivos " +
				"com extensão .tex",
		)
	}

	// Normaliza separadores de caminho para o SO
	filePath = filepath.Clean(filePath)

	// Verifica se o diretório pai existe
	parentDir := filepath.Dir(filePath)
	if parentDir != "" {
		info, err := os.Stat(parentDir)
		if err != nil || !info.IsDir() {
			return "", fmt.Errorf(
				"o diretório '%s' não existe. Verifique o caminho informado",
				parentDir,
			)
		}
	}

	// Se o arquivo já existe, cria um backup antes de sobrescrever
	if _, err := os.Stat(filePath); err == nil {
		backupPath := filePath + ".bak"

		originalContent, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("erro ao ler arquivo original para backup: %w", err)
		}

		if err := os.WriteFile(backupPath, originalContent, 0644); err != nil {
			return "", fmt.Errorf("erro ao criar backup em '%s': %w", backupPath, err)
		}

		log.Printf("[INFO] Backup criado em: %s", backupPath)
	}

	// Escreve o novo conteúdo em UTF-8
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return "", fmt.Errorf("erro ao escrever arquivo '%s': %w", filePath, err)
	}

	sizeKB := float64(len([]byte(newContent))) / 1024.0
	lineCount := strings.Count(newContent, "\n") + 1

	return fmt.Sprintf(
		"SUCESSO: Arquivo '%s' atualizado com sucesso.\nTamanho: %.1f KB | Linhas: %d",
		filePath, sizeKB, lineCount,
	), nil
}
