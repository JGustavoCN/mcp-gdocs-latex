# Relatório de Mineração de Especificações Técnicas: APIs Google Docs v1 e Google Drive v3 para Normas ABNT (NBR 14724)

**Agente:** `spec_miner_api_2`  
**Data:** 01 de Novembro de 2026  
**Contexto do Projeto:** `mcp-gdocs-latex`  
**Escopo:** Investigação exaustiva de viabilidade técnica, mapeamento de payloads, modelos de dados e restrições de baixo nível das APIs do Google Workspace (Google Docs API v1 e Google Drive API v3) para conformidade estrita com as normas da ABNT (NBR 14724, NBR 6023, NBR 6024, NBR 6027, NBR 10520) e diagnóstico de falhas no documento de teste `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`.

---

## 1. Visão Geral e Fontes Autoritativas Analisadas

A presente especificação técnica foi consolidada a partir da mineração rigorosa das seguintes fontes autoritativas:
1. **Google Docs API v1 Official Go Client SDK:** `google.golang.org/api/docs/v1` (v0.310.0 / v0.280.0), mapeado em `docs/pesquisa-api.md`.
2. **Google Drive API v3 Official SDK:** `google.golang.org/api/drive/v3`.
3. **Google Docs Public MCP Toolset Specification:** `docsmcp.googleapis.com/mcp/v1`, mapeado em `docs/pesquisa.md`.
4. **Normas Técnicas Oficiais da ABNT:**
   - **ABNT NBR 14724:2011:** Informação e documentação — Trabalhos acadêmicos — Apresentação.
   - **ABNT NBR 6023:2018:** Informação e documentação — Referências — Elaboração.
   - **ABNT NBR 6024:2012:** Informação e documentação — Numeração progressiva das seções de um documento.
   - **ABNT NBR 6027:2012:** Informação e documentação — Sumário — Apresentação.
   - **ABNT NBR 10520:2023:** Informação e documentação — Citações em documentos — Apresentação.
5. **Código-fonte em Go do Servidor MCP:** Módulos em `internal/gdocs/` (`abnt_writer.go`, `builder.go`, `linter.go`, `inspector.go`, `cover_generator.go`, `structure_normalizer.go`).

---

## 2. Análise Técnica dos 8 Tópicos ABNT

### Tópico 1: Ajuste de Margens e Dimensões de Página (NBR 14724)

#### Norma Regulamentar
Conforme a ABNT NBR 14724 (item 5.1), trabalhos acadêmicos devem ser impressos ou gerados em papel formato A4 ($21,0\text{ cm} \times 29,7\text{ cm}$). Para o anverso da folha, aplicam-se estritamente as margens:
- Margem Superior: $3,0\text{ cm}$
- Margem Esquerda: $3,0\text{ cm}$
- Margem Inferior: $2,0\text{ cm}$
- Margem Direita: $2,0\text{ cm}$

#### Fórmulas de Conversão Métrica para Pontos Tipográficos (DTP Points)
A API do Google Docs opera exclusivamente com a unidade `PT` (Desktop Publishing Points), baseada na convenção internacional onde 1 polegada equivale a 72 pontos e a 2,54 centímetros:
$$\text{Fator de Conversão} = \frac{72\text{ pt}}{2,54\text{ cm}} \approx 28,346456692913387\text{ pt/cm}$$

1. **Margem Superior ($3,0\text{ cm}$):**
   $$3,0 \times 28,3464567 = 85,03937\text{ pt} \approx \mathbf{85,039\text{ pt}}$$
2. **Margem Esquerda ($3,0\text{ cm}$):**
   $$3,0 \times 28,3464567 = 85,03937\text{ pt} \approx \mathbf{85,039\text{ pt}}$$
3. **Margem Inferior ($2,0\text{ cm}$):**
   $$2,0 \times 28,3464567 = 56,69291\text{ pt} \approx \mathbf{56,693\text{ pt}}$$
4. **Margem Direita ($2,0\text{ cm}$):**
   $$2,0 \times 28,3464567 = 56,69291\text{ pt} \approx \mathbf{56,693\text{ pt}}$$
5. **Largura Folha A4 ($21,0\text{ cm}$):**
   $$21,0 \times 28,3464567 = 595,27559\text{ pt} \approx \mathbf{595,276\text{ pt}}$$
6. **Altura Folha A4 ($29,7\text{ cm}$):**
   $$29,7 \times 28,3464567 = 841,88976\text{ pt} \approx \mathbf{841,890\text{ pt}}$$
7. **Largura Útil de Mancha Gráfica:**
   $$21,0 - (3,0 + 2,0) = 16,0\text{ cm} = 16,0 \times 28,3464567 = \mathbf{453,543\text{ pt}}$$

#### Mapeamento de Payload na Google Docs API v1
- **Requisição:** `UpdateDocumentStyleRequest`
- **Mecanismo de Herança de Seções:** A documentação oficial explicita:
  > *"Updating the top/bottom/left/right page margin on the document style clears that margin on all section styles."*
  Ou seja, a aplicação em nível de `DocumentStyle` sobrepõe e limpa automaticamente margens anômalas residuais em todas as seções do documento.

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

### Tópico 2: Alinhamento Justificado e Espaçamento Entrelinhas

#### Norma Regulamentar
Conforme a ABNT NBR 14724 (item 5.2):
- **Corpo do texto:** Alinhamento **justificado**, espaçamento **1,5 entrelinhas**. Os parágrafos de texto corrente são contíguos sem espaçamento adicional entre si (`spaceAbove: 0 pt`, `spaceBelow: 0 pt`). A separação visual entre parágrafos dá-se unicamente pelo recuo de primeira linha.
- **Espaçamento simples (1,0):** Obrigatório em:
  - Citações diretas longas (> 3 linhas);
  - Notas de rodapé;
  - Referências bibliográficas ao final do trabalho (NBR 6023);
  - Legendas e fontes de ilustrações e tabelas;
  - Natureza do trabalho / preâmbulo da folha de rosto;
  - Células internas de tabelas e quadros.

#### Mapeamento de Payload na Google Docs API v1
- **Requisição:** `UpdateParagraphStyleRequest`
- **Semântica do Campo `lineSpacing`:**
  - Na struct Go `docs.ParagraphStyle`: `LineSpacing float64`
  - Na especificação JSON: Porcentagem numérica em que 100.0 representa 100% (espaçamento simples 1,0).
  - Espaçamento 1,5: `lineSpacing: 150` (ou `150.0`).
  - Espaçamento simples 1,0: `lineSpacing: 100` (ou `100.0`).
  - O padrão padrão padrão do Google Docs é `115` (1,15) acompanhado de `spaceBelow: 10 pt`, o que causa infração automática da ABNT quando não normalizado.
- **Campos `spaceAbove` e `spaceBelow`:**
  - Tipo `*docs.Dimension` (`{ "magnitude": number, "unit": "PT" }`).
  - No corpo de texto regular: obrigatoriamente `0 PT`.
  - Em referências (NBR 6023): espaçamento 1,0 simples interno, separado por um espaçamento simples entre referências (`spaceBelow: 6 PT` a `12 PT`).

```json
{
  "updateParagraphStyle": {
    "range": {
      "startIndex": 1050,
      "endIndex": 1820
    },
    "paragraphStyle": {
      "namedStyleType": "NORMAL_TEXT",
      "alignment": "JUSTIFIED",
      "lineSpacing": 150,
      "spaceAbove": { "magnitude": 0, "unit": "PT" },
      "spaceBelow": { "magnitude": 0, "unit": "PT" },
      "avoidWidowAndOrphan": true
    },
    "fields": "namedStyleType,alignment,lineSpacing,spaceAbove,spaceBelow,avoidWidowAndOrphan"
  }
}
```

---

### Tópico 3: Recuo de Primeira Linha de 1,25 cm no Corpo de Texto

#### Norma Regulamentar
A NBR 14724 consagra o recuo de primeira linha de **1,25 cm** para os parágrafos do corpo textual, alinhando visualmente o início de cada ideia. Títulos, citações longas, referências e legendas **não possuem recuo de primeira linha**.

#### Conversão Métrica
$$1,25\text{ cm} \times 28,3464567 = 35,43307\text{ pt} \approx \mathbf{35,433\text{ pt}}$$

#### Aritmética de Indentação na Google Docs API v1
Na API do Google Docs, existem dois campos de indentação concorrentes em `ParagraphStyle`:
1. `indentStart`: Deslocamento da margem inicial do parágrafo inteiro.
2. `indentFirstLine`: Deslocamento da primeira linha.
3. `indentEnd`: Deslocamento da margem final.

**Comportamento Crítico de Acoplamento:**
Se `indentStart` possuir valor residual (ex: 36 pt herdado de listas com marcadores anteriores), definir apenas `indentFirstLine: 35.433 pt` fará com que o parágrafo inteiro fique deslocado para a direita. Para garantir a formatação ABNT rigorosa, **o construtor deve obrigatoriamente resetar `indentStart` para 0 pt concomitantemente**:

```json
{
  "updateParagraphStyle": {
    "range": { "startIndex": 1050, "endIndex": 1820 },
    "paragraphStyle": {
      "indentFirstLine": { "magnitude": 35.433, "unit": "PT" },
      "indentStart": { "magnitude": 0, "unit": "PT" }
    },
    "fields": "indentFirstLine,indentStart"
  }
}
```

---

### Tópico 4: Citação Direta Longa (> 3 Linhas)

#### Norma Regulamentar (NBR 10520:2023 e NBR 14724:2011)
- Citações com mais de 3 linhas de extensão devem ser destacadas com recuo de **4,0 cm** da margem esquerda;
- Fonte reduzida: **10 pt** (quando a fonte do texto principal for 12 pt);
- Espaçamento entrelinhas **simples (1,0 / 100%)**;
- Alinhamento **justificado**;
- **Supressão de aspas:** Não se utilizam aspas no início ou fim;
- Espaçamento anterior e posterior de 1 linha em branco (ou `spaceAbove: 6 PT`, `spaceBelow: 6 PT`).

#### Conversão Métrica
$$4,0\text{ cm} \times 28,3464567 = 113,3858\text{ pt} \approx \mathbf{113,386\text{ pt}}$$

#### Mapeamento de Payload na Google Docs API v1
Requer a composição de duas requisições atômicas:
1. `UpdateParagraphStyleRequest`: Ajusta recuo em bloco (`indentStart: 113.386 PT`), anula recuo de 1ª linha (`indentFirstLine: 0 PT`), ajusta entrelinhas para 100% e insere margens verticais de 6 pt.
2. `UpdateTextStyleRequest`: Reduz a fonte para 10 pt.

```json
[
  {
    "updateParagraphStyle": {
      "range": { "startIndex": 2500, "endIndex": 2850 },
      "paragraphStyle": {
        "namedStyleType": "NORMAL_TEXT",
        "alignment": "JUSTIFIED",
        "lineSpacing": 100,
        "indentStart": { "magnitude": 113.386, "unit": "PT" },
        "indentFirstLine": { "magnitude": 0, "unit": "PT" },
        "spaceAbove": { "magnitude": 6, "unit": "PT" },
        "spaceBelow": { "magnitude": 6, "unit": "PT" },
        "avoidWidowAndOrphan": true
      },
      "fields": "namedStyleType,alignment,lineSpacing,indentStart,indentFirstLine,spaceAbove,spaceBelow,avoidWidowAndOrphan"
    }
  },
  {
    "updateTextStyle": {
      "range": { "startIndex": 2500, "endIndex": 2850 },
      "textStyle": {
        "fontSize": { "magnitude": 10.0, "unit": "PT" },
        "bold": false
      },
      "fields": "fontSize,bold"
    }
  }
]
```

---

### Tópico 5: Capa Institucional e Folha de Rosto (NBR 14724)

#### Norma Regulamentar
- **Capa (Elemento Pré-Textual Obrigatório - Página 1):**
  - Nome da Instituição: centralizado, caixa alta, negrito, 12 pt.
  - Nome do Autor: centralizado, caixa alta, negrito, 12 pt.
  - Título do Trabalho: centralizado, caixa alta, negrito, 12 pt.
  - Local e Ano: rodapé, centralizado, 12 pt.
  - Todos os elementos sem recuo (`indentStart: 0`, `indentFirstLine: 0`).
- **Folha de Rosto (Página 2):**
  - Nome do Autor: centralizado, caixa alta, SEM negrito, 12 pt.
  - Título do Trabalho: centralizado, caixa alta, COM negrito, 12 pt.
  - **Preâmbulo (Nota Institucional):** Recuo de **8,0 cm** a partir da margem esquerda, tamanho **10 pt**, espaçamento simples (1,0), alinhamento **justificado**, sem negrito.
  - Orientador e Coorientador: recuados a 8,0 cm, tamanho 10 pt, espaçamento simples.
  - Local e Ano: rodapé, centralizado, SEM negrito, 12 pt.

#### Conversão Métrica do Preâmbulo (8,0 cm)
$$8,0\text{ cm} \times 28,3464567 = 226,77165\text{ pt} \approx \mathbf{226,772\text{ pt}}$$

#### Mapeamento de Payload na Google Docs API v1
No Google Docs, o preâmbulo é formatado definindo tanto `indentStart` quanto `indentFirstLine` em $226,772\text{ pt}$ para que todas as linhas, inclusive a primeira, alinhem-se à direita a partir do meio da página:

```json
{
  "updateParagraphStyle": {
    "range": { "startIndex": 650, "endIndex": 980 },
    "paragraphStyle": {
      "namedStyleType": "NORMAL_TEXT",
      "alignment": "JUSTIFIED",
      "lineSpacing": 100,
      "indentStart": { "magnitude": 226.772, "unit": "PT" },
      "indentFirstLine": { "magnitude": 226.772, "unit": "PT" },
      "spaceAbove": { "magnitude": 14.17, "unit": "PT" },
      "spaceBelow": { "magnitude": 0, "unit": "PT" },
      "avoidWidowAndOrphan": true
    },
    "fields": "namedStyleType,alignment,lineSpacing,indentStart,indentFirstLine,spaceAbove,spaceBelow,avoidWidowAndOrphan"
  }
}
```

---

### Tópico 6: Numeração de Páginas no Canto Superior Direito e Gerenciamento de Seções

#### Norma Regulamentar (NBR 14724, item 5.3)
1. **Contagem de Páginas:** Todas as páginas do trabalho a partir da folha de rosto são **contadas** sequencialmente. A capa não é contada.
2. **Exibição do Número:** O número da página só deve ser **impresso a partir da primeira folha da parte textual** (Introdução / Capítulo 1).
3. **Localização:** Canto superior direito da folha, a $2,0\text{ cm}$ da borda superior e $2,0\text{ cm}$ da borda direita.
4. **Supressão:** As páginas pré-textuais (Folha de Rosto, Resumo, Abstract, Listas, Sumário) **não recebem número visível**.

#### Investigação de Capacidades e Limitações na Google Docs API v1
1. **Criação de Quebra de Seção:**
   - Requisição: `InsertSectionBreakRequest`
   - Parâmetro: `sectionType: "NEXT_PAGE"`
   - Posicionamento: inserido exatamente antes do título do primeiro capítulo (ex: `1 INTRODUÇÃO`).
   - Efeito: Divide o documento em **Seção 1 (Pré-Textual)** e **Seção 2 (Textual/Pós-Textual)**.

2. **Criação e Isolamento de Cabeçalho por Seção:**
   - Requisição: `CreateHeaderRequest`
   - Parâmetros:
     - `type: "DEFAULT"`
     - `sectionBreakLocation: { "index": sectionBreakIndex }`
   - Regra Autoritativa da API Docs:
     > *"The new header is applied to the SectionStyle at the location of the SectionBreak if specified, otherwise it is applied to the DocumentStyle."*
     > *"If unset, the value inherits from the previous SectionBreak's SectionStyle. If the value is unset in the first SectionBreak, it inherits from DocumentStyle's default_header_id."*
   - **Mecanismo de Desvinculação:**
     - Se o documento tiver um cabeçalho global indesejado na Seção 1, executa-se `DeleteHeaderRequest` para deletá-lo.
     - Cria-se o cabeçalho **exclusivamente na Seção 2** fornecendo o `sectionBreakLocation` da quebra de seção.
     - Como a Seção 1 não possui cabeçalho definido e não herda de nenhuma seção anterior, **as páginas pré-textuais ficam completamente sem cabeçalho e sem numeração**.
     - A Seção 2 ganha um `headerId` próprio, isolado e ativo.

3. **Início da Numeração (`sectionStyle.pageNumberStart`):**
   - Na Seção 2, se não for especificado `pageNumberStart`, a contagem de páginas **continua automaticamente da Seção 1**. Como a Seção 1 possui as páginas pré-textuais, a primeira página da Seção 2 receberá exatamente o número sequencial correto acumulado (ex: página 7 ou 8).
   - Caso seja necessário forçar um número específico (desconsiderando a capa), pode-se atualizar `sectionStyle.pageNumberStart` via `UpdateSectionStyleRequest`.

4. **Limitação Inerente da API v1 a Respeito de AutoText (`PAGE_NUMBER`):**
   - **Constatação Crítica:** A Google Docs API v1 **NÃO POSSUI NENHUMA REQUISIÇÃO PARA INSERIR AUTOTEXT (`PAGE_NUMBER`)** no lote de `batchUpdate`! Não existe `InsertAutoTextRequest` nem `InsertPageNumberRequest`.
   - O tipo `AutoText` com `type: "PAGE_NUMBER"` existe apenas como entidade de leitura (`StructuralElement.Paragraph.Elements[].AutoText`).
   - **Solução Arquitetural de Engenharia:**
     - Para documentos criados dinamicamente via API, o documento deve ser inicializado a partir de um modelo (*template*) do Google Drive que já possua o token nativo de número de página pré-configurado no cabeçalho padrão, ou
     - O MCP deve gerenciar as quebras de seção e isolamento de cabeçalhos (`useFirstPageHeaderFooter` e `SectionBreak`), mantendo o campo de número de página existente na Seção 2 e eliminando-o da Seção 1.

---

### Tópico 7: Sumário Automático (Nativo vs. Gerador Estruturado ABNT)

#### Norma Regulamentar (NBR 6027:2012)
O sumário exige:
- Reprodução exata dos títulos e subtítulos nas mesmas tipografias da NBR 6024;
- Linha pontilhada (*leader dots*) unindo o término do título ao número da página;
- Número da página rigorosamente alinhado à margem direita da folha.

#### Investigação de Capacidades da Google Docs API v1
1. **Inexistência de Requisição de Inserção de TOC Nativo:**
   - A struct `TableOfContents` existe no SDK Go apenas como elemento estrutural de saída (`StructuralElement.TableOfContents`).
   - **A API v1 não disponibiliza `InsertTableOfContentsRequest` nem `UpdateTableOfContentsRequest`**. Não é possível inserir um sumário nativo do Google Docs puramente via chamadas REST da API v1!
   - Além disso, o sumário nativo da UI do Google Docs utiliza estilos de parágrafo rígidos que não atendem às exigências tipográficas da ABNT NBR 6024 (ex: seções primárias em caixa alta e negrito, seções secundárias em caixa alta regular, etc.).

2. **Inexistência de Leader Dots no Objeto `TabStop` da API v1:**
   - A struct `TabStop` da Docs API v1 suporta apenas `Alignment: "START" | "CENTER" | "END"` e `Offset: Dimension`. **Não há propriedade para preenchimento com pontos (leader dots)**.
   - Mais crítico ainda: na struct `ParagraphStyle`, o campo `TabStops []*TabStop` está documentado oficialmente como **`This property is read-only`**. Não é permitido configurar paradas de tabulação via `UpdateParagraphStyleRequest`.

3. **Arquitetura Recomendada: Tabela Invisível ABNT de 2 Colunas:**
   - A solução de engenharia mais robusta e esteticamente precisa na API v1 consiste na geração de uma **Tabela Nativa de 2 colunas sem bordas laterais nem verticais**:
     - Coluna 1 (Largura: $14,0\text{ cm} \approx 396,85\text{ pt}$): Título da seção formatado em conformidade com o nível NBR 6024 correspondente.
     - Coluna 2 (Largura: $2,0\text{ cm} \approx 56,69\text{ pt}$): Número da página alinhado à direita (`Alignment: "END"`).
     - Largura total combinada: $16,0\text{ cm}$ (exatamente a mancha gráfica útil A4).
     - Sem bordas (`borderLeft, borderRight, borderTop, borderBottom` com largura 0 pt).
     - Pontilhamento intermediário gerado no texto da Coluna 1 quando desejado.

---

### Tópico 8: Remoção Programática de Quebras de Linha em Excesso

#### Problema Diagnosticado
No documento de teste e em documentos convertidos de LaTeX/Markdown, é comum a proliferação de:
1. Parágrafos vazios consecutivos (`\n\n\n`) inseridos manualmente para empurrar conteúdo para a página seguinte após a Folha de Rosto;
2. Linhas em branco redundantes entre parágrafos de texto no corpo do trabalho.

#### Restrições Formais de `DeleteContentRangeRequest` na Google Docs API v1
A documentação autoritativa da API do Google Docs impõe restrições intransponíveis que disparam erro HTTP `400 Bad Request: Invalid range`:
1. **Delimitador Final Inamovível:** O último caractere do corpo do documento é um `\n` (`body.EndIndex - 1`). Qualquer tentativa de deletar esse caractere ou ultrapassar `body.EndIndex` é sumariamente rejeitada.
2. **Deleção de Delimitadores de Parágrafo (`\n`):** Deletar o caractere `\n` que separa dois parágrafos faz com que o parágrafo anterior e o posterior se fundam, unificando seus estilos de parágrafo.
3. **Parágrafo Antecessor de Tabela:** É proibido deletar o `\n` imediatamente anterior a uma `Table`, `TableOfContents` ou `SectionBreak`. Toda tabela no Google Docs exige a existência de um parágrafo que a preceda. Tentar deletar esse parágrafo dispara erro 400.
4. **Pares Substitutos UTF-16 (*Surrogate Pairs*):** Deletar apenas 1 unidade de código de um caractere de 4 bytes (como emojis ou símbolos matemáticos compostos) resulta em erro 400.
5. **Deriva Posicional (*Index Drift*):** Deletar múltiplos intervalos em ordem crescente (do início para o fim) desloca as posições de todos os elementos posteriores, corrompendo os índices previamente calculados e causando falhas catastróficas em lote.

#### Algoritmo Canônico de Limpeza com Blindagem
Para higienizar linhas em branco sem violar a integridade estrutural:
1. **Inspeção de Elementos:** Varrer `doc.Body.Content` e mapear os índices de parágrafos puramente vazios (`trimmedText == ""`).
2. **Proteção Pré-Tabela:** Se o parágrafo vazio anteceder imediatamente uma `Table` (`doc.Body.Content[i+1].Table != nil`), **ignorar e não deletar**.
3. **Proteção de Quebras Pré-Textuais:** Preservar quebras estruturais canônicas da Capa e Folha de Rosto, acionando a remoção de vazios apenas a partir do início da parte textual ou após a Folha de Rosto.
4. **Execução em Ordem Reversa (*Reverse Index Order Rule*):** Submeter as deleções partindo do maior índice para o menor índice. Cada deleção afeta apenas índices posteriores que já foram processados, mantendo os offsets anteriores matematicamente válidos.
5. **Clamping Defensivo:** Garantir que `endIndex <= bodyEndIndex - 1`.

---

## 3. Respostas Fundamentadas ao Diagnóstico do Documento de Teste (`1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`)

### Diagnóstico 1: Por que o linter e o formatador não identificaram nem corrigiram variações indesejadas de entrelinhas (ex.: 1,15 vs 1,5)?
- **Causa Raiz no Linter (`internal/gdocs/linter.go`):**
  A struct `LintReport` e a função `LintDocument` sequer possuem campos ou lógica para inspecionar `paragraph.ParagraphStyle.LineSpacing`. O linter verifica `FontSize`, `FontFamily`, `TextColor`, `BoldCascades` e `LatexLeaks`, mas ignora por completo a propriedade de espaçamento entrelinhas! Portanto, parágrafos com entrelinhas padrão do Google Docs (1,15 / 115%) ou qualquer outra variação passam totalmente despercebidos como falsos negativos.
- **Causa Raiz no Formatador (`internal/gdocs/abnt_writer.go`):**
  O formatador aplica `ABNTLineSpacing: 150` quando acionado, porém se a ferramenta de inspeção e linting não sinaliza a anomalia, o usuário ou o agente LLM nunca é informado sobre a discrepância e não dispara a correção localizada.

### Diagnóstico 2: Por que as quebras de linha espúrias após a folha de rosto não foram saneadas?
- **Causa Raiz no Linter (`internal/gdocs/linter.go:730-745`):**
  O linter adota uma máquina de estados: `statePreTextual`, `stateTextual`, `statePostTextual`. Na linha 734, há a guarda:
  ```go
  if trimmed == "" {
      if currentState != statePreTextual { // Linhas vazias na fase pré-textual são ignoradas!
          consecutiveEmptyParas++
          ...
      }
  }
  ```
  A transição para `stateTextual` só ocorre quando o linter encontra o primeiro capítulo numerado (`1 INTRODUÇÃO`). Como a região logo após a Folha de Rosto e antes da Introdução (Resumo, Abstract, etc.) ainda está sob `currentState == statePreTextual`, **todas as linhas em branco consecutivas e quebras espúrias após a folha de rosto são ativamente ignoradas pelo linter**.
- **Causa Raiz no Saneador (`internal/gdocs/abnt_writer.go:180-210`):**
  Na função `HealDocTextAndBlanks`, a detecção do ponto de corte para colapso de linhas em branco depende de encontrar duas vezes o ano `"2026"` ou a palavra `"RESUMO"` de forma estrita. Se o texto contiver quebras de parágrafo adicionais ou variações no ano, o índice de início `targetStartIndex` falha em se posicionar corretamente, abortando a deleção das quebras pré-textuais.

### Diagnóstico 3: O que falta nas ferramentas de leitura e inspeção (`read_doc_content`, `inspect_doc_layout`) para expor metadados suficientes para modelos de IA?
- **Omissão Crítica em `inspect_doc_layout` (`internal/gdocs/inspector.go`):**
  A struct `ElementLayoutDetail` expõe `Alignment`, `IndentStartPt`, `IndentFirstLinePt`, `SpaceAbovePt`, `SpaceBelowPt`, `FontSizePt`, `FontFamily`, `Bold`, `Italic`, mas **NÃO expõe `LineSpacing` (ou `LineSpacingPercent`)**! Um modelo de IA que invoca `inspect_doc_layout` é incapaz de enxergar o espaçamento entrelinhas dos parágrafos inspecionados.
- **Omissão em `read_doc_content`:**
  A ferramenta converte o documento para Markdown plano. Markdown não possui semântica para representar espaçamento entrelinhas (1,15 vs 1,5), recuo de primeira linha (1,25 cm), margens de página ou preâmbulo institucional. Dessa forma, a IA operando apenas com `read_doc_content` não tem visibilidade da conformidade ABNT do documento.

---

## 4. Tabelas Formais de Mineração

### Features Discovered

| # | Category | Feature | Description | Inputs | Outputs | Error Behavior | Discovered Via |
|---|---|---|---|---|---|---|---|
| 1 | Page Margins | `UpdateDocumentStyleRequest` | Aplica margens A4 ABNT (Sup 3cm, Esq 3cm, Inf 2cm, Dir 2cm) e dimensões A4 ($21 \times 29,7\text{ cm}$) | `DocumentStyle` com `marginTop: 85.039pt`, `marginLeft: 85.039pt`, `marginBottom: 56.693pt`, `marginRight: 56.693pt`, `pageSize: 595.276x841.890pt`, `fields` | Atualiza o estilo do documento e limpa margens conflitantes em todas as seções | Erro 400 se `fields` for vazio ou se valores forem negativos | `docs/pesquisa-api.md:7014` |
| 2 | Paragraph Spacing | `UpdateParagraphStyleRequest.lineSpacing` | Define espaçamento entrelinhas percentual (150 = 1,5 entrelinhas ABNT; 100 = 1,0 simples) | `ParagraphStyle.LineSpacing: 150` ou `100`, `fields: "lineSpacing"` e `Range` | Aplica a altura de entrelinhas nos parágrafos do intervalo | Erro 400 se `lineSpacing` for menor que 25% | `docs/pesquisa-api.md:4492`, `docs/pesquisa-api.md:7195` |
| 3 | Paragraph Alignment | `UpdateParagraphStyleRequest.alignment` | Ajusta alinhamento do texto para `JUSTIFIED` no corpo e referências à esquerda (`START`) | `ParagraphStyle.Alignment: "JUSTIFIED"` ou `"START"`, `fields: "alignment"` e `Range` | Alinha o texto das linhas do parágrafo | Erro 400 se o valor de alinhamento não pertencer ao enum | `docs/pesquisa-api.md:4418` |
| 4 | First Line Indent | `UpdateParagraphStyleRequest.indentFirstLine` | Define recuo de primeira linha de 1,25 cm ($35,433\text{ pt}$) em parágrafos do corpo textual | `ParagraphStyle.IndentFirstLine: 35.433pt`, `IndentStart: 0pt`, `fields: "indentFirstLine,indentStart"` | Primeira linha recuada em 1,25 cm, demais linhas na margem zero | Deslocamento cumulativo distorcido se `indentStart` não for zerado | `docs/pesquisa-api.md:4477`, NBR 14724 |
| 5 | Block Indent (Quote) | `UpdateParagraphStyleRequest.indentStart` | Define recuo em bloco de 4,0 cm ($113,386\text{ pt}$) para citações diretas longas | `ParagraphStyle.IndentStart: 113.386pt`, `IndentFirstLine: 0pt`, `lineSpacing: 100`, `fields` | Todo o bloco de parágrafos é recuado a 4 cm com espaçamento simples | Erro 400 se `Range` cruzar limites de tabelas ou seções | `docs/pesquisa-api.md:4480`, NBR 10520 |
| 6 | Pre-textual Note Indent | `UpdateParagraphStyleRequest.indentStart` (8.0 cm) | Define recuo de 8,0 cm ($226,772\text{ pt}$) para preâmbulo/nota institucional na folha de rosto | `ParagraphStyle.IndentStart: 226.772pt`, `IndentFirstLine: 226.772pt`, `lineSpacing: 100`, `fields` | Preâmbulo justificado do meio da página até a margem direita | Erro 400 se `fields` for omitido | `docs/pesquisa-api.md:4480`, NBR 14724 |
| 7 | Section Separation | `InsertSectionBreakRequest` | Cria uma quebra de seção `NEXT_PAGE` antes da Introdução para separar elementos pré-textuais | `sectionType: "NEXT_PAGE"`, `location: { index: introIndex }` | Divide o documento em Seção 1 (pré-textual) e Seção 2 (textual) | Erro 400 se inserido no interior de tabelas ou cabeçalhos | `docs/pesquisa-api.md:3528` |
| 8 | Header Creation | `CreateHeaderRequest` | Cria cabeçalho exclusivo associado ao `SectionBreakLocation` da Seção 2 | `type: "DEFAULT"`, `sectionBreakLocation: { index: sectionBreakIndex }` | Retorna `headerId` e associa cabeçalho isolado à Seção 2 | Erro 400 se já existir cabeçalho do tipo `DEFAULT` na seção | `docs/pesquisa-api.md:1156` |
| 9 | Section Page Numbering | `UpdateSectionStyleRequest.pageNumberStart` | Configura número de início de contagem de páginas para a seção ou mantém herança contínua | `SectionStyle.PageNumberStart: int64`, `fields: "pageNumberStart"`, `Range` | Ajusta o início da numeração da seção | Erro 400 se o número for menor ou igual a zero | `docs/pesquisa-api.md:5604`, `docs/pesquisa-api.md:7225` |
| 10 | Header Suppression | `DeleteHeaderRequest` | Deleta cabeçalho global do `DocumentStyle` ou Seção 1 para suprimir números nas páginas pré-textuais | `headerId: string` | Remove o cabeçalho, deixando a Seção 1 sem numeração visível | Erro 404 se `headerId` não existir ou 400 se já deletado | `docs/pesquisa-api.md:1679` |
| 11 | Content Deletion | `DeleteContentRangeRequest` | Remove faixas de caracteres, quebras duras redundantes ou linhas em branco excessivas | `Range: { startIndex, endIndex }` | Remove o conteúdo compreendido no intervalo | Erro 400 se atingir o último `\n` do body ou o parágrafo antes de tabela | `docs/pesquisa-api.md:1583` |
| 12 | Table Creation (TOC) | `InsertTableRequest` | Cria tabela de 2 colunas para geração do Sumário ABNT com alinhamento rigoroso à direita | `rows: N`, `columns: 2`, `location: { index: sumarioIndex }` | Insere a estrutura de tabela vazia para o sumário estruturado | Erro 400 se o índice não estiver dentro de um parágrafo existente | `docs/pesquisa-api.md:3554` |

---

### Edge Cases

| # | Feature | Input | Observed Behavior |
|---|---|---|---|
| 1 | `DeleteContentRangeRequest` | `endIndex == body.EndIndex` (tentativa de deletar o último parágrafo do documento) | API rejeita com erro `400 Bad Request: The content range must not include the final newline.`. Requer clamping defensivo `endIndex = bodyEndIndex - 1`. |
| 2 | `DeleteContentRangeRequest` | Deletar o parágrafo vazio imediatamente anterior a uma `Table` (`el[i+1].Table != nil`) | API rejeita com erro `400 Bad Request: Deleting the newline character before a Table without deleting the element is not allowed.`. Requer guarda para preservar parágrafo pré-tabela. |
| 3 | `DeleteContentRangeRequest` | Deletar múltiplos intervalos em ordem crescente (`i = 0; i < N; i++`) | O primeiro delete encurta o documento e invalida todos os `startIndex`/`endIndex` seguintes. Requisições posteriores falham com `400 Bad Request: Invalid range` ou deletam texto incorreto. Requer execução em ordem decrescente (*Reverse Index Order*). |
| 4 | `UpdateParagraphStyleRequest` | Aplicar `lineSpacing: 150` sem fornecer máscara `fields` (`fields: ""`) | A requisição é silenciosamente ignorada pela Docs API ou rejeitada com erro 400. `fields` deve conter `"lineSpacing"`. |
| 5 | `UpdateParagraphStyleRequest` | Aplicar `indentFirstLine: 35.433 pt` sem resetar `indentStart` em parágrafo originado de lista | O parágrafo fica com recuo duplo: recuo da lista (ex: 36 pt) + 35,433 pt = 71,433 pt de recuo visual. Requer forçar `indentStart: 0 pt` na mesma requisição. |
| 6 | `CreateHeaderRequest` | Chamar `CreateHeaderRequest` com `type: "DEFAULT"` em seção que já possui cabeçalho `DEFAULT` | API retorna `400 Bad Request: A header of the specified type already exists.`. É mandatório verificar se já existe ou invocar `DeleteHeaderRequest` antes. |
| 7 | `InsertSectionBreakRequest` | Inserir quebra de seção em índice que coincide com o interior de uma tabela ou célula | API rejeita com erro `400 Bad Request: Section breaks can only be inserted inside the body of an existing paragraph.`. O índice deve pertencer a um parágrafo do corpo principal. |
| 8 | `UpdateDocumentStyleRequest` | Aplicar margens em documento configurado no modo `DocumentMode: PAGELESS` | A API aceita a requisição, porém o Google Docs não renderiza margens nem numeração de páginas em modo sem páginas. O documento deve estar em modo paginado (`PAGES`). |
| 9 | `AutoText` | Tentar injetar `AutoText` (`PAGE_NUMBER`) diretamente via requisição `batchUpdate` | Não existe nenhuma struct de requisição `InsertAutoTextRequest` na API v1. O payload é rejeitado no parsing de JSON-RPC / Protobuf como campo desconhecido. |
| 10 | `UpdateParagraphStyleRequest` | Tentar modificar paradas de tabulação via campo `tabStops` | A API ignora ou rejeita a escrita porque o campo `TabStops` na struct `ParagraphStyle` é marcado oficialmente como `read-only`. |

---

## 5. Diretrizes para o Plano de Engenharia (/plan)

Com base nas especificações mineradas, as seguintes diretrizes são mandatórias para a implementação subsequente no código Go:
1. **Atualização do `linter.go`:**
   - Adicionar o campo `LineSpacingAnomalies []LineSpacingAnomalyItem` em `LintReport`.
   - Adicionar auditoria sistemática de `p.ParagraphStyle.LineSpacing`: sinalizar qualquer parágrafo do corpo com `LineSpacing != 150` e qualquer legenda/fonte/citação longa/referência com `LineSpacing != 100`.
   - Ajustar a máquina de estados de seções no linter: permitir que linhas em branco consecutivas espúrias após a Folha de Rosto sejam auditadas e reportadas, mesmo antes do primeiro capítulo numerado.
2. **Atualização do `inspector.go`:**
   - Adicionar `LineSpacing float64` e `SpacingMode string` na struct `ElementLayoutDetail`.
   - Expor as propriedades de entrelinhas na resposta JSON do `inspect_doc_layout` para permitir diagnóstico completo por modelos de IA.
3. **Mecanismo de Proteção de Citações Longas no `abnt_writer.go`:**
   - Ao executar `FormatDocumentABNT` sobre todo o documento, o formatador não deve achatar citações longas para `NORMAL_TEXT` padrão (1.25cm / 1.5). Deve detectar parágrafos com recuo `>= 90 pt` ou com extensão $> 3$ linhas acompanhados de fonte 10 pt e preservar seu recuo de 4,0 cm e entrelinhas 1,0 simples.
4. **Gerador de Sumário Estruturado ABNT:**
   - Utilizar o padrão de Tabela Nativa de 2 colunas sem bordas para garantir alinhamento perfeito dos números de página à direita sem depender de paradas de tabulação ou recursos não suportados na API v1.
