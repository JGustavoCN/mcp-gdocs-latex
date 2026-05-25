package gdocs

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// docURLRegex captura o ID de um documento a partir de uma URL do Google Docs.
// Exemplos que ele reconhece:
//   - https://docs.google.com/document/d/1A2B3C_xyz/edit
//   - https://docs.google.com/document/d/1A2B3C_xyz/edit?usp=sharing
//   - https://docs.google.com/document/d/1A2B3C_xyz
var docURLRegex = regexp.MustCompile(`/document/d/([a-zA-Z0-9_-]+)`)

// ExtractDocID extrai o ID de um documento Google Docs a partir de uma URL
// completa ou retorna a string original se já for um ID puro.
func ExtractDocID(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	// Tenta extrair da URL
	if matches := docURLRegex.FindStringSubmatch(input); len(matches) > 1 {
		log.Printf("[INFO] ID extraído da URL: %s", matches[1])
		return matches[1]
	}

	// Assume que já é um ID
	return input
}

// loadAllowedDocIDs lê a variável de ambiente ALLOWED_DOC_IDS e retorna
// a lista de IDs permitidos (extraindo de URLs se necessário).
func loadAllowedDocIDs() []string {
	raw := os.Getenv("ALLOWED_DOC_IDS")
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	var ids []string
	for _, p := range parts {
		id := ExtractDocID(strings.TrimSpace(p))
		if id != "" {
			ids = append(ids, id)
		}
	}

	log.Printf("[INFO] ALLOWED_DOC_IDS configurado: %d documento(s) na lista de permissão", len(ids))
	return ids
}

// isDocAllowed verifica se um document ID está na lista de documentos permitidos.
// Retorna true se nenhuma restrição estiver configurada (ALLOWED_DOC_IDS vazio).
func isDocAllowed(docID string) bool {
	allowed := loadAllowedDocIDs()
	if len(allowed) == 0 {
		return true // sem restrições
	}

	for _, id := range allowed {
		if id == docID {
			return true
		}
	}
	return false
}

// ExtractAndValidateDocID é a função principal de UX: recebe uma URL ou ID,
// extrai o document ID, verifica a allowlist e faz um ping na API do Google
// para garantir que o documento existe e está acessível.
//
// Fluxo:
//  1. Extrai o ID via regex (ou usa como está)
//  2. Verifica se está na allowlist (ALLOWED_DOC_IDS)
//  3. Faz um GET rápido na API do Docs para validar acesso
//  4. Retorna o ID limpo e validado
func ExtractAndValidateDocID(ctx context.Context, input string) (string, error) {
	docID := ExtractDocID(input)
	if docID == "" {
		return "", fmt.Errorf(
			"não foi possível extrair um ID de documento válido da entrada: '%s'.\n"+
				"Dica: Cole a URL completa do Google Docs ou apenas o ID do documento",
			input,
		)
	}

	// Etapa 2: Verificar allowlist
	if !isDocAllowed(docID) {
		return "", fmt.Errorf(
			"acesso negado: o documento '%s' não está na lista de documentos "+
				"permitidos (ALLOWED_DOC_IDS).\n"+
				"Solução: Adicione este ID à variável de ambiente ALLOWED_DOC_IDS "+
				"ou remova a restrição deixando-a vazia",
			docID,
		)
	}

	// Etapa 3: Validação ativa — ping na API do Google Docs
	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao inicializar serviço do Google Docs para validação: %w", err)
	}

	doc, err := svc.Documents.Get(docID).Do()
	if err != nil {
		return "", fmt.Errorf(
			"o documento '%s' não foi encontrado ou não está acessível.\n"+
				"Verifique se:\n"+
				"  1. O ID ou URL informado está correto\n"+
				"  2. O documento foi compartilhado com o e-mail da Service Account\n"+
				"     (o e-mail está no campo 'client_email' do credentials.json)\n"+
				"  3. A Service Account tem pelo menos permissão de Leitor\n"+
				"Erro da API: %v",
			docID, err,
		)
	}

	log.Printf("[INFO] ✅ Documento validado: '%s' (ID: %s)", doc.Title, docID)
	return docID, nil
}
