package gdocs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/api/googleapi"
)

// ListAvailableDocuments lista todos os documentos Google Docs acessíveis
// pela Service Account, retornando título, ID e link de cada um.
// Se ALLOWED_DOC_IDS estiver configurado, filtra apenas os permitidos.
func ListAvailableDocuments(ctx context.Context) (string, error) {
	svc, err := GetDriveService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Drive: %w", err)
	}

	// Busca apenas Google Docs (não planilhas, slides, etc.)
	query := "mimeType='application/vnd.google-apps.document' and trashed=false"

	type docInfo struct {
		Title string
		ID    string
		Link  string
	}

	var allDocs []docInfo
	pageToken := ""

	for {
		call := svc.Files.List().
			Q(query).
			Fields(googleapi.Field("nextPageToken,files(id,name,webViewLink)")).
			PageSize(100).
			OrderBy("name")

		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		result, err := call.Do()
		if err != nil {
			return "", fmt.Errorf("erro ao listar documentos no Drive: %w", err)
		}

		for _, f := range result.Files {
			link := f.WebViewLink
			if link == "" {
				link = fmt.Sprintf("https://docs.google.com/document/d/%s/edit", f.Id)
			}
			allDocs = append(allDocs, docInfo{
				Title: f.Name,
				ID:    f.Id,
				Link:  link,
			})
		}

		if result.NextPageToken == "" {
			break
		}
		pageToken = result.NextPageToken
	}

	// Filtra pela allowlist se configurada
	allowedIDs := loadAllowedDocIDs()
	if len(allowedIDs) > 0 {
		var filtered []docInfo
		for _, d := range allDocs {
			if isDocAllowed(d.ID) {
				filtered = append(filtered, d)
			}
		}
		log.Printf("[INFO] Filtro ALLOWED_DOC_IDS aplicado: %d de %d documentos permitidos",
			len(filtered), len(allDocs))
		allDocs = filtered
	}

	if len(allDocs) == 0 {
		return "INFO: Nenhum documento Google Docs acessível foi encontrado.\n\n" +
			"Possíveis causas:\n" +
			"  1. Nenhum documento foi compartilhado com o e-mail da Service Account\n" +
			"  2. A variável ALLOWED_DOC_IDS está restringindo os resultados\n\n" +
			"Solução: Compartilhe o Google Doc com o e-mail que aparece no campo\n" +
			"'client_email' do arquivo credentials.json.", nil
	}

	var sb strings.Builder
	sb.WriteString("=== DOCUMENTOS DISPONÍVEIS ===\n")
	sb.WriteString(fmt.Sprintf("Total: %d documento(s) acessível(is)\n\n", len(allDocs)))

	for i, d := range allDocs {
		sb.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, d.Title))
		sb.WriteString(fmt.Sprintf("   ID  : %s\n", d.ID))
		sb.WriteString(fmt.Sprintf("   Link: %s\n\n", d.Link))
	}

	if len(allowedIDs) > 0 {
		sb.WriteString(fmt.Sprintf("ℹ️  Filtro ativo: ALLOWED_DOC_IDS (%d ID(s) configurados)\n", len(allowedIDs)))
	}

	return sb.String(), nil
}
