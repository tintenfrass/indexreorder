package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Version 0.1")
	baptism := []string{"baptism.txt"}
	marriage := []string{"marriage.txt"}

	//Config File einlesen
	data, err := ioutil.ReadFile("config.txt")
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(data), "\r\n")
	for _, line := range lines {
		content := strings.Split(line, "=")
		if len(content) > 1 {
			switch content[0] {
			case "baptism":
				baptism = strings.Split(content[1], ",")
			case "marriage":
				marriage = strings.Split(content[1], ",")
			}
		}
	}

	//Heiraten einlesen
	for _, file := range marriage {
		data, err = ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
		}
		importMarriage(strings.Split(string(data), "\r\n"), file)
	}

	//Taufen einlesen
	for _, file := range baptism {
		data, err = ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
		}
		importBaptism(strings.Split(string(data), "\r\n"), file)
	}

	printOut()

	fmt.Println("\r\nreorder.txt wurde geschrieben")
	fmt.Println("\r\nPress 'Enter' to close...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

var fullData = map[string]map[int]map[string][]string{}

func importMarriage(lines []string, file string) {
	year := 0
	for _, line := range lines {
		content := strings.Split(strings.TrimSpace(line), " ")
		if len(content) < 2 {
			if len(content[0]) >= 4 {
				number, err := strconv.Atoi(content[0][:4])
				if err == nil {
					year = number
				}
			}
			continue
		}
		if year == 0 {
			continue
		}

		//Anfangsbuchstaben bestimmen
		content[len(content)-2] = replaceName(content[len(content)-2])

		if len(content[len(content)-2]) < 1 || len(content[len(content)-1]) < 1 {
			fmt.Println("error bei ", content)
			continue
		}

		key1 := mapPhonetic(content[len(content)-2][:1])
		key2 := mapPhonetic(content[len(content)-1][:1])

		if fullData[key1+key2] == nil {
			fullData[key1+key2] = make(map[int]map[string][]string)
		}
		if fullData[key1+key2][year] == nil {
			fullData[key1+key2][year] = make(map[string][]string)
		}

		fullData[key1+key2][year]["m"] = append(fullData[key1+key2][year]["m"], strconv.Itoa(year)+" "+strings.Join(content, " ")+"\t\t\t\t\t\t\t("+file+")")
	}
}

func importBaptism(lines []string, file string) {
	year := 0
	for _, line := range lines {
		content := strings.Split(line, " ")
		if len(content) < 3 || strings.Contains(strings.TrimSpace(line), "Teil") {
			if len(content[0]) >= 4 {
				number, err := strconv.Atoi(content[0][:4])
				if err == nil {
					year = number
				}
			}
			continue
		}
		if year == 0 {
			continue
		}

		//Anfangsbuchstaben bestimmen
		content[len(content)-2] = replaceName(content[len(content)-2])

		if len(content[len(content)-2]) < 1 || len(content[len(content)-1]) < 1 {
			fmt.Println("error bei ", content)
			continue
		}

		key1 := mapPhonetic(content[len(content)-2][:1])
		key2 := mapPhonetic(content[len(content)-1][:1])

		if fullData[key1+key2] == nil {
			fullData[key1+key2] = make(map[int]map[string][]string)
		}
		if fullData[key1+key2][year] == nil {
			fullData[key1+key2][year] = make(map[string][]string)
		}

		fullData[key1+key2][year]["b"] = append(fullData[key1+key2][year]["b"], "\t"+strconv.Itoa(year)+" "+strings.Join(content, " ")+"\t\t\t\t\t\t("+file+")")
	}
}

func printOut() {
	file, err := os.Create("reordered.txt")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, key := range fullData {
		for i := 1500; i < 2000; i++ {
			if m, ok := key[i]["m"]; ok {
				for _, data := range m {
					datawriter.WriteString(data + "\n")
				}
			}
			if b, ok := key[i]["b"]; ok {
				for _, data := range b {
					datawriter.WriteString(data + "\n")
				}
			}
		}
		datawriter.WriteString("\n")
	}

	datawriter.Flush()
	file.Close()
}

func mapPhonetic(input string) string {
	input = strings.ToLower(input)
	m := map[string]string{
		"t": "d",
		"z": "d",
		"p": "b",
		"c": "k",
		"n": "h",
	}
	if mapped, ok := m[input]; ok {
		return mapped
	}
	return input
}

func replaceName(input string) string {
	m := map[string]string{
		"Jorge":  "Georg",
		"Brosi":  "Ambrosius",
		"Johann": "Hans",
		"Thoni":  "Anton",
	}
	for old, base := range m {
		if strings.Contains(input, old) {
			return base + "(" + input + ")"
		}
	}

	return input
}
