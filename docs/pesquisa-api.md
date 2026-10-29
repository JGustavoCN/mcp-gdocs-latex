
Take the 2026 Go Developer Survey
Go
Skip to Main Content
Search packages or symbols

Why Gosubmenu dropdown icon
Learn
Docssubmenu dropdown icon
Packages
Communitysubmenu dropdown icon
Discover Packages
 
google.golang.org/api
 
docs
 
v1

docs
package
Version: v0.310.0 Latest 
Published: Oct 28, 2026 
License: BSD-3-Clause 
Imports: 18 
Imported by: 112
Details
checked Valid go.mod file 
checked Redistributable license 
checked Tagged version 
unchecked Stable version 
Learn more about best practices
Repository
github.com/googleapis/google-api-go-client
Links
Open Source Insights Logo Open Source Insights
Jump to ...
 Documentation ¶
Overview ¶
Library status
Creating a client
Other authentication options
Package docs provides access to the Google Docs API.

For product documentation, see: https://developers.google.com/workspace/docs/

Library status ¶
These client libraries are officially supported by Google. However, this library is considered complete and is in maintenance mode. This means that we will address critical bugs and security issues but will not add any new features.

When possible, we recommend using our newer [Cloud Client Libraries for Go](https://pkg.go.dev/cloud.google.com/go) that are still actively being worked and iterated on.

Creating a client ¶
Usage example:

import "google.golang.org/api/docs/v1"
...
ctx := context.Background()
docsService, err := docs.NewService(ctx)
In this example, Google Application Default Credentials are used for authentication. For information on how to create and obtain Application Default Credentials, see https://developers.google.com/identity/protocols/application-default-credentials.

Other authentication options ¶
By default, all available scopes (see "Constants") are used to authenticate. To restrict scopes, use google.golang.org/api/option.WithScopes:

docsService, err := docs.NewService(ctx, option.WithScopes(docs.DriveReadonlyScope))
To use an API key for authentication (note: some APIs do not support API keys), use google.golang.org/api/option.WithAPIKey:

docsService, err := docs.NewService(ctx, option.WithAPIKey("AIza..."))
To use an OAuth token (e.g., a user token obtained via a three-legged OAuth flow, use google.golang.org/api/option.WithTokenSource:

config := &oauth2.Config{...}
// ...
token, err := config.Exchange(ctx, ...)
docsService, err := docs.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
See google.golang.org/api/option.ClientOption for details on options.

Index ¶
Constants
type AcceptSuggestionRequest
func (s AcceptSuggestionRequest) MarshalJSON() ([]byte, error)
type AddCommentReplyRequest
func (s AddCommentReplyRequest) MarshalJSON() ([]byte, error)
type AddCommentReplyResponse
func (s AddCommentReplyResponse) MarshalJSON() ([]byte, error)
type AddDocumentTabRequest
func (s AddDocumentTabRequest) MarshalJSON() ([]byte, error)
type AddDocumentTabResponse
func (s AddDocumentTabResponse) MarshalJSON() ([]byte, error)
type AutoText
func (s AutoText) MarshalJSON() ([]byte, error)
type Background
func (s Background) MarshalJSON() ([]byte, error)
type BackgroundSuggestionState
func (s BackgroundSuggestionState) MarshalJSON() ([]byte, error)
type BatchUpdateDocumentRequest
func (s BatchUpdateDocumentRequest) MarshalJSON() ([]byte, error)
type BatchUpdateDocumentResponse
func (s BatchUpdateDocumentResponse) MarshalJSON() ([]byte, error)
type Body
func (s Body) MarshalJSON() ([]byte, error)
type BookmarkLink
func (s BookmarkLink) MarshalJSON() ([]byte, error)
type Bullet
func (s Bullet) MarshalJSON() ([]byte, error)
type BulletSuggestionState
func (s BulletSuggestionState) MarshalJSON() ([]byte, error)
type Color
func (s Color) MarshalJSON() ([]byte, error)
type ColumnBreak
func (s ColumnBreak) MarshalJSON() ([]byte, error)
type CommentAnchor
func (s CommentAnchor) MarshalJSON() ([]byte, error)
type CommentThread
func (s CommentThread) MarshalJSON() ([]byte, error)
type CreateDropdownDefinitionRequest
func (s CreateDropdownDefinitionRequest) MarshalJSON() ([]byte, error)
type CreateDropdownDefinitionResponse
func (s CreateDropdownDefinitionResponse) MarshalJSON() ([]byte, error)
type CreateFooterRequest
func (s CreateFooterRequest) MarshalJSON() ([]byte, error)
type CreateFooterResponse
func (s CreateFooterResponse) MarshalJSON() ([]byte, error)
type CreateFootnoteRequest
func (s CreateFootnoteRequest) MarshalJSON() ([]byte, error)
type CreateFootnoteResponse
func (s CreateFootnoteResponse) MarshalJSON() ([]byte, error)
type CreateHeaderRequest
func (s CreateHeaderRequest) MarshalJSON() ([]byte, error)
type CreateHeaderResponse
func (s CreateHeaderResponse) MarshalJSON() ([]byte, error)
type CreateNamedRangeRequest
func (s CreateNamedRangeRequest) MarshalJSON() ([]byte, error)
type CreateNamedRangeResponse
func (s CreateNamedRangeResponse) MarshalJSON() ([]byte, error)
type CreateParagraphBulletsRequest
func (s CreateParagraphBulletsRequest) MarshalJSON() ([]byte, error)
type CropProperties
func (s CropProperties) MarshalJSON() ([]byte, error)
func (s *CropProperties) UnmarshalJSON(data []byte) error
type CropPropertiesSuggestionState
func (s CropPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type DateElement
func (s DateElement) MarshalJSON() ([]byte, error)
type DateElementProperties
func (s DateElementProperties) MarshalJSON() ([]byte, error)
type DateElementPropertiesSuggestionState
func (s DateElementPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type DeleteCommentReplyRequest
func (s DeleteCommentReplyRequest) MarshalJSON() ([]byte, error)
type DeleteCommentRequest
func (s DeleteCommentRequest) MarshalJSON() ([]byte, error)
type DeleteContentRangeRequest
func (s DeleteContentRangeRequest) MarshalJSON() ([]byte, error)
type DeleteDropdownDefinitionRequest
func (s DeleteDropdownDefinitionRequest) MarshalJSON() ([]byte, error)
type DeleteFooterRequest
func (s DeleteFooterRequest) MarshalJSON() ([]byte, error)
type DeleteHeaderRequest
func (s DeleteHeaderRequest) MarshalJSON() ([]byte, error)
type DeleteNamedRangeRequest
func (s DeleteNamedRangeRequest) MarshalJSON() ([]byte, error)
type DeleteParagraphBulletsRequest
func (s DeleteParagraphBulletsRequest) MarshalJSON() ([]byte, error)
type DeletePositionedObjectRequest
func (s DeletePositionedObjectRequest) MarshalJSON() ([]byte, error)
type DeleteSuggestionRequest
func (s DeleteSuggestionRequest) MarshalJSON() ([]byte, error)
type DeleteTabRequest
func (s DeleteTabRequest) MarshalJSON() ([]byte, error)
type DeleteTableColumnRequest
func (s DeleteTableColumnRequest) MarshalJSON() ([]byte, error)
type DeleteTableRowRequest
func (s DeleteTableRowRequest) MarshalJSON() ([]byte, error)
type Dimension
func (s Dimension) MarshalJSON() ([]byte, error)
func (s *Dimension) UnmarshalJSON(data []byte) error
type Document
func (s Document) MarshalJSON() ([]byte, error)
type DocumentFormat
func (s DocumentFormat) MarshalJSON() ([]byte, error)
type DocumentStyle
func (s DocumentStyle) MarshalJSON() ([]byte, error)
type DocumentStyleSuggestionState
func (s DocumentStyleSuggestionState) MarshalJSON() ([]byte, error)
type DocumentTab
func (s DocumentTab) MarshalJSON() ([]byte, error)
type DocumentsBatchUpdateCall
func (c *DocumentsBatchUpdateCall) Context(ctx context.Context) *DocumentsBatchUpdateCall
func (c *DocumentsBatchUpdateCall) Do(opts ...googleapi.CallOption) (*BatchUpdateDocumentResponse, error)
func (c *DocumentsBatchUpdateCall) Fields(s ...googleapi.Field) *DocumentsBatchUpdateCall
func (c *DocumentsBatchUpdateCall) Header() http.Header
type DocumentsCreateCall
func (c *DocumentsCreateCall) Context(ctx context.Context) *DocumentsCreateCall
func (c *DocumentsCreateCall) Do(opts ...googleapi.CallOption) (*Document, error)
func (c *DocumentsCreateCall) Fields(s ...googleapi.Field) *DocumentsCreateCall
func (c *DocumentsCreateCall) Header() http.Header
type DocumentsGetCall
func (c *DocumentsGetCall) CommentsViewMode(commentsViewMode string) *DocumentsGetCall
func (c *DocumentsGetCall) Context(ctx context.Context) *DocumentsGetCall
func (c *DocumentsGetCall) Do(opts ...googleapi.CallOption) (*Document, error)
func (c *DocumentsGetCall) Fields(s ...googleapi.Field) *DocumentsGetCall
func (c *DocumentsGetCall) Header() http.Header
func (c *DocumentsGetCall) IfNoneMatch(entityTag string) *DocumentsGetCall
func (c *DocumentsGetCall) IncludeTabsContent(includeTabsContent bool) *DocumentsGetCall
func (c *DocumentsGetCall) SuggestionsViewMode(suggestionsViewMode string) *DocumentsGetCall
type DocumentsService
func NewDocumentsService(s *Service) *DocumentsService
func (r *DocumentsService) BatchUpdate(documentId string, batchupdatedocumentrequest *BatchUpdateDocumentRequest) *DocumentsBatchUpdateCall
func (r *DocumentsService) Create(document *Document) *DocumentsCreateCall
func (r *DocumentsService) Get(documentId string) *DocumentsGetCall
type Dropdown
func (s Dropdown) MarshalJSON() ([]byte, error)
type DropdownDefinition
func (s DropdownDefinition) MarshalJSON() ([]byte, error)
type DropdownDefinitionProperties
func (s DropdownDefinitionProperties) MarshalJSON() ([]byte, error)
type DropdownDefinitionPropertiesSuggestionState
func (s DropdownDefinitionPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type DropdownOption
func (s DropdownOption) MarshalJSON() ([]byte, error)
type DropdownProperties
func (s DropdownProperties) MarshalJSON() ([]byte, error)
type DropdownPropertiesSuggestionState
func (s DropdownPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type EmbeddedDrawingProperties
type EmbeddedDrawingPropertiesSuggestionState
type EmbeddedObject
func (s EmbeddedObject) MarshalJSON() ([]byte, error)
type EmbeddedObjectBorder
func (s EmbeddedObjectBorder) MarshalJSON() ([]byte, error)
type EmbeddedObjectBorderSuggestionState
func (s EmbeddedObjectBorderSuggestionState) MarshalJSON() ([]byte, error)
type EmbeddedObjectSuggestionState
func (s EmbeddedObjectSuggestionState) MarshalJSON() ([]byte, error)
type EndOfSegmentLocation
func (s EndOfSegmentLocation) MarshalJSON() ([]byte, error)
type Equation
func (s Equation) MarshalJSON() ([]byte, error)
type Footer
func (s Footer) MarshalJSON() ([]byte, error)
type Footnote
func (s Footnote) MarshalJSON() ([]byte, error)
type FootnoteReference
func (s FootnoteReference) MarshalJSON() ([]byte, error)
type Header
func (s Header) MarshalJSON() ([]byte, error)
type HeadingLink
func (s HeadingLink) MarshalJSON() ([]byte, error)
type HorizontalRule
func (s HorizontalRule) MarshalJSON() ([]byte, error)
type ImageProperties
func (s ImageProperties) MarshalJSON() ([]byte, error)
func (s *ImageProperties) UnmarshalJSON(data []byte) error
type ImagePropertiesSuggestionState
func (s ImagePropertiesSuggestionState) MarshalJSON() ([]byte, error)
type InlineObject
func (s InlineObject) MarshalJSON() ([]byte, error)
type InlineObjectElement
func (s InlineObjectElement) MarshalJSON() ([]byte, error)
type InlineObjectProperties
func (s InlineObjectProperties) MarshalJSON() ([]byte, error)
type InlineObjectPropertiesSuggestionState
func (s InlineObjectPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type InsertCommentRequest
func (s InsertCommentRequest) MarshalJSON() ([]byte, error)
type InsertCommentResponse
func (s InsertCommentResponse) MarshalJSON() ([]byte, error)
type InsertDateRequest
func (s InsertDateRequest) MarshalJSON() ([]byte, error)
type InsertDropdownRequest
func (s InsertDropdownRequest) MarshalJSON() ([]byte, error)
type InsertDropdownResponse
func (s InsertDropdownResponse) MarshalJSON() ([]byte, error)
type InsertInlineImageRequest
func (s InsertInlineImageRequest) MarshalJSON() ([]byte, error)
type InsertInlineImageResponse
func (s InsertInlineImageResponse) MarshalJSON() ([]byte, error)
type InsertInlineSheetsChartResponse
func (s InsertInlineSheetsChartResponse) MarshalJSON() ([]byte, error)
type InsertPageBreakRequest
func (s InsertPageBreakRequest) MarshalJSON() ([]byte, error)
type InsertPersonRequest
func (s InsertPersonRequest) MarshalJSON() ([]byte, error)
type InsertRichLinkRequest
func (s InsertRichLinkRequest) MarshalJSON() ([]byte, error)
type InsertSectionBreakRequest
func (s InsertSectionBreakRequest) MarshalJSON() ([]byte, error)
type InsertTableColumnRequest
func (s InsertTableColumnRequest) MarshalJSON() ([]byte, error)
type InsertTableRequest
func (s InsertTableRequest) MarshalJSON() ([]byte, error)
type InsertTableRowRequest
func (s InsertTableRowRequest) MarshalJSON() ([]byte, error)
type InsertTextRequest
func (s InsertTextRequest) MarshalJSON() ([]byte, error)
type Link
func (s Link) MarshalJSON() ([]byte, error)
type LinkedContentReference
func (s LinkedContentReference) MarshalJSON() ([]byte, error)
type LinkedContentReferenceSuggestionState
func (s LinkedContentReferenceSuggestionState) MarshalJSON() ([]byte, error)
type List
func (s List) MarshalJSON() ([]byte, error)
type ListProperties
func (s ListProperties) MarshalJSON() ([]byte, error)
type ListPropertiesSuggestionState
func (s ListPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type Location
func (s Location) MarshalJSON() ([]byte, error)
type MergeTableCellsRequest
func (s MergeTableCellsRequest) MarshalJSON() ([]byte, error)
type NamedRange
func (s NamedRange) MarshalJSON() ([]byte, error)
type NamedRanges
func (s NamedRanges) MarshalJSON() ([]byte, error)
type NamedStyle
func (s NamedStyle) MarshalJSON() ([]byte, error)
type NamedStyleSuggestionState
func (s NamedStyleSuggestionState) MarshalJSON() ([]byte, error)
type NamedStyles
func (s NamedStyles) MarshalJSON() ([]byte, error)
type NamedStylesSuggestionState
func (s NamedStylesSuggestionState) MarshalJSON() ([]byte, error)
type NestingLevel
func (s NestingLevel) MarshalJSON() ([]byte, error)
type NestingLevelSuggestionState
func (s NestingLevelSuggestionState) MarshalJSON() ([]byte, error)
type ObjectReferences
func (s ObjectReferences) MarshalJSON() ([]byte, error)
type OptionalColor
func (s OptionalColor) MarshalJSON() ([]byte, error)
type PageBreak
func (s PageBreak) MarshalJSON() ([]byte, error)
type Paragraph
func (s Paragraph) MarshalJSON() ([]byte, error)
type ParagraphBorder
func (s ParagraphBorder) MarshalJSON() ([]byte, error)
type ParagraphElement
func (s ParagraphElement) MarshalJSON() ([]byte, error)
type ParagraphStyle
func (s ParagraphStyle) MarshalJSON() ([]byte, error)
func (s *ParagraphStyle) UnmarshalJSON(data []byte) error
type ParagraphStyleSuggestionState
func (s ParagraphStyleSuggestionState) MarshalJSON() ([]byte, error)
type Person
func (s Person) MarshalJSON() ([]byte, error)
type PersonProperties
func (s PersonProperties) MarshalJSON() ([]byte, error)
type PinTableHeaderRowsRequest
func (s PinTableHeaderRowsRequest) MarshalJSON() ([]byte, error)
type PositionedObject
func (s PositionedObject) MarshalJSON() ([]byte, error)
type PositionedObjectPositioning
func (s PositionedObjectPositioning) MarshalJSON() ([]byte, error)
type PositionedObjectPositioningSuggestionState
func (s PositionedObjectPositioningSuggestionState) MarshalJSON() ([]byte, error)
type PositionedObjectProperties
func (s PositionedObjectProperties) MarshalJSON() ([]byte, error)
type PositionedObjectPropertiesSuggestionState
func (s PositionedObjectPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type Post
func (s Post) MarshalJSON() ([]byte, error)
type PostAuthor
func (s PostAuthor) MarshalJSON() ([]byte, error)
type Range
func (s Range) MarshalJSON() ([]byte, error)
type RejectSuggestionRequest
func (s RejectSuggestionRequest) MarshalJSON() ([]byte, error)
type ReplaceAllTextRequest
func (s ReplaceAllTextRequest) MarshalJSON() ([]byte, error)
type ReplaceAllTextResponse
func (s ReplaceAllTextResponse) MarshalJSON() ([]byte, error)
type ReplaceImageRequest
func (s ReplaceImageRequest) MarshalJSON() ([]byte, error)
type ReplaceNamedRangeContentRequest
func (s ReplaceNamedRangeContentRequest) MarshalJSON() ([]byte, error)
type Request
func (s Request) MarshalJSON() ([]byte, error)
type Response
func (s Response) MarshalJSON() ([]byte, error)
type RgbColor
func (s RgbColor) MarshalJSON() ([]byte, error)
func (s *RgbColor) UnmarshalJSON(data []byte) error
type RichLink
func (s RichLink) MarshalJSON() ([]byte, error)
type RichLinkProperties
func (s RichLinkProperties) MarshalJSON() ([]byte, error)
type SectionBreak
func (s SectionBreak) MarshalJSON() ([]byte, error)
type SectionColumnProperties
func (s SectionColumnProperties) MarshalJSON() ([]byte, error)
type SectionStyle
func (s SectionStyle) MarshalJSON() ([]byte, error)
type Service
func New(client *http.Client) (*Service, error)deprecated
func NewService(ctx context.Context, opts ...option.ClientOption) (*Service, error)
type Shading
func (s Shading) MarshalJSON() ([]byte, error)
type ShadingSuggestionState
func (s ShadingSuggestionState) MarshalJSON() ([]byte, error)
type SheetsChartReference
func (s SheetsChartReference) MarshalJSON() ([]byte, error)
type SheetsChartReferenceSuggestionState
func (s SheetsChartReferenceSuggestionState) MarshalJSON() ([]byte, error)
type Size
func (s Size) MarshalJSON() ([]byte, error)
type SizeSuggestionState
func (s SizeSuggestionState) MarshalJSON() ([]byte, error)
type StructuralElement
func (s StructuralElement) MarshalJSON() ([]byte, error)
type SubstringMatchCriteria
func (s SubstringMatchCriteria) MarshalJSON() ([]byte, error)
type SuggestedBullet
func (s SuggestedBullet) MarshalJSON() ([]byte, error)
type SuggestedDateElementProperties
func (s SuggestedDateElementProperties) MarshalJSON() ([]byte, error)
type SuggestedDocumentStyle
func (s SuggestedDocumentStyle) MarshalJSON() ([]byte, error)
type SuggestedDropdownDefinitionProperties
func (s SuggestedDropdownDefinitionProperties) MarshalJSON() ([]byte, error)
type SuggestedDropdownProperties
func (s SuggestedDropdownProperties) MarshalJSON() ([]byte, error)
type SuggestedInlineObjectProperties
func (s SuggestedInlineObjectProperties) MarshalJSON() ([]byte, error)
type SuggestedListProperties
func (s SuggestedListProperties) MarshalJSON() ([]byte, error)
type SuggestedNamedStyles
func (s SuggestedNamedStyles) MarshalJSON() ([]byte, error)
type SuggestedParagraphStyle
func (s SuggestedParagraphStyle) MarshalJSON() ([]byte, error)
type SuggestedPositionedObjectProperties
func (s SuggestedPositionedObjectProperties) MarshalJSON() ([]byte, error)
type SuggestedTableCellStyle
func (s SuggestedTableCellStyle) MarshalJSON() ([]byte, error)
type SuggestedTableRowStyle
func (s SuggestedTableRowStyle) MarshalJSON() ([]byte, error)
type SuggestedTextStyle
func (s SuggestedTextStyle) MarshalJSON() ([]byte, error)
type SuggestionResponse
func (s SuggestionResponse) MarshalJSON() ([]byte, error)
type SuggestionThread
func (s SuggestionThread) MarshalJSON() ([]byte, error)
type Tab
func (s Tab) MarshalJSON() ([]byte, error)
type TabProperties
func (s TabProperties) MarshalJSON() ([]byte, error)
type TabStop
func (s TabStop) MarshalJSON() ([]byte, error)
type Table
func (s Table) MarshalJSON() ([]byte, error)
type TableCell
func (s TableCell) MarshalJSON() ([]byte, error)
type TableCellBorder
func (s TableCellBorder) MarshalJSON() ([]byte, error)
type TableCellLocation
func (s TableCellLocation) MarshalJSON() ([]byte, error)
type TableCellStyle
func (s TableCellStyle) MarshalJSON() ([]byte, error)
type TableCellStyleSuggestionState
func (s TableCellStyleSuggestionState) MarshalJSON() ([]byte, error)
type TableColumnProperties
func (s TableColumnProperties) MarshalJSON() ([]byte, error)
type TableOfContents
func (s TableOfContents) MarshalJSON() ([]byte, error)
type TableRange
func (s TableRange) MarshalJSON() ([]byte, error)
type TableRow
func (s TableRow) MarshalJSON() ([]byte, error)
type TableRowStyle
func (s TableRowStyle) MarshalJSON() ([]byte, error)
type TableRowStyleSuggestionState
func (s TableRowStyleSuggestionState) MarshalJSON() ([]byte, error)
type TableStyle
func (s TableStyle) MarshalJSON() ([]byte, error)
type TabsCriteria
func (s TabsCriteria) MarshalJSON() ([]byte, error)
type TextRun
func (s TextRun) MarshalJSON() ([]byte, error)
type TextStyle
func (s TextStyle) MarshalJSON() ([]byte, error)
type TextStyleSuggestionState
func (s TextStyleSuggestionState) MarshalJSON() ([]byte, error)
type UnmergeTableCellsRequest
func (s UnmergeTableCellsRequest) MarshalJSON() ([]byte, error)
type UpdateCommentPostRequest
func (s UpdateCommentPostRequest) MarshalJSON() ([]byte, error)
type UpdateDocumentStyleRequest
func (s UpdateDocumentStyleRequest) MarshalJSON() ([]byte, error)
type UpdateDocumentTabPropertiesRequest
func (s UpdateDocumentTabPropertiesRequest) MarshalJSON() ([]byte, error)
type UpdateDropdownDefinitionPropertiesRequest
func (s UpdateDropdownDefinitionPropertiesRequest) MarshalJSON() ([]byte, error)
type UpdateDropdownPropertiesRequest
func (s UpdateDropdownPropertiesRequest) MarshalJSON() ([]byte, error)
type UpdateNamedStyleRequest
func (s UpdateNamedStyleRequest) MarshalJSON() ([]byte, error)
type UpdateParagraphStyleRequest
func (s UpdateParagraphStyleRequest) MarshalJSON() ([]byte, error)
type UpdateSectionStyleRequest
func (s UpdateSectionStyleRequest) MarshalJSON() ([]byte, error)
type UpdateTableCellStyleRequest
func (s UpdateTableCellStyleRequest) MarshalJSON() ([]byte, error)
type UpdateTableColumnPropertiesRequest
func (s UpdateTableColumnPropertiesRequest) MarshalJSON() ([]byte, error)
type UpdateTableRowStyleRequest
func (s UpdateTableRowStyleRequest) MarshalJSON() ([]byte, error)
type UpdateTextStyleRequest
func (s UpdateTextStyleRequest) MarshalJSON() ([]byte, error)
type WeightedFontFamily
func (s WeightedFontFamily) MarshalJSON() ([]byte, error)
type WriteControl
func (s WriteControl) MarshalJSON() ([]byte, error)
Constants ¶
View Source
const (
	// See, edit, create, and delete all your Google Docs documents
	DocumentsScope = "https://www.googleapis.com/auth/documents"

	// See all your Google Docs documents
	DocumentsReadonlyScope = "https://www.googleapis.com/auth/documents.readonly"

	// See, edit, create, and delete all of your Google Drive files
	DriveScope = "https://www.googleapis.com/auth/drive"

	// See, edit, create, and delete only the specific Google Drive files you use
	// with this app
	DriveFileScope = "https://www.googleapis.com/auth/drive.file"

	// See and download all your Google Drive files
	DriveReadonlyScope = "https://www.googleapis.com/auth/drive.readonly"
)
OAuth2 scopes used by this API.

Variables ¶
This section is empty.

Functions ¶
This section is empty.

Types ¶
type AcceptSuggestionRequest ¶
type AcceptSuggestionRequest struct {
	// SuggestionId: The ID of the suggestion.
	SuggestionId string `json:"suggestionId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SuggestionId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SuggestionId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
AcceptSuggestionRequest: Accepts a suggestion. Returns a 403 forbidden error if the requesting user does not have edit access to the document.

func (AcceptSuggestionRequest) MarshalJSON ¶
func (s AcceptSuggestionRequest) MarshalJSON() ([]byte, error)
type AddCommentReplyRequest ¶
type AddCommentReplyRequest struct {
	// CommentId: The ID of the CommentThread to add the reply to.
	CommentId string `json:"commentId,omitempty"`
	// Post: The Post representing the reply.
	Post *Post `json:"post,omitempty"`
	// SuggestionId: The ID of the SuggestionThread to add the reply to.
	SuggestionId string `json:"suggestionId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "CommentId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "CommentId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
AddCommentReplyRequest: Inserts a reply Post into a CommentThread or SuggestionThread.

func (AddCommentReplyRequest) MarshalJSON ¶
func (s AddCommentReplyRequest) MarshalJSON() ([]byte, error)
type AddCommentReplyResponse ¶
type AddCommentReplyResponse struct {
	// Post: The newly-inserted reply Post.
	Post *Post `json:"post,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Post") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Post") to include in API requests
	// with the JSON null value. By default, fields with empty values are omitted
	// from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
AddCommentReplyResponse: Response message for adding a reply.

func (AddCommentReplyResponse) MarshalJSON ¶
func (s AddCommentReplyResponse) MarshalJSON() ([]byte, error)
type AddDocumentTabRequest ¶
added in v0.260.0
type AddDocumentTabRequest struct {
	// TabProperties: The properties of the tab to add. All properties are
	// optional.
	TabProperties *TabProperties `json:"tabProperties,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TabProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TabProperties") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
AddDocumentTabRequest: Adds a document tab. When a tab is added at a given index, all subsequent tabs' indexes are incremented.

func (AddDocumentTabRequest) MarshalJSON ¶
added in v0.260.0
func (s AddDocumentTabRequest) MarshalJSON() ([]byte, error)
type AddDocumentTabResponse ¶
added in v0.260.0
type AddDocumentTabResponse struct {
	// TabProperties: The properties of the newly added tab.
	TabProperties *TabProperties `json:"tabProperties,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TabProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TabProperties") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
AddDocumentTabResponse: The result of adding a document tab.

func (AddDocumentTabResponse) MarshalJSON ¶
added in v0.260.0
func (s AddDocumentTabResponse) MarshalJSON() ([]byte, error)
type AutoText ¶
type AutoText struct {
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. An AutoText may have
	// multiple insertion IDs if it's a nested suggested change. If empty, then
	// this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTextStyleChanges: The suggested text style changes to this
	// AutoText, keyed by suggestion ID.
	SuggestedTextStyleChanges map[string]SuggestedTextStyle `json:"suggestedTextStyleChanges,omitempty"`
	// TextStyle: The text style of this AutoText.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// Type: The type of this auto text.
	//
	// Possible values:
	//   "TYPE_UNSPECIFIED" - An unspecified auto text type.
	//   "PAGE_NUMBER" - Type for auto text that represents the current page
	// number.
	//   "PAGE_COUNT" - Type for auto text that represents the total number of
	// pages in the document.
	Type string `json:"type,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SuggestedDeletionIds") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SuggestedDeletionIds") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
AutoText: A ParagraphElement representing a spot in the text that's dynamically replaced with content that can change over time, like a page number.

func (AutoText) MarshalJSON ¶
func (s AutoText) MarshalJSON() ([]byte, error)
type Background ¶
type Background struct {
	// Color: The background color.
	Color *OptionalColor `json:"color,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Color") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Color") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Background: Represents the background of a document.

func (Background) MarshalJSON ¶
func (s Background) MarshalJSON() ([]byte, error)
type BackgroundSuggestionState ¶
type BackgroundSuggestionState struct {
	// BackgroundColorSuggested: Indicates whether the current background color has
	// been modified in this suggestion.
	BackgroundColorSuggested bool `json:"backgroundColorSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "BackgroundColorSuggested")
	// to unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "BackgroundColorSuggested") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
BackgroundSuggestionState: A mask that indicates which of the fields on the base Background have been changed in this suggestion. For any field set to true, the Backgound has a new suggested value.

func (BackgroundSuggestionState) MarshalJSON ¶
func (s BackgroundSuggestionState) MarshalJSON() ([]byte, error)
type BatchUpdateDocumentRequest ¶
type BatchUpdateDocumentRequest struct {
	// Requests: A list of updates to apply to the document.
	Requests []*Request `json:"requests,omitempty"`
	// WriteControl: Provides control over how write requests are executed.
	WriteControl *WriteControl `json:"writeControl,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Requests") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Requests") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
BatchUpdateDocumentRequest: Request message for BatchUpdateDocument.

func (BatchUpdateDocumentRequest) MarshalJSON ¶
func (s BatchUpdateDocumentRequest) MarshalJSON() ([]byte, error)
type BatchUpdateDocumentResponse ¶
type BatchUpdateDocumentResponse struct {
	// CommentUpdateState: Whether comment updates were applied in the batch
	// request.
	//
	// Possible values:
	//   "COMMENT_UPDATE_STATE_UNSPECIFIED" - The status of comment updates is
	// unspecified.
	//   "NO_UPDATES_REQUESTED" - No comment updates were requested in the batch
	// request.
	//   "ALL_SAVED" - All requested comment updates were applied in the batch
	// request.
	//   "ALL_FAILED_UNKNOWN_REASON" - All requested comment updates failed.
	CommentUpdateState string `json:"commentUpdateState,omitempty"`
	// DocumentId: The ID of the document to which the updates were applied to.
	DocumentId string `json:"documentId,omitempty"`
	// Replies: The reply of the updates. This maps 1:1 with the updates, although
	// replies to some requests may be empty.
	Replies []*Response `json:"replies,omitempty"`
	// SuggestionResponses: The suggestions which were affected by each update.
	// This maps 1:1 with the updates.
	SuggestionResponses []*SuggestionResponse `json:"suggestionResponses,omitempty"`
	// WriteControl: The updated write control after applying the request.
	WriteControl *WriteControl `json:"writeControl,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the server.
	googleapi.ServerResponse `json:"-"`
	// ForceSendFields is a list of field names (e.g. "CommentUpdateState") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "CommentUpdateState") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
BatchUpdateDocumentResponse: Response message from a BatchUpdateDocument request.

func (BatchUpdateDocumentResponse) MarshalJSON ¶
func (s BatchUpdateDocumentResponse) MarshalJSON() ([]byte, error)
type Body ¶
type Body struct {
	// Content: The contents of the body. The indexes for the body's content begin
	// at zero.
	Content []*StructuralElement `json:"content,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Content") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Content") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Body: The document body. The body typically contains the full document contents except for headers, footers, and footnotes.

func (Body) MarshalJSON ¶
func (s Body) MarshalJSON() ([]byte, error)
type BookmarkLink ¶
added in v0.190.0
type BookmarkLink struct {
	// Id: The ID of a bookmark in this document.
	Id string `json:"id,omitempty"`
	// TabId: The ID of the tab containing this bookmark.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Id") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Id") to include in API requests
	// with the JSON null value. By default, fields with empty values are omitted
	// from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
BookmarkLink: A reference to a bookmark in this document.

func (BookmarkLink) MarshalJSON ¶
added in v0.190.0
func (s BookmarkLink) MarshalJSON() ([]byte, error)
type Bullet ¶
type Bullet struct {
	// ListId: The ID of the list this paragraph belongs to.
	ListId string `json:"listId,omitempty"`
	// NestingLevel: The nesting level of this paragraph in the list.
	NestingLevel int64 `json:"nestingLevel,omitempty"`
	// TextStyle: The paragraph-specific text style applied to this bullet.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ListId") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ListId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Bullet: Describes the bullet of a paragraph.

func (Bullet) MarshalJSON ¶
func (s Bullet) MarshalJSON() ([]byte, error)
type BulletSuggestionState ¶
type BulletSuggestionState struct {
	// ListIdSuggested: Indicates if there was a suggested change to the list_id.
	ListIdSuggested bool `json:"listIdSuggested,omitempty"`
	// NestingLevelSuggested: Indicates if there was a suggested change to the
	// nesting_level.
	NestingLevelSuggested bool `json:"nestingLevelSuggested,omitempty"`
	// TextStyleSuggestionState: A mask that indicates which of the fields in text
	// style have been changed in this suggestion.
	TextStyleSuggestionState *TextStyleSuggestionState `json:"textStyleSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ListIdSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ListIdSuggested") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
BulletSuggestionState: A mask that indicates which of the fields on the base Bullet have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (BulletSuggestionState) MarshalJSON ¶
func (s BulletSuggestionState) MarshalJSON() ([]byte, error)
type Color ¶
type Color struct {
	// RgbColor: The RGB color value.
	RgbColor *RgbColor `json:"rgbColor,omitempty"`
	// ForceSendFields is a list of field names (e.g. "RgbColor") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "RgbColor") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Color: A solid color.

func (Color) MarshalJSON ¶
func (s Color) MarshalJSON() ([]byte, error)
type ColumnBreak ¶
type ColumnBreak struct {
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. A ColumnBreak may have
	// multiple insertion IDs if it's a nested suggested change. If empty, then
	// this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTextStyleChanges: The suggested text style changes to this
	// ColumnBreak, keyed by suggestion ID.
	SuggestedTextStyleChanges map[string]SuggestedTextStyle `json:"suggestedTextStyleChanges,omitempty"`
	// TextStyle: The text style of this ColumnBreak. Similar to text content, like
	// text runs and footnote references, the text style of a column break can
	// affect content layout as well as the styling of text inserted next to it.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SuggestedDeletionIds") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SuggestedDeletionIds") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ColumnBreak: A ParagraphElement representing a column break. A column break makes the subsequent text start at the top of the next column.

func (ColumnBreak) MarshalJSON ¶
func (s ColumnBreak) MarshalJSON() ([]byte, error)
type CommentAnchor ¶
type CommentAnchor struct {
	// AnchorId: The ID of the comment anchor.
	AnchorId string `json:"anchorId,omitempty"`
	// Ranges: A collection of Ranges in the document which are tied to this
	// anchor.
	Ranges []*Range `json:"ranges,omitempty"`
	// ForceSendFields is a list of field names (e.g. "AnchorId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "AnchorId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CommentAnchor: One or more locations in the document that are tied to CommentThreads with the same anchorId. Note: Multiple anchors may refer to the same location.

func (CommentAnchor) MarshalJSON ¶
func (s CommentAnchor) MarshalJSON() ([]byte, error)
type CommentThread ¶
type CommentThread struct {
	// AnchorId: The ID of the CommentAnchor in the document that this thread is
	// tied to. Multiple comment threads may be anchored to the same CommentAnchor.
	AnchorId string `json:"anchorId,omitempty"`
	// CommentId: The unique ID of the comment thread.
	CommentId string `json:"commentId,omitempty"`
	// HeadPost: The first post in the thread.
	HeadPost *Post `json:"headPost,omitempty"`
	// PlainTextQuote: The quoted text from the document when the comment was
	// created, formatted as plain-text.
	PlainTextQuote string `json:"plainTextQuote,omitempty"`
	// Replies: Replies to the head post.
	Replies []*Post `json:"replies,omitempty"`
	// Status: Whether the thread is open or resolved.
	//
	// Possible values:
	//   "STATUS_UNSPECIFIED" - Default value. This value is unused.
	//   "OPEN" - The comment thread is open.
	//   "RESOLVED" - The comment thread is resolved.
	Status string `json:"status,omitempty"`
	// ForceSendFields is a list of field names (e.g. "AnchorId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "AnchorId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CommentThread: Represents a single comment thread.

func (CommentThread) MarshalJSON ¶
func (s CommentThread) MarshalJSON() ([]byte, error)
type CreateDropdownDefinitionRequest ¶
type CreateDropdownDefinitionRequest struct {
	// DropdownDefinition: Required. The DropdownDefinition to create.
	DropdownDefinition *DropdownDefinition `json:"dropdownDefinition,omitempty"`
	// TabId: The ID of the tab to create the dropdown definition in. When omitted,
	// the request is applied to the first tab. In a document containing a single
	// tab: - If provided, must match the singular tab's ID. - If omitted, the
	// request applies to the singular tab. In a document containing multiple tabs:
	// - If provided, the request applies to the specified tab. - If omitted, the
	// request applies to the first tab in the document.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DropdownDefinition") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DropdownDefinition") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CreateDropdownDefinitionRequest: Creates a new DropdownDefinition in the document.

func (CreateDropdownDefinitionRequest) MarshalJSON ¶
func (s CreateDropdownDefinitionRequest) MarshalJSON() ([]byte, error)
type CreateDropdownDefinitionResponse ¶
type CreateDropdownDefinitionResponse struct {
	// DropdownDefinition: The newly-created DropdownDefinition.
	DropdownDefinition *DropdownDefinition `json:"dropdownDefinition,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DropdownDefinition") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DropdownDefinition") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CreateDropdownDefinitionResponse: Response message for creating a dropdown definition.

func (CreateDropdownDefinitionResponse) MarshalJSON ¶
func (s CreateDropdownDefinitionResponse) MarshalJSON() ([]byte, error)
type CreateFooterRequest ¶
added in v0.12.0
type CreateFooterRequest struct {
	// SectionBreakLocation: The location of the SectionBreak immediately preceding
	// the section whose SectionStyle this footer should belong to. If this is
	// unset or refers to the first section break in the document, the footer
	// applies to the document style.
	SectionBreakLocation *Location `json:"sectionBreakLocation,omitempty"`
	// Type: The type of footer to create.
	//
	// Possible values:
	//   "HEADER_FOOTER_TYPE_UNSPECIFIED" - The header/footer type is unspecified.
	//   "DEFAULT" - A default header/footer.
	Type string `json:"type,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SectionBreakLocation") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SectionBreakLocation") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CreateFooterRequest: Creates a Footer. The new footer is applied to the SectionStyle at the location of the SectionBreak if specified, otherwise it is applied to the DocumentStyle. If a footer of the specified type already exists, a 400 bad request error is returned.

func (CreateFooterRequest) MarshalJSON ¶
added in v0.12.0
func (s CreateFooterRequest) MarshalJSON() ([]byte, error)
type CreateFooterResponse ¶
added in v0.12.0
type CreateFooterResponse struct {
	// FooterId: The ID of the created footer.
	FooterId string `json:"footerId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "FooterId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "FooterId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CreateFooterResponse: The result of creating a footer.

func (CreateFooterResponse) MarshalJSON ¶
added in v0.12.0
func (s CreateFooterResponse) MarshalJSON() ([]byte, error)
type CreateFootnoteRequest ¶
added in v0.16.0
type CreateFootnoteRequest struct {
	// EndOfSegmentLocation: Inserts the footnote reference at the end of the
	// document body. Footnote references cannot be inserted inside a header,
	// footer or footnote. Since footnote references can only be inserted in the
	// body, the segment ID field must be empty.
	EndOfSegmentLocation *EndOfSegmentLocation `json:"endOfSegmentLocation,omitempty"`
	// Location: Inserts the footnote reference at a specific index in the
	// document. The footnote reference must be inserted inside the bounds of an
	// existing Paragraph. For instance, it cannot be inserted at a table's start
	// index (i.e. between the table and its preceding paragraph). Footnote
	// references cannot be inserted inside an equation, header, footer or
	// footnote. Since footnote references can only be inserted in the body, the
	// segment ID field must be empty.
	Location *Location `json:"location,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EndOfSegmentLocation") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EndOfSegmentLocation") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CreateFootnoteRequest: Creates a Footnote segment and inserts a new FootnoteReference to it at the given location. The new Footnote segment will contain a space followed by a newline character.

func (CreateFootnoteRequest) MarshalJSON ¶
added in v0.16.0
func (s CreateFootnoteRequest) MarshalJSON() ([]byte, error)
type CreateFootnoteResponse ¶
added in v0.16.0
type CreateFootnoteResponse struct {
	// FootnoteId: The ID of the created footnote.
	FootnoteId string `json:"footnoteId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "FootnoteId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "FootnoteId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CreateFootnoteResponse: The result of creating a footnote.

func (CreateFootnoteResponse) MarshalJSON ¶
added in v0.16.0
func (s CreateFootnoteResponse) MarshalJSON() ([]byte, error)
type CreateHeaderRequest ¶
added in v0.12.0
type CreateHeaderRequest struct {
	// SectionBreakLocation: The location of the SectionBreak which begins the
	// section this header should belong to. If `section_break_location' is unset
	// or if it refers to the first section break in the document body, the header
	// applies to the DocumentStyle
	SectionBreakLocation *Location `json:"sectionBreakLocation,omitempty"`
	// Type: The type of header to create.
	//
	// Possible values:
	//   "HEADER_FOOTER_TYPE_UNSPECIFIED" - The header/footer type is unspecified.
	//   "DEFAULT" - A default header/footer.
	Type string `json:"type,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SectionBreakLocation") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SectionBreakLocation") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CreateHeaderRequest: Creates a Header. The new header is applied to the SectionStyle at the location of the SectionBreak if specified, otherwise it is applied to the DocumentStyle. If a header of the specified type already exists, a 400 bad request error is returned.

func (CreateHeaderRequest) MarshalJSON ¶
added in v0.12.0
func (s CreateHeaderRequest) MarshalJSON() ([]byte, error)
type CreateHeaderResponse ¶
added in v0.12.0
type CreateHeaderResponse struct {
	// HeaderId: The ID of the created header.
	HeaderId string `json:"headerId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "HeaderId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "HeaderId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CreateHeaderResponse: The result of creating a header.

func (CreateHeaderResponse) MarshalJSON ¶
added in v0.12.0
func (s CreateHeaderResponse) MarshalJSON() ([]byte, error)
type CreateNamedRangeRequest ¶
type CreateNamedRangeRequest struct {
	// Name: The name of the NamedRange. Names do not need to be unique. Names must
	// be at least 1 character and no more than 256 characters, measured in UTF-16
	// code units.
	Name string `json:"name,omitempty"`
	// Range: The range to apply the name to.
	Range *Range `json:"range,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Name") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Name") to include in API requests
	// with the JSON null value. By default, fields with empty values are omitted
	// from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CreateNamedRangeRequest: Creates a NamedRange referencing the given range.

func (CreateNamedRangeRequest) MarshalJSON ¶
func (s CreateNamedRangeRequest) MarshalJSON() ([]byte, error)
type CreateNamedRangeResponse ¶
type CreateNamedRangeResponse struct {
	// NamedRangeId: The ID of the created named range.
	NamedRangeId string `json:"namedRangeId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "NamedRangeId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "NamedRangeId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CreateNamedRangeResponse: The result of creating a named range.

func (CreateNamedRangeResponse) MarshalJSON ¶
func (s CreateNamedRangeResponse) MarshalJSON() ([]byte, error)
type CreateParagraphBulletsRequest ¶
type CreateParagraphBulletsRequest struct {
	// BulletPreset: The kinds of bullet glyphs to be used.
	//
	// Possible values:
	//   "BULLET_GLYPH_PRESET_UNSPECIFIED" - The bullet glyph preset is
	// unspecified.
	//   "BULLET_DISC_CIRCLE_SQUARE" - A bulleted list with a `DISC`, `CIRCLE` and
	// `SQUARE` bullet glyph for the first 3 list nesting levels.
	//   "BULLET_DIAMONDX_ARROW3D_SQUARE" - A bulleted list with a `DIAMONDX`,
	// `ARROW3D` and `SQUARE` bullet glyph for the first 3 list nesting levels.
	//   "BULLET_CHECKBOX" - A bulleted list with `CHECKBOX` bullet glyphs for all
	// list nesting levels.
	//   "BULLET_ARROW_DIAMOND_DISC" - A bulleted list with a `ARROW`, `DIAMOND`
	// and `DISC` bullet glyph for the first 3 list nesting levels.
	//   "BULLET_STAR_CIRCLE_SQUARE" - A bulleted list with a `STAR`, `CIRCLE` and
	// `SQUARE` bullet glyph for the first 3 list nesting levels.
	//   "BULLET_ARROW3D_CIRCLE_SQUARE" - A bulleted list with a `ARROW3D`,
	// `CIRCLE` and `SQUARE` bullet glyph for the first 3 list nesting levels.
	//   "BULLET_LEFTTRIANGLE_DIAMOND_DISC" - A bulleted list with a
	// `LEFTTRIANGLE`, `DIAMOND` and `DISC` bullet glyph for the first 3 list
	// nesting levels.
	//   "BULLET_DIAMONDX_HOLLOWDIAMOND_SQUARE" - A bulleted list with a
	// `DIAMONDX`, `HOLLOWDIAMOND` and `SQUARE` bullet glyph for the first 3 list
	// nesting levels.
	//   "BULLET_DIAMOND_CIRCLE_SQUARE" - A bulleted list with a `DIAMOND`,
	// `CIRCLE` and `SQUARE` bullet glyph for the first 3 list nesting levels.
	//   "NUMBERED_DECIMAL_ALPHA_ROMAN" - A numbered list with `DECIMAL`, `ALPHA`
	// and `ROMAN` numeric glyphs for the first 3 list nesting levels, followed by
	// periods.
	//   "NUMBERED_DECIMAL_ALPHA_ROMAN_PARENS" - A numbered list with `DECIMAL`,
	// `ALPHA` and `ROMAN` numeric glyphs for the first 3 list nesting levels,
	// followed by parenthesis.
	//   "NUMBERED_DECIMAL_NESTED" - A numbered list with `DECIMAL` numeric glyphs
	// separated by periods, where each nesting level uses the previous nesting
	// level's glyph as a prefix. For example: '1.', '1.1.', '2.', '2.2.'.
	//   "NUMBERED_UPPERALPHA_ALPHA_ROMAN" - A numbered list with `UPPERALPHA`,
	// `ALPHA` and `ROMAN` numeric glyphs for the first 3 list nesting levels,
	// followed by periods.
	//   "NUMBERED_UPPERROMAN_UPPERALPHA_DECIMAL" - A numbered list with
	// `UPPERROMAN`, `UPPERALPHA` and `DECIMAL` numeric glyphs for the first 3 list
	// nesting levels, followed by periods.
	//   "NUMBERED_ZERODECIMAL_ALPHA_ROMAN" - A numbered list with `ZERODECIMAL`,
	// `ALPHA` and `ROMAN` numeric glyphs for the first 3 list nesting levels,
	// followed by periods.
	BulletPreset string `json:"bulletPreset,omitempty"`
	// Range: The range to apply the bullet preset to.
	Range *Range `json:"range,omitempty"`
	// ForceSendFields is a list of field names (e.g. "BulletPreset") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "BulletPreset") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CreateParagraphBulletsRequest: Creates bullets for all of the paragraphs that overlap with the given range. The nesting level of each paragraph will be determined by counting leading tabs in front of each paragraph. To avoid excess space between the bullet and the corresponding paragraph, these leading tabs are removed by this request. This may change the indices of parts of the text. If the paragraph immediately before paragraphs being updated is in a list with a matching preset, the paragraphs being updated are added to that preceding list.

func (CreateParagraphBulletsRequest) MarshalJSON ¶
func (s CreateParagraphBulletsRequest) MarshalJSON() ([]byte, error)
type CropProperties ¶
type CropProperties struct {
	// Angle: The clockwise rotation angle of the crop rectangle around its center,
	// in radians. Rotation is applied after the offsets.
	Angle float64 `json:"angle,omitempty"`
	// OffsetBottom: The offset specifies how far inwards the bottom edge of the
	// crop rectangle is from the bottom edge of the original content as a fraction
	// of the original content's height.
	OffsetBottom float64 `json:"offsetBottom,omitempty"`
	// OffsetLeft: The offset specifies how far inwards the left edge of the crop
	// rectangle is from the left edge of the original content as a fraction of the
	// original content's width.
	OffsetLeft float64 `json:"offsetLeft,omitempty"`
	// OffsetRight: The offset specifies how far inwards the right edge of the crop
	// rectangle is from the right edge of the original content as a fraction of
	// the original content's width.
	OffsetRight float64 `json:"offsetRight,omitempty"`
	// OffsetTop: The offset specifies how far inwards the top edge of the crop
	// rectangle is from the top edge of the original content as a fraction of the
	// original content's height.
	OffsetTop float64 `json:"offsetTop,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Angle") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Angle") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CropProperties: The crop properties of an image. The crop rectangle is represented using fractional offsets from the original content's 4 edges. - If the offset is in the interval (0, 1), the corresponding edge of crop rectangle is positioned inside of the image's original bounding rectangle. - If the offset is negative or greater than 1, the corresponding edge of crop rectangle is positioned outside of the image's original bounding rectangle. - If all offsets and rotation angles are 0, the image is not cropped.

func (CropProperties) MarshalJSON ¶
func (s CropProperties) MarshalJSON() ([]byte, error)
func (*CropProperties) UnmarshalJSON ¶
func (s *CropProperties) UnmarshalJSON(data []byte) error
type CropPropertiesSuggestionState ¶
type CropPropertiesSuggestionState struct {
	// AngleSuggested: Indicates if there was a suggested change to angle.
	AngleSuggested bool `json:"angleSuggested,omitempty"`
	// OffsetBottomSuggested: Indicates if there was a suggested change to
	// offset_bottom.
	OffsetBottomSuggested bool `json:"offsetBottomSuggested,omitempty"`
	// OffsetLeftSuggested: Indicates if there was a suggested change to
	// offset_left.
	OffsetLeftSuggested bool `json:"offsetLeftSuggested,omitempty"`
	// OffsetRightSuggested: Indicates if there was a suggested change to
	// offset_right.
	OffsetRightSuggested bool `json:"offsetRightSuggested,omitempty"`
	// OffsetTopSuggested: Indicates if there was a suggested change to offset_top.
	OffsetTopSuggested bool `json:"offsetTopSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "AngleSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "AngleSuggested") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
CropPropertiesSuggestionState: A mask that indicates which of the fields on the base CropProperties have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (CropPropertiesSuggestionState) MarshalJSON ¶
func (s CropPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type DateElement ¶
added in v0.259.0
type DateElement struct {
	// DateElementProperties: The properties of this DateElement.
	DateElementProperties *DateElementProperties `json:"dateElementProperties,omitempty"`
	// DateId: Output only. The unique ID of this date.
	DateId string `json:"dateId,omitempty"`
	// SuggestedDateElementPropertiesChanges: The suggested changes to the date
	// element properties, keyed by suggestion ID.
	SuggestedDateElementPropertiesChanges map[string]SuggestedDateElementProperties `json:"suggestedDateElementPropertiesChanges,omitempty"`
	// SuggestedDeletionIds: IDs for suggestions that remove this date from the
	// document. A DateElement might have multiple deletion IDs if, for example,
	// multiple users suggest deleting it. If empty, then this date isn't suggested
	// for deletion.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: IDs for suggestions that insert this date into the
	// document. A DateElement might have multiple insertion IDs if it's a nested
	// suggested change (a suggestion within a suggestion made by a different user,
	// for example). If empty, then this date isn't a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTextStyleChanges: The suggested text style changes to this
	// DateElement, keyed by suggestion ID.
	SuggestedTextStyleChanges map[string]SuggestedTextStyle `json:"suggestedTextStyleChanges,omitempty"`
	// TextStyle: The text style of this DateElement.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DateElementProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DateElementProperties") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DateElement: A date instance mentioned in a document.

func (DateElement) MarshalJSON ¶
added in v0.259.0
func (s DateElement) MarshalJSON() ([]byte, error)
type DateElementProperties ¶
added in v0.258.0
type DateElementProperties struct {
	// DateFormat: Determines how the date part of the DateElement will be
	// displayed in the document. If unset, the default value is
	// DATE_FORMAT_MONTH_DAY_YEAR_ABBREVIATED, indicating the DateElement will be
	// formatted as `MMM d, y` in `en`, or locale specific equivalent.
	//
	// Possible values:
	//   "DATE_FORMAT_UNSPECIFIED" - The date format is unspecified.
	//   "DATE_FORMAT_CUSTOM" - Output only. The date format is imported from an
	// external source.
	//   "DATE_FORMAT_MONTH_DAY_ABBREVIATED" - The date format is an abbreviated
	// month followed by the day. For example, "Jan 1".
	//   "DATE_FORMAT_MONTH_DAY_FULL" - The date format is a month followed by the
	// day. For example, "January 01".
	//   "DATE_FORMAT_MONTH_DAY_YEAR_ABBREVIATED" - The date format is an
	// abbreviated month followed by the day and the year. For example, "Jan 1,
	// 1970".
	//   "DATE_FORMAT_ISO8601" - The date format is in ISO 8601 format. For
	// example, "1970-01-01".
	DateFormat string `json:"dateFormat,omitempty"`
	// DisplayText: Output only. Indicates how the DateElement is displayed in the
	// document.
	DisplayText string `json:"displayText,omitempty"`
	// Locale: The language code of the DateElement. For example, `en`. If unset,
	// the default locale is `en`. Limited to the following locales: `af`, `am`,
	// `ar`, `as`, `az`, `be`, `bg`, `bn`, `ca`, `cs`, `da`, `de`, `el`, `en`,
	// `en-CA`, `en-GB`, `es`, `es-419`, `et`, `eu`, `fa`, `fi`, `fil`, `fr`,
	// `fr-CA`, `gl`, `gu`, `hi`, `hr`, `hu`, `hy`, `id`, `is`, `it`, `iw`, `ja`,
	// `ka`, `kk`, `km`, `kn`, `ko`, `lo`, `lt`, `lv`, `mk`, `ml`, `mn`, `mr`,
	// `ms`, `ne`, `nl`, `no`, `or`, `pa`, `pl`, `pt-BR`, `pt-PT`, `ro`, `ru`,
	// `si`, `sk`, `sl`, `sq`, `sr`, `sv`, `sw`, `ta`, `te`, `th`, `tr`, `uk`,
	// `ur`, `uz`, `vi`, `zh-CN`, `zh-HK`, `zh-TW`, `zu`, `cy`, `my`.
	Locale string `json:"locale,omitempty"`
	// TimeFormat: Determines how the time part of the DateElement will be
	// displayed in the document. If unset, the default value is
	// TIME_FORMAT_DISABLED, indicating no time should be shown.
	//
	// Possible values:
	//   "TIME_FORMAT_UNSPECIFIED" - The time format is unspecified.
	//   "TIME_FORMAT_DISABLED" - Indicates that the date does not have a time.
	//   "TIME_FORMAT_HOUR_MINUTE" - The time format shows the hour and minute. For
	// example, "Jan 1, 1970 12:00 PM".
	//   "TIME_FORMAT_HOUR_MINUTE_TIMEZONE" - The time format shows the hour,
	// minute, and timezone. For example, "Jan 1, 1970 12:00 PM UTC".
	TimeFormat string `json:"timeFormat,omitempty"`
	// TimeZoneId: The time zone of the DateElement, as defined by the Unicode
	// Common Locale Data Repository (CLDR) project. For example,
	// `America/New_York`. If unset, the default time zone is `etc/UTC`.
	TimeZoneId string `json:"timeZoneId,omitempty"`
	// Timestamp: The point in time to represent, in seconds and nanoseconds since
	// Unix epoch: January 1, 1970 at midnight UTC. Timestamp is expected to be in
	// UTC. If time_zone_id is set, the timestamp is adjusted according to the time
	// zone. For example, a timestamp of `18000` with a date format of
	// `DATE_FORMAT_ISO8601` and time format of `TIME_FORMAT_HOUR_MINUTE` would be
	// displayed as `1970-01-01 5:00 AM`. A timestamp of `18000` with date format
	// of `DATE_FORMAT_ISO8601`, time format of `TIME_FORMAT_HOUR_MINUTE`, and time
	// zone set to `America/New_York` will instead be `1970-01-01 12:00 AM`.
	Timestamp string `json:"timestamp,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DateFormat") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DateFormat") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DateElementProperties: Properties of a DateElement.

func (DateElementProperties) MarshalJSON ¶
added in v0.258.0
func (s DateElementProperties) MarshalJSON() ([]byte, error)
type DateElementPropertiesSuggestionState ¶
added in v0.259.0
type DateElementPropertiesSuggestionState struct {
	// DateFormatSuggested: Indicates if there was a suggested change to
	// date_format.
	DateFormatSuggested bool `json:"dateFormatSuggested,omitempty"`
	// LocaleSuggested: Indicates if there was a suggested change to locale.
	LocaleSuggested bool `json:"localeSuggested,omitempty"`
	// TimeFormatSuggested: Indicates if there was a suggested change to
	// time_format.
	TimeFormatSuggested bool `json:"timeFormatSuggested,omitempty"`
	// TimeZoneIdSuggested: Indicates if there was a suggested change to
	// time_zone_id.
	TimeZoneIdSuggested bool `json:"timeZoneIdSuggested,omitempty"`
	// TimestampSuggested: Indicates if there was a suggested change to timestamp.
	TimestampSuggested bool `json:"timestampSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DateFormatSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DateFormatSuggested") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DateElementPropertiesSuggestionState: A mask that indicates which of the fields on the base DateElementProperties have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (DateElementPropertiesSuggestionState) MarshalJSON ¶
added in v0.259.0
func (s DateElementPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type DeleteCommentReplyRequest ¶
type DeleteCommentReplyRequest struct {
	// CommentId: The ID of the CommentThread which the post belongs to.
	CommentId string `json:"commentId,omitempty"`
	// PostId: The ID of the reply Post being deleted.
	PostId string `json:"postId,omitempty"`
	// SuggestionId: The ID of the SuggestionThread which the post belongs to.
	SuggestionId string `json:"suggestionId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "CommentId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "CommentId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteCommentReplyRequest: Deletes a reply Post from a CommentThread or SuggestionThread. Returns a 400 bad request error if: - The requesting user is not the author of the post. - The reply post contains an action. - The reply post contains an assignee.

func (DeleteCommentReplyRequest) MarshalJSON ¶
func (s DeleteCommentReplyRequest) MarshalJSON() ([]byte, error)
type DeleteCommentRequest ¶
type DeleteCommentRequest struct {
	// CommentId: The ID of the CommentThread that is being deleted.
	CommentId string `json:"commentId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "CommentId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "CommentId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteCommentRequest: Deletes a CommentThread. Returns a 400 bad request error if the requesting user is not the author of the headPost.

func (DeleteCommentRequest) MarshalJSON ¶
func (s DeleteCommentRequest) MarshalJSON() ([]byte, error)
type DeleteContentRangeRequest ¶
type DeleteContentRangeRequest struct {
	// Range: The range of content to delete. Deleting text that crosses a
	// paragraph boundary may result in changes to paragraph styles, lists,
	// positioned objects and bookmarks as the two paragraphs are merged.
	// Attempting to delete certain ranges can result in an invalid document
	// structure in which case a 400 bad request error is returned. Some examples
	// of invalid delete requests include: * Deleting one code unit of a surrogate
	// pair. * Deleting the last newline character of a Body, Header, Footer,
	// Footnote, TableCell or TableOfContents. * Deleting the start or end of a
	// Table, TableOfContents or Equation without deleting the entire element. *
	// Deleting the newline character before a Table, TableOfContents or
	// SectionBreak without deleting the element. * Deleting individual rows or
	// cells of a table. Deleting the content within a table cell is allowed.
	Range *Range `json:"range,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Range") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Range") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteContentRangeRequest: Deletes content from the document.

func (DeleteContentRangeRequest) MarshalJSON ¶
func (s DeleteContentRangeRequest) MarshalJSON() ([]byte, error)
type DeleteDropdownDefinitionRequest ¶
type DeleteDropdownDefinitionRequest struct {
	// DropdownDefinitionId: The ID of the DropdownDefinition to delete.
	DropdownDefinitionId string `json:"dropdownDefinitionId,omitempty"`
	// TabId: The ID of the tab that contains the dropdown definition to delete.
	// When omitted, the request is applied to the first tab. In a document
	// containing a single tab: - If provided, must match the singular tab's ID. -
	// If omitted, the request applies to the singular tab. In a document
	// containing multiple tabs: - If provided, the request applies to the
	// specified tab. - If omitted, the request applies to the first tab in the
	// document.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DropdownDefinitionId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DropdownDefinitionId") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteDropdownDefinitionRequest: Deletes a DropdownDefinition. If the dropdown definition is referenced by any dropdown instances (chips) in the document, a 400 bad request error is returned.

func (DeleteDropdownDefinitionRequest) MarshalJSON ¶
func (s DeleteDropdownDefinitionRequest) MarshalJSON() ([]byte, error)
type DeleteFooterRequest ¶
added in v0.16.0
type DeleteFooterRequest struct {
	// FooterId: The id of the footer to delete. If this footer is defined on
	// DocumentStyle, the reference to this footer is removed, resulting in no
	// footer of that type for the first section of the document. If this footer is
	// defined on a SectionStyle, the reference to this footer is removed and the
	// footer of that type is now continued from the previous section.
	FooterId string `json:"footerId,omitempty"`
	// TabId: The tab that contains the footer to delete. When omitted, the request
	// is applied to the first tab. In a document containing a single tab: - If
	// provided, must match the singular tab's ID. - If omitted, the request
	// applies to the singular tab. In a document containing multiple tabs: - If
	// provided, the request applies to the specified tab. - If omitted, the
	// request applies to the first tab in the document.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "FooterId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "FooterId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteFooterRequest: Deletes a Footer from the document.

func (DeleteFooterRequest) MarshalJSON ¶
added in v0.16.0
func (s DeleteFooterRequest) MarshalJSON() ([]byte, error)
type DeleteHeaderRequest ¶
added in v0.16.0
type DeleteHeaderRequest struct {
	// HeaderId: The id of the header to delete. If this header is defined on
	// DocumentStyle, the reference to this header is removed, resulting in no
	// header of that type for the first section of the document. If this header is
	// defined on a SectionStyle, the reference to this header is removed and the
	// header of that type is now continued from the previous section.
	HeaderId string `json:"headerId,omitempty"`
	// TabId: The tab containing the header to delete. When omitted, the request is
	// applied to the first tab. In a document containing a single tab: - If
	// provided, must match the singular tab's ID. - If omitted, the request
	// applies to the singular tab. In a document containing multiple tabs: - If
	// provided, the request applies to the specified tab. - If omitted, the
	// request applies to the first tab in the document.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "HeaderId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "HeaderId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteHeaderRequest: Deletes a Header from the document.

func (DeleteHeaderRequest) MarshalJSON ¶
added in v0.16.0
func (s DeleteHeaderRequest) MarshalJSON() ([]byte, error)
type DeleteNamedRangeRequest ¶
type DeleteNamedRangeRequest struct {
	// Name: The name of the range(s) to delete. All named ranges with the given
	// name will be deleted.
	Name string `json:"name,omitempty"`
	// NamedRangeId: The ID of the named range to delete.
	NamedRangeId string `json:"namedRangeId,omitempty"`
	// TabsCriteria: Optional. The criteria used to specify which tab(s) the range
	// deletion should occur in. When omitted, the range deletion is applied to all
	// tabs. In a document containing a single tab: - If provided, must match the
	// singular tab's ID. - If omitted, the range deletion applies to the singular
	// tab. In a document containing multiple tabs: - If provided, the range
	// deletion applies to the specified tabs. - If not provided, the range
	// deletion applies to all tabs.
	TabsCriteria *TabsCriteria `json:"tabsCriteria,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Name") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Name") to include in API requests
	// with the JSON null value. By default, fields with empty values are omitted
	// from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteNamedRangeRequest: Deletes a NamedRange.

func (DeleteNamedRangeRequest) MarshalJSON ¶
func (s DeleteNamedRangeRequest) MarshalJSON() ([]byte, error)
type DeleteParagraphBulletsRequest ¶
type DeleteParagraphBulletsRequest struct {
	// Range: The range to delete bullets from.
	Range *Range `json:"range,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Range") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Range") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteParagraphBulletsRequest: Deletes bullets from all of the paragraphs that overlap with the given range. The nesting level of each paragraph will be visually preserved by adding indent to the start of the corresponding paragraph.

func (DeleteParagraphBulletsRequest) MarshalJSON ¶
func (s DeleteParagraphBulletsRequest) MarshalJSON() ([]byte, error)
type DeletePositionedObjectRequest ¶
type DeletePositionedObjectRequest struct {
	// ObjectId: The ID of the positioned object to delete.
	ObjectId string `json:"objectId,omitempty"`
	// TabId: The tab that the positioned object to delete is in. When omitted, the
	// request is applied to the first tab. In a document containing a single tab:
	// - If provided, must match the singular tab's ID. - If omitted, the request
	// applies to the singular tab. In a document containing multiple tabs: - If
	// provided, the request applies to the specified tab. - If omitted, the
	// request applies to the first tab in the document.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ObjectId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ObjectId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeletePositionedObjectRequest: Deletes a PositionedObject from the document.

func (DeletePositionedObjectRequest) MarshalJSON ¶
func (s DeletePositionedObjectRequest) MarshalJSON() ([]byte, error)
type DeleteSuggestionRequest ¶
type DeleteSuggestionRequest struct {
	// SuggestionId: The ID of the suggestion.
	SuggestionId string `json:"suggestionId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SuggestionId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SuggestionId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteSuggestionRequest: Deletes a suggestion. Returns a 403 forbidden error if the requesting user is not the author of the suggestion.

func (DeleteSuggestionRequest) MarshalJSON ¶
func (s DeleteSuggestionRequest) MarshalJSON() ([]byte, error)
type DeleteTabRequest ¶
added in v0.260.0
type DeleteTabRequest struct {
	// TabId: The ID of the tab to delete.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TabId") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TabId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteTabRequest: Deletes a tab. If the tab has child tabs, they are deleted as well.

func (DeleteTabRequest) MarshalJSON ¶
added in v0.260.0
func (s DeleteTabRequest) MarshalJSON() ([]byte, error)
type DeleteTableColumnRequest ¶
type DeleteTableColumnRequest struct {
	// TableCellLocation: The reference table cell location from which the column
	// will be deleted. The column this cell spans will be deleted. If this is a
	// merged cell that spans multiple columns, all columns that the cell spans
	// will be deleted. If no columns remain in the table after this deletion, the
	// whole table is deleted.
	TableCellLocation *TableCellLocation `json:"tableCellLocation,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TableCellLocation") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TableCellLocation") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteTableColumnRequest: Deletes a column from a table.

func (DeleteTableColumnRequest) MarshalJSON ¶
func (s DeleteTableColumnRequest) MarshalJSON() ([]byte, error)
type DeleteTableRowRequest ¶
type DeleteTableRowRequest struct {
	// TableCellLocation: The reference table cell location from which the row will
	// be deleted. The row this cell spans will be deleted. If this is a merged
	// cell that spans multiple rows, all rows that the cell spans will be deleted.
	// If no rows remain in the table after this deletion, the whole table is
	// deleted.
	TableCellLocation *TableCellLocation `json:"tableCellLocation,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TableCellLocation") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TableCellLocation") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DeleteTableRowRequest: Deletes a row from a table.

func (DeleteTableRowRequest) MarshalJSON ¶
func (s DeleteTableRowRequest) MarshalJSON() ([]byte, error)
type Dimension ¶
type Dimension struct {
	// Magnitude: The magnitude.
	Magnitude float64 `json:"magnitude,omitempty"`
	// Unit: The units for magnitude.
	//
	// Possible values:
	//   "UNIT_UNSPECIFIED" - The units are unknown.
	//   "PT" - A point, 1/72 of an inch.
	Unit string `json:"unit,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Magnitude") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Magnitude") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Dimension: A magnitude in a single direction in the specified units.

func (Dimension) MarshalJSON ¶
func (s Dimension) MarshalJSON() ([]byte, error)
func (*Dimension) UnmarshalJSON ¶
func (s *Dimension) UnmarshalJSON(data []byte) error
type Document ¶
type Document struct {
	// Body: Output only. The main body of the document. Legacy field: Instead, use
	// Document.tabs.documentTab.body, which exposes the actual document content
	// from all tabs when the includeTabsContent parameter is set to `true`. If
	// `false` or unset, this field contains information about the first tab in the
	// document.
	Body *Body `json:"body,omitempty"`
	// Comments: Output only. The comments associated with the document. Only
	// populated if the commentsViewMode parameter is set to require comments (such
	// as `COMMENTS_VIEW_MODE_INCLUDED`).
	Comments []*CommentThread `json:"comments,omitempty"`
	// CommentsViewMode: Output only. The comments view mode applied to the
	// document.
	//
	// Possible values:
	//   "COMMENTS_VIEW_MODE_UNSPECIFIED" - The CommentsViewMode is unspecified.
	// COMMENTS_VIEW_MODE_OMITTED is applied.
	//   "COMMENTS_VIEW_MODE_DEFAULT_FOR_CURRENT_ACCESS" - The CommentsViewMode
	// applied to the returned document depends on the user's current access level.
	// If the user only has view access, COMMENTS_VIEW_MODE_OMITTED is applied.
	// Otherwise, COMMENTS_VIEW_MODE_INCLUDED is applied.
	//   "COMMENTS_VIEW_MODE_OMITTED" - The returned document has comments omitted.
	//   "COMMENTS_VIEW_MODE_INCLUDED" - The returned document has comments
	// included. Requests to retrieve a document using this mode will return a 403
	// error if the user does not have permission to view comments. When set,
	// suggestions_view_mode must also be set to SUGGESTIONS_INLINE. Returns a 400
	// bad request error otherwise.
	CommentsViewMode string `json:"commentsViewMode,omitempty"`
	// DocumentId: Output only. The ID of the document.
	DocumentId string `json:"documentId,omitempty"`
	// DocumentStyle: Output only. The style of the document. Legacy field:
	// Instead, use Document.tabs.documentTab.documentStyle, which exposes the
	// actual document content from all tabs when the includeTabsContent parameter
	// is set to `true`. If `false` or unset, this field contains information about
	// the first tab in the document.
	DocumentStyle *DocumentStyle `json:"documentStyle,omitempty"`
	// Footers: Output only. The footers in the document, keyed by footer ID.
	// Legacy field: Instead, use Document.tabs.documentTab.footers, which exposes
	// the actual document content from all tabs when the includeTabsContent
	// parameter is set to `true`. If `false` or unset, this field contains
	// information about the first tab in the document.
	Footers map[string]Footer `json:"footers,omitempty"`
	// Footnotes: Output only. The footnotes in the document, keyed by footnote ID.
	// Legacy field: Instead, use Document.tabs.documentTab.footnotes, which
	// exposes the actual document content from all tabs when the
	// includeTabsContent parameter is set to `true`. If `false` or unset, this
	// field contains information about the first tab in the document.
	Footnotes map[string]Footnote `json:"footnotes,omitempty"`
	// Headers: Output only. The headers in the document, keyed by header ID.
	// Legacy field: Instead, use Document.tabs.documentTab.headers, which exposes
	// the actual document content from all tabs when the includeTabsContent
	// parameter is set to `true`. If `false` or unset, this field contains
	// information about the first tab in the document.
	Headers map[string]Header `json:"headers,omitempty"`
	// InlineObjects: Output only. The inline objects in the document, keyed by
	// object ID. Legacy field: Instead, use
	// Document.tabs.documentTab.inlineObjects, which exposes the actual document
	// content from all tabs when the includeTabsContent parameter is set to
	// `true`. If `false` or unset, this field contains information about the first
	// tab in the document.
	InlineObjects map[string]InlineObject `json:"inlineObjects,omitempty"`
	// Lists: Output only. The lists in the document, keyed by list ID. Legacy
	// field: Instead, use Document.tabs.documentTab.lists, which exposes the
	// actual document content from all tabs when the includeTabsContent parameter
	// is set to `true`. If `false` or unset, this field contains information about
	// the first tab in the document.
	Lists map[string]List `json:"lists,omitempty"`
	// NamedRanges: Output only. The named ranges in the document, keyed by name.
	// Legacy field: Instead, use Document.tabs.documentTab.namedRanges, which
	// exposes the actual document content from all tabs when the
	// includeTabsContent parameter is set to `true`. If `false` or unset, this
	// field contains information about the first tab in the document.
	NamedRanges map[string]NamedRanges `json:"namedRanges,omitempty"`
	// NamedStyles: Output only. The named styles of the document. Legacy field:
	// Instead, use Document.tabs.documentTab.namedStyles, which exposes the actual
	// document content from all tabs when the includeTabsContent parameter is set
	// to `true`. If `false` or unset, this field contains information about the
	// first tab in the document.
	NamedStyles *NamedStyles `json:"namedStyles,omitempty"`
	// PositionedObjects: Output only. The positioned objects in the document,
	// keyed by object ID. Legacy field: Instead, use
	// Document.tabs.documentTab.positionedObjects, which exposes the actual
	// document content from all tabs when the includeTabsContent parameter is set
	// to `true`. If `false` or unset, this field contains information about the
	// first tab in the document.
	PositionedObjects map[string]PositionedObject `json:"positionedObjects,omitempty"`
	// RevisionId: Output only. The revision ID of the document. Can be used in
	// update requests to specify which revision of a document to apply updates to
	// and how the request should behave if the document has been edited since that
	// revision. Only populated if the user has edit access to the document. The
	// revision ID is not a sequential number but an opaque string. The format of
	// the revision ID might change over time. A returned revision ID is only
	// guaranteed to be valid for 24 hours after it has been returned and cannot be
	// shared across users. If the revision ID is unchanged between calls, then the
	// document has not changed. Conversely, a changed ID (for the same document
	// and user) usually means the document has been updated. However, a changed ID
	// can also be due to internal factors such as ID format changes.
	RevisionId string `json:"revisionId,omitempty"`
	// SuggestedDocumentStyleChanges: Output only. The suggested changes to the
	// style of the document, keyed by suggestion ID. Legacy field: Instead, use
	// Document.tabs.documentTab.suggestedDocumentStyleChanges, which exposes the
	// actual document content from all tabs when the includeTabsContent parameter
	// is set to `true`. If `false` or unset, this field contains information about
	// the first tab in the document.
	SuggestedDocumentStyleChanges map[string]SuggestedDocumentStyle `json:"suggestedDocumentStyleChanges,omitempty"`
	// SuggestedNamedStylesChanges: Output only. The suggested changes to the named
	// styles of the document, keyed by suggestion ID. Legacy field: Instead, use
	// Document.tabs.documentTab.suggestedNamedStylesChanges, which exposes the
	// actual document content from all tabs when the includeTabsContent parameter
	// is set to `true`. If `false` or unset, this field contains information about
	// the first tab in the document.
	SuggestedNamedStylesChanges map[string]SuggestedNamedStyles `json:"suggestedNamedStylesChanges,omitempty"`
	// Suggestions: Output only. The suggestions associated with the document. Only
	// populated if the commentsViewMode parameter is set to require comments (such
	// as `COMMENTS_VIEW_MODE_INCLUDED`).
	Suggestions []*SuggestionThread `json:"suggestions,omitempty"`
	// SuggestionsViewMode: Output only. The suggestions view mode applied to the
	// document. Note: When editing a document, changes must be based on a document
	// with SUGGESTIONS_INLINE.
	//
	// Possible values:
	//   "DEFAULT_FOR_CURRENT_ACCESS" - The SuggestionsViewMode applied to the
	// returned document depends on the user's current access level. If the user
	// only has view access, PREVIEW_WITHOUT_SUGGESTIONS is applied. Otherwise,
	// SUGGESTIONS_INLINE is applied. This is the default suggestions view mode.
	//   "SUGGESTIONS_INLINE" - The returned document has suggestions inline.
	// Suggested changes will be differentiated from base content within the
	// document. Requests to retrieve a document using this mode will return a 403
	// error if the user does not have permission to view suggested changes.
	//   "PREVIEW_SUGGESTIONS_ACCEPTED" - The returned document is a preview with
	// all suggested changes accepted. Requests to retrieve a document using this
	// mode will return a 403 error if the user does not have permission to view
	// suggested changes.
	//   "PREVIEW_WITHOUT_SUGGESTIONS" - The returned document is a preview with
	// all suggested changes rejected if there are any suggestions in the document.
	SuggestionsViewMode string `json:"suggestionsViewMode,omitempty"`
	// Tabs: Output only. Tabs that are part of a document. Tabs can contain child
	// tabs, a tab nested within another tab. Child tabs are represented by the
	// Tab.childTabs field.
	Tabs []*Tab `json:"tabs,omitempty"`
	// Title: The title of the document.
	Title string `json:"title,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the server.
	googleapi.ServerResponse `json:"-"`
	// ForceSendFields is a list of field names (e.g. "Body") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Body") to include in API requests
	// with the JSON null value. By default, fields with empty values are omitted
	// from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Document: A Google Docs document.

func (Document) MarshalJSON ¶
func (s Document) MarshalJSON() ([]byte, error)
type DocumentFormat ¶
added in v0.253.0
type DocumentFormat struct {
	// DocumentMode: Whether the document has pages or is pageless.
	//
	// Possible values:
	//   "DOCUMENT_MODE_UNSPECIFIED" - The document mode is unspecified.
	//   "PAGES" - The document has pages.
	//   "PAGELESS" - The document is pageless.
	DocumentMode string `json:"documentMode,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DocumentMode") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DocumentMode") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DocumentFormat: Represents document-level format settings.

func (DocumentFormat) MarshalJSON ¶
added in v0.253.0
func (s DocumentFormat) MarshalJSON() ([]byte, error)
type DocumentStyle ¶
type DocumentStyle struct {
	// Background: The background of the document. Documents cannot have a
	// transparent background color.
	Background *Background `json:"background,omitempty"`
	// DefaultFooterId: The ID of the default footer. If not set, there's no
	// default footer. If DocumentMode is PAGELESS, this property will not be
	// rendered. This property is read-only.
	DefaultFooterId string `json:"defaultFooterId,omitempty"`
	// DefaultHeaderId: The ID of the default header. If not set, there's no
	// default header. If DocumentMode is PAGELESS, this property will not be
	// rendered. This property is read-only.
	DefaultHeaderId string `json:"defaultHeaderId,omitempty"`
	// DocumentFormat: Specifies document-level format settings, such as the
	// document mode (pages vs pageless).
	DocumentFormat *DocumentFormat `json:"documentFormat,omitempty"`
	// EvenPageFooterId: The ID of the footer used only for even pages. The value
	// of use_even_page_header_footer determines whether to use the
	// default_footer_id or this value for the footer on even pages. If not set,
	// there's no even page footer. If DocumentMode is PAGELESS, this property will
	// not be rendered. This property is read-only.
	EvenPageFooterId string `json:"evenPageFooterId,omitempty"`
	// EvenPageHeaderId: The ID of the header used only for even pages. The value
	// of use_even_page_header_footer determines whether to use the
	// default_header_id or this value for the header on even pages. If not set,
	// there's no even page header. If DocumentMode is PAGELESS, this property will
	// not be rendered. This property is read-only.
	EvenPageHeaderId string `json:"evenPageHeaderId,omitempty"`
	// FirstPageFooterId: The ID of the footer used only for the first page. If not
	// set then a unique footer for the first page does not exist. The value of
	// use_first_page_header_footer determines whether to use the default_footer_id
	// or this value for the footer on the first page. If not set, there's no first
	// page footer. If DocumentMode is PAGELESS, this property will not be
	// rendered. This property is read-only.
	FirstPageFooterId string `json:"firstPageFooterId,omitempty"`
	// FirstPageHeaderId: The ID of the header used only for the first page. If not
	// set then a unique header for the first page does not exist. The value of
	// use_first_page_header_footer determines whether to use the default_header_id
	// or this value for the header on the first page. If not set, there's no first
	// page header. If DocumentMode is PAGELESS, this property will not be
	// rendered. This property is read-only.
	FirstPageHeaderId string `json:"firstPageHeaderId,omitempty"`
	// FlipPageOrientation: Optional. Indicates whether to flip the dimensions of
	// the page_size, which allows changing the page orientation between portrait
	// and landscape. If DocumentMode is PAGELESS, this property will not be
	// rendered.
	FlipPageOrientation bool `json:"flipPageOrientation,omitempty"`
	// MarginBottom: The bottom page margin. Updating the bottom page margin on the
	// document style clears the bottom page margin on all section styles. If
	// DocumentMode is PAGELESS, this property will not be rendered.
	MarginBottom *Dimension `json:"marginBottom,omitempty"`
	// MarginFooter: The amount of space between the bottom of the page and the
	// contents of the footer. If DocumentMode is PAGELESS, this property will not
	// be rendered.
	MarginFooter *Dimension `json:"marginFooter,omitempty"`
	// MarginHeader: The amount of space between the top of the page and the
	// contents of the header. If DocumentMode is PAGELESS, this property will not
	// be rendered.
	MarginHeader *Dimension `json:"marginHeader,omitempty"`
	// MarginLeft: The left page margin. Updating the left page margin on the
	// document style clears the left page margin on all section styles. It may
	// also cause columns to resize in all sections. If DocumentMode is PAGELESS,
	// this property will not be rendered.
	MarginLeft *Dimension `json:"marginLeft,omitempty"`
	// MarginRight: The right page margin. Updating the right page margin on the
	// document style clears the right page margin on all section styles. It may
	// also cause columns to resize in all sections. If DocumentMode is PAGELESS,
	// this property will not be rendered.
	MarginRight *Dimension `json:"marginRight,omitempty"`
	// MarginTop: The top page margin. Updating the top page margin on the document
	// style clears the top page margin on all section styles. If DocumentMode is
	// PAGELESS, this property will not be rendered.
	MarginTop *Dimension `json:"marginTop,omitempty"`
	// PageNumberStart: The page number from which to start counting the number of
	// pages. If DocumentMode is PAGELESS, this property will not be rendered.
	PageNumberStart int64 `json:"pageNumberStart,omitempty"`
	// PageSize: The size of a page in the document. If DocumentMode is PAGELESS,
	// this property will not be rendered.
	PageSize *Size `json:"pageSize,omitempty"`
	// UseCustomHeaderFooterMargins: Indicates whether DocumentStyle margin_header,
	// SectionStyle margin_header and DocumentStyle margin_footer, SectionStyle
	// margin_footer are respected. When false, the default values in the Docs
	// editor for header and footer margin is used. If DocumentMode is PAGELESS,
	// this property will not be rendered. This property is read-only.
	UseCustomHeaderFooterMargins bool `json:"useCustomHeaderFooterMargins,omitempty"`
	// UseEvenPageHeaderFooter: Indicates whether to use the even page header /
	// footer IDs for the even pages. If DocumentMode is PAGELESS, this property
	// will not be rendered.
	UseEvenPageHeaderFooter bool `json:"useEvenPageHeaderFooter,omitempty"`
	// UseFirstPageHeaderFooter: Indicates whether to use the first page header /
	// footer IDs for the first page. If DocumentMode is PAGELESS, this property
	// will not be rendered.
	UseFirstPageHeaderFooter bool `json:"useFirstPageHeaderFooter,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Background") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Background") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DocumentStyle: The style of the document.

func (DocumentStyle) MarshalJSON ¶
func (s DocumentStyle) MarshalJSON() ([]byte, error)
type DocumentStyleSuggestionState ¶
type DocumentStyleSuggestionState struct {
	// BackgroundSuggestionState: A mask that indicates which of the fields in
	// background have been changed in this suggestion.
	BackgroundSuggestionState *BackgroundSuggestionState `json:"backgroundSuggestionState,omitempty"`
	// DefaultFooterIdSuggested: Indicates if there was a suggested change to
	// default_footer_id.
	DefaultFooterIdSuggested bool `json:"defaultFooterIdSuggested,omitempty"`
	// DefaultHeaderIdSuggested: Indicates if there was a suggested change to
	// default_header_id.
	DefaultHeaderIdSuggested bool `json:"defaultHeaderIdSuggested,omitempty"`
	// EvenPageFooterIdSuggested: Indicates if there was a suggested change to
	// even_page_footer_id.
	EvenPageFooterIdSuggested bool `json:"evenPageFooterIdSuggested,omitempty"`
	// EvenPageHeaderIdSuggested: Indicates if there was a suggested change to
	// even_page_header_id.
	EvenPageHeaderIdSuggested bool `json:"evenPageHeaderIdSuggested,omitempty"`
	// FirstPageFooterIdSuggested: Indicates if there was a suggested change to
	// first_page_footer_id.
	FirstPageFooterIdSuggested bool `json:"firstPageFooterIdSuggested,omitempty"`
	// FirstPageHeaderIdSuggested: Indicates if there was a suggested change to
	// first_page_header_id.
	FirstPageHeaderIdSuggested bool `json:"firstPageHeaderIdSuggested,omitempty"`
	// FlipPageOrientationSuggested: Optional. Indicates if there was a suggested
	// change to flip_page_orientation.
	FlipPageOrientationSuggested bool `json:"flipPageOrientationSuggested,omitempty"`
	// MarginBottomSuggested: Indicates if there was a suggested change to
	// margin_bottom.
	MarginBottomSuggested bool `json:"marginBottomSuggested,omitempty"`
	// MarginFooterSuggested: Indicates if there was a suggested change to
	// margin_footer.
	MarginFooterSuggested bool `json:"marginFooterSuggested,omitempty"`
	// MarginHeaderSuggested: Indicates if there was a suggested change to
	// margin_header.
	MarginHeaderSuggested bool `json:"marginHeaderSuggested,omitempty"`
	// MarginLeftSuggested: Indicates if there was a suggested change to
	// margin_left.
	MarginLeftSuggested bool `json:"marginLeftSuggested,omitempty"`
	// MarginRightSuggested: Indicates if there was a suggested change to
	// margin_right.
	MarginRightSuggested bool `json:"marginRightSuggested,omitempty"`
	// MarginTopSuggested: Indicates if there was a suggested change to margin_top.
	MarginTopSuggested bool `json:"marginTopSuggested,omitempty"`
	// PageNumberStartSuggested: Indicates if there was a suggested change to
	// page_number_start.
	PageNumberStartSuggested bool `json:"pageNumberStartSuggested,omitempty"`
	// PageSizeSuggestionState: A mask that indicates which of the fields in size
	// have been changed in this suggestion.
	PageSizeSuggestionState *SizeSuggestionState `json:"pageSizeSuggestionState,omitempty"`
	// UseCustomHeaderFooterMarginsSuggested: Indicates if there was a suggested
	// change to use_custom_header_footer_margins.
	UseCustomHeaderFooterMarginsSuggested bool `json:"useCustomHeaderFooterMarginsSuggested,omitempty"`
	// UseEvenPageHeaderFooterSuggested: Indicates if there was a suggested change
	// to use_even_page_header_footer.
	UseEvenPageHeaderFooterSuggested bool `json:"useEvenPageHeaderFooterSuggested,omitempty"`
	// UseFirstPageHeaderFooterSuggested: Indicates if there was a suggested change
	// to use_first_page_header_footer.
	UseFirstPageHeaderFooterSuggested bool `json:"useFirstPageHeaderFooterSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "BackgroundSuggestionState")
	// to unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "BackgroundSuggestionState") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DocumentStyleSuggestionState: A mask that indicates which of the fields on the base DocumentStyle have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (DocumentStyleSuggestionState) MarshalJSON ¶
func (s DocumentStyleSuggestionState) MarshalJSON() ([]byte, error)
type DocumentTab ¶
added in v0.190.0
type DocumentTab struct {
	// Body: The main body of the document tab.
	Body *Body `json:"body,omitempty"`
	// CommentAnchors: The comment anchors in a document tab, keyed by anchor ID.
	// Only populated if the commentsViewMode parameter is set to require comments
	// (such as `COMMENTS_VIEW_MODE_INCLUDED`).
	CommentAnchors map[string]CommentAnchor `json:"commentAnchors,omitempty"`
	// DocumentStyle: The style of the document tab.
	DocumentStyle *DocumentStyle `json:"documentStyle,omitempty"`
	// DropdownDefinitions: The dropdown definitions in a document tab, keyed by
	// dropdown definition ID.
	DropdownDefinitions map[string]DropdownDefinition `json:"dropdownDefinitions,omitempty"`
	// Footers: The footers in the document tab, keyed by footer ID.
	Footers map[string]Footer `json:"footers,omitempty"`
	// Footnotes: The footnotes in the document tab, keyed by footnote ID.
	Footnotes map[string]Footnote `json:"footnotes,omitempty"`
	// Headers: The headers in the document tab, keyed by header ID.
	Headers map[string]Header `json:"headers,omitempty"`
	// InlineObjects: The inline objects in the document tab, keyed by object ID.
	InlineObjects map[string]InlineObject `json:"inlineObjects,omitempty"`
	// Lists: The lists in the document tab, keyed by list ID.
	Lists map[string]List `json:"lists,omitempty"`
	// NamedRanges: The named ranges in the document tab, keyed by name.
	NamedRanges map[string]NamedRanges `json:"namedRanges,omitempty"`
	// NamedStyles: The named styles of the document tab.
	NamedStyles *NamedStyles `json:"namedStyles,omitempty"`
	// PositionedObjects: The positioned objects in the document tab, keyed by
	// object ID.
	PositionedObjects map[string]PositionedObject `json:"positionedObjects,omitempty"`
	// SuggestedDocumentStyleChanges: The suggested changes to the style of the
	// document tab, keyed by suggestion ID.
	SuggestedDocumentStyleChanges map[string]SuggestedDocumentStyle `json:"suggestedDocumentStyleChanges,omitempty"`
	// SuggestedNamedStylesChanges: The suggested changes to the named styles of
	// the document tab, keyed by suggestion ID.
	SuggestedNamedStylesChanges map[string]SuggestedNamedStyles `json:"suggestedNamedStylesChanges,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Body") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Body") to include in API requests
	// with the JSON null value. By default, fields with empty values are omitted
	// from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DocumentTab: A tab with document contents.

func (DocumentTab) MarshalJSON ¶
added in v0.190.0
func (s DocumentTab) MarshalJSON() ([]byte, error)
type DocumentsBatchUpdateCall ¶
type DocumentsBatchUpdateCall struct {
	// contains filtered or unexported fields
}
func (*DocumentsBatchUpdateCall) Context ¶
func (c *DocumentsBatchUpdateCall) Context(ctx context.Context) *DocumentsBatchUpdateCall
Context sets the context to be used in this call's Do method.

func (*DocumentsBatchUpdateCall) Do ¶
func (c *DocumentsBatchUpdateCall) Do(opts ...googleapi.CallOption) (*BatchUpdateDocumentResponse, error)
Do executes the "docs.documents.batchUpdate" call. Any non-2xx status code is an error. Response headers are in either *BatchUpdateDocumentResponse.ServerResponse.Header or (if a response was returned at all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to check whether the returned error was because http.StatusNotModified was returned.

func (*DocumentsBatchUpdateCall) Fields ¶
func (c *DocumentsBatchUpdateCall) Fields(s ...googleapi.Field) *DocumentsBatchUpdateCall
Fields allows partial responses to be retrieved. See https://developers.google.com/gdata/docs/2.0/basics#PartialResponse for more details.

func (*DocumentsBatchUpdateCall) Header ¶
func (c *DocumentsBatchUpdateCall) Header() http.Header
Header returns a http.Header that can be modified by the caller to add headers to the request.

type DocumentsCreateCall ¶
type DocumentsCreateCall struct {
	// contains filtered or unexported fields
}
func (*DocumentsCreateCall) Context ¶
func (c *DocumentsCreateCall) Context(ctx context.Context) *DocumentsCreateCall
Context sets the context to be used in this call's Do method.

func (*DocumentsCreateCall) Do ¶
func (c *DocumentsCreateCall) Do(opts ...googleapi.CallOption) (*Document, error)
Do executes the "docs.documents.create" call. Any non-2xx status code is an error. Response headers are in either *Document.ServerResponse.Header or (if a response was returned at all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to check whether the returned error was because http.StatusNotModified was returned.

func (*DocumentsCreateCall) Fields ¶
func (c *DocumentsCreateCall) Fields(s ...googleapi.Field) *DocumentsCreateCall
Fields allows partial responses to be retrieved. See https://developers.google.com/gdata/docs/2.0/basics#PartialResponse for more details.

func (*DocumentsCreateCall) Header ¶
func (c *DocumentsCreateCall) Header() http.Header
Header returns a http.Header that can be modified by the caller to add headers to the request.

type DocumentsGetCall ¶
type DocumentsGetCall struct {
	// contains filtered or unexported fields
}
func (*DocumentsGetCall) CommentsViewMode ¶
func (c *DocumentsGetCall) CommentsViewMode(commentsViewMode string) *DocumentsGetCall
CommentsViewMode sets the optional parameter "commentsViewMode": The comments view mode to apply to the document. This allows viewing the document with comments omitted or included. If one is not specified, COMMENTS_VIEW_MODE_OMITTED is used. If you set comments_view_mode to any value, you must also set include_tabs_content to `true` or use a field mask that references the Document.tabs field (or any subfield). If you set comments_view_mode to COMMENTS_VIEW_MODE_INCLUDED, you must also explicitly set suggestions_view_mode to SUGGESTIONS_INLINE. If you set comments_view_mode to COMMENTS_VIEW_MODE_INCLUDED or COMMENTS_VIEW_MODE_DEFAULT_FOR_CURRENT_ACCESS, you may not set suggestions_view_mode to PREVIEW_WITHOUT_SUGGESTIONS or PREVIEW_SUGGESTIONS_ACCEPTED.

Possible values:

"COMMENTS_VIEW_MODE_UNSPECIFIED" - The CommentsViewMode is unspecified.
COMMENTS_VIEW_MODE_OMITTED is applied.

"COMMENTS_VIEW_MODE_DEFAULT_FOR_CURRENT_ACCESS" - The CommentsViewMode
applied to the returned document depends on the user's current access level. If the user only has view access, COMMENTS_VIEW_MODE_OMITTED is applied. Otherwise, COMMENTS_VIEW_MODE_INCLUDED is applied.

"COMMENTS_VIEW_MODE_OMITTED" - The returned document has comments omitted.
"COMMENTS_VIEW_MODE_INCLUDED" - The returned document has comments
included. Requests to retrieve a document using this mode will return a 403 error if the user does not have permission to view comments. When set, suggestions_view_mode must also be set to SUGGESTIONS_INLINE. Returns a 400 bad request error otherwise.

func (*DocumentsGetCall) Context ¶
func (c *DocumentsGetCall) Context(ctx context.Context) *DocumentsGetCall
Context sets the context to be used in this call's Do method.

func (*DocumentsGetCall) Do ¶
func (c *DocumentsGetCall) Do(opts ...googleapi.CallOption) (*Document, error)
Do executes the "docs.documents.get" call. Any non-2xx status code is an error. Response headers are in either *Document.ServerResponse.Header or (if a response was returned at all) in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to check whether the returned error was because http.StatusNotModified was returned.

func (*DocumentsGetCall) Fields ¶
func (c *DocumentsGetCall) Fields(s ...googleapi.Field) *DocumentsGetCall
Fields allows partial responses to be retrieved. See https://developers.google.com/gdata/docs/2.0/basics#PartialResponse for more details.

func (*DocumentsGetCall) Header ¶
func (c *DocumentsGetCall) Header() http.Header
Header returns a http.Header that can be modified by the caller to add headers to the request.

func (*DocumentsGetCall) IfNoneMatch ¶
func (c *DocumentsGetCall) IfNoneMatch(entityTag string) *DocumentsGetCall
IfNoneMatch sets an optional parameter which makes the operation fail if the object's ETag matches the given value. This is useful for getting updates only after the object has changed since the last request.

func (*DocumentsGetCall) IncludeTabsContent ¶
added in v0.190.0
func (c *DocumentsGetCall) IncludeTabsContent(includeTabsContent bool) *DocumentsGetCall
IncludeTabsContent sets the optional parameter "includeTabsContent": Whether to populate the `Document.tabs` field instead of the text content fields like `body` and `documentStyle` on `Document`. - When `true`: Document content populates in the `Document.tabs` field instead of the text content fields in `Document`. - When `false`: The content of the document's first tab populates the content fields in `Document` excluding `Document.tabs`. If a document has only one tab, then that tab is used to populate the document content. `Document.tabs` will be empty. If you use a field mask that references the `Document.tabs` field (or any subfield), the API implicitly treats the request as if you set `include_tabs_content` to `true`.

func (*DocumentsGetCall) SuggestionsViewMode ¶
func (c *DocumentsGetCall) SuggestionsViewMode(suggestionsViewMode string) *DocumentsGetCall
SuggestionsViewMode sets the optional parameter "suggestionsViewMode": The suggestions view mode to apply to the document. This allows viewing the document with all suggestions inline, accepted or rejected. If one is not specified, DEFAULT_FOR_CURRENT_ACCESS is used.

Possible values:

"DEFAULT_FOR_CURRENT_ACCESS" - The SuggestionsViewMode applied to the
returned document depends on the user's current access level. If the user only has view access, PREVIEW_WITHOUT_SUGGESTIONS is applied. Otherwise, SUGGESTIONS_INLINE is applied. This is the default suggestions view mode.

"SUGGESTIONS_INLINE" - The returned document has suggestions inline.
Suggested changes will be differentiated from base content within the document. Requests to retrieve a document using this mode will return a 403 error if the user does not have permission to view suggested changes.

"PREVIEW_SUGGESTIONS_ACCEPTED" - The returned document is a preview with
all suggested changes accepted. Requests to retrieve a document using this mode will return a 403 error if the user does not have permission to view suggested changes.

"PREVIEW_WITHOUT_SUGGESTIONS" - The returned document is a preview with
all suggested changes rejected if there are any suggestions in the document.

type DocumentsService ¶
type DocumentsService struct {
	// contains filtered or unexported fields
}
func NewDocumentsService ¶
func NewDocumentsService(s *Service) *DocumentsService
func (*DocumentsService) BatchUpdate ¶
func (r *DocumentsService) BatchUpdate(documentId string, batchupdatedocumentrequest *BatchUpdateDocumentRequest) *DocumentsBatchUpdateCall
BatchUpdate: Applies one or more updates to the document. Each request is validated before being applied. If any request is not valid, then the entire request will fail and nothing will be applied. Some requests have replies to give you some information about how they are applied. Other requests do not need to return information; these each return an empty reply. The order of replies matches that of the requests. For example, suppose you call batchUpdate with four updates, and only the third one returns information. The response would have two empty replies, the reply to the third request, and another empty reply, in that order. Because other users may be editing the document, the document might not exactly reflect your changes: your changes may be altered with respect to collaborator changes. If there are no collaborators, the document should reflect your changes. In any case, the updates in your request are guaranteed to be applied together atomically.

- documentId: The ID of the document to update.

func (*DocumentsService) Create ¶
func (r *DocumentsService) Create(document *Document) *DocumentsCreateCall
Create: Creates a blank document using the title given in the request. Other fields in the request, including any provided content, are ignored. Returns the created document.

func (*DocumentsService) Get ¶
func (r *DocumentsService) Get(documentId string) *DocumentsGetCall
Get: Gets the latest version of the specified document.

- documentId: The ID of the document to retrieve.

type Dropdown ¶
type Dropdown struct {
	// DropdownId: The ID of this dropdown.
	DropdownId string `json:"dropdownId,omitempty"`
	// DropdownProperties: The properties of this dropdown.
	DropdownProperties *DropdownProperties `json:"dropdownProperties,omitempty"`
	// SuggestedDeletionIds: IDs for suggestions that remove this dropdown from the
	// document. If empty, then this dropdown isn't suggested for deletion.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedDropdownPropertiesChanges: The suggested properties changes to this
	// dropdown, keyed by suggestion ID.
	SuggestedDropdownPropertiesChanges map[string]SuggestedDropdownProperties `json:"suggestedDropdownPropertiesChanges,omitempty"`
	// SuggestedInsertionIds: IDs for suggestions that insert this dropdown into
	// the document. If empty, then this dropdown isn't a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTextStyleChanges: The suggested text style changes to this
	// dropdown, keyed by suggestion ID.
	SuggestedTextStyleChanges map[string]SuggestedTextStyle `json:"suggestedTextStyleChanges,omitempty"`
	// TextStyle: The text style of this dropdown.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DropdownId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DropdownId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Dropdown: A dropdown in the document. The chip is displayed as a dropdown menu that allows users to select an option from a configurable list of options.

func (Dropdown) MarshalJSON ¶
func (s Dropdown) MarshalJSON() ([]byte, error)
type DropdownDefinition ¶
type DropdownDefinition struct {
	// DropdownDefinitionId: The ID of this dropdown definition. If you specify an
	// ID, it must be unique among all IDs in the tab. The ID must start with
	// `kix.` and match regex `^kix\.[a-zA-Z0-9_-]{2,14}$` (length 6-18 chars). If
	// you don't specify an ID, a unique one is generated.
	DropdownDefinitionId string `json:"dropdownDefinitionId,omitempty"`
	// DropdownDefinitionProperties: The properties of this dropdown definition.
	DropdownDefinitionProperties *DropdownDefinitionProperties `json:"dropdownDefinitionProperties,omitempty"`
	// SuggestedDeletionId: ID for suggestion that deletes this dropdown
	// definition.
	SuggestedDeletionId string `json:"suggestedDeletionId,omitempty"`
	// SuggestedDropdownDefinitionPropertiesChanges: Suggested property changes to
	// this definition, keyed by suggestion ID.
	SuggestedDropdownDefinitionPropertiesChanges map[string]SuggestedDropdownDefinitionProperties `json:"suggestedDropdownDefinitionPropertiesChanges,omitempty"`
	// SuggestedInsertionId: ID for suggestion that inserts this dropdown
	// definition.
	SuggestedInsertionId string `json:"suggestedInsertionId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DropdownDefinitionId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DropdownDefinitionId") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DropdownDefinition: A dropdown definition in the document.

func (DropdownDefinition) MarshalJSON ¶
func (s DropdownDefinition) MarshalJSON() ([]byte, error)
type DropdownDefinitionProperties ¶
type DropdownDefinitionProperties struct {
	// Options: The list of options defined by this dropdown definition. A dropdown
	// definition must have at least 2 options and at most 50 options.
	Options []*DropdownOption `json:"options,omitempty"`
	// Title: The title of the dropdown definition.
	Title string `json:"title,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Options") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Options") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DropdownDefinitionProperties: Properties of a dropdown definition.

func (DropdownDefinitionProperties) MarshalJSON ¶
func (s DropdownDefinitionProperties) MarshalJSON() ([]byte, error)
type DropdownDefinitionPropertiesSuggestionState ¶
type DropdownDefinitionPropertiesSuggestionState struct {
	// OptionsSuggested: Indicates if there was a suggested change to options.
	OptionsSuggested bool `json:"optionsSuggested,omitempty"`
	// TitleSuggested: Indicates if there was a suggested change to title.
	TitleSuggested bool `json:"titleSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "OptionsSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "OptionsSuggested") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DropdownDefinitionPropertiesSuggestionState: A mask that indicates which of the fields on the base DropdownDefinitionProperties have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (DropdownDefinitionPropertiesSuggestionState) MarshalJSON ¶
func (s DropdownDefinitionPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type DropdownOption ¶
type DropdownOption struct {
	// DisplayValue: The display value of this dropdown option.
	DisplayValue string `json:"displayValue,omitempty"`
	// OptionId: The ID of this dropdown option. If you specify an ID, it must be
	// unique among all options in this dropdown definition. The ID must start with
	// `dropdownItem.` and match regex `^dropdownItem\.[a-zA-Z0-9_-]{2,14}$`
	// (length 15-27 chars). If you don't specify an ID, a unique one is generated.
	OptionId string `json:"optionId,omitempty"`
	// TextStyle: The text style of this dropdown option. Currently, only the
	// `foreground_color` and `background_color` properties are supported. If other
	// properties are set, a 400 bad request error is returned.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DisplayValue") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DisplayValue") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DropdownOption: An option in a Dropdown.

func (DropdownOption) MarshalJSON ¶
func (s DropdownOption) MarshalJSON() ([]byte, error)
type DropdownProperties ¶
type DropdownProperties struct {
	// DisplayValue: The human-readable display text of the currently selected
	// item. This field is populated by the server based on the dropdown definition
	// and the selected option ID. It may differ from
	// `DropdownOption.display_value` if the underlying option definition was
	// modified or deleted, or during pending suggested changes.
	DisplayValue string `json:"displayValue,omitempty"`
	// DropdownDefinitionId: The ID of the DropdownDefinition that defines the
	// options for this dropdown.
	DropdownDefinitionId string `json:"dropdownDefinitionId,omitempty"`
	// SelectedOptionId: The ID of the selected option in this dropdown.
	SelectedOptionId string `json:"selectedOptionId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DisplayValue") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DisplayValue") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DropdownProperties: Properties specific to a dropdown.

func (DropdownProperties) MarshalJSON ¶
func (s DropdownProperties) MarshalJSON() ([]byte, error)
type DropdownPropertiesSuggestionState ¶
type DropdownPropertiesSuggestionState struct {
	// SelectedOptionIdSuggested: Indicates if there was a suggested change to
	// selected_option_id.
	SelectedOptionIdSuggested bool `json:"selectedOptionIdSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SelectedOptionIdSuggested")
	// to unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SelectedOptionIdSuggested") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
DropdownPropertiesSuggestionState: A mask that indicates which of the fields on the base DropdownProperties have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (DropdownPropertiesSuggestionState) MarshalJSON ¶
func (s DropdownPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type EmbeddedDrawingProperties ¶
type EmbeddedDrawingProperties struct {
}
EmbeddedDrawingProperties: The properties of an embedded drawing and used to differentiate the object type. An embedded drawing is one that's created and edited within a document. Note that extensive details are not supported.

type EmbeddedDrawingPropertiesSuggestionState ¶
type EmbeddedDrawingPropertiesSuggestionState struct {
}
EmbeddedDrawingPropertiesSuggestionState: A mask that indicates which of the fields on the base EmbeddedDrawingProperties have been changed in this suggestion. For any field set to true, there's a new suggested value.

type EmbeddedObject ¶
type EmbeddedObject struct {
	// Description: The description of the embedded object. The `title` and
	// `description` are both combined to display alt text.
	Description string `json:"description,omitempty"`
	// EmbeddedDrawingProperties: The properties of an embedded drawing.
	EmbeddedDrawingProperties *EmbeddedDrawingProperties `json:"embeddedDrawingProperties,omitempty"`
	// EmbeddedObjectBorder: The border of the embedded object.
	EmbeddedObjectBorder *EmbeddedObjectBorder `json:"embeddedObjectBorder,omitempty"`
	// ImageProperties: The properties of an image.
	ImageProperties *ImageProperties `json:"imageProperties,omitempty"`
	// LinkedContentReference: A reference to the external linked source content.
	// For example, it contains a reference to the source Google Sheets chart when
	// the embedded object is a linked chart. If unset, then the embedded object is
	// not linked.
	LinkedContentReference *LinkedContentReference `json:"linkedContentReference,omitempty"`
	// MarginBottom: The bottom margin of the embedded object.
	MarginBottom *Dimension `json:"marginBottom,omitempty"`
	// MarginLeft: The left margin of the embedded object.
	MarginLeft *Dimension `json:"marginLeft,omitempty"`
	// MarginRight: The right margin of the embedded object.
	MarginRight *Dimension `json:"marginRight,omitempty"`
	// MarginTop: The top margin of the embedded object.
	MarginTop *Dimension `json:"marginTop,omitempty"`
	// Size: The visible size of the image after cropping.
	Size *Size `json:"size,omitempty"`
	// Title: The title of the embedded object. The `title` and `description` are
	// both combined to display alt text.
	Title string `json:"title,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Description") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Description") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
EmbeddedObject: An embedded object in the document.

func (EmbeddedObject) MarshalJSON ¶
func (s EmbeddedObject) MarshalJSON() ([]byte, error)
type EmbeddedObjectBorder ¶
type EmbeddedObjectBorder struct {
	// Color: The color of the border.
	Color *OptionalColor `json:"color,omitempty"`
	// DashStyle: The dash style of the border.
	//
	// Possible values:
	//   "DASH_STYLE_UNSPECIFIED" - Unspecified dash style.
	//   "SOLID" - Solid line. Corresponds to ECMA-376 ST_PresetLineDashVal value
	// 'solid'. This is the default dash style.
	//   "DOT" - Dotted line. Corresponds to ECMA-376 ST_PresetLineDashVal value
	// 'dot'.
	//   "DASH" - Dashed line. Corresponds to ECMA-376 ST_PresetLineDashVal value
	// 'dash'.
	DashStyle string `json:"dashStyle,omitempty"`
	// PropertyState: The property state of the border property.
	//
	// Possible values:
	//   "RENDERED" - If a property's state is RENDERED, then the element has the
	// corresponding property when rendered in the document. This is the default
	// value.
	//   "NOT_RENDERED" - If a property's state is NOT_RENDERED, then the element
	// does not have the corresponding property when rendered in the document.
	PropertyState string `json:"propertyState,omitempty"`
	// Width: The width of the border.
	Width *Dimension `json:"width,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Color") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Color") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
EmbeddedObjectBorder: A border around an EmbeddedObject.

func (EmbeddedObjectBorder) MarshalJSON ¶
func (s EmbeddedObjectBorder) MarshalJSON() ([]byte, error)
type EmbeddedObjectBorderSuggestionState ¶
type EmbeddedObjectBorderSuggestionState struct {
	// ColorSuggested: Indicates if there was a suggested change to color.
	ColorSuggested bool `json:"colorSuggested,omitempty"`
	// DashStyleSuggested: Indicates if there was a suggested change to dash_style.
	DashStyleSuggested bool `json:"dashStyleSuggested,omitempty"`
	// PropertyStateSuggested: Indicates if there was a suggested change to
	// property_state.
	PropertyStateSuggested bool `json:"propertyStateSuggested,omitempty"`
	// WidthSuggested: Indicates if there was a suggested change to width.
	WidthSuggested bool `json:"widthSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ColorSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ColorSuggested") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
EmbeddedObjectBorderSuggestionState: A mask that indicates which of the fields on the base EmbeddedObjectBorder have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (EmbeddedObjectBorderSuggestionState) MarshalJSON ¶
func (s EmbeddedObjectBorderSuggestionState) MarshalJSON() ([]byte, error)
type EmbeddedObjectSuggestionState ¶
type EmbeddedObjectSuggestionState struct {
	// DescriptionSuggested: Indicates if there was a suggested change to
	// description.
	DescriptionSuggested bool `json:"descriptionSuggested,omitempty"`
	// EmbeddedDrawingPropertiesSuggestionState: A mask that indicates which of the
	// fields in embedded_drawing_properties have been changed in this suggestion.
	EmbeddedDrawingPropertiesSuggestionState *EmbeddedDrawingPropertiesSuggestionState `json:"embeddedDrawingPropertiesSuggestionState,omitempty"`
	// EmbeddedObjectBorderSuggestionState: A mask that indicates which of the
	// fields in embedded_object_border have been changed in this suggestion.
	EmbeddedObjectBorderSuggestionState *EmbeddedObjectBorderSuggestionState `json:"embeddedObjectBorderSuggestionState,omitempty"`
	// ImagePropertiesSuggestionState: A mask that indicates which of the fields in
	// image_properties have been changed in this suggestion.
	ImagePropertiesSuggestionState *ImagePropertiesSuggestionState `json:"imagePropertiesSuggestionState,omitempty"`
	// LinkedContentReferenceSuggestionState: A mask that indicates which of the
	// fields in linked_content_reference have been changed in this suggestion.
	LinkedContentReferenceSuggestionState *LinkedContentReferenceSuggestionState `json:"linkedContentReferenceSuggestionState,omitempty"`
	// MarginBottomSuggested: Indicates if there was a suggested change to
	// margin_bottom.
	MarginBottomSuggested bool `json:"marginBottomSuggested,omitempty"`
	// MarginLeftSuggested: Indicates if there was a suggested change to
	// margin_left.
	MarginLeftSuggested bool `json:"marginLeftSuggested,omitempty"`
	// MarginRightSuggested: Indicates if there was a suggested change to
	// margin_right.
	MarginRightSuggested bool `json:"marginRightSuggested,omitempty"`
	// MarginTopSuggested: Indicates if there was a suggested change to margin_top.
	MarginTopSuggested bool `json:"marginTopSuggested,omitempty"`
	// SizeSuggestionState: A mask that indicates which of the fields in size have
	// been changed in this suggestion.
	SizeSuggestionState *SizeSuggestionState `json:"sizeSuggestionState,omitempty"`
	// TitleSuggested: Indicates if there was a suggested change to title.
	TitleSuggested bool `json:"titleSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DescriptionSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DescriptionSuggested") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
EmbeddedObjectSuggestionState: A mask that indicates which of the fields on the base EmbeddedObject have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (EmbeddedObjectSuggestionState) MarshalJSON ¶
func (s EmbeddedObjectSuggestionState) MarshalJSON() ([]byte, error)
type EndOfSegmentLocation ¶
type EndOfSegmentLocation struct {
	// SegmentId: The ID of the header, footer or footnote the location is in. An
	// empty segment ID signifies the document's body.
	SegmentId string `json:"segmentId,omitempty"`
	// TabId: The tab that the location is in. When omitted, the request is applied
	// to the first tab. In a document containing a single tab: - If provided, must
	// match the singular tab's ID. - If omitted, the request applies to the
	// singular tab. In a document containing multiple tabs: - If provided, the
	// request applies to the specified tab. - If omitted, the request applies to
	// the first tab in the document.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SegmentId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SegmentId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
EndOfSegmentLocation: Location at the end of a body, header, footer or footnote. The location is immediately before the last newline in the document segment.

func (EndOfSegmentLocation) MarshalJSON ¶
func (s EndOfSegmentLocation) MarshalJSON() ([]byte, error)
type Equation ¶
type Equation struct {
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. An Equation may have
	// multiple insertion IDs if it's a nested suggested change. If empty, then
	// this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SuggestedDeletionIds") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SuggestedDeletionIds") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Equation: A ParagraphElement representing an equation.

func (Equation) MarshalJSON ¶
func (s Equation) MarshalJSON() ([]byte, error)
type Footer ¶
type Footer struct {
	// Content: The contents of the footer. The indexes for a footer's content
	// begin at zero.
	Content []*StructuralElement `json:"content,omitempty"`
	// FooterId: The ID of the footer.
	FooterId string `json:"footerId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Content") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Content") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Footer: A document footer.

func (Footer) MarshalJSON ¶
func (s Footer) MarshalJSON() ([]byte, error)
type Footnote ¶
type Footnote struct {
	// Content: The contents of the footnote. The indexes for a footnote's content
	// begin at zero.
	Content []*StructuralElement `json:"content,omitempty"`
	// FootnoteId: The ID of the footnote.
	FootnoteId string `json:"footnoteId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Content") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Content") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Footnote: A document footnote.

func (Footnote) MarshalJSON ¶
func (s Footnote) MarshalJSON() ([]byte, error)
type FootnoteReference ¶
type FootnoteReference struct {
	// FootnoteId: The ID of the footnote that contains the content of this
	// footnote reference.
	FootnoteId string `json:"footnoteId,omitempty"`
	// FootnoteNumber: The rendered number of this footnote.
	FootnoteNumber string `json:"footnoteNumber,omitempty"`
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. A FootnoteReference may
	// have multiple insertion IDs if it's a nested suggested change. If empty,
	// then this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTextStyleChanges: The suggested text style changes to this
	// FootnoteReference, keyed by suggestion ID.
	SuggestedTextStyleChanges map[string]SuggestedTextStyle `json:"suggestedTextStyleChanges,omitempty"`
	// TextStyle: The text style of this FootnoteReference.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "FootnoteId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "FootnoteId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
FootnoteReference: A ParagraphElement representing a footnote reference. A footnote reference is the inline content rendered with a number and is used to identify the footnote.

func (FootnoteReference) MarshalJSON ¶
func (s FootnoteReference) MarshalJSON() ([]byte, error)
type Header ¶
type Header struct {
	// Content: The contents of the header. The indexes for a header's content
	// begin at zero.
	Content []*StructuralElement `json:"content,omitempty"`
	// HeaderId: The ID of the header.
	HeaderId string `json:"headerId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Content") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Content") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Header: A document header.

func (Header) MarshalJSON ¶
func (s Header) MarshalJSON() ([]byte, error)
type HeadingLink ¶
added in v0.190.0
type HeadingLink struct {
	// Id: The ID of a heading in this document.
	Id string `json:"id,omitempty"`
	// TabId: The ID of the tab containing this heading.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Id") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Id") to include in API requests
	// with the JSON null value. By default, fields with empty values are omitted
	// from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
HeadingLink: A reference to a heading in this document.

func (HeadingLink) MarshalJSON ¶
added in v0.190.0
func (s HeadingLink) MarshalJSON() ([]byte, error)
type HorizontalRule ¶
type HorizontalRule struct {
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. A HorizontalRule may
	// have multiple insertion IDs if it is a nested suggested change. If empty,
	// then this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTextStyleChanges: The suggested text style changes to this
	// HorizontalRule, keyed by suggestion ID.
	SuggestedTextStyleChanges map[string]SuggestedTextStyle `json:"suggestedTextStyleChanges,omitempty"`
	// TextStyle: The text style of this HorizontalRule. Similar to text content,
	// like text runs and footnote references, the text style of a horizontal rule
	// can affect content layout as well as the styling of text inserted next to
	// it.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SuggestedDeletionIds") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SuggestedDeletionIds") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
HorizontalRule: A ParagraphElement representing a horizontal line.

func (HorizontalRule) MarshalJSON ¶
func (s HorizontalRule) MarshalJSON() ([]byte, error)
type ImageProperties ¶
type ImageProperties struct {
	// Angle: The clockwise rotation angle of the image, in radians.
	Angle float64 `json:"angle,omitempty"`
	// Brightness: The brightness effect of the image. The value should be in the
	// interval [-1.0, 1.0], where 0 means no effect.
	Brightness float64 `json:"brightness,omitempty"`
	// ContentUri: A URI to the image with a default lifetime of 30 minutes. This
	// URI is tagged with the account of the requester. Anyone with the URI
	// effectively accesses the image as the original requester. Access to the
	// image may be lost if the document's sharing settings change.
	ContentUri string `json:"contentUri,omitempty"`
	// Contrast: The contrast effect of the image. The value should be in the
	// interval [-1.0, 1.0], where 0 means no effect.
	Contrast float64 `json:"contrast,omitempty"`
	// CropProperties: The crop properties of the image.
	CropProperties *CropProperties `json:"cropProperties,omitempty"`
	// SourceUri: The source URI is the URI used to insert the image. The source
	// URI can be empty.
	SourceUri string `json:"sourceUri,omitempty"`
	// Transparency: The transparency effect of the image. The value should be in
	// the interval [0.0, 1.0], where 0 means no effect and 1 means transparent.
	Transparency float64 `json:"transparency,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Angle") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Angle") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ImageProperties: The properties of an image.

func (ImageProperties) MarshalJSON ¶
func (s ImageProperties) MarshalJSON() ([]byte, error)
func (*ImageProperties) UnmarshalJSON ¶
func (s *ImageProperties) UnmarshalJSON(data []byte) error
type ImagePropertiesSuggestionState ¶
type ImagePropertiesSuggestionState struct {
	// AngleSuggested: Indicates if there was a suggested change to angle.
	AngleSuggested bool `json:"angleSuggested,omitempty"`
	// BrightnessSuggested: Indicates if there was a suggested change to
	// brightness.
	BrightnessSuggested bool `json:"brightnessSuggested,omitempty"`
	// ContentUriSuggested: Indicates if there was a suggested change to
	// content_uri.
	ContentUriSuggested bool `json:"contentUriSuggested,omitempty"`
	// ContrastSuggested: Indicates if there was a suggested change to contrast.
	ContrastSuggested bool `json:"contrastSuggested,omitempty"`
	// CropPropertiesSuggestionState: A mask that indicates which of the fields in
	// crop_properties have been changed in this suggestion.
	CropPropertiesSuggestionState *CropPropertiesSuggestionState `json:"cropPropertiesSuggestionState,omitempty"`
	// SourceUriSuggested: Indicates if there was a suggested change to source_uri.
	SourceUriSuggested bool `json:"sourceUriSuggested,omitempty"`
	// TransparencySuggested: Indicates if there was a suggested change to
	// transparency.
	TransparencySuggested bool `json:"transparencySuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "AngleSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "AngleSuggested") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ImagePropertiesSuggestionState: A mask that indicates which of the fields on the base ImageProperties have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (ImagePropertiesSuggestionState) MarshalJSON ¶
func (s ImagePropertiesSuggestionState) MarshalJSON() ([]byte, error)
type InlineObject ¶
type InlineObject struct {
	// InlineObjectProperties: The properties of this inline object.
	InlineObjectProperties *InlineObjectProperties `json:"inlineObjectProperties,omitempty"`
	// ObjectId: The ID of this inline object. Can be used to update an object’s
	// properties.
	ObjectId string `json:"objectId,omitempty"`
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInlineObjectPropertiesChanges: The suggested changes to the inline
	// object properties, keyed by suggestion ID.
	SuggestedInlineObjectPropertiesChanges map[string]SuggestedInlineObjectProperties `json:"suggestedInlineObjectPropertiesChanges,omitempty"`
	// SuggestedInsertionId: The suggested insertion ID. If empty, then this is not
	// a suggested insertion.
	SuggestedInsertionId string `json:"suggestedInsertionId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "InlineObjectProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "InlineObjectProperties") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InlineObject: An object that appears inline with text. An InlineObject contains an EmbeddedObject such as an image.

func (InlineObject) MarshalJSON ¶
func (s InlineObject) MarshalJSON() ([]byte, error)
type InlineObjectElement ¶
type InlineObjectElement struct {
	// InlineObjectId: The ID of the InlineObject this element contains.
	InlineObjectId string `json:"inlineObjectId,omitempty"`
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. An InlineObjectElement
	// may have multiple insertion IDs if it's a nested suggested change. If empty,
	// then this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTextStyleChanges: The suggested text style changes to this
	// InlineObject, keyed by suggestion ID.
	SuggestedTextStyleChanges map[string]SuggestedTextStyle `json:"suggestedTextStyleChanges,omitempty"`
	// TextStyle: The text style of this InlineObjectElement. Similar to text
	// content, like text runs and footnote references, the text style of an inline
	// object element can affect content layout as well as the styling of text
	// inserted next to it.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "InlineObjectId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "InlineObjectId") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InlineObjectElement: A ParagraphElement that contains an InlineObject.

func (InlineObjectElement) MarshalJSON ¶
func (s InlineObjectElement) MarshalJSON() ([]byte, error)
type InlineObjectProperties ¶
type InlineObjectProperties struct {
	// EmbeddedObject: The embedded object of this inline object.
	EmbeddedObject *EmbeddedObject `json:"embeddedObject,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EmbeddedObject") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EmbeddedObject") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InlineObjectProperties: Properties of an InlineObject.

func (InlineObjectProperties) MarshalJSON ¶
func (s InlineObjectProperties) MarshalJSON() ([]byte, error)
type InlineObjectPropertiesSuggestionState ¶
type InlineObjectPropertiesSuggestionState struct {
	// EmbeddedObjectSuggestionState: A mask that indicates which of the fields in
	// embedded_object have been changed in this suggestion.
	EmbeddedObjectSuggestionState *EmbeddedObjectSuggestionState `json:"embeddedObjectSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g.
	// "EmbeddedObjectSuggestionState") to unconditionally include in API requests.
	// By default, fields with empty or default values are omitted from API
	// requests. See https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields
	// for more details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EmbeddedObjectSuggestionState")
	// to include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InlineObjectPropertiesSuggestionState: A mask that indicates which of the fields on the base InlineObjectProperties have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (InlineObjectPropertiesSuggestionState) MarshalJSON ¶
func (s InlineObjectPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type InsertCommentRequest ¶
type InsertCommentRequest struct {
	// AssigneeEmailAddress: Optional. The email address of the assignee of the
	// comment. Leave empty for a non-assigned comment. May not exceed 2048 UTF-8
	// code units.
	AssigneeEmailAddress string `json:"assigneeEmailAddress,omitempty"`
	// Content: The text of the comment, as plain text. This text content will be
	// handled similarly to comments created in the Docs editor. It will have
	// similar behaviors for formatting, notifications, etc. This field cannot be
	// empty, and must not exceed 2048 UTF-8 code units.
	Content string `json:"content,omitempty"`
	// Range: The Range in the document that is tied to this comment.
	Range *Range `json:"range,omitempty"`
	// ForceSendFields is a list of field names (e.g. "AssigneeEmailAddress") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "AssigneeEmailAddress") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertCommentRequest: Inserts a CommentThread into the document.

func (InsertCommentRequest) MarshalJSON ¶
func (s InsertCommentRequest) MarshalJSON() ([]byte, error)
type InsertCommentResponse ¶
type InsertCommentResponse struct {
	// CommentThread: The newly-inserted comment thread.
	CommentThread *CommentThread `json:"commentThread,omitempty"`
	// ForceSendFields is a list of field names (e.g. "CommentThread") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "CommentThread") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertCommentResponse: Response message for inserting a comment.

func (InsertCommentResponse) MarshalJSON ¶
func (s InsertCommentResponse) MarshalJSON() ([]byte, error)
type InsertDateRequest ¶
added in v0.258.0
type InsertDateRequest struct {
	// DateElementProperties: The properties of the date to insert.
	DateElementProperties *DateElementProperties `json:"dateElementProperties,omitempty"`
	// EndOfSegmentLocation: Inserts the date at the end of the given header,
	// footer or document body.
	EndOfSegmentLocation *EndOfSegmentLocation `json:"endOfSegmentLocation,omitempty"`
	// Location: Inserts the date at a specific index in the document. The date
	// must be inserted inside the bounds of an existing Paragraph. For instance,
	// it cannot be inserted at a table's start index (i.e. between an existing
	// table and its preceding paragraph).
	Location *Location `json:"location,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DateElementProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DateElementProperties") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertDateRequest: Inserts a date at the specified location.

func (InsertDateRequest) MarshalJSON ¶
added in v0.258.0
func (s InsertDateRequest) MarshalJSON() ([]byte, error)
type InsertDropdownRequest ¶
type InsertDropdownRequest struct {
	// DropdownDefinitionId: Required. The DropdownDefinition ID.
	DropdownDefinitionId string `json:"dropdownDefinitionId,omitempty"`
	// EndOfSegmentLocation: The EndOfSegmentLocation in the document to insert the
	// dropdown at.
	EndOfSegmentLocation *EndOfSegmentLocation `json:"endOfSegmentLocation,omitempty"`
	// Location: The Location in the document to insert the dropdown at.
	Location *Location `json:"location,omitempty"`
	// SelectedOptionId: Optional initial value for the dropdown. If this field is
	// not specified, the new dropdown will default to selecting the first option
	// defined in the dropdown definition. If this field is specified but does not
	// reference a valid option in the dropdown definition, a 400 bad request error
	// is returned.
	SelectedOptionId string `json:"selectedOptionId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DropdownDefinitionId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DropdownDefinitionId") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertDropdownRequest: Inserts a Dropdown at the specified location.

func (InsertDropdownRequest) MarshalJSON ¶
func (s InsertDropdownRequest) MarshalJSON() ([]byte, error)
type InsertDropdownResponse ¶
type InsertDropdownResponse struct {
	// Dropdown: The newly-inserted Dropdown.
	Dropdown *Dropdown `json:"dropdown,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Dropdown") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Dropdown") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertDropdownResponse: The result of inserting a Dropdown.

func (InsertDropdownResponse) MarshalJSON ¶
func (s InsertDropdownResponse) MarshalJSON() ([]byte, error)
type InsertInlineImageRequest ¶
type InsertInlineImageRequest struct {
	// EndOfSegmentLocation: Inserts the text at the end of a header, footer or the
	// document body. Inline images cannot be inserted inside a footnote.
	EndOfSegmentLocation *EndOfSegmentLocation `json:"endOfSegmentLocation,omitempty"`
	// Location: Inserts the image at a specific index in the document. The image
	// must be inserted inside the bounds of an existing Paragraph. For instance,
	// it cannot be inserted at a table's start index (i.e. between the table and
	// its preceding paragraph). Inline images cannot be inserted inside a footnote
	// or equation.
	Location *Location `json:"location,omitempty"`
	// ObjectSize: The size that the image should appear as in the document. This
	// property is optional and the final size of the image in the document is
	// determined by the following rules: * If neither width nor height is
	// specified, then a default size of the image is calculated based on its
	// resolution. * If one dimension is specified then the other dimension is
	// calculated to preserve the aspect ratio of the image. * If both width and
	// height are specified, the image is scaled to fit within the provided
	// dimensions while maintaining its aspect ratio.
	ObjectSize *Size `json:"objectSize,omitempty"`
	// Uri: The image URI. The image is fetched once at insertion time and a copy
	// is stored for display inside the document. Images must be less than 50MB in
	// size, cannot exceed 25 megapixels, and must be in one of PNG, JPEG, or GIF
	// format. The provided URI must be publicly accessible and at most 2 kB in
	// length. The URI itself is saved with the image, and exposed via the
	// ImageProperties.content_uri field.
	Uri string `json:"uri,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EndOfSegmentLocation") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EndOfSegmentLocation") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertInlineImageRequest: Inserts an InlineObject containing an image at the given location.

func (InsertInlineImageRequest) MarshalJSON ¶
func (s InsertInlineImageRequest) MarshalJSON() ([]byte, error)
type InsertInlineImageResponse ¶
type InsertInlineImageResponse struct {
	// ObjectId: The ID of the created InlineObject.
	ObjectId string `json:"objectId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ObjectId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ObjectId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertInlineImageResponse: The result of inserting an inline image.

func (InsertInlineImageResponse) MarshalJSON ¶
func (s InsertInlineImageResponse) MarshalJSON() ([]byte, error)
type InsertInlineSheetsChartResponse ¶
type InsertInlineSheetsChartResponse struct {
	// ObjectId: The object ID of the inserted chart.
	ObjectId string `json:"objectId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ObjectId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ObjectId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertInlineSheetsChartResponse: The result of inserting an embedded Google Sheets chart.

func (InsertInlineSheetsChartResponse) MarshalJSON ¶
func (s InsertInlineSheetsChartResponse) MarshalJSON() ([]byte, error)
type InsertPageBreakRequest ¶
added in v0.3.0
type InsertPageBreakRequest struct {
	// EndOfSegmentLocation: Inserts the page break at the end of the document
	// body. Page breaks cannot be inserted inside a footnote, header or footer.
	// Since page breaks can only be inserted inside the body, the segment ID field
	// must be empty.
	EndOfSegmentLocation *EndOfSegmentLocation `json:"endOfSegmentLocation,omitempty"`
	// Location: Inserts the page break at a specific index in the document. The
	// page break must be inserted inside the bounds of an existing Paragraph. For
	// instance, it cannot be inserted at a table's start index (i.e. between the
	// table and its preceding paragraph). Page breaks cannot be inserted inside a
	// table, equation, footnote, header or footer. Since page breaks can only be
	// inserted inside the body, the segment ID field must be empty.
	Location *Location `json:"location,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EndOfSegmentLocation") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EndOfSegmentLocation") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertPageBreakRequest: Inserts a page break followed by a newline at the specified location.

func (InsertPageBreakRequest) MarshalJSON ¶
added in v0.3.0
func (s InsertPageBreakRequest) MarshalJSON() ([]byte, error)
type InsertPersonRequest ¶
added in v0.252.0
type InsertPersonRequest struct {
	// EndOfSegmentLocation: Inserts the person mention at the end of a header,
	// footer, footnote or the document body.
	EndOfSegmentLocation *EndOfSegmentLocation `json:"endOfSegmentLocation,omitempty"`
	// Location: Inserts the person mention at a specific index in the document.
	// The person mention must be inserted inside the bounds of an existing
	// Paragraph. For instance, it cannot be inserted at a table's start index
	// (i.e. between the table and its preceding paragraph). Person mentions cannot
	// be inserted inside an equation.
	Location *Location `json:"location,omitempty"`
	// PersonProperties: The properties of the person mention to insert.
	PersonProperties *PersonProperties `json:"personProperties,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EndOfSegmentLocation") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EndOfSegmentLocation") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertPersonRequest: Inserts a person mention.

func (InsertPersonRequest) MarshalJSON ¶
added in v0.252.0
func (s InsertPersonRequest) MarshalJSON() ([]byte, error)
type InsertRichLinkRequest ¶
added in v0.276.0
type InsertRichLinkRequest struct {
	// EndOfSegmentLocation: Inserts the rich link at the end of a header, footer,
	// footnote or the document body.
	EndOfSegmentLocation *EndOfSegmentLocation `json:"endOfSegmentLocation,omitempty"`
	// Location: Inserts the rich link at a specific index in the document. The
	// rich link must be inserted inside the bounds of an existing Paragraph. For
	// instance, it cannot be inserted at a table's start index (i.e. between the
	// table and its preceding paragraph). The rich link cannot be inserted inside
	// an equation.
	Location *Location `json:"location,omitempty"`
	// RichLinkProperties: The properties of the rich link to insert.
	RichLinkProperties *RichLinkProperties `json:"richLinkProperties,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EndOfSegmentLocation") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EndOfSegmentLocation") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertRichLinkRequest: Inserts a RichLink at the specified location.

func (InsertRichLinkRequest) MarshalJSON ¶
added in v0.276.0
func (s InsertRichLinkRequest) MarshalJSON() ([]byte, error)
type InsertSectionBreakRequest ¶
added in v0.15.0
type InsertSectionBreakRequest struct {
	// EndOfSegmentLocation: Inserts a newline and a section break at the end of
	// the document body. Section breaks cannot be inserted inside a footnote,
	// header or footer. Because section breaks can only be inserted inside the
	// body, the segment ID field must be empty.
	EndOfSegmentLocation *EndOfSegmentLocation `json:"endOfSegmentLocation,omitempty"`
	// Location: Inserts a newline and a section break at a specific index in the
	// document. The section break must be inserted inside the bounds of an
	// existing Paragraph. For instance, it cannot be inserted at a table's start
	// index (i.e. between the table and its preceding paragraph). Section breaks
	// cannot be inserted inside a table, equation, footnote, header, or footer.
	// Since section breaks can only be inserted inside the body, the segment ID
	// field must be empty.
	Location *Location `json:"location,omitempty"`
	// SectionType: The type of section to insert.
	//
	// Possible values:
	//   "SECTION_TYPE_UNSPECIFIED" - The section type is unspecified.
	//   "CONTINUOUS" - The section starts immediately after the last paragraph of
	// the previous section.
	//   "NEXT_PAGE" - The section starts on the next page.
	SectionType string `json:"sectionType,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EndOfSegmentLocation") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EndOfSegmentLocation") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertSectionBreakRequest: Inserts a section break at the given location. A newline character will be inserted before the section break.

func (InsertSectionBreakRequest) MarshalJSON ¶
added in v0.15.0
func (s InsertSectionBreakRequest) MarshalJSON() ([]byte, error)
type InsertTableColumnRequest ¶
added in v0.5.0
type InsertTableColumnRequest struct {
	// InsertRight: Whether to insert new column to the right of the reference cell
	// location. - `True`: insert to the right. - `False`: insert to the left.
	InsertRight bool `json:"insertRight,omitempty"`
	// TableCellLocation: The reference table cell location from which columns will
	// be inserted. A new column will be inserted to the left (or right) of the
	// column where the reference cell is. If the reference cell is a merged cell,
	// a new column will be inserted to the left (or right) of the merged cell.
	TableCellLocation *TableCellLocation `json:"tableCellLocation,omitempty"`
	// ForceSendFields is a list of field names (e.g. "InsertRight") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "InsertRight") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertTableColumnRequest: Inserts an empty column into a table.

func (InsertTableColumnRequest) MarshalJSON ¶
added in v0.5.0
func (s InsertTableColumnRequest) MarshalJSON() ([]byte, error)
type InsertTableRequest ¶
added in v0.3.0
type InsertTableRequest struct {
	// Columns: The number of columns in the table.
	Columns int64 `json:"columns,omitempty"`
	// EndOfSegmentLocation: Inserts the table at the end of the given header,
	// footer or document body. A newline character will be inserted before the
	// inserted table. Tables cannot be inserted inside a footnote.
	EndOfSegmentLocation *EndOfSegmentLocation `json:"endOfSegmentLocation,omitempty"`
	// Location: Inserts the table at a specific model index. A newline character
	// will be inserted before the inserted table, therefore the table start index
	// will be at the specified location index + 1. The table must be inserted
	// inside the bounds of an existing Paragraph. For instance, it cannot be
	// inserted at a table's start index (i.e. between an existing table and its
	// preceding paragraph). Tables cannot be inserted inside a footnote or
	// equation.
	Location *Location `json:"location,omitempty"`
	// Rows: The number of rows in the table.
	Rows int64 `json:"rows,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Columns") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Columns") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertTableRequest: Inserts a table at the specified location. A newline character will be inserted before the inserted table.

func (InsertTableRequest) MarshalJSON ¶
added in v0.3.0
func (s InsertTableRequest) MarshalJSON() ([]byte, error)
type InsertTableRowRequest ¶
type InsertTableRowRequest struct {
	// InsertBelow: Whether to insert new row below the reference cell location. -
	// `True`: insert below the cell. - `False`: insert above the cell.
	InsertBelow bool `json:"insertBelow,omitempty"`
	// TableCellLocation: The reference table cell location from which rows will be
	// inserted. A new row will be inserted above (or below) the row where the
	// reference cell is. If the reference cell is a merged cell, a new row will be
	// inserted above (or below) the merged cell.
	TableCellLocation *TableCellLocation `json:"tableCellLocation,omitempty"`
	// ForceSendFields is a list of field names (e.g. "InsertBelow") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "InsertBelow") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertTableRowRequest: Inserts an empty row into a table.

func (InsertTableRowRequest) MarshalJSON ¶
func (s InsertTableRowRequest) MarshalJSON() ([]byte, error)
type InsertTextRequest ¶
type InsertTextRequest struct {
	// EndOfSegmentLocation: Inserts the text at the end of a header, footer,
	// footnote or the document body.
	EndOfSegmentLocation *EndOfSegmentLocation `json:"endOfSegmentLocation,omitempty"`
	// Location: Inserts the text at a specific index in the document. Text must be
	// inserted inside the bounds of an existing Paragraph. For instance, text
	// cannot be inserted at a table's start index (i.e. between the table and its
	// preceding paragraph). The text must be inserted in the preceding paragraph.
	Location *Location `json:"location,omitempty"`
	// Text: The text to be inserted. Inserting a newline character will implicitly
	// create a new Paragraph at that index. The paragraph style of the new
	// paragraph will be copied from the paragraph at the current insertion index,
	// including lists and bullets. Text styles for inserted text will be
	// determined automatically, generally preserving the styling of neighboring
	// text. In most cases, the text style for the inserted text will match the
	// text immediately before the insertion index. Some control characters
	// (U+0000-U+0008, U+000C-U+001F) and characters from the Unicode Basic
	// Multilingual Plane Private Use Area (U+E000-U+F8FF) will be stripped out of
	// the inserted text.
	Text string `json:"text,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EndOfSegmentLocation") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EndOfSegmentLocation") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
InsertTextRequest: Inserts text at the specified location.

func (InsertTextRequest) MarshalJSON ¶
func (s InsertTextRequest) MarshalJSON() ([]byte, error)
type Link ¶
type Link struct {
	// Bookmark: A bookmark in this document. In documents containing a single tab,
	// links to bookmarks within the singular tab continue to return
	// Link.bookmarkId when the includeTabsContent parameter is set to `false` or
	// unset. Otherwise, this field is returned.
	Bookmark *BookmarkLink `json:"bookmark,omitempty"`
	// BookmarkId: The ID of a bookmark in this document. Legacy field: Instead,
	// set includeTabsContent to `true` and use Link.bookmark for read and write
	// operations. This field is only returned when includeTabsContent is set to
	// `false` in documents containing a single tab and links to a bookmark within
	// the singular tab. Otherwise, Link.bookmark is returned. If this field is
	// used in a write request, the bookmark is considered to be from the tab ID
	// specified in the request. If a tab ID is not specified in the request, it is
	// considered to be from the first tab in the document.
	BookmarkId string `json:"bookmarkId,omitempty"`
	// Heading: A heading in this document. In documents containing a single tab,
	// links to headings within the singular tab continue to return Link.headingId
	// when the includeTabsContent parameter is set to `false` or unset. Otherwise,
	// this field is returned.
	Heading *HeadingLink `json:"heading,omitempty"`
	// HeadingId: The ID of a heading in this document. Legacy field: Instead, set
	// includeTabsContent to `true` and use Link.heading for read and write
	// operations. This field is only returned when includeTabsContent is set to
	// `false` in documents containing a single tab and links to a heading within
	// the singular tab. Otherwise, Link.heading is returned. If this field is used
	// in a write request, the heading is considered to be from the tab ID
	// specified in the request. If a tab ID is not specified in the request, it is
	// considered to be from the first tab in the document.
	HeadingId string `json:"headingId,omitempty"`
	// TabId: The ID of a tab in this document.
	TabId string `json:"tabId,omitempty"`
	// Url: An external URL.
	Url string `json:"url,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Bookmark") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Bookmark") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Link: A reference to another portion of a document or an external URL resource.

func (Link) MarshalJSON ¶
func (s Link) MarshalJSON() ([]byte, error)
type LinkedContentReference ¶
type LinkedContentReference struct {
	// SheetsChartReference: A reference to the linked chart.
	SheetsChartReference *SheetsChartReference `json:"sheetsChartReference,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SheetsChartReference") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SheetsChartReference") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
LinkedContentReference: A reference to the external linked source content.

func (LinkedContentReference) MarshalJSON ¶
func (s LinkedContentReference) MarshalJSON() ([]byte, error)
type LinkedContentReferenceSuggestionState ¶
type LinkedContentReferenceSuggestionState struct {
	// SheetsChartReferenceSuggestionState: A mask that indicates which of the
	// fields in sheets_chart_reference have been changed in this suggestion.
	SheetsChartReferenceSuggestionState *SheetsChartReferenceSuggestionState `json:"sheetsChartReferenceSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g.
	// "SheetsChartReferenceSuggestionState") to unconditionally include in API
	// requests. By default, fields with empty or default values are omitted from
	// API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g.
	// "SheetsChartReferenceSuggestionState") to include in API requests with the
	// JSON null value. By default, fields with empty values are omitted from API
	// requests. See https://pkg.go.dev/google.golang.org/api#hdr-NullFields for
	// more details.
	NullFields []string `json:"-"`
}
LinkedContentReferenceSuggestionState: A mask that indicates which of the fields on the base LinkedContentReference have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (LinkedContentReferenceSuggestionState) MarshalJSON ¶
func (s LinkedContentReferenceSuggestionState) MarshalJSON() ([]byte, error)
type List ¶
type List struct {
	// ListProperties: The properties of the list.
	ListProperties *ListProperties `json:"listProperties,omitempty"`
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this list.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionId: The suggested insertion ID. If empty, then this is not
	// a suggested insertion.
	SuggestedInsertionId string `json:"suggestedInsertionId,omitempty"`
	// SuggestedListPropertiesChanges: The suggested changes to the list
	// properties, keyed by suggestion ID.
	SuggestedListPropertiesChanges map[string]SuggestedListProperties `json:"suggestedListPropertiesChanges,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ListProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ListProperties") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
List: A List represents the list attributes for a group of paragraphs that all belong to the same list. A paragraph that's part of a list has a reference to the list's ID in its bullet.

func (List) MarshalJSON ¶
func (s List) MarshalJSON() ([]byte, error)
type ListProperties ¶
type ListProperties struct {
	// NestingLevels: Describes the properties of the bullets at the associated
	// level. A list has at most 9 levels of nesting with nesting level 0
	// corresponding to the top-most level and nesting level 8 corresponding to the
	// most nested level. The nesting levels are returned in ascending order with
	// the least nested returned first.
	NestingLevels []*NestingLevel `json:"nestingLevels,omitempty"`
	// ForceSendFields is a list of field names (e.g. "NestingLevels") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "NestingLevels") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ListProperties: The properties of a list that describe the look and feel of bullets belonging to paragraphs associated with a list.

func (ListProperties) MarshalJSON ¶
func (s ListProperties) MarshalJSON() ([]byte, error)
type ListPropertiesSuggestionState ¶
type ListPropertiesSuggestionState struct {
	// NestingLevelsSuggestionStates: A mask that indicates which of the fields on
	// the corresponding NestingLevel in nesting_levels have been changed in this
	// suggestion. The nesting level suggestion states are returned in ascending
	// order of the nesting level with the least nested returned first.
	NestingLevelsSuggestionStates []*NestingLevelSuggestionState `json:"nestingLevelsSuggestionStates,omitempty"`
	// ForceSendFields is a list of field names (e.g.
	// "NestingLevelsSuggestionStates") to unconditionally include in API requests.
	// By default, fields with empty or default values are omitted from API
	// requests. See https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields
	// for more details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "NestingLevelsSuggestionStates")
	// to include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ListPropertiesSuggestionState: A mask that indicates which of the fields on the base ListProperties have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (ListPropertiesSuggestionState) MarshalJSON ¶
func (s ListPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type Location ¶
type Location struct {
	// Index: The zero-based index, in UTF-16 code units. The index is relative to
	// the beginning of the segment specified by segment_id.
	Index int64 `json:"index,omitempty"`
	// SegmentId: The ID of the header, footer or footnote the location is in. An
	// empty segment ID signifies the document's body.
	SegmentId string `json:"segmentId,omitempty"`
	// TabId: The tab that the location is in. When omitted, the request is applied
	// to the first tab. In a document containing a single tab: - If provided, must
	// match the singular tab's ID. - If omitted, the request applies to the
	// singular tab. In a document containing multiple tabs: - If provided, the
	// request applies to the specified tab. - If omitted, the request applies to
	// the first tab in the document.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Index") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Index") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Location: A particular location in the document.

func (Location) MarshalJSON ¶
func (s Location) MarshalJSON() ([]byte, error)
type MergeTableCellsRequest ¶
added in v0.10.0
type MergeTableCellsRequest struct {
	// TableRange: The table range specifying which cells of the table to merge.
	// Any text in the cells being merged will be concatenated and stored in the
	// "head" cell of the range. This is the upper-left cell of the range when the
	// content direction is left to right, and the upper-right cell of the range
	// otherwise. If the range is non-rectangular (which can occur in some cases
	// where the range covers cells that are already merged or where the table is
	// non-rectangular), a 400 bad request error is returned.
	TableRange *TableRange `json:"tableRange,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TableRange") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TableRange") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
MergeTableCellsRequest: Merges cells in a Table.

func (MergeTableCellsRequest) MarshalJSON ¶
added in v0.10.0
func (s MergeTableCellsRequest) MarshalJSON() ([]byte, error)
type NamedRange ¶
type NamedRange struct {
	// Name: The name of the named range.
	Name string `json:"name,omitempty"`
	// NamedRangeId: The ID of the named range.
	NamedRangeId string `json:"namedRangeId,omitempty"`
	// Ranges: The ranges that belong to this named range.
	Ranges []*Range `json:"ranges,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Name") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Name") to include in API requests
	// with the JSON null value. By default, fields with empty values are omitted
	// from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
NamedRange: A collection of Ranges with the same named range ID. Named ranges allow developers to associate parts of a document with an arbitrary user-defined label so their contents can be programmatically read or edited later. A document can contain multiple named ranges with the same name, but every named range has a unique ID. A named range is created with a single Range, and content inserted inside a named range generally expands that range. However, certain document changes can cause the range to be split into multiple ranges. Named ranges are not private. All applications and collaborators that have access to the document can see its named ranges.

func (NamedRange) MarshalJSON ¶
func (s NamedRange) MarshalJSON() ([]byte, error)
type NamedRanges ¶
type NamedRanges struct {
	// Name: The name that all the named ranges share.
	Name string `json:"name,omitempty"`
	// NamedRanges: The NamedRanges that share the same name.
	NamedRanges []*NamedRange `json:"namedRanges,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Name") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Name") to include in API requests
	// with the JSON null value. By default, fields with empty values are omitted
	// from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
NamedRanges: A collection of all the NamedRanges in the document that share a given name.

func (NamedRanges) MarshalJSON ¶
func (s NamedRanges) MarshalJSON() ([]byte, error)
type NamedStyle ¶
type NamedStyle struct {
	// NamedStyleType: The type of this named style.
	//
	// Possible values:
	//   "NAMED_STYLE_TYPE_UNSPECIFIED" - The type of named style is unspecified.
	//   "NORMAL_TEXT" - Normal text.
	//   "TITLE" - Title.
	//   "SUBTITLE" - Subtitle.
	//   "HEADING_1" - Heading 1.
	//   "HEADING_2" - Heading 2.
	//   "HEADING_3" - Heading 3.
	//   "HEADING_4" - Heading 4.
	//   "HEADING_5" - Heading 5.
	//   "HEADING_6" - Heading 6.
	NamedStyleType string `json:"namedStyleType,omitempty"`
	// ParagraphStyle: The paragraph style of this named style.
	ParagraphStyle *ParagraphStyle `json:"paragraphStyle,omitempty"`
	// TextStyle: The text style of this named style.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "NamedStyleType") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "NamedStyleType") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
NamedStyle: A named style. Paragraphs in the document can inherit their TextStyle and ParagraphStyle from this named style when they have the same named style type.

func (NamedStyle) MarshalJSON ¶
func (s NamedStyle) MarshalJSON() ([]byte, error)
type NamedStyleSuggestionState ¶
type NamedStyleSuggestionState struct {
	// NamedStyleType: The named style type that this suggestion state corresponds
	// to. This field is provided as a convenience for matching the
	// NamedStyleSuggestionState with its corresponding NamedStyle.
	//
	// Possible values:
	//   "NAMED_STYLE_TYPE_UNSPECIFIED" - The type of named style is unspecified.
	//   "NORMAL_TEXT" - Normal text.
	//   "TITLE" - Title.
	//   "SUBTITLE" - Subtitle.
	//   "HEADING_1" - Heading 1.
	//   "HEADING_2" - Heading 2.
	//   "HEADING_3" - Heading 3.
	//   "HEADING_4" - Heading 4.
	//   "HEADING_5" - Heading 5.
	//   "HEADING_6" - Heading 6.
	NamedStyleType string `json:"namedStyleType,omitempty"`
	// ParagraphStyleSuggestionState: A mask that indicates which of the fields in
	// paragraph style have been changed in this suggestion.
	ParagraphStyleSuggestionState *ParagraphStyleSuggestionState `json:"paragraphStyleSuggestionState,omitempty"`
	// TextStyleSuggestionState: A mask that indicates which of the fields in text
	// style have been changed in this suggestion.
	TextStyleSuggestionState *TextStyleSuggestionState `json:"textStyleSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "NamedStyleType") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "NamedStyleType") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
NamedStyleSuggestionState: A suggestion state of a NamedStyle message.

func (NamedStyleSuggestionState) MarshalJSON ¶
func (s NamedStyleSuggestionState) MarshalJSON() ([]byte, error)
type NamedStyles ¶
type NamedStyles struct {
	// Styles: The named styles. There's an entry for each of the possible named
	// style types.
	Styles []*NamedStyle `json:"styles,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Styles") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Styles") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
NamedStyles: The named styles. Paragraphs in the document can inherit their TextStyle and ParagraphStyle from these named styles.

func (NamedStyles) MarshalJSON ¶
func (s NamedStyles) MarshalJSON() ([]byte, error)
type NamedStylesSuggestionState ¶
type NamedStylesSuggestionState struct {
	// StylesSuggestionStates: A mask that indicates which of the fields on the
	// corresponding NamedStyle in styles have been changed in this suggestion. The
	// order of these named style suggestion states matches the order of the
	// corresponding named style within the named styles suggestion.
	StylesSuggestionStates []*NamedStyleSuggestionState `json:"stylesSuggestionStates,omitempty"`
	// ForceSendFields is a list of field names (e.g. "StylesSuggestionStates") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "StylesSuggestionStates") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
NamedStylesSuggestionState: The suggestion state of a NamedStyles message.

func (NamedStylesSuggestionState) MarshalJSON ¶
func (s NamedStylesSuggestionState) MarshalJSON() ([]byte, error)
type NestingLevel ¶
type NestingLevel struct {
	// BulletAlignment: The alignment of the bullet within the space allotted for
	// rendering the bullet.
	//
	// Possible values:
	//   "BULLET_ALIGNMENT_UNSPECIFIED" - The bullet alignment is unspecified.
	//   "START" - The bullet is aligned to the start of the space allotted for
	// rendering the bullet. Left-aligned for LTR text, right-aligned otherwise.
	//   "CENTER" - The bullet is aligned to the center of the space allotted for
	// rendering the bullet.
	//   "END" - The bullet is aligned to the end of the space allotted for
	// rendering the bullet. Right-aligned for LTR text, left-aligned otherwise.
	BulletAlignment string `json:"bulletAlignment,omitempty"`
	// GlyphFormat: The format string used by bullets at this level of nesting. The
	// glyph format contains one or more placeholders, and these placeholders are
	// replaced with the appropriate values depending on the glyph_type or
	// glyph_symbol. The placeholders follow the pattern `%[nesting_level]`.
	// Furthermore, placeholders can have prefixes and suffixes. Thus, the glyph
	// format follows the pattern `%[nesting_level]`. Note that the prefix and
	// suffix are optional and can be arbitrary strings. For example, the glyph
	// format `%0.` indicates that the rendered glyph will replace the placeholder
	// with the corresponding glyph for nesting level 0 followed by a period as the
	// suffix. So a list with a glyph type of UPPER_ALPHA and glyph format `%0.` at
	// nesting level 0 will result in a list with rendered glyphs `A.` `B.` `C.`
	// The glyph format can contain placeholders for the current nesting level as
	// well as placeholders for parent nesting levels. For example, a list can have
	// a glyph format of `%0.` at nesting level 0 and a glyph format of `%0.%1.` at
	// nesting level 1. Assuming both nesting levels have DECIMAL glyph types, this
	// would result in a list with rendered glyphs `1.` `2.` ` 2.1.` ` 2.2.` `3.`
	// For nesting levels that are ordered, the string that replaces a placeholder
	// in the glyph format for a particular paragraph depends on the paragraph's
	// order within the list.
	GlyphFormat string `json:"glyphFormat,omitempty"`
	// GlyphSymbol: A custom glyph symbol used by bullets when paragraphs at this
	// level of nesting is unordered. The glyph symbol replaces placeholders within
	// the glyph_format. For example, if the glyph_symbol is the solid circle
	// corresponding to Unicode U+25cf code point and the glyph_format is `%0`, the
	// rendered glyph would be the solid circle.
	GlyphSymbol string `json:"glyphSymbol,omitempty"`
	// GlyphType: The type of glyph used by bullets when paragraphs at this level
	// of nesting is ordered. The glyph type determines the type of glyph used to
	// replace placeholders within the glyph_format when paragraphs at this level
	// of nesting are ordered. For example, if the nesting level is 0, the
	// glyph_format is `%0.` and the glyph type is DECIMAL, then the rendered glyph
	// would replace the placeholder `%0` in the glyph format with a number
	// corresponding to the list item's order within the list.
	//
	// Possible values:
	//   "GLYPH_TYPE_UNSPECIFIED" - The glyph type is unspecified or unsupported.
	//   "NONE" - An empty string.
	//   "DECIMAL" - A number, like `1`, `2`, or `3`.
	//   "ZERO_DECIMAL" - A number where single digit numbers are prefixed with a
	// zero, like `01`, `02`, or `03`. Numbers with more than one digit are not
	// prefixed with a zero.
	//   "UPPER_ALPHA" - An uppercase letter, like `A`, `B`, or `C`.
	//   "ALPHA" - A lowercase letter, like `a`, `b`, or `c`.
	//   "UPPER_ROMAN" - An uppercase Roman numeral, like `I`, `II`, or `III`.
	//   "ROMAN" - A lowercase Roman numeral, like `i`, `ii`, or `iii`.
	GlyphType string `json:"glyphType,omitempty"`
	// IndentFirstLine: The amount of indentation for the first line of paragraphs
	// at this level of nesting.
	IndentFirstLine *Dimension `json:"indentFirstLine,omitempty"`
	// IndentStart: The amount of indentation for paragraphs at this level of
	// nesting. Applied to the side that corresponds to the start of the text,
	// based on the paragraph's content direction.
	IndentStart *Dimension `json:"indentStart,omitempty"`
	// StartNumber: The number of the first list item at this nesting level. A
	// value of 0 is treated as a value of 1 for lettered lists and Roman numeral
	// lists. For values of both 0 and 1, lettered and Roman numeral lists will
	// begin at `a` and `i` respectively. This value is ignored for nesting levels
	// with unordered glyphs.
	StartNumber int64 `json:"startNumber,omitempty"`
	// TextStyle: The text style of bullets at this level of nesting.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "BulletAlignment") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "BulletAlignment") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
NestingLevel: Contains properties describing the look and feel of a list bullet at a given level of nesting.

func (NestingLevel) MarshalJSON ¶
func (s NestingLevel) MarshalJSON() ([]byte, error)
type NestingLevelSuggestionState ¶
type NestingLevelSuggestionState struct {
	// BulletAlignmentSuggested: Indicates if there was a suggested change to
	// bullet_alignment.
	BulletAlignmentSuggested bool `json:"bulletAlignmentSuggested,omitempty"`
	// GlyphFormatSuggested: Indicates if there was a suggested change to
	// glyph_format.
	GlyphFormatSuggested bool `json:"glyphFormatSuggested,omitempty"`
	// GlyphSymbolSuggested: Indicates if there was a suggested change to
	// glyph_symbol.
	GlyphSymbolSuggested bool `json:"glyphSymbolSuggested,omitempty"`
	// GlyphTypeSuggested: Indicates if there was a suggested change to glyph_type.
	GlyphTypeSuggested bool `json:"glyphTypeSuggested,omitempty"`
	// IndentFirstLineSuggested: Indicates if there was a suggested change to
	// indent_first_line.
	IndentFirstLineSuggested bool `json:"indentFirstLineSuggested,omitempty"`
	// IndentStartSuggested: Indicates if there was a suggested change to
	// indent_start.
	IndentStartSuggested bool `json:"indentStartSuggested,omitempty"`
	// StartNumberSuggested: Indicates if there was a suggested change to
	// start_number.
	StartNumberSuggested bool `json:"startNumberSuggested,omitempty"`
	// TextStyleSuggestionState: A mask that indicates which of the fields in text
	// style have been changed in this suggestion.
	TextStyleSuggestionState *TextStyleSuggestionState `json:"textStyleSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "BulletAlignmentSuggested")
	// to unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "BulletAlignmentSuggested") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
NestingLevelSuggestionState: A mask that indicates which of the fields on the base NestingLevel have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (NestingLevelSuggestionState) MarshalJSON ¶
func (s NestingLevelSuggestionState) MarshalJSON() ([]byte, error)
type ObjectReferences ¶
type ObjectReferences struct {
	// ObjectIds: The object IDs.
	ObjectIds []string `json:"objectIds,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ObjectIds") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ObjectIds") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ObjectReferences: A collection of object IDs.

func (ObjectReferences) MarshalJSON ¶
func (s ObjectReferences) MarshalJSON() ([]byte, error)
type OptionalColor ¶
type OptionalColor struct {
	// Color: If set, this will be used as an opaque color. If unset, this
	// represents a transparent color.
	Color *Color `json:"color,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Color") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Color") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
OptionalColor: A color that can either be fully opaque or fully transparent.

func (OptionalColor) MarshalJSON ¶
func (s OptionalColor) MarshalJSON() ([]byte, error)
type PageBreak ¶
type PageBreak struct {
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. A PageBreak may have
	// multiple insertion IDs if it's a nested suggested change. If empty, then
	// this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTextStyleChanges: The suggested text style changes to this
	// PageBreak, keyed by suggestion ID.
	SuggestedTextStyleChanges map[string]SuggestedTextStyle `json:"suggestedTextStyleChanges,omitempty"`
	// TextStyle: The text style of this PageBreak. Similar to text content, like
	// text runs and footnote references, the text style of a page break can affect
	// content layout as well as the styling of text inserted next to it.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SuggestedDeletionIds") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SuggestedDeletionIds") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
PageBreak: A ParagraphElement representing a page break. A page break makes the subsequent text start at the top of the next page.

func (PageBreak) MarshalJSON ¶
func (s PageBreak) MarshalJSON() ([]byte, error)
type Paragraph ¶
type Paragraph struct {
	// Bullet: The bullet for this paragraph. If not present, the paragraph does
	// not belong to a list.
	Bullet *Bullet `json:"bullet,omitempty"`
	// Elements: The content of the paragraph, broken down into its component
	// parts.
	Elements []*ParagraphElement `json:"elements,omitempty"`
	// ParagraphStyle: The style of this paragraph.
	ParagraphStyle *ParagraphStyle `json:"paragraphStyle,omitempty"`
	// PositionedObjectIds: The IDs of the positioned objects tethered to this
	// paragraph.
	PositionedObjectIds []string `json:"positionedObjectIds,omitempty"`
	// SuggestedBulletChanges: The suggested changes to this paragraph's bullet.
	SuggestedBulletChanges map[string]SuggestedBullet `json:"suggestedBulletChanges,omitempty"`
	// SuggestedParagraphStyleChanges: The suggested paragraph style changes to
	// this paragraph, keyed by suggestion ID.
	SuggestedParagraphStyleChanges map[string]SuggestedParagraphStyle `json:"suggestedParagraphStyleChanges,omitempty"`
	// SuggestedPositionedObjectIds: The IDs of the positioned objects suggested to
	// be attached to this paragraph, keyed by suggestion ID.
	SuggestedPositionedObjectIds map[string]ObjectReferences `json:"suggestedPositionedObjectIds,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Bullet") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Bullet") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Paragraph: A StructuralElement representing a paragraph. A paragraph is a range of content that's terminated with a newline character.

func (Paragraph) MarshalJSON ¶
func (s Paragraph) MarshalJSON() ([]byte, error)
type ParagraphBorder ¶
type ParagraphBorder struct {
	// Color: The color of the border.
	Color *OptionalColor `json:"color,omitempty"`
	// DashStyle: The dash style of the border.
	//
	// Possible values:
	//   "DASH_STYLE_UNSPECIFIED" - Unspecified dash style.
	//   "SOLID" - Solid line. Corresponds to ECMA-376 ST_PresetLineDashVal value
	// 'solid'. This is the default dash style.
	//   "DOT" - Dotted line. Corresponds to ECMA-376 ST_PresetLineDashVal value
	// 'dot'.
	//   "DASH" - Dashed line. Corresponds to ECMA-376 ST_PresetLineDashVal value
	// 'dash'.
	DashStyle string `json:"dashStyle,omitempty"`
	// Padding: The padding of the border.
	Padding *Dimension `json:"padding,omitempty"`
	// Width: The width of the border.
	Width *Dimension `json:"width,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Color") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Color") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ParagraphBorder: A border around a paragraph.

func (ParagraphBorder) MarshalJSON ¶
func (s ParagraphBorder) MarshalJSON() ([]byte, error)
type ParagraphElement ¶
type ParagraphElement struct {
	// AutoText: An auto text paragraph element.
	AutoText *AutoText `json:"autoText,omitempty"`
	// ColumnBreak: A column break paragraph element.
	ColumnBreak *ColumnBreak `json:"columnBreak,omitempty"`
	// DateElement: A paragraph element that represents a date.
	DateElement *DateElement `json:"dateElement,omitempty"`
	// Dropdown: A paragraph element that represents a dropdown menu.
	Dropdown *Dropdown `json:"dropdown,omitempty"`
	// EndIndex: The zero-base end index of this paragraph element, exclusive, in
	// UTF-16 code units.
	EndIndex int64 `json:"endIndex,omitempty"`
	// Equation: An equation paragraph element.
	Equation *Equation `json:"equation,omitempty"`
	// FootnoteReference: A footnote reference paragraph element.
	FootnoteReference *FootnoteReference `json:"footnoteReference,omitempty"`
	// HorizontalRule: A horizontal rule paragraph element.
	HorizontalRule *HorizontalRule `json:"horizontalRule,omitempty"`
	// InlineObjectElement: An inline object paragraph element.
	InlineObjectElement *InlineObjectElement `json:"inlineObjectElement,omitempty"`
	// PageBreak: A page break paragraph element.
	PageBreak *PageBreak `json:"pageBreak,omitempty"`
	// Person: A paragraph element that links to a person or email address.
	Person *Person `json:"person,omitempty"`
	// RichLink: A paragraph element that links to a Google resource (such as a
	// file in Google Drive, a YouTube video, or a Calendar event.)
	RichLink *RichLink `json:"richLink,omitempty"`
	// StartIndex: The zero-based start index of this paragraph element, in UTF-16
	// code units.
	StartIndex int64 `json:"startIndex,omitempty"`
	// TextRun: A text run paragraph element.
	TextRun *TextRun `json:"textRun,omitempty"`
	// ForceSendFields is a list of field names (e.g. "AutoText") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "AutoText") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ParagraphElement: A ParagraphElement describes content within a Paragraph.

func (ParagraphElement) MarshalJSON ¶
func (s ParagraphElement) MarshalJSON() ([]byte, error)
type ParagraphStyle ¶
type ParagraphStyle struct {
	// Alignment: The text alignment for this paragraph.
	//
	// Possible values:
	//   "ALIGNMENT_UNSPECIFIED" - The paragraph alignment is inherited from the
	// parent.
	//   "START" - The paragraph is aligned to the start of the line. Left-aligned
	// for LTR text, right-aligned otherwise.
	//   "CENTER" - The paragraph is centered.
	//   "END" - The paragraph is aligned to the end of the line. Right-aligned for
	// LTR text, left-aligned otherwise.
	//   "JUSTIFIED" - The paragraph is justified.
	Alignment string `json:"alignment,omitempty"`
	// AvoidWidowAndOrphan: Whether to avoid widows and orphans for the paragraph.
	// If unset, the value is inherited from the parent.
	AvoidWidowAndOrphan bool `json:"avoidWidowAndOrphan,omitempty"`
	// BorderBetween: The border between this paragraph and the next and previous
	// paragraphs. If unset, the value is inherited from the parent. The between
	// border is rendered when the adjacent paragraph has the same border and
	// indent properties. Paragraph borders cannot be partially updated. When
	// changing a paragraph border, the new border must be specified in its
	// entirety.
	BorderBetween *ParagraphBorder `json:"borderBetween,omitempty"`
	// BorderBottom: The border at the bottom of this paragraph. If unset, the
	// value is inherited from the parent. The bottom border is rendered when the
	// paragraph below has different border and indent properties. Paragraph
	// borders cannot be partially updated. When changing a paragraph border, the
	// new border must be specified in its entirety.
	BorderBottom *ParagraphBorder `json:"borderBottom,omitempty"`
	// BorderLeft: The border to the left of this paragraph. If unset, the value is
	// inherited from the parent. Paragraph borders cannot be partially updated.
	// When changing a paragraph border, the new border must be specified in its
	// entirety.
	BorderLeft *ParagraphBorder `json:"borderLeft,omitempty"`
	// BorderRight: The border to the right of this paragraph. If unset, the value
	// is inherited from the parent. Paragraph borders cannot be partially updated.
	// When changing a paragraph border, the new border must be specified in its
	// entirety.
	BorderRight *ParagraphBorder `json:"borderRight,omitempty"`
	// BorderTop: The border at the top of this paragraph. If unset, the value is
	// inherited from the parent. The top border is rendered when the paragraph
	// above has different border and indent properties. Paragraph borders cannot
	// be partially updated. When changing a paragraph border, the new border must
	// be specified in its entirety.
	BorderTop *ParagraphBorder `json:"borderTop,omitempty"`
	// Direction: The text direction of this paragraph. If unset, the value
	// defaults to LEFT_TO_RIGHT since paragraph direction is not inherited.
	//
	// Possible values:
	//   "CONTENT_DIRECTION_UNSPECIFIED" - The content direction is unspecified.
	//   "LEFT_TO_RIGHT" - The content goes from left to right.
	//   "RIGHT_TO_LEFT" - The content goes from right to left.
	Direction string `json:"direction,omitempty"`
	// HeadingId: The heading ID of the paragraph. If empty, then this paragraph is
	// not a heading. This property is read-only.
	HeadingId string `json:"headingId,omitempty"`
	// IndentEnd: The amount of indentation for the paragraph on the side that
	// corresponds to the end of the text, based on the current paragraph
	// direction. If unset, the value is inherited from the parent.
	IndentEnd *Dimension `json:"indentEnd,omitempty"`
	// IndentFirstLine: The amount of indentation for the first line of the
	// paragraph. If unset, the value is inherited from the parent.
	IndentFirstLine *Dimension `json:"indentFirstLine,omitempty"`
	// IndentStart: The amount of indentation for the paragraph on the side that
	// corresponds to the start of the text, based on the current paragraph
	// direction. If unset, the value is inherited from the parent.
	IndentStart *Dimension `json:"indentStart,omitempty"`
	// KeepLinesTogether: Whether all lines of the paragraph should be laid out on
	// the same page or column if possible. If unset, the value is inherited from
	// the parent.
	KeepLinesTogether bool `json:"keepLinesTogether,omitempty"`
	// KeepWithNext: Whether at least a part of this paragraph should be laid out
	// on the same page or column as the next paragraph if possible. If unset, the
	// value is inherited from the parent.
	KeepWithNext bool `json:"keepWithNext,omitempty"`
	// LineSpacing: The amount of space between lines, as a percentage of normal,
	// where normal is represented as 100.0. If unset, the value is inherited from
	// the parent.
	LineSpacing float64 `json:"lineSpacing,omitempty"`
	// NamedStyleType: The named style type of the paragraph. Since updating the
	// named style type affects other properties within ParagraphStyle, the named
	// style type is applied before the other properties are updated.
	//
	// Possible values:
	//   "NAMED_STYLE_TYPE_UNSPECIFIED" - The type of named style is unspecified.
	//   "NORMAL_TEXT" - Normal text.
	//   "TITLE" - Title.
	//   "SUBTITLE" - Subtitle.
	//   "HEADING_1" - Heading 1.
	//   "HEADING_2" - Heading 2.
	//   "HEADING_3" - Heading 3.
	//   "HEADING_4" - Heading 4.
	//   "HEADING_5" - Heading 5.
	//   "HEADING_6" - Heading 6.
	NamedStyleType string `json:"namedStyleType,omitempty"`
	// PageBreakBefore: Whether the current paragraph should always start at the
	// beginning of a page. If unset, the value is inherited from the parent.
	// Attempting to update page_break_before for paragraphs in unsupported
	// regions, including Table, Header, Footer and Footnote, can result in an
	// invalid document state that returns a 400 bad request error.
	PageBreakBefore bool `json:"pageBreakBefore,omitempty"`
	// Shading: The shading of the paragraph. If unset, the value is inherited from
	// the parent.
	Shading *Shading `json:"shading,omitempty"`
	// SpaceAbove: The amount of extra space above the paragraph. If unset, the
	// value is inherited from the parent.
	SpaceAbove *Dimension `json:"spaceAbove,omitempty"`
	// SpaceBelow: The amount of extra space below the paragraph. If unset, the
	// value is inherited from the parent.
	SpaceBelow *Dimension `json:"spaceBelow,omitempty"`
	// SpacingMode: The spacing mode for the paragraph.
	//
	// Possible values:
	//   "SPACING_MODE_UNSPECIFIED" - The spacing mode is inherited from the
	// parent.
	//   "NEVER_COLLAPSE" - Paragraph spacing is always rendered.
	//   "COLLAPSE_LISTS" - Paragraph spacing is skipped between list elements.
	SpacingMode string `json:"spacingMode,omitempty"`
	// TabStops: A list of the tab stops for this paragraph. The list of tab stops
	// is not inherited. This property is read-only.
	TabStops []*TabStop `json:"tabStops,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Alignment") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Alignment") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ParagraphStyle: Styles that apply to a whole paragraph. Inherited paragraph styles are represented as unset fields in this message. A paragraph style's parent depends on where the paragraph style is defined: * The ParagraphStyle on a Paragraph inherits from the paragraph's corresponding named style type. * The ParagraphStyle on a named style inherits from the normal text named style. * The ParagraphStyle of the normal text named style inherits from the default paragraph style in the Docs editor. * The ParagraphStyle on a Paragraph element that's contained in a table may inherit its paragraph style from the table style. If the paragraph style does not inherit from a parent, unsetting fields will revert the style to a value matching the defaults in the Docs editor.

func (ParagraphStyle) MarshalJSON ¶
func (s ParagraphStyle) MarshalJSON() ([]byte, error)
func (*ParagraphStyle) UnmarshalJSON ¶
func (s *ParagraphStyle) UnmarshalJSON(data []byte) error
type ParagraphStyleSuggestionState ¶
type ParagraphStyleSuggestionState struct {
	// AlignmentSuggested: Indicates if there was a suggested change to alignment.
	AlignmentSuggested bool `json:"alignmentSuggested,omitempty"`
	// AvoidWidowAndOrphanSuggested: Indicates if there was a suggested change to
	// avoid_widow_and_orphan.
	AvoidWidowAndOrphanSuggested bool `json:"avoidWidowAndOrphanSuggested,omitempty"`
	// BorderBetweenSuggested: Indicates if there was a suggested change to
	// border_between.
	BorderBetweenSuggested bool `json:"borderBetweenSuggested,omitempty"`
	// BorderBottomSuggested: Indicates if there was a suggested change to
	// border_bottom.
	BorderBottomSuggested bool `json:"borderBottomSuggested,omitempty"`
	// BorderLeftSuggested: Indicates if there was a suggested change to
	// border_left.
	BorderLeftSuggested bool `json:"borderLeftSuggested,omitempty"`
	// BorderRightSuggested: Indicates if there was a suggested change to
	// border_right.
	BorderRightSuggested bool `json:"borderRightSuggested,omitempty"`
	// BorderTopSuggested: Indicates if there was a suggested change to border_top.
	BorderTopSuggested bool `json:"borderTopSuggested,omitempty"`
	// DirectionSuggested: Indicates if there was a suggested change to direction.
	DirectionSuggested bool `json:"directionSuggested,omitempty"`
	// HeadingIdSuggested: Indicates if there was a suggested change to heading_id.
	HeadingIdSuggested bool `json:"headingIdSuggested,omitempty"`
	// IndentEndSuggested: Indicates if there was a suggested change to indent_end.
	IndentEndSuggested bool `json:"indentEndSuggested,omitempty"`
	// IndentFirstLineSuggested: Indicates if there was a suggested change to
	// indent_first_line.
	IndentFirstLineSuggested bool `json:"indentFirstLineSuggested,omitempty"`
	// IndentStartSuggested: Indicates if there was a suggested change to
	// indent_start.
	IndentStartSuggested bool `json:"indentStartSuggested,omitempty"`
	// KeepLinesTogetherSuggested: Indicates if there was a suggested change to
	// keep_lines_together.
	KeepLinesTogetherSuggested bool `json:"keepLinesTogetherSuggested,omitempty"`
	// KeepWithNextSuggested: Indicates if there was a suggested change to
	// keep_with_next.
	KeepWithNextSuggested bool `json:"keepWithNextSuggested,omitempty"`
	// LineSpacingSuggested: Indicates if there was a suggested change to
	// line_spacing.
	LineSpacingSuggested bool `json:"lineSpacingSuggested,omitempty"`
	// NamedStyleTypeSuggested: Indicates if there was a suggested change to
	// named_style_type.
	NamedStyleTypeSuggested bool `json:"namedStyleTypeSuggested,omitempty"`
	// PageBreakBeforeSuggested: Indicates if there was a suggested change to
	// page_break_before.
	PageBreakBeforeSuggested bool `json:"pageBreakBeforeSuggested,omitempty"`
	// ShadingSuggestionState: A mask that indicates which of the fields in shading
	// have been changed in this suggestion.
	ShadingSuggestionState *ShadingSuggestionState `json:"shadingSuggestionState,omitempty"`
	// SpaceAboveSuggested: Indicates if there was a suggested change to
	// space_above.
	SpaceAboveSuggested bool `json:"spaceAboveSuggested,omitempty"`
	// SpaceBelowSuggested: Indicates if there was a suggested change to
	// space_below.
	SpaceBelowSuggested bool `json:"spaceBelowSuggested,omitempty"`
	// SpacingModeSuggested: Indicates if there was a suggested change to
	// spacing_mode.
	SpacingModeSuggested bool `json:"spacingModeSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "AlignmentSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "AlignmentSuggested") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ParagraphStyleSuggestionState: A mask that indicates which of the fields on the base ParagraphStyle have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (ParagraphStyleSuggestionState) MarshalJSON ¶
func (s ParagraphStyleSuggestionState) MarshalJSON() ([]byte, error)
type Person ¶
added in v0.41.0
type Person struct {
	// PersonId: Output only. The unique ID of this link.
	PersonId string `json:"personId,omitempty"`
	// PersonProperties: Output only. The properties of this Person. This field is
	// always present.
	PersonProperties *PersonProperties `json:"personProperties,omitempty"`
	// SuggestedDeletionIds: IDs for suggestions that remove this person link from
	// the document. A Person might have multiple deletion IDs if, for example,
	// multiple users suggest deleting it. If empty, then this person link isn't
	// suggested for deletion.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: IDs for suggestions that insert this person link into
	// the document. A Person might have multiple insertion IDs if it's a nested
	// suggested change (a suggestion within a suggestion made by a different user,
	// for example). If empty, then this person link isn't a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTextStyleChanges: The suggested text style changes to this Person,
	// keyed by suggestion ID.
	SuggestedTextStyleChanges map[string]SuggestedTextStyle `json:"suggestedTextStyleChanges,omitempty"`
	// TextStyle: The text style of this Person.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "PersonId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "PersonId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Person: A person or email address mentioned in a document. These mentions behave as a single, immutable element containing the person's name or email address.

func (Person) MarshalJSON ¶
added in v0.41.0
func (s Person) MarshalJSON() ([]byte, error)
type PersonProperties ¶
added in v0.41.0
type PersonProperties struct {
	// Email: The email address linked to this Person. This field is always
	// present.
	Email string `json:"email,omitempty"`
	// Name: The name of the person if it's displayed in the link text instead of
	// the person's email address.
	Name string `json:"name,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Email") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Email") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
PersonProperties: Properties specific to a linked Person.

func (PersonProperties) MarshalJSON ¶
added in v0.41.0
func (s PersonProperties) MarshalJSON() ([]byte, error)
type PinTableHeaderRowsRequest ¶
added in v0.81.0
type PinTableHeaderRowsRequest struct {
	// PinnedHeaderRowsCount: The number of table rows to pin, where 0 implies that
	// all rows are unpinned.
	PinnedHeaderRowsCount int64 `json:"pinnedHeaderRowsCount,omitempty"`
	// TableStartLocation: The location where the table starts in the document.
	TableStartLocation *Location `json:"tableStartLocation,omitempty"`
	// ForceSendFields is a list of field names (e.g. "PinnedHeaderRowsCount") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "PinnedHeaderRowsCount") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
PinTableHeaderRowsRequest: Updates the number of pinned table header rows in a table.

func (PinTableHeaderRowsRequest) MarshalJSON ¶
added in v0.81.0
func (s PinTableHeaderRowsRequest) MarshalJSON() ([]byte, error)
type PositionedObject ¶
type PositionedObject struct {
	// ObjectId: The ID of this positioned object.
	ObjectId string `json:"objectId,omitempty"`
	// PositionedObjectProperties: The properties of this positioned object.
	PositionedObjectProperties *PositionedObjectProperties `json:"positionedObjectProperties,omitempty"`
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionId: The suggested insertion ID. If empty, then this is not
	// a suggested insertion.
	SuggestedInsertionId string `json:"suggestedInsertionId,omitempty"`
	// SuggestedPositionedObjectPropertiesChanges: The suggested changes to the
	// positioned object properties, keyed by suggestion ID.
	SuggestedPositionedObjectPropertiesChanges map[string]SuggestedPositionedObjectProperties `json:"suggestedPositionedObjectPropertiesChanges,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ObjectId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ObjectId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
PositionedObject: An object that's tethered to a Paragraph and positioned relative to the beginning of the paragraph. A PositionedObject contains an EmbeddedObject such as an image.

func (PositionedObject) MarshalJSON ¶
func (s PositionedObject) MarshalJSON() ([]byte, error)
type PositionedObjectPositioning ¶
type PositionedObjectPositioning struct {
	// Layout: The layout of this positioned object.
	//
	// Possible values:
	//   "POSITIONED_OBJECT_LAYOUT_UNSPECIFIED" - The layout is unspecified.
	//   "WRAP_TEXT" - The text wraps around the positioned object.
	//   "BREAK_LEFT" - Breaks text such that the positioned object is on the left
	// and text is on the right.
	//   "BREAK_RIGHT" - Breaks text such that the positioned object is on the
	// right and text is on the left.
	//   "BREAK_LEFT_RIGHT" - Breaks text such that there's no text on the left or
	// right of the positioned object.
	//   "IN_FRONT_OF_TEXT" - The positioned object is in front of the text.
	//   "BEHIND_TEXT" - The positioned object is behind the text.
	Layout string `json:"layout,omitempty"`
	// LeftOffset: The offset of the left edge of the positioned object relative to
	// the beginning of the Paragraph it's tethered to. The exact positioning of
	// the object can depend on other content in the document and the document's
	// styling.
	LeftOffset *Dimension `json:"leftOffset,omitempty"`
	// TopOffset: The offset of the top edge of the positioned object relative to
	// the beginning of the Paragraph it's tethered to. The exact positioning of
	// the object can depend on other content in the document and the document's
	// styling.
	TopOffset *Dimension `json:"topOffset,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Layout") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Layout") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
PositionedObjectPositioning: The positioning of a PositionedObject. The positioned object is positioned relative to the beginning of the Paragraph it's tethered to.

func (PositionedObjectPositioning) MarshalJSON ¶
func (s PositionedObjectPositioning) MarshalJSON() ([]byte, error)
type PositionedObjectPositioningSuggestionState ¶
type PositionedObjectPositioningSuggestionState struct {
	// LayoutSuggested: Indicates if there was a suggested change to layout.
	LayoutSuggested bool `json:"layoutSuggested,omitempty"`
	// LeftOffsetSuggested: Indicates if there was a suggested change to
	// left_offset.
	LeftOffsetSuggested bool `json:"leftOffsetSuggested,omitempty"`
	// TopOffsetSuggested: Indicates if there was a suggested change to top_offset.
	TopOffsetSuggested bool `json:"topOffsetSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "LayoutSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "LayoutSuggested") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
PositionedObjectPositioningSuggestionState: A mask that indicates which of the fields on the base PositionedObjectPositioning have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (PositionedObjectPositioningSuggestionState) MarshalJSON ¶
func (s PositionedObjectPositioningSuggestionState) MarshalJSON() ([]byte, error)
type PositionedObjectProperties ¶
type PositionedObjectProperties struct {
	// EmbeddedObject: The embedded object of this positioned object.
	EmbeddedObject *EmbeddedObject `json:"embeddedObject,omitempty"`
	// Positioning: The positioning of this positioned object relative to the
	// newline of the Paragraph that references this positioned object.
	Positioning *PositionedObjectPositioning `json:"positioning,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EmbeddedObject") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EmbeddedObject") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
PositionedObjectProperties: Properties of a PositionedObject.

func (PositionedObjectProperties) MarshalJSON ¶
func (s PositionedObjectProperties) MarshalJSON() ([]byte, error)
type PositionedObjectPropertiesSuggestionState ¶
type PositionedObjectPropertiesSuggestionState struct {
	// EmbeddedObjectSuggestionState: A mask that indicates which of the fields in
	// embedded_object have been changed in this suggestion.
	EmbeddedObjectSuggestionState *EmbeddedObjectSuggestionState `json:"embeddedObjectSuggestionState,omitempty"`
	// PositioningSuggestionState: A mask that indicates which of the fields in
	// positioning have been changed in this suggestion.
	PositioningSuggestionState *PositionedObjectPositioningSuggestionState `json:"positioningSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g.
	// "EmbeddedObjectSuggestionState") to unconditionally include in API requests.
	// By default, fields with empty or default values are omitted from API
	// requests. See https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields
	// for more details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EmbeddedObjectSuggestionState")
	// to include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
PositionedObjectPropertiesSuggestionState: A mask that indicates which of the fields on the base PositionedObjectProperties have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (PositionedObjectPropertiesSuggestionState) MarshalJSON ¶
func (s PositionedObjectPropertiesSuggestionState) MarshalJSON() ([]byte, error)
type Post ¶
type Post struct {
	// AssigneeEmail: Optional. The email of the user who is being newly assigned
	// to the thread as part of this post. Returns a 400 bad request error if: -
	// The parent thread is a CommentThread whose headPost does not have an
	// assignee. - The parent thread is a SuggestionThread. - commentAction is
	// specified as `RESOLVE` or `REOPEN`. - `assigneeEmail` exceeds 2048 UTF-8
	// code units.
	AssigneeEmail string `json:"assigneeEmail,omitempty"`
	// Author: Output only. The user who created the post.
	Author *PostAuthor `json:"author,omitempty"`
	// CommentAction: Optional. The action type for comment posts.
	//
	// Possible values:
	//   "COMMENT_ACTION_TYPE_UNSPECIFIED" - Default value. This value is unused.
	//   "NO_COMMENT_ACTION_CHANGE" - No action change in this post.
	//   "RESOLVE" - This post resolves the thread.
	//   "REOPEN" - This post reopens the thread.
	CommentAction string `json:"commentAction,omitempty"`
	// Content: The content of the post. Required to be non-empty if commentAction
	// is not `RESOLVE` or `REOPEN`. This text content will be handled similarly to
	// comments created in the Docs editor. It will have similar behaviors for
	// formatting, notifications, etc. May not exceed 2048 UTF-8 code units.
	Content string `json:"content,omitempty"`
	// ContentHtml: Output only. The content of the post as HTML.
	ContentHtml string `json:"contentHtml,omitempty"`
	// CreateTime: Output only. The time the post was created.
	CreateTime string `json:"createTime,omitempty"`
	// Deleted: Output only. Whether the post is deleted. If `true`, content and
	// author fields will be empty.
	Deleted bool `json:"deleted,omitempty"`
	// FromCopiedDocument: Output only. Whether the post is from a copied document.
	// This field cannot be set directly by callers.
	FromCopiedDocument bool `json:"fromCopiedDocument,omitempty"`
	// FromDocumentComparison: Output only. Whether the post is from a document
	// comparison. This field cannot be set directly by callers.
	FromDocumentComparison bool `json:"fromDocumentComparison,omitempty"`
	// FromImportedDocument: Output only. Whether the post is from an imported
	// document. This field cannot be set directly by callers.
	FromImportedDocument bool `json:"fromImportedDocument,omitempty"`
	// PostId: Output only. The unique ID of the post.
	PostId string `json:"postId,omitempty"`
	// SuggestionAction: Output only. The action type for suggestion posts.
	//
	// Possible values:
	//   "SUGGESTION_ACTION_TYPE_UNSPECIFIED" - Default value. This value is
	// unused.
	//   "NO_SUGGESTION_ACTION_CHANGE" - No action change in this post.
	//   "ACCEPT" - This post accepts the suggestion.
	//   "REJECT" - This post rejects the suggestion.
	SuggestionAction string `json:"suggestionAction,omitempty"`
	// UpdateTime: Output only. The time the post was last updated.
	UpdateTime string `json:"updateTime,omitempty"`
	// ForceSendFields is a list of field names (e.g. "AssigneeEmail") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "AssigneeEmail") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Post: Represents a single post in a comment or suggestion thread.

func (Post) MarshalJSON ¶
func (s Post) MarshalJSON() ([]byte, error)
type PostAuthor ¶
type PostAuthor struct {
	// Anonymous: Whether the user is anonymous.
	Anonymous bool `json:"anonymous,omitempty"`
	// DisplayName: The display name of the user. May be absent if the author is
	// anonymous.
	DisplayName string `json:"displayName,omitempty"`
	// Me: Whether the user is the authenticated user making the request.
	Me bool `json:"me,omitempty"`
	// User: The resource name of the post author user, which can also be used to
	// identify the user in the Google People API
	// (https://developers.google.com/people/api/rest/v1/people). Format:
	// `users/{user}`. Will not be populated if the anonymous field is `true` or if
	// the post is from an imported document.
	User string `json:"user,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Anonymous") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Anonymous") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
PostAuthor: Represents a user who authored a comment or suggestion post.

func (PostAuthor) MarshalJSON ¶
func (s PostAuthor) MarshalJSON() ([]byte, error)
type Range ¶
type Range struct {
	// EndIndex: The zero-based end index of this range, exclusive, in UTF-16 code
	// units. In all current uses, an end index must be provided. This field is an
	// Int32Value in order to accommodate future use cases with open-ended ranges.
	EndIndex int64 `json:"endIndex,omitempty"`
	// SegmentId: The ID of the header, footer, or footnote that this range is
	// contained in. An empty segment ID signifies the document's body.
	SegmentId string `json:"segmentId,omitempty"`
	// StartIndex: The zero-based start index of this range, in UTF-16 code units.
	// In all current uses, a start index must be provided. This field is an
	// Int32Value in order to accommodate future use cases with open-ended ranges.
	StartIndex int64 `json:"startIndex,omitempty"`
	// TabId: The tab that contains this range. When omitted, the request applies
	// to the first tab. In a document containing a single tab: - If provided, must
	// match the singular tab's ID. - If omitted, the request applies to the
	// singular tab. In a document containing multiple tabs: - If provided, the
	// request applies to the specified tab. - If omitted, the request applies to
	// the first tab in the document.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EndIndex") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EndIndex") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Range: Specifies a contiguous range of text.

func (Range) MarshalJSON ¶
func (s Range) MarshalJSON() ([]byte, error)
type RejectSuggestionRequest ¶
type RejectSuggestionRequest struct {
	// SuggestionId: The ID of the suggestion.
	SuggestionId string `json:"suggestionId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SuggestionId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SuggestionId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
RejectSuggestionRequest: Rejects a suggestion. Returns a 403 forbidden error if the requesting user does not have edit access to the document and is not the author of the suggestion.

func (RejectSuggestionRequest) MarshalJSON ¶
func (s RejectSuggestionRequest) MarshalJSON() ([]byte, error)
type ReplaceAllTextRequest ¶
type ReplaceAllTextRequest struct {
	// ContainsText: Finds text in the document matching this substring.
	ContainsText *SubstringMatchCriteria `json:"containsText,omitempty"`
	// ReplaceText: The text that will replace the matched text.
	ReplaceText string `json:"replaceText,omitempty"`
	// TabsCriteria: Optional. The criteria used to specify in which tabs the
	// replacement occurs. When omitted, the replacement applies to all tabs. In a
	// document containing a single tab: - If provided, must match the singular
	// tab's ID. - If omitted, the replacement applies to the singular tab. In a
	// document containing multiple tabs: - If provided, the replacement applies to
	// the specified tabs. - If omitted, the replacement applies to all tabs.
	TabsCriteria *TabsCriteria `json:"tabsCriteria,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ContainsText") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ContainsText") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ReplaceAllTextRequest: Replaces all instances of text matching a criteria with replace text.

func (ReplaceAllTextRequest) MarshalJSON ¶
func (s ReplaceAllTextRequest) MarshalJSON() ([]byte, error)
type ReplaceAllTextResponse ¶
type ReplaceAllTextResponse struct {
	// OccurrencesChanged: The number of occurrences changed by replacing all text.
	OccurrencesChanged int64 `json:"occurrencesChanged,omitempty"`
	// ForceSendFields is a list of field names (e.g. "OccurrencesChanged") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "OccurrencesChanged") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ReplaceAllTextResponse: The result of replacing text.

func (ReplaceAllTextResponse) MarshalJSON ¶
func (s ReplaceAllTextResponse) MarshalJSON() ([]byte, error)
type ReplaceImageRequest ¶
added in v0.10.0
type ReplaceImageRequest struct {
	// ImageObjectId: The ID of the existing image that will be replaced. The ID
	// can be retrieved from the response of a get request.
	ImageObjectId string `json:"imageObjectId,omitempty"`
	// ImageReplaceMethod: The replacement method.
	//
	// Possible values:
	//   "IMAGE_REPLACE_METHOD_UNSPECIFIED" - Unspecified image replace method.
	// This value must not be used.
	//   "CENTER_CROP" - Scales and centers the image to fill the bounds of the
	// original image. The image may be cropped in order to fill the original
	// image's bounds. The rendered size of the image will be the same as the
	// original image.
	ImageReplaceMethod string `json:"imageReplaceMethod,omitempty"`
	// TabId: The tab that the image to be replaced is in. When omitted, the
	// request is applied to the first tab. In a document containing a single tab:
	// - If provided, must match the singular tab's ID. - If omitted, the request
	// applies to the singular tab. In a document containing multiple tabs: - If
	// provided, the request applies to the specified tab. - If omitted, the
	// request applies to the first tab in the document.
	TabId string `json:"tabId,omitempty"`
	// Uri: The URI of the new image. The image is fetched once at insertion time
	// and a copy is stored for display inside the document. Images must be less
	// than 50MB, cannot exceed 25 megapixels, and must be in PNG, JPEG, or GIF
	// format. The provided URI can't surpass 2 KB in length. The URI is saved with
	// the image, and exposed through the ImageProperties.source_uri field.
	Uri string `json:"uri,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ImageObjectId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ImageObjectId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ReplaceImageRequest: Replaces an existing image with a new image. Replacing an image removes some image effects from the existing image in order to mirror the behavior of the Docs editor.

func (ReplaceImageRequest) MarshalJSON ¶
added in v0.10.0
func (s ReplaceImageRequest) MarshalJSON() ([]byte, error)
type ReplaceNamedRangeContentRequest ¶
added in v0.12.0
type ReplaceNamedRangeContentRequest struct {
	// NamedRangeId: The ID of the named range whose content will be replaced. If
	// there is no named range with the given ID a 400 bad request error is
	// returned.
	NamedRangeId string `json:"namedRangeId,omitempty"`
	// NamedRangeName: The name of the NamedRanges whose content will be replaced.
	// If there are multiple named ranges with the given name, then the content of
	// each one will be replaced. If there are no named ranges with the given name,
	// then the request will be a no-op.
	NamedRangeName string `json:"namedRangeName,omitempty"`
	// TabsCriteria: Optional. The criteria used to specify in which tabs the
	// replacement occurs. When omitted, the replacement applies to all tabs. In a
	// document containing a single tab: - If provided, must match the singular
	// tab's ID. - If omitted, the replacement applies to the singular tab. In a
	// document containing multiple tabs: - If provided, the replacement applies to
	// the specified tabs. - If omitted, the replacement applies to all tabs.
	TabsCriteria *TabsCriteria `json:"tabsCriteria,omitempty"`
	// Text: Replaces the content of the specified named range(s) with the given
	// text.
	Text string `json:"text,omitempty"`
	// ForceSendFields is a list of field names (e.g. "NamedRangeId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "NamedRangeId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ReplaceNamedRangeContentRequest: Replaces the contents of the specified NamedRange or NamedRanges with the given replacement content. Note that an individual NamedRange may consist of multiple discontinuous ranges. In this case, only the content in the first range will be replaced. The other ranges and their content will be deleted. In cases where replacing or deleting any ranges would result in an invalid document structure, a 400 bad request error is returned.

func (ReplaceNamedRangeContentRequest) MarshalJSON ¶
added in v0.12.0
func (s ReplaceNamedRangeContentRequest) MarshalJSON() ([]byte, error)
type Request ¶
type Request struct {
	// AcceptSuggestion: Accepts a suggestion.
	AcceptSuggestion *AcceptSuggestionRequest `json:"acceptSuggestion,omitempty"`
	// AddCommentReply: Adds a reply to a CommentThread or SuggestionThread.
	AddCommentReply *AddCommentReplyRequest `json:"addCommentReply,omitempty"`
	// AddDocumentTab: Adds a document tab.
	AddDocumentTab *AddDocumentTabRequest `json:"addDocumentTab,omitempty"`
	// CreateDropdownDefinition: Creates a DropdownDefinition.
	CreateDropdownDefinition *CreateDropdownDefinitionRequest `json:"createDropdownDefinition,omitempty"`
	// CreateFooter: Creates a footer.
	CreateFooter *CreateFooterRequest `json:"createFooter,omitempty"`
	// CreateFootnote: Creates a footnote.
	CreateFootnote *CreateFootnoteRequest `json:"createFootnote,omitempty"`
	// CreateHeader: Creates a header.
	CreateHeader *CreateHeaderRequest `json:"createHeader,omitempty"`
	// CreateNamedRange: Creates a named range.
	CreateNamedRange *CreateNamedRangeRequest `json:"createNamedRange,omitempty"`
	// CreateParagraphBullets: Creates bullets for paragraphs.
	CreateParagraphBullets *CreateParagraphBulletsRequest `json:"createParagraphBullets,omitempty"`
	// DeleteComment: Deletes a CommentThread.
	DeleteComment *DeleteCommentRequest `json:"deleteComment,omitempty"`
	// DeleteCommentReply: Deletes a reply Post from a CommentThread or
	// SuggestionThread.
	DeleteCommentReply *DeleteCommentReplyRequest `json:"deleteCommentReply,omitempty"`
	// DeleteContentRange: Deletes content from the document.
	DeleteContentRange *DeleteContentRangeRequest `json:"deleteContentRange,omitempty"`
	// DeleteDropdownDefinition: Deletes a DropdownDefinition.
	DeleteDropdownDefinition *DeleteDropdownDefinitionRequest `json:"deleteDropdownDefinition,omitempty"`
	// DeleteFooter: Deletes a footer from the document.
	DeleteFooter *DeleteFooterRequest `json:"deleteFooter,omitempty"`
	// DeleteHeader: Deletes a header from the document.
	DeleteHeader *DeleteHeaderRequest `json:"deleteHeader,omitempty"`
	// DeleteNamedRange: Deletes a named range.
	DeleteNamedRange *DeleteNamedRangeRequest `json:"deleteNamedRange,omitempty"`
	// DeleteParagraphBullets: Deletes bullets from paragraphs.
	DeleteParagraphBullets *DeleteParagraphBulletsRequest `json:"deleteParagraphBullets,omitempty"`
	// DeletePositionedObject: Deletes a positioned object from the document.
	DeletePositionedObject *DeletePositionedObjectRequest `json:"deletePositionedObject,omitempty"`
	// DeleteSuggestion: Deletes a suggestion.
	DeleteSuggestion *DeleteSuggestionRequest `json:"deleteSuggestion,omitempty"`
	// DeleteTab: Deletes a document tab.
	DeleteTab *DeleteTabRequest `json:"deleteTab,omitempty"`
	// DeleteTableColumn: Deletes a column from a table.
	DeleteTableColumn *DeleteTableColumnRequest `json:"deleteTableColumn,omitempty"`
	// DeleteTableRow: Deletes a row from a table.
	DeleteTableRow *DeleteTableRowRequest `json:"deleteTableRow,omitempty"`
	// InsertComment: Inserts a CommentThread into the document.
	InsertComment *InsertCommentRequest `json:"insertComment,omitempty"`
	// InsertDate: Inserts a date.
	InsertDate *InsertDateRequest `json:"insertDate,omitempty"`
	// InsertDropdown: Inserts a Dropdown at the specified location.
	InsertDropdown *InsertDropdownRequest `json:"insertDropdown,omitempty"`
	// InsertInlineImage: Inserts an inline image at the specified location.
	InsertInlineImage *InsertInlineImageRequest `json:"insertInlineImage,omitempty"`
	// InsertPageBreak: Inserts a page break at the specified location.
	InsertPageBreak *InsertPageBreakRequest `json:"insertPageBreak,omitempty"`
	// InsertPerson: Inserts a person mention.
	InsertPerson *InsertPersonRequest `json:"insertPerson,omitempty"`
	// InsertRichLink: Insert a rich link.
	InsertRichLink *InsertRichLinkRequest `json:"insertRichLink,omitempty"`
	// InsertSectionBreak: Inserts a section break at the specified location.
	InsertSectionBreak *InsertSectionBreakRequest `json:"insertSectionBreak,omitempty"`
	// InsertTable: Inserts a table at the specified location.
	InsertTable *InsertTableRequest `json:"insertTable,omitempty"`
	// InsertTableColumn: Inserts an empty column into a table.
	InsertTableColumn *InsertTableColumnRequest `json:"insertTableColumn,omitempty"`
	// InsertTableRow: Inserts an empty row into a table.
	InsertTableRow *InsertTableRowRequest `json:"insertTableRow,omitempty"`
	// InsertText: Inserts text at the specified location.
	InsertText *InsertTextRequest `json:"insertText,omitempty"`
	// MergeTableCells: Merges cells in a table.
	MergeTableCells *MergeTableCellsRequest `json:"mergeTableCells,omitempty"`
	// PinTableHeaderRows: Updates the number of pinned header rows in a table.
	PinTableHeaderRows *PinTableHeaderRowsRequest `json:"pinTableHeaderRows,omitempty"`
	// RejectSuggestion: Rejects a suggestion.
	RejectSuggestion *RejectSuggestionRequest `json:"rejectSuggestion,omitempty"`
	// ReplaceAllText: Replaces all instances of the specified text.
	ReplaceAllText *ReplaceAllTextRequest `json:"replaceAllText,omitempty"`
	// ReplaceImage: Replaces an image in the document.
	ReplaceImage *ReplaceImageRequest `json:"replaceImage,omitempty"`
	// ReplaceNamedRangeContent: Replaces the content in a named range.
	ReplaceNamedRangeContent *ReplaceNamedRangeContentRequest `json:"replaceNamedRangeContent,omitempty"`
	// UnmergeTableCells: Unmerges cells in a table.
	UnmergeTableCells *UnmergeTableCellsRequest `json:"unmergeTableCells,omitempty"`
	// UpdateCommentPost: Updates an existing post (head post or reply) of a
	// CommentThread or SuggestionThread.
	UpdateCommentPost *UpdateCommentPostRequest `json:"updateCommentPost,omitempty"`
	// UpdateDocumentStyle: Updates the style of the document.
	UpdateDocumentStyle *UpdateDocumentStyleRequest `json:"updateDocumentStyle,omitempty"`
	// UpdateDocumentTabProperties: Updates the properties of a document tab.
	UpdateDocumentTabProperties *UpdateDocumentTabPropertiesRequest `json:"updateDocumentTabProperties,omitempty"`
	// UpdateDropdownDefinitionProperties: Updates the properties of a
	// DropdownDefinition.
	UpdateDropdownDefinitionProperties *UpdateDropdownDefinitionPropertiesRequest `json:"updateDropdownDefinitionProperties,omitempty"`
	// UpdateDropdownProperties: Updates the properties of a Dropdown.
	UpdateDropdownProperties *UpdateDropdownPropertiesRequest `json:"updateDropdownProperties,omitempty"`
	// UpdateNamedStyle: Updates a named style.
	UpdateNamedStyle *UpdateNamedStyleRequest `json:"updateNamedStyle,omitempty"`
	// UpdateParagraphStyle: Updates the paragraph style at the specified range.
	UpdateParagraphStyle *UpdateParagraphStyleRequest `json:"updateParagraphStyle,omitempty"`
	// UpdateSectionStyle: Updates the section style of the specified range.
	UpdateSectionStyle *UpdateSectionStyleRequest `json:"updateSectionStyle,omitempty"`
	// UpdateTableCellStyle: Updates the style of table cells.
	UpdateTableCellStyle *UpdateTableCellStyleRequest `json:"updateTableCellStyle,omitempty"`
	// UpdateTableColumnProperties: Updates the properties of columns in a table.
	UpdateTableColumnProperties *UpdateTableColumnPropertiesRequest `json:"updateTableColumnProperties,omitempty"`
	// UpdateTableRowStyle: Updates the row style in a table.
	UpdateTableRowStyle *UpdateTableRowStyleRequest `json:"updateTableRowStyle,omitempty"`
	// UpdateTextStyle: Updates the text style at the specified range.
	UpdateTextStyle *UpdateTextStyleRequest `json:"updateTextStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "AcceptSuggestion") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "AcceptSuggestion") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Request: A single update to apply to a document.

func (Request) MarshalJSON ¶
func (s Request) MarshalJSON() ([]byte, error)
type Response ¶
type Response struct {
	// AddCommentReply: The result of adding a reply to a comment or suggestion.
	AddCommentReply *AddCommentReplyResponse `json:"addCommentReply,omitempty"`
	// AddDocumentTab: The result of adding a document tab.
	AddDocumentTab *AddDocumentTabResponse `json:"addDocumentTab,omitempty"`
	// CreateDropdownDefinition: The result of creating a dropdown definition.
	CreateDropdownDefinition *CreateDropdownDefinitionResponse `json:"createDropdownDefinition,omitempty"`
	// CreateFooter: The result of creating a footer.
	CreateFooter *CreateFooterResponse `json:"createFooter,omitempty"`
	// CreateFootnote: The result of creating a footnote.
	CreateFootnote *CreateFootnoteResponse `json:"createFootnote,omitempty"`
	// CreateHeader: The result of creating a header.
	CreateHeader *CreateHeaderResponse `json:"createHeader,omitempty"`
	// CreateNamedRange: The result of creating a named range.
	CreateNamedRange *CreateNamedRangeResponse `json:"createNamedRange,omitempty"`
	// InsertComment: The result of inserting a comment.
	InsertComment *InsertCommentResponse `json:"insertComment,omitempty"`
	// InsertDropdown: The result of inserting a dropdown.
	InsertDropdown *InsertDropdownResponse `json:"insertDropdown,omitempty"`
	// InsertInlineImage: The result of inserting an inline image.
	InsertInlineImage *InsertInlineImageResponse `json:"insertInlineImage,omitempty"`
	// InsertInlineSheetsChart: The result of inserting an inline Google Sheets
	// chart.
	InsertInlineSheetsChart *InsertInlineSheetsChartResponse `json:"insertInlineSheetsChart,omitempty"`
	// ReplaceAllText: The result of replacing text.
	ReplaceAllText *ReplaceAllTextResponse `json:"replaceAllText,omitempty"`
	// ForceSendFields is a list of field names (e.g. "AddCommentReply") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "AddCommentReply") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Response: A single response from an update.

func (Response) MarshalJSON ¶
func (s Response) MarshalJSON() ([]byte, error)
type RgbColor ¶
type RgbColor struct {
	// Blue: The blue component of the color, from 0.0 to 1.0.
	Blue float64 `json:"blue,omitempty"`
	// Green: The green component of the color, from 0.0 to 1.0.
	Green float64 `json:"green,omitempty"`
	// Red: The red component of the color, from 0.0 to 1.0.
	Red float64 `json:"red,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Blue") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Blue") to include in API requests
	// with the JSON null value. By default, fields with empty values are omitted
	// from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
RgbColor: An RGB color.

func (RgbColor) MarshalJSON ¶
func (s RgbColor) MarshalJSON() ([]byte, error)
func (*RgbColor) UnmarshalJSON ¶
func (s *RgbColor) UnmarshalJSON(data []byte) error
type RichLink ¶
added in v0.51.0
type RichLink struct {
	// RichLinkId: Output only. The ID of this link.
	RichLinkId string `json:"richLinkId,omitempty"`
	// RichLinkProperties: Output only. The properties of this RichLink. This field
	// is always present.
	RichLinkProperties *RichLinkProperties `json:"richLinkProperties,omitempty"`
	// SuggestedDeletionIds: IDs for suggestions that remove this link from the
	// document. A RichLink might have multiple deletion IDs if, for example,
	// multiple users suggest deleting it. If empty, then this person link isn't
	// suggested for deletion.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: IDs for suggestions that insert this link into the
	// document. A RichLink might have multiple insertion IDs if it's a nested
	// suggested change (a suggestion within a suggestion made by a different user,
	// for example). If empty, then this person link isn't a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTextStyleChanges: The suggested text style changes to this
	// RichLink, keyed by suggestion ID.
	SuggestedTextStyleChanges map[string]SuggestedTextStyle `json:"suggestedTextStyleChanges,omitempty"`
	// TextStyle: The text style of this RichLink.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "RichLinkId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "RichLinkId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
RichLink: A link to a Google resource (such as a file in Drive, a YouTube video, or a Calendar event).

func (RichLink) MarshalJSON ¶
added in v0.51.0
func (s RichLink) MarshalJSON() ([]byte, error)
type RichLinkProperties ¶
added in v0.51.0
type RichLinkProperties struct {
	// MimeType: The MIME type
	// (https://developers.google.com/drive/api/v3/mime-types) of the RichLink, if
	// there's one (for example, when it's a file in Drive).
	MimeType string `json:"mimeType,omitempty"`
	// Title: The title of the RichLink as displayed in the link. This title
	// matches the title of the linked resource at the time of the insertion or
	// last update of the link. This field is always present.
	Title string `json:"title,omitempty"`
	// Uri: The URI to the RichLink. This is always present.
	Uri string `json:"uri,omitempty"`
	// ForceSendFields is a list of field names (e.g. "MimeType") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "MimeType") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
RichLinkProperties: Properties specific to a RichLink.

func (RichLinkProperties) MarshalJSON ¶
added in v0.51.0
func (s RichLinkProperties) MarshalJSON() ([]byte, error)
type SectionBreak ¶
type SectionBreak struct {
	// SectionStyle: The style of the section after this section break.
	SectionStyle *SectionStyle `json:"sectionStyle,omitempty"`
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. A SectionBreak may have
	// multiple insertion IDs if it's a nested suggested change. If empty, then
	// this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// ForceSendFields is a list of field names (e.g. "SectionStyle") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "SectionStyle") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SectionBreak: A StructuralElement representing a section break. A section is a range of content that has the same SectionStyle. A section break represents the start of a new section, and the section style applies to the section after the section break. The document body always begins with a section break.

func (SectionBreak) MarshalJSON ¶
func (s SectionBreak) MarshalJSON() ([]byte, error)
type SectionColumnProperties ¶
type SectionColumnProperties struct {
	// PaddingEnd: The padding at the end of the column.
	PaddingEnd *Dimension `json:"paddingEnd,omitempty"`
	// Width: Output only. The width of the column.
	Width *Dimension `json:"width,omitempty"`
	// ForceSendFields is a list of field names (e.g. "PaddingEnd") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "PaddingEnd") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SectionColumnProperties: Properties that apply to a section's column.

func (SectionColumnProperties) MarshalJSON ¶
func (s SectionColumnProperties) MarshalJSON() ([]byte, error)
type SectionStyle ¶
type SectionStyle struct {
	// ColumnProperties: The section's columns properties. If empty, the section
	// contains one column with the default properties in the Docs editor. A
	// section can be updated to have no more than 3 columns. When updating this
	// property, setting a concrete value is required. Unsetting this property will
	// result in a 400 bad request error.
	ColumnProperties []*SectionColumnProperties `json:"columnProperties,omitempty"`
	// ColumnSeparatorStyle: The style of column separators. This style can be set
	// even when there's one column in the section. When updating this property,
	// setting a concrete value is required. Unsetting this property results in a
	// 400 bad request error.
	//
	// Possible values:
	//   "COLUMN_SEPARATOR_STYLE_UNSPECIFIED" - An unspecified column separator
	// style.
	//   "NONE" - No column separator lines between columns.
	//   "BETWEEN_EACH_COLUMN" - Renders a column separator line between each
	// column.
	ColumnSeparatorStyle string `json:"columnSeparatorStyle,omitempty"`
	// ContentDirection: The content direction of this section. If unset, the value
	// defaults to LEFT_TO_RIGHT. When updating this property, setting a concrete
	// value is required. Unsetting this property results in a 400 bad request
	// error.
	//
	// Possible values:
	//   "CONTENT_DIRECTION_UNSPECIFIED" - The content direction is unspecified.
	//   "LEFT_TO_RIGHT" - The content goes from left to right.
	//   "RIGHT_TO_LEFT" - The content goes from right to left.
	ContentDirection string `json:"contentDirection,omitempty"`
	// DefaultFooterId: The ID of the default footer. If unset, the value inherits
	// from the previous SectionBreak's SectionStyle. If the value is unset in the
	// first SectionBreak, it inherits from DocumentStyle's default_footer_id. If
	// DocumentMode is PAGELESS, this property will not be rendered. This property
	// is read-only.
	DefaultFooterId string `json:"defaultFooterId,omitempty"`
	// DefaultHeaderId: The ID of the default header. If unset, the value inherits
	// from the previous SectionBreak's SectionStyle. If the value is unset in the
	// first SectionBreak, it inherits from DocumentStyle's default_header_id. If
	// DocumentMode is PAGELESS, this property will not be rendered. This property
	// is read-only.
	DefaultHeaderId string `json:"defaultHeaderId,omitempty"`
	// EvenPageFooterId: The ID of the footer used only for even pages. If the
	// value of DocumentStyle's use_even_page_header_footer is true, this value is
	// used for the footers on even pages in the section. If it is false, the
	// footers on even pages use the default_footer_id. If unset, the value
	// inherits from the previous SectionBreak's SectionStyle. If the value is
	// unset in the first SectionBreak, it inherits from DocumentStyle's
	// even_page_footer_id. If DocumentMode is PAGELESS, this property will not be
	// rendered. This property is read-only.
	EvenPageFooterId string `json:"evenPageFooterId,omitempty"`
	// EvenPageHeaderId: The ID of the header used only for even pages. If the
	// value of DocumentStyle's use_even_page_header_footer is true, this value is
	// used for the headers on even pages in the section. If it is false, the
	// headers on even pages use the default_header_id. If unset, the value
	// inherits from the previous SectionBreak's SectionStyle. If the value is
	// unset in the first SectionBreak, it inherits from DocumentStyle's
	// even_page_header_id. If DocumentMode is PAGELESS, this property will not be
	// rendered. This property is read-only.
	EvenPageHeaderId string `json:"evenPageHeaderId,omitempty"`
	// FirstPageFooterId: The ID of the footer used only for the first page of the
	// section. If use_first_page_header_footer is true, this value is used for the
	// footer on the first page of the section. If it's false, the footer on the
	// first page of the section uses the default_footer_id. If unset, the value
	// inherits from the previous SectionBreak's SectionStyle. If the value is
	// unset in the first SectionBreak, it inherits from DocumentStyle's
	// first_page_footer_id. If DocumentMode is PAGELESS, this property will not be
	// rendered. This property is read-only.
	FirstPageFooterId string `json:"firstPageFooterId,omitempty"`
	// FirstPageHeaderId: The ID of the header used only for the first page of the
	// section. If use_first_page_header_footer is true, this value is used for the
	// header on the first page of the section. If it's false, the header on the
	// first page of the section uses the default_header_id. If unset, the value
	// inherits from the previous SectionBreak's SectionStyle. If the value is
	// unset in the first SectionBreak, it inherits from DocumentStyle's
	// first_page_header_id. If DocumentMode is PAGELESS, this property will not be
	// rendered. This property is read-only.
	FirstPageHeaderId string `json:"firstPageHeaderId,omitempty"`
	// FlipPageOrientation: Optional. Indicates whether to flip the dimensions of
	// DocumentStyle's page_size for this section, which allows changing the page
	// orientation between portrait and landscape. If unset, the value inherits
	// from DocumentStyle's flip_page_orientation. If DocumentMode is PAGELESS,
	// this property will not be rendered. When updating this property, setting a
	// concrete value is required. Unsetting this property results in a 400 bad
	// request error.
	FlipPageOrientation bool `json:"flipPageOrientation,omitempty"`
	// MarginBottom: The bottom page margin of the section. If unset, the value
	// defaults to margin_bottom from DocumentStyle. If DocumentMode is PAGELESS,
	// this property will not be rendered. When updating this property, setting a
	// concrete value is required. Unsetting this property results in a 400 bad
	// request error.
	MarginBottom *Dimension `json:"marginBottom,omitempty"`
	// MarginFooter: The footer margin of the section. If unset, the value defaults
	// to margin_footer from DocumentStyle. If updated,
	// use_custom_header_footer_margins is set to true on DocumentStyle. The value
	// of use_custom_header_footer_margins on DocumentStyle indicates if a footer
	// margin is being respected for this section If DocumentMode is PAGELESS, this
	// property will not be rendered. When updating this property, setting a
	// concrete value is required. Unsetting this property results in a 400 bad
	// request error.
	MarginFooter *Dimension `json:"marginFooter,omitempty"`
	// MarginHeader: The header margin of the section. If unset, the value defaults
	// to margin_header from DocumentStyle. If updated,
	// use_custom_header_footer_margins is set to true on DocumentStyle. The value
	// of use_custom_header_footer_margins on DocumentStyle indicates if a header
	// margin is being respected for this section. If DocumentMode is PAGELESS,
	// this property will not be rendered. When updating this property, setting a
	// concrete value is required. Unsetting this property results in a 400 bad
	// request error.
	MarginHeader *Dimension `json:"marginHeader,omitempty"`
	// MarginLeft: The left page margin of the section. If unset, the value
	// defaults to margin_left from DocumentStyle. Updating the left margin causes
	// columns in this section to resize. Since the margin affects column width,
	// it's applied before column properties. If DocumentMode is PAGELESS, this
	// property will not be rendered. When updating this property, setting a
	// concrete value is required. Unsetting this property results in a 400 bad
	// request error.
	MarginLeft *Dimension `json:"marginLeft,omitempty"`
	// MarginRight: The right page margin of the section. If unset, the value
	// defaults to margin_right from DocumentStyle. Updating the right margin
	// causes columns in this section to resize. Since the margin affects column
	// width, it's applied before column properties. If DocumentMode is PAGELESS,
	// this property will not be rendered. When updating this property, setting a
	// concrete value is required. Unsetting this property results in a 400 bad
	// request error.
	MarginRight *Dimension `json:"marginRight,omitempty"`
	// MarginTop: The top page margin of the section. If unset, the value defaults
	// to margin_top from DocumentStyle. If DocumentMode is PAGELESS, this property
	// will not be rendered. When updating this property, setting a concrete value
	// is required. Unsetting this property results in a 400 bad request error.
	MarginTop *Dimension `json:"marginTop,omitempty"`
	// PageNumberStart: The page number from which to start counting the number of
	// pages for this section. If unset, page numbering continues from the previous
	// section. If the value is unset in the first SectionBreak, refer to
	// DocumentStyle's page_number_start. If DocumentMode is PAGELESS, this
	// property will not be rendered. When updating this property, setting a
	// concrete value is required. Unsetting this property results in a 400 bad
	// request error.
	PageNumberStart int64 `json:"pageNumberStart,omitempty"`
	// SectionType: Output only. The type of section.
	//
	// Possible values:
	//   "SECTION_TYPE_UNSPECIFIED" - The section type is unspecified.
	//   "CONTINUOUS" - The section starts immediately after the last paragraph of
	// the previous section.
	//   "NEXT_PAGE" - The section starts on the next page.
	SectionType string `json:"sectionType,omitempty"`
	// UseFirstPageHeaderFooter: Indicates whether to use the first page header /
	// footer IDs for the first page of the section. If unset, it inherits from
	// DocumentStyle's use_first_page_header_footer for the first section. If the
	// value is unset for subsequent sectors, it should be interpreted as false. If
	// DocumentMode is PAGELESS, this property will not be rendered. When updating
	// this property, setting a concrete value is required. Unsetting this property
	// results in a 400 bad request error.
	UseFirstPageHeaderFooter bool `json:"useFirstPageHeaderFooter,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ColumnProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ColumnProperties") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SectionStyle: The styling that applies to a section.

func (SectionStyle) MarshalJSON ¶
func (s SectionStyle) MarshalJSON() ([]byte, error)
type Service ¶
type Service struct {
	BasePath  string // API endpoint base URL
	UserAgent string // optional additional User-Agent fragment

	Documents *DocumentsService
	// contains filtered or unexported fields
}
func
New
deprecated
func NewService ¶
added in v0.3.0
func NewService(ctx context.Context, opts ...option.ClientOption) (*Service, error)
NewService creates a new Service.

type Shading ¶
type Shading struct {
	// BackgroundColor: The background color of this paragraph shading.
	BackgroundColor *OptionalColor `json:"backgroundColor,omitempty"`
	// ForceSendFields is a list of field names (e.g. "BackgroundColor") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "BackgroundColor") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Shading: The shading of a paragraph.

func (Shading) MarshalJSON ¶
func (s Shading) MarshalJSON() ([]byte, error)
type ShadingSuggestionState ¶
type ShadingSuggestionState struct {
	// BackgroundColorSuggested: Indicates if there was a suggested change to the
	// Shading.
	BackgroundColorSuggested bool `json:"backgroundColorSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "BackgroundColorSuggested")
	// to unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "BackgroundColorSuggested") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
ShadingSuggestionState: A mask that indicates which of the fields on the base Shading have been changed in this suggested change. For any field set to true, there's a new suggested value.

func (ShadingSuggestionState) MarshalJSON ¶
func (s ShadingSuggestionState) MarshalJSON() ([]byte, error)
type SheetsChartReference ¶
type SheetsChartReference struct {
	// ChartId: The ID of the specific chart in the Google Sheets spreadsheet
	// that's embedded.
	ChartId int64 `json:"chartId,omitempty"`
	// SpreadsheetId: The ID of the Google Sheets spreadsheet that contains the
	// source chart.
	SpreadsheetId string `json:"spreadsheetId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ChartId") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ChartId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SheetsChartReference: A reference to a linked chart embedded from Google Sheets.

func (SheetsChartReference) MarshalJSON ¶
func (s SheetsChartReference) MarshalJSON() ([]byte, error)
type SheetsChartReferenceSuggestionState ¶
type SheetsChartReferenceSuggestionState struct {
	// ChartIdSuggested: Indicates if there was a suggested change to chart_id.
	ChartIdSuggested bool `json:"chartIdSuggested,omitempty"`
	// SpreadsheetIdSuggested: Indicates if there was a suggested change to
	// spreadsheet_id.
	SpreadsheetIdSuggested bool `json:"spreadsheetIdSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ChartIdSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ChartIdSuggested") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SheetsChartReferenceSuggestionState: A mask that indicates which of the fields on the base SheetsChartReference have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (SheetsChartReferenceSuggestionState) MarshalJSON ¶
func (s SheetsChartReferenceSuggestionState) MarshalJSON() ([]byte, error)
type Size ¶
type Size struct {
	// Height: The height of the object.
	Height *Dimension `json:"height,omitempty"`
	// Width: The width of the object.
	Width *Dimension `json:"width,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Height") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Height") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Size: A width and height.

func (Size) MarshalJSON ¶
func (s Size) MarshalJSON() ([]byte, error)
type SizeSuggestionState ¶
type SizeSuggestionState struct {
	// HeightSuggested: Indicates if there was a suggested change to height.
	HeightSuggested bool `json:"heightSuggested,omitempty"`
	// WidthSuggested: Indicates if there was a suggested change to width.
	WidthSuggested bool `json:"widthSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "HeightSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "HeightSuggested") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SizeSuggestionState: A mask that indicates which of the fields on the base Size have been changed in this suggestion. For any field set to true, the Size has a new suggested value.

func (SizeSuggestionState) MarshalJSON ¶
func (s SizeSuggestionState) MarshalJSON() ([]byte, error)
type StructuralElement ¶
type StructuralElement struct {
	// EndIndex: The zero-based end index of this structural element, exclusive, in
	// UTF-16 code units.
	EndIndex int64 `json:"endIndex,omitempty"`
	// Paragraph: A paragraph type of structural element.
	Paragraph *Paragraph `json:"paragraph,omitempty"`
	// SectionBreak: A section break type of structural element.
	SectionBreak *SectionBreak `json:"sectionBreak,omitempty"`
	// StartIndex: The zero-based start index of this structural element, in UTF-16
	// code units.
	StartIndex int64 `json:"startIndex,omitempty"`
	// Table: A table type of structural element.
	Table *Table `json:"table,omitempty"`
	// TableOfContents: A table of contents type of structural element.
	TableOfContents *TableOfContents `json:"tableOfContents,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EndIndex") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EndIndex") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
StructuralElement: A StructuralElement describes content that provides structure to the document.

func (StructuralElement) MarshalJSON ¶
func (s StructuralElement) MarshalJSON() ([]byte, error)
type SubstringMatchCriteria ¶
type SubstringMatchCriteria struct {
	// MatchCase: Indicates whether the search should respect case: - `True`: the
	// search is case sensitive. - `False`: the search is case insensitive.
	MatchCase bool `json:"matchCase,omitempty"`
	// SearchByRegex: Optional. True if the find value should be treated as a
	// regular expression. Any backslashes in the pattern should be escaped. -
	// `True`: the search text is treated as a regular expressions. - `False`: the
	// search text is treated as a substring for matching.
	SearchByRegex bool `json:"searchByRegex,omitempty"`
	// Text: The text to search for in the document.
	Text string `json:"text,omitempty"`
	// ForceSendFields is a list of field names (e.g. "MatchCase") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "MatchCase") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SubstringMatchCriteria: A criteria that matches a specific string of text in the document.

func (SubstringMatchCriteria) MarshalJSON ¶
func (s SubstringMatchCriteria) MarshalJSON() ([]byte, error)
type SuggestedBullet ¶
type SuggestedBullet struct {
	// Bullet: A Bullet that only includes the changes made in this suggestion.
	// This can be used along with the bullet_suggestion_state to see which fields
	// have changed and their new values.
	Bullet *Bullet `json:"bullet,omitempty"`
	// BulletSuggestionState: A mask that indicates which of the fields on the base
	// Bullet have been changed in this suggestion.
	BulletSuggestionState *BulletSuggestionState `json:"bulletSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Bullet") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Bullet") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedBullet: A suggested change to a Bullet.

func (SuggestedBullet) MarshalJSON ¶
func (s SuggestedBullet) MarshalJSON() ([]byte, error)
type SuggestedDateElementProperties ¶
added in v0.259.0
type SuggestedDateElementProperties struct {
	// DateElementProperties: DateElementProperties that only includes the changes
	// made in this suggestion. This can be used along with the
	// date_element_properties_suggestion_state to see which fields have changed
	// and their new values.
	DateElementProperties *DateElementProperties `json:"dateElementProperties,omitempty"`
	// DateElementPropertiesSuggestionState: A mask that indicates which of the
	// fields on the base DateElementProperties have been changed in this
	// suggestion.
	DateElementPropertiesSuggestionState *DateElementPropertiesSuggestionState `json:"dateElementPropertiesSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DateElementProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DateElementProperties") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedDateElementProperties: A suggested change to a DateElementProperties.

func (SuggestedDateElementProperties) MarshalJSON ¶
added in v0.259.0
func (s SuggestedDateElementProperties) MarshalJSON() ([]byte, error)
type SuggestedDocumentStyle ¶
type SuggestedDocumentStyle struct {
	// DocumentStyle: A DocumentStyle that only includes the changes made in this
	// suggestion. This can be used along with the document_style_suggestion_state
	// to see which fields have changed and their new values.
	DocumentStyle *DocumentStyle `json:"documentStyle,omitempty"`
	// DocumentStyleSuggestionState: A mask that indicates which of the fields on
	// the base DocumentStyle have been changed in this suggestion.
	DocumentStyleSuggestionState *DocumentStyleSuggestionState `json:"documentStyleSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DocumentStyle") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DocumentStyle") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedDocumentStyle: A suggested change to the DocumentStyle.

func (SuggestedDocumentStyle) MarshalJSON ¶
func (s SuggestedDocumentStyle) MarshalJSON() ([]byte, error)
type SuggestedDropdownDefinitionProperties ¶
type SuggestedDropdownDefinitionProperties struct {
	// DropdownDefinitionProperties: A DropdownDefinitionProperties that only
	// includes the changes made in this suggestion. This can be used along with
	// the dropdown_definition_properties_suggestion_state to see which fields have
	// changed and their new values.
	DropdownDefinitionProperties *DropdownDefinitionProperties `json:"dropdownDefinitionProperties,omitempty"`
	// DropdownDefinitionPropertiesSuggestionState: A mask that indicates which of
	// the fields on the base DropdownDefinitionProperties have been changed in
	// this suggestion.
	DropdownDefinitionPropertiesSuggestionState *DropdownDefinitionPropertiesSuggestionState `json:"dropdownDefinitionPropertiesSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g.
	// "DropdownDefinitionProperties") to unconditionally include in API requests.
	// By default, fields with empty or default values are omitted from API
	// requests. See https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields
	// for more details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DropdownDefinitionProperties") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedDropdownDefinitionProperties: A suggested change to dropdown definition properties.

func (SuggestedDropdownDefinitionProperties) MarshalJSON ¶
func (s SuggestedDropdownDefinitionProperties) MarshalJSON() ([]byte, error)
type SuggestedDropdownProperties ¶
type SuggestedDropdownProperties struct {
	// DropdownProperties: A DropdownProperties that only includes the changes made
	// in this suggestion. This can be used along with the
	// dropdown_properties_suggestion_state to see which fields have changed and
	// their new values.
	DropdownProperties *DropdownProperties `json:"dropdownProperties,omitempty"`
	// DropdownPropertiesSuggestionState: A mask that indicates which of the fields
	// on the base DropdownProperties have been changed in this suggestion.
	DropdownPropertiesSuggestionState *DropdownPropertiesSuggestionState `json:"dropdownPropertiesSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DropdownProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DropdownProperties") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedDropdownProperties: A suggested change to dropdown properties.

func (SuggestedDropdownProperties) MarshalJSON ¶
func (s SuggestedDropdownProperties) MarshalJSON() ([]byte, error)
type SuggestedInlineObjectProperties ¶
type SuggestedInlineObjectProperties struct {
	// InlineObjectProperties: An InlineObjectProperties that only includes the
	// changes made in this suggestion. This can be used along with the
	// inline_object_properties_suggestion_state to see which fields have changed
	// and their new values.
	InlineObjectProperties *InlineObjectProperties `json:"inlineObjectProperties,omitempty"`
	// InlineObjectPropertiesSuggestionState: A mask that indicates which of the
	// fields on the base InlineObjectProperties have been changed in this
	// suggestion.
	InlineObjectPropertiesSuggestionState *InlineObjectPropertiesSuggestionState `json:"inlineObjectPropertiesSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "InlineObjectProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "InlineObjectProperties") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedInlineObjectProperties: A suggested change to InlineObjectProperties.

func (SuggestedInlineObjectProperties) MarshalJSON ¶
func (s SuggestedInlineObjectProperties) MarshalJSON() ([]byte, error)
type SuggestedListProperties ¶
type SuggestedListProperties struct {
	// ListProperties: A ListProperties that only includes the changes made in this
	// suggestion. This can be used along with the list_properties_suggestion_state
	// to see which fields have changed and their new values.
	ListProperties *ListProperties `json:"listProperties,omitempty"`
	// ListPropertiesSuggestionState: A mask that indicates which of the fields on
	// the base ListProperties have been changed in this suggestion.
	ListPropertiesSuggestionState *ListPropertiesSuggestionState `json:"listPropertiesSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ListProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ListProperties") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedListProperties: A suggested change to ListProperties.

func (SuggestedListProperties) MarshalJSON ¶
func (s SuggestedListProperties) MarshalJSON() ([]byte, error)
type SuggestedNamedStyles ¶
type SuggestedNamedStyles struct {
	// NamedStyles: A NamedStyles that only includes the changes made in this
	// suggestion. This can be used along with the named_styles_suggestion_state to
	// see which fields have changed and their new values.
	NamedStyles *NamedStyles `json:"namedStyles,omitempty"`
	// NamedStylesSuggestionState: A mask that indicates which of the fields on the
	// base NamedStyles have been changed in this suggestion.
	NamedStylesSuggestionState *NamedStylesSuggestionState `json:"namedStylesSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "NamedStyles") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "NamedStyles") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedNamedStyles: A suggested change to the NamedStyles.

func (SuggestedNamedStyles) MarshalJSON ¶
func (s SuggestedNamedStyles) MarshalJSON() ([]byte, error)
type SuggestedParagraphStyle ¶
type SuggestedParagraphStyle struct {
	// ParagraphStyle: A ParagraphStyle that only includes the changes made in this
	// suggestion. This can be used along with the paragraph_style_suggestion_state
	// to see which fields have changed and their new values.
	ParagraphStyle *ParagraphStyle `json:"paragraphStyle,omitempty"`
	// ParagraphStyleSuggestionState: A mask that indicates which of the fields on
	// the base ParagraphStyle have been changed in this suggestion.
	ParagraphStyleSuggestionState *ParagraphStyleSuggestionState `json:"paragraphStyleSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ParagraphStyle") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ParagraphStyle") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedParagraphStyle: A suggested change to a ParagraphStyle.

func (SuggestedParagraphStyle) MarshalJSON ¶
func (s SuggestedParagraphStyle) MarshalJSON() ([]byte, error)
type SuggestedPositionedObjectProperties ¶
type SuggestedPositionedObjectProperties struct {
	// PositionedObjectProperties: A PositionedObjectProperties that only includes
	// the changes made in this suggestion. This can be used along with the
	// positioned_object_properties_suggestion_state to see which fields have
	// changed and their new values.
	PositionedObjectProperties *PositionedObjectProperties `json:"positionedObjectProperties,omitempty"`
	// PositionedObjectPropertiesSuggestionState: A mask that indicates which of
	// the fields on the base PositionedObjectProperties have been changed in this
	// suggestion.
	PositionedObjectPropertiesSuggestionState *PositionedObjectPropertiesSuggestionState `json:"positionedObjectPropertiesSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "PositionedObjectProperties")
	// to unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "PositionedObjectProperties") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedPositionedObjectProperties: A suggested change to PositionedObjectProperties.

func (SuggestedPositionedObjectProperties) MarshalJSON ¶
func (s SuggestedPositionedObjectProperties) MarshalJSON() ([]byte, error)
type SuggestedTableCellStyle ¶
type SuggestedTableCellStyle struct {
	// TableCellStyle: A TableCellStyle that only includes the changes made in this
	// suggestion. This can be used along with the
	// table_cell_style_suggestion_state to see which fields have changed and their
	// new values.
	TableCellStyle *TableCellStyle `json:"tableCellStyle,omitempty"`
	// TableCellStyleSuggestionState: A mask that indicates which of the fields on
	// the base TableCellStyle have been changed in this suggestion.
	TableCellStyleSuggestionState *TableCellStyleSuggestionState `json:"tableCellStyleSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TableCellStyle") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TableCellStyle") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedTableCellStyle: A suggested change to a TableCellStyle.

func (SuggestedTableCellStyle) MarshalJSON ¶
func (s SuggestedTableCellStyle) MarshalJSON() ([]byte, error)
type SuggestedTableRowStyle ¶
type SuggestedTableRowStyle struct {
	// TableRowStyle: A TableRowStyle that only includes the changes made in this
	// suggestion. This can be used along with the table_row_style_suggestion_state
	// to see which fields have changed and their new values.
	TableRowStyle *TableRowStyle `json:"tableRowStyle,omitempty"`
	// TableRowStyleSuggestionState: A mask that indicates which of the fields on
	// the base TableRowStyle have been changed in this suggestion.
	TableRowStyleSuggestionState *TableRowStyleSuggestionState `json:"tableRowStyleSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TableRowStyle") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TableRowStyle") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedTableRowStyle: A suggested change to a TableRowStyle.

func (SuggestedTableRowStyle) MarshalJSON ¶
func (s SuggestedTableRowStyle) MarshalJSON() ([]byte, error)
type SuggestedTextStyle ¶
type SuggestedTextStyle struct {
	// TextStyle: A TextStyle that only includes the changes made in this
	// suggestion. This can be used along with the text_style_suggestion_state to
	// see which fields have changed and their new values.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// TextStyleSuggestionState: A mask that indicates which of the fields on the
	// base TextStyle have been changed in this suggestion.
	TextStyleSuggestionState *TextStyleSuggestionState `json:"textStyleSuggestionState,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TextStyle") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TextStyle") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestedTextStyle: A suggested change to a TextStyle.

func (SuggestedTextStyle) MarshalJSON ¶
func (s SuggestedTextStyle) MarshalJSON() ([]byte, error)
type SuggestionResponse ¶
type SuggestionResponse struct {
	// AcceptedSuggestionIds: The IDs of suggestions which were accepted during the
	// update.
	AcceptedSuggestionIds []string `json:"acceptedSuggestionIds,omitempty"`
	// CreatedSuggestionIds: The IDs of suggestions which were created during the
	// update.
	CreatedSuggestionIds []string `json:"createdSuggestionIds,omitempty"`
	// DeletedSuggestionIds: The IDs of suggestions which were deleted during the
	// update.
	DeletedSuggestionIds []string `json:"deletedSuggestionIds,omitempty"`
	// RejectedSuggestionIds: The IDs of suggestions which were rejected during the
	// update.
	RejectedSuggestionIds []string `json:"rejectedSuggestionIds,omitempty"`
	// UpdatedSummarySuggestionIds: The IDs of suggestions whose summaries were
	// updated during the update.
	UpdatedSummarySuggestionIds []string `json:"updatedSummarySuggestionIds,omitempty"`
	// ForceSendFields is a list of field names (e.g. "AcceptedSuggestionIds") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "AcceptedSuggestionIds") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestionResponse: The suggestions which were affected by a given update.

func (SuggestionResponse) MarshalJSON ¶
func (s SuggestionResponse) MarshalJSON() ([]byte, error)
type SuggestionThread ¶
type SuggestionThread struct {
	// HeadPost: The first post in the thread.
	HeadPost *Post `json:"headPost,omitempty"`
	// Replies: Replies to the head post.
	Replies []*Post `json:"replies,omitempty"`
	// Status: Whether the thread is open, accepted, or rejected.
	//
	// Possible values:
	//   "STATUS_UNSPECIFIED" - Default value. This value is unused.
	//   "OPEN" - The suggestion thread is open.
	//   "ACCEPTED" - The suggestion thread is accepted.
	//   "REJECTED" - The suggestion thread is rejected.
	Status string `json:"status,omitempty"`
	// SuggestionId: The unique ID of the suggestion.
	SuggestionId string `json:"suggestionId,omitempty"`
	// SummaryHtml: Summary of the suggested differences in the document, in HTML.
	// May be empty.
	SummaryHtml string `json:"summaryHtml,omitempty"`
	// SummaryText: Summary of the suggested differences in the document, in plain
	// text. May be empty.
	SummaryText string `json:"summaryText,omitempty"`
	// ForceSendFields is a list of field names (e.g. "HeadPost") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "HeadPost") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
SuggestionThread: Represents a single suggestion thread. Suggestion threads are created as a byproduct of saving changes to the document while in suggestion mode, and cannot be created directly.

func (SuggestionThread) MarshalJSON ¶
func (s SuggestionThread) MarshalJSON() ([]byte, error)
type Tab ¶
added in v0.190.0
type Tab struct {
	// ChildTabs: The child tabs nested within this tab.
	ChildTabs []*Tab `json:"childTabs,omitempty"`
	// DocumentTab: A tab with document contents, like text and images.
	DocumentTab *DocumentTab `json:"documentTab,omitempty"`
	// TabProperties: The properties of the tab, like ID and title.
	TabProperties *TabProperties `json:"tabProperties,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ChildTabs") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ChildTabs") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Tab: A tab in a document.

func (Tab) MarshalJSON ¶
added in v0.190.0
func (s Tab) MarshalJSON() ([]byte, error)
type TabProperties ¶
added in v0.190.0
type TabProperties struct {
	// IconEmoji: Optional. The emoji icon displayed with the tab. A valid emoji
	// icon is represented by a non-empty Unicode string. Any set of characters
	// that don't represent a single emoji is invalid. If an emoji is invalid, a
	// 400 bad request error is returned. If this value is unset or empty, the tab
	// will display the default tab icon.
	IconEmoji string `json:"iconEmoji,omitempty"`
	// Index: The zero-based index of the tab within the parent.
	Index int64 `json:"index,omitempty"`
	// NestingLevel: Output only. The depth of the tab within the document.
	// Root-level tabs start at 0.
	NestingLevel int64 `json:"nestingLevel,omitempty"`
	// ParentTabId: Optional. The ID of the parent tab. Empty when the current tab
	// is a root-level tab, which means it doesn't have any parents.
	ParentTabId string `json:"parentTabId,omitempty"`
	// TabId: The immutable ID of the tab.
	TabId string `json:"tabId,omitempty"`
	// Title: The user-visible name of the tab.
	Title string `json:"title,omitempty"`
	// ForceSendFields is a list of field names (e.g. "IconEmoji") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "IconEmoji") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TabProperties: Properties of a tab.

func (TabProperties) MarshalJSON ¶
added in v0.190.0
func (s TabProperties) MarshalJSON() ([]byte, error)
type TabStop ¶
type TabStop struct {
	// Alignment: The alignment of this tab stop. If unset, the value defaults to
	// START.
	//
	// Possible values:
	//   "TAB_STOP_ALIGNMENT_UNSPECIFIED" - The tab stop alignment is unspecified.
	//   "START" - The tab stop is aligned to the start of the line. This is the
	// default.
	//   "CENTER" - The tab stop is aligned to the center of the line.
	//   "END" - The tab stop is aligned to the end of the line.
	Alignment string `json:"alignment,omitempty"`
	// Offset: The offset between this tab stop and the start margin.
	Offset *Dimension `json:"offset,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Alignment") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Alignment") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TabStop: A tab stop within a paragraph.

func (TabStop) MarshalJSON ¶
func (s TabStop) MarshalJSON() ([]byte, error)
type Table ¶
type Table struct {
	// Columns: Number of columns in the table. It's possible for a table to be
	// non-rectangular, so some rows may have a different number of cells.
	Columns int64 `json:"columns,omitempty"`
	// Rows: Number of rows in the table.
	Rows int64 `json:"rows,omitempty"`
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. A Table may have
	// multiple insertion IDs if it's a nested suggested change. If empty, then
	// this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// TableRows: The contents and style of each row.
	TableRows []*TableRow `json:"tableRows,omitempty"`
	// TableStyle: The style of the table.
	TableStyle *TableStyle `json:"tableStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Columns") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Columns") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
Table: A StructuralElement representing a table.

func (Table) MarshalJSON ¶
func (s Table) MarshalJSON() ([]byte, error)
type TableCell ¶
type TableCell struct {
	// Content: The content of the cell.
	Content []*StructuralElement `json:"content,omitempty"`
	// EndIndex: The zero-based end index of this cell, exclusive, in UTF-16 code
	// units.
	EndIndex int64 `json:"endIndex,omitempty"`
	// StartIndex: The zero-based start index of this cell, in UTF-16 code units.
	StartIndex int64 `json:"startIndex,omitempty"`
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. A TableCell may have
	// multiple insertion IDs if it's a nested suggested change. If empty, then
	// this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTableCellStyleChanges: The suggested changes to the table cell
	// style, keyed by suggestion ID.
	SuggestedTableCellStyleChanges map[string]SuggestedTableCellStyle `json:"suggestedTableCellStyleChanges,omitempty"`
	// TableCellStyle: The style of the cell.
	TableCellStyle *TableCellStyle `json:"tableCellStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Content") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Content") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableCell: The contents and style of a cell in a Table.

func (TableCell) MarshalJSON ¶
func (s TableCell) MarshalJSON() ([]byte, error)
type TableCellBorder ¶
type TableCellBorder struct {
	// Color: The color of the border. This color cannot be transparent.
	Color *OptionalColor `json:"color,omitempty"`
	// DashStyle: The dash style of the border.
	//
	// Possible values:
	//   "DASH_STYLE_UNSPECIFIED" - Unspecified dash style.
	//   "SOLID" - Solid line. Corresponds to ECMA-376 ST_PresetLineDashVal value
	// 'solid'. This is the default dash style.
	//   "DOT" - Dotted line. Corresponds to ECMA-376 ST_PresetLineDashVal value
	// 'dot'.
	//   "DASH" - Dashed line. Corresponds to ECMA-376 ST_PresetLineDashVal value
	// 'dash'.
	DashStyle string `json:"dashStyle,omitempty"`
	// Width: The width of the border.
	Width *Dimension `json:"width,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Color") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Color") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableCellBorder: A border around a table cell. Table cell borders cannot be transparent. To hide a table cell border, make its width 0.

func (TableCellBorder) MarshalJSON ¶
func (s TableCellBorder) MarshalJSON() ([]byte, error)
type TableCellLocation ¶
type TableCellLocation struct {
	// ColumnIndex: The zero-based column index. For example, the second column in
	// the table has a column index of 1.
	ColumnIndex int64 `json:"columnIndex,omitempty"`
	// RowIndex: The zero-based row index. For example, the second row in the table
	// has a row index of 1.
	RowIndex int64 `json:"rowIndex,omitempty"`
	// TableStartLocation: The location where the table starts in the document.
	TableStartLocation *Location `json:"tableStartLocation,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ColumnIndex") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ColumnIndex") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableCellLocation: Location of a single cell within a table.

func (TableCellLocation) MarshalJSON ¶
func (s TableCellLocation) MarshalJSON() ([]byte, error)
type TableCellStyle ¶
type TableCellStyle struct {
	// BackgroundColor: The background color of the cell.
	BackgroundColor *OptionalColor `json:"backgroundColor,omitempty"`
	// BorderBottom: The bottom border of the cell.
	BorderBottom *TableCellBorder `json:"borderBottom,omitempty"`
	// BorderLeft: The left border of the cell.
	BorderLeft *TableCellBorder `json:"borderLeft,omitempty"`
	// BorderRight: The right border of the cell.
	BorderRight *TableCellBorder `json:"borderRight,omitempty"`
	// BorderTop: The top border of the cell.
	BorderTop *TableCellBorder `json:"borderTop,omitempty"`
	// ColumnSpan: The column span of the cell. This property is read-only.
	ColumnSpan int64 `json:"columnSpan,omitempty"`
	// ContentAlignment: The alignment of the content in the table cell. The
	// default alignment matches the alignment for newly created table cells in the
	// Docs editor.
	//
	// Possible values:
	//   "CONTENT_ALIGNMENT_UNSPECIFIED" - An unspecified content alignment. The
	// content alignment is inherited from the parent if one exists.
	//   "CONTENT_ALIGNMENT_UNSUPPORTED" - An unsupported content alignment.
	//   "TOP" - An alignment that aligns the content to the top of the content
	// holder. Corresponds to ECMA-376 ST_TextAnchoringType 't'.
	//   "MIDDLE" - An alignment that aligns the content to the middle of the
	// content holder. Corresponds to ECMA-376 ST_TextAnchoringType 'ctr'.
	//   "BOTTOM" - An alignment that aligns the content to the bottom of the
	// content holder. Corresponds to ECMA-376 ST_TextAnchoringType 'b'.
	ContentAlignment string `json:"contentAlignment,omitempty"`
	// PaddingBottom: The bottom padding of the cell.
	PaddingBottom *Dimension `json:"paddingBottom,omitempty"`
	// PaddingLeft: The left padding of the cell.
	PaddingLeft *Dimension `json:"paddingLeft,omitempty"`
	// PaddingRight: The right padding of the cell.
	PaddingRight *Dimension `json:"paddingRight,omitempty"`
	// PaddingTop: The top padding of the cell.
	PaddingTop *Dimension `json:"paddingTop,omitempty"`
	// RowSpan: The row span of the cell. This property is read-only.
	RowSpan int64 `json:"rowSpan,omitempty"`
	// ForceSendFields is a list of field names (e.g. "BackgroundColor") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "BackgroundColor") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableCellStyle: The style of a TableCell. Inherited table cell styles are represented as unset fields in this message. A table cell style can inherit from the table's style.

func (TableCellStyle) MarshalJSON ¶
func (s TableCellStyle) MarshalJSON() ([]byte, error)
type TableCellStyleSuggestionState ¶
type TableCellStyleSuggestionState struct {
	// BackgroundColorSuggested: Indicates if there was a suggested change to
	// background_color.
	BackgroundColorSuggested bool `json:"backgroundColorSuggested,omitempty"`
	// BorderBottomSuggested: Indicates if there was a suggested change to
	// border_bottom.
	BorderBottomSuggested bool `json:"borderBottomSuggested,omitempty"`
	// BorderLeftSuggested: Indicates if there was a suggested change to
	// border_left.
	BorderLeftSuggested bool `json:"borderLeftSuggested,omitempty"`
	// BorderRightSuggested: Indicates if there was a suggested change to
	// border_right.
	BorderRightSuggested bool `json:"borderRightSuggested,omitempty"`
	// BorderTopSuggested: Indicates if there was a suggested change to border_top.
	BorderTopSuggested bool `json:"borderTopSuggested,omitempty"`
	// ColumnSpanSuggested: Indicates if there was a suggested change to
	// column_span.
	ColumnSpanSuggested bool `json:"columnSpanSuggested,omitempty"`
	// ContentAlignmentSuggested: Indicates if there was a suggested change to
	// content_alignment.
	ContentAlignmentSuggested bool `json:"contentAlignmentSuggested,omitempty"`
	// PaddingBottomSuggested: Indicates if there was a suggested change to
	// padding_bottom.
	PaddingBottomSuggested bool `json:"paddingBottomSuggested,omitempty"`
	// PaddingLeftSuggested: Indicates if there was a suggested change to
	// padding_left.
	PaddingLeftSuggested bool `json:"paddingLeftSuggested,omitempty"`
	// PaddingRightSuggested: Indicates if there was a suggested change to
	// padding_right.
	PaddingRightSuggested bool `json:"paddingRightSuggested,omitempty"`
	// PaddingTopSuggested: Indicates if there was a suggested change to
	// padding_top.
	PaddingTopSuggested bool `json:"paddingTopSuggested,omitempty"`
	// RowSpanSuggested: Indicates if there was a suggested change to row_span.
	RowSpanSuggested bool `json:"rowSpanSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "BackgroundColorSuggested")
	// to unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "BackgroundColorSuggested") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableCellStyleSuggestionState: A mask that indicates which of the fields on the base TableCellStyle have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (TableCellStyleSuggestionState) MarshalJSON ¶
func (s TableCellStyleSuggestionState) MarshalJSON() ([]byte, error)
type TableColumnProperties ¶
type TableColumnProperties struct {
	// Width: The width of the column. Set when the column's `width_type` is
	// FIXED_WIDTH.
	Width *Dimension `json:"width,omitempty"`
	// WidthType: The width type of the column.
	//
	// Possible values:
	//   "WIDTH_TYPE_UNSPECIFIED" - The column width type is unspecified.
	//   "EVENLY_DISTRIBUTED" - The column width is evenly distributed among the
	// other evenly distributed columns. The width of the column is automatically
	// determined and will have an equal portion of the width remaining for the
	// table after accounting for all columns with specified widths.
	//   "FIXED_WIDTH" - A fixed column width. The width property contains the
	// column's width.
	WidthType string `json:"widthType,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Width") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Width") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableColumnProperties: The properties of a column in a table.

func (TableColumnProperties) MarshalJSON ¶
func (s TableColumnProperties) MarshalJSON() ([]byte, error)
type TableOfContents ¶
type TableOfContents struct {
	// Content: The content of the table of contents.
	Content []*StructuralElement `json:"content,omitempty"`
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. A TableOfContents may
	// have multiple insertion IDs if it is a nested suggested change. If empty,
	// then this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Content") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Content") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableOfContents: A StructuralElement representing a table of contents.

func (TableOfContents) MarshalJSON ¶
func (s TableOfContents) MarshalJSON() ([]byte, error)
type TableRange ¶
added in v0.8.0
type TableRange struct {
	// ColumnSpan: The column span of the table range.
	ColumnSpan int64 `json:"columnSpan,omitempty"`
	// RowSpan: The row span of the table range.
	RowSpan int64 `json:"rowSpan,omitempty"`
	// TableCellLocation: The cell location where the table range starts.
	TableCellLocation *TableCellLocation `json:"tableCellLocation,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ColumnSpan") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ColumnSpan") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableRange: A table range represents a reference to a subset of a table. It's important to note that the cells specified by a table range do not necessarily form a rectangle. For example, let's say we have a 3 x 3 table where all the cells of the last row are merged together. The table looks like this: [ ] A table range with table cell location = (table_start_location, row = 0, column = 0), row span = 3 and column span = 2 specifies the following cells: x x [ x x x ]

func (TableRange) MarshalJSON ¶
added in v0.8.0
func (s TableRange) MarshalJSON() ([]byte, error)
type TableRow ¶
type TableRow struct {
	// EndIndex: The zero-based end index of this row, exclusive, in UTF-16 code
	// units.
	EndIndex int64 `json:"endIndex,omitempty"`
	// StartIndex: The zero-based start index of this row, in UTF-16 code units.
	StartIndex int64 `json:"startIndex,omitempty"`
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. A TableRow may have
	// multiple insertion IDs if it's a nested suggested change. If empty, then
	// this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTableRowStyleChanges: The suggested style changes to this row,
	// keyed by suggestion ID.
	SuggestedTableRowStyleChanges map[string]SuggestedTableRowStyle `json:"suggestedTableRowStyleChanges,omitempty"`
	// TableCells: The contents and style of each cell in this row. It's possible
	// for a table to be non-rectangular, so some rows may have a different number
	// of cells than other rows in the same table.
	TableCells []*TableCell `json:"tableCells,omitempty"`
	// TableRowStyle: The style of the table row.
	TableRowStyle *TableRowStyle `json:"tableRowStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "EndIndex") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "EndIndex") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableRow: The contents and style of a row in a Table.

func (TableRow) MarshalJSON ¶
func (s TableRow) MarshalJSON() ([]byte, error)
type TableRowStyle ¶
type TableRowStyle struct {
	// MinRowHeight: The minimum height of the row. The row will be rendered in the
	// Docs editor at a height equal to or greater than this value in order to show
	// all the content in the row's cells.
	MinRowHeight *Dimension `json:"minRowHeight,omitempty"`
	// PreventOverflow: Whether the row cannot overflow across page or column
	// boundaries.
	PreventOverflow bool `json:"preventOverflow,omitempty"`
	// TableHeader: Whether the row is a table header.
	TableHeader bool `json:"tableHeader,omitempty"`
	// ForceSendFields is a list of field names (e.g. "MinRowHeight") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "MinRowHeight") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableRowStyle: Styles that apply to a table row.

func (TableRowStyle) MarshalJSON ¶
func (s TableRowStyle) MarshalJSON() ([]byte, error)
type TableRowStyleSuggestionState ¶
type TableRowStyleSuggestionState struct {
	// MinRowHeightSuggested: Indicates if there was a suggested change to
	// min_row_height.
	MinRowHeightSuggested bool `json:"minRowHeightSuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "MinRowHeightSuggested") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "MinRowHeightSuggested") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableRowStyleSuggestionState: A mask that indicates which of the fields on the base TableRowStyle have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (TableRowStyleSuggestionState) MarshalJSON ¶
func (s TableRowStyleSuggestionState) MarshalJSON() ([]byte, error)
type TableStyle ¶
type TableStyle struct {
	// TableColumnProperties: The properties of each column. Note that in Docs,
	// tables contain rows and rows contain cells, similar to HTML. So the
	// properties for a row can be found on the row's table_row_style.
	TableColumnProperties []*TableColumnProperties `json:"tableColumnProperties,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TableColumnProperties") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TableColumnProperties") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TableStyle: Styles that apply to a table.

func (TableStyle) MarshalJSON ¶
func (s TableStyle) MarshalJSON() ([]byte, error)
type TabsCriteria ¶
added in v0.190.0
type TabsCriteria struct {
	// TabIds: The list of tab IDs in which the request executes.
	TabIds []string `json:"tabIds,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TabIds") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TabIds") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TabsCriteria: A criteria that specifies in which tabs a request executes.

func (TabsCriteria) MarshalJSON ¶
added in v0.190.0
func (s TabsCriteria) MarshalJSON() ([]byte, error)
type TextRun ¶
type TextRun struct {
	// Content: The text of this run. Any non-text elements in the run are replaced
	// with the Unicode character U+E907.
	Content string `json:"content,omitempty"`
	// SuggestedDeletionIds: The suggested deletion IDs. If empty, then there are
	// no suggested deletions of this content.
	SuggestedDeletionIds []string `json:"suggestedDeletionIds,omitempty"`
	// SuggestedInsertionIds: The suggested insertion IDs. A TextRun may have
	// multiple insertion IDs if it's a nested suggested change. If empty, then
	// this is not a suggested insertion.
	SuggestedInsertionIds []string `json:"suggestedInsertionIds,omitempty"`
	// SuggestedTextStyleChanges: The suggested text style changes to this run,
	// keyed by suggestion ID.
	SuggestedTextStyleChanges map[string]SuggestedTextStyle `json:"suggestedTextStyleChanges,omitempty"`
	// TextStyle: The text style of this run.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Content") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Content") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TextRun: A ParagraphElement that represents a run of text that all has the same styling.

func (TextRun) MarshalJSON ¶
func (s TextRun) MarshalJSON() ([]byte, error)
type TextStyle ¶
type TextStyle struct {
	// BackgroundColor: The background color of the text. If set, the color is
	// either an RGB color or transparent, depending on the `color` field.
	BackgroundColor *OptionalColor `json:"backgroundColor,omitempty"`
	// BaselineOffset: The text's vertical offset from its normal position. Text
	// with `SUPERSCRIPT` or `SUBSCRIPT` baseline offsets is automatically rendered
	// in a smaller font size, computed based on the `font_size` field. Changes in
	// this field don't affect the `font_size`.
	//
	// Possible values:
	//   "BASELINE_OFFSET_UNSPECIFIED" - The text's baseline offset is inherited
	// from the parent.
	//   "NONE" - The text is not vertically offset.
	//   "SUPERSCRIPT" - The text is vertically offset upwards (superscript).
	//   "SUBSCRIPT" - The text is vertically offset downwards (subscript).
	BaselineOffset string `json:"baselineOffset,omitempty"`
	// Bold: Whether or not the text is rendered as bold.
	Bold bool `json:"bold,omitempty"`
	// FontSize: The size of the text's font.
	FontSize *Dimension `json:"fontSize,omitempty"`
	// ForegroundColor: The foreground color of the text. If set, the color is
	// either an RGB color or transparent, depending on the `color` field.
	ForegroundColor *OptionalColor `json:"foregroundColor,omitempty"`
	// Italic: Whether or not the text is italicized.
	Italic bool `json:"italic,omitempty"`
	// Link: The hyperlink destination of the text. If unset, there's no link.
	// Links are not inherited from parent text. Changing the link in an update
	// request causes some other changes to the text style of the range: * When
	// setting a link, the text foreground color will be updated to the default
	// link color and the text will be underlined. If these fields are modified in
	// the same request, those values will be used instead of the link defaults. *
	// Setting a link on a text range that overlaps with an existing link will also
	// update the existing link to point to the new URL. * Links are not settable
	// on newline characters. As a result, setting a link on a text range that
	// crosses a paragraph boundary, such as "ABC\n123", will separate the
	// newline character(s) into their own text runs. The link will be applied
	// separately to the runs before and after the newline. * Removing a link will
	// update the text style of the range to match the style of the preceding text
	// (or the default text styles if the preceding text is another link) unless
	// different styles are being set in the same request.
	Link *Link `json:"link,omitempty"`
	// SmallCaps: Whether or not the text is in small capital letters.
	SmallCaps bool `json:"smallCaps,omitempty"`
	// Strikethrough: Whether or not the text is struck through.
	Strikethrough bool `json:"strikethrough,omitempty"`
	// Underline: Whether or not the text is underlined.
	Underline bool `json:"underline,omitempty"`
	// WeightedFontFamily: The font family and rendered weight of the text. If an
	// update request specifies values for both `weighted_font_family` and `bold`,
	// the `weighted_font_family` is applied first, then `bold`. If
	// `weighted_font_family#weight` is not set, it defaults to `400`. If
	// `weighted_font_family` is set, then `weighted_font_family#font_family` must
	// also be set with a non-empty value. Otherwise, a 400 bad request error is
	// returned.
	WeightedFontFamily *WeightedFontFamily `json:"weightedFontFamily,omitempty"`
	// ForceSendFields is a list of field names (e.g. "BackgroundColor") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "BackgroundColor") to include in
	// API requests with the JSON null value. By default, fields with empty values
	// are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TextStyle: Represents the styling that can be applied to text. Inherited text styles are represented as unset fields in this message. A text style's parent depends on where the text style is defined: * The TextStyle of text in a Paragraph inherits from the paragraph's corresponding named style type. * The TextStyle on a named style inherits from the normal text named style. * The TextStyle of the normal text named style inherits from the default text style in the Docs editor. * The TextStyle on a Paragraph element that's contained in a table may inherit its text style from the table style. If the text style does not inherit from a parent, unsetting fields will revert the style to a value matching the defaults in the Docs editor.

func (TextStyle) MarshalJSON ¶
func (s TextStyle) MarshalJSON() ([]byte, error)
type TextStyleSuggestionState ¶
type TextStyleSuggestionState struct {
	// BackgroundColorSuggested: Indicates if there was a suggested change to
	// background_color.
	BackgroundColorSuggested bool `json:"backgroundColorSuggested,omitempty"`
	// BaselineOffsetSuggested: Indicates if there was a suggested change to
	// baseline_offset.
	BaselineOffsetSuggested bool `json:"baselineOffsetSuggested,omitempty"`
	// BoldSuggested: Indicates if there was a suggested change to bold.
	BoldSuggested bool `json:"boldSuggested,omitempty"`
	// FontSizeSuggested: Indicates if there was a suggested change to font_size.
	FontSizeSuggested bool `json:"fontSizeSuggested,omitempty"`
	// ForegroundColorSuggested: Indicates if there was a suggested change to
	// foreground_color.
	ForegroundColorSuggested bool `json:"foregroundColorSuggested,omitempty"`
	// ItalicSuggested: Indicates if there was a suggested change to italic.
	ItalicSuggested bool `json:"italicSuggested,omitempty"`
	// LinkSuggested: Indicates if there was a suggested change to link.
	LinkSuggested bool `json:"linkSuggested,omitempty"`
	// SmallCapsSuggested: Indicates if there was a suggested change to small_caps.
	SmallCapsSuggested bool `json:"smallCapsSuggested,omitempty"`
	// StrikethroughSuggested: Indicates if there was a suggested change to
	// strikethrough.
	StrikethroughSuggested bool `json:"strikethroughSuggested,omitempty"`
	// UnderlineSuggested: Indicates if there was a suggested change to underline.
	UnderlineSuggested bool `json:"underlineSuggested,omitempty"`
	// WeightedFontFamilySuggested: Indicates if there was a suggested change to
	// weighted_font_family.
	WeightedFontFamilySuggested bool `json:"weightedFontFamilySuggested,omitempty"`
	// ForceSendFields is a list of field names (e.g. "BackgroundColorSuggested")
	// to unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "BackgroundColorSuggested") to
	// include in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
TextStyleSuggestionState: A mask that indicates which of the fields on the base TextStyle have been changed in this suggestion. For any field set to true, there's a new suggested value.

func (TextStyleSuggestionState) MarshalJSON ¶
func (s TextStyleSuggestionState) MarshalJSON() ([]byte, error)
type UnmergeTableCellsRequest ¶
added in v0.10.0
type UnmergeTableCellsRequest struct {
	// TableRange: The table range specifying which cells of the table to unmerge.
	// All merged cells in this range will be unmerged, and cells that are already
	// unmerged will not be affected. If the range has no merged cells, the request
	// will do nothing. If there is text in any of the merged cells, the text will
	// remain in the "head" cell of the resulting block of unmerged cells. The
	// "head" cell is the upper-left cell when the content direction is from left
	// to right, and the upper-right otherwise.
	TableRange *TableRange `json:"tableRange,omitempty"`
	// ForceSendFields is a list of field names (e.g. "TableRange") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "TableRange") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UnmergeTableCellsRequest: Unmerges cells in a Table.

func (UnmergeTableCellsRequest) MarshalJSON ¶
added in v0.10.0
func (s UnmergeTableCellsRequest) MarshalJSON() ([]byte, error)
type UpdateCommentPostRequest ¶
type UpdateCommentPostRequest struct {
	// CommentId: The ID of the CommentThread which the post belongs to.
	CommentId string `json:"commentId,omitempty"`
	// Content: The new text of the comment, as plain text. This text content will
	// be handled similarly to comments created in the Docs editor. It will have
	// similar behaviors for formatting, notifications, etc. This field cannot be
	// empty, and must not exceed 2048 UTF-8 code units.
	Content string `json:"content,omitempty"`
	// PostId: The ID of the post being updated.
	PostId string `json:"postId,omitempty"`
	// SuggestionId: The ID of the SuggestionThread which the post belongs to.
	SuggestionId string `json:"suggestionId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "CommentId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "CommentId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateCommentPostRequest: Updates a Post in a CommentThread or SuggestionThread. Returns a 400 bad request error if: - The post is the headPost of a SuggestionThread. - The requesting user is not the author of the post.

func (UpdateCommentPostRequest) MarshalJSON ¶
func (s UpdateCommentPostRequest) MarshalJSON() ([]byte, error)
type UpdateDocumentStyleRequest ¶
added in v0.8.0
type UpdateDocumentStyleRequest struct {
	// DocumentStyle: The styles to set on the document. Certain document style
	// changes may cause other changes in order to mirror the behavior of the Docs
	// editor. See the documentation of DocumentStyle for more information.
	DocumentStyle *DocumentStyle `json:"documentStyle,omitempty"`
	// Fields: The fields that should be updated. At least one field must be
	// specified. The root `document_style` is implied and should not be specified.
	// A single "*" can be used as short-hand for listing every field. For
	// example to update the background, set `fields` to "background".
	Fields string `json:"fields,omitempty"`
	// TabId: The tab that contains the style to update. When omitted, the request
	// applies to the first tab. In a document containing a single tab: - If
	// provided, must match the singular tab's ID. - If omitted, the request
	// applies to the singular tab. In a document containing multiple tabs: - If
	// provided, the request applies to the specified tab. - If not provided, the
	// request applies to the first tab in the document.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DocumentStyle") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DocumentStyle") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateDocumentStyleRequest: Updates the DocumentStyle.

func (UpdateDocumentStyleRequest) MarshalJSON ¶
added in v0.8.0
func (s UpdateDocumentStyleRequest) MarshalJSON() ([]byte, error)
type UpdateDocumentTabPropertiesRequest ¶
added in v0.260.0
type UpdateDocumentTabPropertiesRequest struct {
	// Fields: The fields that should be updated. At least one field must be
	// specified. The root `tab_properties` is implied and should not be specified.
	// A single "*" can be used as short-hand for listing every field.
	Fields string `json:"fields,omitempty"`
	// TabProperties: The tab properties to update.
	TabProperties *TabProperties `json:"tabProperties,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Fields") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Fields") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateDocumentTabPropertiesRequest: Update the properties of a document tab.

func (UpdateDocumentTabPropertiesRequest) MarshalJSON ¶
added in v0.260.0
func (s UpdateDocumentTabPropertiesRequest) MarshalJSON() ([]byte, error)
type UpdateDropdownDefinitionPropertiesRequest ¶
type UpdateDropdownDefinitionPropertiesRequest struct {
	// DropdownDefinitionId: The ID of the DropdownDefinition to update.
	DropdownDefinitionId string `json:"dropdownDefinitionId,omitempty"`
	// DropdownDefinitionProperties: The properties to update.
	DropdownDefinitionProperties *DropdownDefinitionProperties `json:"dropdownDefinitionProperties,omitempty"`
	// Fields: The fields that should be updated. At least one field must be
	// specified. The root `dropdown_definition_properties` is implied and should
	// not be specified. A single "*" can be used as short-hand for listing every
	// field. When `dropdown_definition_properties.options` is included in the
	// field mask, the full, complete list of desired options must be provided in
	// `dropdown_definition_properties.options`.
	Fields string `json:"fields,omitempty"`
	// SelectedOptionIdReplacements: A map of option IDs to their replacements,
	// used to automatically reassign orphaned Dropdown chips when an option is
	// deleted. The keys are the IDs of the options being deleted, and the values
	// are the IDs of their replacement options. If an option being deleted is
	// selected in one or more Dropdown chips in the document, a replacement entry
	// for that option must be provided in this map, and the replacement option ID
	// must exist in the updated DropdownDefinition. If a replacement is required
	// but not provided, a 400 bad request error is returned. Options being deleted
	// that are not selected in any Dropdown chips do not require a replacement.
	// For example, if option A is being replaced by option B, the map should be
	// `{"A": "B"}`.
	SelectedOptionIdReplacements map[string]string `json:"selectedOptionIdReplacements,omitempty"`
	// TabId: The ID of the tab that contains the dropdown definition to update.
	// When omitted, the request is applied to the first tab. In a document
	// containing a single tab: - If provided, must match the singular tab's ID. -
	// If omitted, the request applies to the singular tab. In a document
	// containing multiple tabs: - If provided, the request applies to the
	// specified tab. - If omitted, the request applies to the first tab in the
	// document.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DropdownDefinitionId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DropdownDefinitionId") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateDropdownDefinitionPropertiesRequest: Updates the properties of a DropdownDefinition.

func (UpdateDropdownDefinitionPropertiesRequest) MarshalJSON ¶
func (s UpdateDropdownDefinitionPropertiesRequest) MarshalJSON() ([]byte, error)
type UpdateDropdownPropertiesRequest ¶
type UpdateDropdownPropertiesRequest struct {
	// DropdownId: Required. The Dropdown ID.
	DropdownId string `json:"dropdownId,omitempty"`
	// DropdownProperties: The properties to update.
	DropdownProperties *DropdownProperties `json:"dropdownProperties,omitempty"`
	// Fields: The fields that should be updated. At least one field must be
	// specified. The root `dropdown_properties` is implied and should not be
	// specified. A single "*" can be used as short-hand for listing every field.
	Fields string `json:"fields,omitempty"`
	// TabId: The ID of the tab that contains the dropdown to update. When omitted,
	// the request is applied to the first tab. In a document containing a single
	// tab: - If provided, must match the singular tab's ID. - If omitted, the
	// request applies to the singular tab. In a document containing multiple tabs:
	// - If provided, the request applies to the specified tab. - If omitted, the
	// request applies to the first tab in the document.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "DropdownId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "DropdownId") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateDropdownPropertiesRequest: Updates the properties of a Dropdown.

func (UpdateDropdownPropertiesRequest) MarshalJSON ¶
func (s UpdateDropdownPropertiesRequest) MarshalJSON() ([]byte, error)
type UpdateNamedStyleRequest ¶
added in v0.277.0
type UpdateNamedStyleRequest struct {
	// Fields: The NamedStyle fields that should be updated. At least
	// `named_style_type` must be specified. The root `named_style` is implied and
	// should not be specified. A single "*" can be used as short-hand for
	// listing every field. For example, to update the text style to bold, set
	// `fields` to include "text_style" and "text_style.bold". To update the
	// paragraph style's alignment property, set `fields` to include
	// "paragraph_style" and "paragraph_style.alignment". To reset a property
	// to its default value, include its field name in the field mask but leave the
	// field itself unset. Specifying "text_style" or "paragraph_style" with an
	// empty TextStyle or ParagraphStyle will reset all of its nested fields.
	Fields string `json:"fields,omitempty"`
	// NamedStyle: The document style to update.
	NamedStyle *NamedStyle `json:"namedStyle,omitempty"`
	// TabId: The document tab to update. By default, the update is applied to the
	// first tab.
	TabId string `json:"tabId,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Fields") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Fields") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateNamedStyleRequest: Updates a named style.

func (UpdateNamedStyleRequest) MarshalJSON ¶
added in v0.277.0
func (s UpdateNamedStyleRequest) MarshalJSON() ([]byte, error)
type UpdateParagraphStyleRequest ¶
type UpdateParagraphStyleRequest struct {
	// Fields: The fields that should be updated. At least one field must be
	// specified. The root `paragraph_style` is implied and should not be
	// specified. A single "*" can be used as short-hand for listing every field.
	// For example, to update the paragraph style's alignment property, set
	// `fields` to "alignment". To reset a property to its default value, include
	// its field name in the field mask but leave the field itself unset.
	Fields string `json:"fields,omitempty"`
	// ParagraphStyle: The styles to set on the paragraphs. Certain paragraph style
	// changes may cause other changes in order to mirror the behavior of the Docs
	// editor. See the documentation of ParagraphStyle for more information.
	ParagraphStyle *ParagraphStyle `json:"paragraphStyle,omitempty"`
	// Range: The range overlapping the paragraphs to style.
	Range *Range `json:"range,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Fields") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Fields") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateParagraphStyleRequest: Update the styling of all paragraphs that overlap with the given range.

func (UpdateParagraphStyleRequest) MarshalJSON ¶
func (s UpdateParagraphStyleRequest) MarshalJSON() ([]byte, error)
type UpdateSectionStyleRequest ¶
added in v0.15.0
type UpdateSectionStyleRequest struct {
	// Fields: The fields that should be updated. At least one field must be
	// specified. The root `section_style` is implied and must not be specified. A
	// single "*" can be used as short-hand for listing every field. For example
	// to update the left margin, set `fields` to "margin_left".
	Fields string `json:"fields,omitempty"`
	// Range: The range overlapping the sections to style. Because section breaks
	// can only be inserted inside the body, the segment ID field must be empty.
	Range *Range `json:"range,omitempty"`
	// SectionStyle: The styles to be set on the section. Certain section style
	// changes may cause other changes in order to mirror the behavior of the Docs
	// editor. See the documentation of SectionStyle for more information.
	SectionStyle *SectionStyle `json:"sectionStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Fields") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Fields") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateSectionStyleRequest: Updates the SectionStyle.

func (UpdateSectionStyleRequest) MarshalJSON ¶
added in v0.15.0
func (s UpdateSectionStyleRequest) MarshalJSON() ([]byte, error)
type UpdateTableCellStyleRequest ¶
added in v0.8.0
type UpdateTableCellStyleRequest struct {
	// Fields: The fields that should be updated. At least one field must be
	// specified. The root `tableCellStyle` is implied and should not be specified.
	// A single "*" can be used as short-hand for listing every field. For
	// example to update the table cell background color, set `fields` to
	// "backgroundColor". To reset a property to its default value, include its
	// field name in the field mask but leave the field itself unset.
	Fields string `json:"fields,omitempty"`
	// TableCellStyle: The style to set on the table cells. When updating borders,
	// if a cell shares a border with an adjacent cell, the corresponding border
	// property of the adjacent cell is updated as well. Borders that are merged
	// and invisible are not updated. Since updating a border shared by adjacent
	// cells in the same request can cause conflicting border updates, border
	// updates are applied in the following order: - `border_right` - `border_left`
	// - `border_bottom` - `border_top`
	TableCellStyle *TableCellStyle `json:"tableCellStyle,omitempty"`
	// TableRange: The table range representing the subset of the table to which
	// the updates are applied.
	TableRange *TableRange `json:"tableRange,omitempty"`
	// TableStartLocation: The location where the table starts in the document.
	// When specified, the updates are applied to all the cells in the table.
	TableStartLocation *Location `json:"tableStartLocation,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Fields") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Fields") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateTableCellStyleRequest: Updates the style of a range of table cells.

func (UpdateTableCellStyleRequest) MarshalJSON ¶
added in v0.8.0
func (s UpdateTableCellStyleRequest) MarshalJSON() ([]byte, error)
type UpdateTableColumnPropertiesRequest ¶
added in v0.5.0
type UpdateTableColumnPropertiesRequest struct {
	// ColumnIndices: The list of zero-based column indices whose property should
	// be updated. If no indices are specified, all columns will be updated.
	ColumnIndices []int64 `json:"columnIndices,omitempty"`
	// Fields: The fields that should be updated. At least one field must be
	// specified. The root `tableColumnProperties` is implied and should not be
	// specified. A single "*" can be used as short-hand for listing every field.
	// For example to update the column width, set `fields` to "width".
	Fields string `json:"fields,omitempty"`
	// TableColumnProperties: The table column properties to update. If the value
	// of `table_column_properties#width` is less than 5 points (5/72 inch), a 400
	// bad request error is returned.
	TableColumnProperties *TableColumnProperties `json:"tableColumnProperties,omitempty"`
	// TableStartLocation: The location where the table starts in the document.
	TableStartLocation *Location `json:"tableStartLocation,omitempty"`
	// ForceSendFields is a list of field names (e.g. "ColumnIndices") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "ColumnIndices") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateTableColumnPropertiesRequest: Updates the TableColumnProperties of columns in a table.

func (UpdateTableColumnPropertiesRequest) MarshalJSON ¶
added in v0.5.0
func (s UpdateTableColumnPropertiesRequest) MarshalJSON() ([]byte, error)
type UpdateTableRowStyleRequest ¶
added in v0.5.0
type UpdateTableRowStyleRequest struct {
	// Fields: The fields that should be updated. At least one field must be
	// specified. The root `tableRowStyle` is implied and should not be specified.
	// A single "*" can be used as short-hand for listing every field. For
	// example to update the minimum row height, set `fields` to
	// "min_row_height".
	Fields string `json:"fields,omitempty"`
	// RowIndices: The list of zero-based row indices whose style should be
	// updated. If no indices are specified, all rows will be updated.
	RowIndices []int64 `json:"rowIndices,omitempty"`
	// TableRowStyle: The styles to be set on the rows.
	TableRowStyle *TableRowStyle `json:"tableRowStyle,omitempty"`
	// TableStartLocation: The location where the table starts in the document.
	TableStartLocation *Location `json:"tableStartLocation,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Fields") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Fields") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateTableRowStyleRequest: Updates the TableRowStyle of rows in a table.

func (UpdateTableRowStyleRequest) MarshalJSON ¶
added in v0.5.0
func (s UpdateTableRowStyleRequest) MarshalJSON() ([]byte, error)
type UpdateTextStyleRequest ¶
type UpdateTextStyleRequest struct {
	// Fields: The fields that should be updated. At least one field must be
	// specified. The root `text_style` is implied and should not be specified. A
	// single "*" can be used as short-hand for listing every field. For example,
	// to update the text style to bold, set `fields` to "bold". To reset a
	// property to its default value, include its field name in the field mask but
	// leave the field itself unset.
	Fields string `json:"fields,omitempty"`
	// Range: The range of text to style. The range may be extended to include
	// adjacent newlines. If the range fully contains a paragraph belonging to a
	// list, the paragraph's bullet is also updated with the matching text style.
	// Ranges cannot be inserted inside a relative UpdateTextStyleRequest.
	Range *Range `json:"range,omitempty"`
	// TextStyle: The styles to set on the text. If the value for a particular
	// style matches that of the parent, that style will be set to inherit. Certain
	// text style changes may cause other changes in order to to mirror the
	// behavior of the Docs editor. See the documentation of TextStyle for more
	// information.
	TextStyle *TextStyle `json:"textStyle,omitempty"`
	// ForceSendFields is a list of field names (e.g. "Fields") to unconditionally
	// include in API requests. By default, fields with empty or default values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "Fields") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
UpdateTextStyleRequest: Update the styling of text.

func (UpdateTextStyleRequest) MarshalJSON ¶
func (s UpdateTextStyleRequest) MarshalJSON() ([]byte, error)
type WeightedFontFamily ¶
type WeightedFontFamily struct {
	// FontFamily: The font family of the text. The font family can be any font
	// from the Font menu in Docs or from [Google Fonts]
	// (https://fonts.google.com/). If the font name is unrecognized, the text is
	// rendered in `Arial`.
	FontFamily string `json:"fontFamily,omitempty"`
	// Weight: The weight of the font. This field can have any value that's a
	// multiple of `100` between `100` and `900`, inclusive. This range corresponds
	// to the numerical values described in the CSS 2.1 Specification, section 15.6
	// (https://www.w3.org/TR/CSS21/fonts.html#font-boldness), with non-numerical
	// values disallowed. The default value is `400` ("normal"). The font weight
	// makes up just one component of the rendered font weight. A combination of
	// the `weight` and the text style's resolved `bold` value determine the
	// rendered weight, after accounting for inheritance: * If the text is bold and
	// the weight is less than `400`, the rendered weight is 400. * If the text is
	// bold and the weight is greater than or equal to `400` but is less than
	// `700`, the rendered weight is `700`. * If the weight is greater than or
	// equal to `700`, the rendered weight is equal to the weight. * If the text is
	// not bold, the rendered weight is equal to the weight.
	Weight int64 `json:"weight,omitempty"`
	// ForceSendFields is a list of field names (e.g. "FontFamily") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "FontFamily") to include in API
	// requests with the JSON null value. By default, fields with empty values are
	// omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
WeightedFontFamily: Represents a font family and weight of text.

func (WeightedFontFamily) MarshalJSON ¶
func (s WeightedFontFamily) MarshalJSON() ([]byte, error)
type WriteControl ¶
type WriteControl struct {
	// RequiredRevisionId: The optional revision ID of the document the write
	// request is applied to. If this is not the latest revision of the document,
	// the request is not processed and returns a 400 bad request error. When a
	// required revision ID is returned in a response, it indicates the revision ID
	// of the document after the request was applied.
	RequiredRevisionId string `json:"requiredRevisionId,omitempty"`
	// TargetRevisionId: The optional target revision ID of the document the write
	// request is applied to. If collaborator changes have occurred after the
	// document was read using the API, the changes produced by this write request
	// are applied against the collaborator changes. This results in a new revision
	// of the document that incorporates both the collaborator changes and the
	// changes in the request, with the Docs server resolving conflicting changes.
	// When using target revision ID, the API client can be thought of as another
	// collaborator of the document. The target revision ID can only be used to
	// write to recent versions of a document. If the target revision is too far
	// behind the latest revision, the request is not processed and returns a 400
	// bad request error. The request should be tried again after retrieving the
	// latest version of the document. Usually a revision ID remains valid for use
	// as a target revision for several minutes after it's read, but for frequently
	// edited documents this window might be shorter.
	TargetRevisionId string `json:"targetRevisionId,omitempty"`
	// WriteMode: How the request updates should be applied to the document. If
	// unspecified, the request updates will be applied as normal edits.
	//
	// Possible values:
	//   "WRITE_MODE_UNSPECIFIED" - The write mode is unspecified. Defaults to EDIT
	// behavior.
	//   "EDIT" - Apply all updates as normal edits.
	//   "SUGGEST" - Apply all updates as suggestions.
	WriteMode string `json:"writeMode,omitempty"`
	// ForceSendFields is a list of field names (e.g. "RequiredRevisionId") to
	// unconditionally include in API requests. By default, fields with empty or
	// default values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-ForceSendFields for more
	// details.
	ForceSendFields []string `json:"-"`
	// NullFields is a list of field names (e.g. "RequiredRevisionId") to include
	// in API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. See
	// https://pkg.go.dev/google.golang.org/api#hdr-NullFields for more details.
	NullFields []string `json:"-"`
}
WriteControl: Provides control over how write requests are executed.

func (WriteControl) MarshalJSON ¶
func (s WriteControl) MarshalJSON() ([]byte, error)
 Source Files ¶
View all Source files
docs-gen.go
Why Go
Use Cases
Case Studies
Get Started
Playground
Tour
Stack Overflow
Help
Packages
Standard Library
Sub-repositories
About Go Packages
pkg.go.dev API
About
Download
Blog
Issue Tracker
Release Notes
Brand Guidelines
Code of Conduct
Connect
Twitter
GitHub
Slack
r/golang
Meetup
Golang Weekly
Gopher in flight goggles
Copyright
Terms of Service
Privacy Policy
Report an Issue
System theme
Theme Toggle


Shortcuts Modal

Google logo