// Package gdocs fornece funções para autenticação e interação
// com as APIs do Google (Docs v1 e Drive v3) via Service Account.
package gdocs

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Scopes de acesso às APIs do Google.
// - documents: leitura e ESCRITA no Google Docs (necessário para replace_text_in_doc).
// - drive: leitura e ESCRITA no Google Drive (necessário para comentários).
var scopes = []string{
	"https://www.googleapis.com/auth/documents",
	"https://www.googleapis.com/auth/drive",
}

// Singletons thread-safe para os serviços Google.
var (
	docsService  *docs.Service
	driveService *drive.Service
	docsOnce     sync.Once
	driveOnce    sync.Once
	docsErr      error
	driveErr     error
)

// GetCredentialsPath retorna o caminho do arquivo de credenciais.
// Prioridade:
//  1. Variável de ambiente GOOGLE_APPLICATION_CREDENTIALS
//  2. Arquivo "credentials.json" na mesma pasta do executável
func GetCredentialsPath() string {
	// 1. Verifica variável de ambiente
	if envPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); envPath != "" {
		log.Printf("[INFO] Usando credenciais da variável de ambiente: %s", envPath)
		return envPath
	}

	// 2. Busca ao lado do executável
	exe, err := os.Executable()
	if err != nil {
		log.Printf("[WARN] Não foi possível determinar o caminho do executável: %v", err)
		return "credentials.json"
	}

	credPath := filepath.Join(filepath.Dir(exe), "credentials.json")
	log.Printf("[INFO] Buscando credenciais em: %s", credPath)
	return credPath
}

// loadCredentials carrega e valida as credenciais da Service Account.
func loadCredentials() (*google.Credentials, error) {
	credPath := GetCredentialsPath()

	data, err := os.ReadFile(credPath)
	if err != nil {
		return nil, fmt.Errorf(
			"arquivo de credenciais não encontrado em '%s'.\n"+
				"Solução: Defina a variável GOOGLE_APPLICATION_CREDENTIALS ou "+
				"coloque o arquivo 'credentials.json' na mesma pasta do executável",
			credPath,
		)
	}

	creds, err := google.CredentialsFromJSON(context.Background(), data, scopes...)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear credenciais JSON: %w", err)
	}

	log.Println("[INFO] Credenciais da Service Account carregadas com sucesso.")
	return creds, nil
}

// GetDocsService retorna o client da API Google Docs v1 (singleton).
func GetDocsService(ctx context.Context) (*docs.Service, error) {
	docsOnce.Do(func() {
		creds, err := loadCredentials()
		if err != nil {
			docsErr = err
			return
		}

		docsService, docsErr = docs.NewService(ctx, option.WithCredentials(creds), option.WithScopes(docs.DocumentsScope))
		if docsErr == nil {
			log.Println("[INFO] Google Docs service inicializado com sucesso.")
		}
	})
	return docsService, docsErr
}

// GetDriveService retorna o client da API Google Drive v3 (singleton).
func GetDriveService(ctx context.Context) (*drive.Service, error) {
	driveOnce.Do(func() {
		creds, err := loadCredentials()
		if err != nil {
			driveErr = err
			return
		}

		driveService, driveErr = drive.NewService(ctx, option.WithCredentials(creds), option.WithScopes(drive.DriveScope))
		if driveErr == nil {
			log.Println("[INFO] Google Drive service inicializado com sucesso.")
		}
	})
	return driveService, driveErr
}
