package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// CreatedTask 構造体は、レスポンスとして返ってくる
// JSONメッセージを受け取るための構造体です。
type CreatedTask struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

// TaskCreationReturnedStatusCodeUnexpected はタスク作成のリクエストを行った際に、
// 201 Created 以外のレスポンスコードが返ってきた場合のエラーメッセージです。
const TaskCreationReturnedStatusCodeUnexpected = "期待したレスポンスステータスコード(201 Created)ではありません。"

// ResponseBodyReadFailure はリクエストに対する
// レスポンスでボディの読み込みに失敗した場合のエラーメセージです。
const ResponseBodyReadFailure = "レスポンスボディの読み込みに失敗しました"

// ResponseBodyParseFailure はリクエストに対する
// レスポンスでボディのパースに失敗した場合のエラーメッセージです。
const ResponseBodyParseFailure = "レスポンスボディのパースに失敗しました"

// CreateTask はToDoクライアントからToDoサーバへのTask作成を行います。
func CreateTask(protocol string, host string, port int, token string, title string, description string) (CreatedTask, error) {
	path := "/api/task"
	url := protocol + "://" + host + ":" + strconv.Itoa(port) + path

	client := &http.Client{}

	authHeader := "JWT " + token

	taskInfo := fmt.Sprintf(`{
		"title": "%s",
		"description": "%s"
	}`, title, description)

	req, err := http.NewRequest("POST", url, strings.NewReader(taskInfo))
	req.Header.Set("Authorization", authHeader)

	if err != nil {
		log.Fatal(err)
	}

	// Task作成のリクエストを発行します。
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		log.Println(res)
		log.Println(TaskCreationReturnedStatusCodeUnexpected)
		err := errors.New(TaskCreationReturnedStatusCodeUnexpected)

		return CreatedTask{}, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		log.Println(ResponseBodyReadFailure)
		err := errors.New(ResponseBodyReadFailure)
		return CreatedTask{}, err
	}

	var task CreatedTask
	if err := json.Unmarshal(body, &task); err != nil {
		log.Println(err)
		log.Println(ResponseBodyParseFailure)
		err := errors.New(ResponseBodyParseFailure)
		return CreatedTask{}, err
	}

	log.Printf("Task created. TaskID: %d\n", task.ID)
	log.Printf("Task Title: %s\n", task.Title)
	return task, err
}

// Task 構造体は、タスク取得リクエストのレスポンスとして返ってくる
// JSONメッセージを受け取るための構造体です。
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// TaskGetReturnedStatusCodeUnexpected はタスク取得リクエスト実行時に、
// ステータスコードとして200 OK以外が返ってきた場合に表示するメッセージです。
const TaskGetReturnedStatusCodeUnexpected = "期待したレスポンスステータスコード(200 OK)ではありません。"

// TaskGetReturnedNotFoundStatusCode はタスク取得リクエスト実行時に、
// ステータスコードとして404 Not Foundが返ってきた場合に表示するメッセージです。
const TaskGetReturnedNotFoundStatusCode = "指定したタスクのIDに対応するタスクが見つかりません。"

// GetTask は指定したtask_idの値に対応したTaskを返します。
func GetTask(protocol string, host string, port int, token string, taskID int) ([]Task, error) {
	path := "/api/task/" + strconv.Itoa(taskID)
	url := protocol + "://" + host + ":" + strconv.Itoa(port) + path

	client := &http.Client{}

	authHeader := "JWT " + token

	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	req.Header.Set("Authorization", authHeader)

	// Task作成のリクエストを発行します。
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	// レスポンスのステータスコードが200 OK以外だったときの処理
	if res.StatusCode != http.StatusOK {
		log.Println(TaskGetReturnedStatusCodeUnexpected)
		if res.StatusCode == http.StatusNotFound {
			log.Println(TaskGetReturnedNotFoundStatusCode)
			err := errors.New(TaskGetReturnedNotFoundStatusCode)
			return []Task{Task{}}, err
		}
		err := errors.New(TaskGetReturnedStatusCodeUnexpected)
		return []Task{Task{}}, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		log.Println(ResponseBodyReadFailure)
		return []Task{Task{}}, err
	}

	var tasks []Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		log.Printf("%s\n", body)
		log.Println(err)
		log.Println(ResponseBodyParseFailure)
		err := errors.New(ResponseBodyParseFailure)
		return []Task{Task{}}, err
	}

	return tasks, err
}

// GetTasks はリクエストしたユーザに紐づくタスク全てを返します
func GetTasks(protocol string, host string, port int, token string) ([]Task, error) {
	path := "/api/task"
	url := protocol + "://" + host + ":" + strconv.Itoa(port) + path

	client := &http.Client{}

	authHeader := "JWT " + token

	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	req.Header.Set("Authorization", authHeader)

	// Task取得のリクエストを発行します。
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	// レスポンスのステータスコードが200 OK以外だったときの処理
	if res.StatusCode != http.StatusOK {
		log.Println(TaskGetReturnedStatusCodeUnexpected)
		err := errors.New(TaskGetReturnedStatusCodeUnexpected)
		return []Task{Task{}}, err
	}

	// レスポンスのステータスコードが200 OK以外だったときの処理
	if res.StatusCode != http.StatusOK {
		log.Println(TaskGetReturnedStatusCodeUnexpected)
		if res.StatusCode == http.StatusNotFound {
			err := errors.New(TaskGetReturnedNotFoundStatusCode)
			return []Task{Task{}}, err
		}
		err := errors.New(TaskGetReturnedStatusCodeUnexpected)
		return []Task{Task{}}, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		log.Println(ResponseBodyReadFailure)
		return []Task{Task{}}, err
	}

	var tasks []Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		log.Printf("%s\n", body)
		log.Println(err)
		log.Println(ResponseBodyParseFailure)
		err := errors.New(ResponseBodyParseFailure)
		return []Task{Task{}}, err
	}

	return tasks, err
}

const TaskDeleteReturnedStatusCodeUnexpected = "指定されたIDに対応するTaskを削除しようとしましたが、想定外のステータスコードが返されました。"

const TaskDeleteReturnedNotFoundStatusCode = "削除の為に指定したIDに対応するTaskが見つかりませんでした。"

// DeleteTask は指定されたIDをもつタスクの削除を試みます
func DeleteTask(protocol string, host string, port int, token string, taskID int) (Task, error) {
	path := "/api/task/" + strconv.Itoa(taskID)
	url := protocol + "://" + host + ":" + strconv.Itoa(port) + path

	client := &http.Client{}

	authHeader := "JWT " + token

	req, err := http.NewRequest("DELETE", url, strings.NewReader(""))
	req.Header.Set("Authorization", authHeader)

	// Task削除のリクエストを発行します。
	res, err := client.Do(req)

	// レスポンスのステータスコードが200 OK以外だったときの処理
	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			log.Println(TaskDeleteReturnedNotFoundStatusCode)
			err := errors.New(TaskDeleteReturnedNotFoundStatusCode)
			return Task{}, err
		}
		log.Println(TaskDeleteReturnedStatusCodeUnexpected)
		err := errors.New(TaskDeleteReturnedStatusCodeUnexpected)
		return Task{}, err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		log.Println(ResponseBodyReadFailure)
		return Task{}, err
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil {
		log.Printf("%s\n", body)
		log.Println(err)
		log.Println(ResponseBodyParseFailure)
		err := errors.New(ResponseBodyParseFailure)
		return Task{}, err
	}

	return task, err
}

// TaskUpdateReturnedStatusCodeUnexpected は200,400, 404以外のステータスコードが返ってきたときに指定
const TaskUpdateReturnedStatusCodeUnexpected = "指定されたIDに対応するTaskを更新しようとしましたが、想定外のステータスコードが返されました。"

// TaskUpdateReturnedNotFoundStatusCode は指定したタスクが存在しない、または他のユーザに紐付いているなどの理由で
// 404 Not Foundが返された場合に
const TaskUpdateReturnedNotFoundStatusCode = "更新の為に指定したIDに対応するTaskが見つかりませんでした。"

const TaskUpdateReturnedBadRequestStatusCode = "更新の為に指定したステータス情報が不正です"

// UpdateTask は指定されたIDを持つタスクの情報更新を行います。
func UpdateTask(protocol string, host string, port int, token string, taskID int, title string, description string, status string) (Task, error) {
	path := "/api/task/" + strconv.Itoa(taskID)
	url := protocol + "://" + host + ":" + strconv.Itoa(port) + path

	client := &http.Client{}

	task, err := GetTask(protocol, host, port, token, taskID)
	if err != nil {
		if err.Error() == TaskGetReturnedNotFoundStatusCode {
			log.Println(err)
			err = errors.New(TaskUpdateReturnedNotFoundStatusCode)
			log.Println(err)
			return task[0], err
		}
		if err.Error() == TaskGetReturnedStatusCodeUnexpected {
			log.Println(err)
			err = errors.New(TaskUpdateReturnedStatusCodeUnexpected)
			log.Println(err)
			return task[0], err
		}
	}

	authHeader := "JWT " + token

	if title == "" {
		title = task[0].Title
	}

	if description == "" {
		description = task[0].Description
	}

	if status == "" {
		status = task[0].Status
	}

	taskInfo := fmt.Sprintf(`{
		"title": "%s",
		"description": "%s",
		"status": "%s"
	}`, title, description, status)

	req, err := http.NewRequest("PATCH", url, strings.NewReader(taskInfo))
	req.Header.Set("Authorization", authHeader)

	// Task更新のリクエストを発行します。
	res, err := client.Do(req)

	// レスポンスのステータスコードが200 OK以外だったときの処理
	if res.StatusCode != http.StatusOK {
		// 404 Not Foundが返ってきた場合(基本的には到達不可能なコード)
		if res.StatusCode == http.StatusNotFound {
			log.Println(TaskUpdateReturnedNotFoundStatusCode)
			err := errors.New(TaskUpdateReturnedNotFoundStatusCode)
			return Task{}, err
		}

		// 400 Bad Requestが返ってきた場合
		if res.StatusCode == http.StatusBadRequest {
			log.Println(TaskUpdateReturnedBadRequestStatusCode)
			err := errors.New(TaskUpdateReturnedBadRequestStatusCode)
			return Task{}, err
		}

		log.Println(TaskUpdateReturnedStatusCodeUnexpected)
		err := errors.New(TaskUpdateReturnedStatusCodeUnexpected)
		return Task{}, err
	}

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
		log.Println(ResponseBodyReadFailure)
		err := errors.New(ResponseBodyReadFailure)
		return Task{}, err
	}

	var updatedTask Task
	if err := json.Unmarshal(body, &updatedTask); err != nil {
		log.Println(err)
		log.Println(ResponseBodyParseFailure)
		err := errors.New(ResponseBodyParseFailure)
		return Task{}, err
	}

	log.Printf("Task(ID=%d) is updated.\n", updatedTask.ID)

	return updatedTask, err
}
