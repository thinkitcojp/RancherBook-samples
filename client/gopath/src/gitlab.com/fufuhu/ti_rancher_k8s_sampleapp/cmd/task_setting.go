package cmd

import (
	"errors"
	"log"
)

// TaskRequestSetting はToDoクライアントが
// タスクを作成する際のタスクの情報を格納します。
type TaskRequestSetting struct {
	// Title 作成したいタスクの名前
	Title func() (string, error)
	// Description 作成したいタスクの概要
	Description func() (string, error)
	// ID 取得したいタスク
	ID func() (int, error)
	// Status 取得したいタスクのステータス
	Status func() (string, error)
}

// SettingTaskTitleNotFound はタスクの名前がtitleオプションで設定
// されていない場合に発生するエラーに含まれるエラーメッセージです。
const SettingTaskTitleNotFound = "タスクの名前が指定されていません"

// SettingTaskDescriptionNotFound はタスクの名前がdesriptionオプションで設定
// されていない場合に発生するエラーに含まれるエラーメッセージです。
const SettingTaskDescriptionNotFound = "タスクの概要が指定されていません"

var taskRequestSetting TaskRequestSetting

func init() {
	taskRequestSetting.Title = func() (string, error) {
		title, err := updateCmd.Flags().GetString("title")
		if err != nil {
			log.Println(err)
		}
		if title != "" {
			return title, err
		}

		title, err = createCmd.Flags().GetString("title")
		if err != nil {
			log.Println(err)
		}
		if title != "" {
			return title, err
		}

		if title == "" {
			err = errors.New(SettingTaskTitleNotFound)
		}

		return title, err
	}

	taskRequestSetting.Description = func() (string, error) {
		description, err := createCmd.Flags().GetString("description")
		if err != nil {
			log.Println(err)
		} else if description != "" {
			return description, err
		}

		// タスクの更新時の--descriptionを取得
		description, err = updateCmd.Flags().GetString("description")
		if err != nil {
			log.Println(err)
		}

		if description == "" {
			err = errors.New(SettingTaskDescriptionNotFound)
		}
		return description, err
	}

	taskRequestSetting.ID = func() (int, error) {
		// デフォルトでIDの値が0で入ることを逆手に取る
		getID, err := getCmd.Flags().GetInt("id")
		if err != nil {
			log.Println(err)
		}
		deleteID, err := deleteCmd.Flags().GetInt("id")
		if err != nil {
			log.Println(err)
		}
		updateID, err := updateCmd.Flags().GetInt("id")
		if err != nil {
			log.Println(err)
		}

		id := getID + deleteID + updateID

		return id, err
	}

	taskRequestSetting.Status = func() (string, error) {
		// タスク更新時の--statusオプションの値を取得
		status, err := updateCmd.Flags().GetString("status")

		if err != nil {
			log.Println(err)
			return "", err
		}

		return status, err
	}
}
