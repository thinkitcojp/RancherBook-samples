from rest_framework import serializers
from rest.models.task import Task

class TaskSerializer(serializers.Serializer):
    id = serializers.IntegerField(required=False)
    title = serializers.CharField(required=True, max_length=100)
    description = serializers.CharField(required=False)
    created_at = serializers.DateTimeField(required=False)
    status = serializers.CharField(max_length=20, required=False)