> [!WARNING]
> **Developer Preview:** Available as part of the [Google Workspace Developer Preview Program](https://developers.google.com/workspace/preview), which grants early access to certain features.

A [Model Context Protocol (MCP) server](https://modelcontextprotocol.io/docs/learn/server-concepts) acts as a proxy between an external service that provides context, data, or capabilities to a Large Language Model (LLM) or AI application. MCP servers connect AI applications to external systems such as databases and web services, translating their responses into a format that the AI application can understand.

## MCP Tools

An [MCP tool](https://modelcontextprotocol.io/legacy/concepts/tools) is a function or executable capability that an MCP server exposes to a LLM or AI application to perform an action in the real world. Toolsets are a group of tools that are useful for a specific task. Using toolsets can help your agent perform better because it decreases the number of tools available to your agent.

### Toolsets

The docsmcp.googleapis.com MCP server has the following toolsets:

| MCP Toolsets |||
| Endpoint | Description | Tools |
|---|---|---|
| /mcp/v1 | Google Docs Public MCP v1 toolset. | ` - https://developers.google.com/workspace/docs/api/reference/mcp/tools_list/read_doc - https://developers.google.com/workspace/docs/api/reference/mcp/tools_list/update_doc ` |

### Get MCP tool specifications


To get the MCP tool specifications for all tools in an MCP server, use the `tools/list` method. The following example demonstrates how to use `curl` to list all tools and their specifications currently available within the MCP server.

| Curl Request |
|---|
| ```bash curl --location 'https://docsmcp.googleapis.com/mcp/v1' \ --header 'content-type: application/json' \ --header 'accept: application/json, text/event-stream' \ --data '{ "method": "tools/list", "jsonrpc": "2.0", "id": 1 }' ``` |

## Tool: `update_doc`

> [!WARNING]
> **Developer Preview:** Available as part of the [Google Workspace Developer Preview Program](https://developers.google.com/workspace/preview), which grants early access to certain features.

Updates a document using batch update requests. It accepts a `documents.batchUpdate` request (as JSON).

Corresponds to [`documents.batchUpdate`](https://developers.google.com/workspace/docs/api/reference/rest/v1/documents/batchUpdate) in the REST API.

The list of possible updates is:

- `replaceAllText`: Replaces all instances of the specified text. Schema:
  - `replaceText` (string, required): The text that will replace the matched text.
  - Criteria (exactly one required):
  - `containsText` (object): Criteria for matching text:
    - `text` (string, required): The text to search for in the document.
    - `matchCase` (boolean, optional): Whether the search should respect case.
    - `searchByRegex` (boolean, optional): Whether the search text is treated as a regular expression.
  - `tabsCriteria` (object, optional): Criteria used to specify in which tabs the replacement occurs:
  - `tabIds` (array of strings, optional): The list of tab IDs in which the request executes.
- `insertText`: Inserts text at the specified location. Schema:
  - `text` (string, required): The text to be inserted.
  - Insertion location (exactly one required):
  - `location` (object): Inserts text at a specific index:
    - `index` (integer, required): Zero-based index in UTF-16 code units.
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
  - `endOfSegmentLocation` (object): Inserts text at the end of a header, footer, footnote, or document body:
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
- `updateTextStyle`: Updates the text style at the specified range. Schema:
  - `textStyle` (object, required): The styles to set on the text:
  - `bold` (boolean, optional): Whether the text is rendered as bold.
  - `italic` (boolean, optional): Whether the text is italicized.
  - `underline` (boolean, optional): Whether the text is underlined.
  - `strikethrough` (boolean, optional): Whether the text is struck through.
  - `smallCaps` (boolean, optional): Whether the text is in small capital letters.
  - `backgroundColor` (object, optional): Text background color: `{"color": {"rgbColor": {"red": number, "green": number, "blue": number}}}`.
  - `foregroundColor` (object, optional): Text foreground color: `{"color": {"rgbColor": {"red": number, "green": number, "blue": number}}}`.
  - `fontSize` (object, optional): Font size: `{"magnitude": number, "unit": "PT"}`.
  - `weightedFontFamily` (object, optional): Font family and weight: `{"fontFamily": string, "weight": integer}`.
  - `baselineOffset` (string, optional): Vertical offset: `"NONE"`, `"SUPERSCRIPT"`, or `"SUBSCRIPT"`.
  - `link` (object, optional): Hyperlink destination: `{"url": string}` or `{"tabId": string}`.
  - `fields` (string, required): Comma-separated list of fields to update (e.g. `"bold,foregroundColor"` or `"*"`).
  - Insertion location (exactly one required):
  - `range` (object): The range of text to style:
    - `startIndex` (integer, optional): Zero-based start index of this range.
    - `endIndex` (integer, optional): Zero-based end index of this range.
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that contains this range.
- `createParagraphBullets`: Creates bullets for paragraphs. Schema:
  - `range` (object, required): The range to apply the bullet preset to:
  - `startIndex` (integer, optional): Zero-based start index of this range.
  - `endIndex` (integer, optional): Zero-based end index of this range.
  - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
  - `tabId` (string, optional): The tab that contains this range.
  - `bulletPreset` (string, required): Preset pattern of bullet glyphs for lists: `"BULLET_DISC_CIRCLE_SQUARE"`, `"BULLET_DIAMONDX_ARROW3D_SQUARE"`, `"BULLET_CHECKBOX"`, `"BULLET_ARROW_DIAMOND_DISC"`, `"BULLET_STAR_CIRCLE_SQUARE"`, `"BULLET_ARROW3D_CIRCLE_SQUARE"`, `"BULLET_LEFTTRIANGLE_DIAMOND_DISC"`, `"BULLET_DIAMONDX_HOLLOWDIAMOND_SQUARE"`, `"BULLET_DIAMOND_CIRCLE_SQUARE"`, `"NUMBERED_DECIMAL_ALPHA_ROMAN"`, `"NUMBERED_DECIMAL_ALPHA_ROMAN_PARENS"`, `"NUMBERED_DECIMAL_NESTED"`, `"NUMBERED_UPPERALPHA_ALPHA_ROMAN"`, `"NUMBERED_UPPERROMAN_UPPERALPHA_DECIMAL"`, `"NUMBERED_ZERODECIMAL_ALPHA_ROMAN"`.
- `deleteParagraphBullets`: Deletes bullets from paragraphs. Schema:
  - `range` (object, required): The range to delete bullets from:
  - `startIndex` (integer, optional): Zero-based start index of this range.
  - `endIndex` (integer, optional): Zero-based end index of this range.
  - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
  - `tabId` (string, optional): The tab that contains this range.
- `createNamedRange`: Creates a named range. Schema:
  - `name` (string, required): The name of the NamedRange (1 to 256 characters).
  - `range` (object, required): The range to apply the name to:
  - `startIndex` (integer, optional): Zero-based start index of this range.
  - `endIndex` (integer, optional): Zero-based end index of this range.
  - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
  - `tabId` (string, optional): The tab that contains this range.
- `deleteNamedRange`: Deletes a named range. Schema:
  - Named range reference (exactly one required):
  - `namedRangeId` (string): The ID of the named range to delete.
  - `name` (string): The name of the range(s) to delete. All named ranges with this name will be deleted.
  - `tabsCriteria` (object, optional): Criteria used to specify which tab(s) the range deletion applies to:
  - `tabIds` (array of strings, optional): The list of tab IDs.
- `updateParagraphStyle`: Updates the paragraph style at the specified range. Schema:
  - `paragraphStyle` (object, required): The styles to set on the paragraphs:
  - `alignment` (string, optional): Text alignment: `"START"`, `"CENTER"`, `"END"`, `"JUSTIFIED"`.
  - `lineSpacing` (number, optional): Line spacing percentage (e.g. 100 for 100%).
  - `direction` (string, optional): Text direction: `"LEFT_TO_RIGHT"`, `"RIGHT_TO_LEFT"`.
  - `spacingMode` (string, optional): Spacing mode: `"NEVER_COLLAPSE"`, `"COLLAPSE_LISTS"`.
  - `spaceAbove` (object, optional): Space above paragraph: `{"magnitude": number, "unit": "PT"}`.
  - `spaceBelow` (object, optional): Space below paragraph: `{"magnitude": number, "unit": "PT"}`.
  - `borderBetween` / `borderTop` / `borderBottom` / `borderLeft` / `borderRight` (object, optional): `{"color": {"color": {"rgbColor": ...}}, "width": {"magnitude": number, "unit": "PT"}, "padding": {"magnitude": number, "unit": "PT"}, "dashStyle": "SOLID"|"DOT"|"DASH"}`.
  - `indentFirstLine` (object, optional): First line indent: `{"magnitude": number, "unit": "PT"}`.
  - `indentStart` (object, optional): Indent at start: `{"magnitude": number, "unit": "PT"}`.
  - `indentEnd` (object, optional): Indent at end: `{"magnitude": number, "unit": "PT"}`.
  - `namedStyleType` (string, optional): Named style type: `"NORMAL_TEXT"`, `"TITLE"`, `"SUBTITLE"`, `"HEADING_1"`, `"HEADING_2"`, `"HEADING_3"`, `"HEADING_4"`, `"HEADING_5"`, `"HEADING_6"`.
  - `keepWithNext` (boolean, optional): Whether to keep paragraph on same page as next.
  - `keepLinesTogether` (boolean, optional): Whether to keep all lines on same page.
  - `avoidWidowAndOrphan` (boolean, optional): Whether to avoid widow/orphan lines.
  - `shading` (object, optional): Paragraph background shading: `{"backgroundColor": {"color": {"rgbColor": {"red": number, "green": number, "blue": number}}}}`.
  - `pageBreakBefore` (boolean, optional): Whether the paragraph should always start at the beginning of a page.
  - `fields` (string, required): Comma-separated list of fields to update (e.g. `"alignment,namedStyleType"` or `"*"`).
  - Insertion location (exactly one required):
  - `range` (object): The range overlapping the paragraphs to style:
    - `startIndex` (integer, optional): Zero-based start index of this range.
    - `endIndex` (integer, optional): Zero-based end index of this range.
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that contains this range.
- `deleteContentRange`: Deletes content from the document. Schema:
  - `range` (object, required): The range of content to delete:
  - `startIndex` (integer, optional): Zero-based start index of this range.
  - `endIndex` (integer, optional): Zero-based end index of this range.
  - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
  - `tabId` (string, optional): The tab that contains this range.
- `insertInlineImage`: Inserts an inline image at the specified location. Schema:
  - `uri` (string, required): The image URI. Must be publicly accessible and at most 2 kB in length.
  - `objectSize` (object, optional): Size image should appear in document: `{"width": {"magnitude": number, "unit": "PT"}, "height": {"magnitude": number, "unit": "PT"}}`.
  - Insertion location (exactly one required):
  - `location` (object): Inserts image at a specific location:
    - `index` (integer, required): Zero-based index in UTF-16 code units.
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
  - `endOfSegmentLocation` (object): Inserts image at end of header, footer, or document body:
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
- `insertTable`: Inserts a table at the specified location. Schema:
  - `rows` (integer, required): Number of rows in the table.
  - `columns` (integer, required): Number of columns in the table.
  - Insertion location (exactly one required):
  - `location` (object): Inserts table at a specific location:
    - `index` (integer, required): Zero-based index in UTF-16 code units.
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
  - `endOfSegmentLocation` (object): Inserts table at end of header, footer, or document body:
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
- `insertTableRow`: Inserts an empty row into a table. Schema:
  - `tableCellLocation` (object, required): Reference table cell location from which rows will be inserted:
  - `tableStartLocation` (object, required): Location where table starts in document:
    - `index` (integer, required): Zero-based index of table start in UTF-16 code units.
    - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
    - `tabId` (string, optional): Tab that contains the table.
  - `rowIndex` (integer, required): Zero-based row index.
  - `columnIndex` (integer, required): Zero-based column index.
  - `insertBelow` (boolean, required): Whether to insert new row below the reference cell (`true` for below, `false` for above).
- `insertTableColumn`: Inserts an empty column into a table. Schema:
  - `tableCellLocation` (object, required): Reference table cell location from which columns will be inserted:
  - `tableStartLocation` (object, required): Location where table starts in document:
    - `index` (integer, required): Zero-based index of table start in UTF-16 code units.
    - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
    - `tabId` (string, optional): Tab that contains the table.
  - `rowIndex` (integer, required): Zero-based row index.
  - `columnIndex` (integer, required): Zero-based column index.
  - `insertRight` (boolean, required): Whether to insert new column to the right of the reference cell (`true` for right, `false` for left).
- `deleteTableRow`: Deletes a row from a table. Schema:
  - `tableCellLocation` (object, required): Reference table cell location from which row will be deleted:
  - `tableStartLocation` (object, required): Location where table starts in document:
    - `index` (integer, required): Zero-based index of table start in UTF-16 code units.
    - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
    - `tabId` (string, optional): Tab that contains the table.
  - `rowIndex` (integer, required): Zero-based row index.
  - `columnIndex` (integer, required): Zero-based column index.
- `deleteTableColumn`: Deletes a column from a table. Schema:
  - `tableCellLocation` (object, required): Reference table cell location from which column will be deleted:
  - `tableStartLocation` (object, required): Location where table starts in document:
    - `index` (integer, required): Zero-based index of table start in UTF-16 code units.
    - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
    - `tabId` (string, optional): Tab that contains the table.
  - `rowIndex` (integer, required): Zero-based row index.
  - `columnIndex` (integer, required): Zero-based column index.
- `insertPageBreak`: Inserts a page break at the specified location. Schema:
  - Insertion location (exactly one required):
  - `location` (object): Inserts page break at a specific index:
    - `index` (integer, required): Zero-based index in UTF-16 code units.
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
  - `endOfSegmentLocation` (object): Inserts page break at the end of the document body:
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
- `deletePositionedObject`: Deletes a positioned object from the document. Schema:
  - `objectId` (string, required): ID of positioned object to delete.
  - `tabId` (string, optional): Tab that contains the positioned object.
- `updateTableColumnProperties`: Updates the properties of columns in a table. Schema:
  - `tableStartLocation` (object, required): Location where table starts in document:
  - `index` (integer, required): Zero-based index of table start in UTF-16 code units.
  - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
  - `tabId` (string, optional): Tab that contains the table.
  - `columnIndices` (array of integers, optional): Zero-based column indices to update. If omitted, all columns are updated.
  - `tableColumnProperties` (object, required): Column properties to update:
  - `width` (object, optional): Column width: `{"magnitude": number, "unit": "PT"}`.
  - `widthType` (string, optional): `"EVENLY_DISTRIBUTED"`, `"FIXED_WIDTH"`.
  - `fields` (string, required): Comma-separated list of fields to update (e.g. `"width"` or `"*"`).
- `updateTableCellStyle`: Updates the style of table cells. Schema:
  - Cells (exactly one required):
  - `tableRange` (object): Subset of table cells to update:
    - `tableCellLocation` (object, required): Reference cell location:
    - `tableStartLocation` (object, required): Location where table starts in document:
      - `index` (integer, required): Zero-based index of table start in UTF-16 code units.
      - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
      - `tabId` (string, optional): Tab that contains the table.
    - `rowIndex` (integer, required): Zero-based row index.
    - `columnIndex` (integer, required): Zero-based column index.
    - `rowSpan` (integer, required): Row span of table range.
    - `columnSpan` (integer, required): Column span of table range.
  - `tableStartLocation` (object): Location where table starts to apply to all cells in the table:
    - `index` (integer, required): Zero-based index of table start in UTF-16 code units.
    - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
    - `tabId` (string, optional): Tab that contains the table.
  - `tableCellStyle` (object, required): Cell style to set:
  - `backgroundColor` (object, optional): `{"color": {"rgbColor": {"red": number, "green": number, "blue": number}}}`.
  - `paddingTop` / `paddingBottom` / `paddingLeft` / `paddingRight` (object, optional): `{"magnitude": number, "unit": "PT"}`.
  - `contentAlignment` (string, optional): `"TOP"`, `"MIDDLE"`, `"BOTTOM"`.
  - `borderTop` / `borderBottom` / `borderLeft` / `borderRight` (object, optional): `{"color": {"color": {"rgbColor": ...}}, "width": {"magnitude": number, "unit": "PT"}, "dashStyle": "SOLID"|"DOT"|"DASH"}`.
  - `fields` (string, required): Comma-separated list of fields to update (e.g. `"backgroundColor,contentAlignment"` or `"*"`).
- `updateTableRowStyle`: Updates the row style in a table. Schema:
  - `tableStartLocation` (object, required): Location where table starts in document:
  - `index` (integer, required): Zero-based index of table start in UTF-16 code units.
  - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
  - `tabId` (string, optional): Tab that contains the table.
  - `rowIndices` (array of integers, optional): Zero-based row indices to update. If omitted, all rows are updated.
  - `tableRowStyle` (object, required): Row style to set:
  - `minRowHeight` (object, optional): Minimum row height: `{"magnitude": number, "unit": "PT"}`.
  - `preventOverflow` (boolean, optional): Whether row cannot overflow across pages.
  - `fields` (string, required): Comma-separated list of fields to update (e.g. `"minRowHeight"` or `"*"`).
- `replaceImage`: Replaces an image in the document. Schema:
  - `imageObjectId` (string, required): ID of existing image to replace.
  - `uri` (string, required): URI of new image. Must be publicly accessible.
  - `imageReplaceMethod` (string, optional): `"CENTER_CROP"`.
  - `tabId` (string, optional): Tab that contains the image to replace.
- `updateDocumentStyle`: Updates the style of the document. Schema:
  - `documentStyle` (object, required): Document style to set:
  - `background` (object, optional): `{"color": {"color": {"rgbColor": {"red": number, "green": number, "blue": number}}}}`.
  - `marginTop` / `marginBottom` / `marginLeft` / `marginRight` (object, optional): `{"magnitude": number, "unit": "PT"}`.
  - `pageSize` (object, optional): `{"width": {"magnitude": number, "unit": "PT"}, "height": {"magnitude": number, "unit": "PT"}}`.
  - `useCustomHeaderFooterMargins` (boolean, optional): Whether to use custom header/footer margins.
  - `marginHeader` / `marginFooter` (object, optional): `{"magnitude": number, "unit": "PT"}`.
  - `fields` (string, required): Comma-separated list of fields to update (e.g. `"background,marginTop"` or `"*"`).
  - `tabId` (string, optional): Tab containing the document style to update.
- `mergeTableCells`: Merges cells in a table. Schema:
  - `tableRange` (object, required): Table range specifying which cells to merge:
  - `tableCellLocation` (object, required): Location of reference cell:
    - `tableStartLocation` (object, required): Location where table starts in document:
    - `index` (integer, required): Zero-based index of table start in UTF-16 code units.
    - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
    - `tabId` (string, optional): Tab that contains the table.
    - `rowIndex` (integer, required): Zero-based row index.
    - `columnIndex` (integer, required): Zero-based column index.
  - `rowSpan` (integer, required): Row span of table range.
  - `columnSpan` (integer, required): Column span of table range.
- `unmergeTableCells`: Unmerges cells in a table. Schema:
  - `tableRange` (object, required): Table range specifying which cells to unmerge:
  - `tableCellLocation` (object, required): Location of reference cell:
    - `tableStartLocation` (object, required): Location where table starts in document:
    - `index` (integer, required): Zero-based index of table start in UTF-16 code units.
    - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
    - `tabId` (string, optional): Tab that contains the table.
    - `rowIndex` (integer, required): Zero-based row index.
    - `columnIndex` (integer, required): Zero-based column index.
  - `rowSpan` (integer, required): Row span of table range.
  - `columnSpan` (integer, required): Column span of table range.
- `createHeader`: Creates a header. Schema:
  - `type` (string, required): Type of header: `"DEFAULT"`, `"FIRST_PAGE"`.
  - `sectionBreakLocation` (object, optional): Location of section break which begins the section this header belongs to. If omitted, applies to document style:
  - `index` (integer, required): Zero-based index in UTF-16 code units.
  - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
  - `tabId` (string, optional): Tab that contains the section break.
- `createFooter`: Creates a footer. Schema:
  - `type` (string, required): Type of footer: `"DEFAULT"`, `"FIRST_PAGE"`.
  - `sectionBreakLocation` (object, optional): Location of section break preceding the section this footer belongs to. If omitted, applies to document style:
  - `index` (integer, required): Zero-based index in UTF-16 code units.
  - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
  - `tabId` (string, optional): Tab that contains the section break.
- `createFootnote`: Creates a footnote. Schema:
  - Footnote reference location (exactly one required):
  - `location` (object): Inserts footnote reference at a specific index:
    - `index` (integer, required): Zero-based index in UTF-16 code units.
    - `tabId` (string, optional): The tab that the location is in.
  - `endOfSegmentLocation` (object): Inserts footnote reference at the end of the document body:
    - `tabId` (string, optional): The tab that the location is in.
- `replaceNamedRangeContent`: Replaces the content in a named range. Schema:
  - `text` (string, required): Text to replace the content with.
  - Named range reference (exactly one required):
  - `namedRangeId` (string): ID of named range to replace.
  - `namedRangeName` (string): Name of named range(s) to replace.
  - `tabsCriteria` (object, optional):
  - `tabIds` (array of strings, optional): The list of tab IDs in which replacement occurs.
- `updateSectionStyle`: Updates the section style of the specified range. Schema:
  - `range` (object, required): Range overlapping sections to style:
  - `startIndex` (integer, optional): Zero-based start index of range.
  - `endIndex` (integer, optional): Zero-based end index of range.
  - `tabId` (string, optional): Tab that contains this range.
  - `sectionStyle` (object, required): Section style properties:
  - `columnProperties` (array of objects, optional): Properties for each column (up to 3 columns):
    - `paddingEnd` (object, optional): Padding at the end of the column: `{"magnitude": number, "unit": "PT"}`.
    - `width` (object, optional): Column width: `{"magnitude": number, "unit": "PT"}`.
  - `columnSeparatorStyle` (string, optional): `"NONE"`, `"BETWEEN_EACH_COLUMN"`.
  - `contentDirection` (string, optional): `"LEFT_TO_RIGHT"`, `"RIGHT_TO_LEFT"`.
  - `marginTop` / `marginBottom` / `marginLeft` / `marginRight` / `marginHeader` / `marginFooter` (object, optional): `{"magnitude": number, "unit": "PT"}`.
  - `sectionType` (string, optional): `"CONTINUOUS"`, `"NEXT_PAGE"`.
  - `fields` (string, required): Comma-separated list of fields to update (e.g. `"marginTop,sectionType"` or `"*"`).
- `insertSectionBreak`: Inserts a section break at the specified location. Schema:
  - `sectionType` (string, required): Type of section: `"CONTINUOUS"` or `"NEXT_PAGE"`.
  - Insertion location (exactly one required):
  - `location` (object): Inserts section break at a specific index in the document:
    - `index` (integer, required): Zero-based index in UTF-16 code units.
    - `tabId` (string, optional): The tab that the location is in.
  - `endOfSegmentLocation` (object): Inserts section break at the end of the document body:
    - `tabId` (string, optional): The tab that the location is in.
- `deleteHeader`: Deletes a header from the document. Schema:
  - `headerId` (string, required): ID of header to delete.
  - `tabId` (string, optional): Tab containing the header to delete.
- `deleteFooter`: Deletes a footer from the document. Schema:
  - `footerId` (string, required): ID of footer to delete.
  - `tabId` (string, optional): Tab containing the footer to delete.
- `pinTableHeaderRows`: Updates the number of pinned header rows in a table. Schema:
  - `tableStartLocation` (object, required): Location where table starts in document:
  - `index` (integer, required): Zero-based index of table start in UTF-16 code units.
  - `segmentId` (string, optional): ID of header, footer, or footnote. Empty or omitted for document body.
  - `tabId` (string, optional): Tab that contains the table.
  - `pinnedHeaderRowsCount` (integer, required): Number of table rows to pin (0 unpins all rows).
- `addDocumentTab`: Adds a document tab. Schema:
  - `tabProperties` (object, optional): Properties of tab to add:
  - `tabId` (string, optional): User-supplied tab ID.
  - `title` (string, optional): Title of the tab.
  - `index` (integer, optional): Zero-based index where tab will be inserted.
  - `parentTabId` (string, optional): Parent tab ID for nested tabs.
  - `iconEmoji` (string, optional): Emoji icon displayed with the tab.
- `deleteTab`: Deletes a document tab. Schema:
  - `tabId` (string, required): ID of tab to delete.
- `updateDocumentTabProperties`: Updates the properties of a document tab. Schema:
  - `tabProperties` (object, required): Tab properties to update:
  - `tabId` (string, required): ID of tab to update.
  - `title` (string, optional): Updated title of the tab.
  - `index` (integer, optional): Updated zero-based index of the tab.
  - `parentTabId` (string, optional): Updated parent tab ID.
  - `iconEmoji` (string, optional): Updated emoji icon displayed with the tab.
  - `fields` (string, required): Comma-separated list of fields to update (e.g. `"title,iconEmoji"` or `"*"`).
- `insertPerson`: Inserts a person mention. Schema:
  - `personProperties` (object, required): Person properties:
  - `email` (string, required): Email address of the person.
  - `name` (string, optional): Name of the person.
  - Insertion location (exactly one required):
  - `location` (object): Inserts person mention at a specific location:
    - `index` (integer, required): Zero-based index in UTF-16 code units.
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
  - `endOfSegmentLocation` (object): Inserts person mention at end of header, footer, footnote, or document body:
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
- `updateNamedStyle`: Updates a named style. Schema:
  - `namedStyle` (object, required): Named style to update:
  - `namedStyleType` (string, required): Named style type: `"NORMAL_TEXT"`, `"TITLE"`, `"SUBTITLE"`, `"HEADING_1"`, `"HEADING_2"`, `"HEADING_3"`, `"HEADING_4"`, `"HEADING_5"`, `"HEADING_6"`.
  - `textStyle` (object, optional): Text style properties for this named style.
  - `paragraphStyle` (object, optional): Paragraph style properties for this named style.
  - `fields` (string, required): Comma-separated list of fields to update (e.g. `"textStyle,paragraphStyle"` or `"*"`).
  - `tabId` (string, optional): Tab to update.
- `insertRichLink`: Insert a rich link. Schema:
  - `richLinkProperties` (object, required): Rich link properties:
  - `uri` (string, required): URI of the rich link.
  - `title` (string, optional): Title of the rich link.
  - `mimeType` (string, optional): MIME type of the rich link.
  - Insertion location (exactly one required):
  - `location` (object): Inserts rich link at a specific location:
    - `index` (integer, required): Zero-based index in UTF-16 code units.
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
  - `endOfSegmentLocation` (object): Inserts rich link at end of header, footer, footnote, or document body:
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
- `insertDate`: Inserts a date. Schema:
  - `dateElementProperties` (object, required): Date properties:
  - `timestamp` (string, required): UTC timestamp.
  - `timeZoneId` (string, optional): Time zone (e.g. `"America/New_York"`). Defaults to `"etc/UTC"`.
  - `locale` (string, optional): Locale code (e.g. `"en"`).
  - `dateFormat` (string, optional): `"DATE_FORMAT_MONTH_DAY_ABBREVIATED"`, `"DATE_FORMAT_MONTH_DAY_FULL"`, `"DATE_FORMAT_MONTH_DAY_YEAR_ABBREVIATED"`, `"DATE_FORMAT_ISO8601"`.
  - `timeFormat` (string, optional): `"TIME_FORMAT_DISABLED"`, `"TIME_FORMAT_HOUR_MINUTE"`, `"TIME_FORMAT_HOUR_MINUTE_TIMEZONE"`.
  - Insertion location (exactly one required):
  - `location` (object): Inserts date at a specific location:
    - `index` (integer, required): Zero-based index in UTF-16 code units.
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
  - `endOfSegmentLocation` (object): Inserts date at end of header, footer, or document body:
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that the location is in.
- `insertComment`: Inserts a comment into the document. Schema:
  - `content` (string, required): Plain text comment content.
  - `assigneeEmailAddress` (string, optional): The email address of the assignee of the comment. Leave empty for a non-assigned comment.
  - Anchor (exactly one required):
  - `range` (object): The range in the document to anchor the comment to:
    - `startIndex` (integer, optional): The zero-based start index of this range.
    - `endIndex` (integer, optional): The zero-based end index of this range.
    - `segmentId` (string, optional): The ID of the header, footer, or footnote. Empty for document body.
    - `tabId` (string, optional): The tab that contains this range.
- `addCommentReply`: Adds a reply to an existing comment or suggestion thread. Schema:
  - Thread ID (exactly one required):
  - `commentId` (string): ID of the comment thread.
  - `suggestionId` (string): ID of the suggestion thread.
  - `post` (object, required):
  - `content` (string, required unless resolving/reopening): Plain text reply content.
  - `commentAction` (string, optional): Action for comment threads (`"NO_COMMENT_ACTION_CHANGE"`, `"RESOLVE"`, `"REOPEN"`).
  - `assigneeEmail` (string, optional): The email of the user who is being newly assigned to the thread as part of this post.
- `updateCommentPost`: Updates the content of a comment post in a comment thread or suggestion thread. Schema:
  - Thread ID (exactly one required):
  - `commentId` (string): ID of the comment thread.
  - `suggestionId` (string): ID of the suggestion thread.
  - `postId` (string, required): ID of the comment post to update.
  - `content` (string, required): The updated content of the comment post.
- `deleteComment`: Deletes a comment thread. Schema:
  - `commentId` (string, required): ID of the comment thread to delete.
- `deleteCommentReply`: Deletes a reply post from a comment thread or suggestion thread. Schema:
  - Thread ID (exactly one required):
  - `commentId` (string): ID of the comment thread.
  - `suggestionId` (string): ID of the suggestion thread.
  - `postId` (string, required): ID of the reply post to delete.
- `acceptSuggestion`: Accepts a suggestion. Schema:
  - `suggestionId` (string, required): ID of the suggestion to accept.
- `rejectSuggestion`: Rejects a suggestion. Schema:
  - `suggestionId` (string, required): ID of the suggestion to reject.
- `deleteSuggestion`: Deletes a suggestion. Schema:
  - `suggestionId` (string, required): ID of the suggestion to delete.

The following code sample shows how to use `curl` to call the `update_doc` MCP tool.

| Curl Request |
|---|
| ```bash curl --location 'https://docsmcp.googleapis.com/mcp/v1' \ --header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \ --header 'content-type: application/json' \ --header 'accept: application/json, text/event-stream' \ --data '{ "method": "tools/call", "params": { "name": "update_doc", "arguments": { // Provide these details according to the MCP tool specification. } }, "jsonrpc": "2.0", "id": 1 }' ``` |

## Input Schema

### UpdateDocRequest

| JSON representation |
|---|
| ``` { "documentId": string, "requests": [ { object } ], "writeControl": { object (`https://developers.google.com/workspace/docs/api/reference/mcp/tools_list/update_doc#Input.Schema.WriteControl`) } } ``` |

| Fields ||
|---|---|
| `documentId` | `string` Required. The ID of the document to update. This is the same as file_id from Drive tools. |
| `requests[]` | ``object (`https://protobuf.dev/reference/protobuf/google.protobuf#struct` format)`` A list of updates to apply to the document. Each request should be a valid documents.batchUpdate Request object, using the schema documented in: <https://developers.google.com/workspace/docs/api/reference/rest/v1/documents/request> Requests will be applied in the order they are specified. If any request is not valid, no requests will be applied. |
| `writeControl` | ``object (`https://developers.google.com/workspace/docs/api/reference/mcp/tools_list/update_doc#Input.Schema.WriteControl`)`` Optional. Provides control over how write requests are executed. |

### Struct

| JSON representation |
|---|
| ``` { "fields": { string: value, ... } } ``` |

| Fields ||
|---|---|
| `fields` | ``map (key: string, value: value (`https://protobuf.dev/reference/protobuf/google.protobuf#value` format))`` Unordered map of dynamically typed values. An object containing a list of `"key": value` pairs. Example: `{ "name": "wrench", "mass": "1.3kg", "count": "3" }`. |

### FieldsEntry

| JSON representation |
|---|
| ``` { "key": string, "value": value } ``` |

| Fields ||
|---|---|
| `key` | `string` |
| `value` | ``value (`https://protobuf.dev/reference/protobuf/google.protobuf#value` format)`` |

### Value

| JSON representation |
|---|
| ``` { "nullValue": null, "numberValue": number, "stringValue": string, "boolValue": boolean, "structValue": { object }, "listValue": array } ``` |

| Fields ||
|---|---|
| Union field `kind`. The kind of value. `kind` can be only one of the following: ||
| `nullValue` | `null` Represents a JSON `null`. |
| `numberValue` | `number` Represents a JSON number. Must not be `NaN`, `Infinity` or `-Infinity`, since those are not supported in JSON. This also cannot represent large Int64 values, since JSON format generally does not support them in its number type. |
| `stringValue` | `string` Represents a JSON string. |
| `boolValue` | `boolean` Represents a JSON boolean (`true` or `false` literal in JSON). |
| `structValue` | ``object (`https://protobuf.dev/reference/protobuf/google.protobuf#struct` format)`` Represents a JSON object. |
| `listValue` | ``array (`https://protobuf.dev/reference/protobuf/google.protobuf#list-value` format)`` Represents a JSON array. |
|   ||

### ListValue

| JSON representation |
|---|
| ``` { "values": [ value ] } ``` |

| Fields ||
|---|---|
| `values[]` | ``value (`https://protobuf.dev/reference/protobuf/google.protobuf#value` format)`` Repeated field of dynamically typed values. |

### WriteControl

| JSON representation |
|---|
| ``` { "writeMode": enum (`https://developers.google.com/workspace/docs/api/reference/mcp/tools_list/update_doc#Input.Schema.WriteMode`), "requiredRevisionId": string, "targetRevisionId": string } ``` |

| Fields ||
|---|---|
| `writeMode` | ``enum (`https://developers.google.com/workspace/docs/api/reference/mcp/tools_list/update_doc#Input.Schema.WriteMode`)`` Optional. The write mode of the request. |
| Union field `control`. Determines the revision of the document to write to and how the request should behave if that revision is not the current revision of the document. `control` can be only one of the following: ||
| `requiredRevisionId` | `string` The revision ID of the document that the write request will be applied to. If this is not the latest revision of the document, the request will not be processed and will return a 400 bad request error. When a required revision ID is returned in a response, it indicates the revision ID of the document after the request was applied. |
| `targetRevisionId` | `string` The target revision ID of the document that the write request will be applied to. If collaborator changes have occurred after the target revision ID was retrieved, the changes may produce a different result than if they were applied to the target revision ID. When a target revision ID is specified, the write request is applied based on the target revision ID and the revision ID of the document is incremented by 1. If you don't specify a target revision ID, the write request is applied to the latest revision ID and the revision ID of the document is incremented by 1. When a target revision ID is returned in a response, it indicates the revision ID of the document after the request was applied. |
|   ||

### NullValue

Represents a JSON `null`.

`NullValue` is a sentinel, using an enum with only one value to represent the null value for the `Value` type union.

A field of type `NullValue` with any value other than `0` is considered invalid. Most ProtoJSON serializers will emit a `Value` with a `null_value` set as a JSON `null` regardless of the integer value, and so will round trip to a `0` value.

| Enums ||
|---|---|
| `NULL_VALUE` | Null value. |

### WriteMode

The mode for write requests.

| Enums ||
|---|---|
| `WRITE_MODE_UNSPECIFIED` | The default write mode. A request without a write mode specified is treated as `EDIT`. |
| `EDIT` | Changes are applied directly to the document. |
| `SUGGEST` | Changes are made as suggestions. |

## Output Schema

### UpdateDocResponse

| JSON representation |
|---|
| ``` { "replies": [ { object } ], "writeControl": { object (`https://developers.google.com/workspace/docs/api/reference/mcp/tools_list/update_doc#Input.Schema.WriteControl`) } } ``` |

| Fields ||
|---|---|
| `replies[]` | ``object (`https://protobuf.dev/reference/protobuf/google.protobuf#struct` format)`` The replies from executing the batch update requests. Includes errors, warnings, and raw API responses to help the model make adjustments. The raw API responses use the schema documented in: <https://developers.google.com/workspace/docs/api/reference/rest/v1/documents/response> |
| `writeControl` | ``object (`https://developers.google.com/workspace/docs/api/reference/mcp/tools_list/update_doc#Input.Schema.WriteControl`)`` The revision of the document after the updates were applied. Pass `write_control.required_revision_id` in the next `update_doc` request to reject the write if the document changed in between. |

### Struct

| JSON representation |
|---|
| ``` { "fields": { string: value, ... } } ``` |

| Fields ||
|---|---|
| `fields` | ``map (key: string, value: value (`https://protobuf.dev/reference/protobuf/google.protobuf#value` format))`` Unordered map of dynamically typed values. An object containing a list of `"key": value` pairs. Example: `{ "name": "wrench", "mass": "1.3kg", "count": "3" }`. |

### FieldsEntry

| JSON representation |
|---|
| ``` { "key": string, "value": value } ``` |

| Fields ||
|---|---|
| `key` | `string` |
| `value` | ``value (`https://protobuf.dev/reference/protobuf/google.protobuf#value` format)`` |

### Value

| JSON representation |
|---|
| ``` { "nullValue": null, "numberValue": number, "stringValue": string, "boolValue": boolean, "structValue": { object }, "listValue": array } ``` |

| Fields ||
|---|---|
| Union field `kind`. The kind of value. `kind` can be only one of the following: ||
| `nullValue` | `null` Represents a JSON `null`. |
| `numberValue` | `number` Represents a JSON number. Must not be `NaN`, `Infinity` or `-Infinity`, since those are not supported in JSON. This also cannot represent large Int64 values, since JSON format generally does not support them in its number type. |
| `stringValue` | `string` Represents a JSON string. |
| `boolValue` | `boolean` Represents a JSON boolean (`true` or `false` literal in JSON). |
| `structValue` | ``object (`https://protobuf.dev/reference/protobuf/google.protobuf#struct` format)`` Represents a JSON object. |
| `listValue` | ``array (`https://protobuf.dev/reference/protobuf/google.protobuf#list-value` format)`` Represents a JSON array. |
|   ||

### ListValue

| JSON representation |
|---|
| ``` { "values": [ value ] } ``` |

| Fields ||
|---|---|
| `values[]` | ``value (`https://protobuf.dev/reference/protobuf/google.protobuf#value` format)`` Repeated field of dynamically typed values. |

### WriteControl

| JSON representation |
|---|
| ``` { "writeMode": enum (`https://developers.google.com/workspace/docs/api/reference/mcp/tools_list/update_doc#Input.Schema.WriteMode`), "requiredRevisionId": string, "targetRevisionId": string } ``` |

| Fields ||
|---|---|
| `writeMode` | ``enum (`https://developers.google.com/workspace/docs/api/reference/mcp/tools_list/update_doc#Input.Schema.WriteMode`)`` Optional. The write mode of the request. |
| Union field `control`. Determines the revision of the document to write to and how the request should behave if that revision is not the current revision of the document. `control` can be only one of the following: ||
| `requiredRevisionId` | `string` The revision ID of the document that the write request will be applied to. If this is not the latest revision of the document, the request will not be processed and will return a 400 bad request error. When a required revision ID is returned in a response, it indicates the revision ID of the document after the request was applied. |
| `targetRevisionId` | `string` The target revision ID of the document that the write request will be applied to. If collaborator changes have occurred after the target revision ID was retrieved, the changes may produce a different result than if they were applied to the target revision ID. When a target revision ID is specified, the write request is applied based on the target revision ID and the revision ID of the document is incremented by 1. If you don't specify a target revision ID, the write request is applied to the latest revision ID and the revision ID of the document is incremented by 1. When a target revision ID is returned in a response, it indicates the revision ID of the document after the request was applied. |
|   ||

### NullValue

Represents a JSON `null`.

`NullValue` is a sentinel, using an enum with only one value to represent the null value for the `Value` type union.

A field of type `NullValue` with any value other than `0` is considered invalid. Most ProtoJSON serializers will emit a `Value` with a `null_value` set as a JSON `null` regardless of the integer value, and so will round trip to a `0` value.

| Enums ||
|---|---|
| `NULL_VALUE` | Null value. |

### WriteMode

The mode for write requests.

| Enums ||
|---|---|
| `WRITE_MODE_UNSPECIFIED` | The default write mode. A request without a write mode specified is treated as `EDIT`. |
| `EDIT` | Changes are applied directly to the document. |
| `SUGGEST` | Changes are made as suggestions. |

### Tool Annotations


[Tool annotations](https://modelcontextprotocol.io/specification/latest/schema#toolannotations) are sent to MCP clients to describe the basic risk of a given tool. Most clients treat these hints as untrusted, but they can be used to decide when a confirmation prompt might be sent to a user.


Along with the title string, the following boolean hints are defined as follows:

- `readOnlyHint`: If true, the tool doesn't modify its environment. Default: false.
- `destructiveHint`: If true, then the tool can perform destructive actions. If false, then the tool can only perform additive actions. Default: true.
- `idempotentHint`: If true, then calling the tool repeatedly with the same arguments will have no additional effect on its environment. Default: false.
- `openWorldHint`: If true, then the tool can interact with an 'open world' of external entities. If false, then the tool can only interact with internal entities. For example, a web search tool would be open world, while a memory tool would not be open world.

Destructive Hint: ❌ \| Idempotent Hint: ❌ \| Read Only Hint: ❌ \| Open World Hint: ✅

## Authorization Scopes

Requires one of the following OAuth scopes:

- `https://www.googleapis.com/auth/drive`
- `https://www.googleapis.com/auth/drive.file`
- `https://www.googleapis.com/auth/documents`


## Tool: `read_doc`

> [!WARNING]
> **Developer Preview:** Available as part of the [Google Workspace Developer Preview Program](https://developers.google.com/workspace/preview), which grants early access to certain features.

This tool retrieves a JSON representation of the Google Doc given its document ID.

The JSON representation includes both the text as well as structural information about the document.

Corresponds to [`documents.get`](https://developers.google.com/workspace/docs/api/reference/rest/v1/documents/get) in the REST API.

The following code sample shows how to use `curl` to call the `read_doc` MCP tool.

| Curl Request |
|---|
| ```bash curl --location 'https://docsmcp.googleapis.com/mcp/v1' \ --header 'Authorization: Bearer YOUR_ACCESS_TOKEN' \ --header 'content-type: application/json' \ --header 'accept: application/json, text/event-stream' \ --data '{ "method": "tools/call", "params": { "name": "read_doc", "arguments": { // Provide these details according to the MCP tool specification. } }, "jsonrpc": "2.0", "id": 1 }' ``` |

## Input Schema

### ReadDocRequest

| JSON representation |
|---|
| ``` { "documentId": string, "commentsIncluded": boolean } ``` |

| Fields ||
|---|---|
| `documentId` | `string` Required. The ID of the document to read. This is the same as file_id from Drive tools. |
| Union field `_comments_included`. `_comments_included` can be only one of the following: ||
| `commentsIncluded` | `boolean` Optional. If true, comments will be included in the response. By default, comments are not included. |
|   ||

## Output Schema

### ReadDocResponse

| JSON representation |
|---|
| ``` { "content": { object } } ``` |

| Fields ||
|---|---|
| `content` | ``object (`https://protobuf.dev/reference/protobuf/google.protobuf#struct` format)`` The verbalized text content of the document. |

### Struct

| JSON representation |
|---|
| ``` { "fields": { string: value, ... } } ``` |

| Fields ||
|---|---|
| `fields` | ``map (key: string, value: value (`https://protobuf.dev/reference/protobuf/google.protobuf#value` format))`` Unordered map of dynamically typed values. An object containing a list of `"key": value` pairs. Example: `{ "name": "wrench", "mass": "1.3kg", "count": "3" }`. |

### FieldsEntry

| JSON representation |
|---|
| ``` { "key": string, "value": value } ``` |

| Fields ||
|---|---|
| `key` | `string` |
| `value` | ``value (`https://protobuf.dev/reference/protobuf/google.protobuf#value` format)`` |

### Value

| JSON representation |
|---|
| ``` { "nullValue": null, "numberValue": number, "stringValue": string, "boolValue": boolean, "structValue": { object }, "listValue": array } ``` |

| Fields ||
|---|---|
| Union field `kind`. The kind of value. `kind` can be only one of the following: ||
| `nullValue` | `null` Represents a JSON `null`. |
| `numberValue` | `number` Represents a JSON number. Must not be `NaN`, `Infinity` or `-Infinity`, since those are not supported in JSON. This also cannot represent large Int64 values, since JSON format generally does not support them in its number type. |
| `stringValue` | `string` Represents a JSON string. |
| `boolValue` | `boolean` Represents a JSON boolean (`true` or `false` literal in JSON). |
| `structValue` | ``object (`https://protobuf.dev/reference/protobuf/google.protobuf#struct` format)`` Represents a JSON object. |
| `listValue` | ``array (`https://protobuf.dev/reference/protobuf/google.protobuf#list-value` format)`` Represents a JSON array. |
|   ||

### ListValue

| JSON representation |
|---|
| ``` { "values": [ value ] } ``` |

| Fields ||
|---|---|
| `values[]` | ``value (`https://protobuf.dev/reference/protobuf/google.protobuf#value` format)`` Repeated field of dynamically typed values. |

### NullValue

Represents a JSON `null`.

`NullValue` is a sentinel, using an enum with only one value to represent the null value for the `Value` type union.

A field of type `NullValue` with any value other than `0` is considered invalid. Most ProtoJSON serializers will emit a `Value` with a `null_value` set as a JSON `null` regardless of the integer value, and so will round trip to a `0` value.

| Enums ||
|---|---|
| `NULL_VALUE` | Null value. |

### Tool Annotations


[Tool annotations](https://modelcontextprotocol.io/specification/latest/schema#toolannotations) are sent to MCP clients to describe the basic risk of a given tool. Most clients treat these hints as untrusted, but they can be used to decide when a confirmation prompt might be sent to a user.


Along with the title string, the following boolean hints are defined as follows:

- `readOnlyHint`: If true, the tool doesn't modify its environment. Default: false.
- `destructiveHint`: If true, then the tool can perform destructive actions. If false, then the tool can only perform additive actions. Default: true.
- `idempotentHint`: If true, then calling the tool repeatedly with the same arguments will have no additional effect on its environment. Default: false.
- `openWorldHint`: If true, then the tool can interact with an 'open world' of external entities. If false, then the tool can only interact with internal entities. For example, a web search tool would be open world, while a memory tool would not be open world.

Destructive Hint: ❌ \| Idempotent Hint: ✅ \| Read Only Hint: ✅ \| Open World Hint: ✅

## Authorization Scopes

Requires one of the following OAuth scopes:

- `https://www.googleapis.com/auth/drive.readonly`
- `https://www.googleapis.com/auth/documents.readonly`
- `https://www.googleapis.com/auth/drive`
- `https://www.googleapis.com/auth/drive.file`
- `https://www.googleapis.com/auth/documents`