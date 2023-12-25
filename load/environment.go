package load

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)


type EnvironmentVariables struct {
	APP_KEY 		string
	IMAGE_KEY 		string
	PROCESS_KEY 	string
}

func LoadEnvVariables(logger *log.Logger) EnvironmentVariables {
	if err := godotenv.Load("G:\\Programs\\discord\\dbd-rpc-bridge\\.env"); err != nil {
		logger.Println("Le fichier .env n'a pas été trouvé")
		panic(err)
	}

	logger.Println("Les variables d'environnement ont été chargées avec succès")

	return EnvironmentVariables{
		APP_KEY: os.Getenv("APP_KEY"),
		IMAGE_KEY: os.Getenv("IMAGE_KEY"),
		PROCESS_KEY: os.Getenv("PROCESS_KEY"),
	};
}
