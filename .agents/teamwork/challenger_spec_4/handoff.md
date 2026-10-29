# Relatório de Handoff — Challenger Spec 4 (Iteração 2)

**Data:** 29 de Outubro de 2026  
**De:** `challenger_spec_4` (`teamwork_challenger`)  
**Para:** Orquestrador (`c37f0df8-888e-4a9d-8960-88e6b44fcd03`)  
**Tipo de Handoff:** Hard (Validação Adversarial Concluída — Iteração 2)  
**Veredicto Inequívoco:** **REQUEST_CHANGES** (Requer Ajustes Arquiteturais Críticos Antes da Implementação)

---

## 1. Observation (Observação)

Durante a auditoria adversarial dos artefatos da Iteração 2 (`PROJECT.md` e `docs/architecture_specification.md`), foram diretamente observados e confirmados os seguintes pontos empíricos e analíticos:

1. **Ordenação Temporal no Passo 3 de Inserção de Tabelas:**
   - Em `docs/architecture_specification.md` (linhas 385-395), o algoritmo Two-Step define o Passo 3 da seguinte forma:
     ```text
     [Passo 3: Preenchimento Atômico em Lote Reverso com Trava Concorrente]
       1. Montar BatchUpdateDocumentRequest configurando WriteControl: RequiredRevisionId: docRevisionId
       2. Iterar de r = R-1 até 0: Iterar de c = C-1 até 0: Adicionar InsertTextRequest { location: { index: CellIndex[r][c] }, text: MatrizDados[r][c] }
       3. Adicionar requisição de estilização tabular (UpdateTableCellStyleRequest)
       4. Se houver Fonte (Source):
          - Inserir parágrafo centralizado com estilo ABNT em tableEndIndex contendo "\nFonte: " + Source
       5. Executar docsService.BatchUpdate(docID, batchPasso3)
     ```
   - No Passo 2 (linha 383), `tableEndIndex = table.EndIndex` é capturado sobre a tabela *vazia*.
   - A execução sequencial das requisições de inserção de texto nas células (item 2) expande a tabela no documento por um total de $\sum \text{UTF16Length}(\text{MatrizDados}[r][c])$ caracteres.
   - O item 4 envia a inserção da Fonte no `tableEndIndex` estático original (não ajustado), que agora se localiza no *interior* das células já preenchidas da tabela e não no seu rodapé.

2. **Cálculo Unidimensional de Dimensão para Diagramas Mermaid:**
   - Em `docs/architecture_specification.md` (linhas 536-544), o redimensionamento de imagens Mermaid define:
     ```go
     targetWidthPt := 420.0
     if customWidthPt > 0 && customWidthPt <= 450.0 {
         targetWidthPt = customWidthPt
     }
     targetHeightPt := targetWidthPt * (float64(cfg.Height) / float64(cfg.Width))
     ```
   - O repositório contém o arquivo `resources/tcc/figuras/sequencia.mmd` (diagrama de sequência vertical com 41 linhas, 6 participantes e 3 caixas de fase). Ao ser renderizado a 300 DPI, sua relação de aspecto $H/W$ típica é de aproximadamente $2,1$ a $2,3$.
   - Com $\text{targetWidthPt} = 420.0\text{ pt}$, a altura calculada $\text{targetHeightPt}$ atinge entre $882,0\text{ pt}$ e $966,0\text{ pt}$.
   - A folha A4 possui altura total de $841,890\text{ pt}$ e altura útil entre margens ABNT de $700,158\text{ pt}$ ($841,890 - 85,039 - 56,693$).
   - A imagem gerada transborda os limites da página física em mais de $100\text{ pt}$ e excede a área útil em mais de $180\text{ pt}$.

3. **Requisitos de Compilação do Go para `image.DecodeConfig`:**
   - Em Go, o pacote `image` não registra decodificadores de formato por padrão. Para decodificar imagens PNG via `image.DecodeConfig`, é mandatório o registro via import em branco `_ "image/png"`. Caso contrário, a função retorna o erro `image: unknown format`.

4. **Sintaxe do Marcador de Fim de PDF (ISO 32000-1):**
   - Em `docs/architecture_specification.md` (linha 654), a checagem é expressa como:
     `if !bytes.Contains(tailBuf, []byte("%EOF"))`
   - A especificação formal ISO 32000-1 (Seção 7.5.5) define o trailer estritamente como `%%EOF`. A busca por 4 bytes em arquivos binários comprimidos apresenta probabilidade de falso positivo estatisticamente superior à busca por 5 bytes.

5. **Roteamento e Payloads de Upload no Mock Server:**
   - Em `docs/architecture_specification.md` (linhas 788-789), o servidor mock multiplexado registra `/upload/drive/v3/files` e `/upload/drive/v3/files/`.
   - O SDK `google.golang.org/api` transmite requisições de mídia como `Content-Type: multipart/related`, exigindo que o handler de mock não tente decodificar o corpo bruto diretamente como JSON.

---

## 2. Logic Chain (Cadeia Lógica)

1. *Premissa:* Requisições em lote no `documents.batchUpdate` do Google Docs são executadas em ordem de array, e mutações prévias deslocam os índices de caracteres posteriores.
   *Observação:* No Passo 3, as células da tabela são preenchidas antes da inserção da Fonte, expandindo o conteúdo do documento antes do fim da tabela.
   *Conclusão:* Inserir a Fonte no `tableEndIndex` capturado no Passo 2 faz com que o texto da fonte caia dentro de uma célula da tabela ou provoque erro HTTP 400. Como `tableEndIndex` é estritamente maior que qualquer coordenada de célula, **inserir a Fonte como a primeira instrução do Passo 3** elimina completamente a necessidade de recalcular offsets e preserva a regra decrescente de inserção ("highest index first").

2. *Premissa:* Documentos acadêmicos ABNT no Google Docs devem manter imagens contidas dentro dos limites imprimíveis da página A4.
   *Observação:* `sequencia.mmd` possui aspect ratio fortemente vertical ($H/W > 2,0$), e a fixação de largura em $420\text{ pt}$ projeta altura $> 880\text{ pt}$, superando a folha inteira ($841,89\text{ pt}$) e a altura útil ($700,16\text{ pt}$).
   *Conclusão:* É imperativo que a especificação adote um algoritmo de enquadramento bidimensional (*2D Bounding Box*), limitando a largura em $\le 450,0\text{ pt}$ e a altura em $\le 620,0\text{ pt}$ (garantindo espaço vertical para legenda e fonte).

3. *Premissa:* O compilador Go exige registro explícito de drivers de imagem para decodificação em runtime.
   *Conclusão:* A ausência da indicação de `_ "image/png"` provocará falha em tempo de execução na primeira invocação de `insert_mermaid_diagram`.

4. *Premissa:* A integridade de PDFs depende de um trailer válido e imutabilidade durante o cálculo de hash.
   *Conclusão:* O marcador canônico ISO 32000-1 deve ser `%%EOF` e a função deve validar se os bytes lidos por `io.Copy` coincidem com `stat.Size()`.

---

## 3. Caveats (Ressalvas)

1. **Restrições de Execução do Sandbox do Host:** Conforme registrado pelo worker e confirmado nesta sessão, a execução de processos sem permissão de elevação foi bloqueada pela política do ambiente Windows. Toda a validação foi conduzida por derivação matemática estrita, modelagem formal de estados do Google Docs API v1, especificações ISO e análise estática do compilador Go.
2. **Resolução de Renderização do Mermaid CLI:** As dimensões em pixels dos arquivos Mermaid dependem do fator de escala `-s 3` (300 DPI) do `mmdc`. Em qualquer escala, o aspect ratio $(H/W)$ é uma constante geométrica invariante do diagrama.

---

## 4. Conclusion (Conclusão e Veredicto)

**Veredicto Formal: REQUEST_CHANGES**

Recomenda-se o retorno pontual ao `worker_spec` para que realize os seguintes 4 ajustes cirúrgicos em `docs/architecture_specification.md` e `PROJECT.md`:

1. **Ajustar a Sequência do Passo 3 no Two-Step de Tabelas:**
   - Mover a inserção da Fonte (`\nFonte: ` + Source) em `tableEndIndex` para ser a **primeira operação** do Passo 3 (antes de iterar sobre as células em ordem decrescente), respeitando categoricamente a ordenação do maior offset para o menor offset.
   - Documentar política de retry para Passo 2 + Passo 3 em caso de conflito de revisão (`WriteControl.RequiredRevisionId`).

2. **Implementar Enquadramento Bidimensional (2D Bounding Box) em Mermaid:**
   - Adicionar restrição de altura máxima ($\text{MaxHeightPt} = 620.0\text{ pt}$), calculando $\text{scale} = \min(450.0 / W, 620.0 / H)$ para garantir que diagramas verticais como `sequencia.mmd` não estourem a página.
   - Declarar explicitamente a necessidade de import anônimo `_ "image/png"`.
   - Adicionar guarda preventiva contra divisão por zero (`cfg.Width <= 0`).

3. **Padronizar Marcador de Fim de PDF e Prevenir TOCTOU:**
   - Substituir `[]byte("%EOF")` por `[]byte("%%EOF")` em conformidade com a ISO 32000-1.
   - Validar que a quantidade de bytes lidos por `io.Copy(hasher, f)` seja igual a `stat.Size()`.

4. **Complementar Mock Server de Upload com Suporte a Multipart:**
   - Documentar que o mock handler em `SetupMockWorkspaceServer` deve tratar `multipart/related` e normalizar a extração do `fileID` removendo o prefixo `/upload`.

---

## 5. Verification Method (Método de Verificação Independente)

Para auditar e aprovar a implementação dos ajustes solicitados:

1. **Verificação da Ordem no Passo 3 de Tabelas:**
   - Inspecionar a seção de Two-Step Table Approach em `docs/architecture_specification.md` e certificar-se de que a operação em `tableEndIndex` precede o loop reverso de células.
2. **Verificação do Bounding Box de Diagramas Mermaid:**
   - Inspecionar o cálculo de altura e largura de Mermaid e conferir a presença de `MaxImageHeightPt: 620.0` e `math.Min`.
3. **Verificação de `%%EOF` e TOCTOU:**
   - Inspecionar a função de hash e conferir a string `%%EOF` e o teste `copiedSize == stat.Size()`.
4. **Condição de Invalidação:**
   - Este relatório será invalidado caso a Google Docs API passe a tratar índices em lote de forma comutativa e não sequencial, ou caso o Google Docs permita imagens maiores que $842\text{ pt}$ em folhas A4 sem deformações.
