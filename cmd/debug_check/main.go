package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
)

func main() {
	credPath, _ := filepath.Abs("credentials.json")
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credPath)
	ctx := context.Background()

	docID := "1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8"
	opts := gdocs.LintOptions{
		CheckLatexLeaks:      true,
		CheckSpacingAndText:  true,
		CheckStructuralOrder: true,
	}

	report, err := gdocs.LintDocumentABNT(ctx, docID, false, false, opts)
	if err != nil {
		fmt.Printf("Erro: %v\n", err)
		return
	}

	fmt.Printf("LatexSyntaxLeaks count: %d\n", len(report.LatexSyntaxLeaks))
	for idx, leak := range report.LatexSyntaxLeaks {
		if idx < 10 {
			fmt.Printf("  [%d] Type: %s | Snippet: %q\n", idx, leak.LeakType, leak.Snippet)
		}
	}
}
