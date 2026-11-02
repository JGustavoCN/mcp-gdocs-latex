# Relatório de Análise Adversarial e Testes de Estresse Conceitual: Plano de Engenharia ABNT

**Agente:** `challenger_plan_1` (Critic & Specialist - Adversarial Challenger)  
**Data:** 01 de Novembro de 2026  
**Documentos Auditados:**
- `docs/plano_engenharia_abnt.md` (v2.0.0)
- `PROJECT.md` (v1.3.0)
- `.agents/teamwork/worker_plan_1/handoff.md`
- `.agents/teamwork/ORIGINAL_REQUEST.md`

**Veredicto Oficial:** **REQUEST_CHANGES (Alterações Críticas Necessárias)**  
**Classificação Geral de Risco:** **HIGH (Alto Risco de Regressões e Falha de Implementação sem Ajustes)**

---

## 1. Resumo Executivo da Auditoria Adversarial

O plano de engenharia elaborado pelo `worker_plan_1` em `docs/plano_engenharia_abnt.md` e os incrementos documentados em `PROJECT.md` apresentam excelente profundidade teórica na conversão de pontos DTP, mapeamento de payloads JSON da Google Docs API v1 e identificação de causas raízes no código atual (`linter.go`, `abnt_writer.go`, `inspector.go`, `reader.go`, `cover_generator.go`).

Contudo, sob a ótica estritamente adversarial e empírica da Google Docs API e das normas regulamentares da ABNT, a análise identificou **6 vulnerabilidades de alto impacto**, incluindo:
1. **Risco Crítico de Destruição de Quebras de Página Reais (`PageBreak`):** A regra de saneamento de linhas em branco `consecutiveBlanks >= 1` não discrimina elementos estruturais sem texto, correndo o risco iminente de excluir o próprio parágrafo que abriga a quebra de página da Folha de Rosto ou entre capítulos.
2. **Lacuna Arquitetural na Origem dos Números de Página do Sumário:** A Google Docs API REST v1 **não fornece a página em que um parágrafo está renderizado**. O plano propõe preencher a Coluna 2 da tabela do Sumário com o número de página, mas omite completamente a fonte e o mecanismo para calcular ou obter esses números.
3. **Mutilação de Citações Diretas Longas com Múltiplos Parágrafos:** A heurística de detecção analisa parágrafos isoladamente sem estado de continuidade, rebaixando parágrafos curtos subsequentes de uma mesma citação longa para o formato de corpo de texto (`NORMAL_TEXT` 1,25 cm / 1,5), além de aplicar espaçamento vertical excessivo entre eles.
4. **Cegueira e Falsos Positivos em Seções Pós-Textuais (Apêndices e Anexos):** A máquina de estados e o escopo do documento limitam o pós-textual às Referências. Apêndices e Anexos (com identificadores alfabéticos NBR 6024) não são reconhecidos, gerando alertas indevidos de ilustrações fora de lugar e exclusão no sumário.
5. **Corrupção de Listas dentro de Células de Tabelas:** O zeramento irrestrito de `indentStart` e a deleção compulsória de bullets (`DeleteParagraphBulletsRequest`) em parágrafos de tabelas destrói listas internas e sobrepõe marcadores.
6. **Inconsistência de Governança de Ferramentas no `PROJECT.md`:** Divergência entre as ferramentas registradas no código (`tools.go`), a consolidação de `apply_abnt_page_margins` em `format_document_abnt` e o catálogo formal de 17 ferramentas listado no `PROJECT.md`.

---

## 2. Análise Detalhada por Eixo de Desafio

---

### Eixo 1: Cenários Extremos em ABNT

#### 1.1 Apêndices e Anexos com Numeração Própria (ABNT NBR 14724:2011 item 4.3 e NBR 6024:2012 item 3.8)
- **Cenário Adversarial:**
  O documento acadêmico possui seções pós-textuais além das Referências: `APÊNDICE A – Roteiro de Entrevista` e `ANEXO A – Parecer do Comitê de Ética`. Tais seções não utilizam algarismos arábicos (não são "8 APÊNDICE"), mas letras maiúsculas ordenadas, travessão e título em caixa alta.
- **Falhas Encontradas no Plano e no Codebase:**
  1. *Definição Redutiva de Seção 2 no Plano (`plano_engenharia_abnt.md:270`):*  
     O plano define: `Seção 2 (Textual e Pós-Textual): Introdução até as Referências`. Essa formulação esquece formalmente Apêndices e Anexos, o que induz a erros no particionamento de cabeçalho e contagem de páginas. As páginas de apêndices e anexos **devem ser contadas e numeradas continuamente** até o encerramento do volume.
  2. *Incompetência de Reconhecimento no Código Go (`structure_normalizer.go:56`):*  
     `regexPrimaryHeading` exige `[1-9][0-9]*\s+` ou strings fixas (`INTRODUÇÃO|...|REFERÊNCIAS|RESUMO|ABSTRACT`). Ele **ignora completamente** títulos iniciados por `APÊNDICE` e `ANEXO`.
  3. *Geração de Falsos Positivos no Linter (`linter.go:544-550`):*  
     Se um Apêndice ou Anexo contiver uma tabela ou figura (ex.: formulário com `Tabela A1`), o linter categoriza como pós-textual indiscriminado e dispara o falso positivo:  
     `"Legenda de ilustração '...' posicionada indevidamente dentro da seção de REFERÊNCIAS (NBR 6023)"`.
  4. *Ausência no Sumário ABNT (NBR 6027:2012):*  
     Apêndices e Anexos são de inclusão obrigatória no sumário se presentes no documento. O algoritmo de geração de tabela do sumário precisa rastreá-los explicitamente.
- **Mitigação Exigida:**
  - Estender o escopo da Seção 2 para: "Da Introdução até a última página do último Anexo/Apêndice".
  - Incluir padrões `^APÊNDICE\s+[A-Z]\s+[–\-]\s+` e `^ANEXO\s+[A-Z]\s+[–\-]\s+` nos parsers de títulos primários.
  - Ajustar a máquina de estados do `linter.go` criando estados específicos `stateReferences` e `stateAppendixAnnex`, desativando a proibição de figuras/tabelas nos apêndices e anexos.

---

#### 1.2 Citações Diretas Longas com Múltiplos Parágrafos (ABNT NBR 10520:2023)
- **Cenário Adversarial:**
  O autor transcreve uma citação direta longa de 2 ou 3 parágrafos contíguos (ex.: múltiplos parágrafos de um artigo de lei ou transcrição textual com troca de interlocutor). O primeiro parágrafo possui 4 linhas, mas o segundo parágrafo possui apenas 2 linhas.
- **Falhas Encontradas no Plano:**
  1. *Heurística Isolada e Mutilação de Parágrafos Curtos (`plano_engenharia_abnt.md:628-634`):*  
     A regra proposta classifica cada parágrafo de forma atômica e independente:
     `existingIndent >= 90.0 PT OU (extensão > 3 linhas E fonte <= 10.5 pt)`.
     Em um documento recém-gerado, importado de Markdown/LaTeX ou ainda desformatado (onde `existingIndent == 0` e a fonte ainda é 12 pt ou 11 pt):
     - O parágrafo 1 tem > 3 linhas.
     - O parágrafo 2 tem apenas 2 linhas e fonte 12 pt.
     **Resultado:** A heurística falha para o parágrafo 2! O parágrafo 2 é classificado como "Corpo de Texto Regular", recebendo `NORMAL_TEXT`, recuo de 1,25 cm, entrelinhas 1,5 e fonte 12 pt. A citação longa é partida ao meio e desfigurada!
  2. *Espaçamento Vertical Excessivo Interno:*  
     O payload proposto para citação longa (linhas 199-200) aplica:  
     `spaceAbove: 6.0 PT` e `spaceBelow: 6.0 PT`.  
     Se aplicado a parágrafos consecutivos da mesma citação longa, haverá um espaçamento de `6 + 6 = 12 pt` entre eles. Conforme a NBR 10520, o espaçamento entre parágrafos internos de uma citação direta deve ser simples (1,0 contíguo), aplicando-se o espaçamento extra apenas antes do primeiro parágrafo e depois do último.
  3. *Inexistência de Métrica Computacional para "> 3 linhas":*  
     A Google Docs API v1 opera sobre streams de texto e **não expõe o número de linhas renderizadas**. O plano menciona "extensão > 3 linhas" sem definir a regra algorítmica em Go. Na mancha gráfica padrão ABNT (16 cm) com recuo de 4 cm (largura útil = 12 cm), uma linha em fonte 10 pt comporta em média 75 a 85 caracteres. Três linhas correspondem a aproximadamente 240-250 caracteres. O plano precisa formalizar essa métrica para evitar ambiguidade de implementação.
- **Mitigação Exigida:**
  - Implementar máquina de estados com **memória de contexto de citação contínua**: se o parágrafo anterior é uma citação longa e o parágrafo atual não é título nem elemento pré-textual e possui indicação contextual (ex.: fechamento de aspas ou citação parentética de autor-data ao final), ele herda a classe de citação longa.
  - Calibrar o espaçamento vertical: `spaceAbove: 6 pt` somente no primeiro parágrafo da citação e `spaceBelow: 6 pt` somente no último parágrafo da citação (parágrafos intermediários com `spaceAbove: 0` e `spaceBelow: 0`).
  - Definir limiar determinístico em caracteres UTF-16: `utf8.RuneCountInString(text) >= 240` para a heurística textual na ausência de recuo prévio.

---

#### 1.3 Listas com Marcadores dentro de Células de Tabelas
- **Cenário Adversarial:**
  Tabelas comparativas e matrizes de requisitos contêm células com listas pontuadas (bullets `•`, `-` ou marcadores nativos da Docs API).
- **Falhas Encontradas no Plano:**
  1. *Destruição de Recuos de Lista nas Células (`plano_engenharia_abnt.md:468-470`):*  
     O plano prescreve:  
     `UpdateParagraphStyleRequest: lineSpacing: 100, indentStart: 0, indentFirstLine: 0`.  
     Em parágrafos com lista ativa (`p.Bullet != nil`), forçar `indentStart: 0` e `indentFirstLine: 0` elimina o recuo padrão do marcador do Google Docs (hanging indent), fazendo com que o glifo do marcador sobreponha o primeiro caractere do texto ou transborde para a margem da célula.
  2. *Deleção Compulsória de Bullets em `replace_block_abnt`:*  
     No componente `builder.go:140`, o método `BuildABNTReplaceSequence` emite incondicionalmente `DeleteParagraphBulletsRequest`. Se o usuário invocar `replace_block_abnt` para editar o texto de uma célula que contém lista, a lista é destruída.
  3. *Inadequação de Alinhamento Justificado em Tabelas:*  
     O plano propõe `alignment: "JUSTIFIED" ou "START"`. Em células de tabela, o alinhamento justificado cria rios de espaços em branco (*typographical rivers*) devido à pequena largura das colunas.
- **Mitigação Exigida:**
  - Na formatação de células de tabela, verificar se `p.Bullet != nil`: se houver lista, preservar o recuo de lista proporcional (`indentStart: 14.17 pt` / 0,5 cm, `indentFirstLine: -14.17 pt`) em vez de zerar.
  - Fixar como norma absoluta em células de tabela o alinhamento à esquerda (`START`).
  - Adicionar flag `preserve_bullets: bool` em `replace_block_abnt`.

---

### Eixo 2: Desafio do Saneamento de Quebras de Linha e Regras da Docs API

#### 2.1 Risco de Deleção Acidental de Quebras de Página Reais (`PageBreak`)
- **Cenário Adversarial:**
  A folha de rosto encerra-se com uma quebra de página formal (`PageBreak`). Na página seguinte, inicia-se o "RESUMO". O parágrafo que abriga a `PageBreak` é um elemento do tipo `Paragraph` que contém apenas `pel.PageBreak != nil` e nenhum `TextRun` com texto visível.
- **Mecanismo da Falha Crítica:**
  1. A Google Docs API v1 modela `PageBreak` como um filho `ParagraphElement` de um parágrafo (`Paragraph.Elements[*].PageBreak`).
  2. O código de varredura (ex.: `abnt_writer.go:334-342`) inspeciona o parágrafo fazendo:
     ```go
     txt := ""
     for _, pel := range el.Paragraph.Elements {
         if pel.TextRun != nil {
             txt += pel.TextRun.Content
         }
     }
     t := strings.TrimSpace(txt)
     if t == "" { ... }
     ```
  3. Para o parágrafo da `PageBreak`, `pel.TextRun` é `nil`. Logo, `t` é avaliado como string vazia `""`!
  4. Se a nova regra de saneamento do plano (`consecutiveBlanks >= 1` imediatamente após transição de página) for acionada:
     O parágrafo contendo a PRÓPRIA `PageBreak` é interpretado como uma "linha em branco espúria" e agendado para deleção (`toDeleteRanges = append(...)`)!
  5. Ao executar o `DeleteContentRangeRequest`, a `PageBreak` é EXCLUÍDA! A folha de rosto e o resumo colapsam na mesma folha!
  6. **Evidência no Codebase Atual:**  
     Este exato defeito já ocorre silenciosamente em `cover_generator.go:891-912`, onde o loop de limpeza reversa tenta deletar o parágrafo que contém a `PageBreak` da folha de rosto, disparando erro de violação de range na API que é então silenciado por `_, _ = docsSvc.Documents.BatchUpdate(...)`.
- **Mitigação Exigida:**
  O plano deve impor uma regra de **Guarda Absoluta de Elementos Estruturais Não-Textuais**:  
  Um parágrafo NUNCA pode ser considerado vazio para deleção se possuir qualquer um dos seguintes elementos:
  - `pel.PageBreak != nil`
  - `pel.HorizontalRule != nil`
  - `pel.InlineObjectElement != nil`
  - `pel.FootnoteReference != nil`
  - `pel.Person != nil`

---

#### 2.2 Proteção Estrita do Delimitador Pré-Tabela e do Final do Body (`body.EndIndex - 1`)
- **Cenário Adversarial:**
  Uma página termina com `PageBreak` e a página seguinte começa diretamente com uma `Table` nativa (muito comum em capítulos de resultados). O parágrafo imediatamente entre a quebra de página e a tabela é o único parágrafo existente antes da tabela.
- **Comportamento da Docs API v1:**
  A Google Docs API v1 **exige obrigatoriamente um parágrafo antes de qualquer tabela**. Tentar deletar esse parágrafo dispara erro HTTP `400 Bad Request: Table must be preceded by a paragraph.`
- **Mitigação Exigida no Plano:**
  A checagem pré-tabela (`content[i+1].Table != nil`) deve ser uma cláusula de exclusão obrigatória e prioritária em todas as rotinas de limpeza (`HealDocTextAndBlanks`, `cover_generator.go` e `replace_block_abnt`), mesmo que o parágrafo seja classificado como vazio pós-`PageBreak`.
  Além disso, a proteção de `body.EndIndex - 1` deve impedir que parágrafos finais com `el.EndIndex >= body.EndIndex` sejam agendados para exclusão integral.

---

### Eixo 3: Desafio do Sumário em Tabela Nativa de 2 Colunas

#### 3.1 A Grande Lacuna: De Onde Vêm os Números de Página?
- **O Fato Inegável da Google Docs API v1:**
  A Google Docs API REST v1 **NÃO POSSUI NENHUM ENDPOINT, MÉTODO OU CAMPO QUE RETORNE O NÚMERO DA PÁGINA ONDE UM PARÁGRAFO OU SEÇÃO ESTÁ RENDERIZADO**. O Google Docs é internamente um fluxo linear de caracteres, e a quebra de páginas física é executada exclusivamente pelo motor de renderização da interface Web ou no momento do export para PDF.
- **A Omissão Crítica no Plano de Engenharia:**
  O plano (`plano_engenharia_abnt.md:307-311`) especifica com clareza:
  - Coluna 1: Título completo formatado na NBR 6024;
  - Coluna 2: "Número da página correspondente, alinhado à direita".
  **Porém, o plano não responde:** *De onde a ferramenta Go vai tirar esse número da página?*
- **Solução Arquitetural Concreta que Deve Constatar no Plano:**
  O plano deve especificar formalmente uma das três estratégias viáveis (ou uma composição delas):
  1. **Estratégia Primária (LaTeX-Assisted / Canonical TCC):**  
     Como o projeto `mcp-gdocs-latex` opera acoplado ao TCC em `resources/tcc`, a ferramenta `get_tcc_latex_structure` ou o compilador LaTeX local gera os arquivos `main.toc` e `main.pdf`. A ferramenta extrai a paginação real do arquivo `.toc` gerado pelo LaTeX.
  2. **Estratégia Secundária (Drive PDF Render Pipeline):**  
     O servidor executa uma chamada em tempo de execução via `DriveService.Files.Export(docID, "application/pdf")` obtendo o PDF renderizado pelo Google Docs em memória, processa os fluxos de texto via parser de PDF em Go (ex.: `pdfcpu` ou parser nativo) e extrai a página física exata de cada título.
  3. **Estratégia Simbólica (Placeholder com Alerta):**  
     Caso o documento seja gerado desacompanhado do LaTeX ou do PDF, preencher a Coluna 2 com marcadores `[PÁG]` e instruir o usuário/LLM a executar a sincronização pós-compilação.
  *Deixar essa questão em aberto tornaria a ferramenta impossível de ser implementada com sucesso.*

---

#### 3.2 Comportamento em Múltiplas Páginas e Efeito Cascata (*Feedback Loop*)
- **Cenário Adversarial:**
  O sumário possui 45 seções e subdivisões, ocupando 2 páginas inteiras.
- **Análise do Comportamento no Google Docs:**
  - Tabelas nativas de muitas linhas quebram naturalmente entre páginas sem problemas, o que é positivo.
  - Como a quebra de seção `InsertSectionBreakRequest` está posicionada antes de `1 INTRODUÇÃO`, todas as páginas do sumário permanecem na Seção 1 (pré-textual), sem cabeçalho aparente, conforme manda a ABNT.
- **O Efeito Colateral Oculto:**
  Ao inserir um sumário de 2 páginas em vez de 1 página, o início de `1 INTRODUÇÃO` é empurrado em +1 página (ex.: da página 7 para a página 8). Todos os capítulos seguintes também são empurrados. Se a numeração foi extraída antes da criação da tabela, **todos os números da Coluna 2 estarão errados por um deslocamento de +1 ou +2 páginas**!
- **Mitigação Exigida:**
  O plano deve prever que a tabela do sumário seja inserida primeiro (mesmo que com placeholders), para que o documento atinja seu volume paginado final antes da resolução definitiva dos números de página.

---

#### 3.3 Hiperlinks de Navegação Interna (Heading Links / Bookmarks)
- **Norma e Boas Práticas (ABNT NBR 6027:2012 item 4.2):**
  Em suporte eletrônico, os títulos do sumário devem permitir navegação direta aos capítulos correspondentes.
- **Capacidade Real da Google Docs API v1:**
  A Google Docs API v1 permite vincular links internos diretamente ao `headingId` gerado para cada parágrafo com `namedStyleType` (`HEADING_1`, `HEADING_2`, etc.):
  ```json
  {
    "updateTextStyle": {
      "range": { "startIndex": startTitle, "endIndex": endTitle },
      "textStyle": {
        "link": { "headingId": "h.abcdef12345" },
        "underline": false,
        "foregroundColor": { "color": { "rgbColor": { "blue": 0, "green": 0, "red": 0 } } }
      },
      "fields": "link,underline,foregroundColor"
    }
  }
  ```
- **Omissão no Plano:**
  O plano não mencionou a criação desses links internos nem a blindagem contra o sublinhado azul padrão da Web, o que desvaloriza a experiência do usuário e a acessibilidade do documento acadêmico.
- **Mitigação Exigida:**
  Incluir no Passo de Preenchimento da Tabela do Sumário a extração do `p.ParagraphStyle.HeadingId` de cada título do documento e a aplicação do `link.headingId` estilizado em preto e sem sublinhado.

---

### Eixo 4: Desalinhamentos de Governança e Contratos em `PROJECT.md`

1. **Inconsistência no Catálogo de Ferramentas:**  
   Em `PROJECT.md:416-435`, a lista das 17 ferramentas lista `apply_abnt_page_margins` (FEAT-01). No entanto, o próprio código Go em `internal/mcp/tools.go:36` documenta:  
   `// (apply_abnt_page_margins foi consolidada em format_document_abnt para eliminar redundância)`.  
   Além disso, ferramentas essenciais como `format_document_abnt`, `lint_document_abnt`, `inspect_doc_layout` e `get_tcc_latex_structure` não constam na contagem oficial das 17 ferramentas do item 4.2 do `PROJECT.md`.
2. **Definição de Métricas de Teste em Go:**  
   O plano cita testes com mocks em `httptest.Server`, mas não estabelece a cobertura mínima nem valida a idempotência do formatador ao ser executado duas vezes seguidas no mesmo documento (*idempotency pass*).

---

## 3. Matriz Consolidada de Desafios Adversariais

| ID | Severidade | Eixo | Hipótese / Suposição Desafiada | Cenário de Falha / Modo de Colapso | Blast Radius | Ação Mitigadora Exigida no Plano |
|---|---|---|---|---|---|---|
| **ADV-01** | **CRITICAL** | Saneamento de Quebras | Parágrafos com `trimmedText == ""` após `PageBreak` são sempre linhas vazias deletáveis. | O parágrafo que abriga a própria `PageBreak` tem `TextRun == nil`. O algoritmo agenda sua deleção e remove a quebra de página da folha de rosto. | Colapso estrutural da capa e páginas pré-textuais. | Adicionar guarda explícita de `pel.PageBreak != nil` e outros elementos inline antes de agendar deleções. |
| **ADV-02** | **CRITICAL** | Sumário de 2 Colunas | Os números de página para a Coluna 2 da tabela de sumário estão prontamente disponíveis via Docs API. | A Google Docs API REST v1 não fornece números de página física de nenhum elemento estrutural. | Ferramenta inoperante ou preenchida com valores nulos/fictícios. | Definir pipeline canônico de resolução de páginas (via `.toc` do LaTeX compilado ou export PDF temporário). |
| **ADV-03** | **HIGH** | Citações Longas | Heurística isolada por parágrafo (`indent >= 90 pt` ou `> 3 linhas e 10 pt`) é suficiente. | Citação com múltiplos parágrafos onde o 2º parágrafo tem 2 linhas é rebaixado para `NORMAL_TEXT` (1,25 cm / 1,5), fatiando a citação. | Mutilação de formatação de citações diretas longas. | Implementar máquina de estados com memória de continuidade de citação e métrica em runes (`>= 240 runes`). |
| **ADV-04** | **HIGH** | Elementos Pós-Textuais | O pós-textual se resume à seção de REFERÊNCIAS. | Documentos com Apêndices e Anexos sofrem falsos positivos de legendas deslocadas no linter e são excluídos do sumário e seções. | Inconformidade estrita com NBR 14724 item 4.3 e NBR 6027. | Expandir Seção 2 até o fim dos anexos, adicionar regex para `APÊNDICE` e `ANEXO` e ajustar máquina de estados do linter. |
| **ADV-05** | **MEDIUM** | Tabelas Nativas | Aplicar `indentStart: 0` e `indentFirstLine: 0` em todas as células de tabelas é seguro. | Células contendo itens de listas com marcadores perdem o recuo e os marcadores colidem com o texto ou invadem as bordas. | Desalinhamento visual grave em matrizes de requisitos e tabelas comparativas. | Preservar `indentStart` quando `p.Bullet != nil` e forçar estritamente alinhamento `START` em células. |
| **ADV-06** | **MEDIUM** | Sumário de 2 Colunas | Tabela estática de 2 colunas atende aos padrões de documentos eletrônicos sem navegação. | Documento digital não permite clique de hiperlink para os capítulos e sofre de feedback loop de paginação multi-página. | Baixa usabilidade e descumprimento de acessibilidade NBR 6027 item 4.2. | Vincular `headingId` nos títulos da Coluna 1 e ordenar inserção da tabela antes da resolução de páginas. |
| **ADV-07** | **LOW** | Governança | `PROJECT.md` reflete fielmente o catálogo e status de ferramentas. | Desalinhamento entre a listagem de 17 ferramentas de `PROJECT.md` e a consolidação de `apply_abnt_page_margins` em `format_document_abnt`. | Confusão de escopo entre agentes e desenvolvedores. | Sincronizar catálogo do `PROJECT.md` com a arquitetura real consolidada. |

---

## 4. Conclusão e Veredicto

### Veredicto: **REQUEST_CHANGES**

O plano de engenharia elaborado pelo `worker_plan_1` possui sólidos fundamentos teóricos e resolve a grande maioria das limitações do sistema atual. No entanto, sua aprovação imediata para implementação introduziria riscos críticos de regressão na integridade de documentos reais (especialmente o risco de exclusão de quebras de página nativas e a lacuna de implementação dos números do sumário).

Para que o plano seja formalmente aprovado e liberado para a fase de codificação, os seguintes ajustes devem ser incorporados em uma revisão (v2.1.0) de `docs/plano_engenharia_abnt.md` e no `PROJECT.md`:
1. **Blindagem anti-deleção de `PageBreak`** explicitada no algoritmo de saneamento de linhas em branco.
2. **Definição da estratégia de resolução de números de página** para a tabela de sumário.
3. **Máquina de estados de continuidade para citações diretas longas multi-parágrafo** com limiar numérico formal (`>= 240 runes`).
4. **Tratamento de Apêndices e Anexos** na máquina de estados do linter, na divisão de seções e na montagem do sumário.
5. **Preservação de listas e forçamento de alinhamento à esquerda (`START`)** nas células de tabelas.
6. **Inclusão de navegação interna (`headingId`)** na Coluna 1 da tabela de sumário.
7. **Harmonização do inventário de ferramentas** no `PROJECT.md`.
