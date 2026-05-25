// Package mcp registra e conecta as ferramentas do servidor MCP,
// servindo como ponte entre o protocolo MCP e os módulos internos.
package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
	"github.com/JGustavoCN/mcp-gdocs-latex/internal/latex"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registra todas as 9 ferramentas no servidor MCP.
func RegisterTools(s *server.MCPServer) {
	registerListAvailableDocuments(s)
	registerGetDocSkeleton(s)
	registerReadDocContent(s)
	registerSearchInDoc(s)
	registerListDocComments(s)
	registerReplyToComment(s)
	registerResolveComment(s)
	registerReplaceTextInDoc(s)
	registerMultiReplaceDocBlock(s)
	registerSyncPDFToDrive(s)
	registerUpdateLocalLatex(s)
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

	// Extrai, valida allowlist e faz ping na API
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

	// Extrai, valida allowlist e faz ping na API
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

// ─── Tool 8: replace_text_in_doc (Cirúrgico) ────────────────────────────────

func registerReplaceTextInDoc(s *server.MCPServer) {
	tool := mcplib.NewTool("replace_text_in_doc",
		mcplib.WithDescription(
			"Realiza a substituição de termos de texto no documento. Permite informar o parâmetro 'occurrence_index' para alterar uma ocorrência específica (ex: 0 para a primeira ocorrência, 1 para a segunda). Se definido como -1 ou omitido, realiza a substituição de todas as correspondências encontradas."),
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
		mcplib.WithInteger("occurrence_index",
			mcplib.Description("Índice da ocorrência a substituir (0 = primeira, 1 = segunda, etc). Se -1, substitui todas."),
		),
	)
	s.AddTool(tool, handleReplaceTextInDoc)
	log.Println("[INFO] Ferramenta registrada: replace_text_in_doc")
}

func handleReplaceTextInDoc(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
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

	// occurrence_index é opcional, se não enviado, é -1
	occIndex := request.GetInt("occurrence_index", -1)

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] replace_text_in_doc → doc: %s | substituir: '%s' → '%s' (ocorrência: %d)",
		docID, oldText, newText, occIndex)

	result, err := gdocs.ReplaceTextInDoc(ctx, docID, oldText, newText, occIndex)
	if err != nil {
		log.Printf("[ERRO] replace_text_in_doc: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 5: multi_replace_doc_block ────────────────────────────────────────

func registerMultiReplaceDocBlock(s *server.MCPServer) {
	tool := mcplib.NewTool("multi_replace_doc_block",
		mcplib.WithDescription(
			"Aplica em lote substituições cirúrgicas de texto em índices absolutos de caracteres. Requer validação exata do texto esperado (expected_text) em cada intervalo informado. Executa as alterações de trás para frente no documento para neutralizar o deslocamento de índices (offset shift) durante a operação."),
		mcplib.WithString("document_url_or_id",
			mcplib.Required(),
			mcplib.Description("URL completa do Google Docs ou apenas o ID do documento"),
		),
		// MCP-GO API requires chunks to be an array, but the SDK may not fully expose WithArray nicely without custom schemas.
		// We can define it as an object that contains an array, or a stringified JSON.
		// Let's use a string argument for the JSON chunks to ensure wide compatibility with the Go SDK.
		mcplib.WithString("chunks_json",
			mcplib.Required(),
			mcplib.Description("JSON array stringified com blocos: [{'start_index': 10, 'end_index': 15, 'expected_text': 'abc', 'replacement_text': 'def'}]"),
		),
	)
	s.AddTool(tool, handleMultiReplaceDocBlock)
	log.Println("[INFO] Ferramenta registrada: multi_replace_doc_block")
}

func handleMultiReplaceDocBlock(ctx context.Context, request mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	input, err := request.RequireString("document_url_or_id")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'document_url_or_id' é obrigatório."), nil
	}

	chunksJSON, err := request.RequireString("chunks_json")
	if err != nil {
		return mcplib.NewToolResultError("Parâmetro 'chunks_json' é obrigatório."), nil
	}

	var chunks []gdocs.Chunk
	if err := json.Unmarshal([]byte(chunksJSON), &chunks); err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: Falha ao fazer parse do chunks_json: %v", err)), nil
	}

	docID, err := gdocs.ExtractAndValidateDocID(ctx, input)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	log.Printf("[INFO] multi_replace_doc_block → doc: %s | chunks: %d", docID, len(chunks))

	result, err := gdocs.MultiReplaceDocBlock(ctx, docID, chunks)
	if err != nil {
		log.Printf("[ERRO] multi_replace_doc_block: %v", err)
		return mcplib.NewToolResultError(fmt.Sprintf("ERRO: %v", err)), nil
	}

	return mcplib.NewToolResultText(result), nil
}

// ─── Tool 9: sync_pdf_to_drive ──────────────────────────────────────────────

func registerSyncPDFToDrive(s *server.MCPServer) {
	tool := mcplib.NewTool("sync_pdf_to_drive",
		mcplib.WithDescription(
			"Faz o upload do PDF compilado sobrescrevendo os bytes de um arquivo existente. Esta ferramenta não cria arquivos. O LLM deve extrair o ID do PDF usando a ferramenta 'list_available_documents' e usá-lo no parâmetro 'target_pdf_file_id'."),
		mcplib.WithString("local_pdf_path",
			mcplib.Required(),
			mcplib.Description("Caminho absoluto ou relativo do arquivo PDF local (ex: 'src/main.pdf')"),
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
