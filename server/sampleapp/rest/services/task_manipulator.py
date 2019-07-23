from django.db import transaction
from rest.models.task import Task
from rest.models.task_status import TaskStatus
from rest.models.task_history import TaskHistory
import logging

class TaskManipulator():

    def __init__(self):
        self.logger = logging.getLogger(__name__)

    @transaction.atomic
    def create(self, title, user, description=None):        
        """
        Taskのレコードを作成するためのメソッドです。
        """
        task = Task(
            title=title,
            description=description,
            user = user
        )

        task.save()

        # Task作成の時点では、特定の状態に固定する。
        TaskHistory(
            status=TaskStatus.objects.get(pk=1),
            task=task
        ).save()

        return task

    def get_task(self, user, task_id):
        """
        特定のIDを持つタスクを取得するためのメソッドです。
        Taskと最新のTaskStatusを返します。
        指定したtask_idが存在しても、
        task_idに対応するTaskのuserが
        一致しない限りはNoneを返します。
        task_idに対応するTaskが存在しない場合も
        Task, TaskStatusともにNoneを返します。
        """
        try:
            # DoesNotExist例外が発生した場合もNone, Noneを返す
            task = Task.objects.get(id=task_id)
        except Task.DoesNotExist:
            return None, None

        if task.user != user:
            # 取得できたTaskに対応するUserが異なる場合もNone, Noneを返す
            return None, None
        latest =  TaskHistory.objects.filter(task=task).latest('created_at')
        # tasks = Task.objects.filter(id=task_id).filter(user=user).\
        return task, latest.status
            
    def get_tasks(self, user):
        """
        特定のユーザが持つ全タスクを取得するためのメソッドです。
        [(TaskA, TaskAの最新TaskStatus),(TaskB,TaskBの最新TaskStatus),...]
        の様な形で格納された値が返ります。
        ユーザに対応するタスクがない場合は空配列が返ります。
        """
        tasks = Task.objects.filter(user=user).order_by('id')

        task_list = []
        for task in tasks:
            # ここの処理をget_taskに寄せるかはお好みで
            latest_status = TaskHistory.objects.filter(task=task).latest('created_at').status
            task_list.append((task, latest_status))
        
        return task_list
    
    @transaction.atomic
    def update_task(self, user, task_id, title=None, description=None, status=None):

        if task_id is None:
            self.logger.error("task_idが指定されていません。")
            return None, None

        task, _ = self.get_task(user=user, task_id=task_id)

        if task is None:
            self.logger.warning("指定したUserに紐づくTaskが存在しません。")
            return None, None
        
        try:
            task_status = TaskStatus.objects.get(name=status)
        except TaskStatus.DoesNotExist:
            self.logger.warning("指定したTaskStatusが存在しません。")
            return task, None

        if title is not None:
            task.title = title

        if description is not None:
            task.description = description

        task_history = TaskHistory(
            status=task_status,
            task=task
        ) 

        task.save()
        task_history.save() 

        return task, task_history.status
    
    @transaction.atomic
    def delete_task(self, user, task_id):

        try:
            task = Task.objects.filter(user=user).get(pk=task_id)
        except Task.DoesNotExist:
            return None
        
        task_id = task.id
        task.delete()

        return task_id
