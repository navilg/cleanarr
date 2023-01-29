package internal

import "fmt"

func Job(isDryRun bool) error {

	if Config.Radarr.Enabled {
		// fmt.Println(config.Radarr.B64APIKey)
		ignoreTagId, err := GetTagIdFromRadarr(Config.IgnoreTag)
		if err != nil {
			return err
		}
		if ignoreTagId == nil {
			ignoreTagId, err = CreateTagInRadarr(Config.IgnoreTag)
			if err != nil {
				return err
			}
		}
		moviesdata, _ := GetMoviesData()
		moviesMarkedForDeletion, err := MarkMoviesForDeletion(moviesdata, *ignoreTagId, isDryRun)
		if err != nil {
			return err
		}

		moviesDeleted, err := DeleteExpiredMovies(moviesdata, *ignoreTagId, isDryRun)
		if err != nil {
			return err
		}

		if Config.Radarr.Notification && !isDryRun {
			body := `Movies deleted --> ` + fmt.Sprint(moviesDeleted) + `

			Movies Marked for deletion --> ` + fmt.Sprint(moviesMarkedForDeletion) + `

			Movies marked for deletion will be deleted in next maintenance run.
			To protect them from deletion, Add tag "` + Config.IgnoreTag + `" to movies in Radarr`

			SendEmailNotification("ALERT: [Cleanmyarr] [RADARR] Movies deleted", body)
		}
	}

	// fmt.Println(*config)

	return nil
}
