package mysql

import (
	"errors"
	// "fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// Load mysql options.
// See also https://dev.mysql.com/doc/refman/5.5/en/option-files.html
func LoadOptions(conf string) (map[string]map[string]string, error) {
	props := make(map[string]map[string]string)
	content, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}

	var group map[string]string
	groupName := ""
	commentMatcher, _ := regexp.Compile("^[#;].*")
	groupMatcher, _ := regexp.Compile("^\\[(.+?)\\]")
	keyOnlyMatcher, _ := regexp.Compile("^([\\w\\d\\_\\-]+)\\s*$")
	keyValueMatcher, _ := regexp.Compile("^(.+?)\\s*[=]\\s*(.+)")

	key, value := "", ""
	for _, line := range strings.Split(string(content), "\n") {
		if len(strings.TrimSpace(line)) == 0 || commentMatcher.MatchString(line) {
			continue
		}
		if groupMatcher.MatchString(line) {
			groupName = groupMatcher.FindString(line)
			groupName = strings.ToLower(groupName[1 : len(groupName)-1])
		} else {
			if keyOnlyMatcher.MatchString(line) {
				key = strings.TrimSpace(line)
				value = "true"
			} else if keyValueMatcher.MatchString(line) {
				keyValue := strings.SplitN(line, "=", 2)
				key = strings.TrimSpace(keyValue[0])
				value = strings.TrimSpace(keyValue[1])
			} else {
				return nil, errors.New("fail to parse: " + line)
			}
			group = props[groupName]
			if group == nil {
				group = make(map[string]string)
				props[groupName] = group
			}
			group[key] = value
		}
	}
	return props, nil
}

// Change the content of the origin conf with the properties and save to a new conf file.
// You may use the same conf file and it will replace the origin conf.
func SaveOptions(originConf, newConf string, props map[string]map[string]string) error {
	content, err := ioutil.ReadFile(originConf)
	if err != nil {
		return err
	}

	groupName, nextGroupName := "", ""
	commentMatcher, _ := regexp.Compile("^[#;].*")
	groupMatcher, _ := regexp.Compile("^\\[(.+?)\\]")
	keyOnlyMatcher, _ := regexp.Compile("^([\\w\\d\\_\\-]+)\\s*$")
	keyValueMatcher, _ := regexp.Compile("^(.+?)\\s*[=]\\s*(.+)")

	key := ""
	newContent := ""
	for _, line := range strings.Split(string(content), "\n") {
		if len(strings.TrimSpace(line)) == 0 || commentMatcher.MatchString(line) {
			if len(groupName) > 0 && len(props[groupName]) == 0 {
				continue
			}

			newContent += line + "\n"
			continue
		}
		if groupMatcher.MatchString(line) {
			nextGroupName = groupMatcher.FindString(line)
			nextGroupName = strings.ToLower(nextGroupName[1 : len(nextGroupName)-1])
			if len(groupName) != 0 {
				for restKey, restValue := range props[groupName] {
					newContent += restKey + " = " + restValue + "\n"
				}
				delete(props, groupName)
			}
			if len(props[nextGroupName]) != 0 {
				newContent += "[" + nextGroupName + "]\n"
			}
			groupName = nextGroupName
			continue
		}
		if len(groupName) != 0 {
			group := props[groupName]
			if group == nil {
				continue
			}
			if keyOnlyMatcher.MatchString(line) {
				key = strings.TrimSpace(line)
			} else if keyValueMatcher.MatchString(line) {
				keyValue := strings.SplitN(line, "=", 2)
				key = strings.TrimSpace(keyValue[0])
			} else {
				return errors.New("fail to parse: " + line)
			}
			if len(group[key]) > 0 {
				newContent += key + " = " + group[key] + "\n"
				delete(group, key)
			}
			if len(group) == 0 {
				delete(props, groupName)
			}
		}
	}

	for moreName, moreGroup := range props {
		if len(moreGroup) == 0 {
			continue
		}
		newContent += "[" + moreName + "]\n"
		for moreKey, moreValue := range moreGroup {
			newContent += moreKey + " = " + moreValue + "\n"
		}
	}
	ioutil.WriteFile(newConf, []byte(newContent), 0644)
	return nil
}
