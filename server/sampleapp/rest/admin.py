from django.contrib import admin
from rest.models.task import Task
from rest.models.task_status import TaskStatus
from rest.models.task_history import TaskHistory

# Register your models here.

admin.site.register(Task)
admin.site.register(TaskStatus)
admin.site.register(TaskHistory)