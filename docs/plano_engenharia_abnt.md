# Plano de Engenharia e Relatório Técnico Definitivo: Conformidade ABNT Estrita para o Servidor MCP `mcp-gdocs-latex`

**Documento:** Plano de Engenharia e Arquitetura de Conformidade ABNT  
**Versão:** 2.0.0  
**Data:** 01 de Novembro de 2026  
**Autor:** `worker_plan_1` (Engenheiro de Software & Arquiteto de Sistemas Go / Google Workspace APIs)  
**Status:** Especificação Técnica e Plano de Implementação Aprovado  
**Normas Regulamentares ABNT Alvo:** ABNT NBR 14724:2011, ABNT NBR 6023:2018, ABNT NBR 6024:2012, ABNT NBR 6027:2012, ABNT NBR 10520:2023  
**Documento de Teste Auditado:** `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` (*TCC Teste IA*)  

---

## Índice Analítico
1. [Seção 1: Relatório Técnico de Viabilidade nas APIs Google Docs v1 e Google Drive v3](#seção-1-relatório-técnico-de-viabilidade-nas-apis-google-docs-v1-e-google-drive-v3)
   - 1.1 Tópico 1: Ajuste de Margens e Dimensões de Página (NBR 14724)
   - 1.2 Tópico 2: Alinhamento Justificado e Espaçamento Entrelinhas
   - 1.3 Tópico 3: Recuo de Primeira Linha (1,25 cm) e Blindagem de Indentação
   - 1.4 Tópico 4: Citações Diretas Longas (> 3 Linhas) e Preservação Tipográfica
   - 1.5 Tópico 5: Capa Institucional e Folha de Rosto (NBR 14724)
   - 1.6 Tópico 6: Numeração de Páginas no Canto Superior Direito e Gerenciamento de Seções
   - 1.7 Tópico 7: Sumário Automático (Nativo vs. Gerador Estruturado ABNT)
   - 1.8 Tópico 8: Remoção Programática de Quebras de Linha e Parágrafos Espúrios
2. [Seção 2: Diagnóstico Detalhado do Codebase Go e do Documento de Teste](#seção-2-diagnóstico-detalhado-do-codebase-go-e-do-documento-de-teste)
   - 2.1 Causa Raiz da Não Detecção e Correção de Entrelinhas (1,15 vs 1,5)
   - 2.2 Causa Raiz da Permanência de Quebras de Linha Espúrias após a Folha de Rosto
   - 2.3 Deficiências Estruturais das Ferramentas de Leitura e Inspeção (`read_doc_content` e `inspect_doc_layout`)
   - 2.4 Matriz de Falsos Negativos do Documento de Teste `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`
3. [Seção 3: Plano Detalhado de Engenharia (/plan) por Ferramenta MCP](#seção-3-plano-detalhado-de-engenharia-plan-por-ferramenta-mcp)
   - 3.1 `format_document_abnt`: Refatoração do Motor de Formatação Atômica
   - 3.2 `lint_document_abnt`: Novos Checkers e Fim da Cegueira Pré-Textual
   - 3.3 `inspect_doc_layout`: Exposição Completa de Geometria, Entrelinhas e Seções
   - 3.4 `read_doc_content`: Preservação Estrutural e Modo Enriquecido para IA
   - 3.5 `cover_generator.go`: Eliminação de `\n` Consecutivos e Correção de Limpeza Reversa
4. [Seção 4: Matriz de Cooperação entre Ferramentas e Preservação Estrutural](#seção-4-matriz-de-cooperação-entre-ferramentas-e-preservação-estrutural)
   - 4.1 Heurística Determinística de Classificação de Parágrafos
   - 4.2 Blindagem de Citações Longas e Elementos Pré-Textuais
   - 4.3 Ordem Canônica de Execução do Pipeline MCP
5. [Seção 5: Roteiro de Implementação em Go e Critérios de Aceite](#seção-5-roteiro-de-implementação-em-go-e-critérios-de-aceite)
   - 5.1 Inventário de Arquivos e Funções a Alterar
   - 5.2 Novas Structs e Modelos de Dados em Go
   - 5.3 Suíte de Testes Automatizados com Mocks `httptest.Server`
   - 5.4 Critérios de Aceite Formais e Checklist de Conclusão

---

# Seção 1: Relatório Técnico de Viabilidade nas APIs Google Docs v1 e Google Drive v3

Esta seção analisa exaustivamente a viabilidade matemática, tipográfica e computacional dos 8 requisitos de conformidade ABNT no ecossistema da **Google Docs API v1** (`google.golang.org/api/docs/v1`) e **Google Drive API v3** (`google.golang.org/api/drive/v3`).

---

### 1.1 Tópico 1: Ajuste de Margens e Dimensões de Página (NBR 14724)

#### Norma Regulamentar
Segundo a ABNT NBR 14724:2011 (item 5.1), o documento impresso ou digitalizado deve ser apresentado em papel formato A4 ($21,0\text{ cm} \times 29,7\text{ cm}$) com margens:
- **Margem Superior:** $3,0\text{ cm}$
- **Margem Esquerda:** $3,0\text{ cm}$
- **Margem Inferior:** $2,0\text{ cm}$
- **Margem Direita:** $2,0\text{ cm}$

#### Fórmulas de Conversão DTP Points (Desktop Publishing Points)
A Google Docs API v1 utiliza exclusivamente a unidade tipográfica `PT` baseada na convenção internacional onde 1 polegada ($2,54\text{ cm}$) equivale a 72 pontos:
$$\text{Fator de Conversão} = \frac{72\text{ pt}}{2,54\text{ cm}} \approx 28,346456692913387\text{ pt/cm}$$

1. **Margem Superior ($3,0\text{ cm}$):**
   $$3,0 \times 28,3464567 = 85,03937\text{ pt} \approx \mathbf{85,039\text{ pt}}$$
2. **Margem Esquerda ($3,0\text{ cm}$):**
   $$3,0 \times 28,3464567 = 85,03937\text{ pt} \approx \mathbf{85,039\text{ pt}}$$
3. **Margem Inferior ($2,0\text{ cm}$):**
   $$2,0 \times 28,3464567 = 56,69291\text{ pt} \approx \mathbf{56,693\text{ pt}}$$
4. **Margem Direita ($2,0\text{ cm}$):**
   $$2,0 \times 28,3464567 = 56,69291\text{ pt} \approx \mathbf{56,693\text{ pt}}$$
5. **Dimensões do Papel A4 ($21,0 \times 29,7\text{ cm}$):**
   $$\text{Largura} = 21,0 \times 28,3464567 = 595,27559\text{ pt} \approx \mathbf{595,276\text{ pt}}$$
   $$\text{Altura} = 29,7 \times 28,3464567 = 841,88976\text{ pt} \approx \mathbf{841,890\text{ pt}}$$
6. **Mancha Gráfica Útil (Área Textual):**
   $$\text{Largura Útil} = 595,276 - (85,039 + 56,693) = \mathbf{453,544\text{ pt}} \ (16,0\text{ cm})$$

#### Payload Atômico na Google Docs API v1
O ajuste é realizado via `UpdateDocumentStyleRequest`. Conforme a especificação autoritativa da Google Docs API:
> *"Updating the top/bottom/left/right page margin on the document style clears that margin on all section styles."*

Essa propriedade garante que a chamada global no `DocumentStyle` sobreponha e elimine automaticamente quaisquer margens anômalas residuais em todas as seções existentes do documento:

```json
{
  "updateDocumentStyle": {
    "documentStyle": {
      "marginTop": { "magnitude": 85.039, "unit": "PT" },
      "marginLeft": { "magnitude": 85.039, "unit": "PT" },
      "marginBottom": { "magnitude": 56.693, "unit": "PT" },
      "marginRight": { "magnitude": 56.693, "unit": "PT" },
      "pageSize": {
        "width": { "magnitude": 595.276, "unit": "PT" },
        "height": { "magnitude": 841.890, "unit": "PT" }
      }
    },
    "fields": "marginTop,marginLeft,marginBottom,marginRight,pageSize"
  }
}
```

---

### 1.2 Tópico 2: Alinhamento Justificado e Espaçamento Entrelinhas

#### Norma Regulamentar
Conforme a ABNT NBR 14724:2011 (item 5.2):
- **Corpo do texto:** Alinhamento **justificado** com espaçamento entrelinhas de **1,5**. Os parágrafos normais contíguos não devem ter espaçamento vertical extra (`spaceAbove: 0 pt` e `spaceBelow: 0 pt`). A separação visual entre parágrafos dá-se estritamente pelo recuo de primeira linha.
- **Espaçamento simples (1,0):** Obrigatório em citações diretas longas (> 3 linhas), notas de rodapé, referências bibliográficas ao final do trabalho (NBR 6023), legendas e fontes de ilustrações/tabelas, preâmbulo da folha de rosto e células internas de tabelas.

#### Semântica do Modelo Google Docs
Na struct Go `docs.ParagraphStyle`:
- `LineSpacing`: valor `float64` representando porcentagem.
  - $150.0 = 150\%$ (Espaçamento 1,5 da ABNT).
  - $100.0 = 100\%$ (Espaçamento simples 1,0 da ABNT).
  - Padrão do Google Docs: $115.0$ acompanhado de `spaceBelow: 10 pt` (o que fere a ABNT se não for sobrescrito).
- `SpaceAbove` e `SpaceBelow`: do tipo `*docs.Dimension` (`{ "magnitude": float64, "unit": "PT" }`). No corpo de texto regular, devem ser forçados a `0 PT`. Em títulos primários e secundários, utiliza-se `spaceAbove: 18 PT` e `spaceBelow: 12 PT` para emular o espaçamento canônico de 1 linha em branco antes e depois.

#### Payload Atômico para Corpo de Texto (1,5 Justificado)
```json
{
  "updateParagraphStyle": {
    "range": { "startIndex": 2994, "endIndex": 3646 },
    "paragraphStyle": {
      "namedStyleType": "NORMAL_TEXT",
      "alignment": "JUSTIFIED",
      "lineSpacing": 150.0,
      "spaceAbove": { "magnitude": 0, "unit": "PT" },
      "spaceBelow": { "magnitude": 0, "unit": "PT" },
      "avoidWidowAndOrphan": true
    },
    "fields": "namedStyleType,alignment,lineSpacing,spaceAbove,spaceBelow,avoidWidowAndOrphan"
  }
}
```

---

### 1.3 Tópico 3: Recuo de Primeira Linha (1,25 cm) e Blindagem de Indentação

#### Norma Regulamentar
A NBR 14724 estabelece o recuo de primeira linha de **1,25 cm** para os parágrafos do corpo textual. Títulos, citações longas, referências e legendas **não possuem recuo de primeira linha**.

#### Conversão Métrica
$$1,25\text{ cm} \times 28,3464567 = 35,43307\text{ pt} \approx \mathbf{35,433\text{ pt}}$$

#### Aritmética de Indentação e Blindagem Residual na Google Docs API
A API do Google Docs expõe em `ParagraphStyle`:
- `indentStart`: deslocamento da margem esquerda do parágrafo inteiro em relação à margem da página.
- `indentFirstLine`: deslocamento da primeira linha do parágrafo.

**Vulnerabilidade Residual Identificada:** Se o parágrafo tiver origem em um item de lista (`Bullet`) ou citação anterior, `indentStart` mantém um valor residual (ex.: 36 pt). Definir apenas `indentFirstLine: 35.433 pt` fará com que a primeira linha comece em $36 + 35,433 = 71,433\text{ pt}$ e as demais em $36\text{ pt}$.
**Blindagem Obrigatória:** Para garantir conformidade irrestrita, o builder deve obrigatoriamente forçar `indentStart: 0 PT` no mesmo payload em que define `indentFirstLine: 35.433 PT`:

```json
{
  "updateParagraphStyle": {
    "range": { "startIndex": 2994, "endIndex": 3646 },
    "paragraphStyle": {
      "indentFirstLine": { "magnitude": 35.433, "unit": "PT" },
      "indentStart": { "magnitude": 0, "unit": "PT" }
    },
    "fields": "indentFirstLine,indentStart"
  }
}
```

---

### 1.4 Tópico 4: Citações Diretas Longas (> 3 Linhas) e Preservação Tipográfica

#### Norma Regulamentar (NBR 10520:2023 e NBR 14724:2011)
- Citações com mais de 3 linhas de extensão textual:
  - Recuo em bloco de **4,0 cm** a partir da margem esquerda;
  - Fonte reduzida: **10 pt** (quando o corpo do trabalho for 12 pt);
  - Espaçamento entrelinhas **simples (1,0 / 100%)**;
  - Alinhamento **justificado**;
  - **Sem aspas** no início ou fim;
  - Espaçamento anterior e posterior de 1 linha em branco (ou `spaceAbove: 6 PT`, `spaceBelow: 6 PT`).

#### Conversão Métrica
$$4,0\text{ cm} \times 28,3464567 = 113,3858\text{ pt} \approx \mathbf{113,386\text{ pt}}$$

#### Payload Atômico Composto na Google Docs API
Exige composição atômica de `UpdateParagraphStyleRequest` e `UpdateTextStyleRequest`:

```json
[
  {
    "updateParagraphStyle": {
      "range": { "startIndex": 4500, "endIndex": 4820 },
      "paragraphStyle": {
        "namedStyleType": "NORMAL_TEXT",
        "alignment": "JUSTIFIED",
        "lineSpacing": 100.0,
        "indentStart": { "magnitude": 113.386, "unit": "PT" },
        "indentFirstLine": { "magnitude": 0, "unit": "PT" },
        "spaceAbove": { "magnitude": 6.0, "unit": "PT" },
        "spaceBelow": { "magnitude": 6.0, "unit": "PT" },
        "avoidWidowAndOrphan": true
      },
      "fields": "namedStyleType,alignment,lineSpacing,indentStart,indentFirstLine,spaceAbove,spaceBelow,avoidWidowAndOrphan"
    }
  },
  {
    "updateTextStyle": {
      "range": { "startIndex": 4500, "endIndex": 4820 },
      "textStyle": {
        "fontSize": { "magnitude": 10.0, "unit": "PT" },
        "bold": false,
        "italic": false
      },
      "fields": "fontSize,bold,italic"
    }
  }
]
```

---

### 1.5 Tópico 5: Capa Institucional e Folha de Rosto (NBR 14724)

#### Norma Regulamentar
- **Capa (Página 1):** Nome da Instituição (12 pt, caixa alta, negrito, centralizado), Autor (12 pt, caixa alta, negrito, centralizado), Título (12 pt, caixa alta, negrito, centralizado), Local e Ano no rodapé (12 pt, negrito, centralizado). Sem recuos (`indentStart: 0`, `indentFirstLine: 0`).
- **Folha de Rosto (Página 2):** Autor (12 pt, regular, centralizado), Título (12 pt, negrito, centralizado), Local e Ano (12 pt, regular, centralizado).
- **Preâmbulo / Nota de Apresentação:** Recuo de **8,0 cm** a partir da margem esquerda, fonte **10 pt**, espaçamento simples (100%), alinhamento justificado.
- Orientador e Coorientador: recuados a 8,0 cm, fonte 10 pt, espaçamento simples.

#### Conversão Métrica do Preâmbulo
$$8,0\text{ cm} \times 28,3464567 = 226,77165\text{ pt} \approx \mathbf{226,772\text{ pt}}$$

#### Payload Atômico na Google Docs API
Para formatar o bloco de preâmbulo e orientador de modo que todo o parágrafo se alinhe a partir do meio da página até a margem direita, define-se `indentStart: 226.772 PT` e `indentFirstLine: 226.772 PT`:

```json
{
  "updateParagraphStyle": {
    "range": { "startIndex": 437, "endIndex": 750 },
    "paragraphStyle": {
      "namedStyleType": "NORMAL_TEXT",
      "alignment": "JUSTIFIED",
      "lineSpacing": 100.0,
      "indentStart": { "magnitude": 226.772, "unit": "PT" },
      "indentFirstLine": { "magnitude": 226.772, "unit": "PT" },
      "spaceAbove": { "magnitude": 12.0, "unit": "PT" },
      "spaceBelow": { "magnitude": 0, "unit": "PT" },
      "avoidWidowAndOrphan": true
    },
    "fields": "namedStyleType,alignment,lineSpacing,indentStart,indentFirstLine,spaceAbove,spaceBelow,avoidWidowAndOrphan"
  }
}
```

---

### 1.6 Tópico 6: Numeração de Páginas no Canto Superior Direito e Gerenciamento de Seções

#### Norma Regulamentar (NBR 14724, item 5.3)
1. **Contagem:** Todas as folhas a partir da folha de rosto são **contadas** sequencialmente (a capa não conta).
2. **Exibição:** O número da página só deve figurar visivelmente **a partir da primeira página da parte textual** (Introdução / Capítulo 1).
3. **Localização:** Canto superior direito, a 2,0 cm da borda superior e 2,0 cm da borda direita da folha.
4. **Supressão:** Folha de Rosto, Resumo, Abstract, Listas e Sumário **não possuem numeração aparente**.

#### Arquitetura de Seções e Cabeçalhos na Google Docs API v1
1. **Particionamento de Seção via `InsertSectionBreakRequest`:**
   - Para dissociar as páginas pré-textuais das textuais, é imperativo inserir uma quebra de seção `SECTION_TYPE_NEXT_PAGE` imediatamente antes do título do primeiro capítulo (`1 INTRODUÇÃO`).
   - O documento passa a ter:
     - **Seção 1 (Pré-Textual):** Capa, Folha de Rosto, Resumo, Sumário.
     - **Seção 2 (Textual e Pós-Textual):** Introdução até as Referências.
2. **Isolamento de Cabeçalho via `CreateHeaderRequest`:**
   - Na API v1, um cabeçalho pode ser criado vinculado diretamente a uma seção específica fornecendo o parâmetro `sectionBreakLocation`:
   ```json
   {
     "createHeader": {
       "type": "DEFAULT",
       "sectionBreakLocation": {
         "index": 2962
       }
     }
   }
   ```
   - Regra da API: Ao vincular o cabeçalho exclusivamente à Seção 2 e garantir que a Seção 1 não possua cabeçalho padrão (`DeleteHeaderRequest` no `default_header_id` do documento, se existir), **a Seção 1 não renderiza cabeçalho algum**.
3. **Contagem Contínua de Páginas (`pageNumberStart`):**
   - Na Seção 2, ao manter `pageNumberStart` vazio na `SectionStyle`, o Google Docs continua automaticamente a contagem a partir das páginas da Seção 1.
4. **Constatação Crítica sobre `AutoText` na API v1 REST:**
   - **Restrição Intransponível:** A Google Docs API v1 **NÃO POSSUI NENHUMA REQUISIÇÃO PARA INSERIR AUTOTEXT (`PAGE_NUMBER`) VIA `batchUpdate`**. Não existe `InsertAutoTextRequest` nem `InsertPageNumberRequest`.
   - O elemento `AutoText` com `type: "PAGE_NUMBER"` é uma estrutura somente-leitura exposta pela API quando lida de documentos pré-existentes.
   - **Solução de Engenharia Validada:** O documento base gerado pelo MCP deve ser originado ou clonado a partir de um *Template de Referência* do Google Drive (`resources/template/` ou ID de template canônico) onde o cabeçalho padrão já contenha o campo dinâmico de número de página configurado com alinhamento à direita. O motor Go do MCP então executa o isolamento de seções (`InsertSectionBreakRequest` e limpeza do cabeçalho da Seção 1), garantindo a conformidade sem violar a limitação da API REST.

---

### 1.7 Tópico 7: Sumário Automático (Nativo vs. Gerador Estruturado ABNT)

#### Norma Regulamentar (NBR 6027:2012)
- Reprodução idêntica dos títulos das divisões e subdivisões, com a mesma tipografia da NBR 6024;
- Linha pontilhada contínua (*leader dots*) unindo o final do título ao número da página correspondente;
- Números de página rigorosamente alinhados à margem direita da página.

#### Análise Comparativa de Viabilidade na Google Docs API v1
1. **Sumário Nativo do Google Docs (`TableOfContents`):**
   - **Inviabilidade via API:** A Google Docs API v1 **não possui `InsertTableOfContentsRequest`**. A inserção de sumários nativos só é possível manualmente pela interface gráfica Web do Docs.
   - **Inconformidade Tipográfica:** O TOC nativo do Docs não formata as entradas em estrita conformidade com a NBR 6024 (ex.: seções primárias em caixa alta e negrito, seções secundárias em caixa alta regular, etc.).
2. **TabStops com Leader Dots:**
   - **Inviabilidade via API:** Na struct `docs.ParagraphStyle`, o campo `TabStops` é **somente-leitura (`read-only`)**. A API rejeita requisições que tentem escrever paradas de tabulação. Além disso, o enum de tabulação do Docs suporta apenas `START`, `CENTER` e `END`, **sem suporte a preenchimento com pontos**.
3. **Solução Canônica Recomendada: Gerador Estruturado ABNT com Tabela Nativa Invisível de 2 Colunas:**
   - Criação de uma tabela nativa via `InsertTableRequest` com 2 colunas e N linhas (uma linha para cada seção):
     - **Coluna 1 (Largura: $14,0\text{ cm} \approx 396,85\text{ pt}$):** Título completo com indicativo numérico, formatado na tipografia NBR 6024 apropriada (ex.: negrito para primários, normal para secundários) e alinhado à esquerda (`START`).
     - **Coluna 2 (Largura: $2,0\text{ cm} \approx 56,69\text{ pt}$):** Número da página correspondente, alinhado à direita (`END`).
     - **Estilização Atômica:** Aplica `UpdateTableCellStyleRequest` em lote zerando todas as bordas (`borderTop`, `borderBottom`, `borderLeft`, `borderRight` com largura 0 pt) e zerando os paddings das células.
   - **Vantagem Absoluta:** Alinhamento matemático perfeito na margem direita de 16,0 cm, imune a variações de fontes proporcionais e totalmente suportado pela API v1.

---

### 1.8 Tópico 8: Remoção Programática de Quebras de Linha e Parágrafos Espúrios

#### Restrições Formais de `DeleteContentRangeRequest` na Google Docs API v1
Qualquer violação das regras abaixo resulta em erro fatal HTTP `400 Bad Request: Invalid range`:
1. **Delimitador Final Inamovível:** O caractere `body.EndIndex - 1` é o delimitador de parágrafo final obrigatório do Google Docs. É estritamente proibido deletá-lo.
2. **Parágrafo Antecessor de Tabela:** Toda tabela no Google Docs deve obrigatoriamente ser precedida por um parágrafo. Deletar o `\n` imediatamente anterior a um elemento `Table` dispara erro 400.
3. **Deriva Posicional (*Index Drift*):** Deletar trechos em ordem crescente desloca os índices de todos os caracteres subsequentes, corrompendo as coordenadas dos deletes seguintes.

#### Algoritmo Canônico com Ordem Reversa (*Reverse Index Order Rule*)
1. Varrer `doc.Body.Content` e coletar os intervalos de parágrafos puramente vazios (`trimmedText == ""`).
2. **Guarda Pré-Tabela:** Se o parágrafo vazio for o elemento imediatamente anterior a uma `Table` (`doc.Body.Content[i+1].Table != nil`), ignorá-lo e não agendar sua deleção.
3. **Execução Reversa:** Ordenar a lista de intervalos do maior índice para o menor índice antes de submeter o lote de `DeleteContentRangeRequest`.
4. **Clamping Defensivo:** Garantir que `endIndex <= bodyEndIndex - 1`.

---

# Seção 2: Diagnóstico Detalhado do Codebase Go e do Documento de Teste

Auditoria empírica realizada no código-fonte em `c:\Projetos\golang\mcp-gdocs-latex\internal\gdocs` e confrontada com o documento real `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` (*TCC Teste IA*, 541 elementos estruturais).

---

### 2.1 Causa Raiz da Não Detecção e Correção de Entrelinhas (1,15 vs 1,5)

#### 1. Cegueira Total no Linter (`internal/gdocs/linter.go`)
- No arquivo `linter.go` (1.257 linhas), não existe **nenhuma menção a `LineSpacing`, `SpaceAbove` ou `SpaceBelow`**.
- O linter audita cores de texto, tamanhos de fonte, famílias tipográficas e tabelas vazias, mas **não valida o espaçamento entrelinhas**. Parágrafos com entrelinhas padrão do Google Docs (1,15 / $115\%$) ou simples (1,0) passam com classificação de conformidade de 100%.

#### 2. Bypass no Formatador (`internal/gdocs/abnt_writer.go`)
- O formatador aplica `ABNTLineSpacing = 150` no Caso J (`abnt_writer.go:1260-1277`), porém:
  - **Tabelas Nativas são Ignoradas:** Na linha 700:
    ```go
    if el.Table != nil {
        continue
    }
    ```
    O comando `continue` pula todos os parágrafos dentro de células de tabela. As células mantêm o espaçamento pré-existente (1,15), violando a NBR 14724 que prescreve espaçamento 1,0 simples no interior de tabelas.
  - **Parágrafos Vazios do Corpo são Pulados:** Na linha 746, parágrafos vazios pós-capa são ignorados com `continue` sem redefinir `LineSpacing` ou `SpaceAbove`/`SpaceBelow`.

---

### 2.2 Causa Raiz da Permanência de Quebras de Linha Espúrias após a Folha de Rosto

#### 1. Injeção Excessiva em `cover_generator.go`
- O gerador de capa emula o `\vfill` do LaTeX concatenando quebras mecânicas repetidas (`strings.Repeat("\n", 18)` na capa, `strings.Repeat("\n", 14)` na folha de rosto).

#### 2. Regra de Limiar `consecutiveBlanks >= 2` em `HealDocTextAndBlanks`
- Em `internal/gdocs/abnt_writer.go:343-356`:
  ```go
  if t == "" {
      consecutiveBlanks++
      if consecutiveBlanks >= 2 {
          toDeleteRanges = append(toDeleteRanges, &docs.Range{StartIndex: el.StartIndex, EndIndex: el.EndIndex})
      }
  } else {
      consecutiveBlanks = 0
  }
  ```
  A primeira linha em branco após a quebra de página da folha de rosto tem `consecutiveBlanks == 1`. Ela **nunca é adicionada para deleção**. O título `"RESUMO"` é portanto empurrado para baixo, criando uma linha em branco indesejada no topo da página 3.

#### 3. Erro Silenciado na Limpeza Reversa de `cover_generator.go:891-912`
- O loop de limpeza reversa ativa `inResumeTransition = true` ao detectar `"RESUMO"`. O parágrafo imediatamente anterior contém a `PageBreak` da folha de rosto.
- O código tenta deletar esse parágrafo, provocando rejeição pela API ou tentativa de excluir a própria quebra de página.
- Na linha 911, o código silencia o erro com `_, _ = docsSvc.Documents.BatchUpdate(...)`, deixando todas as quebras intactas.

#### 4. Filtro por String Hardcoded `"2026"`
- Em `abnt_writer.go:191-196`, a identificação da folha de rosto depende da ocorrência exata da string `"2026"` duas vezes seguidas. Se o texto contiver `"Lagarto -- SE / 2026"`, a condição falha.

---

### 2.3 Deficiências Estruturais das Ferramentas de Leitura e Inspeção (`read_doc_content` e `inspect_doc_layout`)

#### 1. Omissão Crítica em `inspect_doc_layout` (`internal/gdocs/inspector.go`)
- A struct `ElementLayoutDetail` lista fontes, alinhamentos e recuos, mas **omite completamente o campo `LineSpacing`**!
- Além disso, não inspeciona `DocumentStyle` (margens da página) nem propriedades de seções e cabeçalhos (`SectionStyle`, `HeaderId`). Modelos de IA invocando essa ferramenta são cegos quanto ao espaçamento entrelinhas e margens.

#### 2. Omissão em `read_doc_content` (`internal/gdocs/reader.go`)
- Processa elementos de parágrafo e filtra com `if elem.TextRun == nil { continue }` (linha 280).
- Elementos que contêm `PageBreak` possuem `TextRun == nil`. O `continue` descarta silenciosamente as quebras de página.
- O Markdown resultante não indica onde páginas começam ou terminam, nem expõe metadados de entrelinhas e recuos.

---

### 2.4 Matriz de Falsos Negativos do Documento de Teste `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`

| Defeito Real no Documento | Diagnóstico do `lint_document_abnt` | Causa Raiz Técnica |
|---|---|---|
| **Espaçamento Entrelinhas Desconforme** | `0 anomalias identificadas` | Ausência completa de verificação de `LineSpacing` no linter |
| **Margens da Página não Auditadas** | `0 anomalias identificadas` | `DocumentStyle` não é lido pelo linter |
| **Ausência de Numeração de Página** | `0 anomalias identificadas` | Não valida presença de `CreateHeader` nem `SectionBreak` |
| **18 Quebras de Linha Vazias na Capa** | `Capa em plena conformidade` | Zona pré-textual é blindada contra checagem de linhas em branco |
| **13 Quebras Vazias na Folha de Rosto** | `Folha de rosto em plena conformidade` | Linter audita apenas distâncias parciais e ignora o bloco final |
| **Linhas em Branco Isoladas entre Títulos e Seções** | `0 anomalias identificadas` | Linter ignora elementos antes da Introdução e exige `consecutive >= 2` |
| **Citações Longas não Auditadas** | `0 anomalias identificadas` | Não existe heurística de auditoria de citações longas no linter |

---

# Seção 3: Plano Detalhado de Engenharia (/plan) por Ferramenta MCP

---

### 3.1 `format_document_abnt`: Refatoração do Motor de Formatação Atômica

#### Objetivo
Garantir formatação estrita e idempotente em todo o documento, abrangendo margens ABNT, corpo 1,5, células de tabela em 1,0 simples, preservação de citações longas e saneamento de quebras espúrias.

#### Assinatura e Parâmetros da Tool
```json
{
  "name": "format_document_abnt",
  "description": "Formata todo o documento Google Docs ou trecho específico em conformidade estrita com normas ABNT (NBR 14724, NBR 6023, NBR 6024, NBR 10520).",
  "parameters": {
    "type": "object",
    "properties": {
      "document_url_or_id": {
        "type": "string",
        "description": "ID ou URL do documento Google Docs"
      },
      "start_index": {
        "type": "integer",
        "description": "Offset inicial UTF-16 (opcional, padrão: 1 para o documento completo)"
      },
      "end_index": {
        "type": "integer",
        "description": "Offset final UTF-16 (opcional, padrão: fim do documento)"
      },
      "apply_page_margins": {
        "type": "boolean",
        "description": "Se verdadeiro, aplica margens ABNT (Sup 3cm, Esq 3cm, Inf 2cm, Dir 2cm) no DocumentStyle",
        "default": true
      },
      "format_tables": {
        "type": "boolean",
        "description": "Se verdadeiro, formata células internas de tabelas com espaçamento simples (1,0) e fonte 10pt",
        "default": true
      },
      "sanitize_blank_lines": {
        "type": "boolean",
        "description": "Se verdadeiro, remove linhas em branco espúrias em transições de página e entre parágrafos",
        "default": true
      }
    },
    "required": ["document_url_or_id"]
  }
}
```

#### Regras de Implementação no Código Go (`abnt_writer.go`)
1. **Margens Globais:** Se `apply_page_margins == true`, injetar como primeira operação do lote o `UpdateDocumentStyleRequest` com 85,039 pt (superior/esquerda) e 56,693 pt (inferior/direita).
2. **Varredura e Formatação de Tabelas Nativas:**
   - Remover a cláusula excludente `if el.Table != nil { continue }`.
   - Iterar sobre `el.Table.TableRows[*].TableCells[*].Content[*]`.
   - Para cada parágrafo de célula de tabela, emitir:
     - `UpdateParagraphStyleRequest`: `lineSpacing: 100`, `alignment: "JUSTIFIED"` ou `"START"`, `indentStart: 0`, `indentFirstLine: 0`, `spaceAbove: 0`, `spaceBelow: 0`.
     - `UpdateTextStyleRequest`: `fontSize: 10.0 PT`.
3. **Preservação de Citações Longas:**
   - Parágrafos identificados como Citação Longa (ver Seção 4.1) recebem:
     - `indentStart: 113.386 PT` (4,0 cm), `indentFirstLine: 0 PT`, `lineSpacing: 100`, `fontSize: 10.0 PT`, `spaceAbove: 6.0 PT`, `spaceBelow: 6.0 PT`.
   - O formatador **não deve** sobrescrever esses parágrafos com `NORMAL_TEXT` (1,25 cm / 1,5).
4. **Saneamento Inteligente de Quebras:**
   - Em transições após `PageBreak`, parágrafos puramente vazios imediatamente anteriores ao próximo título textual devem ser deletados mesmo que `consecutiveBlanks == 1`.

---

### 3.2 `lint_document_abnt`: Novos Checkers e Fim da Cegueira Pré-Textual

#### Objetivo
Identificar com 100% de precisão todas as infrações normativas da ABNT, eliminando falsos negativos no documento real.

#### Novos Checkers a Implementar em `internal/gdocs/linter.go`
1. **Checker de Entrelinhas (`checkLineSpacing`):**
   - Comparar `p.ParagraphStyle.LineSpacing`:
     - Parágrafos do corpo textual: deve ser exatamente $150.0$ (tolerância: $148.0$ a $152.0$).
     - Parágrafos em citações longas, notas, referências, células de tabelas e legendas: deve ser $100.0$ (tolerância: $98.0$ a $102.0$).
     - Sinalizar anomalias no novo campo `LineSpacingAnomalies []LineSpacingAnomalyItem`.
2. **Checker de Espaçamento Vertical (`checkParagraphSpacing`):**
   - No corpo de texto regular: se `spaceAbove > 0` ou `spaceBelow > 0`, emitir alerta no novo campo `SpacingAnomalies []ParagraphSpacingAnomalyItem`.
3. **Checker de Margens de Página (`checkPageMargins`):**
   - Inspecionar `doc.DocumentStyle`:
     - Superior: $85,039\text{ pt} \pm 1\text{ pt}$
     - Esquerda: $85,039\text{ pt} \pm 1\text{ pt}$
     - Inferior: $56,693\text{ pt} \pm 1\text{ pt}$
     - Direita: $56,693\text{ pt} \pm 1\text{ pt}$
   - Relatar inconformidades no campo `MarginsConformity`.
4. **Checker de Seções e Numeração (`checkSectionsAndNumbering`):**
   - Verificar se existe `SECTION_BREAK` do tipo `NEXT_PAGE` antes do capítulo `1 INTRODUÇÃO`.
   - Verificar se a Seção 2 possui cabeçalho ativo e se a Seção 1 está com cabeçalho suprimido.
5. **Ajuste na Máquina de Estados:**
   - Remover a blindagem incondicional `if currentState != statePreTextual` para parágrafos vazios. Linhas em branco redundantes na transição da folha de rosto para o Resumo passam a ser auditadas e relatadas.

---

### 3.3 `inspect_doc_layout`: Exposição Completa de Geometria, Entrelinhas e Seções

#### Objetivo
Expor todas as dimensões tipográficas e estruturais para consumo analítico de clientes MCP e modelos de IA.

#### Refatoração da Struct `ElementLayoutDetail` (`inspector.go`)
```go
type ElementLayoutDetail struct {
    Index              int     `json:"index"`
    StartIndex         int64   `json:"start_index"`
    EndIndex           int64   `json:"end_index"`
    Type               string  `json:"type"`
    StyleName          string  `json:"style_name,omitempty"`
    Alignment          string  `json:"alignment,omitempty"`
    LineSpacingPercent float64 `json:"line_spacing_percent,omitempty"` // NOVO: 150 = 1,5; 100 = 1,0
    IndentStartPt      float64 `json:"indent_start_pt,omitempty"`
    IndentFirstLinePt  float64 `json:"indent_first_line_pt,omitempty"`
    SpaceAbovePt       float64 `json:"space_above_pt,omitempty"`
    SpaceBelowPt       float64 `json:"space_below_pt,omitempty"`
    KeepWithNext       bool    `json:"keep_with_next,omitempty"`        // NOVO: proteção anti-órfão
    AvoidWidowOrphan   bool    `json:"avoid_widow_orphan,omitempty"`   // NOVO
    FontSizePt         float64 `json:"font_size_pt,omitempty"`
    FontFamily         string  `json:"font_family,omitempty"`
    Bold               bool    `json:"bold,omitempty"`
    Italic             bool    `json:"italic,omitempty"`
    HasPageBreak       bool    `json:"has_page_break,omitempty"`
    HasInlineObject    bool    `json:"has_inline_object,omitempty"`
    TextSnippet        string  `json:"text_snippet,omitempty"`
    TableRows          int     `json:"table_rows,omitempty"`
    TableCols          int     `json:"table_cols,omitempty"`
    EmptyCells         int     `json:"empty_cells,omitempty"`
}
```

#### Enriquecimento de `InspectLayoutReport`
Adicionar blocos de nível de documento:
```go
type InspectLayoutReport struct {
    DocumentID       string                   `json:"document_id"`
    DocumentTitle    string                   `json:"document_title"`
    TotalElements    int                      `json:"total_elements"`
    DocumentGeometry DocumentGeometryDetail   `json:"document_geometry"` // NOVO
    Sections         []SectionLayoutDetail    `json:"sections"`          // NOVO
    Elements         []ElementLayoutDetail    `json:"elements"`
}
```

---

### 3.4 `read_doc_content`: Preservação Estrutural e Modo Enriquecido para IA

#### Objetivo
Permitir que agentes LLM visualizem quebras de página nativas e metadados geométricos do documento.

#### Modificações em `internal/gdocs/reader.go`
1. **Preservação de `PageBreak`:**
   Ao iterar sobre `paragraph.Elements`, se `elem.PageBreak != nil`, emitir o marcador textual:
   ```markdown
   \n\n[PAGE_BREAK]\n\n
   ```
2. **Preservação de `SectionBreak`:**
   Se `elem.SectionBreak != nil`, emitir:
   ```markdown
   \n\n[SECTION_BREAK: NEXT_PAGE]\n\n
   ```
3. **Parâmetro `include_metadata` (booleano, opcional):**
   Quando verdadeiro, cada parágrafo é precedido por um comentário estrutural:
   ```markdown
   <!-- [OFFSET: 2994-3646 | STYLE: NORMAL_TEXT | LINE_SPACING: 1.5 | INDENT: 1.25cm] -->
   A gestão acadêmica no contexto do Instituto Federal...
   ```

---

### 3.5 `cover_generator.go`: Eliminação de `\n` Consecutivos e Correção de Limpeza Reversa

#### Objetivo
Eliminar a injeção maciça de caracteres `\n` na geração da capa e folha de rosto, e corrigir a falha silenciosa na transição para o Resumo.

#### Modificações em `internal/gdocs/cover_generator.go`
1. **Espaçamento Paramétrico por Parágrafo:**
   - Substituir `strings.Repeat("\n", 18)` no pós-título da Capa por um parágrafo estruturado com `spaceAbove: 140 PT` (equivalente ao espaçamento visual desejado) ou quebra de página limpa.
   - Na Folha de Rosto, substituir os blocos de `\n` por `spaceAbove` calibrado nos parágrafos do Preâmbulo e Local/Ano.
2. **Correção do Loop de Limpeza Reversa (`cover_generator.go:891-912`):**
   - Ao varrer de trás para frente a partir de `"RESUMO"`, **verificar se o parágrafo contém `PageBreak`**:
     ```go
     if hasPageBreak(el) {
         // Não deletar o parágrafo da quebra de página da folha de rosto!
         continue
     }
     ```
   - Tratar formalmente o erro de `docsSvc.Documents.BatchUpdate`: se retornar erro, registrar no log e propagar para a camada superior em vez de silenciar com `_, _ =`.

---

# Seção 4: Matriz de Cooperação entre Ferramentas e Preservação Estrutural

Esta seção define as salvaguardas que impedem ferramentas de formatação de desconfigurarem elementos especializados como citações longas, tabelas e cabeçalhos de seção.

---

### 4.1 Heurística Determinística de Classificação de Parágrafos

Para garantir que o formatador `format_document_abnt` ou o substituidor `replace_block_abnt` não destruam formatações prévias, cada parágrafo é classificado de acordo com a seguinte máquina determinística de regras:

```
┌────────────────────────────────────────────────────────────────────────┐
│                        Classificação de Parágrafo                      │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │
          ┌─────────────────────────┼─────────────────────────┐
          ▼                         ▼                         ▼
   [Interior de Tabela]       [Pré-Textual]             [Corpo Textual]
          │                         │                         │
          ▼                         ▼                         ▼
   LineSpacing: 100          Recuo: 8,0 cm             Heurística Citação
   Fonte: 10 pt              (Preâmbulo/Orientador)    Longa vs Normal
   Sem recuo 1ª linha        LineSpacing: 100
```

#### Regra de Classificação de Citações Diretas Longas:
Um parágrafo é classificado como **Citação Direta Longa** se atender cumulativamente a:
1. `existingIndent >= 90.0 PT` (indicando intenção de recuo em bloco próximo aos 113,386 pt / 4,0 cm) **OU** extensão de texto $> 3$ linhas acompanhada de fonte $\le 10.5\text{ pt}$;
2. **NÃO** ser título de seção (sem indicativo numérico de Heading);
3. **NÃO** ser elemento do pré-textual (não conter palavras-chave de preâmbulo institucional);
4. **NÃO** pertencer a células de tabelas.

**Ação de Preservação:** Parágrafos classificados como citação longa recebem estritamente `indentStart: 113.386 PT`, `indentFirstLine: 0 PT`, `lineSpacing: 100.0` e `fontSize: 10.0 PT`.

---

### 4.2 Blindagem de Citações Longas e Elementos Pré-Textuais

| Elemento Estrutural | Recuo Inicial (`indentStart`) | Recuo 1ª Linha (`indentFirstLine`) | Entrelinhas (`lineSpacing`) | Tamanho da Fonte | Ação do Formatador Global |
|---|---|---|---|---|---|
| **Corpo de Texto Regular** | `0 PT` | `35.433 PT` ($1,25\text{ cm}$) | `150.0` ($1,5$) | `12.0 PT` | Aplica `NORMAL_TEXT` ABNT |
| **Citação Direta Longa** | `113.386 PT` ($4,0\text{ cm}$) | `0 PT` | `100.0` ($1,0$) | `10.0 PT` | **Preserva e calibra a 4cm/10pt/1.0** |
| **Preâmbulo da Folha de Rosto** | `226.772 PT` ($8,0\text{ cm}$) | `226.772 PT` ($8,0\text{ cm}$) | `100.0` ($1,0$) | `10.0 PT` | **Preserva e calibra a 8cm/10pt/1.0** |
| **Célula de Tabela Nativa** | `0 PT` | `0 PT` | `100.0` ($1,0$) | `10.0 PT` | Formata para simples/10pt |
| **Legenda / Fonte de Ilustração** | `0 PT` | `0 PT` | `100.0` ($1,0$) | `10.0 PT` | Centralizado, simples/10pt |
| **Referências (NBR 6023)** | `0 PT` | `0 PT` | `100.0` ($1,0$) | `12.0 PT` | Alinhado à esquerda, simples/12pt |

---

### 4.3 Ordem Canônica de Execução do Pipeline MCP

Para formatar ou atualizar um documento sem efeitos colaterais indesejados, o cliente ou orquestrador MCP deve seguir o fluxo canônico:

1. **`inspect_doc_layout` / `get_doc_skeleton`:** Leitura diagnóstica completa dos offsets e geometria atual.
2. **`lint_document_abnt`:** Auditoria prévia gerando lista formal de desvios normativos.
3. **`format_document_abnt`:**
   - Passo 3.1: Ajuste de margens de página (`UpdateDocumentStyleRequest`).
   - Passo 3.2: Particionamento de seção antes da Introdução e configuração do cabeçalho.
   - Passo 3.3: Saneamento reverso de quebras de parágrafo espúrias.
   - Passo 3.4: Formatação dos parágrafos textuais e preservação de citações longas.
   - Passo 3.5: Formatação de células de tabelas.
4. **`lint_document_abnt` (Verificação Pós-Execução):** Confirmação de que as anomalias foram reduzidas a zero.

---

# Seção 5: Roteiro de Implementação em Go e Critérios de Aceite

---

### 5.1 Inventário de Arquivos e Funções a Alterar

| Pacote | Arquivo | Responsabilidade da Modificação | Funções / Métodos Envolvidos |
|---|---|---|---|
| `internal/gdocs` | `inspector.go` | Extrair `LineSpacing`, margens do documento e seções | `InspectDocLayout`, `extractElementDetails` |
| `internal/gdocs` | `linter.go` | Adicionar checagem de `LineSpacing`, `SpaceAbove/Below`, margens e seções | `LintDocument`, `checkLineSpacing`, `checkMargins` |
| `internal/gdocs` | `abnt_writer.go` | Formatar células de tabelas, calibrar citações longas e sanear quebras pós-página | `FormatDocumentABNT`, `formatTableCells`, `HealDocTextAndBlanks` |
| `internal/gdocs` | `reader.go` | Preservar `[PAGE_BREAK]` e `[SECTION_BREAK]` no Markdown | `ReadDocContent`, `buildMarkdownFromDoc` |
| `internal/gdocs` | `cover_generator.go` | Eliminar dezenas de `\n` manuais e tratar erro na transição de página | `GenerateCoverAndPreText`, `cleanSpuriousCoverBreaks` |
| `internal/mcp` | `tools.go` | Atualizar esquemas JSON e descrições das ferramentas enriquecidas | `registerInspectDocLayout`, `registerFormatDocumentAbnt` |

---

### 5.2 Novas Structs e Modelos de Dados em Go

```go
package gdocs

// DocumentGeometryDetail armazena as dimensões físicas da folha e margens
type DocumentGeometryDetail struct {
    PageWidthPt  float64 `json:"page_width_pt"`
    PageHeightPt float64 `json:"page_height_pt"`
    MarginTopPt  float64 `json:"margin_top_pt"`
    MarginTopCm  float64 `json:"margin_top_cm"`
    MarginLeftPt float64 `json:"margin_left_pt"`
    MarginLeftCm float64 `json:"margin_left_cm"`
    MarginBottomPt float64 `json:"margin_bottom_pt"`
    MarginBottomCm float64 `json:"margin_bottom_cm"`
    MarginRightPt  float64 `json:"margin_right_pt"`
    MarginRightCm  float64 `json:"margin_right_cm"`
}

// LineSpacingAnomalyItem representa uma discrepância de entrelinhas auditada pelo linter
type LineSpacingAnomalyItem struct {
    ElementIndex     int     `json:"element_index"`
    StartIndex       int64   `json:"start_index"`
    EndIndex         int64   `json:"end_index"`
    TextSnippet      string  `json:"text_snippet"`
    ObservedSpacing  float64 `json:"observed_spacing"`
    ExpectedSpacing  float64 `json:"expected_spacing"`
    ContextCategory  string  `json:"context_category"` // "CORPO_TEXTUAL", "TABELA", "CITACAO_LONGA", etc.
}

// SectionLayoutDetail descreve a segmentação de seções do documento
type SectionLayoutDetail struct {
    SectionIndex int64  `json:"section_index"`
    StartIndex   int64  `json:"start_index"`
    EndIndex     int64  `json:"end_index"`
    HasHeader    bool   `json:"has_header"`
    HeaderID     string `json:"header_id,omitempty"`
    PageNumberStart *int64 `json:"page_number_start,omitempty"`
}
```

---

### 5.3 Suíte de Testes Automatizados com Mocks `httptest.Server`

Para viabilizar a verificação sem credenciais ativas do Google Workspace, a suíte de testes deve ser ampliada em `internal/gdocs`:

1. **`inspector_test.go`:**
   - Testar extração de `LineSpacing: 150` e `LineSpacing: 115`.
   - Testar leitura de `DocumentStyle` com margens 3/3/2/2 cm.
2. **`linter_test.go`:**
   - Testar detecção de parágrafos com entrelinhas 1,15 em documento simulado.
   - Testar sinalização de quebras de linha espúrias após quebra de página.
3. **`abnt_writer_test.go`:**
   - Testar lote de `batchUpdate` contendo `UpdateParagraphStyleRequest` com `lineSpacing: 100` para células de tabela.
   - Testar preservação de parágrafo de citação longa com recuo de 4,0 cm (113,386 pt).
   - Testar emissão de `UpdateDocumentStyleRequest` com margens ABNT.

---

### 5.4 Critérios de Aceite Formais e Checklist de Conclusão

- [ ] **Margens ABNT:** A ferramenta `format_document_abnt` aplica margens de 3cm (superior/esquerda) e 2cm (inferior/direita) e tamanho A4 via `UpdateDocumentStyleRequest`.
- [ ] **Entrelinhas 1,5 e Simples:** O corpo de texto recebe 150% justificado; citações longas e células de tabelas recebem 100% simples e fonte 10pt.
- [ ] **Auditoria Sem Falsos Negativos:** O `lint_document_abnt` aponta desvios de `LineSpacing` (1,15 vs 1,5), margens anômalas e quebras de linha em transições de página.
- [ ] **Inspeção Completa para IA:** `inspect_doc_layout` retorna `line_spacing_percent`, margens do documento e propriedades de seção.
- [ ] **Leitura com Quebras:** `read_doc_content` exibe explicitamente `[PAGE_BREAK]` e suporta anotações de metadados para LLMs.
- [ ] **Preservação de Citações Longas:** Parágrafos com recuo em bloco ou >3 linhas não são rebaixados para corpo normal durante a formatação global.
- [ ] **Testes 100% Automatizados:** Toda a nova suíte executa e passa em `go test ./...` sem dependência de credenciais externas.
