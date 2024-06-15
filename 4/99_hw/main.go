package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"strings"
)

type SearchServer struct {
	server *httptest.Server
	client SearchClient
}

type XmlData struct {
	XMLName xml.Name     `xml:"root"`
	Row     []XmlDataRow `xml:"row"`
}

type XmlDataRow struct {
	XMLName   xml.Name `xml:"row"`
	Id        int      `xml:"id"`
	FirstName string   `xml:"first_name"`
	LastName  string   `xml:"last_name"`
	About     string   `xml:"about"`
	Age       int      `xml:"age"`
	Gender    string   `xml:"gender"`
}

const ACCESS_TOKEN = "g_max"

func SearchServerFunc(w http.ResponseWriter, r *http.Request) {
	// проверяем текущий токен
	if r.Header.Get("AccessToken") != ACCESS_TOKEN {
		http.Error(w, "Incorrect AccessToken", http.StatusUnauthorized)
		return
	}

	// открываем файл, если токен валидный
	xml_file, err := os.Open("dataset.xml")
	if err != nil {
		http.Error(w, "Incorrect Open file", http.StatusExpectationFailed)
		return
	}
	var (
		users    []User
		xml_data XmlData
	)
	// прочитаем данные из xml-файл файла и распакуем их в нужный тип
	data, err := ioutil.ReadAll(xml_file)
	err = xml.Unmarshal(data, &xml_data)
	if err != nil {
		http.Error(w, "Incorrect Unmarshalling", http.StatusInternalServerError)
	}

	// обрабатываем запрос запрос
	query := r.URL.Query()
	for _, a := range xml_data.Row {
		//fmt.Println(a.FirstName+" ", a.LastName)
		if query.Get("query") != "" {
			match := strings.Contains(query.Get("query"), a.FirstName) ||
				strings.Contains(query.Get("query"), a.LastName) || strings.Contains(query.Get("query"), a.About)
			if !match {
				continue
			} else {
				users = append(users, User{
					Id:     a.Id,
					Name:   a.FirstName + " " + a.LastName,
					Age:    a.Age,
					About:  a.About,
					Gender: a.Gender,
				})
			}
		}
	}

	// обрабатываем order_field и order_by
	order_by, err := strconv.Atoi(query.Get("order_by"))
	if err != nil {
		http.Error(w, "Incorrect type of variable", http.StatusInternalServerError)
	}
	if order_by != OrderByAsIs {
		var comp func(u1, u2 User) bool
		str := query.Get("order_field")
		if order_by == 1 {
			if str == "Name" || str == "" {
				comp = func(u1, u2 User) bool {
					return u1.Name > u2.Name
				}
			} else if str == "Id" {
				comp = func(u1, u2 User) bool {
					return u1.Id > u2.Id
				}
			} else if str == "Age" {
				comp = func(u1, u2 User) bool {
					return u1.Age > u2.Age
				}
			}
			sort.Slice(users, func(i, j int) bool {
				return comp(users[i], users[j])
			})
		} else {
			if str == "Name" || str == "" {
				comp = func(u1, u2 User) bool {
					return u1.Name < u2.Name
				}
			} else if str == "Id" {
				comp = func(u1, u2 User) bool {
					return u1.Id < u2.Id
				}
			} else if str == "Age" {
				comp = func(u1, u2 User) bool {
					return u1.Age < u2.Age
				}
			}
			sort.Slice(users, func(i, j int) bool {
				return comp(users[i], users[j])
			})
		}
	}

	// limit
	// offset
	limit, _ := strconv.Atoi("limit")
	offset, _ := strconv.Atoi("offset")
	if limit > 0 {
		from := offset
		if from > len(users)-1 {
			users = []User{}
		} else {
			to := offset + limit
			if to > len(users) {
				to = len(users)
			}
			users = users[from:to]
		}
	}

	// упаковываем в json и отпавляем
	json_, err := json.Marshal(&users)
	if err != nil {
		http.Error(w, "Error in creating JSON", http.StatusInternalServerError)
	}

	/*for _, i := range users {
		fmt.Println(i)
	}*/

	w.Header().Set("Content", "application/json")
	_, err_ := w.Write(json_)
	if err_ != nil {
		http.Error(w, "Error in write JSON", http.StatusInternalServerError)
	}
}

// кастомный тестовый сервер
func TestingServer(token string) SearchServer {
	srv := httptest.NewServer(http.HandlerFunc(SearchServerFunc))
	clt := SearchClient{token, srv.URL}
	return SearchServer{srv, clt}
}

func (server *SearchServer) close() {
	server.server.Close()
}

func main() {
	server := TestingServer(ACCESS_TOKEN)
	resp, err := server.client.FindUsers(SearchRequest{Query: "Annie Pentihina"})
	if err != nil {
		fmt.Println("Error")
		return
	}
	fmt.Println(resp)
}
