# Relatório de Handoff — worker_spec_1

**Data:** 29 de Outubro de 2026  
**De:** `worker_spec_1` (`teamwork_worker`)  
**Para:** Orquestrador (`c37f0df8-888e-4a9d-8960-88e6b44fcd03`)  
**Tipo de Handoff:** Hard (Especificação Arquitetural e Planejamento Concluídos)  
**Artefatos Produzidos e de Propriedade Exclusiva:**
1. `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md`
2. `c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md`
3. `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_spec_1\handoff.md`

---

## 1. Observation (Observação)

Durante a execução da tarefa e consolidação das investigações dos três agentes exploradores (`spec_miner_api_1`, `explorer_codebase_1`, `explorer_tcc_1`), foram diretamente observados e validados os seguintes pontos:

1. **Código-Fonte e Diagnósticos Críticos do Servidor Go:**
   - Em `internal/gdocs/writer.go` (linhas 292–305), a cláusula `if bloco.ParagraphStyle != "NORMAL_TEXT" && bloco.ParagraphStyle != ""` descarta sumariamente a geração de `UpdateParagraphStyleRequest` quando `bloco.ParagraphStyle == "NORMAL_TEXT"`, impedindo a normalização de cabeçalhos e a remoção de heranças de estilo de parágrafo.
   - Em `internal/gdocs/writer.go` (linhas 29–41), a ferramenta `SearchAndReplaceText` faz uso exclusivo de `ReplaceAllTextRequest`, a qual copia passivamente o `TextStyle` do primeiro caractere casado para todo o bloco injetado, provocando o fenômeno de Bold Bleeding observado em `resources/Avaliação de uma Arquitetura Híbrida... (2).md` (linhas 49–68).
   - Em `internal/gdocs/reader.go` (linhas 171–240) e `writer.go`, não há suporte a marcadores de lista (`paragraph.Bullet`), nem execução de `DeleteParagraphBulletsRequest`. Isso provocou o surgimento de marcadores de lista vazios (`* `) ao injetar múltiplas quebras de linha sobre blocos com marcadores (linha 226 do documento exportado).
   - Em `internal/gdocs/upload.go` (linhas 11–30), a ferramenta `SyncPDFToDrive` (31 linhas de código) executa `svc.Files.Update` cego de mídia, sem cálculo de hash MD5 local (`crypto/md5`) e sem consulta ao metadado `md5Checksum` da Google Drive API v3 (`files.get`).
   - A base Go contém 11 arquivos de código e **zero arquivos de teste** (`*_test.go`), resultando em 0% de cobertura automatizada. O código existente acopla-se diretamente aos singletons `*docs.Service` e `*drive.Service` retornados em `internal/gdocs/auth.go`.

2. **Repositório do TCC Canônico (`resources/tcc`):**
   - O arquivo mestre `main.tex` e o arquivo de configuração `config.tex` estabelecem as normas ABNT NBR 14724 / NBR 6023 (`abntex2`, `12pt`, Palatino via `mathpazo`, margens 3/3/2/2 cm, recuo de 1,25 cm e entrelinhas de 1,5 linha).
   - Foram mapeadas exatamente **22 tabelas estruturadas** nos capítulos do TCC (1 no Cap 1, 1 no Cap 2, 5 no Cap 3 e 15 no Cap 4), sendo 21 em ambiente `tabularx` e 1 em `tabular`, todas utilizando réguas `booktabs` (`\toprule`, `\midrule`, `\bottomrule`).
   - Foram mapeados **4 diagramas conceituais Mermaid** em `resources/tcc/figuras/` (`arquitetura.mmd`, `casos_uso.mmd`, `procedimento_experimental.mmd`, `sequencia.mmd`), todos com vetores PDF compilados correspondentes (`.pdf`), porém sem versões rasterizadas PNG (as quais são exigidas pela Docs API).
   - O arquivo PDF compilado do TCC `resources/tcc/main.pdf` possui exatamente **1.098.114 bytes** e corresponde ao arquivo remoto no Google Drive com ID `1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH`.

3. **Especificações das APIs Oficiais:**
   - Na Google Docs API v1, toda mutação estrutural ocorre via `POST /v1/documents/{documentId}:batchUpdate`, cuja transação possui atomicidade estrita (rollback em caso de erro em qualquer operação).
   - Os índices do Google Docs são medidos estritamente em **unidades UTF-16**, exigindo conversão precisa a partir de strings UTF-8 em Go via `unicode/utf16`.
   - A Google Docs API não permite inserção de tabelas e preenchimento de suas células no mesmo lote sem prévia descoberta dos índices estruturais de cada célula via `GET /v1/documents/{documentId}` (Two-Step Approach).
   - A Google Drive API v3 sempre fornece o campo `md5Checksum` para arquivos binários (`application/pdf`, imagens), permitindo verificação de paridade sem tráfego de dados.

---

## 2. Logic Chain (Cadeia Lógica)

1. **Da Causa do Bold Bleeding à Solução Canônica:**
   - *Premissa:* `ReplaceAllTextRequest` herda as formatações do primeiro caractere substituído.
   - *Premissa:* Substituições frequentes ancoraram-se em títulos em negrito (`HEADING_1`, `bold: true`).
   - *Conclusão:* A substituição de blocos multi-linha deve abandonar `ReplaceAllTextRequest` e adotar a sequência atômica de 5 passos via `BatchUpdateBuilder`: `DeleteContentRangeRequest` + `InsertTextRequest` + `DeleteParagraphBulletsRequest` + `UpdateParagraphStyleRequest(NORMAL_TEXT, ABNT)` + `UpdateTextStyleRequest(bold:false)`.

2. **Da Correção do Bug de Parágrafo (`writer.go:292`):**
   - *Premissa:* A expressão `bloco.ParagraphStyle != "NORMAL_TEXT"` impedia que requisições com estilo `"NORMAL_TEXT"` fossem geradas.
   - *Conclusão:* A remoção dessa exclusão e a geração explícita de `UpdateParagraphStyleRequest` com `namedStyleType: "NORMAL_TEXT"` restaura a capacidade de desarmar cabeçalhos indesejados.

3. **Da Eliminação do List Bleeding e Bullets Vazios:**
   - *Premissa:* Parágrafos no Google Docs retêm a propriedade `bullet` se não forem explicitamente desvinculados, convertendo cada `\n` inserido em um novo item da lista.
   - *Conclusão:* A inclusão obrigatória de `DeleteParagraphBulletsRequest` sobre a faixa editada desarma o marcador de lista e extingue bullets órfãos.

4. **Da Métrica ABNT para a Google Docs API:**
   - *Premissa:* $1\text{ cm} = 28,3464567\text{ pt}$.
   - *Conclusão:*
     - Recuo de primeira linha ($1,25\text{ cm}$): $35,433\text{ pt}$.
     - Entrelinhas ($1,5\text{ linha}$): valor relativo `150`.
     - Alinhamento: `JUSTIFIED`.
     - Margens Superior e Esquerda ($3,0\text{ cm}$): $85,039\text{ pt}$.
     - Margens Inferior e Direita ($2,0\text{ cm}$): $56,693\text{ pt}$ via `UpdateDocumentStyleRequest`.

5. **Da Inserção de Tabelas Nativas (Two-Step Approach):**
   - *Premissa:* Representações em ASCII Art são destruídas por fontes proporcionais do Google Docs e o cálculo antecipado de células antes da inserção falha na API.
   - *Conclusão:* O método robusto deve ser executado em duas fases: (1) `InsertTableRequest` criando a grade $R \times C$; (2) leitura dos índices de células via `documents.get`; (3) lote atômico de `InsertTextRequest` executado em ordem reversa decrescente de índice de início, finalizando com `UpdateTableCellStyleRequest` para zerar bordas laterais (norma IBGE).

6. **Do Pipeline de Diagramas Mermaid:**
   - *Premissa:* A Docs API exige imagem rasterizada pública e não aceita código Mermaid nem vetor PDF.
   - *Conclusão:* O pipeline converte `.mmd`/`.pdf` em PNG 300 DPI, realiza upload temporário no Google Drive via `drive.Files.Create`, concede leitura `anyone/reader` e injeta a imagem inline via `InsertInlineImageRequest`. O Google Docs persiste o binário internamente.

7. **Da Sincronização Inteligente de PDF por Hash MD5:**
   - *Premissa:* `resources/tcc/main.pdf` é binário de 1.098.114 bytes e o Drive armazena o hash MD5 no campo `md5Checksum`.
   - *Conclusão:* O servidor deve computar o MD5 local por streaming (`crypto/md5`) e confrontar com `files.get?fields=md5Checksum`. Se iguais, o upload de 1,1 MB é dispensado (`IN_SYNC`). Se divergentes, realiza upload e confere o hash pós-gravação.

8. **Da Testabilidade e Desacoplamento:**
   - *Premissa:* A ausência de testes decorre do acoplamento a singletons do SDK oficial do Google.
   - *Conclusão:* A introdução do Adapter Pattern (`DocsService`, `DriveService`) e servidores locais em memória via `httptest.Server` viabiliza cobertura total de testes automatizados sem credenciais reais e sem internet.

---

## 3. Caveats (Ressalvas e Suposições)

1. **Políticas de Compartilhamento Corporativo no Drive:**
   - A inserção de imagem inline via `InsertInlineImageRequest` consome a URL direta `webContentLink` do Google Drive, a qual exige que o arquivo tenha permissão pública de leitura (`anyone/reader`). Se a conta Service Account pertencer a uma organização Google Workspace com restrição severa que proíba compartilhamento com "qualquer pessoa com o link", a API do Docs rejeitará a URL. No ambiente de desenvolvimento e contas Google padrão, o fluxo é totalmente funcional.
2. **Modo do Documento (Paged vs Pageless):**
   - As margens de página definidas em `UpdateDocumentStyleRequest` só produzem efeito visual quando o documento estiver configurado no modo `PAGED` (com paginação física). Em modo `PAGELESS`, as margens são ignoradas pelo Docs.
3. **Equações Matemáticas em Tabelas LaTeX:**
   - As tabelas do TCC que contenham notações matemáticas simples (como $2 \times 2$) devem ter sua notação convertida em texto simples ou caracteres Unicode adequados, dado que a Google Docs API v1 não possui endpoint para inserção de blocos matemáticos LaTeX inline via texto comum.

---

## 4. Conclusion (Conclusão)

A especificação técnica formal e o plano de projeto foram integralmente elaborados e registrados:
- `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md` estabelece o plano geral de projeto, a arquitetura do servidor Go e do ecossistema MCP, a Feature Inventory exaustiva com 19 funcionalidades vinculadas aos Milestones M1 a M5, os contratos completos de interface em Go, schemas de ferramentas MCP e o layout estrutural do código.
- `c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md` documenta a especificação técnica formal para cada um dos cinco requisitos (R1 a R5), com payloads canônicos da Google Docs API v1, Google Drive API v3, cálculo preciso de unidades UTF-16, parametrização ABNT exata, algoritmo Two-Step para tabelas, pipeline para diagramas Mermaid, sincronização criptográfica de PDF e framework de mocks com `httptest.Server`.

O projeto encontra-se em estado pronto para que os agentes implementadores iniciem o Milestone 1 (M1).

---

## 5. Verification Method (Método de Verificação Independente)

Para auditar e validar independentemente a integridade e precisão dos artefatos gerados:

1. **Inspeção do `PROJECT.md`:**
   - Inspecionar `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md`.
   - Verificar a presença da tabela Feature Inventory (Seção 2), dos Milestones M1 a M5 (Seção 3), dos contratos Go e JSON schemas (Seção 4) e da árvore de Code Layout (Seção 5).

2. **Inspeção de `docs/architecture_specification.md`:**
   - Inspecionar `c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md`.
   - Confirmar a presença das seções correspondentes a R1, R2, R3, R4 e R5.
   - Conferir os cálculos métricos ABNT: $1,25\text{ cm} \rightarrow 35,433\text{ pt}$; $3,0\text{ cm} \rightarrow 85,039\text{ pt}$; $2,0\text{ cm} \rightarrow 56,693\text{ pt}$.
   - Conferir a citação cirúrgica da linha defeituosa `internal/gdocs/writer.go:292` e sua correção.
   - Conferir o mapeamento das 22 tabelas do TCC e dos 4 diagramas Mermaid.
   - Conferir a especificação do hash MD5 do `resources/tcc/main.pdf` (1.098.114 bytes) contra o ID `1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH`.

3. **Condições de Invalidação:**
   - Esta especificação será invalidada caso a Google Docs API passe a aceitar nativamente tabelas populadas em chamada única de `InsertTableRequest` sem necessidade de prévia consulta de índices de célula.
   - A especificação será invalidada se a Google Drive API v3 deixar de retornar o campo `md5Checksum` para arquivos com MIME `application/pdf`.
