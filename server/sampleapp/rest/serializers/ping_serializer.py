from rest_framework import serializers


class PingSerializer(serializers.Serializer):
    message = serializers.CharField(default='pong')