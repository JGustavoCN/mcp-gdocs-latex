package gdocs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/latex"
	"github.com/ledongthuc/pdf"
	"google.golang.org/api/docs/v1"
)

// PreTextualPaginationReport detalha a auditoria física real das páginas extraídas do PDF renderizado pelo Google Docs.
type PreTextualPaginationReport struct {
	IsConform           bool     `json:"is_conform"`
	TotalPDFPages       int      `json:"total_pdf_pages"`
	CapaPage            int      `json:"capa_page"`
	CapaAnoPage         int      `json:"capa_ano_page"`
	FolhaRostoPage      int      `json:"folha_rosto_page"`
	FolhaRostoAnoPage   int      `json:"folha_rosto_ano_page"`
	ResumoPage          int      `json:"resumo_page"`
	CapaOverflow        bool     `json:"capa_overflow"`
	FolhaRostoOverflow  bool     `json:"folha_rosto_overflow"`
	ResumoPremature     bool     `json:"resumo_premature"`
	Details             []string `json:"details"`
}

// ExportDocPDF exporta o documento Google Docs como PDF binário oficial renderizado pelo motor do Google Drive.
func ExportDocPDF(ctx context.Context, docID string) ([]byte, error) {
	driveSvc, err := GetDriveService(ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter drive service: %w", err)
	}

	resp, err := driveSvc.Files.Export(docID, "application/pdf").Download()
	if err != nil {
		return nil, fmt.Errorf("erro ao exportar documento para PDF via Drive API: %w", err)
	}
	defer resp.Body.Close()

	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler stream do PDF: %w", err)
	}
	if len(pdfBytes) == 0 {
		return nil, fmt.Errorf("PDF exportado retornou vazio")
	}

	return pdfBytes, nil
}

// ExtractPDFPages extrai o texto de cada página física do PDF.
func ExtractPDFPages(pdfBytes []byte) ([]string, error) {
	if len(pdfBytes) == 0 {
		return nil, fmt.Errorf("buffer PDF vazio")
	}

	reader, err := pdf.NewReader(bytes.NewReader(pdfBytes), int64(len(pdfBytes)))
	if err != nil {
		return nil, fmt.Errorf("falha ao instanciar reader de PDF: %w", err)
	}

	numPages := reader.NumPage()
	pages := make([]string, numPages)

	for i := 1; i <= numPages; i++ {
		page := reader.Page(i)
		if page.V.IsNull() {
			continue
		}
		txt, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}
		pages[i-1] = txt
	}

	return pages, nil
}

// VerifyPreTextualPagination analisa as páginas reais do PDF e verifica se a Capa está 100% contida na Página 1,
// a Folha de Rosto na Página 2 e o RESUMO iniciando na Página 3.
func VerifyPreTextualPagination(ctx context.Context, docID string, meta *latex.TCCMetadata, doc *docs.Document) (*PreTextualPaginationReport, error) {
	pdfBytes, err := ExportDocPDF(ctx, docID)
	if err != nil {
		return nil, err
	}

	pages, err := ExtractPDFPages(pdfBytes)
	if err != nil {
		return nil, err
	}

	report := &PreTextualPaginationReport{
		IsConform:     true,
		TotalPDFPages: len(pages),
	}

	st := identifyPreTextualStructure(doc, meta)

	// Termos canônicos de busca extraídos da estrutura real do documento
	anoText := ""
	if st.AnoCapaEl != nil {
		anoText = strings.TrimSpace(extractParagraphText(st.AnoCapaEl))
	}
	if anoText == "" && meta != nil {
		anoText = strings.TrimSpace(meta.Data)
	}

	cidadeText := ""
	if st.CidadeCapaEl != nil {
		cidadeText = strings.TrimSpace(extractParagraphText(st.CidadeCapaEl))
	}
	if cidadeText == "" && meta != nil {
		cidadeText = strings.TrimSpace(meta.Local)
	}

	autorText := ""
	if st.AutorCapaEl != nil {
		autorText = strings.TrimSpace(extractParagraphText(st.AutorCapaEl))
	}
	if autorText == "" && meta != nil {
		autorText = strings.TrimSpace(meta.Autor)
	}

	// 1. Rastreia em quais páginas ocorrem o ano e o autor
	type matchLoc struct {
		page int
	}

	var anoMatches []matchLoc
	var autorMatches []matchLoc
	resumoPage := 0

	for pageNum := 1; pageNum <= len(pages); pageNum++ {
		pageTxt := pages[pageNum-1]
		upperTxt := strings.ToUpper(pageTxt)

		if anoText != "" && strings.Contains(pageTxt, anoText) {
			anoMatches = append(anoMatches, matchLoc{page: pageNum})
		}
		if autorText != "" && strings.Contains(upperTxt, strings.ToUpper(autorText)) {
			autorMatches = append(autorMatches, matchLoc{page: pageNum})
		}
		if resumoPage == 0 && (strings.Contains(upperTxt, "RESUMO") || strings.Contains(upperTxt, "ABSTRACT")) {
			// Garante que não é uma citação acidental na capa
			if pageNum >= 2 {
				resumoPage = pageNum
			}
		}
	}

	// Avalia a Capa
	report.CapaPage = 1
	if len(anoMatches) >= 1 {
		report.CapaAnoPage = anoMatches[0].page
	}
	if report.CapaAnoPage > 1 {
		report.CapaOverflow = true
		report.IsConform = false
		report.Details = append(report.Details, fmt.Sprintf("Estouro na Capa: o Ano (%s) vazou para a Página %d (esperado: Página 1)", anoText, report.CapaAnoPage))
	}

	// Avalia a Folha de Rosto
	if len(autorMatches) >= 2 {
		report.FolhaRostoPage = autorMatches[1].page
	} else if len(autorMatches) == 1 && autorMatches[0].page == 2 {
		report.FolhaRostoPage = 2
	}
	if len(anoMatches) >= 2 {
		report.FolhaRostoAnoPage = anoMatches[1].page
	}

	if report.FolhaRostoAnoPage > 2 {
		report.FolhaRostoOverflow = true
		report.IsConform = false
		report.Details = append(report.Details, fmt.Sprintf("Estouro na Folha de Rosto: o Ano (%s) vazou para a Página %d (esperado: Página 2)", anoText, report.FolhaRostoAnoPage))
	}

	// Avalia o RESUMO
	report.ResumoPage = resumoPage
	if resumoPage > 0 && resumoPage < 3 {
		report.ResumoPremature = true
		report.IsConform = false
		report.Details = append(report.Details, fmt.Sprintf("RESUMO iniciou prematuramente na Página %d (esperado: Página 3)", resumoPage))
	} else if resumoPage > 3 {
		report.Details = append(report.Details, fmt.Sprintf("Aviso: RESUMO iniciou na Página %d (páginas intermediárias detectadas)", resumoPage))
	}

	if report.IsConform {
		report.Details = append(report.Details, fmt.Sprintf("Conformidade Física Perfeita: Capa (Página 1), Folha de Rosto (Página 2) e RESUMO (Página %d). Total de páginas: %d.", report.ResumoPage, report.TotalPDFPages))
	}

	return report, nil
}

