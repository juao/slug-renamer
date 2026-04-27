Unicode True

!define APP_NAME      "Slug Renamer"
!define APP_VERSION   "1.0.0"
!define APP_EXE       "slug-renamer.exe"
!define INSTALL_DIR   "$PROGRAMFILES64\SlugRenamer"
!define UNINSTALL_KEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\SlugRenamer"

Name "${APP_NAME} ${APP_VERSION}"
OutFile "slug-renamer-setup.exe"
InstallDir "${INSTALL_DIR}"
RequestExecutionLevel admin
ShowInstDetails show
ShowUninstDetails show

Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Section "Instalar ${APP_NAME}" SecInstall
  SetOutPath "$INSTDIR"
  File "slug-renamer.exe"
  ExecWait '"$INSTDIR\${APP_EXE}" --install' $0
  WriteRegStr HKLM "${UNINSTALL_KEY}" "DisplayName"     "${APP_NAME}"
  WriteRegStr HKLM "${UNINSTALL_KEY}" "DisplayVersion"  "${APP_VERSION}"
  WriteRegStr HKLM "${UNINSTALL_KEY}" "Publisher"       "Seu Nome"
  WriteRegStr HKLM "${UNINSTALL_KEY}" "InstallLocation" "$INSTDIR"
  WriteRegStr HKLM "${UNINSTALL_KEY}" "UninstallString" '"$INSTDIR\uninstall.exe"'
  WriteRegStr HKLM "${UNINSTALL_KEY}" "DisplayIcon"     "$INSTDIR\${APP_EXE}"
  WriteRegDWORD HKLM "${UNINSTALL_KEY}" "NoModify"      1
  WriteRegDWORD HKLM "${UNINSTALL_KEY}" "NoRepair"      1
  WriteUninstaller "$INSTDIR\uninstall.exe"
  MessageBox MB_ICONINFORMATION "Instalacao concluida! Clique com botao direito dentro de uma pasta e use Renomear arquivos como slug."
SectionEnd

Section "Uninstall"
  ExecWait '"$INSTDIR\${APP_EXE}" --uninstall'
  DeleteRegKey HKLM "${UNINSTALL_KEY}"
  Delete "$INSTDIR\${APP_EXE}"
  Delete "$INSTDIR\uninstall.exe"
  RMDir  "$INSTDIR"
  MessageBox MB_ICONINFORMATION "${APP_NAME} foi removido com sucesso."
SectionEnd