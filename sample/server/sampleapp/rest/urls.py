from django.urls import path
from django.contrib import admin
from rest_framework_jwt.views import obtain_jwt_token

from rest.views.ping_view import PingView
from rest.views.task_view import TaskView

urlpatterns = [
    path('ping', PingView.as_view(), name='pingpong'),
    path('auth', obtain_jwt_token, name='auth'),
    path('task', TaskView.as_view(), name='task'),
    path('task/<int:task_id>', TaskView.as_view(), name='specific_task'),
]