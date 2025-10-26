package main

import (
	"log"
	"time"

	miniflux "miniflux.app/v2/client"
)

var cronjobUpdate = make(chan bool, 1)

func requestCronConfigUpdate() {
	cronjobUpdate <- true
}

func startCronjob() {
	interval, _ := getCronInterval()
	ticker := time.NewTicker(interval)
	requestCronConfigUpdate()
	go func() {
		var usernames []string
		configs := make(map[string]Config)
		highestIds := make(map[string]int64)
		updateConfig := true // init cronjob data on first run

		for {
			if updateConfig { // update cronjob data
				updateConfig = false
				unames, err := getConfigUsernames()
				if err != nil {
					log.Fatal("get stored user names: ", err)
				}
				usernames = unames
				for _, username := range usernames {
					var config Config
					if err := readConfig(username, &config); err != nil {
						log.Fatal("err reading config:", err)
					}
					configs[username] = config
					highestIds[username] = -1 // reset known ids so all entries will be processed with new rules
				}
				log.Println("updated cron job")
			}

			for _, username := range usernames { // do filter
				client := miniflux.NewClient(getMinifluxUrl(), configs[username].ApiKey)
				val := getHighestId(client)
				if val > highestIds[username] {
					if err := runCronUpdate(client, configs[username], highestIds[username]); err != nil {
						log.Fatal("cron job failed:", err)
					}
					highestIds[username] = val
				}
			}
			select {
			case <-ticker.C:
			case <-cronjobUpdate:
				updateConfig = true
			}
		}
	}()
}

func runCronUpdate(client *miniflux.Client, config Config, oldHighestId int64) error {
	var readIds []int64
	var removedIds []int64
	entryResultSet, err := client.Entries(&miniflux.Filter{Starred: miniflux.FilterNotStarred, Direction: "desc", Status: "unread", AfterEntryID: oldHighestId, Limit: EntryProcessingLimit, Order: "id"})
	if err != nil {
		return err
	}
	for _, entry := range entryResultSet.Entries {
		for _, rule := range config.Rules {
			if isEntryMatchingRule(entry, &rule) {
				if !contains(readIds, entry.ID) && !contains(removedIds, entry.ID) {
					switch rule.State {
					case "read":
						readIds = append(readIds, entry.ID)
					case "removed":
						removedIds = append(removedIds, entry.ID)
					}
				}
			}
		}
	}
	if err := updateEntries(client, &readIds, &removedIds); err != nil {
		return err
	}
	return nil
}
