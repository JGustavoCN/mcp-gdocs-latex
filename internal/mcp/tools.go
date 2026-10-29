// Package mcp registra e conecta as ferramentas do servidor MCP,
// servindo como ponte entre o protocolo MCP e os módulos internos.
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
	"github.com/JGustavoCN/mcp-gdocs-latex/internal/latex"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registra todas as ferramentas no servidor MCP.
func RegisterTools(s *server.MCPServer) {
	// Ferramentas pré-existentes
	registerListAvailableDocuments(s)
	registerGetDocSkeleton(s)
	registerReadDocContent(s)
	registerSearchInDoc(s)
	registerListDocComments(s)
	registerReplyToComment(s)
	registerResolveComment(s)
	registerApplyTextStyle(s)
	registerSearchAndReplaceText(s)
	registerUpdateBlockByIndex(s)
	registerSyncPDFToDrive(s)
	registerUpdateLocalLatex(s)

	// Novas ferramentas ABNT / Tabelas / Imagens / Hash
	registerApplyAbntPageMargins(s)
	registerReplaceBlockAbnt(s)
	registerFormatDocumentAbnt(s)
	registerInsertNativeTable(s)
	registerInsertMermaidDiagram(s)
	registerCheckPdfSync(s)
	registerLintDocumentAbnt(s)
}

// ─── Tool X: get_doc_skeleton ───────────────────────────────────────────────

func registerGetDocSkeleton(s *server.MCPServer) {
	tool := mcplib.NewTool("get_doc_skeleton",
		mcplib.WithDescription(
			"Retorna a estrutura hierárquica e lógica de tópicos, capítulos e parágrafos marcados como cabeçalhos ou destaques do documento, contendo seus índices absolutos de caracteres (start_index e end_index). Essencial para planejar leituras parciais de documentos acadêmicos extensos antes de ler o texto completo."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
	)
	s.AddTool(tool, handleGetDocSkeleton)
	log.Println("[INFO] Ferramenta registrada: get_doc_skeleton")
}

func handleGetDocSkeleton(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] get_doc_skeleton → documento: %s", docID)

	result, err := gdocs.GetDocSkeleton(ctx, docID)
	if err != nil {
		log.Printf("[ERRO] get_doc_skeleton: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 1: read_doc_content ───────────────────────────────────────────────

func registerReadDocContent(s *server.MCPServer) {
	tool := mcplib.NewTool("read_doc_content",
		mcplib.WithDescription(
			"Lê o texto do Google Doc e o retorna estruturado em formato Markdown. Aceita parâmetros opcionais de início (start_index) e fim (end_index) para ler fatias específicas do documento, reduzindo o consumo de tokens em arquivos grandes."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithInteger("start_index",
			mcplib.Description("Índice inicial da leitura (retornado por get_doc_skeleton). Se omitido, lê desde o começo."),
		),
		mcplib.WithInteger("end_index",
			mcplib.Description("Índice final da leitura (retornado por get_doc_skeleton). Se omitido, lê até o fim."),
		),
	)
	s.AddTool(tool, handleReadDocContent)
	log.Println("[INFO] Ferramenta registrada: read_doc_content")
}

func handleReadDocContent(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	startIdx := int64(request.GetInt("start_index", -1))
	endIdx := int64(request.GetInt("end_index", -1))

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		log.Printf("[ERRO] read_doc_content — validação falhou: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] read_doc_content → documento: %s (índices: %d a %d)", docID, startIdx, endIdx)

	result, err := gdocs.ReadDocContent(ctx, docID, startIdx, endIdx)
	if err != nil {
		log.Printf("[ERRO] read_doc_content: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 2: list_doc_comments ──────────────────────────────────────────────

func registerListDocComments(s *server.MCPServer) {
	tool := mcplib.NewTool("list_doc_comments",
		mcplib.WithDescription(
			"Lista todos os comentários de um arquivo no Google Drive,\n"+
				"incluindo o trecho de texto original ao qual cada comentário está ancorado.\n"+
				"Ideal para capturar correções e sugestões do orientador.\n\n"+
				"Aceita tanto o ID puro do documento quanto a URL completa do Google Docs.\n"+
				"O ID do Google Docs é o mesmo usado no Google Drive."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
	)
	s.AddTool(tool, handleListDocComments)
	log.Println("[INFO] Ferramenta registrada: list_doc_comments")
}

func handleListDocComments(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		log.Printf("[ERRO] list_doc_comments — validação falhou: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] list_doc_comments → documento: %s", docID)

	result, err := gdocs.ListDocComments(ctx, docID)
	if err != nil {
		log.Printf("[ERRO] list_doc_comments: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 3: search_in_doc ──────────────────────────────────────────────────

func registerSearchInDoc(s *server.MCPServer) {
	tool := mcplib.NewTool("search_in_doc",
		mcplib.WithDescription(
			"Pesquisa um termo no Google Doc e retorna os índices inicial e final de todas as ocorrências encontradas junto a um snippet de contexto. Se o parâmetro 'is_regex' for verdadeiro, processa a consulta como uma Expressão Regular para buscas complexas ou de padrões de formatação (ABNT)."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithString("query",
			mcplib.Required(),
			mcplib.Description("O texto a ser buscado no documento"),
		),
		mcplib.WithBoolean("is_regex",
			mcplib.Description("Se true, trata o query como uma Expressão Regular em vez de busca literal (default: false)"),
		),
	)
	s.AddTool(tool, handleSearchInDoc)
	log.Println("[INFO] Ferramenta registrada: search_in_doc")
}

func handleSearchInDoc(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	query, err := request.RequireString("query")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'query' é obrigatório."), nil
	}

	isRegex := request.GetBool("is_regex", false)

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] search_in_doc → doc: %s | query: '%s' | regex: %v", docID, query, isRegex)

	result, err := gdocs.SearchInDoc(ctx, docID, query, isRegex)
	if err != nil {
		log.Printf("[ERRO] search_in_doc: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 6: reply_to_comment ───────────────────────────────────────────────

func registerReplyToComment(s *server.MCPServer) {
	tool := mcplib.NewTool("reply_to_comment",
		mcplib.WithDescription(
			"Adiciona uma resposta a um comentário existente no Google Docs via Drive API."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithString("comment_id",
			mcplib.Required(),
			mcplib.Description("ID do comentário original"),
		),
		mcplib.WithString("reply_text",
			mcplib.Required(),
			mcplib.Description("Texto da resposta (ex: 'Corrigido na nova versão do TCC.')"),
		),
	)
	s.AddTool(tool, handleReplyToComment)
	log.Println("[INFO] Ferramenta registrada: reply_to_comment")
}

func handleReplyToComment(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	commentID, err := request.RequireString("comment_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'comment_id' é obrigatório."), nil
	}

	replyText, err := request.RequireString("reply_text")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'reply_text' é obrigatório."), nil
	}

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] reply_to_comment → doc: %s | comment_id: %s", docID, commentID)

	result, err := gdocs.ReplyToComment(ctx, docID, commentID, replyText)
	if err != nil {
		log.Printf("[ERRO] reply_to_comment: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 7: resolve_comment ────────────────────────────────────────────────

func registerResolveComment(s *server.MCPServer) {
	tool := mcplib.NewTool("resolve_comment",
		mcplib.WithDescription(
			"Marca um comentário como resolvido no Google Docs via Drive API."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithString("comment_id",
			mcplib.Required(),
			mcplib.Description("ID do comentário a ser resolvido"),
		),
	)
	s.AddTool(tool, handleResolveComment)
	log.Println("[INFO] Ferramenta registrada: resolve_comment")
}

func handleResolveComment(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	commentID, err := request.RequireString("comment_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'comment_id' é obrigatório."), nil
	}

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] resolve_comment → doc: %s | comment_id: %s", docID, commentID)

	result, err := gdocs.ResolveComment(ctx, docID, commentID)
	if err != nil {
		log.Printf("[ERRO] resolve_comment: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 8: search_and_replace_text ────────────────────────────────────────

func registerSearchAndReplaceText(s *server.MCPServer) {
	tool := mcplib.NewTool("search_and_replace_text",
		mcplib.WithDescription(
			"Usada apenas para correções ortográficas rápidas e troca global de termos. Procura e substitui ocorrências exatas de texto globalmente no Google Docs."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithString("old_text",
			mcplib.Required(),
			mcplib.Description("O texto exato a ser encontrado e substituído"),
		),
		mcplib.WithString("new_text",
			mcplib.Required(),
			mcplib.Description("O novo texto que substituirá o antigo"),
		),
	)
	s.AddTool(tool, handleSearchAndReplaceText)
	log.Println("[INFO] Ferramenta registrada: search_and_replace_text")
}

func handleSearchAndReplaceText(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	oldText, err := request.RequireString("old_text")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'old_text' é obrigatório."), nil
	}

	newText, err := request.RequireString("new_text")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'new_text' é obrigatório."), nil
	}

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] search_and_replace_text → doc: %s | substituir: '%s' → '%s'",
		docID, oldText, newText)

	result, err := gdocs.SearchAndReplaceText(ctx, docID, oldText, newText)
	if err != nil {
		log.Printf("[ERRO] search_and_replace_text: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 5: update_block_by_index ──────────────────────────────────────────

func registerUpdateBlockByIndex(s *server.MCPServer) {
	tool := mcplib.NewTool("update_block_by_index",
		mcplib.WithDescription(
			"Ferramenta determinística para aplicar formatação ABNT baseada em JSON. Substitui um bloco exato por uma lista de trechos de texto estruturados, com suporte a estilos ricos (negrito, itálico) e parágrafos (HEADING_1, NORMAL_TEXT)."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithString("payload_json",
			mcplib.Required(),
			mcplib.Description("String JSON contendo {start_index, end_index, expected_text, content: [{text, paragraph_style, text_style}]}"),
		),
	)
	s.AddTool(tool, handleUpdateBlockByIndex)
	log.Println("[INFO] Ferramenta registrada: update_block_by_index")
}

func handleUpdateBlockByIndex(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	payloadJSON, err := request.RequireString("payload_json")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'payload_json' é obrigatório."), nil
	}

	var req gdocs.ReplaceRequest
	if err := json.Unmarshal([]byte(payloadJSON), &req); err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: Falha ao fazer parse do payload_json: %v", err)), nil
	}

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] update_block_by_index → doc: %s", docID)

	result, err := gdocs.UpdateBlockByIndex(ctx, docID, req)
	if err != nil {
		log.Printf("[ERRO] update_block_by_index: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 11: apply_text_style ───────────────────────────────────────────────

func registerApplyTextStyle(s *server.MCPServer) {
	tool := mcplib.NewTool("apply_text_style",
		mcplib.WithDescription(
			"Use esta ferramenta para aplicar estilos em palavras específicas após encontrar seus índices com search_in_doc. Não apaga texto e preserva comentários."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithInteger("start_index",
			mcplib.Required(),
			mcplib.Description("Índice de início"),
		),
		mcplib.WithInteger("end_index",
			mcplib.Required(),
			mcplib.Description("Índice de fim"),
		),
		mcplib.WithString("style",
			mcplib.Required(),
			mcplib.Description("Estilo a ser aplicado. Valores suportados: BOLD, ITALIC, NONE"),
		),
	)
	s.AddTool(tool, handleApplyTextStyle)
	log.Println("[INFO] Ferramenta registrada: apply_text_style")
}

func handleApplyTextStyle(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}
	
	startIndex := request.GetInt("start_index", -1)
	endIndex := request.GetInt("end_index", -1)
	if startIndex < 0 || endIndex < 0 || startIndex >= endIndex {
		return mcplib.NewToolResultError("Parâmetros 'start_index' e 'end_index' devem ser válidos e maiores que zero (com start_index < end_index)."), nil
	}

	style, err := request.RequireString("style")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'style' é obrigatório."), nil
	}

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] apply_text_style → doc: %s | [%d-%d] | style: %s", docID, startIndex, endIndex, style)

	result, err := gdocs.ApplyTextStyle(ctx, docID, int64(startIndex), int64(endIndex), style)
	if err != nil {
		log.Printf("[ERRO] apply_text_style: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 9: sync_pdf_to_drive ──────────────────────────────────────────────

func registerSyncPDFToDrive(s *server.MCPServer) {
	tool := mcplib.NewTool("sync_pdf_to_drive",
		mcplib.WithDescription(
			"Faz o upload do PDF compilado de forma idempotente, verificando o hash MD5 e integridade prévia %%EOF (ISO 32000-1). Se idêntico, o upload é suprimido mantendo a revisão perene."),
		mcplib.WithString("local_pdf_path",
			mcplib.Required(),
			mcplib.Description("Caminho absoluto ou relativo do arquivo PDF local (ex: 'resources/tcc/main.pdf')"),
		),
		mcplib.WithString("target_pdf_file_id",
			mcplib.Required(),
			mcplib.Description("Obrigatório. ID do arquivo PDF existente no drive para atualização contínua."),
		),
	)
	s.AddTool(tool, handleSyncPDFToDrive)
	log.Println("[INFO] Ferramenta registrada: sync_pdf_to_drive")
}

func handleSyncPDFToDrive(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	localPDFPath, err := request.RequireString("local_pdf_path")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'local_pdf_path' é obrigatório."), nil
	}

	targetPdfFileID, err := request.RequireString("target_pdf_file_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'target_pdf_file_id' é obrigatório."), nil
	}

	log.Printf("[INFO] sync_pdf_to_drive → pdf: '%s' | updateID: '%s'", localPDFPath, targetPdfFileID)

	resultStr, err := gdocs.SyncPDFToDrive(ctx, targetPdfFileID, localPDFPath)
	if err != nil {
		log.Printf("[ERRO] sync_pdf_to_drive: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(resultStr), nil
}

// ─── Tool 4: update_local_latex ─────────────────────────────────────────────

func registerUpdateLocalLatex(s *server.MCPServer) {
	tool := mcplib.NewTool("update_local_latex",
		mcplib.WithDescription(
			"Reescreve o conteúdo de um arquivo .tex local no disco.\n"+
				"Cria automaticamente um backup (.bak) do original antes de sobrescrever.\n"+
				"Por segurança, aceita APENAS arquivos com extensão .tex.\n\n"+
				"Parâmetros:\n"+
				"- filepath: Caminho absoluto do arquivo .tex (ex: 'C:/Users/.../main.tex')\n"+
				"- new_content: O conteúdo LaTeX completo que será escrito no arquivo."),
		mcplib.WithString("filepath",
			mcplib.Required(),
			mcplib.Description("Caminho absoluto do arquivo .tex"),
		),
		mcplib.WithString("new_content",
			mcplib.Required(),
			mcplib.Description("O conteúdo LaTeX completo que será escrito no arquivo"),
		),
	)
	s.AddTool(tool, handleUpdateLocalLatex)
	log.Println("[INFO] Ferramenta registrada: update_local_latex")
}

func handleUpdateLocalLatex(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	filePath, err := request.RequireString("filepath")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'filepath' é obrigatório."), nil
	}

	newContent, err := request.RequireString("new_content")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'new_content' é obrigatório."), nil
	}

	log.Printf("[INFO] update_local_latex → arquivo: %s", filePath)

	result, err := latex.UpdateLocalLatex(filePath, newContent)
	if err != nil {
		log.Printf("[ERRO] update_local_latex: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 5: list_available_documents ───────────────────────────────────────

func registerListAvailableDocuments(s *server.MCPServer) {
	tool := mcplib.NewTool("list_available_documents",
		mcplib.WithDescription(
			"Lista os documentos do Google Docs e os arquivos PDF compartilhados com o robô. Funciona como um Health Check. Se a seção de PDFs retornar vazia, o LLM DEVE parar o fluxo e instruir o usuário a criar um PDF no Drive e compartilhá-lo como Editor com o e-mail do robô."),
	)
	s.AddTool(tool, handleListAvailableDocuments)
	log.Println("[INFO] Ferramenta registrada: list_available_documents")
}

func handleListAvailableDocuments(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	log.Println("[INFO] list_available_documents → listando documentos acessíveis...")

	result, err := gdocs.ListAvailableDocuments(ctx)
	if err != nil {
		log.Printf("[ERRO] list_available_documents: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 12: apply_abnt_page_margins ───────────────────────────────────────

func registerApplyAbntPageMargins(s *server.MCPServer) {
	tool := mcplib.NewTool("apply_abnt_page_margins",
		mcplib.WithDescription(
			"Aplica as margens oficiais da ABNT NBR 14724 no Google Docs (Superior 3cm, Esquerda 3cm, Inferior 2cm, Direita 2cm) e formato de página A4."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
	)
	s.AddTool(tool, handleApplyAbntPageMargins)
	log.Println("[INFO] Ferramenta registrada: apply_abnt_page_margins")
}

func handleApplyAbntPageMargins(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] apply_abnt_page_margins → doc: %s", docID)

	result, err := gdocs.ApplyABNTPageMargins(ctx, docID)
	if err != nil {
		log.Printf("[ERRO] apply_abnt_page_margins: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 13: replace_block_abnt ────────────────────────────────────────────

func registerReplaceBlockAbnt(s *server.MCPServer) {
	tool := mcplib.NewTool("replace_block_abnt",
		mcplib.WithDescription(
			"Substitui um bloco de texto aplicando a sequência canônica anti-bleeding (elimina bold bleeding e cascata de marcadores) e formatação ABNT (recuo 1,25cm, 1.5 entrelinhas, justificado)."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithInteger("start_index",
			mcplib.Required(),
			mcplib.Description("Índice de início"),
		),
		mcplib.WithInteger("end_index",
			mcplib.Required(),
			mcplib.Description("Índice de fim"),
		),
		mcplib.WithString("new_text",
			mcplib.Required(),
			mcplib.Description("Novo texto a ser inserido e higienizado com regras ABNT"),
		),
	)
	s.AddTool(tool, handleReplaceBlockAbnt)
	log.Println("[INFO] Ferramenta registrada: replace_block_abnt")
}

func handleReplaceBlockAbnt(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	startIdx := int64(request.GetInt("start_index", -1))
	endIdx := int64(request.GetInt("end_index", -1))
	if startIdx < 0 || endIdx < startIdx {
		return mcplib.NewToolResultError("Parâmetros 'start_index' e 'end_index' devem ser válidos (start_index <= end_index)."), nil
	}

	newText, err := request.RequireString("new_text")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'new_text' é obrigatório."), nil
	}

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] replace_block_abnt → doc: %s | [%d-%d]", docID, startIdx, endIdx)

	result, err := gdocs.ReplaceBlockABNT(ctx, docID, startIdx, endIdx, newText)
	if err != nil {
		log.Printf("[ERRO] replace_block_abnt: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 14: format_document_abnt ──────────────────────────────────────────

func registerFormatDocumentAbnt(s *server.MCPServer) {
	tool := mcplib.NewTool("format_document_abnt",
		mcplib.WithDescription(
			"Formata parágrafos do documento no padrão ABNT NBR 14724 (alinhamento justificado, espaçamento 1.5, recuo de 1ª linha de 1.25cm). Aceita opcionalmente start_index e end_index para formatação parcial."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithInteger("start_index",
			mcplib.Description("Índice de início (opcional, padrão: início do documento)"),
		),
		mcplib.WithInteger("end_index",
			mcplib.Description("Índice de fim (opcional, padrão: fim do documento)"),
		),
		mcplib.WithBoolean("apply_margins",
			mcplib.Description("Se true, também configura margens ABNT 3/3/2/2 cm na folha A4 (padrão: true)"),
		),
	)
	s.AddTool(tool, handleFormatDocumentAbnt)
	log.Println("[INFO] Ferramenta registrada: format_document_abnt")
}

func handleFormatDocumentAbnt(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	startIdx := int64(request.GetInt("start_index", -1))
	endIdx := int64(request.GetInt("end_index", -1))
	applyMargins := request.GetBool("apply_margins", true)

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] format_document_abnt → doc: %s | [%d-%d] | margins: %v", docID, startIdx, endIdx, applyMargins)

	result, err := gdocs.FormatDocumentABNT(ctx, docID, startIdx, endIdx, applyMargins)
	if err != nil {
		log.Printf("[ERRO] format_document_abnt: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 15: insert_native_table ───────────────────────────────────────────

func registerInsertNativeTable(s *server.MCPServer) {
	tool := mcplib.NewTool("insert_native_table",
		mcplib.WithDescription(
			"Insere uma tabela nativa estruturada no Google Docs segundo o padrão ABNT/IBGE (sem bordas verticais, cabeçalho em negrito, legenda no topo e fonte no rodapé). Aceita payload JSON ou código de tabela LaTeX."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithInteger("target_index",
			mcplib.Required(),
			mcplib.Description("Índice de inserção da tabela no documento"),
		),
		mcplib.WithString("latex_tabular_or_json",
			mcplib.Required(),
			mcplib.Description("Bloco LaTeX da tabela (\\begin{table}... ou \\begin{tabular}...) ou JSON estruturado com {caption, source, headers, data}"),
		),
	)
	s.AddTool(tool, handleInsertNativeTable)
	log.Println("[INFO] Ferramenta registrada: insert_native_table")
}

func handleInsertNativeTable(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	targetIndex := int64(request.GetInt("target_index", -1))
	if targetIndex < 0 {
		return mcplib.NewToolResultError("Parâmetro 'target_index' deve ser maior ou igual a zero."), nil
	}

	rawInput, err := request.RequireString("latex_tabular_or_json")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'latex_tabular_or_json' é obrigatório."), nil
	}

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	var tableData latex.LaTeXTableData
	trimmed := strings.TrimSpace(rawInput)

	if strings.Contains(trimmed, "\\begin{") || strings.Contains(trimmed, "&") {
		tableData = latex.ParseSingleLaTeXTable(trimmed)
	} else {
		if errJSON := json.Unmarshal([]byte(trimmed), &tableData); errJSON != nil {
			return mcplib.NewToolResultError(fmt.Sprintf("ERRO: Entrada não é LaTeX válido nem JSON estruturado: %v", errJSON)), nil
		}
	}

	if tableData.RowsCount <= 0 {
		tableData.RowsCount = int64(len(tableData.Data))
		if len(tableData.Headers) > 0 {
			tableData.RowsCount++
		}
	}
	if tableData.ColsCount <= 0 {
		if len(tableData.Headers) > 0 {
			tableData.ColsCount = int64(len(tableData.Headers))
		} else if len(tableData.Data) > 0 {
			tableData.ColsCount = int64(len(tableData.Data[0]))
		}
	}

	log.Printf("[INFO] insert_native_table → doc: %s | target: %d | matriz: %dx%d", docID, targetIndex, tableData.RowsCount, tableData.ColsCount)

	result, err := gdocs.InsertNativeTable(ctx, docID, targetIndex, tableData)
	if err != nil {
		log.Printf("[ERRO] insert_native_table: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	resJSON, _ := json.MarshalIndent(result, "", "  ")
	return mcplib.NewToolResultText(string(resJSON)), nil
}

// ─── Tool 16: insert_mermaid_diagram ────────────────────────────────────────

func registerInsertMermaidDiagram(s *server.MCPServer) {
	tool := mcplib.NewTool("insert_mermaid_diagram",
		mcplib.WithDescription(
			"Insere um diagrama Mermaid como imagem no Google Docs com cálculo de 2D Bounding Box (proporcional até 450x620 pt), upload no Drive e fallback para Deep Link caso haja restrição DLP."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithInteger("target_index",
			mcplib.Required(),
			mcplib.Description("Índice de inserção da imagem no documento"),
		),
		mcplib.WithString("png_path",
			mcplib.Required(),
			mcplib.Description("Caminho absoluto ou relativo do arquivo PNG da imagem renderizada"),
		),
		mcplib.WithString("caption",
			mcplib.Description("Legenda da figura (opcional)"),
		),
		mcplib.WithString("source",
			mcplib.Description("Fonte da figura (ex: 'Autor (2026)', opcional)"),
		),
	)
	s.AddTool(tool, handleInsertMermaidDiagram)
	log.Println("[INFO] Ferramenta registrada: insert_mermaid_diagram")
}

func handleInsertMermaidDiagram(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	targetIndex := int64(request.GetInt("target_index", -1))
	if targetIndex < 0 {
		return mcplib.NewToolResultError("Parâmetro 'target_index' deve ser maior ou igual a zero."), nil
	}

	pngPath, err := request.RequireString("png_path")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'png_path' é obrigatório."), nil
	}

	caption := request.GetString("caption", "")
	source := request.GetString("source", "")

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] insert_mermaid_diagram → doc: %s | target: %d | png: %s", docID, targetIndex, pngPath)

	result, err := gdocs.InsertMermaidDiagram(ctx, docID, targetIndex, pngPath, caption, source)
	if err != nil {
		log.Printf("[ERRO] insert_mermaid_diagram: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	resJSON, _ := json.MarshalIndent(result, "", "  ")
	return mcplib.NewToolResultText(string(resJSON)), nil
}

// ─── Tool 17: check_pdf_sync ────────────────────────────────────────────────

func registerCheckPdfSync(s *server.MCPServer) {
	tool := mcplib.NewTool("check_pdf_sync",
		mcplib.WithDescription(
			"Verifica a paridade criptográfica entre o PDF compilado localmente (TCC) e o arquivo no Google Drive via hash MD5 com validação prévia de integridade ISO 32000-1 (%%EOF) e TOCTOU."),
		mcplib.WithString("local_pdf_path",
			mcplib.Required(),
			mcplib.Description("Caminho absoluto ou relativo do arquivo PDF local compilado (ex: 'resources/tcc/main.pdf')"),
		),
		mcplib.WithString("target_pdf_file_id",
			mcplib.Required(),
			mcplib.Description("ID do arquivo PDF no Google Drive"),
		),
	)
	s.AddTool(tool, handleCheckPdfSync)
	log.Println("[INFO] Ferramenta registrada: check_pdf_sync")
}

func handleCheckPdfSync(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	localPDFPath, err := request.RequireString("local_pdf_path")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'local_pdf_path' é obrigatório."), nil
	}

	targetPdfFileID, err := request.RequireString("target_pdf_file_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'target_pdf_file_id' é obrigatório."), nil
	}

	log.Printf("[INFO] check_pdf_sync → local: %s | drive: %s", localPDFPath, targetPdfFileID)

	report, err := gdocs.CheckPDFSync(ctx, localPDFPath, targetPdfFileID)
	if err != nil {
		log.Printf("[ERRO] check_pdf_sync: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	resJSON, _ := json.MarshalIndent(report, "", "  ")
	return mcplib.NewToolResultText(string(resJSON)), nil
}

// ─── Tool 18: lint_document_abnt ───────────────────────────────────────────

func registerLintDocumentAbnt(s *server.MCPServer) {
	tool := mcplib.NewTool("lint_document_abnt",
		mcplib.WithDescription(
			"Diagnostica e higieniza anomalias no Google Docs segundo as normas ABNT e boas práticas editoriais: detecta tabelas vazias ou com células faltantes, cascatas de negrito (>150 caracteres contínuos no corpo do texto), quebras manuais de linha (\\u000b), quebras múltiplas consecutivas e marcadores soltos. O parâmetro 'auto_fix' (default: false) permite expurgar automaticamente tabelas vazias e normalizar quebras."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		mcplib.WithBoolean("auto_fix",
			mcplib.Description("Se true, executa automaticamente a exclusão de tabelas vazias e a normalização de quebras de linha manuais (default: false)"),
		),
	)
	s.AddTool(tool, handleLintDocumentAbnt)
	log.Println("[INFO] Ferramenta registrada: lint_document_abnt")
}

func handleLintDocumentAbnt(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	autoFix := request.GetBool("auto_fix", false)

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] lint_document_abnt → doc: %s | autoFix: %v", docID, autoFix)

	report, err := gdocs.LintDocumentABNT(ctx, docID, autoFix)
	if err != nil {
		log.Printf("[ERRO] lint_document_abnt: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	resJSON, _ := json.MarshalIndent(report, "", "  ")
	return mcplib.NewToolResultText(string(resJSON)), nil
}

