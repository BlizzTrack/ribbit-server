/*
 * Copyright (c) 2020. BlizzTrack
 */

package managers

import (
	"errors"
	"github.com/blizztrack/ribbit-go"
	"github.com/jhillyerd/enmime"
	"github.com/mitchellh/mapstructure"
	"strings"
)

func ToSummaryList(raw string) []ribbit.SummaryItem {
	raw, _ = parseMime(raw)

	var result []ribbit.SummaryItem
	mapstructure.Decode(parseManifest(raw), &result)

	return result
}

func parseMime(raw string) (string, error) {
	content := raw
	r := strings.NewReader(content)
	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return "", err
	}

	if env.Root == nil || env.Root.FirstChild == nil {
		return "", errors.New("root or firstchild of root is empty")
	}

	return string(env.Root.FirstChild.Content), nil
}

func parseManifest(file string) []map[string]string {
	lines := strings.Split(file, "\n")
	keys := strings.Split(lines[0], `|`)
	keysList := make([]string, len(keys))

	for i := 0; i < len(keys); i++ {
		keyList := strings.Split(keys[i], `!`)

		keysList[i] = strings.ToLower(keyList[0])
	}

	var data []map[string]string
	for i := 1; i < len(lines); i++ {
		if len(strings.TrimSpace(lines[i])) > 0 {
			if !strings.HasPrefix(lines[i], "#") {
				local := make(map[string]string)

				lineData := strings.Split(lines[i], `|`)

				for x := 0; x < len(keysList); x++ {
					if len(lineData[x]) > 0 {
						local[keysList[x]] = lineData[x]
					}
				}

				data = append(data, local)
			}
		}
	}

	return data
}