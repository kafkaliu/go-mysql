package mysql

import (
	"encoding/json"
	"errors"
	// "fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

// see also https://dev.mysql.com/doc/refman/5.5/en/option-files.html
type MySqlOptionFile struct {
	FileName   string
	Properties map[string]map[string]string
}

func clone(src map[string]map[string]string) map[string]map[string]string {
	dest := make(map[string]map[string]string)
	for name, group := range src {
		destGroup := dest[name]
		if destGroup == nil {
			destGroup = make(map[string]string)
			dest[name] = destGroup
		}
		for key, value := range group {
			destGroup[key] = value
		}
	}
	return dest
}

func (this *MySqlOptionFile) String() string {
	str, _ := json.Marshal(this.Properties)
	return string(str)
}

func (this *MySqlOptionFile) Load() error {
	props := make(map[string]map[string]string)
	content, err := ioutil.ReadFile(this.FileName)
	if err != nil {
		return err
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
				return errors.New("fail to parse: " + line)
			}
			group = props[groupName]
			if group == nil {
				group = make(map[string]string)
				props[groupName] = group
			}
			group[key] = value
		}
	}
	this.Properties = props
	return nil
}

func (this *MySqlOptionFile) Save(newFile string) error {
	props := clone(this.Properties)
	content, err := ioutil.ReadFile(this.FileName)
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
	ioutil.WriteFile(newFile, []byte(newContent), 0644)
	return nil
}
