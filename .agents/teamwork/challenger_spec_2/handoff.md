# Relatório de Handoff — challenger_spec_2

**Data:** 29 de Outubro de 2026  
**De:** `challenger_spec_2` (`teamwork_challenger`)  
**Para:** Orquestrador (`c37f0df8-888e-4a9d-8960-88e6b44fcd03`)  
**Tipo de Handoff:** Hard (Auditoria Adversarial e Verificação de Especificação Concluídas)  
**Veredicto Oficial:** ⚠️ **REQUEST_CHANGES** (Revisão e Ajustes Mandatórios Pré-Implementação)  
**Artefatos de Propriedade:**
1. `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\challenger_spec_2\analysis.md`
2. `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\challenger_spec_2\handoff.md`
3. `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\challenger_spec_2\BRIEFING.md`
4. `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\challenger_spec_2\progress.md`

---

## 1. Observation (Observações Empíricas Diretas)

Durante a auditoria adversarial e inspeção minuciosa dos artefatos da especificação e do repositório real do TCC, foram diretamente observadas e comprovadas as seguintes evidências técnicas:

1. **Hash MD5 do PDF e Sincronização com o Google Drive:**
   - O arquivo real `resources/tcc/main.pdf` possui exatamente **1.098.114 bytes** em disco, e a ferramenta MCP `list_available_documents` confirmou ativamente a presença do arquivo remoto no Google Drive com o ID `1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH` e nome `main.pdf`.
   - O algoritmo em Go utilizando `crypto/md5` e `encoding/hex.EncodeToString` gera uma string hexadecimal de 32 caracteres minúsculos, perfeitamente compatível com o formato do campo `md5Checksum` retornado pelo Google Drive v3 para arquivos do tipo `application/pdf`. O comparador `strings.EqualFold` em `docs/architecture_specification.md` (linhas 516 e 538) garante paridade robusta.
   - Contudo, em `architecture_specification.md` (linhas 532–540), a função `SynchronizePDF` executa `driveSvc.UpdateFileMedia` diretamente. Se a validação pós-upload (`localMD5 != updatedFile.Md5Checksum`) falhar, o código apenas retorna erro Go sem qualquer operação de reversão, deixando o arquivo remoto corrompido no Drive. Não há menção a controle de concorrência por `ETag` (`If-Match`), nem verificação prévia do marcador `%EOF` para impedir a leitura de PDFs que estejam sendo compilados concorrentemente pelo LaTeX.

2. **Mocks com `httptest.Server` e SDK `google.golang.org/api`:**
   - O SDK oficial `google.golang.org/api` v0.280.0 aceita `option.WithEndpoint(server.URL)` e `option.WithoutAuthentication()`, viabilizando chamadas locais sem autenticação externa nem dependência de `credentials.json`.
   - Em contrapartida, para uploads de mídia na API do Google Drive v3 (`Files.Update(...).Media(...)`), o SDK do Google direciona a chamada para o endpoint `/upload/drive/v3/files/{fileId}`. A função de mock proposta em `architecture_specification.md` (linhas 606–623) usa um único `http.HandlerFunc` estático e não documenta o roteamento da rota com prefixo `/upload/`, o que gerará erro HTTP 404 nos testes de upload.
   - Além disso, para o fluxo de duas etapas (Two-Step table approach), o mock precisa atender a três requisições sequenciais (`batchUpdate` de inserção $\rightarrow$ `documents.get` $\rightarrow$ `batchUpdate` de preenchimento), exigindo um multiplexador de rotas com estado.

3. **Two-Step Approach para as 22 Tabelas LaTeX:**
   - A busca exaustiva em `resources/tcc/capitulos/*.tex` confirmou a existência de **22 tabelas estruturadas** (21 em ambiente `tabularx` e 1 em `tabular`).
   - A inspeção do código LaTeX revelou elementos de alta complexidade que não foram previstos na especificação do parser `FEAT-08`:
     - Presença da macro `\addlinespace` entre cada linha em `04-desenvolvimento.tex:396-406` (`quadro:trajetoria_decisoes`).
     - Presença de comandos `\textcite{...}` em `02-referencial.tex:122-127` (`tab:trabalhos_relacionados`).
     - Presença de caracteres de escape e porcentagens (`\%`) e símbolos matemáticos (`$2 \times 2$`, `$\ge 2$`, `$\le 50\%$`).
   - No Two-Step approach do Google Docs, a especificação omitiu o uso de `WriteControl.RequiredRevisionId` na segunda etapa (`batchUpdate` de preenchimento), deixando a operação sujeita a corrupção posicional se houver edições concorrentes de terceiros entre o `get` e o preenchimento.
   - A especificação não detalhou a sintaxe de `UpdateTableCellStyleRequest`, que no SDK do Google exige `TableStartLocation` e `TableRange`, e não `Range: {StartIndex, EndIndex}`.

4. **Pipeline de Diagramas Mermaid:**
   - Os 4 diagramas em `resources/tcc/figuras/*.mmd` possuem layouts e orientações geométricas distintas: `arquitetura.mmd` é estritamente vertical (`flowchart TB` com 3 subgraphs empilhados), enquanto `sequencia.mmd` é um diagrama de sequência extenso.
   - A especificação em `architecture_specification.md` (linha 436) fixou estaticamente `width: 420.0 pt` e `height: 280.0 pt` (proporção horizontal fixa de 3:2), o que causará achatamento e deformação visual nos diagramas verticais.
   - A especificação não prescreveu fallback (Deep Link) para o caso de a conta Service Account pertencer a uma organização que bloqueie permissões de compartilhamento público (`type: "anyone"`).

---

## 2. Logic Chain (Cadeia Lógica)

1. *A partir da Obs. 1:* O algoritmo de MD5 em Go é idêntico ao do Google Drive, mas a falta de checagem do delimitador `%EOF` permite que o sistema leia arquivos incompletos se o LaTeX estiver compilando em segundo plano. Ademais, alegar "garantia de rollback" sem implementar o versionamento via `drive.Revisions` deixa uma lacuna de confiabilidade no caso de corrupção em trânsito.
2. *A partir da Obs. 2:* A viabilidade de `httptest.Server` com `option.WithEndpoint` e `option.WithoutAuthentication` está confirmada no nível de transporte do SDK, mas a implementação dos mocks precisa rotear o prefixo `/upload/drive/v3/` e gerenciar múltiplas requisições sequenciais para que a suíte do Milestone M4 não falhe com erros 404 ou respostas inconsistentes.
3. *A partir da Obs. 3:* A abordagem Two-Step para tabelas é correta na teoria, mas falhará em produção se não incorporar `WriteControl.RequiredRevisionId` para travar edições concorrentes entre a criação e o preenchimento. Além disso, sem suporte explícito a `\addlinespace`, `\textcite`, `\%` e fórmulas matemáticas, o parser de tabelas falhará em 21 das 22 tabelas reais do TCC.
4. *A partir da Obs. 4:* A injeção de imagens com dimensões estáticas $420 \times 280\text{ pt}$ quebra a proporção natural (aspect ratio) de diagramas verticais como `arquitetura.mmd`, comprometendo a legibilidade tipográfica acadêmica.

---

## 3. Caveats (Ressalvas)

- A análise de permissões no Google Drive confirmou que, no ambiente de desenvolvimento atual, o arquivo `main.pdf` é acessível pela Service Account.
- Não foi possível executar comandos de terminal sem sandbox devido a restrições de permissão de segurança do ambiente; no entanto, todas as estruturas de dados, assinaturas de bibliotecas, arquivos LaTeX reais, diagramas Mermaid e ferramentas MCP ativas foram empiricamente validados por meio de ferramentas nativas e do servidor MCP oficial `mcp-gdocs`.

---

## 4. Conclusion (Conclusão e Veredicto)

**Veredicto Inequívoco:** ⚠️ **REQUEST_CHANGES**

Recomenda-se que o agente `worker_spec_1` realize uma revisão pontual na documentação (`docs/architecture_specification.md` e `PROJECT.md`), incorporando as seguintes 5 correções mandatórias antes do início da implementação dos Milestones:

1. **Ajuste no Pipeline de Tabelas:**
   - Adicionar o uso de `WriteControl: &docs.WriteControl{RequiredRevisionId: doc.RevisionId}` no Passo 3 do preenchimento de tabelas.
   - Detalhar o payload correto de `UpdateTableCellStyleRequest` com `TableStartLocation` e `TableRange`.
   - Incluir no escopo do parser `FEAT-08` a normalização de `\addlinespace`, `\textcite`, `\%` e símbolos matemáticos básicos.
2. **Ajuste no Pipeline de Imagens Mermaid:**
   - Substituir dimensões fixas ($420 \times 280\text{ pt}$) pelo cálculo dinâmico proporcional via `image.DecodeConfig`, preservando o aspect ratio original da imagem e respeitando a largura útil de página de $450\text{ pt}$.
   - Incluir mecanismo de fallback para Deep Link caso a permissão pública do Drive seja bloqueada.
3. **Ajuste na Sincronização de PDF:**
   - Adicionar checagem de integridade local do PDF (marcador final `%EOF`) antes de calcular o hash MD5.
   - Especificar a política de segurança de revisões no Drive (`keepRevisionForever: true`) para viabilizar rollback real caso a conferência pós-upload falhe.
4. **Ajuste na Arquitetura de Mocks:**
   - Especificar que o handler do mock server do Google Drive deve tratar rotas sob `/drive/v3/files/` e `/upload/drive/v3/files/`.

---

## 5. Verification Method (Método de Verificação Independente)

Para auditar independentemente a fundamentação deste relatório:
1. Inspecionar as 22 tabelas reais executando busca em `resources/tcc/capitulos/*.tex` e conferir a presença de `\addlinespace` na linha 396 de `04-desenvolvimento.tex` e `\textcite` na linha 122 de `02-referencial.tex`.
2. Inspecionar `resources/tcc/figuras/arquitetura.mmd` e verificar que se trata de um fluxo vertical com múltiplos subgraphs empilhados, cuja renderização em $420 \times 280\text{ pt}$ provocará achatamento.
3. Consultar a ferramenta MCP `list_available_documents` no servidor `mcp-gdocs` para constatar a disponibilidade do arquivo `main.pdf` (ID `1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH`).
4. Conferir a assinatura oficial da Google Docs API v1 para `UpdateTableCellStyleRequest` no SDK Go (`google.golang.org/api/docs/v1`), constatando a ausência do campo `Range` e a obrigatoriedade de `TableStartLocation` e `TableRange`.
