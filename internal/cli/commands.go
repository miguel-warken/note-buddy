package cli

import (
	"fmt"
	"note-buddy-main/internal/store/sqlite"
)

type Command struct {
	Name string
	Args []string
}

func (a *App) handleCapture(args []string) error {
	db, err := sqlite.Open(a.dbPath, true)
	if err != nil {
		return fmt.Errorf("erro ao abrir o banco no capture: %w", err)
	}
	defer db.Close()

	fmt.Println("Executando o capture... Banco aberto com sucesso!")
	return nil
}

