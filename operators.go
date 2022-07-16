package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func List(file *os.File, writer io.Writer) error {
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		_, err := writer.Write(sc.Bytes())
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	}

	return nil
}

func Add(args Arguments, file *os.File, writer io.Writer) error {
	if args["item"] == "" {
		return fmt.Errorf("-item flag has to be specified")
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	var fileItems []Item
	if len(data) > 1 {
		err = json.Unmarshal(data, &fileItems)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	var items Item
	err = json.Unmarshal([]byte(args["item"]), &items)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for i := range fileItems {
		if fileItems[i].Id == items.Id {
			writer.Write([]byte("Item with id " + items.Id + " already exists"))
			return nil
		}
	}

	file.Write([]byte("[" + args["item"] + "]"))
	return nil
}

func FindById(args Arguments, file *os.File, writer io.Writer) error {
	if args["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	var fileItems []Item
	if len(data) > 1 {
		err = json.Unmarshal(data, &fileItems)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	for i := range fileItems {
		if fileItems[i].Id == args["id"] {
			res, err := json.Marshal(fileItems[i])
			if err != nil {
				fmt.Errorf("marshalling is failed: %v", err.Error())
			}

			writer.Write(res)
			return nil
		}
	}

	return nil
}

func Remove(args Arguments, file *os.File, writer io.Writer) error {
	if args["id"] == "" {
		return fmt.Errorf("-id flag has to be specified")
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	var fileItems []Item
	if len(data) > 1 {
		err = json.Unmarshal(data, &fileItems)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	for i := range fileItems {
		if fileItems[i].Id == args["id"] {
			fileItems = append(fileItems[:i], fileItems[i+1:]...)

			res, err := json.Marshal(fileItems)
			if err != nil {
				fmt.Errorf("marshalling is failed: %v", err.Error())
			}

			file.Seek(0, io.SeekStart)
			file.Truncate(0)

			file.Write(res)
			return nil
		}
	}

	writer.Write([]byte("Item with id " + args["id"] + " not found"))
	return nil
}
