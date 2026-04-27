//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const (
	menuLabel   = "Renomear arquivos como slug"
	menuKeyPath = `Directory\Background\shell\SlugRenamer`
	cmdKeyPath  = `Directory\Background\shell\SlugRenamer\command`
)

// registerContextMenu registra a entrada no menu de contexto do Windows Explorer.
// Aparece ao clicar com botão direito DENTRO de uma pasta (no fundo/background).
func registerContextMenu() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("não foi possível obter o caminho do executável: %w", err)
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return err
	}

	// Cria a chave principal com o label do menu
	menuKey, _, err := registry.CreateKey(
		registry.CLASSES_ROOT,
		menuKeyPath,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar chave de registro: %w", err)
	}
	defer menuKey.Close()

	if err := menuKey.SetStringValue("", menuLabel); err != nil {
		return err
	}
	// Ícone opcional (usa o próprio exe como ícone)
	if err := menuKey.SetStringValue("Icon", exePath); err != nil {
		return err
	}

	// Cria a subchave \command com o executável e o placeholder da pasta
	cmdKey, _, err := registry.CreateKey(
		registry.CLASSES_ROOT,
		cmdKeyPath,
		registry.SET_VALUE,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar subchave command: %w", err)
	}
	defer cmdKey.Close()

	// %V = caminho completo da pasta atual no Explorer
	command := fmt.Sprintf(`"%s" "%%V"`, exePath)
	return cmdKey.SetStringValue("", command)
}

// unregisterContextMenu remove a entrada do Registry.
func unregisterContextMenu() error {
	err := registry.DeleteKey(registry.CLASSES_ROOT, cmdKeyPath)
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("erro ao remover subchave command: %w", err)
	}

	err = registry.DeleteKey(registry.CLASSES_ROOT, menuKeyPath)
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("erro ao remover chave principal: %w", err)
	}

	return nil
}
