package update

import (
	"os/exec"
	"runtime"
)

// openBrowser intenta abrir url en el navegador por defecto, según el SO. Es
// best-effort: cualquier fallo se devuelve como error para que el llamante lo
// ignore; la URL siempre se imprime como alternativa (R8).
func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
