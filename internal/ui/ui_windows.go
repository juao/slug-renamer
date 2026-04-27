//go:build windows

package ui

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/juao/slug-renamer/internal/slugify"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procMessageBoxW = user32.NewProc("MessageBoxW")
)

const (
	mbOK              = 0x00000000
	mbOKCancel        = 0x00000001
	mbIconInformation = 0x00000040
	mbIconError       = 0x00000010
	mbIconQuestion    = 0x00000020
	mbDefButton2      = 0x00000100
	idOK              = 1
)

func messageBox(title, text string, flags uint32) int32 {
	t, _ := syscall.UTF16PtrFromString(title)
	m, _ := syscall.UTF16PtrFromString(text)
	ret, _, _ := procMessageBoxW.Call(
		0,
		uintptr(unsafe.Pointer(m)),
		uintptr(unsafe.Pointer(t)),
		uintptr(flags),
	)
	return int32(ret)
}

func ShowPreview(previews []slugify.RenamePreview) bool {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d arquivo(s) serao renomeados:\n\n", len(previews)))
	limit := len(previews)
	if limit > 20 {
		limit = 20
	}
	for _, p := range previews[:limit] {
		sb.WriteString(fmt.Sprintf("%s\n  -> %s\n\n", p.Original, p.New))
	}
	if len(previews) > 20 {
		sb.WriteString(fmt.Sprintf("... e mais %d arquivo(s).\n\n", len(previews)-20))
	}
	sb.WriteString("Clique OK para renomear ou Cancelar para abortar.")
	ret := messageBox("Slug Renamer - Preview", sb.String(), mbOKCancel|mbIconQuestion|mbDefButton2)
	return ret == idOK
}

func ShowError(msg string) {
	messageBox("Slug Renamer - Erro", msg, mbOK|mbIconError)
}

func ShowInfo(msg string) {
	messageBox("Slug Renamer", msg, mbOK|mbIconInformation)
}
