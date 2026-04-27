package main

import (
	"fmt"
	"os"

	"github.com/juao/slug-renamer/internal/slugify"
	"github.com/juao/slug-renamer/internal/ui"
)

func main() {
	args := os.Args[1:]

	// Modo: --install ou --uninstall (chamado pelo instalador)
	if len(args) == 1 {
		switch args[0] {
		case "--install":
			if err := registerContextMenu(); err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao instalar: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Instalado com sucesso!")
			return
		case "--uninstall":
			if err := unregisterContextMenu(); err != nil {
				fmt.Fprintf(os.Stderr, "Erro ao desinstalar: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Desinstalado com sucesso!")
			return
		}
	}

	// Modo normal: chamado pelo menu de contexto com o caminho da pasta
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Uso: slug-renamer.exe <caminho-da-pasta>")
		os.Exit(1)
	}

	folderPath := args[0]

	// Verifica se é uma pasta válida
	info, err := os.Stat(folderPath)
	if err != nil || !info.IsDir() {
		ui.ShowError("Caminho inválido ou não é uma pasta:\n" + folderPath)
		os.Exit(1)
	}

	// Lê os arquivos da pasta
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		ui.ShowError("Erro ao ler a pasta:\n" + err.Error())
		os.Exit(1)
	}

	// Monta preview dos renames
	var previews []slugify.RenamePreview
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		original := entry.Name()
		renamed := slugify.ToSlug(original)
		if original != renamed {
			previews = append(previews, slugify.RenamePreview{
				Original: original,
				New:      renamed,
				Path:     folderPath,
			})
		}
	}

	if len(previews) == 0 {
		ui.ShowInfo("Todos os arquivos já estão no formato slug.\nNada a renomear.")
		return
	}

	// Exibe UI de preview e aguarda confirmação
	confirmed := ui.ShowPreview(previews)
	if !confirmed {
		return
	}

	// Executa o rename
	errors := slugify.RenameFiles(previews)
	if len(errors) > 0 {
		msg := fmt.Sprintf("Concluído com %d erro(s):\n", len(errors))
		for _, e := range errors {
			msg += "• " + e + "\n"
		}
		ui.ShowError(msg)
	} else {
		ui.ShowInfo(fmt.Sprintf("%d arquivo(s) renomeado(s) com sucesso!", len(previews)))
	}
}
