package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// Представим дерево каталогов как B-дерево и будем обходить его в глубину,
// тем самым получая рисунок этого дерева каталогов, если последний параметр - true,
// выводим все файлы и каталоги, иначе только каталоги

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err_ := dirTree(out, path, printFiles)
	if err_ != nil {
		panic(err_.Error())
	}
}

// узел дерева

type TreeNode struct {
	info     os.FileInfo
	children []TreeNode
}

// Рекурсивно обходим все директории

func getNodes(path string, withFiles bool) ([]TreeNode, error) {
	// возвращаем все файлы из данной директории
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var nodes []TreeNode
	// просматриваем каждый элемент директории
	for _, file := range files {
		// не директория и "false" - переходим к следующему элементу в данной директории
		if !withFiles && !file.IsDir() {
			continue
		}

		// создаём новую ноду
		node := TreeNode{
			info: file,
		}

		// если директория - переходим в неё
		if file.IsDir() {
			// рекурсивно переходим к следующему узлу(следующий слой в дереве каталога)
			children, err := getNodes(path+string(os.PathSeparator)+file.Name(), withFiles)
			if err != nil {
				return nil, err
			}
			// инициализируем ноду массивом нод, полученных вследствие рекурсии
			node.children = children
		}
		// ( если нет , массив занулится(детей нет) ) и добавляем ноду в массив
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// рекурсивно записываем ноды в интерфейсный объект io.Writer, начиная с префикса
func printNodes(out io.Writer, nodes []TreeNode, parentPrefix string, flag bool) {
	// создаем множество переменных в одном блоке
	var (
		lastIdx     = len(nodes) - 1
		prefix      = "├───"
		childPrefix = "│\t"
	)

	// проходимся по полученному массиву нод
	for i, node := range nodes {
		// последняя нода в дереве
		if i == lastIdx {
			prefix = "└───"
			childPrefix = "\t"
		}
		if !node.info.IsDir() && flag {
			if node.info.Size() == 0 {
				fmt.Fprint(out, parentPrefix, prefix, node.info.Name(), " ", "(", "empty", ")", "\n")
			} else {
				// записываем в объект out: общий префикс для каталога, текущий префикс и имя файла с его размером
				fmt.Fprint(out, parentPrefix, prefix, node.info.Name(), " ", "(", node.info.Size(), "b)", "\n")
			}
		} else {
			fmt.Fprint(out, parentPrefix, prefix, node.info.Name(), "\n")
		}

		//  если директория, рекурсивно переходим в узел ребенка, изменяя состояние  общего префикса
		// , добавляя табуляцию в строку
		if node.info.IsDir() {
			printNodes(out, node.children, parentPrefix+childPrefix, flag)
		}
	}
}

func dirTree(writer io.Writer, path string, printFiles bool) error {
	nodes, err := getNodes(path, printFiles)
	if err != nil {
		return err
	}
	printNodes(writer, nodes, "", printFiles)
	return nil
}
