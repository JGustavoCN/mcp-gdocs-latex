package gdocs

import (
	"context"
	"fmt"
	"html"
	"regexp"
	"strings"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// ListDocComments lista todos os comentários de um arquivo no Google Drive,
// incluindo o trecho de texto original ancorado, respostas e status de resolução.
func ListDocComments(ctx context.Context, fileID string) (string, error) {
	svc, err := GetDriveService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Drive: %w", err)
	}

	// Campos que queremos na resposta (obrigatório especificar no Drive API)
	fields := "nextPageToken,comments(id,content,author/displayName," +
		"quotedFileContent/value,resolved,createdTime,anchor," +
		"replies(content,author/displayName,createdTime))"

	// Coleta todos os comentários com paginação
	type comment = struct {
		ID                string
		Content           string
		AuthorDisplayName string
		QuotedValue       string
		Anchor            string
		Resolved          bool
		CreatedTime       string
		Replies           []struct {
			AuthorDisplayName string
			Content           string
			CreatedTime       string
		}
	}

	var allComments []comment
	pageToken := ""

	for {
		call := svc.Comments.List(fileID).
			Fields(googleapi.Field(fields)).
			IncludeDeleted(false).
			PageSize(100)

		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		result, err := call.Do()
		if err != nil {
			return "", fmt.Errorf("erro ao listar comentários do arquivo '%s': %w", fileID, err)
		}

		for _, c := range result.Comments {
			cm := comment{
				ID:          c.Id,
				Content:     c.Content,
				Resolved:    c.Resolved,
				CreatedTime: c.CreatedTime,
				Anchor:      c.Anchor,
			}

			if c.Author != nil {
				cm.AuthorDisplayName = c.Author.DisplayName
			} else {
				cm.AuthorDisplayName = "Desconhecido"
			}

			if c.QuotedFileContent != nil {
				cm.QuotedValue = c.QuotedFileContent.Value
			}

			for _, r := range c.Replies {
				rAuthor := "?"
				if r.Author != nil {
					rAuthor = r.Author.DisplayName
				}
				cm.Replies = append(cm.Replies, struct {
					AuthorDisplayName string
					Content           string
					CreatedTime       string
				}{
					AuthorDisplayName: rAuthor,
					Content:           r.Content,
					CreatedTime:       r.CreatedTime,
				})
			}

			allComments = append(allComments, cm)
		}

		if result.NextPageToken == "" {
			break
		}
		pageToken = result.NextPageToken
	}

	if len(allComments) == 0 {
		return "INFO: Nenhum comentário encontrado neste documento.", nil
	}

	// Formata a saída de forma legível e rica
	var sb strings.Builder
	sb.WriteString("=== COMENTÁRIOS DO DOCUMENTO ===\n")
	sb.WriteString(fmt.Sprintf("Total: %d comentário(s)\n\n", len(allComments)))

	for i, cm := range allComments {
		status := "⏳ Pendente"
		if cm.Resolved {
			status = "✅ Resolvido"
		}

		sb.WriteString(fmt.Sprintf("--- Comentário #%d (ID: %s) ---\n", i+1, cm.ID))
		sb.WriteString(fmt.Sprintf("  Autor     : %s\n", cm.AuthorDisplayName))
		sb.WriteString(fmt.Sprintf("  Data      : %s\n", cm.CreatedTime))
		sb.WriteString(fmt.Sprintf("  Status    : %s\n", status))

		if cm.QuotedValue != "" {
			sb.WriteString(fmt.Sprintf("  Trecho    : \"%s\"\n", cm.QuotedValue))
		}

		sb.WriteString(fmt.Sprintf("  Comentário: %s\n", cm.Content))

		if len(cm.Replies) > 0 {
			sb.WriteString(fmt.Sprintf("  Respostas (%d):\n", len(cm.Replies)))
			for _, r := range cm.Replies {
				sb.WriteString(fmt.Sprintf("    → [%s em %s]: %s\n",
					r.AuthorDisplayName, r.CreatedTime, r.Content))
			}
		}

		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// ReplyToComment adiciona uma resposta a um comentário existente.
func ReplyToComment(ctx context.Context, fileID, commentID, replyText string) (string, error) {
	svc, err := GetDriveService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Drive: %w", err)
	}

	reply := &drive.Reply{
		Content: replyText,
	}

	_, err = svc.Replies.Create(fileID, commentID, reply).Fields("id").Do()
	if err != nil {
		return "", fmt.Errorf("erro ao criar resposta no comentário '%s': %w", commentID, err)
	}

	return fmt.Sprintf("SUCESSO: Resposta adicionada ao comentário %s.", commentID), nil
}

// ResolveComment marca um comentário como resolvido.
func ResolveComment(ctx context.Context, fileID, commentID string) (string, error) {
	svc, err := GetDriveService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Drive: %w", err)
	}

	reply := &drive.Reply{
		Action:  "resolve",
		Content: "Resolvido e aplicado via MCP.",
	}

	// Executa a criação da resposta com a ação de resolução
	_, err = svc.Replies.Create(fileID, commentID, reply).Fields("id").Do()
	if err != nil {
		return "", fmt.Errorf("erro ao resolver o comentário '%s': %w", commentID, err)
	}

	return fmt.Sprintf("SUCESSO: Comentário %s marcado como resolvido.", commentID), nil
}

// CheckCommentCollision valida se o texto esperado cruza com algum comentário não resolvido.
func CheckCommentCollision(ctx context.Context, fileID string, expectedText string) error {
	svc, err := GetDriveService(ctx)
	if err != nil {
		return fmt.Errorf("erro ao obter serviço do Drive para validação de comentários: %w", err)
	}

	pageToken := ""
	for {
		call := svc.Comments.List(fileID).
			Fields("nextPageToken,comments(id,quotedFileContent/value,resolved,deleted)").
			IncludeDeleted(false).
			PageSize(100)

		if pageToken != "" {
			call = call.PageToken(pageToken)
		}

		result, err := call.Do()
		if err != nil {
			return fmt.Errorf("erro ao listar comentários para validação: %w", err)
		}

		for _, c := range result.Comments {
			if c.Deleted || c.Resolved {
				continue
			}

			if c.QuotedFileContent != nil {
				quotedValue := c.QuotedFileContent.Value
				if quotedValue == "" {
					continue
				}

				quotedValue = html.UnescapeString(quotedValue)
				quotedValueNorm := strings.ToLower(strings.TrimSpace(quotedValue))
				expectedTextNorm := strings.ToLower(strings.TrimSpace(expectedText))

				collision := false
				if len(quotedValueNorm) < 4 {
					// Word boundary match para palavras pequenas
					escaped := regexp.QuoteMeta(quotedValueNorm)
					pattern := `(?i)\b` + escaped + `\b`
					matched, _ := regexp.MatchString(pattern, expectedTextNorm)
					if matched {
						collision = true
					}
				} else {
					// Contains normal
					if strings.Contains(expectedTextNorm, quotedValueNorm) {
						collision = true
					}
				}

				if collision {
					return fmt.Errorf("Erro 403 BLOCKED: O trecho a ser substituído contém o comentário não resolvido ('%s'). Instrução OBRIGATÓRIA: Use a tool resolve_comment ou reply_to_comment com o ID %s para liberar a área antes de tentar sobrescrever o texto novamente.", quotedValue, c.Id)
				}
			}
		}

		if result.NextPageToken == "" {
			break
		}
		pageToken = result.NextPageToken
	}

	return nil
}
