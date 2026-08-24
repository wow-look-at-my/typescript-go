package lsproto

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/json"
	"github.com/wow-look-at-my/testify/require"
)

func TestUnmarshalCompletionItem(t *testing.T) {
	t.Parallel()

	const message = `{
    "label": "pageXOffset",
    "insertTextFormat": 1,
    "textEdit": {
        "newText": "pageXOffset",
        "insert": {
            "start": {
                "line": 4,
                "character": 0
            },
            "end": {
                "line": 4,
                "character": 4
            }
        },
        "replace": {
            "start": {
                "line": 4,
                "character": 0
            },
            "end": {
                "line": 4,
                "character": 4
            }
        }
    },
    "kind": 6,
    "sortText": "15",
    "commitCharacters": [
        ".",
        ",",
        ";"
    ]
}`

	var result CompletionItem
	err := json.Unmarshal([]byte(message), &result)
	require.NoError(t, err)

	require.Equal(t, result, CompletionItem{
		Label:			"pageXOffset",
		InsertTextFormat:	new(InsertTextFormatPlainText),
		TextEdit: &TextEditOrInsertReplaceEdit{
			InsertReplaceEdit: &InsertReplaceEdit{
				NewText:	"pageXOffset",
				Insert: Range{
					Start: Position{
						Line:		4,
						Character:	0,
					},
					End: Position{
						Line:		4,
						Character:	4,
					},
				},
				Replace: Range{
					Start: Position{
						Line:		4,
						Character:	0,
					},
					End: Position{
						Line:		4,
						Character:	4,
					},
				},
			},
		},
		Kind:			new(CompletionItemKindVariable),
		SortText:		new("15"),
		CommitCharacters:	new([]string{".", ",", ";"}),
	})
}
