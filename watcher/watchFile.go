package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type UserSettings struct {
	Title string
	Description string
}

type Settings struct {
    Titres       []Item `json:"titres"`
    Descriptions []Item `json:"descriptions"`
}

type Item struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    IsDefault bool   `json:"isDefault"`
}

func WatchFile(file string, logger *log.Logger, onUpdate func(UserSettings)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Erreur lors de la création du watcher: ", err)
		return
	}

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					settings := UpdateStatusFromFile(file, logger)
					onUpdate(settings)
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Println("Erreur lors de la surveillance du fichier")
			}
		}
	}()

	if err := watcher.Add(file); err != nil {
		logger.Println("Erreur lors de la surveillance du fichier")
		return
	}

	<-done
}

func UpdateStatusFromFile(file string, logger *log.Logger) UserSettings {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier: ", err)
		logger.Println("Erreur lors de la lecture du fichier")
		return UserSettings{}
	}

	var settings Settings;
	if err := json.Unmarshal(data, &settings); err != nil {
		fmt.Println("Erreur lors de la lecture du fichier: ", err)
		logger.Println("Erreur lors de la lecture du fichier")
		return UserSettings{}
	}

	logger.Println("Les paramètres ont été chargés avec succès")

	defaultTitle := GetUserConfig(settings.Titres)
	defaultDescription := GetUserConfig(settings.Descriptions)

	return UserSettings{
		Title: defaultTitle.Name,
		Description: defaultDescription.Name,
	}
}
