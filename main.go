package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	jsonFile, err := os.Open("_data/json/allCountries.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	scanner := bufio.NewScanner(jsonFile)

	var (
		maxDocs   uint                            = 10000 // 5000000 10512109
		docId     uint                            = 0
		docs      map[uint]map[string]interface{} = make(map[uint]map[string]interface{}, maxDocs)
		revIdxStr map[string]map[string][]uint    = make(map[string]map[string][]uint)
		revIdxNum map[string]map[float64][]uint   = make(map[string]map[float64][]uint)
	)

	for scanner.Scan() {
		var doc map[string]interface{}
		err := json.Unmarshal(scanner.Bytes(), &doc)
		if err != nil {
			//panic("line " + strconv.Itoa(int(docId+1)) + ": " + err.Error())
			fmt.Printf("Ignored line %d : %s\n", docId+1, err.Error())
		}

		docs[docId] = doc

		for key, value := range doc {
			switch value.(type) {
			case string:
				valueStr := value.(string)
				if _, ok := revIdxStr[key]; !ok {
					revIdxStr[key] = make(map[string][]uint)
				}
				revIdxStr[key][valueStr] = append(revIdxStr[key][valueStr], docId)
			case float64:
				valueNum := value.(float64)
				if _, ok := revIdxNum[key]; !ok {
					revIdxNum[key] = make(map[float64][]uint)
				}
				revIdxNum[key][valueNum] = append(revIdxNum[key][valueNum], docId)
			}
		}

		docId++
		if docId >= maxDocs {
			break
		}
	}

	fmt.Println("docs loaded:", docId)

	//fmt.Printf("%+v\n", docs)
	//fmt.Printf("%+v\n", revIdxStr)
	//fmt.Printf("%+v\n", revIdxNum)

	err = scanner.Err()
	if err != nil {
		panic(err)
	}

	var pop float64 = 0
	var i uint

	for i = 0; i < docId; i++ {
		if value, ok := docs[i]["dem"]; ok {
			pop += value.(float64)
		}
	}

	fmt.Printf("pop: %.f\n", pop)

	//bufio.NewReader(os.Stdin).ReadBytes('\n')
}
