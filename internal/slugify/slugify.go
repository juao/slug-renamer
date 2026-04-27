package slugify

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"unicode"
)

// RenamePreview representa um arquivo e como ele ficará após o rename.
type RenamePreview struct {
	Original string
	New      string
	Path     string // pasta pai
}

var (
	// Remove acentos após normalização unicode
	nonSpacingMarks = runes.In(unicode.Mn)

	// Substitui qualquer sequência de caracteres não-alfanuméricos por hífen
	nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

	// Remove hífens no início e no fim
	leadingTrailingHyphen = regexp.MustCompile(`^-+|-+$`)

	// Múltiplos hífens consecutivos → um só
	multipleHyphens = regexp.MustCompile(`-{2,}`)
)

// ToSlug converte um nome de arquivo para o formato slug.
// Preserva a extensão do arquivo.
//
// Exemplos:
//   "Meu Arquivo (1).jpg"  → "meu-arquivo-1.jpg"
//   "Foto da Viagem!.PNG"  → "foto-da-viagem.png"
//   "relatório_2024.pdf"   → "relatorio-2024.pdf"
func ToSlug(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// 1. Normaliza unicode NFD (decompõe caracteres acentuados)
	t := transform.Chain(norm.NFD, runes.Remove(nonSpacingMarks), norm.NFC)
	normalized, _, err := transform.String(t, name)
	if err != nil {
		normalized = name
	}

	// 2. Minúsculas
	lower := strings.ToLower(normalized)

	// 3. Substitui underscores e espaços por hífen explicitamente
	lower = strings.ReplaceAll(lower, "_", "-")
	lower = strings.ReplaceAll(lower, " ", "-")

	// 4. Remove caracteres não-alfanuméricos (exceto hífen)
	slug := nonAlphanumeric.ReplaceAllString(lower, "-")

	// 5. Limpa hífens duplicados e nas bordas
	slug = multipleHyphens.ReplaceAllString(slug, "-")
	slug = leadingTrailingHyphen.ReplaceAllString(slug, "")

	if slug == "" {
		slug = "arquivo"
	}

	return slug + ext
}

// RenameFiles executa os renames e retorna uma lista de erros (se houver).
func RenameFiles(previews []RenamePreview) []string {
	var errors []string

	for _, p := range previews {
		oldPath := filepath.Join(p.Path, p.Original)
		newPath := filepath.Join(p.Path, p.New)

		// Evita sobrescrever arquivo existente
		if _, err := os.Stat(newPath); err == nil {
			errors = append(errors, p.Original+" → já existe: "+p.New)
			continue
		}

		if err := os.Rename(oldPath, newPath); err != nil {
			errors = append(errors, p.Original+": "+err.Error())
		}
	}

	return errors
}
