package service

import (
	"log"
	"os"
	"testing"
)

func TestCreateTask(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	TodoTitle := "TODO"
	TodoDescription := "TODO_DESCRIPTION"

	task, err := CreateTask("http", testTarget, 8000, loginConfig.Token, TodoTitle, TodoDescription)

	if task.Title != TodoTitle {
		t.Fail()
	}

	if task.Description != TodoDescription {
		t.Fail()
	}
}

func TestCreateTaskWithWrongStatusCode(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	TodoTitle := "TODO"
	TodoDescription := "TODO_DESCRIPTION"

	_, err := CreateTask("http", testTarget, 8000, "WrongToken", TodoTitle, TodoDescription)

	if err.Error() != TaskCreationReturnedStatusCodeUnexpected {
		t.Fail()
	}

}

// TestGetTask では単一のTaskを取得できることを確認する。
func TestGetTask(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	TodoTitle := "TODO"
	TodoDescription := "TODO_DESCRIPTION"
	TodoStatus := "TODO"

	createdTask, err := CreateTask("http", testTarget, 8000, loginConfig.Token, TodoTitle, TodoDescription)

	TargetTaskID := createdTask.ID
	log.Printf("TargetTaskID: %d\n", TargetTaskID)
	task, err := GetTask("http", testTarget, 8000, loginConfig.Token, TargetTaskID)

	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	if task[0].Title != TodoTitle {
		log.Printf("Title: %s\n", task[0].Title)
		t.Fail()
	}

	if task[0].Description != TodoDescription {
		log.Printf("Description: %s\n", task[0].Description)
		t.Fail()
	}

	if task[0].Status != TodoStatus {
		log.Printf("Status: %s\n", task[0].Status)
		t.Fail()
	}
}

// TestGetTaskWithoutTask では、存在しないTaskを指定して取得できないことを確認する。
// 具体的にはTaskGetReturnedNotFoundStatusCodeに対応するエラーメッセージが取得できることを
// 確認する。
func TestGetTaskWithoutTask(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	_, err = GetTask("http", testTarget, 8000, loginConfig.Token, 1000)

	if err.Error() != TaskGetReturnedNotFoundStatusCode {
		log.Fatal(err)
		t.Fail()
	}
}

// TestGetTaskWithOtherUsersTask では、Taskは存在するが他のUserに紐づくものである場合に、
// TaskGetReturnedNotFoundStatusCodeに対応するエラーメッセージが取得できることを確認する。
func TestGetTaskWithOtherUsersTask(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	// fujiwaraでログインしてタスクを作成
	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	TodoTitle := "TODO"
	TodoDescription := "TODO_DESCRIPTION"

	createdTask, err := CreateTask("http", testTarget, 8000, loginConfig.Token, TodoTitle, TodoDescription)

	// test_userでログインしてタスク取得を試みる
	loginConfig, err = Login("http", testTarget, 8000, "test_user", "test_password")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	TargetTaskID := createdTask.ID
	_, err = GetTask("http", testTarget, 8000, loginConfig.Token, TargetTaskID)

	if err.Error() != TaskGetReturnedNotFoundStatusCode {
		log.Fatal(err)
		t.Fail()
	}

}

// TestGetTask では、誤った認証情報を与えてリクエストを行うと
// TestGetReturnedStatusCodeUnexpectedに対応するエラーメッセージが取得できることを確認する。
func TestGetTaskWithWrongAuthInfo(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	_, err := GetTask("http", testTarget, 8000, "WrongAuthInfo", 1000)

	if err.Error() != TaskGetReturnedStatusCodeUnexpected {
		log.Fatal(err)
		t.Fail()
	}
}

// TestGetTasks では追加したタスクが格納されていることを確認する。
func TestGetTasks(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	// fujiwaraでログインしてタスクを作成
	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	TodoTitle := "TODO"
	TodoDescription := "TODO_DESCRIPTION"
	TaskLength := 4

	tasks, err := GetTasks("http", testTarget, 8000, loginConfig.Token)
	TaskLengthBeforeCreateTask := len(tasks)

	for i := 0; i < TaskLength; i++ {
		_, err := CreateTask("http", testTarget, 8000, loginConfig.Token, TodoTitle, TodoDescription)
		if err != nil {
			t.Fail()
		}
	}

	tasks, err = GetTasks("http", testTarget, 8000, loginConfig.Token)

	if err != nil {
		t.Fail()
	}

	if len(tasks)-TaskLengthBeforeCreateTask < TaskLength {
		t.Fail()
	}

	for i := 4; i > 0; i-- {
		if tasks[len(tasks)-i].Title != TodoTitle {
			t.Fail()
		}
		if tasks[len(tasks)-i].Description != TodoDescription {
			t.Fail()
		}
	}
}

// TestDeleteTask では作成したタスクが削除できていることを確認する。
func TestDeleteTask(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	TodoTitle := "TODO"
	TodoDescription := "TODO_DESCRIPTION"

	// 削除対象とするタスクを作成
	createdTask, err := CreateTask("http", testTarget, 8000, loginConfig.Token, TodoTitle, TodoDescription)
	// 削除を実施
	deletedTask, err := DeleteTask("http", testTarget, 8000, loginConfig.Token, createdTask.ID)

	// 作成したタスクと削除したタスクが一致することを確認
	if createdTask.ID != deletedTask.ID {
		t.Fail()
	}

	if createdTask.Title != deletedTask.Title {
		t.Fail()
	}

	if createdTask.Description != deletedTask.Description {
		t.Fail()
	}

	// 削除したタスクが取得できないことを確認
	_, err = GetTask("http", testTarget, 8000, loginConfig.Token, deletedTask.ID)

	if err.Error() != TaskGetReturnedNotFoundStatusCode {
		t.Fail()
	}
}

// TestDeleteTaskWithoutTask は存在しないタスクを削除しようとした際に、
// エラー(TaskGetReturnedNotFoundStatusCode)を返すことを確認する。
func TestDeleteTaskWithoutTask(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	_, err = DeleteTask("http", testTarget, 8000, loginConfig.Token, 10000)

	if err.Error() != TaskDeleteReturnedNotFoundStatusCode {
		log.Println(err)
		t.Fail()
	}
}

// TestDeleteTaskWithWrongAuth では、誤った認証情報を利用して、
// 404 Not Found以外のステータスコードが返ってきていることを確認する。
func TestDeleteTaskWithWrongAuth(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	_, err := DeleteTask("http", testTarget, 8000, "WrontToken", 100000)

	if err.Error() != TaskDeleteReturnedStatusCodeUnexpected {
		t.Fail()
	}
}

// TestDeleteTaskWithOtherUsersTask では、他のユーザに紐づく
// タスクを削除しようとして404 Not Foundが返ってきていることを確認する。
func TestDeleteTaskWithOtherUsersTask(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	TodoTitle := "TODO"
	TodoDescription := "TODO_DESCRIPTION"

	// 削除対象とするタスクを作成
	createdTask, err := CreateTask("http", testTarget, 8000, loginConfig.Token, TodoTitle, TodoDescription)

	// test_userでログインして作成したタスクの削除を試みる
	loginConfig, err = Login("http", testTarget, 8000, "test_user", "test_password")

	_, err = DeleteTask("http", testTarget, 8000, loginConfig.Token, createdTask.ID)

	if err.Error() != TaskDeleteReturnedNotFoundStatusCode {
		t.Fail()
	}
}

// TestUpdateTask では、更新対象のタスクを更新した際に更新できていることを確認する。
func TestUpdateTask(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	TodoTitle := "TODO"
	TodoDescription := "TODO_DESCRIPTION"

	// 更新対象とするタスクを作成
	createdTask, err := CreateTask("http", testTarget, 8000, loginConfig.Token, TodoTitle, TodoDescription)

	UpdatedTitle := "UPDATED"
	UpdatedDescription := "UPDATED_DESCRIPTION"
	UpdatedStatus := "RUNNING"
	// 更新を実施
	updatedTask, err := UpdateTask("http", testTarget, 8000, loginConfig.Token, createdTask.ID, UpdatedTitle, UpdatedDescription, UpdatedStatus)

	// 更新したタスクと削除したタスクが一致することを確認
	if createdTask.ID != updatedTask.ID {
		t.Fail()
	}

	if updatedTask.Title != UpdatedTitle {
		t.Fail()
	}

	if updatedTask.Description != UpdatedDescription {
		t.Fail()
	}

	// 更新したタスクを取得して意図した通りに更新されていることを確認
	task, err := GetTask("http", testTarget, 8000, loginConfig.Token, updatedTask.ID)
	if err != nil {
		t.Fail()
	}

	if task[0].ID != updatedTask.ID {
		t.Fail()
	}

	if task[0].Title != updatedTask.Title {
		t.Fail()
	}

	if task[0].Description != updatedTask.Description {
		t.Fail()
	}

	if task[0].Status != updatedTask.Status {
		t.Fail()
	}
}

// TestUpdateTaskWithBadRequest 不正なステータス情報をタスク更新時に渡して
// 400 Bad Requestが返ってくることを確認する。
func TestUpdateTaskWithBadRequest(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	TodoTitle := "TODO"
	TodoDescription := "TODO_DESCRIPTION"

	// 更新対象とするタスクを作成
	createdTask, err := CreateTask("http", testTarget, 8000, loginConfig.Token, TodoTitle, TodoDescription)

	UpdatedTitle := "UPDATED"
	UpdatedDescription := "UPDATED_DESCRIPTION"
	UpdatedStatus := "WRONG_STATUS" // 存在しないステータスを設定
	// 更新を実施
	_, err = UpdateTask("http", testTarget, 8000, loginConfig.Token, createdTask.ID, UpdatedTitle, UpdatedDescription, UpdatedStatus)

	if err.Error() != TaskUpdateReturnedBadRequestStatusCode {
		log.Println(err)
		t.Fail()
	}

}

// TestUpdateTaskWithoutTask 存在しないタスクに対して更新処理を試みて
// 404 Not Foundが返ってくることを確認する。
func TestUpdateTaskWithoutTask(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	UpdatedTitle := "UPDATED"
	UpdatedDescription := "UPDATED_DESCRIPTION"
	UpdatedStatus := "RUNNING"
	// 更新を実施
	_, err = UpdateTask("http", testTarget, 8000, loginConfig.Token, 1000000, UpdatedTitle, UpdatedDescription, UpdatedStatus)

	if err.Error() != TaskUpdateReturnedNotFoundStatusCode {
		log.Println(err)
		t.Fail()
	}
}

// TestUpdateTaskWithOtherUsersTask 他のユーザのタスクにたいして更新処理を試みて
// 404 Not Foundが返ってくることを確認する。
func TestUpdateTaskWithOtherUsersTask(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	TodoTitle := "TODO"
	TodoDescription := "TODO_DESCRIPTION"

	// 更新対象とするタスクを作成
	createdTask, err := CreateTask("http", testTarget, 8000, loginConfig.Token, TodoTitle, TodoDescription)

	if err != nil {
		t.Fail()
	}

	// test_userでログインして作成したタスクの更新を試みる
	loginConfig, err = Login("http", testTarget, 8000, "test_user", "test_password")

	UpdatedTitle := "UPDATED"
	UpdatedDescription := "UPDATED_DESCRIPTION"
	UpdatedStatus := "RUNNING"
	_, err = UpdateTask("http", testTarget, 8000, loginConfig.Token, createdTask.ID, UpdatedTitle, UpdatedDescription, UpdatedStatus)

	if err.Error() != TaskUpdateReturnedNotFoundStatusCode {
		log.Println(err)
		t.Fail()
	}
}

// TestUpdateTaskWithWrongAuth は誤った認証情報を与えた際に、
// エラーメッセージとしてTaskUpdateReturnedStatusCodeUnexpectedが返ってくることを期待する。

func TestUpdateTaskWithWrongAuth(t *testing.T) {
	testTarget := os.Getenv("TODO_TESTSERVER")
	if testTarget == "" {
		testTarget = "127.0.0.1"
	}

	loginConfig, err := Login("http", testTarget, 8000, "fujiwara", "fujiwara")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	TodoTitle := "TODO"
	TodoDescription := "TODO_DESCRIPTION"

	// 更新対象とするタスクを作成
	createdTask, err := CreateTask("http", testTarget, 8000, loginConfig.Token, TodoTitle, TodoDescription)

	UpdatedTitle := "UPDATED"
	UpdatedDescription := "UPDATED_DESCRIPTION"
	UpdatedStatus := "RUNNING"

	// 誤った認証トークン
	WrongToken := "wrong token"
	// 誤った認証トークンを与えて更新を実施
	_, err = UpdateTask("http", testTarget, 8000, WrongToken, createdTask.ID, UpdatedTitle, UpdatedDescription, UpdatedStatus)

	if err.Error() != TaskUpdateReturnedStatusCodeUnexpected {
		t.Fail()
	}

}
