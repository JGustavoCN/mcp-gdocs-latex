# Relatório de Handoff Técnico: Auditoria Adversarial e Desafio do Plano ABNT

**Agente:** `challenger_plan_1` (Critic & Specialist - Adversarial Challenger)  
**Data:** 01 de Novembro de 2026  
**Destinatário:** Orquestrador (`parent` - ID `8522a030-5a02-4454-bc22-525b48f7de69`)  
**Missão Concluída:** Submissão do plano de engenharia ABNT (`docs/plano_engenharia_abnt.md`, `PROJECT.md`, `worker_plan_1/handoff.md`) a testes adversariais empíricos, conceituais e de limites operacionais na Google Docs API v1 e normas ABNT.  
**Veredicto Formal:** **REQUEST_CHANGES (Ajustes Críticos Obrigatórios Antes da Implementação)**  

---

## 1. Observation (Observações Diretas)

1. **Restrições Formais da Google Docs API v1:**
   - A Google Docs API REST v1 **não possui nenhum endpoint, campo ou estrutura que retorne o número da página física** em que um elemento ou parágrafo está renderizado. O documento é exposto puramente como uma sequência linear de offsets UTF-16 (`startIndex` e `endIndex`).
   - A Google Docs API v1 modela `PageBreak` não como elemento raiz de `Body.Content`, mas como um filho `ParagraphElement` dentro de um `Paragraph` (`paragraph.elements[*].pageBreak`). Em parágrafos contendo exclusivamente uma quebra de página, `pel.TextRun` é nulo (`nil`), e o texto acumulado avaliado por `strings.TrimSpace` resulta em string vazia `""`.
   - A Docs API rejeita com erro HTTP 400 fatal (`Bad Request`) qualquer tentativa de excluir o parágrafo delimitador imediatamente anterior a um elemento `Table`, bem como qualquer requisição que tente deletar o caractere final `body.EndIndex - 1`.
   - A Docs API suporta hiperlinks internos para cabeçalhos através da propriedade `textStyle.link.headingId` em qualquer trecho de texto formatado.

2. **Evidências Diretas no Codebase Atual (`internal/gdocs` e `internal/latex`):**
   - `abnt_writer.go:334-342`: O método `HealDocTextAndBlanks` itera sobre `el.Paragraph.Elements` concatenando apenas `pel.TextRun.Content`. Se o parágrafo contiver `pel.PageBreak`, `t` é avaliado como `""`.
   - `cover_generator.go:891-912`: O loop reverso de saneamento tenta deletar parágrafos na transição para o Resumo e colide com o parágrafo da `PageBreak`, sendo forçado a silenciar o erro da Docs API com `_, _ = docsSvc.Documents.BatchUpdate(...)` para não quebrar a execução.
   - `structure_normalizer.go:56`: A expressão regular `regexPrimaryHeading` contempla apenas numerais seguidos de texto ou títulos fixos (`INTRODUÇÃO|...|REFERÊNCIAS|RESUMO|ABSTRACT`), omitindo `APÊNDICE` e `ANEXO`.
   - `linter.go:544-550`: Toda e qualquer legenda de ilustração ou tabela posicionada após a seção de REFERÊNCIAS dispara falso positivo de `MisplacedElements` ("posicionada indevidamente dentro da seção de REFERÊNCIAS"), inviabilizando tabelas/figuras legítimas em Apêndices e Anexos.
   - `builder.go:140`: O método `BuildABNTReplaceSequence` emite incondicionalmente `DeleteParagraphBulletsRequest`, destruindo listas com marcadores caso o bloco substituído esteja dentro de células de tabelas.
   - `tools.go:36`: Documenta que `apply_abnt_page_margins` foi incorporada em `format_document_abnt`, em divergência com a listagem formal de 17 ferramentas do `PROJECT.md:416-435`.

3. **Validação do Ambiente de Build:**
   - O comando `go test ./...` foi executado no workspace e finalizou com sucesso (código de saída 0 em todos os módulos).

---

## 2. Logic Chain (Cadeia de Raciocínio Lógico)

1. **Do Modelo de `PageBreak` ao Risco Crítico de Colapso de Páginas:**
   - *Premissa 1:* No Google Docs, parágrafos que contêm `PageBreak` possuem `TextRun == nil` e seu conteúdo textual acumulado é `""`.
   - *Premissa 2:* O plano propõe que na transição pós-quebra de página, parágrafos vazios sejam deletados mesmo com `consecutiveBlanks >= 1`.
   - *Dedução:* Se o código não inspecionar expressamente se `pel.PageBreak != nil` antes de considerar o parágrafo "vazio", o parágrafo que abriga a própria quebra de página da Folha de Rosto ou entre capítulos será selecionado para deleção. Ao executar a exclusão, a quebra de página física desaparece e as páginas colapsam em uma só.
   - *Conclusão:* É mandatório incluir uma cláusula de guarda estrutural (`pel.PageBreak == nil && pel.InlineObjectElement == nil && pel.FootnoteReference == nil`) no algoritmo de saneamento.

2. **Da Ausência de Numeração Física de Páginas na Docs API ao Sumário:**
   - *Premissa 1:* A Google Docs API REST v1 não fornece números de página física de nenhum elemento.
   - *Premissa 2:* O plano prescreve uma tabela de 2 colunas para o Sumário com a Coluna 2 exibindo o "número da página correspondente".
   - *Dedução:* Como a API não expõe esse valor, o código Go não tem fonte de dados intrínseca no Google Docs para saber em qual página o capítulo "2 REVISÃO BIBLIOGRÁFICA" se encontra.
   - *Conclusão:* O plano possui uma lacuna arquitetural e deve especificar formalmente o pipeline de resolução: via arquivo `.toc` gerado pela compilação do LaTeX canônico em `resources/tcc`, via exportação em memória de PDF pelo Drive (`files.export?mimeType=application/pdf`), ou via placeholder controlado.

3. **Da Heurística Atômica de Parágrafo à Fragmentação de Citações Longas:**
   - *Premissa 1:* Citações diretas longas frequentemente possuem mais de um parágrafo (transcrições de entrevistas, diálogos, múltiplos parágrafos de dispositivos legais).
   - *Premissa 2:* A heurística proposta no plano avalia cada parágrafo isoladamente: `existingIndent >= 90 pt OU (> 3 linhas E fonte <= 10.5 pt)`.
   - *Dedução:* Em um documento desformatado (sem recuo prévio e com fonte padrão 12 pt), se o segundo parágrafo da citação tiver apenas 2 linhas, ele falha em todos os critérios. O formatador o rebaixa para texto regular (`NORMAL_TEXT`, 1,25 cm, 1,5 entrelinhas, 12 pt), mutilando a citação direta.
   - *Conclusão:* A classificação precisa operar como máquina de estados com memória de continuidade de bloco de citação, aplicando espaçamento de 1 linha (`spaceAbove`/`spaceBelow: 6 pt`) apenas nas extremidades do bloco e nunca entre parágrafos internos.

4. **Da Exclusão de Apêndices e Anexos aos Falsos Positivos e Quebra de Sumário:**
   - *Premissa 1:* Apêndices e Anexos são elementos pós-textuais essenciais regulados pela ABNT NBR 14724 item 4.3 e NBR 6024 item 3.8, identificados por letras maiúsculas ordenadas (ex.: `APÊNDICE A`).
   - *Premissa 2:* O plano limita a Seção 2 textual/pós-textual até as Referências, e o linter assume que qualquer elemento após as referências pertence a elas.
   - *Dedução:* Além de não estarem previstos no particionamento de seções e paginação contínua, elementos legítimos dentro de apêndices (figuras, tabelas) disparam erros indevidos no linter.
   - *Conclusão:* O plano deve expandir formalmente o escopo da Seção 2 até o último anexo e ajustar as expressões regulares e estados do linter.

---

## 3. Caveats (Ressalvas e Limitações)

1. **Modo Pageless do Google Docs:**
   Se um documento estiver configurado pelo usuário na interface do Google Docs no modo "Sem páginas" (*Pageless*), cabeçalhos, quebras de página físicas e margens não são exibidos visualmente na interface Web, embora continuem preservados e manipuláveis via API REST para fins de exportação em PDF.
2. **Ambiente com Falta de LaTeX Instalado Localmente:**
   Caso a estratégia de extração de páginas do Sumário dependa do `.toc` compilado, em ambientes onde o compilador LaTeX (`pdflatex` ou `latexmk`) não esteja instalado no host, o sistema deverá obrigatoriamente recorrer ao fallback de exportação em PDF via Google Drive v3 ou gerar marcadores simbólicos de revisão.

---

## 4. Conclusion (Conclusão e Veredicto)

**Veredicto Oficial:** **REQUEST_CHANGES**

O plano de engenharia de `worker_plan_1` possui arquitetura sólida e conversões matemáticas precisas, mas **não pode ser aprovado em sua forma atual** sem os seguintes ajustes mandatórios em `docs/plano_engenharia_abnt.md` e `PROJECT.md`:

1. **Blindagem de `PageBreak` no Saneador:** Formalizar que nenhum parágrafo contendo `pel.PageBreak != nil` ou elementos de mídia/linha seja considerado vazio no algoritmo `consecutiveBlanks >= 1`.
2. **Definição do Pipeline de Paginação do Sumário:** Especificar explicitamente no plano como a Coluna 2 da tabela de sumário obtém os números de página (extração via `.toc` do LaTeX local em `resources/tcc/` ou via export de PDF no Drive).
3. **Máquina de Estados para Citações Longas Multi-Parágrafo:** Tratar blocos de citação contínua, aplicando espaçamento de 1 linha somente nas bordas externas do bloco e adotando métrica em runes (`>= 240 runes`) para o critério de > 3 linhas.
4. **Modelagem de Apêndices e Anexos:** Incorporar regex para identificadores alfabéticos (`APÊNDICE A`, `ANEXO A`), estender a Seção 2 de paginação até o término dos anexos e eliminar falsos positivos de figuras/tabelas no linter.
5. **Preservação de Listas em Tabelas:** Preservar `indentStart` quando `p.Bullet != nil` em células e consagrar o alinhamento estrito à esquerda (`START`).
6. **Inclusão de Navegação Interna (`headingId`):** Prever hyperlinks navegáveis e sem sublinhado nos títulos da Coluna 1 do sumário.
7. **Sincronização do `PROJECT.md`:** Harmonizar o catálogo oficial de ferramentas, consolidando formalmente `format_document_abnt` e `lint_document_abnt`.

O relatório analítico exaustivo com todos os cenários e matriz de impacto foi gravado em:  
`c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\challenger_plan_1\analysis.md`.

---

## 5. Verification Method (Método de Verificação Independente)

Para auditar e reproduzir as constatações deste relatório:

1. **Inspeção do Relatório de Análise Adversarial:**
   ```powershell
   Get-Content "c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\challenger_plan_1\analysis.md"
   ```

2. **Verificação Empírica da Cegueira a `PageBreak` no Código Atual:**
   Inspecionar as linhas 334 a 355 de `internal/gdocs/abnt_writer.go`:
   ```powershell
   Select-String -Path "c:\Projetos\golang\mcp-gdocs-latex\internal\gdocs\abnt_writer.go" -Pattern "PageBreak|TextRun" -Context 3,3
   ```
   *Resultado Observável:* Observa-se que apenas `pel.TextRun` é lido; `pel.PageBreak` é completamente ignorado na verificação de parágrafos vazios.

3. **Verificação da Inexistência de Regex para Apêndices e Anexos:**
   ```powershell
   Select-String -Path "c:\Projetos\golang\mcp-gdocs-latex\internal\gdocs\structure_normalizer.go" -Pattern "regexPrimaryHeading"
   ```
   *Resultado Observável:* Apenas `[1-9][0-9]*` e títulos textuais fixos estão contemplados; `APÊNDICE` e `ANEXO` estão ausentes.

4. **Verificação do Falso Positivo de Figuras Pós-Referências no Linter:**
   Inspecionar as linhas 544 a 550 de `internal/gdocs/linter.go`:
   ```powershell
   Select-String -Path "c:\Projetos\golang\mcp-gdocs-latex\internal\gdocs\linter.go" -Pattern "statePostTextual" -Context 2,5
   ```
   *Resultado Observável:* Qualquer legenda pós-referências é classificada indevidamente como violação de `MisplacedElements`.
