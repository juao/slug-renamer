# Slug Renamer

> Renomeie todos os arquivos de uma pasta para o formato **slug** com um clique direito.

Clique com o botão direito dentro de qualquer pasta no Windows Explorer e escolha **"Renomear arquivos como slug"**.

Uma janela de preview mostra exatamente como os arquivos ficarão **antes** de qualquer mudança.

---

## Exemplos de conversão

| Antes | Depois |
|---|---|
| `Meu Arquivo (1).jpg` | `meu-arquivo-1.jpg` |
| `Foto da Viagem!.PNG` | `foto-da-viagem.png` |
| `relatório_2024.pdf` | `relatorio-2024.pdf` |
| `Ação & Aventura — Ep 3.mp4` | `acao-aventura-ep-3.mp4` |
| `TUDO_EM_MAIÚSCULAS.DOC` | `tudo-em-maiusculas.doc` |

---

## Instalação (usuário final)

1. Baixe `slug-renamer-setup.exe` na [página de releases](../../releases)
2. Execute o instalador (ele pede permissão de administrador — necessário para registrar o menu)
3. Pronto! Clique direito dentro de qualquer pasta e use a opção

Para desinstalar: Painel de Controle → Adicionar/Remover Programas → Slug Renamer

---

## Build (desenvolvedor)

### Pré-requisitos

| Ferramenta | Versão | Link |
|---|---|---|
| Go | 1.22+ | https://go.dev/dl/ |
| MinGW (windres) | qualquer | `choco install mingw` |
| NSIS | 3.x | `choco install nsis` |

> Dica: instale o [Chocolatey](https://chocolatey.org/) e rode `choco install golang mingw nsis` de uma vez.

### Compilar

```powershell
# Clone o repositório
git clone https://github.com/juao/slug-renamer
cd slug-renamer

# Instale as dependências Go
go mod download

# Build completo (exe + instalador)
.\build.ps1
```

Os artefatos ficam em `.\dist\`:
- `slug-renamer.exe` — executável standalone
- `slug-renamer-setup.exe` — instalador para distribuição

### Compilar só o exe (sem instalador)

```powershell
$env:GOOS="windows"; $env:GOARCH="amd64"
# Antes, gere o .syso para embutir manifesto + version info:
pushd cmd; windres slug-renamer.rc -O coff -o slug-renamer.syso; popd
go build -o dist\slug-renamer.exe -ldflags "-s -w -H windowsgui" .\cmd\...
```

### Rodar os testes

```powershell
go test ./internal/slugify/...
```

---

## Estrutura do projeto

```
slug-renamer/
├── cmd/
│   ├── main.go               # Entry point
│   ├── registry_windows.go   # Registro no Windows Registry
│   ├── manifest.xml          # Manifesto de admin
│   └── slug-renamer.rc       # Resource file (ícone + versão)
├── internal/
│   ├── slugify/
│   │   ├── slugify.go        # Lógica de conversão para slug
│   │   └── slugify_test.go   # Testes unitários
│   └── ui/
│       └── ui_windows.go     # Interface gráfica (walk)
├── installer/
│   └── installer.nsi         # Script NSIS
├── build.ps1                 # Script de build
├── go.mod
└── README.md
```

---

## Como funciona

1. O instalador copia o `slug-renamer.exe` para `%ProgramFiles%\SlugRenamer\`
2. Roda `slug-renamer.exe --install` que escreve em:
   ```
   HKEY_CLASSES_ROOT\Directory\Background\shell\SlugRenamer
   HKEY_CLASSES_ROOT\Directory\Background\shell\SlugRenamer\command
   ```
3. Quando o usuário clica com botão direito **dentro** de uma pasta, o Windows chama:
   ```
   slug-renamer.exe "C:\caminho\da\pasta"
   ```
4. O app lê os arquivos, gera o preview e executa o rename após confirmação

---

## Distribuição

Para publicar:

1. Faça o build com `.\build.ps1`
2. (Opcional mas recomendado) Assine o executável com um certificado de Code Signing para evitar avisos do Windows Defender
3. Publique o `dist\slug-renamer-setup.exe` no GitHub Releases

---

## Licença

MIT
