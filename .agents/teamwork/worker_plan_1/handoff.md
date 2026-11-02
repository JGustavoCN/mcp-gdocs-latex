# Relatório de Handoff Técnico: Plano de Engenharia e Arquitetura de Conformidade ABNT

**Agente:** `worker_plan_1`  
**Data:** 01 de Novembro de 2026  
**Destinatário:** Orquestrador (`parent` - ID `8522a030-5a02-4454-bc22-525b48f7de69`)  
**Missão Concluída:** Elaboração do relatório técnico definitivo de viabilidade nas APIs Google Docs v1 e Drive v3 e do plano detalhado de engenharia (/plan) para conformidade estrita com as normas da ABNT no servidor MCP `mcp-gdocs-latex`, além da atualização formal do `PROJECT.md`.

---

## 1. Observation (Observações Diretas)

1. **Insumos Técnicos e Diagnósticos da Fase 0:**
   - `.agents/teamwork/spec_miner_api_2/analysis.md`: Mapeamento completo dos 8 tópicos ABNT NBR 14724 nas APIs Google Workspace, fórmulas de conversão métrica ($1\text{ cm} \approx 28,3464567\text{ pt}$), mapeamento de payloads `UpdateDocumentStyleRequest`, `UpdateParagraphStyleRequest`, `InsertSectionBreakRequest` e identificação das restrições fundamentais da Google Docs API v1 (ausência de `InsertAutoTextRequest` e `InsertTableOfContentsRequest`, e campo `TabStops` marcado como `read-only`).
   - `.agents/teamwork/explorer_diag_1/analysis.md`: Diagnóstico preciso das causas raízes no codebase Go (`internal/gdocs`):
     - `linter.go:259-768`: ausência absoluta de inspeção de `LineSpacing`, `SpaceAbove` e `SpaceBelow`, gerando falsos negativos em parágrafos com entrelinhas 1,15 ou espaçamento simples;
     - `linter.go:730-745`: regra `if currentState != statePreTextual` que silencia e ignora parágrafos vazios espúrios em transições de página pré-textuais;
     - `abnt_writer.go:700`: exclusão incondicional de tabelas `if el.Table != nil { continue }`, deixando parágrafos de células com espaçamento desconforme;
     - `abnt_writer.go:343-356`: limiar `consecutiveBlanks >= 2` em `HealDocTextAndBlanks`, preservando a primeira quebra de linha após quebra de página;
     - `cover_generator.go:308-356`: injeção de 18 `\n` na Capa e 14 `\n` na Folha de Rosto para emular `\vfill`, e silenciamento de erro na limpeza reversa (`cover_generator.go:911`: `_, _ = docsSvc.Documents.BatchUpdate(...)`);
     - `inspector.go:9-31, 117-132`: struct `ElementLayoutDetail` omite o campo `LineSpacing`, e a ferramenta não inspeciona `DocumentStyle` (margens) nem `SectionStyle`;
     - `reader.go:280`: descarta elementos com `PageBreak` porque `elem.TextRun == nil`.
   - `.agents/teamwork/explorer_doc_1/analysis.md`: Diagnóstico empírico do documento vivo `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` (*TCC Teste IA*, 541 elementos estruturais):
     - Execução real do `lint_document_abnt` reportou `0 anomalias identificadas` apesar de dezenas de quebras de linha manuais (elementos 14 a 31 e 56 a 68) e ausência de quebra de seção para paginação;
     - Preâmbulo institucional está geometricamente perfeito em 8,0 cm (`indent_start_pt: 226.772` e `indent_first_line_pt: 226.772`) com fonte 10 pt e alinhamento justificado;
     - Documento possui apenas 1 seção global e nenhuma quebra de seção precedente a `"1 INTRODUÇÃO"`.

2. **Integridade de Compilação e Testes Atuais:**
   - Comando executado: `go test ./...`
   - Resultado: Sucesso absoluto (código de saída 0 em todos os pacotes Go `internal/gdocs` e `internal/latex`).

---

## 2. Logic Chain (Cadeia de Raciocínio Lógico)

1. **Da Cegueira do Linter à Omissão do Formatador:**
   - *Observação 1.1:* O arquivo `linter.go` não possui leitura de `LineSpacing`, e `inspector.go` não extrai essa propriedade para a struct `ElementLayoutDetail`.
   - *Dedução:* Modelos de IA e usuários clientes do MCP não recebem qualquer sinalização de que o documento contém parágrafos formatados com 1,15 (padrão web do Docs) ou sem espaçamento ABNT 1,5. Como a ferramenta de auditoria afirma que há "0 anomalias", ações corretivas não são disparadas.
   - *Solução Arquitetada:* Implementar em `linter.go` o checker `checkLineSpacing` com o campo `LineSpacingAnomalies` em `LintReport`, e enriquecer `ElementLayoutDetail` em `inspector.go` com `LineSpacingPercent float64`.

2. **Da Limitação da Docs API v1 à Arquitetura de Paginação e Sumário:**
   - *Observação 1.1:* A API v1 do Google Docs não possui chamadas para inserção de `AutoText` (`PAGE_NUMBER`) dinâmico via REST nem para criação de `TableOfContents` nativo, além de marcar `TabStops` como somente-leitura sem suporte a leader dots.
   - *Dedução:* Qualquer tentativa de enviar `InsertAutoTextRequest` ou editar paradas de tabulação com pontos falha com erro 400.
   - *Solução Arquitetada:*
     - Para **Sumário ABNT**: Criar uma Tabela Nativa de 2 colunas sem bordas visíveis (Coluna 1: $14,0\text{ cm}$ justificado/esquerda; Coluna 2: $2,0\text{ cm}$ alinhado à direita), garantindo alinhamento perfeito sem depender de tabulações.
     - Para **Paginação ABNT**: Isolar as páginas pré-textuais inserindo `InsertSectionBreakRequest` (`SECTION_TYPE_NEXT_PAGE`) imediatamente antes de `1 INTRODUÇÃO`, mantendo a Seção 1 sem cabeçalho e vinculando o cabeçalho numerado exclusivamente à Seção 2 via `CreateHeaderRequest` com `sectionBreakLocation`.

3. **Das Quebras Espúrias ao Saneamento Cirúrgico:**
   - *Observação 1.1 e 1.2:* O limiar `consecutiveBlanks >= 2` em `abnt_writer.go` preserva a primeira quebra após a `PageBreak` da Folha de Rosto, rebaixando o título `"RESUMO"`. Ao mesmo tempo, `cover_generator.go` injeta dezenas de `\n` manuais e silencia o erro ao tentar deletar o parágrafo que abriga a quebra de página.
   - *Solução Arquitetada:* Modificar `HealDocTextAndBlanks` para aplicar `consecutiveBlanks >= 1` quando um parágrafo vazio anteceder diretamente o próximo título textual após uma quebra de página. Substituir os `\n` repetidos em `cover_generator.go` por espaçamento vertical paramétrico `spaceAbove` e proteger o parágrafo da `PageBreak` contra exclusão indevida.

4. **Da Preservação Estrutural de Citações Longas e Tabelas:**
   - *Observação 1.1 e 1.2:* Parágrafos contidos em tabelas são atualmente ignorados com `continue`, e citações diretas longas correm o risco de serem formatadas como texto normal (1,25 cm / 1,5).
   - *Solução Arquitetada:* Definir uma heurística determinística de classificação: parágrafos com recuo `>= 90 pt` ou com extensão $> 3$ linhas e fonte 10 pt são identificados como citação longa e recebem estritamente `indentStart: 113.386 PT` (4,0 cm), `indentFirstLine: 0 PT`, `lineSpacing: 100` e fonte 10 pt. As células de tabelas são iteradas e formatadas com espaçamento simples 1,0 e fonte 10 pt.

---

## 3. Caveats (Ressalvas e Limitações)

1. **Dependência de Template para Injeção Dinâmica de `AutoText`:**
   Como a API Google Docs v1 não suporta a criação de novos tokens de número de página via chamadas REST (`InsertAutoTextRequest`), novos documentos gerados do zero via MCP devem utilizar o template de referência canônico (`resources/template/` ou ID de template do Drive) onde o campo dinâmico de número de página já existe no cabeçalho padrão. O MCP encarrega-se do particionamento de seções e desvinculação da Seção 1.
2. **Modo de Exibição Paginado (`PAGES` vs `PAGELESS`):**
   A renderização de margens ABNT e números de página no cabeçalho só é visível se o documento estiver configurado no modo paginado. Caso o usuário configure o documento no Google Docs como "Sem páginas" (*Pageless*), as margens e cabeçalhos continuam salvos nos metadados da API, porém não serão renderizados visualmente na interface web.

---

## 4. Conclusion (Conclusão e Entregas)

A especificação técnica e o plano de engenharia foram completamente consolidados, detalhados e formalizados nos artefatos do repositório:
1. **Artefato Principal:** `c:\Projetos\golang\mcp-gdocs-latex\docs\plano_engenharia_abnt.md`
   - Documento definitivo com 5 seções exaustivas:
     - *Seção 1:* Relatório técnico de viabilidade nas APIs Google Workspace para os 8 tópicos ABNT, com conversões matemáticas exatas e payloads JSON reais.
     - *Seção 2:* Diagnóstico detalhado do codebase Go e do documento de teste `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` com mapeamento linha por linha de causas raízes.
     - *Seção 3:* Plano detalhado de engenharia (/plan) para `format_document_abnt`, `lint_document_abnt`, `inspect_doc_layout`, `read_doc_content` e `cover_generator.go`.
     - *Seção 4:* Matriz de cooperação e heurística determinística de preservação de citações longas e tabelas.
     - *Seção 5:* Roteiro de implementação em Go com novas structs e suíte de testes com mocks `httptest.Server`.
2. **Atualização Arquitetural:** `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md`
   - Atualizado para a versão **1.3.0** (01 de Novembro de 2026).
   - Inclusão das funcionalidades **FEAT-20 a FEAT-24** no Feature Inventory oficial.
   - Adição do **Milestone 6 (M6): Conformidade Estrita ABNT (Entrelinhas, Margens, Quebras e Auditoria)** com entregáveis técnicos detalhados e critérios de aceite formais.
3. **Validação do Sistema:** Todos os testes unitários do workspace executam com sucesso via `go test ./...`.

---

## 5. Verification Method (Método de Verificação Independente)

Para auditar e verificar de forma independente as entregas realizadas:

1. **Inspeção dos Artefatos Produzidos:**
   - Verificar a existência e integridade do plano detalhado:
     ```powershell
     Get-Item "c:\Projetos\golang\mcp-gdocs-latex\docs\plano_engenharia_abnt.md"
     ```
   - Inspecionar a presença das 5 seções obrigatórias no plano:
     ```powershell
     Select-String -Path "c:\Projetos\golang\mcp-gdocs-latex\docs\plano_engenharia_abnt.md" -Pattern "^# Seção"
     ```
   - Inspecionar a atualização do `PROJECT.md` contendo FEAT-20 a FEAT-24 e o Milestone 6:
     ```powershell
     Select-String -Path "c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md" -Pattern "FEAT-20|Milestone 6"
     ```

2. **Verificação de Compilação e Testes no Repositório Go:**
   - Executar a suíte de testes automatizados do projeto:
     ```powershell
     go test ./...
     ```
   - *Condição de Sucesso:* Saída `ok` nos pacotes `internal/gdocs` e `internal/latex` com código de retorno 0.
