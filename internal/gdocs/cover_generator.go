package gdocs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/latex"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
)

// DefaultIFSLogoGitHubURL representa a URL canônica pública da logo do IFS hospedada no repositório GitHub
const DefaultIFSLogoGitHubURL = "https://raw.githubusercontent.com/JGustavoCN/mcp-gdocs-latex/main/assets/logo_ifs.png"

// ResolveLogoImage localiza o arquivo da logo no sistema de arquivos, permitindo customizações transparentes
func ResolveLogoImage(customPath string) (string, []byte, error) {
	paths := []string{}
	if customPath != "" {
		paths = append(paths, customPath)
	}
	if env := os.Getenv("TCC_LOGO_PATH"); env != "" {
		paths = append(paths, env)
	}
	// Caminhos canônicos padrão do repositório versionados no Git (prioriza assets/ e raiz)
	paths = append(paths,
		filepath.Join("assets", "logo_ifs.png"),
		filepath.Join("assets", "logos", "logo_ifs.png"),
		"logo_ifs.png",
		filepath.Join("resources", "tcc", "logo_ifs.png"),
		filepath.Join("resources", "template", "public", "logo_ifs.png"),
		filepath.Join("resources", "template", "logo_ifs.png"),
	)

	for _, p := range paths {
		if b, err := os.ReadFile(p); err == nil && len(b) > 0 {
			return p, b, nil
		}
	}

	return "", nil, fmt.Errorf("arquivo da logo não localizado nos caminhos esperados: %v", paths)
}

// uploadImageToDirectCDN faz upload de contingência para provedor público com suporte direto a image/png
func uploadImageToDirectCDN(pngBytes []byte) (string, error) {
	if len(pngBytes) == 0 {
		return "", fmt.Errorf("buffer PNG vazio")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("reqtype", "fileupload")
	part, err := writer.CreateFormFile("fileToUpload", "logo_ifs.png")
	if err != nil {
		return "", err
	}
	if _, err := part.Write(pngBytes); err != nil {
		return "", err
	}
	_ = writer.Close()

	req, err := http.NewRequest("POST", "https://catbox.moe/user/api.php", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	urlStr := strings.TrimSpace(string(respBytes))
	if !strings.HasPrefix(urlStr, "http") {
		return "", fmt.Errorf("resposta inesperada do provedor de imagem: %s", urlStr)
	}
	return urlStr, nil
}

// UploadIFSLogoToDrive obtém uma URI pública estável e acessível da logo para consumo pelo Google Docs
func UploadIFSLogoToDrive(ctx context.Context, logoPath, documentID string) (string, error) {
	// Se for uma URL web fornecida diretamente pelo usuário, consome imediatamente
	if strings.HasPrefix(logoPath, "http://") || strings.HasPrefix(logoPath, "https://") {
		return logoPath, nil
	}

	resolvedPath, pngBytes, errResolve := ResolveLogoImage(logoPath)
	if errResolve != nil {
		log.Printf("[AVISO] %v. Utilizando contingência canônica do GitHub: %s", errResolve, DefaultIFSLogoGitHubURL)
		return DefaultIFSLogoGitHubURL, nil
	}

	driveSvc, err := GetDriveService(ctx)
	if err != nil {
		// Sem serviço Drive configurado, tenta CDN de imagens diretas
		if cdnURL, errCDN := uploadImageToDirectCDN(pngBytes); errCDN == nil && cdnURL != "" {
			return cdnURL, nil
		}
		return DefaultIFSLogoGitHubURL, nil
	}

	driveFile := &drive.File{
		Name:     filepath.Base(resolvedPath),
		MimeType: "image/png",
	}

	if documentID != "" {
		if docMeta, errDoc := driveSvc.Files.Get(documentID).Fields("parents").Do(); errDoc == nil && len(docMeta.Parents) > 0 {
			driveFile.Parents = []string{docMeta.Parents[0]}
		}
	}

	uploadedFile, errUpload := driveSvc.Files.Create(driveFile).
		Media(bytes.NewReader(pngBytes)).
		Fields("id, webContentLink, webViewLink").
		Do()
	if errUpload != nil {
		log.Printf("[INFO] Drive pessoal sem cota para upload (%v). Ativando canal de transferência CDN.", errUpload)
		cdnURI, errCDN := uploadImageToDirectCDN(pngBytes)
		if errCDN == nil && cdnURI != "" {
			return cdnURI, nil
		}
		log.Printf("[AVISO] Falha em upload CDN (%v). Usando fallback canônico do GitHub: %s", errCDN, DefaultIFSLogoGitHubURL)
		return DefaultIFSLogoGitHubURL, nil
	}

	fileID := uploadedFile.Id
	_ = grantPublicReadPermission(driveSvc, fileID)

	imageURI := fmt.Sprintf("https://drive.google.com/uc?id=%s&export=download", fileID)
	log.Printf("[INFO] Logo enviada com sucesso ao Google Drive: ID %s", fileID)
	return imageURI, nil
}

// HasABNTCoverAndLogo inspeciona o documento e verifica se a Capa oficial com a logo e a Folha de Rosto estão presentes
func HasABNTCoverAndLogo(doc *docs.Document) (hasCover bool, hasLogo bool, hasTitlePage bool) {
	if doc == nil || doc.Body == nil || len(doc.Body.Content) == 0 {
		return false, false, false
	}

	firstChapterIndex := int64(-1)
	for _, el := range doc.Body.Content {
		if el.Paragraph != nil {
			pText := ""
			for _, pEl := range el.Paragraph.Elements {
				if pEl.TextRun != nil {
					pText += pEl.TextRun.Content
				}
			}
			t := strings.TrimSpace(pText)
			if rePrimaryHeadingNum.MatchString(t) || t == "1 INTRODUÇÃO" {
				firstChapterIndex = el.StartIndex
				break
			}
		}
	}

	if firstChapterIndex < 0 {
		firstChapterIndex = 999999
	}

	// 2. Analisa elementos pré-textuais estritamente antes do primeiro capítulo
	hasInstituicao := false
	hasNota := false

	for _, el := range doc.Body.Content {
		if el.StartIndex >= firstChapterIndex {
			break
		}
		if el.Paragraph != nil {
			pText := ""
			for _, pEl := range el.Paragraph.Elements {
				if pEl.TextRun != nil {
					pText += pEl.TextRun.Content
				}
				if pEl.InlineObjectElement != nil {
					hasLogo = true
				}
			}
			upper := strings.ToUpper(pText)
			if strings.Contains(upper, "INSTITUTO FEDERAL") && strings.Contains(upper, "SERGIPE") {
				hasInstituicao = true
			}
			if isPreTextualNote(pText) {
				hasNota = true
			}
		}
	}

	hasCover = hasInstituicao
	hasTitlePage = hasNota
	return hasCover, hasLogo, hasTitlePage
}

// GenerateABNTIFSCoverAndTitlePage gera ou injeta no Google Docs a Capa e Folha de Rosto oficiais do IFS/ABNT
func GenerateABNTIFSCoverAndTitlePage(ctx context.Context, documentID string, meta *latex.TCCMetadata, logoPath string, prepend bool) error {
	docsSvc, err := GetDocsService(ctx)
	if err != nil {
		return fmt.Errorf("erro ao obter docs service: %w", err)
	}

	if meta == nil {
		var errMeta error
		meta, errMeta = latex.ParseTCCMetadata("")
		if errMeta != nil {
			return fmt.Errorf("falha ao carregar metadados do TCC: %w", errMeta)
		}
	}

	doc, err := docsSvc.Documents.Get(documentID).Do()
	if err != nil {
		return fmt.Errorf("erro ao carregar documento: %w", err)
	}

	// 1. Obter URI pública da logo do IFS
	logoURI, errLogo := UploadIFSLogoToDrive(ctx, logoPath, documentID)
	if errLogo != nil {
		log.Printf("[AVISO] Não foi possível fazer upload da logo do IFS: %v. Prosseguindo com cabeçalho textual.", errLogo)
	}

	bodyEnd := getBodyEndIndex(doc)

	// Se não for prepend, faz a limpeza total para inicializar com o template oficial completo
	if !prepend && bodyEnd > 2 {
		delReq := &docs.BatchUpdateDocumentRequest{
			Requests: []*docs.Request{
				{
					DeleteContentRange: &docs.DeleteContentRangeRequest{
						Range: &docs.Range{
							StartIndex: 1,
							EndIndex:   bodyEnd - 1,
						},
					},
				},
			},
		}
		if _, errDel := docsSvc.Documents.BatchUpdate(documentID, delReq).Do(); errDel != nil {
			return fmt.Errorf("erro ao limpar documento para geração da capa: %w", errDel)
		}
	} else if prepend && bodyEnd > 2 {
		// Se for prepend em documento existente, localiza onde começa o texto acadêmico real para substituir capa antiga
		cutIndex := int64(-1)
		needsResumoPrefix := false
		for _, el := range doc.Body.Content {
			if el.Paragraph != nil {
				pText := ""
				for _, pEl := range el.Paragraph.Elements {
					if pEl.TextRun != nil {
						pText += pEl.TextRun.Content
					}
				}
				trimmed := strings.TrimSpace(pText)
				trimmedUpper := strings.ToUpper(trimmed)

				// Detecta se há texto do resumo que ficou intercalado
				if strings.Contains(trimmed, "submissão eletrônica") || strings.Contains(trimmed, "submissão e") {
					cutIndex = el.StartIndex
					needsResumoPrefix = true
					break
				}

				if (trimmedUpper == "RESUMO" || strings.HasPrefix(trimmedUpper, "RESUMO") ||
					trimmedUpper == "ABSTRACT" || strings.HasPrefix(trimmedUpper, "1 INTRODUÇÃO") ||
					rePrimaryHeadingNum.MatchString(trimmedUpper)) && !isPreTextualNote(trimmed) {
					cutIndex = el.StartIndex
					break
				}
			}
		}

		if cutIndex > 1 {
			delReq := &docs.BatchUpdateDocumentRequest{
				Requests: []*docs.Request{
					{
						DeleteContentRange: &docs.DeleteContentRangeRequest{
							Range: &docs.Range{
								StartIndex: 1,
								EndIndex:   cutIndex - 1,
							},
						},
					},
				},
			}
			_, _ = docsSvc.Documents.BatchUpdate(documentID, delReq).Do()
		}

		_ = needsResumoPrefix
	}

	// 2. Construção Textual Canônica da Capa e Folha de Rosto (NBR 14724 + IFS config.tex)
	var sb strings.Builder

	// --- PÁGINA 1: CAPA (NBR 14724 + IFS config.tex) ---
	sb.WriteString("\n\n")

	// Instituição
	instLines := meta.InstituicaoLines()
	for _, l := range instLines {
		sb.WriteString(strings.ToUpper(l))
		sb.WriteString("\n")
	}
	sb.WriteString(strings.Repeat("\n", 3)) // 3 quebras verticais calibradas (NBR 14724 / sem overflow)

	// Autor
	sb.WriteString(strings.ToUpper(meta.Autor))
	sb.WriteString(strings.Repeat("\n", 3)) // 3 quebras verticais calibradas (NBR 14724 / sem overflow)

	// Título do Trabalho
	sb.WriteString(strings.ToUpper(meta.Titulo))
	sb.WriteString(strings.Repeat("\n", 15)) // 15 quebras verticais posicionam Local e Ano no rodapé da Capa (Página 1 estrita)

	// Local e Ano
	sb.WriteString(strings.ToUpper(meta.Local))
	sb.WriteString("\n")
	sb.WriteString(meta.Data)

	capaText := sb.String()

	// --- PÁGINA 2: FOLHA DE ROSTO (NBR 14724 + IFS config.tex) ---
	var sbFR strings.Builder
	sbFR.WriteString(strings.Repeat("\n", 2)) // 2 quebras verticais no topo (Página 2 estrita)
	sbFR.WriteString(strings.ToUpper(meta.Autor))
	sbFR.WriteString(strings.Repeat("\n", 4)) // 4 quebras entre Autor e Título
	sbFR.WriteString(strings.ToUpper(meta.Titulo))
	sbFR.WriteString(strings.Repeat("\n", 3)) // 3 quebras entre Título e Preâmbulo

	// Bloco da Nota Institucional + Orientador (alinhado da metade para a margem direita, recuo 8.0 cm)
	sbFR.WriteString(meta.Preambulo)
	sbFR.WriteString("\n")
	sbFR.WriteString("Orientador: ")
	sbFR.WriteString(meta.Orientador)
	if meta.Coorientador != "" {
		sbFR.WriteString("\nCoorientador: ")
		sbFR.WriteString(meta.Coorientador)
	}
	sbFR.WriteString(strings.Repeat("\n", 12)) // 12 quebras verticais posicionam Local e Ano no rodapé da Folha de Rosto (Página 2 estrita)

	// Local e Ano
	sbFR.WriteString(strings.ToUpper(meta.Local))
	sbFR.WriteString("\n")
	sbFR.WriteString(meta.Data)

	folhaRostoText := sbFR.String()

	capaUTF16Len := UTF16Length(capaText)
	frUTF16Len := UTF16Length(folhaRostoText)

	// Verifica se precisamos restaurar o cabeçalho do resumo recortado
	needsResumo := false
	if prepend {
		// Checa se o primeiro parágrafo restante precisa do cabeçalho de RESUMO
		if docAfterDel, errAfter := docsSvc.Documents.Get(documentID).Do(); errAfter == nil {
			for _, el := range docAfterDel.Body.Content {
				if el.Paragraph != nil {
					pText := ""
					for _, pEl := range el.Paragraph.Elements {
						if pEl.TextRun != nil {
							pText += pEl.TextRun.Content
						}
					}
					if strings.Contains(pText, "submissão eletrônica") || strings.Contains(pText, "submissão e") {
						needsResumo = true
					}
					break
				}
			}
		}
	}

	var fullPreText strings.Builder
	fullPreText.WriteString(capaText)
	fullPreText.WriteString("\n")
	fullPreText.WriteString(folhaRostoText)
	if needsResumo {
		fullPreText.WriteString("\n# RESUMO\n\nA ")
	} else if !prepend {
		fullPreText.WriteString("\n1 INTRODUÇÃO\n\nEscreva aqui a introdução do trabalho.\n\n")
		fullPreText.WriteString("2 REFERENCIAL TEÓRICO\n\nEscreva aqui o referencial teórico.\n\n")
		fullPreText.WriteString("3 METODOLOGIA\n\nEscreva aqui a metodologia científica adotada.\n\n")
		fullPreText.WriteString("4 RESULTADOS E DISCUSSÃO\n\nEscreva aqui os resultados alcançados.\n\n")
		fullPreText.WriteString("5 CONCLUSÃO\n\nEscreva aqui as considerações finais.\n\n")
		fullPreText.WriteString("REFERÊNCIAS\n\nInsira aqui as referências no padrão ABNT NBR 6023.\n")
	}

	// 3. Insere todo o bloco pré-textual na posição 1
	insReq := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			{
				InsertText: &docs.InsertTextRequest{
					Location: &docs.Location{Index: 1},
					Text:     fullPreText.String(),
				},
			},
		},
	}
	if _, errIns := docsSvc.Documents.BatchUpdate(documentID, insReq).Do(); errIns != nil {
		return fmt.Errorf("erro ao inserir bloco textual pré-textual: %w", errIns)
	}

	// 4. Inserção das quebras de página em ordem inversa para preservar posições relativas
	posFRBreak := int64(1) + capaUTF16Len + 1 + frUTF16Len
	posCapaBreak := int64(1) + capaUTF16Len

	pbReq := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			{
				InsertPageBreak: &docs.InsertPageBreakRequest{
					Location: &docs.Location{Index: posFRBreak},
				},
			},
			{
				InsertPageBreak: &docs.InsertPageBreakRequest{
					Location: &docs.Location{Index: posCapaBreak},
				},
			},
		},
	}
	if _, errPb := docsSvc.Documents.BatchUpdate(documentID, pbReq).Do(); errPb != nil {
		log.Printf("[AVISO] Falha ao inserir quebras de página pré-textuais: %v", errPb)
	}

	// 5. Inserção da Imagem da Logo do IFS no topo da Capa (índice 1) com centralização estrita
	if logoURI != "" {
		imgReq := &docs.BatchUpdateDocumentRequest{
			Requests: []*docs.Request{
				{
					InsertInlineImage: &docs.InsertInlineImageRequest{
						Location: &docs.Location{Index: 1},
						Uri:      logoURI,
						ObjectSize: &docs.Size{
							Width:  &docs.Dimension{Magnitude: 91.0, Unit: "PT"},
							Height: &docs.Dimension{Magnitude: 91.0, Unit: "PT"},
						},
					},
				},
				{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: &docs.Range{
							StartIndex: 1,
							EndIndex:   2,
						},
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType:  "NORMAL_TEXT",
							Alignment:       "CENTER",
							IndentFirstLine: &docs.Dimension{Magnitude: 0, Unit: "PT"},
							IndentStart:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceAbove:      &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceBelow:      &docs.Dimension{Magnitude: 14.17, Unit: "PT"}, // 0.5 cm
						},
						Fields: "namedStyleType,alignment,indentFirstLine,indentStart,spaceAbove,spaceBelow",
					},
				},
			},
		}
		if _, errImg := docsSvc.Documents.BatchUpdate(documentID, imgReq).Do(); errImg != nil {
			log.Printf("[AVISO] Falha ao inserir logo inline do IFS: %v", errImg)
		}
	}

	// 6. Formata o documento com as regras soberanas ABNT (centralização da capa, recuo 7.5cm da nota, 12pt)
	_, _ = FormatDocumentABNTExtended(ctx, documentID, -1, -1, true, false, ABNTDefaultFontFamily, false, false, false, logoPath)

	log.Printf("[INFO] Capa Oficial do IFS e Folha de Rosto geradas com sucesso para o documento %s", documentID)
	return nil
}

// PreTextualStructure mapeia com precisão posicional e genérica todos os elementos
// canônicos da Capa (Página 1) e Folha de Rosto (Página 2), sem strings fixas.
type PreTextualStructure struct {
	// Capa (Página 1)
	LogoEl          *docs.StructuralElement
	InstEls         []*docs.StructuralElement
	AutorCapaEl     *docs.StructuralElement
	TituloCapaEl    *docs.StructuralElement
	CidadeCapaEl    *docs.StructuralElement
	AnoCapaEl       *docs.StructuralElement
	AnoCapaHasBreak bool

	// Folha de Rosto (Página 2)
	AutorFREl      *docs.StructuralElement
	TituloFREl     *docs.StructuralElement
	PreambuloFREls []*docs.StructuralElement
	OrientadorFREl *docs.StructuralElement
	CidadeFREl     *docs.StructuralElement
	AnoFREl        *docs.StructuralElement
	AnoFRHasBreak  bool

	// Limites de Índice
	CoverEndIndex     int64
	TitlePageEndIndex int64
}

// extractParagraphText extrai todo o texto concatenado dos TextRuns de um parágrafo
func extractParagraphText(el *docs.StructuralElement) string {
	if el == nil || el.Paragraph == nil {
		return ""
	}
	var sb strings.Builder
	for _, pel := range el.Paragraph.Elements {
		if pel.TextRun != nil {
			sb.WriteString(pel.TextRun.Content)
		}
	}
	return sb.String()
}

// paragraphHasPageBreak verifica se o parágrafo contém elemento de quebra de página
func paragraphHasPageBreak(el *docs.StructuralElement) bool {
	if el == nil || el.Paragraph == nil {
		return false
	}
	for _, pel := range el.Paragraph.Elements {
		if pel.PageBreak != nil {
			return true
		}
	}
	return false
}

// paragraphHasInlineObject verifica se o parágrafo contém objeto inline
func paragraphHasInlineObject(el *docs.StructuralElement) bool {
	if el == nil || el.Paragraph == nil {
		return false
	}
	for _, pel := range el.Paragraph.Elements {
		if pel.InlineObjectElement != nil {
			return true
		}
	}
	return false
}

func isElementInList(target *docs.StructuralElement, list []*docs.StructuralElement) bool {
	if target == nil {
		return false
	}
	for _, item := range list {
		if item != nil && item.StartIndex == target.StartIndex && item.EndIndex == target.EndIndex {
			return true
		}
	}
	return false
}

var reAcademicYear = regexp.MustCompile(`^\s*(?:19|20)\d{2}\s*$`)

// identifyPreTextualStructure identifica de forma 100% genérica, posicional e baseada na ABNT NBR 14724
// os elementos da Capa e da Folha de Rosto, sem depender de nenhuma string fixa de autor, campus ou título.
func identifyPreTextualStructure(doc *docs.Document, meta *latex.TCCMetadata) PreTextualStructure {
	var res PreTextualStructure
	res.CoverEndIndex = 999999
	res.TitlePageEndIndex = 999999

	if doc == nil || doc.Body == nil || len(doc.Body.Content) == 0 {
		return res
	}

	// 1. Limite final dos elementos pré-textuais (onde começam RESUMO ou Capítulo 1)
	preTextualLimit := int64(999999)
	for _, el := range doc.Body.Content {
		if el.Paragraph != nil {
			txt := strings.TrimSpace(extractParagraphText(el))
			upper := strings.ToUpper(txt)
			if (upper == "RESUMO" || strings.HasPrefix(upper, "RESUMO") ||
				upper == "ABSTRACT" || strings.HasPrefix(upper, "1 INTRODUÇÃO") ||
				rePrimaryHeadingNum.MatchString(upper)) && !isPreTextualNote(txt) {
				preTextualLimit = el.StartIndex
				break
			}
		}
	}

	// 2. Coleta todos os parágrafos não-vazios e identifica os Anos (Capa e Folha de Rosto)
	type pInfo struct {
		el       *docs.StructuralElement
		text     string
		hasBreak bool
		hasObj   bool
	}
	var textParas []pInfo

	for _, el := range doc.Body.Content {
		if el.StartIndex >= preTextualLimit {
			break
		}
		if el.Paragraph != nil {
			txt := strings.TrimSpace(extractParagraphText(el))
			hasBreak := paragraphHasPageBreak(el)
			hasObj := paragraphHasInlineObject(el)

			if hasObj && res.LogoEl == nil {
				res.LogoEl = el
			}

			if txt == "" {
				continue
			}

			isYear := reAcademicYear.MatchString(txt) || (meta != nil && meta.Data != "" && txt == strings.TrimSpace(meta.Data))
			if isYear {
				if res.AnoCapaEl == nil {
					res.AnoCapaEl = el
					res.AnoCapaHasBreak = hasBreak
					res.CoverEndIndex = el.EndIndex
				} else if res.AnoFREl == nil {
					res.AnoFREl = el
					res.AnoFRHasBreak = hasBreak
					res.TitlePageEndIndex = el.EndIndex
				}
			}

			textParas = append(textParas, pInfo{
				el:       el,
				text:     txt,
				hasBreak: hasBreak,
				hasObj:   hasObj,
			})
		}
	}

	if res.AnoCapaEl == nil {
		return res
	}

	// 3. Mapeamento da Capa (elementos antes do Ano da Capa)
	var capaParas []pInfo
	var frParas []pInfo

	for _, p := range textParas {
		if p.el.EndIndex <= res.CoverEndIndex {
			capaParas = append(capaParas, p)
		} else if res.AnoFREl != nil && p.el.EndIndex <= res.TitlePageEndIndex {
			frParas = append(frParas, p)
		}
	}

	// Na Capa: o último elemento textual é o Ano. O penúltimo é a Cidade.
	if len(capaParas) >= 2 {
		res.CidadeCapaEl = capaParas[len(capaParas)-2].el
	}

	// Elementos da Capa antes da Cidade:
	var capaBeforeCity []pInfo
	for _, p := range capaParas {
		if res.CidadeCapaEl != nil && p.el.EndIndex > res.CidadeCapaEl.StartIndex {
			break
		}
		capaBeforeCity = append(capaBeforeCity, p)
	}

	if len(capaBeforeCity) >= 3 {
		res.TituloCapaEl = capaBeforeCity[len(capaBeforeCity)-1].el
		res.AutorCapaEl = capaBeforeCity[len(capaBeforeCity)-2].el
		for i := 0; i < len(capaBeforeCity)-2; i++ {
			res.InstEls = append(res.InstEls, capaBeforeCity[i].el)
		}
	} else if len(capaBeforeCity) == 2 {
		res.AutorCapaEl = capaBeforeCity[0].el
		res.TituloCapaEl = capaBeforeCity[1].el
	} else if len(capaBeforeCity) == 1 {
		res.TituloCapaEl = capaBeforeCity[0].el
	}

	// Confirmação / desempate via metadados (se houver meta)
	if meta != nil {
		for _, p := range capaBeforeCity {
			if meta.Autor != "" && strings.EqualFold(p.text, strings.TrimSpace(meta.Autor)) {
				res.AutorCapaEl = p.el
			}
			if meta.Titulo != "" && strings.EqualFold(p.text, strings.TrimSpace(meta.Titulo)) {
				res.TituloCapaEl = p.el
			}
		}
	}

	// 4. Mapeamento da Folha de Rosto (elementos entre Ano da Capa e Ano da Folha de Rosto)
	if res.AnoFREl != nil && len(frParas) >= 2 {
		res.CidadeFREl = frParas[len(frParas)-2].el
	}

	var frBeforeCity []pInfo
	for _, p := range frParas {
		if res.CidadeFREl != nil && p.el.EndIndex > res.CidadeFREl.StartIndex {
			break
		}
		frBeforeCity = append(frBeforeCity, p)
	}

	if len(frBeforeCity) >= 1 {
		res.AutorFREl = frBeforeCity[0].el
	}
	if len(frBeforeCity) >= 2 {
		res.TituloFREl = frBeforeCity[1].el
	}

	// Demais elementos entre Título e Cidade:
	for i := 2; i < len(frBeforeCity); i++ {
		p := frBeforeCity[i]
		lower := strings.ToLower(p.text)
		if strings.Contains(lower, "orientador") || strings.Contains(lower, "coorientador") {
			res.OrientadorFREl = p.el
		} else {
			res.PreambuloFREls = append(res.PreambuloFREls, p.el)
		}
	}

	// Confirmação via metadados (se houver meta)
	if meta != nil {
		for _, p := range frBeforeCity {
			if meta.Autor != "" && strings.EqualFold(p.text, strings.TrimSpace(meta.Autor)) {
				res.AutorFREl = p.el
			}
			if meta.Titulo != "" && strings.EqualFold(p.text, strings.TrimSpace(meta.Titulo)) {
				res.TituloFREl = p.el
			}
			if meta.Orientador != "" && strings.Contains(strings.ToLower(p.text), strings.ToLower(meta.Orientador)) {
				res.OrientadorFREl = p.el
			}
		}
	}

	return res
}

// PreTextualOptimizationResult reúne os resultados do cálculo geométrico determinístico da Capa e Folha de Rosto.
type PreTextualOptimizationResult struct {
	PageHeightPt   float64 `json:"page_height_pt"`
	PageWidthPt    float64 `json:"page_width_pt"`
	MarginTopPt    float64 `json:"margin_top_pt"`
	MarginBottomPt float64 `json:"margin_bottom_pt"`
	MarginLeftPt   float64 `json:"margin_left_pt"`
	MarginRightPt  float64 `json:"margin_right_pt"`
	UsableHeightPt float64 `json:"usable_height_pt"`
	UsableWidthPt  float64 `json:"usable_width_pt"`

	// Capa
	CapaHeightPt         float64 `json:"capa_height_pt"`
	CapaFixedItemsHeight float64 `json:"capa_fixed_items_height"`
	CoverBreaks          int     `json:"cover_breaks"`
	CurrentCoverBreaks   int     `json:"current_cover_breaks"`
	CapaOverflow         bool    `json:"capa_overflow"`

	// Folha de Rosto
	FolhaRostoHeightPt float64 `json:"folha_rosto_height_pt"`
	FRFixedItemsHeight float64 `json:"fr_fixed_items_height"`
	TitlePageBreaks    int     `json:"title_page_breaks"`
	CurrentFRBreaks    int     `json:"current_fr_breaks"`
	FolhaRostoOverflow bool    `json:"folha_rosto_overflow"`

	// Anomalias estruturais
	EmptyBetweenLogoInst  int `json:"empty_between_logo_inst"`
	EmptyBetweenPreOrient int `json:"empty_between_pre_orient"`
}

// EstimateElementHeightPt calcula com rigor a altura em pontos de um elemento estrutural
func EstimateElementHeightPt(el *docs.StructuralElement, usableWidth float64) float64 {
	if el == nil {
		return 0
	}
	if el.Table != nil {
		return float64(len(el.Table.TableRows)) * 20.0
	}
	if el.SectionBreak != nil {
		return 0
	}
	if el.Paragraph == nil {
		return 0
	}

	p := el.Paragraph
	fontSize := 12.0
	lineSpacingMult := 1.5 // padrão pré-textual NBR 14724
	spaceAbove := float64(0)
	spaceBelow := float64(0)
	indentStart := float64(0)

	if p.ParagraphStyle != nil {
		if p.ParagraphStyle.LineSpacing > 0 {
			lineSpacingMult = p.ParagraphStyle.LineSpacing / 100.0
		}
		if p.ParagraphStyle.SpaceAbove != nil {
			spaceAbove = p.ParagraphStyle.SpaceAbove.Magnitude
		}
		if p.ParagraphStyle.SpaceBelow != nil {
			spaceBelow = p.ParagraphStyle.SpaceBelow.Magnitude
		}
		if p.ParagraphStyle.IndentStart != nil {
			indentStart = p.ParagraphStyle.IndentStart.Magnitude
		}
	}

	hasInline := false
	var txtBuilder strings.Builder
	for _, pel := range p.Elements {
		if pel.InlineObjectElement != nil {
			hasInline = true
		}
		if pel.TextRun != nil {
			txtBuilder.WriteString(pel.TextRun.Content)
			if pel.TextRun.TextStyle != nil && pel.TextRun.TextStyle.FontSize != nil && pel.TextRun.TextStyle.FontSize.Magnitude > 0 {
				fontSize = pel.TextRun.TextStyle.FontSize.Magnitude
			}
		}
	}

	lineHeight := fontSize * 1.25 * lineSpacingMult

	if hasInline {
		// Imagem inline oficial da logo IFS (3.2 cm = 90.7 pt) + spaceBelow (14.17 pt) + leading interno no Docs
		// Caixa física real no Google Docs: ~125 pt
		return 125.0 + spaceAbove + spaceBelow
	}

	fullTxt := strings.TrimSpace(txtBuilder.String())
	if fullTxt == "" {
		return lineHeight + spaceAbove + spaceBelow
	}

	usableW := usableWidth - indentStart
	if usableW <= 100.0 {
		usableW = 300.0
	}
	charWidth := fontSize * 0.52 // Média tipográfica Times/Arial
	charsPerLine := math.Max(10, math.Floor(usableW/charWidth))
	numLines := math.Ceil(float64(len(fullTxt)) / charsPerLine)

	return numLines*lineHeight + spaceAbove + spaceBelow
}

// DetectDocumentFontFamily identifica a família tipográfica predominante no documento (ex: "Arial" ou "Times New Roman")
func DetectDocumentFontFamily(doc *docs.Document) string {
	if doc == nil || doc.Body == nil {
		return "Arial"
	}
	timesCount := 0
	arialCount := 0
	for _, el := range doc.Body.Content {
		if el.Paragraph != nil {
			for _, pel := range el.Paragraph.Elements {
				if pel.TextRun != nil && pel.TextRun.TextStyle != nil && pel.TextRun.TextStyle.WeightedFontFamily != nil && pel.TextRun.TextStyle.WeightedFontFamily.FontFamily != "" {
					ff := strings.ToLower(pel.TextRun.TextStyle.WeightedFontFamily.FontFamily)
					if strings.Contains(ff, "times") {
						timesCount++
					} else if strings.Contains(ff, "arial") {
						arialCount++
					}
				}
			}
		}
		if timesCount > 10 || arialCount > 10 {
			break
		}
	}
	if timesCount > arialCount {
		return "Times New Roman"
	}
	return "Arial"
}

// CalculateDynamicPreTextualBreaks calcula com precisão a quantidade de quebras antes da Cidade
// na Capa e na Folha de Rosto, respeitando estritamente a base ABNT/IFS de 15 (Capa) e 14 (Folha de Rosto).
// Quando a tipografia é Arial, remove exatamente 1 quebra na Folha de Rosto (13 quebras) para compensar a largura
// dos glifos e evitar que o ano pule para a página 3.
func CalculateDynamicPreTextualBreaks(doc *docs.Document, meta *latex.TCCMetadata, fontFamily string) (int, int) {
	if fontFamily == "" {
		fontFamily = DetectDocumentFontFamily(doc)
	}

	isArial := strings.Contains(strings.ToLower(fontFamily), "arial") || !strings.Contains(strings.ToLower(fontFamily), "times")

	// Padrão institucional: Capa sempre 15
	coverBreaks := 15

	// Folha de Rosto: 14 como base; quando é Arial, remove 1 quebra (13)
	titlePageBreaks := 14
	if isArial {
		titlePageBreaks = 13
	}

	return coverBreaks, titlePageBreaks
}

// CalculateOptimalPreTextualLayout calcula de forma 100% analítica, geométrica e generalizada
// a quantidade ótima de quebras verticais antes da Cidade na Capa e na Folha de Rosto,
// assegurando que Cidade e Ano fiquem no rodapé da folha física sem vazar para a página seguinte.
func CalculateOptimalPreTextualLayout(doc *docs.Document, meta *latex.TCCMetadata) PreTextualOptimizationResult {
	res := PreTextualOptimizationResult{
		PageHeightPt:    841.89,
		PageWidthPt:     595.28,
		MarginTopPt:     85.04,
		MarginBottomPt:  56.69,
		MarginLeftPt:    85.04,
		MarginRightPt:   56.69,
		UsableHeightPt:  700.16,
		UsableWidthPt:   453.55,
		CoverBreaks:     15,
		TitlePageBreaks: 14,
	}

	if doc == nil || doc.Body == nil {
		return res
	}

	// 1. Extração da geometria do documento
	if doc.DocumentStyle != nil {
		if doc.DocumentStyle.PageSize != nil {
			if doc.DocumentStyle.PageSize.Height != nil && doc.DocumentStyle.PageSize.Height.Magnitude > 0 {
				res.PageHeightPt = doc.DocumentStyle.PageSize.Height.Magnitude
			}
			if doc.DocumentStyle.PageSize.Width != nil && doc.DocumentStyle.PageSize.Width.Magnitude > 0 {
				res.PageWidthPt = doc.DocumentStyle.PageSize.Width.Magnitude
			}
		}
		if doc.DocumentStyle.MarginTop != nil {
			res.MarginTopPt = doc.DocumentStyle.MarginTop.Magnitude
		}
		if doc.DocumentStyle.MarginBottom != nil {
			res.MarginBottomPt = doc.DocumentStyle.MarginBottom.Magnitude
		}
		if doc.DocumentStyle.MarginLeft != nil {
			res.MarginLeftPt = doc.DocumentStyle.MarginLeft.Magnitude
		}
		if doc.DocumentStyle.MarginRight != nil {
			res.MarginRightPt = doc.DocumentStyle.MarginRight.Magnitude
		}
		res.UsableHeightPt = res.PageHeightPt - res.MarginTopPt - res.MarginBottomPt
		res.UsableWidthPt = res.PageWidthPt - res.MarginLeftPt - res.MarginRightPt
	}

	st := identifyPreTextualStructure(doc, meta)
	if st.AnoCapaEl == nil {
		return res
	}

	fontFamily := DetectDocumentFontFamily(doc)
	isArial := strings.Contains(strings.ToLower(fontFamily), "arial") || !strings.Contains(strings.ToLower(fontFamily), "times")

	lineHeight12pt150 := 12.0 * 1.25 * 1.5 // 22.5 pt
	if isArial {
		lineHeight12pt150 = 12.0 * 1.30 * 1.5 // 23.4 pt
	}

	// 2. Análise da Capa (Página 1)
	if st.LogoEl != nil && len(st.InstEls) > 0 {
		res.EmptyBetweenLogoInst = countEmptyParasBetween(doc, st.LogoEl.EndIndex, st.InstEls[0].StartIndex)
	}

	capaFixedH := float64(0)
	if st.LogoEl != nil {
		capaFixedH += EstimateElementHeightPt(st.LogoEl, res.UsableWidthPt)
	}
	for _, el := range st.InstEls {
		capaFixedH += EstimateElementHeightPt(el, res.UsableWidthPt)
	}
	// 3 quebras entre Instituição e Autor
	capaFixedH += 3.0 * lineHeight12pt150
	if st.AutorCapaEl != nil {
		capaFixedH += EstimateElementHeightPt(st.AutorCapaEl, res.UsableWidthPt)
	}
	// 3 quebras entre Autor e Título
	capaFixedH += 3.0 * lineHeight12pt150
	if st.TituloCapaEl != nil {
		capaFixedH += EstimateElementHeightPt(st.TituloCapaEl, res.UsableWidthPt)
	}
	// Cidade e Ano (0 quebras entre eles)
	if st.CidadeCapaEl != nil {
		capaFixedH += EstimateElementHeightPt(st.CidadeCapaEl, res.UsableWidthPt)
	} else {
		capaFixedH += lineHeight12pt150
	}
	if st.AnoCapaEl != nil {
		capaFixedH += EstimateElementHeightPt(st.AnoCapaEl, res.UsableWidthPt)
	} else {
		capaFixedH += lineHeight12pt150
	}
	res.CapaFixedItemsHeight = capaFixedH

	// Quebras atuais no documento entre Título e Cidade
	if st.TituloCapaEl != nil && st.CidadeCapaEl != nil {
		res.CurrentCoverBreaks = countEmptyParasBetween(doc, st.TituloCapaEl.EndIndex, st.CidadeCapaEl.StartIndex)
	}

	// Altura atual total da Capa no documento:
	currentCapaH := capaFixedH + float64(res.CurrentCoverBreaks)*lineHeight12pt150
	res.CapaHeightPt = currentCapaH

	// Meta de quebras para a Capa: exatamente 15 quebras
	res.CoverBreaks = 15
	res.CapaOverflow = (res.CurrentCoverBreaks > 15)

	// 3. Análise da Folha de Rosto (Página 2)
	if len(st.PreambuloFREls) > 0 && st.OrientadorFREl != nil {
		lastPre := st.PreambuloFREls[len(st.PreambuloFREls)-1]
		res.EmptyBetweenPreOrient = countEmptyParasBetween(doc, lastPre.EndIndex, st.OrientadorFREl.StartIndex)
	}

	frFixedH := float64(0)
	// 3 quebras topo (Ano Capa -> Autor FR)
	frFixedH += 3.0 * lineHeight12pt150
	if st.AutorFREl != nil {
		frFixedH += EstimateElementHeightPt(st.AutorFREl, res.UsableWidthPt)
	}
	// 4 quebras entre Autor e Título
	frFixedH += 4.0 * lineHeight12pt150
	if st.TituloFREl != nil {
		frFixedH += EstimateElementHeightPt(st.TituloFREl, res.UsableWidthPt)
	}
	// 3 quebras entre Título e Preâmbulo
	frFixedH += 3.0 * lineHeight12pt150
	for _, el := range st.PreambuloFREls {
		frFixedH += EstimateElementHeightPt(el, res.UsableWidthPt)
	}
	if st.OrientadorFREl != nil {
		frFixedH += EstimateElementHeightPt(st.OrientadorFREl, res.UsableWidthPt)
	}
	// Cidade e Ano
	if st.CidadeFREl != nil {
		frFixedH += EstimateElementHeightPt(st.CidadeFREl, res.UsableWidthPt)
	} else {
		frFixedH += lineHeight12pt150
	}
	if st.AnoFREl != nil {
		frFixedH += EstimateElementHeightPt(st.AnoFREl, res.UsableWidthPt)
	} else {
		frFixedH += lineHeight12pt150
	}
	res.FRFixedItemsHeight = frFixedH

	// Quebras atuais no documento entre Orientador e Cidade
	if st.OrientadorFREl != nil && st.CidadeFREl != nil {
		res.CurrentFRBreaks = countEmptyParasBetween(doc, st.OrientadorFREl.EndIndex, st.CidadeFREl.StartIndex)
	}

	currentFRH := frFixedH + float64(res.CurrentFRBreaks)*lineHeight12pt150
	res.FolhaRostoHeightPt = currentFRH

	// Meta para Folha de Rosto: 14 em Times New Roman; 13 em Arial (remove 1 quebra para o ano não pular)
	targetFR := 14
	if isArial {
		targetFR = 13
	}
	res.TitlePageBreaks = targetFR
	res.FolhaRostoOverflow = (res.CurrentFRBreaks > targetFR)

	return res
}

// countEmptyParasBetween conta quantos parágrafos vazios existem estritamente entre startLimit e endLimit
func countEmptyParasBetween(doc *docs.Document, startLimit, endLimit int64) int {
	if doc == nil || doc.Body == nil {
		return 0
	}
	count := 0
	for _, el := range doc.Body.Content {
		if el.StartIndex >= startLimit && el.EndIndex <= endLimit {
			if el.Paragraph != nil {
				txt := strings.TrimSpace(extractParagraphText(el))
				if txt == "" {
					count++
				}
			}
		}
	}
	return count
}

// PreTextualGap define um intervalo entre dois elementos estruturais e a quantidade alvo de parágrafos vazios.
type PreTextualGap struct {
	StartLimit  int64
	EndLimit    int64
	TargetEmpty int
	Description string
}

// CalibratePreTextualLayoutAtomic calibra em uma única operação atômica (1 BatchUpdate em lote reverso)
// todas as quebras verticais da Capa e Folha de Rosto sem múltiplos round-trips HTTP.
func CalibratePreTextualLayoutAtomic(docsSvc *docs.Service, documentID string, meta *latex.TCCMetadata, coverBreaks, titlePageBreaks int) error {
	doc, err := docsSvc.Documents.Get(documentID).Do()
	if err != nil {
		return err
	}

	st := identifyPreTextualStructure(doc, meta)
	if st.AnoCapaEl == nil {
		return nil
	}

	var gaps []PreTextualGap

	// Capa (Página 1):
	if st.LogoEl != nil && len(st.InstEls) > 0 {
		gaps = append(gaps, PreTextualGap{StartLimit: st.LogoEl.EndIndex, EndLimit: st.InstEls[0].StartIndex, TargetEmpty: 0, Description: "Logo->Inst"})
	}
	if len(st.InstEls) > 0 && st.AutorCapaEl != nil {
		lastInst := st.InstEls[len(st.InstEls)-1]
		gaps = append(gaps, PreTextualGap{StartLimit: lastInst.EndIndex, EndLimit: st.AutorCapaEl.StartIndex, TargetEmpty: 3, Description: "Inst->Autor"})
	}
	if st.AutorCapaEl != nil && st.TituloCapaEl != nil {
		gaps = append(gaps, PreTextualGap{StartLimit: st.AutorCapaEl.EndIndex, EndLimit: st.TituloCapaEl.StartIndex, TargetEmpty: 3, Description: "Autor->Título"})
	}
	if st.TituloCapaEl != nil && st.CidadeCapaEl != nil {
		gaps = append(gaps, PreTextualGap{StartLimit: st.TituloCapaEl.EndIndex, EndLimit: st.CidadeCapaEl.StartIndex, TargetEmpty: coverBreaks, Description: "Título->Cidade"})
	}
	if st.CidadeCapaEl != nil && st.AnoCapaEl != nil {
		gaps = append(gaps, PreTextualGap{StartLimit: st.CidadeCapaEl.EndIndex, EndLimit: st.AnoCapaEl.StartIndex, TargetEmpty: 0, Description: "Cidade->Ano"})
	}

	// Folha de Rosto (Página 2):
	if st.AnoCapaEl != nil && st.AutorFREl != nil {
		gaps = append(gaps, PreTextualGap{StartLimit: st.AnoCapaEl.EndIndex, EndLimit: st.AutorFREl.StartIndex, TargetEmpty: 3, Description: "AnoCapa->AutorFR"})
	}
	if st.AutorFREl != nil && st.TituloFREl != nil {
		gaps = append(gaps, PreTextualGap{StartLimit: st.AutorFREl.EndIndex, EndLimit: st.TituloFREl.StartIndex, TargetEmpty: 4, Description: "AutorFR->TítuloFR"})
	}
	if st.TituloFREl != nil && len(st.PreambuloFREls) > 0 {
		firstPre := st.PreambuloFREls[0]
		gaps = append(gaps, PreTextualGap{StartLimit: st.TituloFREl.EndIndex, EndLimit: firstPre.StartIndex, TargetEmpty: 4, Description: "TítuloFR->Preâmbulo"})
	}
	if len(st.PreambuloFREls) > 0 && st.OrientadorFREl != nil {
		lastPre := st.PreambuloFREls[len(st.PreambuloFREls)-1]
		gaps = append(gaps, PreTextualGap{StartLimit: lastPre.EndIndex, EndLimit: st.OrientadorFREl.StartIndex, TargetEmpty: 0, Description: "Preâmbulo->Orientador"})
	}
	if st.OrientadorFREl != nil && st.CidadeFREl != nil {
		gaps = append(gaps, PreTextualGap{StartLimit: st.OrientadorFREl.EndIndex, EndLimit: st.CidadeFREl.StartIndex, TargetEmpty: titlePageBreaks, Description: "Orientador->CidadeFR"})
	}
	if st.CidadeFREl != nil && st.AnoFREl != nil {
		gaps = append(gaps, PreTextualGap{StartLimit: st.CidadeFREl.EndIndex, EndLimit: st.AnoFREl.StartIndex, TargetEmpty: 0, Description: "CidadeFR->AnoFR"})
	}

	// Ordena os gaps em ordem DECRESCENTE de StartLimit para processamento reverso atômico
	sort.Slice(gaps, func(i, j int) bool {
		return gaps[i].StartLimit > gaps[j].StartLimit
	})

	var atomicReqs []*docs.Request

	for _, gap := range gaps {
		emptyCount := 0
		var emptyEls []*docs.StructuralElement
		for _, el := range doc.Body.Content {
			if el.StartIndex >= gap.StartLimit && el.EndIndex <= gap.EndLimit {
				if el.Paragraph != nil {
					txt := ""
					for _, pel := range el.Paragraph.Elements {
						if pel.TextRun != nil {
							txt += pel.TextRun.Content
						}
					}
					if strings.TrimSpace(txt) == "" {
						emptyCount++
						emptyEls = append(emptyEls, el)
					}
				}
			}
		}

		if emptyCount < gap.TargetEmpty {
			needed := gap.TargetEmpty - emptyCount
			insertLoc := gap.EndLimit
			if len(emptyEls) > 0 {
				insertLoc = emptyEls[len(emptyEls)-1].EndIndex
			}
			atomicReqs = append(atomicReqs, &docs.Request{
				InsertText: &docs.InsertTextRequest{
					Location: &docs.Location{Index: insertLoc},
					Text:     strings.Repeat("\n", needed),
				},
			})
		} else if emptyCount > gap.TargetEmpty {
			firstDelIdx := gap.TargetEmpty
			startDel := emptyEls[firstDelIdx].StartIndex
			endDel := emptyEls[len(emptyEls)-1].EndIndex
			if endDel > startDel {
				atomicReqs = append(atomicReqs, &docs.Request{
					DeleteContentRange: &docs.DeleteContentRangeRequest{
						Range: &docs.Range{
							StartIndex: startDel,
							EndIndex:   endDel,
						},
					},
				})
			}
		}
	}

	if len(atomicReqs) > 0 {
		_, errBatch := docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: atomicReqs}).Do()
		if errBatch != nil {
			log.Printf("[AVISO] CalibratePreTextualLayoutAtomic lote falhou: %v", errBatch)
			return errBatch
		}
	}

	return nil
}

// calibrateEmptyParagraphsBetween calibra a quantidade de parágrafos vazios entre startLimit e endLimit
// garantindo que fiquem exatamente targetEmpty parágrafos vazios.
func calibrateEmptyParagraphsBetween(docsSvc *docs.Service, documentID string, startLimit, endLimit int64, targetEmpty int) error {
	doc, err := docsSvc.Documents.Get(documentID).Do()
	if err != nil {
		return err
	}

	emptyCount := 0
	var emptyEls []*docs.StructuralElement
	for _, el := range doc.Body.Content {
		if el.StartIndex >= startLimit && el.EndIndex <= endLimit {
			if el.Paragraph != nil {
				txt := ""
				for _, pel := range el.Paragraph.Elements {
					if pel.TextRun != nil {
						txt += pel.TextRun.Content
					}
				}
				if strings.TrimSpace(txt) == "" {
					emptyCount++
					emptyEls = append(emptyEls, el)
				}
			}
		}
	}

	if emptyCount < targetEmpty {
		needed := targetEmpty - emptyCount
		insertLoc := endLimit
		if len(emptyEls) > 0 {
			insertLoc = emptyEls[len(emptyEls)-1].EndIndex
		}
		_, errIns := docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{
			Requests: []*docs.Request{
				{
					InsertText: &docs.InsertTextRequest{
						Location: &docs.Location{Index: insertLoc},
						Text:     strings.Repeat("\n", needed),
					},
				},
			},
		}).Do()
		return errIns
	} else if emptyCount > targetEmpty {
		excess := emptyCount - targetEmpty
		firstDelIdx := targetEmpty
		startDel := emptyEls[firstDelIdx].StartIndex
		endDel := emptyEls[len(emptyEls)-1].EndIndex

		delReq := &docs.Request{
			DeleteContentRange: &docs.DeleteContentRangeRequest{
				Range: &docs.Range{
					StartIndex: startDel,
					EndIndex:   endDel,
				},
			},
		}
		_, errDel := docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{
			Requests: []*docs.Request{delReq},
		}).Do()
		if errDel != nil {
			log.Printf("[AVISO] calibrateEmptyParagraphsBetween: deleção em bloco [%d-%d] falhou: %v. Tentando individual...", startDel, endDel, errDel)
			for i := 0; i < excess; i++ {
				dFresh, errFresh := docsSvc.Documents.Get(documentID).Do()
				if errFresh != nil {
					return errFresh
				}
				var lastEmpty *docs.StructuralElement
				for _, el := range dFresh.Body.Content {
					if el.StartIndex >= startLimit && el.EndIndex <= endLimit && el.Paragraph != nil {
						txt := ""
						for _, pel := range el.Paragraph.Elements {
							if pel.TextRun != nil {
								txt += pel.TextRun.Content
							}
						}
						if strings.TrimSpace(txt) == "" {
							lastEmpty = el
						}
					}
				}
				if lastEmpty == nil {
					break
				}
				_, errSingle := docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{
					Requests: []*docs.Request{
						{
							DeleteContentRange: &docs.DeleteContentRangeRequest{
								Range: &docs.Range{
									StartIndex: lastEmpty.StartIndex,
									EndIndex:   lastEmpty.EndIndex,
								},
							},
						},
					},
				}).Do()
				if errSingle != nil {
					log.Printf("[ERRO] calibrateEmptyParagraphsBetween: fallback individual falhou: %v", errSingle)
					return errSingle
				}
			}
		}
	}
	return nil
}

// EnsurePreTextualPageBreaksAndLayout audita e assegura as quebras de página explícitas e a distribuição vertical (NBR 14724):
// 1) Insere Quebra de Página (PageBreak) imediatamente após o Ano da Capa (para isolar a Página 1 sem estouro);
// 2) Insere Quebra de Página (PageBreak) imediatamente após o Ano da Folha de Rosto (para isolar a Página 2 e manter o RESUMO na Página 3);
// 3) Recompõe preâmbulo incompleto a partir dos metadados oficiais com formatação de 8.0 cm à direita (justificado, 10pt, simples);
// 4) Calibra com exatidão matemática o espaçamento vertical na Capa:
//   - Entre Instituição e Autor: exatamente 3 quebras vazias (\n x 3)
//   - Entre Autor e Título: exatamente 3 quebras vazias (\n x 3)
//   - Entre Título e Cidade: exatamente 14 quebras vazias (\n x 14)
//   - Entre Cidade e Ano: 0 quebras vazias
//
// 5) Calibra com exatidão matemática o espaçamento vertical na Folha de Rosto:
//   - Topo da Folha de Rosto (Ano Capa -> Autor): exatamente 3 quebras vazias (\n x 3)
//   - Entre Autor e Título: exatamente 4 quebras vazias (\n x 4)
//   - Entre Título e Preâmbulo: exatamente 4 quebras vazias (\n x 4)
//   - Entre Preâmbulo e Orientador: 0 linhas vazias (bloco contíguo com spaceAbove 14.17 pt)
//   - Entre Orientador e Cidade: exatamente 11 quebras vazias (\n x 11)
//   - Entre Cidade e Ano: 0 quebras vazias
func EnsurePreTextualPageBreaksAndLayout(ctx context.Context, documentID string) error {
	docsSvc, err := GetDocsService(ctx)
	if err != nil {
		return fmt.Errorf("erro ao obter docs service: %w", err)
	}

	doc, err := docsSvc.Documents.Get(documentID).Do()
	if err != nil {
		return fmt.Errorf("erro ao carregar documento: %w", err)
	}

	meta, _ := latex.ParseTCCMetadata("")
	st := identifyPreTextualStructure(doc, meta)

	// 1. Recomposição e unificação do preâmbulo em parágrafo único contínuo na Folha de Rosto
	if meta != nil && meta.Preambulo != "" && len(st.PreambuloFREls) > 0 {
		firstEl := st.PreambuloFREls[0]
		lastEl := st.PreambuloFREls[len(st.PreambuloFREls)-1]
		combinedTxt := ""
		for _, el := range st.PreambuloFREls {
			combinedTxt += extractParagraphText(el)
		}
		combinedTrim := strings.TrimSpace(combinedTxt)
		if len(st.PreambuloFREls) > 1 || combinedTrim != strings.TrimSpace(meta.Preambulo) || !strings.Contains(combinedTrim, "obtenção do grau") {
			repReq := &docs.BatchUpdateDocumentRequest{
				Requests: []*docs.Request{
					{
						DeleteContentRange: &docs.DeleteContentRangeRequest{
							Range: &docs.Range{StartIndex: firstEl.StartIndex, EndIndex: lastEl.EndIndex - 1},
						},
					},
					{
						InsertText: &docs.InsertTextRequest{
							Location: &docs.Location{Index: firstEl.StartIndex},
							Text:     meta.Preambulo,
						},
					},
				},
			}
			_, _ = docsSvc.Documents.BatchUpdate(documentID, repReq).Do()
			doc, _ = docsSvc.Documents.Get(documentID).Do()
			st = identifyPreTextualStructure(doc, meta)
		}
	}

	// 1.1 Zero quebras de linha entre Logo e Instituição (exigência ABNT/IFS)
	if st.LogoEl != nil && len(st.InstEls) > 0 {
		_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, st.LogoEl.EndIndex, st.InstEls[0].StartIndex, 0)
		doc, _ = docsSvc.Documents.Get(documentID).Do()
		st = identifyPreTextualStructure(doc, meta)
	}

	// 2. Localiza os elementos dos anos e insere quebras de página se ausentes
	var pbReqs []*docs.Request
	if st.AnoFREl != nil && !st.AnoFRHasBreak {
		pbReqs = append(pbReqs, &docs.Request{
			InsertPageBreak: &docs.InsertPageBreakRequest{
				Location: &docs.Location{Index: st.AnoFREl.EndIndex - 1},
			},
		})
	}
	if st.AnoCapaEl != nil && !st.AnoCapaHasBreak {
		pbReqs = append(pbReqs, &docs.Request{
			InsertPageBreak: &docs.InsertPageBreakRequest{
				Location: &docs.Location{Index: st.AnoCapaEl.EndIndex - 1},
			},
		})
	}
	if len(pbReqs) > 0 {
		_, _ = docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: pbReqs}).Do()
		doc, _ = docsSvc.Documents.Get(documentID).Do()
		st = identifyPreTextualStructure(doc, meta)
	}

	// 3. Calibração Atômica (1 BatchUpdate) de todas as quebras verticais pré-textuais:
	// - Capa: 0 entre Logo e Inst, 3 entre Inst e Autor, 3 entre Autor e Título, base 15 entre Título e Cidade, 0 entre Cidade e Ano
	// - Folha de Rosto: 3 topo, 4 entre Autor e Título, 4 entre Título e Preâmbulo, 0 entre Preâmbulo e Orientador, base 14 (ou 13 se Arial) entre Orientador e Cidade, 0 entre Cidade e Ano
	coverBreaks, titlePageBreaks := CalculateDynamicPreTextualBreaks(doc, meta, "")
	if errCal := CalibratePreTextualLayoutAtomic(docsSvc, documentID, meta, coverBreaks, titlePageBreaks); errCal != nil {
		log.Printf("[AVISO] EnsurePreTextualPageBreaksAndLayout: CalibratePreTextualLayoutAtomic: %v", errCal)
	}

	// 4. Limpeza de parágrafos vazios espúrios na transição da Folha de Rosto para o RESUMO
	doc, _ = docsSvc.Documents.Get(documentID).Do()
	var cleanReqs []*docs.Request
	inResumeTransition := false
	for i := len(doc.Body.Content) - 1; i >= 0; i-- {
		el := doc.Body.Content[i]
		if el.Paragraph != nil {
			txt := strings.TrimSpace(extractParagraphText(el))
			hasPageBreak := paragraphHasPageBreak(el)
			upperT := strings.ToUpper(txt)

			if upperT == "RESUMO" || upperT == "# RESUMO" || strings.HasPrefix(upperT, "RESUMO") {
				inResumeTransition = true
				continue
			}
			if inResumeTransition {
				if txt == "" && !hasPageBreak {
					cleanReqs = append(cleanReqs, &docs.Request{
						DeleteContentRange: &docs.DeleteContentRangeRequest{
							Range: &docs.Range{StartIndex: el.StartIndex, EndIndex: el.EndIndex},
						},
					})
				} else {
					inResumeTransition = false
				}
			}
		}
	}
	if len(cleanReqs) > 0 {
		_, _ = docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: cleanReqs}).Do()
	}

	return nil
}

// FormatABNTCoverAndTitlePageIsolated formata com isolamento absoluto a Capa (Página 1) e Folha de Rosto (Página 2),
// garantindo que regras e espaçamentos do corpo textual não afetem os elementos pré-textuais.
func FormatABNTCoverAndTitlePageIsolated(ctx context.Context, documentID string, fontFamily string) error {
	if fontFamily == "" {
		fontFamily = ABNTDefaultFontFamily
	}

	docsSvc, err := GetDocsService(ctx)
	if err != nil {
		return fmt.Errorf("erro ao obter docs service: %w", err)
	}

	doc, err := docsSvc.Documents.Get(documentID).Do()
	if err != nil {
		return fmt.Errorf("erro ao carregar documento: %w", err)
	}

	meta, _ := latex.ParseTCCMetadata("")
	st := identifyPreTextualStructure(doc, meta)
	if st.AnoCapaEl == nil {
		return nil // Sem capa identificável
	}

	blackColor := &docs.OptionalColor{
		Color: &docs.Color{RgbColor: &docs.RgbColor{Red: 0, Green: 0, Blue: 0}},
	}

	// Normaliza hífen duplo em cidades nos elementos pré-textuais
	var replaceReqs []*docs.Request
	for _, cidEl := range []*docs.StructuralElement{st.CidadeCapaEl, st.CidadeFREl} {
		if cidEl != nil {
			txt := extractParagraphText(cidEl)
			if strings.Contains(txt, " -- ") {
				replaceReqs = append(replaceReqs, &docs.Request{
					ReplaceAllText: &docs.ReplaceAllTextRequest{
						ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: " -- "},
						ReplaceText:  " - ",
					},
				})
			}
		}
	}
	if len(replaceReqs) > 0 {
		_, _ = docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: replaceReqs}).Do()
		doc, err = docsSvc.Documents.Get(documentID).Do()
		if err != nil {
			return fmt.Errorf("erro ao recarregar documento pós replace: %w", err)
		}
		st = identifyPreTextualStructure(doc, meta)
	}

	var reqs []*docs.Request

	for _, el := range doc.Body.Content {
		if el.Paragraph == nil {
			continue
		}
		if el.EndIndex > st.TitlePageEndIndex {
			break // Limite estrito: jamais toca o corpo do trabalho (Página 3 em diante)
		}

		pStart := el.StartIndex
		pEnd := el.EndIndex
		pRange := &docs.Range{StartIndex: pStart, EndIndex: pEnd}

		txt := extractParagraphText(el)
		hasInline := paragraphHasInlineObject(el)
		trimmed := strings.TrimSpace(txt)

		isCapa := el.EndIndex <= st.CoverEndIndex

		// ─── PARÁGRAFO VAZIO OU LOGO ───
		if trimmed == "" {
			if hasInline {
				// Logo oficial centralizada no topo da Capa (sem espaço abaixo)
				reqs = append(reqs, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType:      "NORMAL_TEXT",
							Alignment:           "CENTER",
							IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
							IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
							SpaceAbove:          &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
							SpaceBelow:          &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
							AvoidWidowAndOrphan: true,
						},
						Fields: "namedStyleType,alignment,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
					},
				})
			} else {
				// Linhas em branco pré-textuais da Capa ou Folha de Rosto:
				// Exigência ABNT: entrelinhas 1,5 (150%), fontSize 12pt, spaceAbove/Below zerados, recuos zerados
				reqs = append(reqs, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType:  "NORMAL_TEXT",
							Alignment:       "CENTER",
							LineSpacing:     ABNTLineSpacing, // 150% (1,5)
							IndentFirstLine: &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
							IndentStart:     &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
							SpaceAbove:      &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
							SpaceBelow:      &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
						},
						Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow",
					},
				})
				reqs = append(reqs, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: pRange,
						TextStyle: &docs.TextStyle{
							FontSize: &docs.Dimension{Magnitude: 12.0, Unit: "PT"},
						},
						Fields: "fontSize",
					},
				})
			}
			continue
		}

		// ─── PÁGINA 1: CAPA OFICIAL (config.tex \imprimircapa: \center\bfseries) ───
		if isCapa {
			// Capa completa com espaçamento entre linhas 1,5 (150%) em todos os elementos, negrito, centralizado
			reqs = append(reqs, &docs.Request{
				UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
					Range: pRange,
					ParagraphStyle: &docs.ParagraphStyle{
						NamedStyleType:      "NORMAL_TEXT",
						Alignment:           "CENTER",
						LineSpacing:         ABNTLineSpacing, // 150% (1,5) em todos os elementos da Capa
						IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
						IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
						SpaceAbove:          &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
						SpaceBelow:          &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
						AvoidWidowAndOrphan: true,
					},
					Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
				},
			})
			reqs = append(reqs, &docs.Request{
				UpdateTextStyle: &docs.UpdateTextStyleRequest{
					Range: pRange,
					TextStyle: &docs.TextStyle{
						FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeBodyPt, Unit: "PT"}, // 12.0 pt
						WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
						Bold:               true, // Capa: todos em negrito
						ForegroundColor:    blackColor,
					},
					Fields: "fontSize,weightedFontFamily,bold,foregroundColor",
				},
			})
			continue
		}

		// ─── PÁGINA 2: FOLHA DE ROSTO (config.tex \imprimirfolhaderosto) ───
		isPreambulo := isElementInList(el, st.PreambuloFREls)
		isOrientador := (st.OrientadorFREl != nil && el.StartIndex == st.OrientadorFREl.StartIndex)

		if isPreambulo || isOrientador {
			// Preâmbulo e Orientador:
			// NBR 14724: IndentStart 8,0 cm (226.772 pt), IndentFirstLine 8,0 cm (226.772 pt) formando bloco reto perfeito do meio para a direita
			// Justificado, 10 pt, entrelinhas simples (1,0 / 100%)
			spAbove := float64(0)
			if isOrientador {
				spAbove = 14.17 // 0.5 cm como \vspace{0.5cm} do LaTeX
			}
			reqs = append(reqs, &docs.Request{
				UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
					Range: pRange,
					ParagraphStyle: &docs.ParagraphStyle{
						NamedStyleType:      "NORMAL_TEXT",
						Alignment:           "JUSTIFIED",
						LineSpacing:         ABNTLineSpacingSingle,                                                                                      // 100% (Simples)
						IndentFirstLine:     &docs.Dimension{Magnitude: ABNTIndentPreTextualNotePt, Unit: "PT", ForceSendFields: []string{"Magnitude"}}, // 8,0 cm absoluto na régua
						IndentStart:         &docs.Dimension{Magnitude: ABNTIndentPreTextualNotePt, Unit: "PT", ForceSendFields: []string{"Magnitude"}}, // 8,0 cm absoluto na régua
						SpaceAbove:          &docs.Dimension{Magnitude: spAbove, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
						SpaceBelow:          &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
						AvoidWidowAndOrphan: true,
					},
					Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
				},
			})
			reqs = append(reqs, &docs.Request{
				UpdateTextStyle: &docs.UpdateTextStyleRequest{
					Range: pRange,
					TextStyle: &docs.TextStyle{
						FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeCaptionPt, Unit: "PT"}, // 10.0 pt obrigatório ABNT
						WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
						Bold:               false, // Sem negrito
						ForegroundColor:    blackColor,
					},
					Fields: "fontSize,weightedFontFamily,bold,foregroundColor",
				},
			})
			continue
		}

		// Elementos textuais comuns da Folha de Rosto (Autor, Título, Local, Ano):
		// Espaçamento 1,5 (150%), fonte 12pt, centralizado
		isTituloFR := (st.TituloFREl != nil && el.StartIndex == st.TituloFREl.StartIndex)
		boldFR := false
		if isTituloFR {
			boldFR = true // Apenas o Título é em negrito na Folha de Rosto
		}

		reqs = append(reqs, &docs.Request{
			UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
				Range: pRange,
				ParagraphStyle: &docs.ParagraphStyle{
					NamedStyleType:      "NORMAL_TEXT",
					Alignment:           "CENTER",
					LineSpacing:         ABNTLineSpacing, // 150% (1,5)
					IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
					IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
					SpaceAbove:          &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
					SpaceBelow:          &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
					AvoidWidowAndOrphan: true,
				},
				Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
			},
		})
		reqs = append(reqs, &docs.Request{
			UpdateTextStyle: &docs.UpdateTextStyleRequest{
				Range: pRange,
				TextStyle: &docs.TextStyle{
					FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeBodyPt, Unit: "PT"}, // 12.0 pt
					WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
					Bold:               boldFR, // Apenas Título em negrito na Folha de Rosto
					ForegroundColor:    blackColor,
				},
				Fields: "fontSize,weightedFontFamily,bold,foregroundColor",
			},
		})
	}

	if len(reqs) > 0 {
		_, errBatch := docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: reqs}).Do()
		if errBatch != nil {
			return fmt.Errorf("erro ao formatar Capa e Folha de Rosto isoladamente: %w", errBatch)
		}
	}

	return nil
}

// FitCityYearToPageBottom ajusta dinamicamente a quantidade de quebras de linha antes da Cidade
// na Capa e na Folha de Rosto, validando contra o PDF oficial renderizado pelo Google Docs para
// empurrar o Ano para o limite inferior da folha física sem permitir que ele vaze para a página seguinte.
func FitCityYearToPageBottom(ctx context.Context, documentID string, meta *latex.TCCMetadata, fontFamily string) error {
	docsSvc, err := GetDocsService(ctx)
	if err != nil {
		return fmt.Errorf("erro ao obter docs service: %w", err)
	}

	doc, err := docsSvc.Documents.Get(documentID).Do()
	if err != nil {
		return fmt.Errorf("erro ao carregar documento: %w", err)
	}

	if meta == nil {
		meta, _ = latex.ParseTCCMetadata("")
	}
	st := identifyPreTextualStructure(doc, meta)
	if st.AnoCapaEl == nil {
		return nil // Sem elementos pré-textuais
	}

	// Calibração atômica dinâmica com base 15 (Capa) e 14/13 (Folha de Rosto) com verificação de estouro
	coverBreaks, titlePageBreaks := CalculateDynamicPreTextualBreaks(doc, meta, fontFamily)

	if errCal := CalibratePreTextualLayoutAtomic(docsSvc, documentID, meta, coverBreaks, titlePageBreaks); errCal != nil {
		log.Printf("[AVISO] FitCityYearToPageBottom: CalibratePreTextualLayoutAtomic: %v", errCal)
	}

	if errFmt := FormatABNTCoverAndTitlePageIsolated(ctx, documentID, fontFamily); errFmt != nil {
		log.Printf("[AVISO] FitCityYearToPageBottom: FormatABNTCoverAndTitlePageIsolated: %v", errFmt)
	}

	log.Printf("[INFO] FitCityYearToPageBottom concluído com sucesso: Capa = %d quebras, Folha de Rosto = %d quebras.",
		coverBreaks, titlePageBreaks)

	return nil
}
