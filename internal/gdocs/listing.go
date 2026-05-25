package gdocs

import (
	"context"
	"fmt"
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

	// Busca Google Docs e PDFs
	query := "(mimeType='application/vnd.google-apps.document' or mimeType='application/pdf') and trashed=false"

	type itemInfo struct {
		Title    string
		ID       string
		Link     string
		MimeType string
	}

	var allItems []itemInfo
	pageToken := ""

	for {
		call := svc.Files.List().
			Q(query).
			Fields(googleapi.Field("nextPageToken,files(id,name,webViewLink,mimeType)")).
			PageSize(100).
			OrderBy("name")

		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		result, err := call.Do()
		if err != nil {
			return "", fmt.Errorf("erro ao listar itens no Drive: %w", err)
		}

		for _, f := range result.Files {
			link := f.WebViewLink
			if link == "" && f.MimeType == "application/vnd.google-apps.document" {
				link = fmt.Sprintf("https://docs.google.com/document/d/%s/edit", f.Id)
			}
			allItems = append(allItems, itemInfo{
				Title:    f.Name,
				ID:       f.Id,
				Link:     link,
				MimeType: f.MimeType,
			})
		}

		if result.NextPageToken == "" {
			break
		}
		pageToken = result.NextPageToken
	}

	// Filtra pela allowlist se configurada (apenas para documentos, não afeta PDFs)
	allowedIDs := loadAllowedDocIDs()
	var docs []itemInfo
	var pdfs []itemInfo

	for _, item := range allItems {
		if item.MimeType == "application/pdf" {
			pdfs = append(pdfs, item)
		} else {
			if len(allowedIDs) > 0 && !isDocAllowed(item.ID) {
				continue
			}
			docs = append(docs, item)
		}
	}

	var sb strings.Builder

	sb.WriteString("📄 Documentos Editáveis (Google Docs):\n")
	if len(docs) == 0 {
		sb.WriteString("Nenhum documento encontrado (Verifique se compartilhou o Doc com a Service Account).\n")
	} else {
		for _, d := range docs {
			sb.WriteString(fmt.Sprintf("- ID: %s | Nome: %s\n", d.ID, d.Title))
		}
	}

	sb.WriteString("\n📕 PDFs Compartilhados (Para Sincronização via Update):\n")
	if len(pdfs) == 0 {
		sb.WriteString("(Atenção IA: Se esta lista de PDFs estiver vazia, o ambiente NÃO está configurado para upload).\n")
	} else {
		for _, f := range pdfs {
			sb.WriteString(fmt.Sprintf("- ID: %s | Nome: %s\n", f.ID, f.Title))
		}
	}

	if len(allowedIDs) > 0 {
		sb.WriteString(fmt.Sprintf("\nℹ️  Filtro ativo: ALLOWED_DOC_IDS limitou os resultados de documentos.\n"))
	}

	return sb.String(), nil
}
