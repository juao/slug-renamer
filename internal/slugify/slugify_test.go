package slugify_test

import (
	"testing"

	"github.com/juao/slug-renamer/internal/slugify"
)

func TestToSlug(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"Meu Arquivo (1).jpg", "meu-arquivo-1.jpg"},
		{"Foto da Viagem!.PNG", "foto-da-viagem.png"},
		{"relatório_2024.pdf", "relatorio-2024.pdf"},
		{"já estava em slug.txt", "ja-estava-em-slug.txt"},
		{"  espaços  nas  bordas  .md", "espacos-nas-bordas.md"},
		{"Ação & Aventura — Episódio 3.mp4", "acao-aventura-episodio-3.mp4"},
		{"arquivo.com.dupla.extensao.tar.gz", "arquivo-com-dupla-extensao-tar.gz"},
		{"TUDO_EM_MAIÚSCULAS.DOC", "tudo-em-maiusculas.doc"},
		{"já-esta-em-slug.txt", "ja-esta-em-slug.txt"},
		{"###$$$%%%", "arquivo"}, // nome vazio após limpeza
		{"noext", "noext"},
		{"múltiplos   espaços.jpg", "multiplos-espacos.jpg"},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			got := slugify.ToSlug(c.input)
			if got != c.expected {
				t.Errorf("\n  input:    %q\n  expected: %q\n  got:      %q", c.input, c.expected, got)
			}
		})
	}
}
