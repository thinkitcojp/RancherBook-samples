from django.test import TestCase
from django.contrib.auth.models import User
from rest.models.task import Task
from rest.models.task_status import TaskStatus
from rest.models.task_history import TaskHistory
from rest.services.task_manipulator import TaskManipulator

class TestTaskManipulator(TestCase):

    # フィクスチャとしてtask_statusと、
    # sampleapp側のtest_userを利用する
    fixtures = [
        'task_status',
        'sampleapp/fixtures/test_user'
    ]
    def test_create_task(self):

        # 作成されることを期待しているTaskの情報を
        # 定数的に記述しておく。
        TODO_TITLE = "test-task"
        TODO_USER = User.objects.get(pk=2) 
        TODO_DESCRIPTION = "test-description"

        manipulator = TaskManipulator()

        task = manipulator.create(
            title=TODO_TITLE,
            user=TODO_USER,
            description=TODO_DESCRIPTION
        )

        self.assertIsNotNone(task.id)
        self.assertEqual(task.title, TODO_TITLE)
        self.assertEqual(task.description, TODO_DESCRIPTION)
        self.assertEqual(task.user, TODO_USER)

        # Task作成時に一緒に作成されているTaskStatusが期待したものかを確認する。

        EXPECTED_TASK_STATUS = TaskStatus.objects.get(pk=1)
        task_history = TaskHistory.objects.filter(task=task)
        self.assertEqual(len(task_history), 1)
        self.assertEqual(task_history[0].status, EXPECTED_TASK_STATUS)

    def test_get_task(self):
        # 当該Userに紐づくTaskが1つある場合に
        # Task, TaskStatus(TODO状態)が
        # 返ってくることを確認します。

        manipulator = TaskManipulator()
        
        # 事前にタスクを作成
        TODO_TITLE = "test-task"
        TODO_USER = User.objects.get(pk=2) 
        TODO_DESCRIPTION = "test-description"
        TODO_STATUS = TaskStatus.objects.get(pk=1)

        manipulator = TaskManipulator()

        TODO_TASK = manipulator.create(
            title=TODO_TITLE,
            user=TODO_USER,
            description=TODO_DESCRIPTION
        )
        
        task, status = manipulator.get_task(user=TODO_USER, task_id=TODO_TASK.id)

        self.assertEqual(task.title, TODO_TITLE)
        self.assertEqual(task.user, TODO_USER)
        self.assertEqual(task.description, TODO_DESCRIPTION)
        self.assertEqual(status, TODO_STATUS) 
    
    def test_get_task_updated(self):
        # 当該タスクに対してより新しいTaskStatusが追加された場合に
        # get_taskで最新のTaskStatusが取得できることを確認します。

        manipulator = TaskManipulator()
        
        # 事前にタスクを作成
        TODO_TITLE = "test-task"
        TODO_USER = User.objects.get(pk=2) 
        TODO_DESCRIPTION = "test-description"
        TODO_STATUS = TaskStatus.objects.get(pk=1)

        manipulator = TaskManipulator()

        TODO_TASK = manipulator.create(
            title=TODO_TITLE,
            user=TODO_USER,
            description=TODO_DESCRIPTION
        )
        
        # TaskHistoryにRUNNNINGの状態のものを追加
        RUNNING_STATUS = TaskStatus.objects.get(pk=2)
        TaskHistory(
            status=RUNNING_STATUS,
            task=TODO_TASK
        ).save()

        task, status = manipulator.get_task(user=TODO_USER, task_id=TODO_TASK.id) 

        self.assertEqual(task, TODO_TASK)
        self.assertEqual(status, RUNNING_STATUS)
    
    def test_get_task_not_exist(self):
        # 存在しないTaskのidを指定して
        # None, Noneが返ってくることを期待する。
        manipulator = TaskManipulator()

        TODO_USER = User.objects.get(pk=2) 

        # 明らかに存在しないTaskのidを指定する
        task, status = manipulator.get_task(user=TODO_USER, task_id=100)

        self.assertIsNone(task)
        self.assertIsNone(status)
    
    def test_get_task_with_other_users(self):
        # 指定したidのTaskは存在するが、
        # 他UserのTaskである場合。
        manipulator = TaskManipulator()
        
        # 事前にタスクを作成
        TODO_TITLE = "test-task"
        TODO_USER = User.objects.get(pk=2)
        TODO_DESCRIPTION = "test-description"

        manipulator = TaskManipulator()

        TODO_TASK = manipulator.create(
            title=TODO_TITLE,
            user=TODO_USER,
            description=TODO_DESCRIPTION
        )

        USER = User.objects.get(pk=1) # Task作成時のUserとは別Userを指定
        task, status = manipulator.get_task(user=USER, task_id=TODO_TASK.id)
        
        self.assertIsNone(task)
        self.assertIsNone(status)

    def test_get_tasks(self):

        manipulator = TaskManipulator()

        # 事前にタスクを作成
        TODO_USER = User.objects.get(pk=2)
        TODO_DESCRIPTION = "test-description"
        TODO_STATUS = TaskStatus.objects.get(pk=1)

        TODO_TASKS = []
        for i in range(3):
            task = manipulator.create(
                title=str(i),
                user=TODO_USER,
                description=TODO_DESCRIPTION
            )
            TODO_TASKS.append(task)

        task_list = manipulator.get_tasks(user=TODO_USER)

        self.assertEqual(len(task_list), 3)

        for i in range(3):
            self.assertEqual(task_list[i][0].id, TODO_TASKS[i].id)
            self.assertEqual(task_list[i][0].title, TODO_TASKS[i].title)
            self.assertEqual(task_list[i][0].description, TODO_TASKS[i].description)
            self.assertEqual(task_list[i][0].user, TODO_USER)
            self.assertEqual(task_list[i][1], TODO_STATUS)
    
    def test_get_tasks_without_tasks(self):

        manipulator = TaskManipulator()
        TODO_USER = User.objects.get(pk=2)

        task_list = manipulator.get_tasks(user=TODO_USER)

        self.assertEqual(len(task_list), 0)

    def test_get_tasks_with_other_users(self):

        manipulator = TaskManipulator()

        # 事前にタスクを作成
        TODO_USER = User.objects.get(pk=2)
        TODO_DESCRIPTION = "test-description"
        TODO_STATUS = TaskStatus.objects.get(pk=1)

        TODO_TASKS = []
        for i in range(3):
            task = manipulator.create(
                title=str(i),
                user=TODO_USER,
                description=TODO_DESCRIPTION
            )
            TODO_TASKS.append(task)

        USER = User.objects.get(pk=1)

        task_list = manipulator.get_tasks(user=USER)

        self.assertEqual(len(task_list), 0)
    
    def test_update_task(self):
        """
        正常にTaskおよびTaskStatusの更新ができることを確認する
        """

        manipulator = TaskManipulator()

        # 事前にタスクを作成
        TODO_USER = User.objects.get(pk=2)
        TODO_TITLE = "todo"
        TODO_DESCRIPTION = "test-description"
        
        task = manipulator.create(
            title=TODO_TITLE,
            user=TODO_USER,
            description=TODO_DESCRIPTION
        ) 

        # 更新処理タスクが期待したとおりに更新されているかを確認
        RUNNING_STATUS = TaskStatus.objects.get(pk=2)
        UPDATED_TITLE = "updated"
        UPDATED_DESCRIPTION = "updated-description"

        updated_task, updated_task_status = manipulator.update_task(
            user=TODO_USER,
            task_id=task.id,
            title=UPDATED_TITLE,
            description=UPDATED_DESCRIPTION,
            status=RUNNING_STATUS.name
        )

        self.assertEqual(updated_task.id, task.id)
        self.assertEqual(updated_task.title, UPDATED_TITLE)
        self.assertEqual(updated_task.description, UPDATED_DESCRIPTION)
        self.assertEqual(updated_task.user, TODO_USER)
        self.assertEqual(updated_task_status.name, RUNNING_STATUS.name)
    
    def test_update_task_with_status_only(self):
        """
        TaskStatusのみをアップデート
        """
        manipulator = TaskManipulator()

        # 事前にタスクを作成
        TODO_USER = User.objects.get(pk=2)
        TODO_TITLE = "todo"
        TODO_DESCRIPTION = "test-description"
        
        task = manipulator.create(
            title=TODO_TITLE,
            user=TODO_USER,
            description=TODO_DESCRIPTION
        ) 

        # 更新処理タスクが期待したとおりに更新されているかを確認
        RUNNING_STATUS = TaskStatus.objects.get(pk=2)

        updated_task, updated_task_status = manipulator.update_task(
            user=TODO_USER,
            task_id=task.id,
            status=RUNNING_STATUS.name
        )

        self.assertEqual(updated_task.id, task.id)
        self.assertEqual(updated_task.title, TODO_TITLE)
        self.assertEqual(updated_task.description, TODO_DESCRIPTION)
        self.assertEqual(updated_task.user, TODO_USER)
        self.assertEqual(updated_task_status.name, RUNNING_STATUS.name)
    
    def test_update_task_with_wrong_status(self):
        
        manipulator = TaskManipulator()

        # 事前にタスクを作成
        TODO_USER = User.objects.get(pk=2)
        TODO_TITLE = "todo"
        TODO_DESCRIPTION = "test-description"
        
        task = manipulator.create(
            title=TODO_TITLE,
            user=TODO_USER,
            description=TODO_DESCRIPTION
        ) 

        
        # 更新処理タスクが期待したとおりに更新されているかを確認
        UPDATED_TITLE = "updated"
        UPDATED_DESCRIPTION = "updated-description"
        WRONG_STATUS = "wrong"

        updated_task, updated_task_status = manipulator.update_task(
            user=TODO_USER,
            task_id=task.id,
            title=UPDATED_TITLE,
            description=UPDATED_DESCRIPTION,
            status=WRONG_STATUS
        )

        # Taskが更新されていないこと、
        # TaskStatusとしてNoneが返ってきていることを確認
        self.assertEqual(updated_task, task)
        self.assertIsNone(updated_task_status)

    def test_update_task_status_without_task(self):
        """
        指定したtask_idに対応するTaskが存在しない場合
        """
        manipulator = TaskManipulator()

        TODO_USER = User.objects.get(pk=2)
        # 更新処理タスクが期待したとおりに更新されているかを確認
        UPDATED_TITLE = "updated"
        UPDATED_DESCRIPTION = "updated-description"
        RUNNING_STATUS = TaskStatus.objects.get(pk=2)

        updated_task, updated_task_status = manipulator.update_task(
            user=TODO_USER,
            task_id=140000000,
            title=UPDATED_TITLE,
            description=UPDATED_DESCRIPTION,
            status=RUNNING_STATUS
        )

        # Task, TaskStatusともにNoneであることを確認
        self.assertIsNone(updated_task)
        self.assertIsNone(updated_task_status)
    
    def test_update_task_status_with_other_user(self):
        """
        他のユーザに紐づくタスクを指定した場合
        に更新されずにNone, Noneが返ってくることを確認
        """

        manipulator = TaskManipulator()

        # 事前にタスクを作成
        TODO_USER = User.objects.get(pk=2)
        TODO_TITLE = "todo"
        TODO_DESCRIPTION = "test-description"
        
        task = manipulator.create(
            title=TODO_TITLE,
            user=TODO_USER,
            description=TODO_DESCRIPTION
        ) 

        # 更新処理タスクが期待したとおりに更新されているかを確認
        USER = User.objects.get(pk=1)
        RUNNING_STATUS = TaskStatus.objects.get(pk=2)
        UPDATED_TITLE = "updated"
        UPDATED_DESCRIPTION = "updated-description"

        updated_task, updated_task_status = manipulator.update_task(
            user=USER,
            task_id=task.id,
            title=UPDATED_TITLE,
            description=UPDATED_DESCRIPTION,
            status=RUNNING_STATUS.name
        )

        self.assertIsNone(updated_task)
        self.assertIsNone(updated_task_status)

        task = Task.objects.get(pk=task.id)
        self.assertEqual(task.title, TODO_TITLE)
        self.assertEqual(task.description, TODO_DESCRIPTION)
    
    def test_update_task_without_task_id(self):

        manipulator = TaskManipulator()
        user = User.objects.get(pk=2)

        task, task_status = manipulator.update_task(
            user=user,
            task_id=None
        )

        self.assertIsNone(task)
        self.assertIsNone(task_status)

    
    def test_delete_task(self):
        """
        Taskが正常に削除され、かつ、
        削除されたTaskのIDが返ってくることを期待する。
        """
        
        manipulator = TaskManipulator()
        TODO_USER = User.objects.get(pk=2)
        TODO_TITLE = "todo"
        TODO_DESCRIPTION = "todo-description"

        TODO_TASK = manipulator.create(
            user=TODO_USER,
            title=TODO_TITLE,
            description=TODO_DESCRIPTION
        )

        task_id = manipulator.delete_task(
            user=TODO_USER,
            task_id=TODO_TASK.id
        )
        

        self.assertEqual(task_id, TODO_TASK.id)

        # 正常に削除できていることを確認する。
        task, task_status = manipulator.get_task(user=TODO_USER, task_id=task_id)

        self.assertIsNone(task)
        self.assertIsNone(task_status)



    def test_delete_task_with_no_task(self):
        """
        存在しないTaskを指定する。
        Noneが返ってくることを期待する。
        """
        manipulator = TaskManipulator()
        TODO_USER = User.objects.get(pk=2)

        task_id = manipulator.delete_task(
            user=TODO_USER,
            task_id = 100
        )

        self.assertIsNone(task_id)
    
    def test_delete_task_with_other_users(self):
        """
        存在するが、他Userに紐付いているTaskを
        削除しようとしてNoneが返ってくることを期待する。
        """
        manipulator = TaskManipulator()
        TODO_USER = User.objects.get(pk=2)
        TODO_TITLE = "todo"
        TODO_DESCRIPTION = "todo-description"

        TODO_TASK = manipulator.create(
            user=TODO_USER,
            title=TODO_TITLE,
            description=TODO_DESCRIPTION
        )

        user = User.objects.get(pk=1)
        task_id = manipulator.delete_task(
            user=user,
            task_id=TODO_TASK.id
        )

        self.assertIsNone(task_id)

        # 当該タスクが削除されていないことを確認する。
        task, _ = manipulator.get_task(user=TODO_USER, task_id=TODO_TASK.id)

        self.assertEqual(task, TODO_TASK)

