package main

import (
	"os"
	"fmt"
	"io"
	"encoding/json"
	"io/ioutil"
)

type Post struct {
	Id       int       `json:"id"`
	Content  string    `json:"content"`
	Author   Author    `json:"author"`
	Comments []Comment `json:"comments"`
}

type Author struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

func jsonParse() {
	jsonFile, err := os.Open("/Users/onepice2015/Desktop/project/code/go_web/src/seven_webapp/json_parsing_1/post.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}

	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return
	}

	var post Post
	json.Unmarshal(jsonData, &post)
	fmt.Println(post)
}

func decoderJsonParse() {
	fmt.Println("===========Decoder Json 解析==========")

	jsonFile1, err1 := os.Open("/Users/onepice2015/Desktop/project/code/go_web/src/seven_webapp/json_parsing_1/post.json")
	if err1 != nil {
		fmt.Println("Error opening JSON file:", err1)
		return
	}

	defer jsonFile1.Close()

	decoder := json.NewDecoder(jsonFile1)
	for {
		var post1 Post
		err := decoder.Decode(&post1)
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		fmt.Println(post1)
	}
}

func createJson() {
	fmt.Println("===========创建 Json 解析==========")

	post := Post{
		Id:      1,
		Content: "Hello World!",
		Author: Author{
			Id:   2,
			Name: "Sau Sheong",
		},
		Comments: []Comment{
			Comment{
				Id:      3,
				Content: "Have a greate day!",
				Author:  "Adam",
			},
			Comment{
				Id:      4,
				Content: "Have are you today?",
				Author:  "Betty",
			},
		},
	}

	output, err := json.MarshalIndent(&post, "", "\t\t")

	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	err = ioutil.WriteFile("/Users/onepice2015/Desktop/project/code/go_web/src/seven_webapp/json_parsing_1/output.json", output, 0644)

	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println(post)
}

func createEncoderJson() {
	fmt.Println("===========创建 Encoder Json 解析==========")

	post := Post{
		Id:      1,
		Content: "Hello World!",
		Author: Author{
			Id:   2,
			Name: "Sau Sheong",
		},
		Comments: []Comment{
			Comment{
				Id:      3,
				Content: "Have a greate day!",
				Author:  "Adam",
			},
			Comment{
				Id:      4,
				Content: "Have are you today?",
				Author:  "Betty",
			},
		},
	}

	outPutFile, err := os.Create("/Users/onepice2015/Desktop/project/code/go_web/src/seven_webapp/json_parsing_1/outputencode.json")

	if err != nil {
		fmt.Println("Error createing JSON file:", err)
		return
	}

	encoder := json.NewEncoder(outPutFile)
	err = encoder.Encode(&post)
	if err != nil {
		fmt.Println("Error encoding JSON to file:", err)
		return
	}
	fmt.Println(post)
}

func main() {
	jsonParse()
	decoderJsonParse()
	createJson()
	createEncoderJson()
}
