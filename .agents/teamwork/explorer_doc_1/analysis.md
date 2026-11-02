# Relatório de Diagnóstico Empírico do Documento Google Docs Real

**Documento Alvo**: `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`  
**Título Oficial no Google Docs**: `TCC Teste IA`  
**Data da Investigação**: 2026-11-01  
**Autor**: `explorer_doc_1` (Teamwork Explorer / Investigator)  
**Marco**: R2 — Diagnóstico Empírico do Documento Real e Correlação com Codebase Go  

---

## 1. Sumário Executivo

A inspeção empírica direta do documento real Google Docs `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`, confrontada com a execução das ferramentas MCP (`lint_document_abnt`, `inspect_doc_layout`, `get_doc_skeleton`, `read_doc_content`) e com a análise estática do código-fonte Go (`internal/gdocs`), revelou um cenário de **extrema desconexão entre a realidade do documento e o diagnóstico das ferramentas de auditoria**:

1. **Falso Negativo Crítico no Linter**: O `lint_document_abnt` reportou `0 anomalias/destaques identificados` e declarou que a Capa e Folha de Rosto estão em *"plena conformidade"*. No entanto, o documento real possui **dezenas de linhas em branco espúrias**, **não possui numeração de página configurada**, **possui parágrafos vazios intercalando títulos e seções**, e o linter é **100% cego a variações de entrelinhas (LineSpacing)** e **margens de página**.
2. **Cegueira Estrutural das Ferramentas de Inspeção**: A ferramenta `inspect_doc_layout` não extrai nem expõe o campo `LineSpacing`, nem inspeciona o `DocumentStyle` (margens e cabeçalhos da página). A ferramenta `read_doc_content` exporta o conteúdo em Markdown plano, descartando completamente entrelinhas, recuos exatos em centímetros, notas de cabeçalho e quebras de seção. Com isso, qualquer modelo de IA que consuma essas ferramentas torna-se incapaz de enxergar defeitos de diagramação visual e tipográfica da ABNT.
3. **Capa e Folha de Rosto Empíricas**: O preâmbulo institucional está geometricamente cravado no recuo de 8,0 cm (`indent_start_pt: 226.772` e `indent_first_line_pt: 226.772`) com fonte 10 pt e alinhamento justificado, respeitando a NBR 14724. Contudo, a distribuição vertical tanto da capa quanto da folha de rosto foi feita mediante a injeção mecânica de **40 quebras de linha (`\n`) manuais**, em vez de espaçamentos tipográficos nativos (`spaceAbove`/`spaceBelow`).
4. **Quebras Espúrias na Transição Pré-Textual -> Textual**: Entre a Folha de Rosto e a Introdução, foram mapeadas quebras de linha vazias isoladas após cada título e parágrafo (elementos 72, 74, 76, 78, 80). Elas passaram desapercebidas porque o linter desativa a contagem de quebras na zona pré-textual (`currentState == statePreTextual`) e exige no mínimo duas linhas vazias consecutivas (`consecutiveEmptyParas >= 2`) para disparar alerta.
5. **Inexistência de Seções e Paginação ABNT**: O documento possui apenas uma seção global (`SECTION_BREAK` no índice 0). Não existe a divisão canônica em seções exigida pela API do Google Docs para omitir o cabeçalho nas páginas pré-textuais e ativar o número da página a partir da Introdução no canto superior direito.

---

## 2. Metadados Gerais do Documento Inspecionado

Os dados a seguir foram extraídos via chamadas MCP ao documento vivo:

| Propriedade | Valor Observado | Detalhes Técnicos |
|---|---|---|
| **ID do Documento** | `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` | Google Docs v1 |
| **Título do Documento** | `TCC Teste IA` | Retornado por `inspect_doc_layout` e `get_doc_skeleton` |
| **Total de Elementos Estruturais** | `541` | `len(doc.Body.Content)` |
| **Total de Tabelas Nativas** | `7` | Elementos 128, 214, 227, 282, 299, 319, 416 |
| **Status das Tabelas** | 100% Preenchidas | 0 células vazias detectadas pelo linter |
| **Quebras de Seção (`SECTION_BREAK`)** | `1` única | Elemento 0 [índices 0–1] |
| **Quebras de Página (`PageBreak`)** | Múltiplas | Elemento 33 [286–292] (Capa), Elemento 70 [777–783] (Folha de Rosto) |

---

## 3. Investigação Empírica Ponto a Ponto

### 3.1 Variações de Entrelinhas (Line Spacing) e Espaçamento Vertical (`SpaceAbove`/`SpaceBelow`)

#### Observações Empíricas no Documento
Ao inspecionar os elementos de texto e títulos do documento:
- **Títulos Primários e Secundários**:
  - Elemento 71 (`RESUMO`): `space_above_pt: 18`, `space_below_pt: 12`, `font_size_pt: 12`, `bold: true`, centralizado.
  - Elemento 77 (`1 INTRODUÇÃO`): `space_above_pt: 18`, `space_below_pt: 12`, `font_size_pt: 12`, `bold: true`, alinhamento à esquerda (`START`).
  - Elemento 79 (`Contextualização` - Nível 2): `space_above_pt: 12`, `space_below_pt: 6`, `font_size_pt: 12`.
  - Elemento 85 (`Problema de Pesquisa` - Nível 2): `space_above_pt: 12`, `space_below_pt: 6`, `font_size_pt: 12`.
  - Elemento 93 (`Objetivos Específicos` - Nível 3): `space_above_pt: 12`, `space_below_pt: 6`, `font_size_pt: 12`, `bold: true`.
- **Corpo de Texto**:
  - Elemento 73 (Resumo), Elementos 81, 82, 83, 84, 86, 88, 89, 91 (Corpo da Introdução):
    Possuem `alignment: JUSTIFIED`, `indent_first_line_pt: 35.433` (exatamente 1,25 cm: $1{,}25 \times 72 / 2{,}54 = 35{,}433\text{ pt}$), `font_size_pt: 12`.
- **Omissão do Atributo `LineSpacing`**:
  O JSON retornado por `inspect_doc_layout` **não possui a chave `line_spacing`**.
  Verificando `internal/gdocs/inspector.go` (linhas 9–31):
  ```go
  type ElementLayoutDetail struct {
      Index             int     `json:"index"`
      StartIndex        int64   `json:"start_index"`
      EndIndex          int64   `json:"end_index"`
      Type              string  `json:"type"`
      StyleName         string  `json:"style_name,omitempty"`
      Alignment         string  `json:"alignment,omitempty"`
      IndentStartPt     float64 `json:"indent_start_pt,omitempty"`
      IndentFirstLinePt float64 `json:"indent_first_line_pt,omitempty"`
      SpaceAbovePt      float64 `json:"space_above_pt,omitempty"`
      SpaceBelowPt      float64 `json:"space_below_pt,omitempty"`
      FontSizePt        float64 `json:"font_size_pt,omitempty"`
      FontFamily        string  `json:"font_family,omitempty"`
      Bold              bool    `json:"bold,omitempty"`
      Italic            bool    `json:"italic,omitempty"`
      HasPageBreak      bool    `json:"has_page_break,omitempty"`
      HasInlineObject   bool    `json:"has_inline_object,omitempty"`
      TextSnippet       string  `json:"text_snippet,omitempty"`
      TableRows         int     `json:"table_rows,omitempty"`
      TableCols         int     `json:"table_cols,omitempty"`
      EmptyCells        int     `json:"empty_cells,omitempty"`
  }
  ```
  O campo `LineSpacing` da API do Docs (`p.ParagraphStyle.LineSpacing`) foi **completamente esquecido** na implementação do inspector.

#### Por que o Linter e o Formatador Falharam nas Entrelinhas?
1. **No Linter (`internal/gdocs/linter.go`)**:
   Uma busca por `LineSpacing` ou `line_spacing` no arquivo de 1.257 linhas do linter resulta em **0 ocorrências**. O linter audita cores, famílias de fontes e tamanhos de fontes, mas não possui nenhuma regra para verificar se o espaçamento entrelinhas é 150% no corpo ou 100% nas notas/citações. Documentos com entrelinhas 1,15 (padrão web do Docs) ou com espaçamento simples passam ilesos com zero anomalias.
2. **No Formatador (`internal/gdocs/abnt_writer.go`)**:
   - O formatador aplica `LineSpacing: 150` no Caso J (Corpo do texto).
   - Porém, **parágrafos contidos em tabelas são deliberadamente ignorados** (`if el.Table != nil { continue }`, linha 700). As células de tabela preservam qualquer entrelinhas arbitrário pré-existente (violando a ABNT NBR 14724, que exige entrelinhas 1,0 simples no interior de tabelas).
   - A detecção de citações longas é realizada pela checagem `existingIndent >= 90.0` (linha 819). Se uma citação longa for inserida sem recuo prévio pelo usuário, o formatador a classifica como texto normal e aplica entrelinhas 1,5 e recuo de 1,25 cm, destruindo a estrutura de citação direta da ABNT.

---

### 3.2 Quebras de Linha Espúrias e Linhas em Branco

#### Mapeamento Exato dos Elementos no Documento Real

| Elemento | Intervalo (Índices) | Tipo / Conteúdo | Diagnóstico de Espaçamento |
|---|---|---|---|
| **Elemento 13** | 158–254 | Título da Capa | Texto centralizado 12pt negrito |
| **Elementos 14 a 31** | **254–272** | **18 Parágrafos Vazios Consecutivos** | Cada elemento tem tamanho de 1 caractere (`\n`). Espaçamento vertical empírico manual para empurrar Local/Ano para o rodapé da Capa |
| **Elemento 32** | 272–286 | "LAGARTO -- SE" | Local na Capa |
| **Elemento 33** | 286–292 | "2026" + `PageBreak` | Ano da Capa seguido de Quebra de Página |
| **Elementos 34 a 38** | **292–297** | **5 Parágrafos Vazios Consecutivos** | Linhas em branco no topo da Página 2 antes do Autor da Folha de Rosto |
| **Elemento 39** | 297–329 | "JOSÉ GUSTAVO CORREIA NASCIMENTO" | Autor da Folha de Rosto |
| **Elementos 40 a 46** | **329–336** | **7 Parágrafos Vazios Consecutivos** | Espaçamento manual entre Autor e Título da Folha de Rosto |
| **Elemento 47** | 336–432 | Título na Folha de Rosto | Texto centralizado 12pt negrito |
| **Elementos 48 a 52** | **432–437** | **5 Parágrafos Vazios Consecutivos** | Espaçamento manual entre Título e Preâmbulo da Folha de Rosto |
| **Elemento 53** | 437–697 | Preâmbulo / Nota Institucional | Recuo 8,0 cm (226,772 pt), fonte 10pt, justificado |
| **Elemento 54** | **697–698** | **1 Parágrafo Vazio** | Espaço entre Preâmbulo e Orientador |
| **Elemento 55** | 698–750 | Orientador | Recuo 8,0 cm (226,772 pt), fonte 10pt |
| **Elementos 56 a 68** | **750–763** | **13 Parágrafos Vazios Consecutivos** | Espaçamento manual entre Orientador e Local da Folha de Rosto |
| **Elemento 69** | 763–777 | "LAGARTO -- SE" | Local na Folha de Rosto |
| **Elemento 70** | 777–783 | "2026" + `PageBreak` | Ano da Folha de Rosto seguido de Quebra de Página |
| **Elemento 71** | 783–790 | "RESUMO" | Título Primário (Página 3) |
| **Elemento 72** | **790–791** | **1 Parágrafo Vazio** | Quebra de linha espúria entre RESUMO e o texto do Resumo |
| **Elemento 73** | 791–2816 | Texto do Resumo | Corpo de texto justificado |
| **Elemento 74** | **2816–2817** | **1 Parágrafo Vazio** | Quebra de linha espúria entre o Resumo e Palavras-chave |
| **Elemento 75** | 2817–2961 | "Palavras-chave: ..." | Texto das palavras-chave |
| **Elemento 76** | **2961–2962** | **1 Parágrafo Vazio** | Quebra de linha espúria antes da Introdução |
| **Elemento 77** | 2962–2975 | "1 INTRODUÇÃO" | Primeiro capítulo textual |
| **Elemento 78** | **2975–2976** | **1 Parágrafo Vazio** | Quebra de linha espúria entre o capítulo e a seção 1.1 |
| **Elemento 79** | 2976–2993 | "Contextualização" | Seção secundária |
| **Elemento 80** | **2993–2994** | **1 Parágrafo Vazio** | Quebra de linha espúria entre o subtítulo e o primeiro parágrafo |
| **Elemento 81** | 2994–3646 | "A gestão acadêmica..." | Primeiro parágrafo textual efetivo |

#### Causa do Não Saneamento
1. **Regra de Blindagem Pré-Textual no Linter**:
   No arquivo `internal/gdocs/linter.go` (linhas 733–740):
   ```go
   if trimmed == "" {
       if currentState != statePreTextual {
           consecutiveEmptyParas++
           if consecutiveEmptyParas >= 2 {
               report.ConsecutiveBlankLines = append(report.ConsecutiveBlankLines, ...)
           }
       }
   }
   ```
   Como `currentState` só transita para `stateTextual` ao atingir `"1 INTRODUÇÃO"` (Elemento 77), **todas as quebras de linha dos elementos 0 a 76 são completamente invisíveis para o linter**.
2. **Omissão de Quebras Unitárias**:
   O linter exige `consecutiveEmptyParas >= 2`. Ou seja, linhas em branco únicas inseridas entre seções (Elementos 72, 74, 76, 78, 80) **não são consideradas anomalias**. Na norma ABNT, a separação entre títulos e texto deve ser feita exclusivamente pela métrica de parágrafo (`SpaceAbove`/`SpaceBelow`), nunca pela inserção de parágrafos vazios avulsos.
3. **Comportamento do `HealDocTextAndBlanks`**:
   O formatador saneia apenas quebras duplas consecutivas e preserva a primeira linha em branco após a quebra de página da folha de rosto, fazendo com que o título `"RESUMO"` ou o primeiro texto apareça rebaixado na página seguinte.

---

### 3.3 Renderização da Capa e Folha de Rosto (ABNT NBR 14724)

#### Avaliação Empírica dos Elementos
- **Capa Oficial**:
  - Elemento 1: Logo institucional centralizada (`alignment: CENTER`, `space_below_pt: 14.17`).
  - Elementos 3 a 5: Nome da Instituição, Campus e Curso em caixa alta, centralizados, 12 pt, negrito.
  - Elemento 9: Nome do Autor em caixa alta, centralizado, 12 pt, negrito.
  - Elemento 13: Título do trabalho em caixa alta, centralizado, 12 pt, negrito.
  - Elementos 32 e 33: Local e Ano centralizados, 12 pt, negrito, finalizando com `PageBreak`.
  - *Conformidade*: Os dados textuais e a ordem canônica respeitam a NBR 14724 e as diretrizes institucionais do IFS. A falha técnica reside no uso de 18 quebras de linha repetidas em vez de formatação de espaçamento vertical proporcional.
- **Folha de Rosto**:
  - Elemento 39: Nome do Autor centralizado, 12 pt, sem negrito (em conformidade com a ABNT).
  - Elemento 47: Título em caixa alta, centralizado, 12 pt, negrito.
  - Elemento 53: Preâmbulo / Nota de Apresentação:
    - `alignment: JUSTIFIED`
    - `indent_start_pt: 226.772`
    - `indent_first_line_pt: 226.772`
    - `font_size_pt: 10.0`
    - **Análise Geométrica**: 8,0 cm convertido para pontos tipográficos ($72\text{ pt} / 2{,}54\text{ cm}$) resulta em exatamente **226,77165 pt**, correspondendo precisamente ao valor `226.772 pt` aferido no documento. O preâmbulo está rigorosamente alinhado do meio da página para a margem direita.
  - Elemento 55: Nome do Orientador com mesmo recuo (`indent_start_pt: 226.772`), corpo 10 pt e alinhamento justificado.
  - Elementos 69 e 70: Local e Ano centralizados, corpo 12 pt, finalizando com `PageBreak`.

---

### 3.4 Configuração de Margens de Página (NBR 14724)

- **Exigência Normativa ABNT NBR 14724**:
  - Margem Superior: 3,0 cm (85,039 pt)
  - Margem Esquerda: 3,0 cm (85,039 pt)
  - Margem Inferior: 2,0 cm (56,693 pt)
  - Margem Direita: 2,0 cm (56,693 pt)
- **Diagnóstico no Documento e nas Ferramentas**:
  - A API do Google Docs armazena margens no objeto `DocumentStyle` (`marginTop`, `marginBottom`, `marginLeft`, `marginRight`) e em cada `SectionStyle`.
  - **Cegueira Completa**:
    - `inspect_doc_layout` varre apenas `doc.Body.Content`. Não lê nem expõe `doc.DocumentStyle`.
    - `lint_document_abnt` não contém nenhum teste para checar se as margens do documento conferem com 3/3/2/2 cm.
    - Se um documento recém-criado possuir margens no padrão Google Docs (2,54 cm em todas as bordas), nenhuma ferramenta atual alertará sobre a inconformidade.

---

### 3.5 Numeração de Páginas e Cabeçalhos

- **Regra ABNT NBR 14724 para Numeração**:
  - Todas as folhas a partir da folha de rosto são **contadas**, mas **não numeradas**.
  - A numeração deve figurar a partir da primeira folha da parte textual (**Introdução**) até a última folha do trabalho.
  - O número deve ser posicionado no **canto superior direito**, a 2,0 cm da borda do topo, em algarismos arábicos (fonte 10 pt).
- **Situação no Documento Real**:
  - O documento `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` possui apenas uma quebra de seção (`SECTION_BREAK` no índice 0).
  - **Não existe nenhuma quebra de seção (`SECTION_BREAK`) inserida entre a Folha de Rosto e a Introdução**.
  - Na Google Docs API, um Cabeçalho (`Header`) aplica-se globalmente à seção. Sem um `SectionBreak` com `SectionType: "NEXT_PAGE"` imediatamente antes do título `"1 INTRODUÇÃO"`, qualquer cabeçalho adicionado ao documento vazará e exibirá numeração de página sobre a Capa, Folha de Rosto e Resumo.
  - **No Código Go**: O repositório não implementa nenhuma rotina de criação de cabeçalho (`CreateHeaderRequest`), inserção de número de página (`PageNumber`) ou particionamento de seções para a Introdução.

---

### 3.6 Auditoria Executada pelo `lint_document_abnt`: Relatório Real vs Falsos Negativos

A execução direta de `lint_document_abnt` sobre o documento `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` retornou o seguinte payload JSON:

```json
{
  "total_elements": 541,
  "total_tables": 7,
  "filled_tables": [
    "Elemento 128 [16247-16744]: Tabela 4x3 (Células Vazias: 0/12) [PREENCHIDA]",
    "Elemento 214 [56557-57685]: Tabela 8x5 (Células Vazias: 0/40) [PREENCHIDA]",
    "Elemento 227 [59994-60767]: Tabela 5x3 (Células Vazias: 0/15) [PREENCHIDA]",
    "Elemento 282 [70054-70433]: Tabela 5x3 (Células Vazias: 0/15) [PREENCHIDA]",
    "Elemento 299 [72431-72828]: Tabela 6x3 (Células Vazias: 0/18) [PREENCHIDA]",
    "Elemento 319 [74553-75174]: Tabela 5x4 (Células Vazias: 0/20) [PREENCHIDA]",
    "Elemento 416 [85028-85879]: Tabela 5x4 (Células Vazias: 0/20) [PREENCHIDA]"
  ],
  "cover_status": "Capa Oficial IFS com Logo (3.2cm centralizada) e Folha de Rosto NBR 14724 em plena conformidade.",
  "message": "Auditoria concluída. 0 anomalias/destaques identificados no documento (autoFix=false)."
}
```

#### Matriz de Falsos Negativos Identificados

| Categoria ABNT | Problema Real Observado no Documento | Diagnóstico Reportado pelo Linter | Classificação | Causa no Código (`linter.go`) |
|---|---|---|---|---|
| **Espaçamento Entrelinhas** | Entrelinhas não auditadas no corpo e interior de tabelas | `0 anomalias` | **Falso Negativo Crítico** | Ausência absoluta de inspeção de `LineSpacing` |
| **Margens da Página** | Margens ABNT 3/3/2/2 cm não verificadas | `0 anomalias` | **Falso Negativo Crítico** | `DocumentStyle` não é inspecionado |
| **Numeração de Páginas** | Ausência de quebra de seção e de paginação no cabeçalho superior direito | `0 anomalias` | **Falso Negativo Crítico** | Não audita presença de `Header` nem `SectionBreak` |
| **Quebras Espúrias na Capa** | 18 parágrafos vazios consecutivos (elementos 14–31) | `Capa em plena conformidade` | **Falso Negativo** | Linhas vazias na Capa são ignoradas por premissa |
| **Quebras Espúrias na Folha de Rosto** | 13 parágrafos vazios entre Orientador e Local | `Folha de rosto em plena conformidade` | **Falso Negativo** | Apenas as distâncias Autor->Título e Título->Preâmbulo são checadas |
| **Linhas Vazias Intercaladas** | Parágrafos vazios unitários entre títulos e textos (72, 74, 76, 78, 80) | `0 anomalias` | **Falso Negativo** | Linter ignora elementos antes da Introdução e exige `consecutive >= 2` |
| **Citações Longas (>3 linhas)** | Inexistência de verificação de recuo de 4,0 cm e entrelinhas 1,0 simples | `0 anomalias` | **Falso Negativo** | Não existe checagem de citações longas no linter |

---

## 4. Correlação com as Deficiências das Ferramentas MCP Atuais

A tabela abaixo resume as deficiências identificadas nas 4 ferramentas centrais do servidor MCP:

```
┌───────────────────────────┬────────────────────────────────────────────────────────────────────────┐
│ Ferramenta MCP            │ Deficiências Técnicas Identificadas                                    │
├───────────────────────────┼────────────────────────────────────────────────────────────────────────┤
│ lint_document_abnt        │ 1. Ignora LineSpacing, SpaceAbove e SpaceBelow.                        │
│                           │ 2. Não audita margens da página (DocumentStyle).                       │
│                           │ 3. Não audita paginação em cabeçalhos nem presença de seções.          │
│                           │ 4. Ignora quebras de linha na fase pré-textual.                        │
│                           │ 5. Não identifica linhas em branco isoladas entre títulos e textos.    │
├───────────────────────────┼────────────────────────────────────────────────────────────────────────┤
│ inspect_doc_layout        │ 1. Struct ElementLayoutDetail omite LineSpacing e LineSpacingMode.     │
│                           │ 2. Omite DocumentStyle (margens superior/inferior/esquerda/direita).   │
│                           │ 3. Omite SectionStyle (tipo de seção, headers associados).             │
│                           │ 4. Cegueira transmitida para os agentes de IA.                         │
├───────────────────────────┼────────────────────────────────────────────────────────────────────────┤
│ read_doc_content          │ 1. Serializa para Markdown simples, descartando entrelinhas e recuos.  │
│                           │ 2. Descarta quebras de página nativas e cabeçalhos.                    │
│                           │ 3. Ilude LLMs ao apresentar texto limpo sem a geometria visual real.  │
├───────────────────────────┼────────────────────────────────────────────────────────────────────────┤
│ format_document_abnt      │ 1. Não particiona o documento com SectionBreak antes da Introdução.    │
│                           │ 2. Blindagem de tabelas impede a correção de entrelinhas nas células.  │
│                           │ 3. Heurística frágil de citação longa (<90pt recai em corpo normal).   │
│                           │ 4. Não remove linhas em branco avulsas entre títulos e parágrafos.     │
└───────────────────────────┴────────────────────────────────────────────────────────────────────────┘
```

---

## 5. Recomendações Técnicas para o Plano de Engenharia (R3)

Para sanar em definitivo as discrepâncias entre a norma ABNT e a automação do MCP:

1. **Enriquecimento do `inspect_doc_layout`**:
   - Adicionar os campos `LineSpacing` (float64, ex: 150 para 1.5, 100 para 1.0) e `LineSpacingMode` à struct `ElementLayoutDetail`.
   - Adicionar ao `InspectLayoutReport` o bloco `DocumentPageStyle` contendo as margens atuais em centímetros e pontos (`marginTop`, `marginBottom`, `marginLeft`, `marginRight`) e dimensões da página.
   - Adicionar metadados de cabeçalho (`HeaderID`, presença de numeração de página e seção associada).
2. **Evolução do `lint_document_abnt`**:
   - Implementar checagem de `LineSpacing` com tolerância (corpo = 150%, citações/legendas/tabelas = 100%).
   - Implementar auditoria de margens (3cm superior/esquerda, 2cm inferior/direita).
   - Auditar a existência de `SectionBreak` precedente à Introdução e a ativação de número de página superior direito.
   - Auditar e sinalizar linhas em branco isoladas desnecessárias entre títulos e corpo de texto.
3. **Aprimoramento do `format_document_abnt`**:
   - Inserir cirurgicamente `SectionBreak` (`NEXT_PAGE`) antes do primeiro capítulo numerado (`1 INTRODUÇÃO`).
   - Configurar cabeçalho na Seção 2 com alinhamento à direita e inserção do campo nativo de número de página.
   - Desvincular o cabeçalho da Seção 1 (`UseFirstPageHeaderFooter = true` ou cabeçalho nulo na Seção 1).
   - Formatar o conteúdo interno das células de tabelas com espaçamento simples (`LineSpacing: 100`) e fonte 10 pt.
   - Aprimorar o algoritmo de detecção de citações diretas longas (bloqueando a aplicação de 150% e 1,25 cm em trechos recuados).
