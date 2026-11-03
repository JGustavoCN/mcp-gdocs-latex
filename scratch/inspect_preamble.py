import json
from google.oauth2 import service_account
from googleapiclient.discovery import build

creds = service_account.Credentials.from_service_account_file('credentials.json')
service = build('docs', 'v1', credentials=creds)
doc = service.documents().get(documentId='1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8').execute()

for idx, el in enumerate(doc['body']['content']):
    p = el.get('paragraph')
    if p:
        text = ''.join(e.get('textRun', {}).get('content', '') for e in p.get('elements', []))
        if any(k in text for k in ['Trabalho de Conclusão', 'Orientador:', 'LAGARTO', '2026']):
            print(f"Element {idx} [{el['startIndex']}-{el['endIndex']}]:")
            print("  Text:", repr(text.strip()[:60]))
            print("  Style:", json.dumps(p.get('paragraphStyle', {}), indent=4))

