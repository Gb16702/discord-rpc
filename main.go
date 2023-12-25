package main

import (
	"fmt"
	"log"
	"os"

	// "os"
	"time"

	"github.com/hugolgst/rich-go/client"
	"github.com/shirou/gopsutil/process"

	"dbd/discord"
	"dbd/load"
	"dbd/logs"
	"dbd/watcher"
)


type ProcessState struct {
	Running 	bool
	StartTime 	time.Time
}

func main() {
	fmt.Println("Lancement de l'application")
	logs.Logs()
	defer logs.CloseLogs()

	env := load.LoadEnvVariables(logs.Logger);

	get_working_directory, err := os.Getwd();
	if err != nil {
		logs.Logger.Println("Erreur lors de la récupération du chemin du fichier exécutable: ", err)
		return
	}

	player_file_path := get_working_directory + "\\player\\player_settings.json"

	logs.Logger.Println("Le fichier de configuration du joueur a été trouvé")

	process_state := ProcessState{}

	currentSettings := watcher.UpdateStatusFromFile(player_file_path, logs.Logger)

	defer func() {
		if process_state.Running {
			discord.ClientLogout(logs.Logger)
		}
	}()

	go watcher.WatchFile(player_file_path, logs.Logger, func(settings watcher.UserSettings) {
		currentSettings  = settings
		if process_state.Running {
			updateRichPresence(env, process_state, currentSettings, logs.Logger)
		}
	})


	for {
		running := isProcessRunning(env.PROCESS_KEY, logs.Logger)
		if running != process_state.Running {
			logs.Logger.Println("Test")
			discord.ClientLogin(env.APP_KEY, logs.Logger)
			process_state.Running = running
			logs.Logger.Println(env.PROCESS_KEY, "est en cours d'exécution")
			process_state.StartTime = time.Now()
			updateRichPresence(env, process_state, currentSettings, logs.Logger)
		}

		time.Sleep(5 * time.Second)
	}
}

func updateRichPresence(env load.EnvironmentVariables, state ProcessState, currentSettings watcher.UserSettings, logger *log.Logger) {
	if state.Running {
		err := client.SetActivity(client.Activity{
			State:      currentSettings.Description,
			Details:    currentSettings.Title,
			LargeImage: env.IMAGE_KEY,
			Timestamps: &client.Timestamps{Start: &state.StartTime},
		})

		if err != nil {
			logger.Println("Erreur lors de la mise à jour de la présence: ", err)
		}
	} else {
		discord.ClientLogout(logger)
	}
}

func isProcessRunning(process_key string, logger *log.Logger) bool {
	if processes, err := process.Processes(); err == nil {
		for _, p := range processes {
			name, err := p.Name()
			if err == nil && name == process_key {
				return true
			}
		}
	} else {
		logger.Println("Erreur lors de la récupération des processus: ", err)
		return false
	}

	return false
}
