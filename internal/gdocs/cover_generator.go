package gdocs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
	sb.WriteString("\n\n\n")

	// Autor
	sb.WriteString(strings.ToUpper(meta.Autor))
	sb.WriteString("\n\n\n")

	// Título do Trabalho
	sb.WriteString(strings.ToUpper(meta.Titulo))
	sb.WriteString(strings.Repeat("\n", 18)) // 18 quebras vazias (calibração \vfill) para empurrar Local e Ano para a margem inferior da Capa

	// Local e Ano
	sb.WriteString(strings.ToUpper(meta.Local))
	sb.WriteString("\n")
	sb.WriteString(meta.Data)

	capaText := sb.String()

	// --- PÁGINA 2: FOLHA DE ROSTO (NBR 14724 + IFS config.tex) ---
	var sbFR strings.Builder
	sbFR.WriteString("\n\n") // \vspace*{2cm} do LaTeX
	sbFR.WriteString(strings.ToUpper(meta.Autor))
	sbFR.WriteString(strings.Repeat("\n", 7)) // Exatamente 7 quebras de linhas entre Autor e Título
	sbFR.WriteString(strings.ToUpper(meta.Titulo))
	sbFR.WriteString(strings.Repeat("\n", 5)) // Exatamente 5 quebras de linhas entre Título e Preâmbulo

	// Bloco da Nota Institucional + Orientador (alinhado da metade para a margem direita, recuo 8.0 cm)
	sbFR.WriteString(meta.Preambulo)
	sbFR.WriteString("\n")
	sbFR.WriteString("Orientador: ")
	sbFR.WriteString(meta.Orientador)
	if meta.Coorientador != "" {
		sbFR.WriteString("\nCoorientador: ")
		sbFR.WriteString(meta.Coorientador)
	}
	sbFR.WriteString(strings.Repeat("\n", 14)) // 14 quebras vazias (calibração \vfill) para empurrar Local e Ano para a margem inferior da Folha de Rosto

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

// EnsurePreTextualPageBreaksAndLayout audita e assegura as quebras de página explícitas e a distribuição vertical (NBR 14724):
// 1) Insere Quebra de Página (PageBreak) imediatamente após o Ano da Capa (para isolar a Página 1);
// 2) Insere Quebra de Página (PageBreak) imediatamente após o Ano da Folha de Rosto (para isolar a Página 2 e manter o RESUMO na Página 3);
// 3) Recompõe preâmbulo incompleto a partir dos metadados oficiais com formatação de 8.0 cm à direita (justificado, 10pt, simples);
// 4) Calibra o espaçamento vertical para empurrar Local e Ano para o limite inferior da folha tanto na Capa quanto na Folha de Rosto.
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

	// 1. Recomposição de preâmbulo se estiver cortado ("Campus\n" sem o restante)
	if meta != nil && meta.Preambulo != "" {
		for _, el := range doc.Body.Content {
			if el.Paragraph != nil {
				txt := ""
				for _, pel := range el.Paragraph.Elements {
					if pel.TextRun != nil {
						txt += pel.TextRun.Content
					}
				}
				trimmed := strings.TrimSpace(txt)
				if strings.HasPrefix(trimmed, "Trabalho de Conclusão de Curso apresentado") && !strings.Contains(trimmed, "obtenção do grau") {
					repReq := &docs.BatchUpdateDocumentRequest{
						Requests: []*docs.Request{
							{
								DeleteContentRange: &docs.DeleteContentRangeRequest{
									Range: &docs.Range{StartIndex: el.StartIndex, EndIndex: el.EndIndex - 1},
								},
							},
							{
								InsertText: &docs.InsertTextRequest{
									Location: &docs.Location{Index: el.StartIndex},
									Text:     meta.Preambulo,
								},
							},
						},
					}
					_, _ = docsSvc.Documents.BatchUpdate(documentID, repReq).Do()
					doc, _ = docsSvc.Documents.Get(documentID).Do()
					break
				}
			}
		}
	}

	// 2. Localiza os elementos dos anos e checa a existência de quebras de página
	var ano1El *docs.StructuralElement
	var ano2El *docs.StructuralElement
	ano1HasBreak := false
	ano2HasBreak := false
	anoCount := 0

	for _, el := range doc.Body.Content {
		if el.Paragraph != nil {
			txt := ""
			pHasBreak := false
			for _, pel := range el.Paragraph.Elements {
				if pel.TextRun != nil {
					txt += pel.TextRun.Content
				}
				if pel.PageBreak != nil {
					pHasBreak = true
				}
			}
			t := strings.TrimSpace(txt)
			if t == "2026" || (meta != nil && t == strings.TrimSpace(meta.Data)) {
				anoCount++
				if anoCount == 1 {
					ano1El = el
					ano1HasBreak = pHasBreak
				} else if anoCount == 2 {
					ano2El = el
					ano2HasBreak = pHasBreak
					break
				}
			}
		}
	}

	// 3. Inserção das quebras de página em ordem inversa (Ano 2 primeiro, depois Ano 1)
	var pbReqs []*docs.Request
	if ano2El != nil && !ano2HasBreak {
		pbReqs = append(pbReqs, &docs.Request{
			InsertPageBreak: &docs.InsertPageBreakRequest{
				Location: &docs.Location{Index: ano2El.EndIndex - 1},
			},
		})
	}

	if ano1El != nil && !ano1HasBreak {
		pbReqs = append(pbReqs, &docs.Request{
			InsertPageBreak: &docs.InsertPageBreakRequest{
				Location: &docs.Location{Index: ano1El.EndIndex - 1},
			},
		})
	}

	if len(pbReqs) > 0 {
		_, errPb := docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: pbReqs}).Do()
		if errPb != nil {
			log.Printf("[AVISO] EnsurePreTextualPageBreaksAndLayout: erro ao inserir quebras: %v", errPb)
		} else {
			log.Printf("[INFO] Inseridas %d quebras de página obrigatórias (NBR 14724) na Capa / Folha de Rosto", len(pbReqs))
		}
		doc, _ = docsSvc.Documents.Get(documentID).Do()
	}

	// 3.5 Calibração rigorosa das quebras internas da Folha de Rosto (NBR 14724):
	// - Exatamente 7 quebras de linhas (\n x 7) entre o Nome do Autor e o Título do Trabalho.
	// - Exatamente 5 quebras de linhas (\n x 5) entre o Título do Trabalho e o Preâmbulo / Tipo de Trabalho.
	doc, _ = docsSvc.Documents.Get(documentID).Do()
	if ano1El != nil {
		ano1EndIndex := ano1El.EndIndex
		var autorFR, tituloFR, preambuloFR *docs.StructuralElement
		for _, el := range doc.Body.Content {
			if el.StartIndex < ano1EndIndex {
				continue
			}
			if el.Paragraph != nil {
				t := ""
				for _, pel := range el.Paragraph.Elements {
					if pel.TextRun != nil {
						t += pel.TextRun.Content
					}
				}
				trim := strings.TrimSpace(t)
				if trim == "" {
					continue
				}
				if autorFR == nil && (strings.Contains(strings.ToUpper(trim), "JOSÉ GUSTAVO") || (meta != nil && strings.EqualFold(trim, strings.ToUpper(meta.Autor)))) {
					autorFR = el
					continue
				}
				if autorFR != nil && tituloFR == nil && (strings.Contains(strings.ToUpper(trim), "AVALIAÇÃO DE UMA") || (meta != nil && strings.EqualFold(trim, strings.ToUpper(meta.Titulo)))) {
					tituloFR = el
					continue
				}
				if tituloFR != nil && preambuloFR == nil && (strings.HasPrefix(trim, "Trabalho de Conclusão") || isPreTextualNote(trim)) {
					preambuloFR = el
					break
				}
			}
		}

		if autorFR != nil && tituloFR != nil && preambuloFR != nil {
			// Ajusta entre Título e Preâmbulo primeiro (para preservar índices anteriores)
			emptyTituloPreambulo := 0
			var emptyTPList []*docs.StructuralElement
			for _, el := range doc.Body.Content {
				if el.StartIndex >= tituloFR.EndIndex && el.EndIndex <= preambuloFR.StartIndex {
					if el.Paragraph != nil {
						t := ""
						for _, pel := range el.Paragraph.Elements {
							if pel.TextRun != nil {
								t += pel.TextRun.Content
							}
						}
						if strings.TrimSpace(t) == "" {
							emptyTituloPreambulo++
							emptyTPList = append(emptyTPList, el)
						}
					}
				}
			}

			if emptyTituloPreambulo < 5 {
				needed := 5 - emptyTituloPreambulo
				_, _ = docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{
					Requests: []*docs.Request{
						{
							InsertText: &docs.InsertTextRequest{
								Location: &docs.Location{Index: preambuloFR.StartIndex},
								Text:     strings.Repeat("\n", needed),
							},
						},
					},
				}).Do()
				doc, _ = docsSvc.Documents.Get(documentID).Do()
			} else if emptyTituloPreambulo > 5 {
				excess := emptyTituloPreambulo - 5
				for k := len(emptyTPList) - 1; k >= len(emptyTPList)-excess; k-- {
					elToDel := emptyTPList[k]
					_, _ = docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{
						Requests: []*docs.Request{
							{
								DeleteContentRange: &docs.DeleteContentRangeRequest{
									Range: &docs.Range{StartIndex: elToDel.StartIndex, EndIndex: elToDel.EndIndex},
								},
							},
						},
					}).Do()
				}
				doc, _ = docsSvc.Documents.Get(documentID).Do()
			}

			// Recarrega referências para ajustar entre Autor e Título
			doc, _ = docsSvc.Documents.Get(documentID).Do()
			autorFR = nil
			tituloFR = nil
			for _, el := range doc.Body.Content {
				if el.StartIndex < ano1EndIndex {
					continue
				}
				if el.Paragraph != nil {
					t := ""
					for _, pel := range el.Paragraph.Elements {
						if pel.TextRun != nil {
							t += pel.TextRun.Content
						}
					}
					trim := strings.TrimSpace(t)
					if trim == "" {
						continue
					}
					if autorFR == nil && (strings.Contains(strings.ToUpper(trim), "JOSÉ GUSTAVO") || (meta != nil && strings.EqualFold(trim, strings.ToUpper(meta.Autor)))) {
						autorFR = el
						continue
					}
					if autorFR != nil && tituloFR == nil && (strings.Contains(strings.ToUpper(trim), "AVALIAÇÃO DE UMA") || (meta != nil && strings.EqualFold(trim, strings.ToUpper(meta.Titulo)))) {
						tituloFR = el
						break
					}
				}
			}

			if autorFR != nil && tituloFR != nil {
				emptyAutorTitulo := 0
				var emptyATList []*docs.StructuralElement
				for _, el := range doc.Body.Content {
					if el.StartIndex >= autorFR.EndIndex && el.EndIndex <= tituloFR.StartIndex {
						if el.Paragraph != nil {
							t := ""
							for _, pel := range el.Paragraph.Elements {
								if pel.TextRun != nil {
									t += pel.TextRun.Content
								}
							}
							if strings.TrimSpace(t) == "" {
								emptyAutorTitulo++
								emptyATList = append(emptyATList, el)
							}
						}
					}
				}

				if emptyAutorTitulo < 7 {
					needed := 7 - emptyAutorTitulo
					_, _ = docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{
						Requests: []*docs.Request{
							{
								InsertText: &docs.InsertTextRequest{
									Location: &docs.Location{Index: tituloFR.StartIndex},
									Text:     strings.Repeat("\n", needed),
								},
							},
						},
					}).Do()
					doc, _ = docsSvc.Documents.Get(documentID).Do()
				} else if emptyAutorTitulo > 7 {
					excess := emptyAutorTitulo - 7
					for k := len(emptyATList) - 1; k >= len(emptyATList)-excess; k-- {
						elToDel := emptyATList[k]
						_, _ = docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{
							Requests: []*docs.Request{
								{
									DeleteContentRange: &docs.DeleteContentRangeRequest{
										Range: &docs.Range{StartIndex: elToDel.StartIndex, EndIndex: elToDel.EndIndex},
									},
								},
							},
						}).Do()
					}
					doc, _ = docsSvc.Documents.Get(documentID).Do()
				}
			}
		}
	}

	// 4. Calibração vertical para empurrar Local e Ano para o limite inferior da folha:
	// Identifica quantos parágrafos vazios existem antes do Local ("LAGARTO -- SE") na Capa e na Folha de Rosto.
	// Capa: meta de 18 quebras vazias. Folha de Rosto: meta de 13 quebras vazias (calibrada após os 7 e 5 espaços).
	var insertVFillReqs []*docs.Request
	doc, _ = docsSvc.Documents.Get(documentID).Do()

	content := doc.Body.Content
	cityIndex := 0
	for i, el := range content {
		if el.Paragraph != nil {
			txt := ""
			for _, pel := range el.Paragraph.Elements {
				if pel.TextRun != nil {
					txt += pel.TextRun.Content
				}
			}
			t := strings.TrimSpace(txt)
			upperT := strings.ToUpper(t)

			// Verifica se é a linha da Cidade (Local), distinguindo rigorosamente de "CAMPUS LAGARTO"
			isCity := (strings.Contains(upperT, "LAGARTO -- SE") || strings.Contains(upperT, "LAGARTO - SE")) && !strings.Contains(upperT, "CAMPUS")
			if meta != nil && strings.EqualFold(t, strings.TrimSpace(meta.Local)) && !strings.Contains(upperT, "CAMPUS") {
				isCity = true
			}

			if isCity {
				cityIndex++
				emptyCount := 0
				var emptyEls []*docs.StructuralElement
				for j := i - 1; j >= 0; j-- {
					prevEl := content[j]
					if prevEl.Paragraph != nil {
						pTxt := ""
						for _, pel := range prevEl.Paragraph.Elements {
							if pel.TextRun != nil {
								pTxt += pel.TextRun.Content
							}
						}
						if strings.TrimSpace(pTxt) == "" {
							emptyCount++
							emptyEls = append(emptyEls, prevEl)
						} else {
							break
						}
					}
				}

				targetEmpty := 18
				if cityIndex >= 2 {
					targetEmpty = 13 // 13 quebras na Folha de Rosto mantêm LAGARTO/2026 no rodapé sem estourar
				}
				needed := targetEmpty - emptyCount
				if needed > 0 {
					insertVFillReqs = append(insertVFillReqs, &docs.Request{
						InsertText: &docs.InsertTextRequest{
							Location: &docs.Location{Index: el.StartIndex},
							Text:     strings.Repeat("\n", needed),
						},
					})
				} else if needed < 0 && cityIndex >= 2 {
					excess := -needed
					for k := 0; k < excess && k < len(emptyEls); k++ {
						elDel := emptyEls[k]
						_, _ = docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{
							Requests: []*docs.Request{
								{
									DeleteContentRange: &docs.DeleteContentRangeRequest{
										Range: &docs.Range{StartIndex: elDel.StartIndex, EndIndex: elDel.EndIndex},
									},
								},
							},
						}).Do()
					}
				}
			}
		}
	}

	// Insere em ordem inversa para não invalidar índices
	if len(insertVFillReqs) > 0 {
		var revReqs []*docs.Request
		for i := len(insertVFillReqs) - 1; i >= 0; i-- {
			revReqs = append(revReqs, insertVFillReqs[i])
		}
		_, _ = docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: revReqs}).Do()
		log.Printf("[INFO] Calibração vertical (\vfill) aplicada: %d expansões de base inseridas", len(revReqs))
	}

	// 5. Limpeza de parágrafos vazios indevidos na Capa e nas transições de página:
	// - Entre "INSTITUTO FEDERAL" e "CAMPUS LAGARTO" não deve haver parágrafos vazios.
	// - Após a Quebra de Página da Folha de Rosto (antes de RESUMO), remove linhas vazias redundantes.
	doc, _ = docsSvc.Documents.Get(documentID).Do()
	var cleanReqs []*docs.Request
	inInstHeader := false
	inResumeTransition := false

	for i := len(doc.Body.Content) - 1; i >= 0; i-- {
		el := doc.Body.Content[i]
		if el.Paragraph != nil {
			txt := ""
			for _, pel := range el.Paragraph.Elements {
				if pel.TextRun != nil {
					txt += pel.TextRun.Content
				}
			}
			t := strings.TrimSpace(txt)
			upperT := strings.ToUpper(t)

			// Detecta bloco institucional da Capa
			if strings.Contains(upperT, "CAMPUS LAGARTO") {
				inInstHeader = true
			}
			if inInstHeader && t == "" {
				// Parágrafo vazio espúrio dentro do cabeçalho institucional da Capa
				cleanReqs = append(cleanReqs, &docs.Request{
					DeleteContentRange: &docs.DeleteContentRangeRequest{
						Range: &docs.Range{StartIndex: el.StartIndex, EndIndex: el.EndIndex},
					},
				})
			}
			if strings.Contains(upperT, "INSTITUTO FEDERAL") {
				inInstHeader = false
			}

			// Detecta transição da Folha de Rosto para o RESUMO
			hasPageBreak := false
			for _, pel := range el.Paragraph.Elements {
				if pel.PageBreak != nil {
					hasPageBreak = true
					break
				}
			}

			if upperT == "RESUMO" || upperT == "# RESUMO" || strings.HasPrefix(upperT, "RESUMO") {
				inResumeTransition = true
				continue
			}
			if inResumeTransition {
				if t == "" && !hasPageBreak {
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
		log.Printf("[INFO] Saneamento de quebras redundantes pré-textuais aplicado (%d removidas)", len(cleanReqs))
	}

	return nil
}
