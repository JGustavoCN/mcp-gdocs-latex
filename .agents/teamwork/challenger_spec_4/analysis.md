# Relatório de Análise e Testes Adversariais — Iteração 2

**Data:** 29 de Outubro de 2026  
**Autor:** `challenger_spec_4` (`teamwork_challenger`)  
**Missão:** Validação adversarial da robustez de sistemas e integridade arquitetural na Iteração 2 (`mcp-gdocs-latex`).  
**Documentos Auditados:**
- `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md`
- `c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md`
- `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_spec_2\handoff.md`
- `c:\Projetos\golang\mcp-gdocs-latex\resources\tcc\figuras\*.mmd`
- `c:\Projetos\golang\mcp-gdocs-latex\resources\tcc\capitulos\*.tex`

---

## Challenge Summary

**Overall risk assessment**: **HIGH**

Embora a especificação da Iteração 2 tenha incorporado avanços significativos (como a introdução de `WriteControl: RequiredRevisionId`, normalização NFC, checagem de integridade e adapters desacoplados), a análise adversarial identificou **2 falhas críticas e 2 vulnerabilidades moderadas** que causariam falhas em tempo de execução ou corrupção de documentos no Google Docs se implementadas estritamente conforme especificado:

1. **[CRITICAL] Corrupção Posicional da Fonte de Tabelas no Passo 3:** A inserção da Fonte no rodapé da tabela foi posicionada após o preenchimento de células no Passo 3 utilizando o `tableEndIndex` obsoleto do Passo 2. As inserções nas células expandem a tabela, fazendo com que o `tableEndIndex` antigo aponte para o *interior* da tabela, inserindo o texto da fonte dentro de células de dados ou disparando erro HTTP 400.
2. **[HIGH] Estouro Vertical de Página em Diagramas Mermaid Verticais:** O dimensionamento baseado exclusivamente na largura ($\text{width} = 420.0\text{ pt}$) faz com que diagramas verticais longos como `sequencia.mmd` alcancem alturas superiores a $750\text{ pt}$ a $945\text{ pt}$, excedendo a altura útil máxima de uma página A4 ABNT ($700,158\text{ pt}$), causando overflow visual e quebras desastrosas.
3. **[MEDIUM] Dependência de Driver PNG e Divisão por Zero em Mermaid:** O uso de `image.DecodeConfig` exige obrigatoriamente a inclusão do import anônimo `_ "image/png"` no Go; sua omissão em tempo de compilação gera erro de execução `image: unknown format`. Adicionalmente, ausência de guarda para `cfg.Width == 0` gera `+Inf`.
4. **[MEDIUM] Marcador `%EOF` vs `%%EOF` e Vulnerabilidade TOCTOU:** O marcador no código está definido como `%EOF` (1 caractere percentual), em desacordo com a ISO 32000-1 (`%%EOF`, 2 percentuais), e inexiste validação pós-hashing contra truncamento concorrente durante o streaming de `io.Copy`.
5. **[MEDIUM] Armadilha de Multipart e Prefixo de URL no Mock Server:** O mock server de upload necessita lidar com corpos `multipart/related` e normalizar a rota `/upload/drive/v3/files/` ao extrair o `fileID`.

---

## Challenges

### [Critical] Desafio 1: Corrupção Posicional da Fonte da Tabela no Passo 3 do Two-Step Approach

- **Suposição Desafiada:** A especificação afirma que o lote reverso de preenchimento de células no Passo 3 seguido da inserção da Fonte em `tableEndIndex` preserva a integridade posicional do documento.
- **Cenário de Ataque / Falha:**
  1. No Passo 1, a tabela vazia é inserida.
  2. No Passo 2, `GetDocument` captura `tableEndIndex` (ex: índice 120, correspondente ao fim da tabela vazia).
  3. No Passo 3, o algoritmo itera sobre as células preenchendo textos:
     - Uma tabela $5 \times 4$ (20 células) recebe em média 25 caracteres por célula, totalizando 500 caracteres UTF-16 inseridos dentro da estrutura da tabela.
     - À medida que cada `InsertTextRequest` é aplicado sequencialmente pela Google Docs API, o final físico da tabela é empurrado de 120 para 620.
  4. Ao término do loop de células, o lote executa a requisição do Item 4:
     `InsertTextRequest { location: { index: tableEndIndex }, text: "\nFonte: " + Source }`
     Como `tableEndIndex` foi fixado em 120 no Passo 2, a Google Docs API tentará inserir `\nFonte: ...` na posição 120!
  5. A posição 120 agora não é mais o fim da tabela: é o meio de uma célula da primeira ou segunda linha da tabela!
- **Raio de Destruição (Blast Radius):**
  - O texto da fonte ("Fonte: Autoria própria") é injetado corrompendo os dados tabulares internos de uma célula, OU
  - A API do Docs rejeita a requisição inteira com erro `400 Bad Request: The insertion index must be inside the bounds of an existing paragraph` ou conflito de segmento estrutural.
- **Mitigação Recomendada:**
  Aplicar o princípio de ordenação decrescente de índices ("highest index first") de maneira estrita:
  - Como `tableEndIndex` capturado no Passo 2 é estritamente maior do que qualquer `CellIndex[r][c]` interno da tabela, a inserção da Fonte deve ser **a PRIMEIRA operação gerada no Passo 3**!
  - Ao ser executada em primeiro lugar, a inserção da Fonte ocorre exatamente no `tableEndIndex` (início do parágrafo subsequente à tabela vazia).
  - Como ela se situa *após* todas as células, essa inserção não altera nenhum dos índices `CellIndex[r][c]`.
  - Em seguida, o preenchimento reverso das células $(R-1, C-1) \dots (0, 0)$ altera apenas os índices já processados, garantindo paridade posicional absoluta.
  - Alternativa: recalcular `tableEndIndexAdjusted = tableEndIndex + totalInsertedChars`. A solução de inverter a ordem no lote (inserir fonte primeiro) é muito mais limpa e imune a erros aritméticos.

---

### [High] Desafio 2: Estouro Vertical de Página em Diagramas Mermaid (`sequencia.mmd`)

- **Suposição Desafiada:** A fórmula $\text{targetHeightPt} = \text{targetWidthPt} \times (\text{cfg.Height} / \text{cfg.Width})$ com $\text{targetWidthPt} = 420.0\text{ pt}$ garante que os diagramas cabem na página A4 sem deformação.
- **Cenário de Ataque / Falha:**
  1. A largura útil de uma página A4 ABNT (margens 3cm esq, 2cm dir) é:
     $$\text{Largura Útil} = 595,276\text{ pt} - 85,039\text{ pt} - 56,693\text{ pt} = 453,544\text{ pt}$$
  2. A altura útil de uma página A4 ABNT (margens 3cm sup, 2cm inf) é:
     $$\text{Altura Útil} = 841,890\text{ pt} - 85,039\text{ pt} - 56,693\text{ pt} = 700,158\text{ pt}$$
  3. Descontando o espaço mínimo ocupado pela Legenda no topo (~25 pt), Fonte no rodapé (~25 pt) e espaçamentos ABNT (~20 pt), a altura livre máxima para a imagem é de aproximadamente **$630,0\text{ pt}$**.
  4. O diagrama `resources/tcc/figuras/sequencia.mmd` possui 41 linhas, 6 participantes horizontais e 15 interações com notas e 3 blocos `rect`. Ao ser renderizado a 300 DPI, suas dimensões típicas são de aproximadamente $900\text{ px}$ de largura por $1.900\text{ px}$ de altura (aspect ratio $H/W \approx 2,11$).
  5. Aplicando a fórmula especificada:
     $$\text{targetWidthPt} = 420.0\text{ pt}$$
     $$\text{targetHeightPt} = 420.0 \times 2,111 = 886,6\text{ pt}$$
  6. $886,6\text{ pt}$ ultrapassa a folha de papel A4 inteira ($841,890\text{ pt}$) e supera a área útil ($700,158\text{ pt}$) em mais de $186\text{ pt}$!
- **Raio de Destruição (Blast Radius):**
  - No Google Docs, uma imagem inline cujo tamanho ultrapassa a altura da página causa quebra visual de página defeituosa, transbordamento sobre o rodapé ou falha `400 Bad Request` na API v1.
- **Mitigação Recomendada:**
  Implementar algoritmo de enquadramento bidimensional (*2D Bounding Box Fit*):
  ```go
  const (
      MaxImageWidthPt  = 450.0 // Largura útil ABNT
      MaxImageHeightPt = 620.0 // Altura útil ABNT deduzidos legenda e rodapé
  )

  scaleW := MaxImageWidthPt / float64(cfg.Width)
  scaleH := MaxImageHeightPt / float64(cfg.Height)
  scale := math.Min(scaleW, scaleH)
  if customWidthPt > 0 && customWidthPt <= MaxImageWidthPt {
      // Se usuário solicitou largura específica, respeitar desde que não estoure a altura
      targetWidthPt = customWidthPt
      targetHeightPt = targetWidthPt * (float64(cfg.Height) / float64(cfg.Width))
      if targetHeightPt > MaxImageHeightPt {
          targetHeightPt = MaxImageHeightPt
          targetWidthPt = targetHeightPt * (float64(cfg.Width) / float64(cfg.Height))
      }
  } else {
      targetWidthPt = float64(cfg.Width) * scale
      targetHeightPt = float64(cfg.Height) * scale
  }
  ```
  Isso preserva $100\%$ do aspect ratio e garante matematicamente que a figura jamais transborde a folha nem na largura ($\le 450.0\text{ pt}$) nem na altura ($\le 620.0\text{ pt}$).

---

### [Medium] Desafio 3: Registro de Formato PNG e Divisão por Zero no Leitor de Imagens

- **Suposição Desafiada:** A chamada `image.DecodeConfig(bytes.NewReader(pngBytes))` é suficiente para extrair as dimensões da imagem PNG.
- **Cenário de Ataque / Falha:**
  - Na biblioteca padrão de Go, o pacote `image` não embute decodificadores automaticamente. Se o arquivo `mermaid_writer.go` contiver apenas `import "image"`, a execução de `image.DecodeConfig` retornará invariavelmente:
    `image: unknown format`
  - Se os bytes forem corrompidos ou a imagem possuir `cfg.Width == 0`, a expressão `float64(cfg.Height) / float64(cfg.Width)` avalia para `+Inf`, gerando pânico no JSON marshaling da requisição do Docs (`json: unsupported value: +Inf`).
- **Raio de Destruição:** Falha catastrófica de todas as chamadas da tool `insert_mermaid_diagram`.
- **Mitigação:**
  - Declarar explicitamente `import _ "image/png"` e `import _ "image/jpeg"`.
  - Adicionar guarda: `if cfg.Width <= 0 || cfg.Height <= 0 { return nil, fmt.Errorf("dimensões inválidas da imagem: %dx%d", cfg.Width, cfg.Height) }`.

---

### [Medium] Desafio 4: Marcador `%EOF` vs `%%EOF` e Corrida TOCTOU na Sincronização de PDF

- **Suposição Desafiada:** A checagem de `%EOF` nos últimos 1024 bytes antes do cálculo do MD5 garante que arquivos incompletos em compilação concorrente não sejam lidos.
- **Cenário de Ataque / Falha:**
  - O padrão oficial PDF (ISO 32000-1, Seção 7.5.5) determina que o marcador de fim de arquivo é `%%EOF` (dois sinais de porcentagem). Embora `%EOF` case com `%%EOF`, a busca por apenas 4 bytes (`%EOF`) em fluxos binários comprimidos (streams FlateDecode, fontes CFF/TrueType embutidas) possui probabilidade de colisão estatística de $1021 / 2^{32} \approx 2.38 \times 10^{-7}$ por bloco de 1024 bytes. O uso de `%%EOF` (5 bytes) reduz esse risco para $< 10^{-9}$.
  - Vulnerabilidade TOCTOU (Time-of-Check to Time-of-Use): Se o compilador LaTeX estiver reescrevendo `main.pdf` e o processo de compilação truncar o arquivo ou adicionar novos bytes *após* a checagem dos 1024 bytes e *durante* o `io.Copy(hasher, f)`, o tamanho lido por `io.Copy` divergirá de `stat.Size()`.
- **Mitigação:**
  - Buscar `[]byte("%%EOF")`.
  - Validar paridade de tamanho: `if copiedBytes != stat.Size() { return nil, fmt.Errorf("arquivo alterado concorrentemente durante o cálculo do hash") }`.

---

### [Medium] Desafio 5: Tratamento de Multipart e Prefixo de Rota no Mock Server do Drive

- **Suposição Desafiada:** O registro de `/upload/drive/v3/files/` no `http.ServeMux` resolve integralmente os testes de upload do SDK Google.
- **Cenário de Ataque / Falha:**
  - As requisições de upload com mídia enviadas por `google.golang.org/api` utilizam `Content-Type: multipart/related; boundary=...`.
  - Se o mock handler registrar apenas um decoder JSON simples (`json.NewDecoder(r.Body).Decode(&meta)`), o teste falhará com erro de sintaxe JSON ao ler os delimitadores multipart (`--boundary`).
  - Se o mock handler extrair o ID do arquivo via `strings.TrimPrefix(r.URL.Path, "/drive/v3/files/")`, em uploads a rota será `/upload/drive/v3/files/{id}`, fazendo com que o `fileID` capturado contenha lixo de rota.
- **Mitigação:**
  - No mock handler, normalizar a rota: `cleanPath := strings.TrimPrefix(r.URL.Path, "/upload")`.
  - Inspecionar `r.Header.Get("Content-Type")`: se contiver `multipart/related`, utilizar `mime/multipart.NewReader` ou responder com o metadado simulado sem forçar decode JSON direto do stream bruto.

---

## Stress Test Results

| Cenário de Teste | Comportamento Esperado | Comportamento Predito / Observado | Status |
|---|---|---|---|
| Inserção de tabela $5 \times 4$ com Legenda e Fonte | Legenda no topo, tabela preenchida e Fonte no rodapé externo | Fonte é inserida dentro da célula interna devido ao `tableEndIndex` desatualizado | **FAIL** |
| Renderização de `sequencia.mmd` ($900 \times 1900\text{ px}$) | Imagem cabe integralmente na área útil da página A4 ($\le 620\text{ pt}$) | Imagem atinge $886,6\text{ pt}$ estourando a página e sobrepondo rodapés | **FAIL** |
| Invocação de `image.DecodeConfig` sem driver | Decodificação instantânea de cabeçalho PNG | Retorno do erro `image: unknown format` | **FAIL** |
| Leitura de PDF incompleto em compilação ativa (0 bytes ou sem trailer) | Rejeição preventiva antes do cálculo de MD5 | Rejeição confirmada pelo marcador nos 1024 bytes | **PASS** |
| Roteamento de upload no Mock Server (`httptest.Server`) | Rota `/upload/drive/v3/files/{id}` atende sem HTTP 404 | Casamento de rota no `ServeMux` bem-sucedido | **PASS** |
| Prevenção de corrida concorrente de terceiros com `WriteControl` | Abortamento com HTTP 400 se `RevisionId` for divergente | Travamento estrito garantido pelo Google Docs | **PASS** |
| Clamping de fim de documento (`min(End, bodyEnd - 1)`) | Prevenção de erro HTTP 400 no `\n` inamovível | Proteção atômica eficaz contra boundary error | **PASS** |

---

## Unchallenged Areas

- **Normalização Unicode NFC e UTF16Length:** A especificação de `norm.NFC` e contagem de code units UTF-16 em `utf16_utils.go` é matematicamente sólida e alinhada às regras da Google Docs API.
- **Anulação de Recuo Residual de Lista (`indentStart: 0 pt`):** O mecanismo de reset de recuo ao deletar bullets foi exaustivamente verificado e atende aos requisitos do Docs.
- **Configuração de Margens ABNT de Página:** Os valores métricos em pontos tipográficos ($85,039\text{ pt}$ e $56,693\text{ pt}$) estão perfeitamente calibrados para A4 ($595,276 \times 841,890\text{ pt}$).
- **ListFiles na Interface DriveService:** A adição do método `ListFiles` resolve a lacuna de desacoplamento identificada na Iteração 1.
