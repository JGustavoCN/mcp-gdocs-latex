package gdocs

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png" // Registro obrigatório do codec PNG para image.DecodeConfig
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
)

// MermaidInsertionResult encapsula o retorno da inserção do diagrama Mermaid
type MermaidInsertionResult struct {
	FileID   string  `json:"file_id,omitempty"`
	WidthPt  float64 `json:"width_pt"`
	HeightPt float64 `json:"height_pt"`
	IsInline bool    `json:"is_inline"`
	LinkURL  string  `json:"link_url,omitempty"`
	Message  string  `json:"message"`
}

// BuildMermaidPakoURL converte código Mermaid em uma URL compactada Pako compatível com mermaid.ink (< 2KB)
func BuildMermaidPakoURL(mermaidCode string) (string, error) {
	payload := map[string]interface{}{
		"code": strings.TrimSpace(mermaidCode),
		"mermaid": map[string]string{
			"theme": "default",
		},
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	var zlibBuf bytes.Buffer
	zw := zlib.NewWriter(&zlibBuf)
	if _, err := zw.Write(jsonBytes); err != nil {
		return "", err
	}
	if err := zw.Close(); err != nil {
		return "", err
	}
	encodedPako := base64.RawURLEncoding.EncodeToString(zlibBuf.Bytes())
	return fmt.Sprintf("https://mermaid.ink/img/pako:%s", encodedPako), nil
}

// InsertMermaidDiagram carrega diagrama Mermaid (código, arquivo .mmd ou imagem PNG) e o insere inline no Google Docs
func InsertMermaidDiagram(ctx context.Context, documentID string, targetIndex int64, pngOrMmdPath, caption, source string) (*MermaidInsertionResult, error) {
	trimmed := strings.TrimSpace(pngOrMmdPath)

	// 1. Verifica se é código Mermaid ou arquivo .mmd
	var mmdCode string
	if strings.HasSuffix(strings.ToLower(trimmed), ".mmd") {
		content, err := os.ReadFile(trimmed)
		if err == nil {
			mmdCode = string(content)
		}
	} else if strings.HasPrefix(trimmed, "flowchart") || strings.HasPrefix(trimmed, "graph") ||
		strings.HasPrefix(trimmed, "sequenceDiagram") || strings.HasPrefix(trimmed, "classDiagram") ||
		strings.HasPrefix(trimmed, "stateDiagram") || strings.HasPrefix(trimmed, "erDiagram") {
		mmdCode = trimmed
	} else if strings.HasSuffix(strings.ToLower(trimmed), ".png") {
		// Tenta encontrar o arquivo .mmd correspondente (ex: arquitetura-1.png -> arquitetura.mmd)
		candidates := []string{
			strings.TrimSuffix(trimmed, ".png") + ".mmd",
			strings.TrimSuffix(trimmed, "-1.png") + ".mmd",
		}
		for _, c := range candidates {
			if content, err := os.ReadFile(c); err == nil {
				mmdCode = string(content)
				log.Printf("[INFO] Localizado código fonte Mermaid correspondente: %s", c)
				break
			}
		}
	}

	// Se temos código Mermaid, usa a renderização direta via Pako (Zero Cota do Drive!)
	if mmdCode != "" {
		pakoURL, err := BuildMermaidPakoURL(mmdCode)
		if err == nil && len(pakoURL) <= 2040 {
			targetWidthPt := 450.0
			targetHeightPt := 320.0
			// Se o arquivo PNG original existe, preserva o aspect ratio exato
			if len(trimmed) > 0 {
				if pngBytes, errRead := os.ReadFile(trimmed); errRead == nil {
					if cfg, _, errCfg := image.DecodeConfig(bytes.NewReader(pngBytes)); errCfg == nil && cfg.Width > 0 && cfg.Height > 0 {
						wPt := float64(cfg.Width)
						hPt := float64(cfg.Height)
						scale := math.Min(1.0, math.Min(450.0/wPt, 620.0/hPt))
						targetWidthPt = wPt * scale
						targetHeightPt = hPt * scale
					}
				}
			}
			log.Printf("[INFO] Inserindo imagem inline via Mermaid Pako URL (%d bytes, dimensões: %.1fx%.1f pt)", len(pakoURL), targetWidthPt, targetHeightPt)
			return insertInlineImageDirect(ctx, documentID, targetIndex, pakoURL, targetWidthPt, targetHeightPt, caption, source)
		}
	}

	// 2. Se for arquivo PNG binário e não houver .mmd, executa o fluxo com upload/fallback
	pngBytes, err := os.ReadFile(trimmed)
	if err != nil {
		return nil, fmt.Errorf("falha ao processar diagrama: arquivo '%s' não encontrado: %w", trimmed, err)
	}

	return InsertMermaidDiagramBytes(ctx, documentID, targetIndex, pngBytes, filepath.Base(trimmed), caption, source)
}

// insertInlineImageDirect insere uma imagem diretamente no Docs a partir de uma URI pública (ex: mermaid.ink)
func insertInlineImageDirect(ctx context.Context, documentID string, targetIndex int64, imageURI string, widthPt, heightPt float64, caption, source string) (*MermaidInsertionResult, error) {
	docsSvc, err := GetDocsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter docs service: %w", err)
	}

	var requests []*docs.Request
	currentOffset := targetIndex

	// 1. Legenda no Topo
	captionText := strings.TrimSpace(caption)
	if captionText != "" {
		upperCap := strings.ToUpper(captionText)
		if !strings.HasPrefix(upperCap, "FIGURA ") && !strings.HasPrefix(upperCap, "FIGURA –") && !strings.HasPrefix(upperCap, "FIGURA -") && !strings.HasPrefix(upperCap, "GRÁFICO") {
			captionText = fmt.Sprintf("Figura – %s", captionText)
		}
		capText := fmt.Sprintf("%s\n", captionText)
		capLen := UTF16Length(capText)
		requests = append(requests, &docs.Request{
			InsertText: &docs.InsertTextRequest{
				Location: &docs.Location{Index: currentOffset},
				Text:     capText,
			},
		})
		requests = append(requests, &docs.Request{
			UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
				Range: &docs.Range{
					StartIndex: currentOffset,
					EndIndex:   currentOffset + capLen,
				},
				ParagraphStyle: &docs.ParagraphStyle{
					Alignment:      "CENTER",
					NamedStyleType: "NORMAL_TEXT",
					KeepWithNext:   true,
				},
				Fields: "alignment,namedStyleType,keepWithNext",
			},
		})
		currentOffset += capLen
	}

	// 2. Inserção Inline da Imagem
	requests = append(requests, &docs.Request{
		InsertInlineImage: &docs.InsertInlineImageRequest{
			Location: &docs.Location{
				Index: currentOffset,
			},
			Uri: imageURI,
			ObjectSize: &docs.Size{
				Width: &docs.Dimension{
					Magnitude: widthPt,
					Unit:      "PT",
				},
				Height: &docs.Dimension{
					Magnitude: heightPt,
					Unit:      "PT",
				},
			},
		},
	})
	requests = append(requests, &docs.Request{
		UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
			Range: &docs.Range{
				StartIndex: currentOffset,
				EndIndex:   currentOffset + 1,
			},
			ParagraphStyle: &docs.ParagraphStyle{
				Alignment:      "CENTER",
				NamedStyleType: "NORMAL_TEXT",
			},
			Fields: "alignment,namedStyleType",
		},
	})
	currentOffset += 1

	// 3. Fonte no Rodapé
	sourceText := strings.TrimSpace(source)
	if sourceText != "" {
		upperSrc := strings.ToUpper(sourceText)
		if !strings.HasPrefix(upperSrc, "FONTE:") && !strings.HasPrefix(upperSrc, "FONTE ") {
			sourceText = fmt.Sprintf("Fonte: %s", sourceText)
		}
		sourceFinal := fmt.Sprintf("\n%s\n\n", sourceText)
		sourceLen := UTF16Length(sourceFinal)
		requests = append(requests, &docs.Request{
			InsertText: &docs.InsertTextRequest{
				Location: &docs.Location{Index: currentOffset},
				Text:     sourceFinal,
			},
		})
		requests = append(requests, &docs.Request{
			UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
				Range: &docs.Range{
					StartIndex: currentOffset,
					EndIndex:   currentOffset + sourceLen,
				},
				ParagraphStyle: &docs.ParagraphStyle{
					Alignment:      "CENTER",
					NamedStyleType: "NORMAL_TEXT",
					KeepWithNext:   true,
				},
				Fields: "alignment,namedStyleType,keepWithNext",
			},
		})
	}

	_, err = docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: requests}).Do()
	if err != nil {
		return nil, fmt.Errorf("erro ao inserir imagem inline no Google Docs: %w", err)
	}

	return &MermaidInsertionResult{
		WidthPt:  widthPt,
		HeightPt: heightPt,
		IsInline: true,
		LinkURL:  imageURI,
		Message:  fmt.Sprintf("SUCESSO: Diagrama Mermaid renderizado e inserido inline nativamente com dimensões %.1f x %.1f pt.", widthPt, heightPt),
	}, nil
}

// InsertMermaidDiagramBytes realiza o processo a partir de um buffer de bytes PNG
func InsertMermaidDiagramBytes(ctx context.Context, documentID string, targetIndex int64, pngBytes []byte, filename, caption, source string) (*MermaidInsertionResult, error) {
	if len(pngBytes) == 0 {
		return nil, fmt.Errorf("buffer PNG de diagrama vazio")
	}

	// 1. Decodificação rápida de dimensões com guarda contra divisão por zero
	cfg, _, err := image.DecodeConfig(bytes.NewReader(pngBytes))
	if err != nil {
		return nil, fmt.Errorf("falha ao decodificar dimensões da imagem PNG: %w", err)
	}

	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, fmt.Errorf("dimensões inválidas da imagem rasterizada: %dx%d", cfg.Width, cfg.Height)
	}

	// 2. Enquadramento Bidimensional (2D Bounding Box) na folha A4 com margens ABNT
	maxWidthPt := 450.0
	maxHeightPt := 620.0

	wPt := float64(cfg.Width)
	hPt := float64(cfg.Height)

	scale := math.Min(1.0, math.Min(maxWidthPt/wPt, maxHeightPt/hPt))
	targetWidthPt := wPt * scale
	targetHeightPt := hPt * scale

	// 3. Upload para o Google Drive
	driveSvc, err := GetDriveService(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter serviço do Google Drive: %w", err)
	}

	driveFile := &drive.File{
		Name:     filename,
		MimeType: "image/png",
	}

	if docMeta, errDoc := driveSvc.Files.Get(documentID).Fields("parents").Do(); errDoc == nil && len(docMeta.Parents) > 0 {
		driveFile.Parents = []string{docMeta.Parents[0]}
	}

	uploadedFile, err := driveSvc.Files.Create(driveFile).
		Media(bytes.NewReader(pngBytes)).
		Fields("id, webContentLink, webViewLink").
		Do()
	if err != nil {
		log.Printf("[WARN] Upload de PNG no Drive falhou (%v). Ativando fallback descritivo ABNT.", err)
		docsSvc, errDocSvc := GetDocsService(ctx)
		if errDocSvc != nil {
			return nil, errDocSvc
		}
		fallbackCap := strings.TrimSpace(caption)
		if fallbackCap != "" {
			upperCap := strings.ToUpper(fallbackCap)
			if !strings.HasPrefix(upperCap, "FIGURA ") && !strings.HasPrefix(upperCap, "FIGURA –") && !strings.HasPrefix(upperCap, "FIGURA -") && !strings.HasPrefix(upperCap, "GRÁFICO") {
				fallbackCap = fmt.Sprintf("Figura – %s", fallbackCap)
			}
		}
		fallbackSrc := strings.TrimSpace(source)
		if fallbackSrc != "" {
			upperSrc := strings.ToUpper(fallbackSrc)
			if !strings.HasPrefix(upperSrc, "FONTE:") && !strings.HasPrefix(upperSrc, "FONTE ") {
				fallbackSrc = fmt.Sprintf("Fonte: %s", fallbackSrc)
			}
		}
		fallbackText := fmt.Sprintf("\n%s\n[Diagrama renderizado no PDF compilado: %s]\n%s\n\n",
			fallbackCap, filename, fallbackSrc)
		req := &docs.BatchUpdateDocumentRequest{
			Requests: []*docs.Request{
				{
					InsertText: &docs.InsertTextRequest{
						Location: &docs.Location{Index: targetIndex},
						Text:     fallbackText,
					},
				},
			},
		}
		_, errBatch := docsSvc.Documents.BatchUpdate(documentID, req).Do()
		if errBatch != nil {
			return nil, fmt.Errorf("erro ao inserir fallback: %w", errBatch)
		}
		return &MermaidInsertionResult{
			WidthPt:  targetWidthPt,
			HeightPt: targetHeightPt,
			IsInline: false,
			Message:  fmt.Sprintf("AVISO: Cota de Service Account no Drive excedida. Fallback ABNT inserido apontando para '%s' compilado no PDF.", filename),
		}, nil
	}

	fileID := uploadedFile.Id
	webViewLink := uploadedFile.WebViewLink
	if webViewLink == "" {
		webViewLink = fmt.Sprintf("https://drive.google.com/file/d/%s/view", fileID)
	}

	// 4. Concessão de permissão temporária de leitura para o Docs conseguir baixar a imagem
	permErr := grantPublicReadPermission(driveSvc, fileID)
	if permErr != nil {
		log.Printf("[WARN] Permissão pública falhou: %v", permErr)
	}

	imageURI := fmt.Sprintf("https://drive.google.com/uc?id=%s&export=download", fileID)
	return insertInlineImageDirect(ctx, documentID, targetIndex, imageURI, targetWidthPt, targetHeightPt, caption, source)
}

func grantPublicReadPermission(driveSvc *drive.Service, fileID string) error {
	perm := &drive.Permission{
		Role: "reader",
		Type: "anyone",
	}
	_, err := driveSvc.Permissions.Create(fileID, perm).Do()
	return err
}
