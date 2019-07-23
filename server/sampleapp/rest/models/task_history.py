from django.db import models
from .task import Task
from .task_status import TaskStatus


class TaskHistory(models.Model):
    created_at = models.DateTimeField(auto_now_add=True)
    status = models.ForeignKey(TaskStatus, null=True, on_delete=models.SET_NULL)
    task = models.ForeignKey(Task, on_delete=models.CASCADE)

