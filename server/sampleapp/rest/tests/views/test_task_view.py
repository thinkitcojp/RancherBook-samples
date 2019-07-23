from django.urls import reverse
from django.test import TestCase
from django.contrib.auth.models import User
from rest_framework import status
from rest_framework.test import APIClient
from rest.services.task_manipulator import TaskManipulator
from rest.models.task import Task
from rest.models.task_status import TaskStatus

class TestTaskView(TestCase):

    fixtures = [
        'task_status',
        'sampleapp/fixtures/test_user'
    ]
    def test_create_task(self):
        """
        認証APIへのリクエストを投げた後に、
        その認証情報を使ってタスクの作成リクエストを投げる。
        """
        # test_userユーザを利用してテストを実行する。
        user = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": user.username,
            "password": "test_password",
        }

        response = client.post(path=path, data=data, format='json')

        # JWTを取得
        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)


        path = reverse('task')

        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)
        TITLE = "title"
        DESCRIPTION = "description"
        TASK_TITLE = '洗濯機から洗濯物を回収する'
        TASK_DESCRIPTION = '洗濯機から乾燥済みの洗濯物を回収する'
        data = {
            TITLE: TASK_TITLE,
            DESCRIPTION: TASK_DESCRIPTION,
        }

        response = client.post(path=path, data=data, format='json')

        # レスポンスコードが意図したものであることを確認する。
        self.assertEqual(response.status_code, status.HTTP_201_CREATED)

        # 返ってきたレスポンスの内容が意図したものであることを確認する
        self.assertEqual(type(response.data.get('id')), int)
        self.assertEqual(response.data.get(TITLE), TASK_TITLE)
        self.assertEqual(response.data.get(DESCRIPTION), TASK_DESCRIPTION)

    def test_get_tasks(self):
        """
        タスク一覧取得の処理を投げる。
        以下の結果を期待する。
        
        [
            {
                "id": "?",
                "title": "0",
                "description": "todo-description",
                "status": "TODO"
            },
            {
                "id": "?",
                "title": "1",
                "description": "todo-description",
                "status": "TODO"
            },
            {
                "id": "?",
                "title": "2",
                "description": "todo-description",
                "status": "TODO"
            }
        ]
        """
        # JWTを取得
        TODO_USER = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": TODO_USER.username,
            "password": "test_password",
        }

        response = client.post(path=path, data=data, format='json')

        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)

        
        # Taskを事前に作成
        TODO_DESCRIPTION = 'todo-description'
        manipulator = TaskManipulator()
        for i in range(3):
            manipulator.create(
                title=str(i),
                description=TODO_DESCRIPTION,
                user=TODO_USER
            )
        
        path = reverse('task')
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)
        response = client.get(path=path, format='json')

        self.assertEqual(response.status_code, status.HTTP_200_OK)
        tasks = response.data

        self.assertEqual(len(tasks), 3)

        for i in range(3):
            self.assertEqual(tasks[i].get('title'), str(i))
            self.assertEqual(tasks[i].get('description'), TODO_DESCRIPTION)
            self.assertEqual(tasks[i].get('status'), "TODO")
    
    def test_get_task_without_tasks(self):
        """
        Task一覧を取得するがリクエストするUserの
        Taskが一見も存在しない場合に空配列を受け取れるかを確認する。
        """
        # JWTを取得
        TODO_USER = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": TODO_USER.username,
            "password": "test_password",
        }

        response = client.post(path=path, data=data, format='json')

        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)


        # 別ユーザでタスクを追加
        OTHER_USER = User.objects.get(pk=1)
        manipulator = TaskManipulator()
        manipulator.create(
            title="todo",
            description="todo-descirption",
            user=OTHER_USER
        )

        # TODO_USERの持つTask一覧を取得
        path = reverse('task')
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)
        response = client.get(path=path, format='json')

        # 200 OKが返ってくることを確認
        self.assertEqual(response.status_code, status.HTTP_200_OK)

        # tasksが空配列であることを確認する
        tasks = response.data
        self.assertEqual(len(tasks), 0)
    
    def test_get_task(self):
        """
        task_idを指定して当該Taskを取得する
        """
        # JWTを取得
        TODO_USER = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": TODO_USER.username,
            "password": "test_password",
        }

        response = client.post(path=path, data=data, format='json')

        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)

        # 取得対象Taskを追加
        TODO_TITLE = 'todo-title'
        TODO_DESCRIPTION = 'todo-description'

        manipulator = TaskManipulator()
        TODO_TASK = manipulator.create(
            title=TODO_TITLE,
            description=TODO_DESCRIPTION,
            user=TODO_USER
        )

        # 追加したTaskを取得
        path = reverse('specific_task', kwargs={ 'task_id': TODO_TASK.id })
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)
        response = client.get(path=path, format='json')

        # レスポンスコードが200 OKであることを確認
        self.assertEqual(response.status_code, status.HTTP_200_OK)

        # 返ってきたTaskが事前に作成したTask(TODO_TASK)と一致することを確認
        tasks = response.data
        self.assertEqual(len(tasks), 1)
        self.assertEqual(tasks[0].get('id'), TODO_TASK.id)
        self.assertEqual(tasks[0].get('title'), TODO_TASK.title)
        self.assertEqual(tasks[0].get('description'), TODO_TASK.description)
        self.assertEqual(tasks[0].get('status'), "TODO")

    def test_get_task_with_no_task(self):
        """
        存在しないtask_idを指定して空配列が返ってくること、
        404 NOT FOUNDとなることを確認する。
        """
        # JWTを取得
        TODO_USER = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": TODO_USER.username,
            "password": "test_password",
        }

        response = client.post(path=path, data=data, format='json')

        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)

        # 存在しないtask_idを指定してリクエスト
        TASK_ID = 10000
        path = reverse('specific_task', kwargs={ 'task_id': TASK_ID })
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)
        response = client.get(path=path, format='json')

        # レスポンスコードが404 NOT FOUNDであることを確認
        self.assertEqual(response.status_code, status.HTTP_404_NOT_FOUND)

        # 空配列が返ってきていることを確認
        tasks = response.data
        self.assertEqual(len(tasks), 0)

    def test_update_task(self):
        """
        存在するtask_idを指定して更新を行い
        更新済みのTaskの情報が返ってくること、
        200 OKとなることを確認する。
        以下のような形式のメッセージが返ってくることを期待する。
        
        {
            "id": "?",
            "title": "0",
            "description": "todo-description",
            "status": "TODO"
        }
        """
        # JWTを取得
        TODO_USER = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": TODO_USER.username,
            "password": "test_password",
        }

        response = client.post(path=path, data=data, format='json')

        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)

        # Taskを事前に作成
        TODO_TITLE = 'todo-title'
        TODO_DESCRIPTION = 'todo-description'
        TODO_USER = User.objects.get(pk=2)
        manipulator = TaskManipulator()
        task = manipulator.create(
            title=TODO_TITLE,
            description=TODO_DESCRIPTION,
            user=TODO_USER
        )

        # 追加したタスクを更新
        path = reverse('specific_task', kwargs={ 'task_id': task.id })
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)

        TITLE = "title"
        DESCRIPTION = "description"
        STATUS = "status"

        UPDATED_TASK_TITLE = '洗濯機から洗濯物を回収する'
        UPDATED_TASK_DESCRIPTION = '洗濯機から乾燥済みの洗濯物を回収する'
        UPDATED_TASK_STATUS = TaskStatus.objects.get(pk=2).name
        data = {
            TITLE: UPDATED_TASK_TITLE,
            DESCRIPTION: UPDATED_TASK_DESCRIPTION,
            STATUS: UPDATED_TASK_STATUS
        }

        response = client.patch(path=path, data=data,format='json')

        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.assertEqual(response.data.get("id"), task.id)
        self.assertEqual(response.data.get("title"), UPDATED_TASK_TITLE)
        self.assertEqual(response.data.get("description"), UPDATED_TASK_DESCRIPTION)
        self.assertEqual(response.data.get("status"), UPDATED_TASK_STATUS)

    def test_update_task_with_wrong_status(self):
        """
        存在しないTaskStatusを指定して更新を試みて
        400 Bad Requestが返ってくることを確認する
        """
        # JWTを取得
        TODO_USER = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": TODO_USER.username,
            "password": "test_password",
        }

        response = client.post(path=path, data=data, format='json')

        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)

        # Taskを事前に作成
        TODO_TITLE = 'todo-title'
        TODO_DESCRIPTION = 'todo-description'
        TODO_USER = User.objects.get(pk=2)
        manipulator = TaskManipulator()
        task = manipulator.create(
            title=TODO_TITLE,
            description=TODO_DESCRIPTION,
            user=TODO_USER
        )

        # 追加したタスクを更新
        path = reverse('specific_task', kwargs={ 'task_id': task.id })
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)

        TITLE = "title"
        DESCRIPTION = "description"
        STATUS = "status"

        UPDATED_TASK_TITLE = '洗濯機から洗濯物を回収する'
        UPDATED_TASK_DESCRIPTION = '洗濯機から乾燥済みの洗濯物を回収する'
        UPDATED_TASK_STATUS = "WRONG"
        data = {
            TITLE: UPDATED_TASK_TITLE,
            DESCRIPTION: UPDATED_TASK_DESCRIPTION,
            STATUS: UPDATED_TASK_STATUS
        }

        response = client.patch(path=path, data=data,format='json')
        
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)
        self.assertIsNone(response.data)
    
    def test_update_task_without_task(self):
        """
        存在しないTaskを指定して更新を試みて
        404 Not Foundが返ってくることを確認する
        """
        # JWTを取得
        TODO_USER = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": TODO_USER.username,
            "password": "test_password",
        }

        response = client.post(path=path, data=data, format='json')

        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)

        # 明らかに存在しないtask_idを指定
        path = reverse('specific_task', kwargs={ 'task_id': 100000 })
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)

        TITLE = "title"
        DESCRIPTION = "description"
        STATUS = "status"

        UPDATED_TASK_TITLE = '洗濯機から洗濯物を回収する'
        UPDATED_TASK_DESCRIPTION = '洗濯機から乾燥済みの洗濯物を回収する'
        UPDATED_TASK_STATUS = TaskStatus.objects.get(pk=2).name
        data = {
            TITLE: UPDATED_TASK_TITLE,
            DESCRIPTION: UPDATED_TASK_DESCRIPTION,
            STATUS: UPDATED_TASK_STATUS
        }

        response = client.patch(path=path, data=data,format='json')
        
        self.assertEqual(response.status_code, status.HTTP_404_NOT_FOUND)
        self.assertIsNone(response.data)

    def test_update_task_with_other_user(self):
        """
        更新しようとしたTaskは存在するものの、
        リクエストしたUserに紐付かないものである場合に、
        404 Not Foundが返ってくることを確認する。
        """
        # JWTを取得
        TODO_USER = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": TODO_USER.username,
            "password": "test_password",
        }

        response = client.post(path=path, data=data, format='json')

        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)

        # Taskを事前に作成
        TODO_TITLE = 'todo-title'
        TODO_DESCRIPTION = 'todo-description'
        USER = User.objects.get(pk=1)
        manipulator = TaskManipulator()
        task = manipulator.create(
            title=TODO_TITLE,
            description=TODO_DESCRIPTION,
            user=USER
        )

        # 追加したタスクを更新
        path = reverse('specific_task', kwargs={ 'task_id': task.id })
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)

        TITLE = "title"
        DESCRIPTION = "description"
        STATUS = "status"

        UPDATED_TASK_TITLE = '洗濯機から洗濯物を回収する'
        UPDATED_TASK_DESCRIPTION = '洗濯機から乾燥済みの洗濯物を回収する'
        UPDATED_TASK_STATUS = TaskStatus.objects.get(pk=2).name
        data = {
            TITLE: UPDATED_TASK_TITLE,
            DESCRIPTION: UPDATED_TASK_DESCRIPTION,
            STATUS: UPDATED_TASK_STATUS
        }

        response = client.patch(path=path, data=data,format='json')

        self.assertEqual(response.status_code, status.HTTP_404_NOT_FOUND)
        self.assertIsNone(response.data)

    def test_delete_task(self):
        # JWTを取得
        TODO_USER = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": TODO_USER.username,
            "password": "test_password",
        }
        response = client.post(path=path, data=data, format='json')

        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)

        # Taskを事前に作成
        TODO_TITLE = 'todo-title'
        TODO_DESCRIPTION = 'todo-description'
        manipulator = TaskManipulator()
        task = manipulator.create(
            title=TODO_TITLE,
            description=TODO_DESCRIPTION,
            user=TODO_USER
        )

        # 追加したタスクを削除
        path = reverse('specific_task', kwargs={ 'task_id': task.id })
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)

        response = client.delete(path=path, format='json')

        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.assertEqual(response.data.get('id'), task.id)
        self.assertEqual(response.data.get('title'), task.title)
        self.assertEqual(response.data.get('description'), task.description)

    def test_delete_task_without_task(self):
        # JWTを取得
        TODO_USER = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": TODO_USER.username,
            "password": "test_password",
        }
        response = client.post(path=path, data=data, format='json')

        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)        


        # 存在しないTaskの削除を試みる
        path = reverse('specific_task', kwargs={ 'task_id': 1000 })
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)

        response = client.delete(path=path, format='json')

        self.assertEqual(response.status_code, status.HTTP_404_NOT_FOUND)
    
    def test_delete_task_with_other_users(self):
        # JWTを取得
        TODO_USER = User.objects.get(pk=2)
        client = APIClient()

        path = reverse('auth')

        data = {
            "username": TODO_USER.username,
            "password": "test_password",
        }
        response = client.post(path=path, data=data, format='json')

        token = response.data.get('token')
        self.assertEqual(response.status_code, status.HTTP_200_OK)        

        # 別Userに紐づくTaskを事前に作成
        TODO_TITLE = 'todo-title'
        TODO_DESCRIPTION = 'todo-description'
        manipulator = TaskManipulator()
        user = User.objects.get(pk=1)

        task = manipulator.create(
            user=user,
            title=TODO_TITLE,
            description=TODO_DESCRIPTION
        )

        # Taskが作成されていることを確認
        self.assertIsNotNone(task)

        # 他Userに紐づくTaskを削除しようと試みる
        path = reverse('specific_task', kwargs={ 'task_id': task.id })
        client = APIClient()
        client.credentials(HTTP_AUTHORIZATION='JWT ' + token)

        response = client.delete(path=path, format='json')

        self.assertEqual(response.status_code, status.HTTP_404_NOT_FOUND)