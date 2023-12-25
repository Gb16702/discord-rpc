package discord

import (
	"log"

	"github.com/hugolgst/rich-go/client"
)

func ClientLogin(appKey string, logger *log.Logger) {
	if err := client.Login(appKey); err != nil {
		logger.Println("Erreur lors de la connexion du client")
		return
	}

	logger.Println("Connexion à Discord réussie")
}

func ClientLogout(logger *log.Logger) {
	client.Logout()
	logger.Println("Déconnexion de Discord réussie")
}
