package main

import (
	"bufio"
	"fmt"
	"hw3/easyjson"
	"io"
	"os"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	user := easyjson.User{}
	map_ := map[string]bool{}
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	fileContents := bufio.NewScanner(file)
	fmt.Fprint(out, "found users:\n")
	counter := 0
	for fileContents.Scan() {
		if err = user.UnmarshalJSON(fileContents.Bytes()); err != nil {
			panic(err)
		}
		Android, MSIE := false, false
		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				Android = true
			} else if strings.Contains(browser, "MSIE") {
				MSIE = true
			} else {
				continue
			}
			map_[browser] = true
		}
		if Android && MSIE {
			email := strings.Replace(user.Email, "@", " [at] ", 1)
			fmt.Fprintln(out, fmt.Sprintf("[%d] %s <%s>", counter, user.Name, email))
		}
		counter++
	}
	fmt.Fprintln(out, "\nTotal unique browsers", len(map_))

}
