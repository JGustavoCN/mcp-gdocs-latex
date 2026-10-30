package gdocs

import (
	"context"
	"fmt"
	"strings"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/latex"
	"google.golang.org/api/docs/v1"
)

// TableInsertionResult representa o resultado auditável da criação e preenchimento da tabela no Docs
type TableInsertionResult struct {
	TableStartIndex int64  `json:"table_start_index"`
	TableEndIndex   int64  `json:"table_end_index"`
	RowsCount       int64  `json:"rows_count"`
	ColsCount       int64  `json:"cols_count"`
	Message         string `json:"message"`
}

// InsertNativeTable executa o algoritmo Two-Step para inserir e estilizar uma tabela nativa ABNT no Docs
func InsertNativeTable(ctx context.Context, documentID string, targetIndex int64, tableData latex.LaTeXTableData) (*TableInsertionResult, error) {
	if tableData.RowsCount <= 0 || tableData.ColsCount <= 0 {
		return nil, fmt.Errorf("tabela vazia ou com dimensões inválidas: %dx%d", tableData.RowsCount, tableData.ColsCount)
	}

	svc, err := GetDocsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	// ─── Passo 1: Inserção da Legenda e Criação da Estrutura Vazia ───
	var step1Requests []*docs.Request
	currentOffset := targetIndex

	captionText := ""
	if strings.TrimSpace(tableData.Caption) != "" {
		captionText = fmt.Sprintf("Tabela – %s\n", strings.TrimSpace(tableData.Caption))
		step1Requests = append(step1Requests, &docs.Request{
			InsertText: &docs.InsertTextRequest{
				Location: &docs.Location{Index: currentOffset},
				Text:     captionText,
			},
		})
		capLen := UTF16Length(captionText)
		step1Requests = append(step1Requests, &docs.Request{
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

	tableStartIndex := currentOffset
	step1Requests = append(step1Requests, &docs.Request{
		InsertTable: &docs.InsertTableRequest{
			Rows:    tableData.RowsCount,
			Columns: tableData.ColsCount,
			Location: &docs.Location{
				Index: tableStartIndex,
			},
		},
	})

	_, err = svc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: step1Requests}).Do()
	if err != nil {
		return nil, fmt.Errorf("erro no Passo 1 da criação de tabela: %w", err)
	}

	// ─── Passo 2: Descoberta de Coordenadas Celulares e Revisão ───
	doc, err := svc.Documents.Get(documentID).Do()
	if err != nil {
		return nil, fmt.Errorf("erro ao reler documento para inspecionar tabela: %w", err)
	}

	docRevisionId := doc.RevisionId

	var targetElement *docs.StructuralElement
	if doc.Body != nil {
		for _, el := range doc.Body.Content {
			if el.Table != nil && el.StartIndex >= tableStartIndex {
				targetElement = el
				break
			}
		}
	}

	if targetElement == nil || targetElement.Table == nil {
		return nil, fmt.Errorf("tabela inserida não foi localizada no corpo do documento em offset >= %d", tableStartIndex)
	}

	targetTable := targetElement.Table
	rowsCount := int64(len(targetTable.TableRows))
	if rowsCount == 0 {
		return nil, fmt.Errorf("tabela encontrada sem linhas")
	}
	colsCount := int64(len(targetTable.TableRows[0].TableCells))

	// Mapeia índices iniciais das células
	cellIndices := make([][]int64, rowsCount)
	for r := int64(0); r < rowsCount; r++ {
		cellIndices[r] = make([]int64, colsCount)
		for c := int64(0); c < colsCount; c++ {
			cell := targetTable.TableRows[r].TableCells[c]
			if len(cell.Content) > 0 {
				cellIndices[r][c] = cell.Content[0].StartIndex
			}
		}
	}
	tableEndIndex := targetElement.EndIndex

	// ─── Passo 3: Preenchimento Atômico em Lote Reverso com Trava de Revisão ───
	var step3Requests []*docs.Request
	actualTableStartIndex := targetElement.StartIndex

	// 3.1 Estilização tabular ABNT/IBGE (Sem bordas verticais, com topo e fundo)
	step3Requests = append(step3Requests, &docs.Request{
		UpdateTableCellStyle: &docs.UpdateTableCellStyleRequest{
			TableStartLocation: &docs.Location{
				Index: actualTableStartIndex,
			},
			TableCellStyle: &docs.TableCellStyle{
				BorderLeft: &docs.TableCellBorder{
					DashStyle: "SOLID",
					Width:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
					Color:     &docs.OptionalColor{Color: &docs.Color{RgbColor: &docs.RgbColor{Red: 0, Green: 0, Blue: 0}}},
				},
				BorderRight: &docs.TableCellBorder{
					DashStyle: "SOLID",
					Width:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
					Color:     &docs.OptionalColor{Color: &docs.Color{RgbColor: &docs.RgbColor{Red: 0, Green: 0, Blue: 0}}},
				},
				BorderTop: &docs.TableCellBorder{
					DashStyle: "SOLID",
					Width:     &docs.Dimension{Magnitude: 1.0, Unit: "PT"},
					Color:     &docs.OptionalColor{Color: &docs.Color{RgbColor: &docs.RgbColor{Red: 0, Green: 0, Blue: 0}}},
				},
				BorderBottom: &docs.TableCellBorder{
					DashStyle: "SOLID",
					Width:     &docs.Dimension{Magnitude: 1.0, Unit: "PT"},
					Color:     &docs.OptionalColor{Color: &docs.Color{RgbColor: &docs.RgbColor{Red: 0, Green: 0, Blue: 0}}},
				},
				PaddingTop:    &docs.Dimension{Magnitude: 4, Unit: "PT"},
				PaddingBottom: &docs.Dimension{Magnitude: 4, Unit: "PT"},
				PaddingLeft:   &docs.Dimension{Magnitude: 6, Unit: "PT"},
				PaddingRight:  &docs.Dimension{Magnitude: 6, Unit: "PT"},
			},
			Fields: "borderLeft,borderRight,borderTop,borderBottom,paddingTop,paddingBottom,paddingLeft,paddingRight",
		},
	})

	// 3.2 Fonte no maior índice (tableEndIndex)
	if strings.TrimSpace(tableData.Source) != "" {
		sourceText := fmt.Sprintf("\nFonte: %s\n", strings.TrimSpace(tableData.Source))
		sourceLen := UTF16Length(sourceText)

		step3Requests = append(step3Requests, &docs.Request{
			InsertText: &docs.InsertTextRequest{
				Location: &docs.Location{Index: tableEndIndex},
				Text:     sourceText,
			},
		})
		step3Requests = append(step3Requests, &docs.Request{
			UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
				Range: &docs.Range{
					StartIndex: tableEndIndex,
					EndIndex:   tableEndIndex + sourceLen,
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

	// Monta matriz de dados completa (Headers + Data)
	fullGrid := make([][]string, 0, rowsCount)
	if len(tableData.Headers) > 0 {
		fullGrid = append(fullGrid, tableData.Headers)
	}
	fullGrid = append(fullGrid, tableData.Data...)

	// 3.3 Preenchimento reverso das células (R-1, C-1 até 0, 0)
	for r := rowsCount - 1; r >= 0; r-- {
		for c := colsCount - 1; c >= 0; c-- {
			var cellVal string
			if r < int64(len(fullGrid)) && c < int64(len(fullGrid[r])) {
				cellVal = SanitizeText(fullGrid[r][c])
			}
			if cellVal != "" {
				step3Requests = append(step3Requests, &docs.Request{
					InsertText: &docs.InsertTextRequest{
						Location: &docs.Location{Index: cellIndices[r][c]},
						Text:     cellVal,
					},
				})
				// Negrito na linha de cabeçalho (r == 0)
				if r == 0 {
					step3Requests = append(step3Requests, &docs.Request{
						UpdateTextStyle: &docs.UpdateTextStyleRequest{
							Range: &docs.Range{
								StartIndex: cellIndices[r][c],
								EndIndex:   cellIndices[r][c] + UTF16Length(cellVal),
							},
							TextStyle: &docs.TextStyle{
								Bold: true,
							},
							Fields: "bold",
						},
					})
				}
			}
		}
	}

	updateStep3 := &docs.BatchUpdateDocumentRequest{
		Requests: step3Requests,
		WriteControl: &docs.WriteControl{
			RequiredRevisionId: docRevisionId,
		},
	}

	_, err = svc.Documents.BatchUpdate(documentID, updateStep3).Do()
	if err != nil {
		return nil, fmt.Errorf("erro no Passo 3 do preenchimento da tabela (trava de concorrência): %w", err)
	}

	return &TableInsertionResult{
		TableStartIndex: actualTableStartIndex,
		TableEndIndex:   tableEndIndex,
		RowsCount:       rowsCount,
		ColsCount:       colsCount,
		Message:         fmt.Sprintf("SUCESSO: Tabela ABNT (%dx%d) inserida e estilizada com norma IBGE/ABNT.", rowsCount, colsCount),
	}, nil
}
