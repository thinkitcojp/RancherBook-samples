from rest_framework.views import APIView
from rest_framework.response import Response

from rest.serializers.ping_serializer import PingSerializer

class PingView(APIView):
    """
    GETアクセスに対してpongの文言をボディに含んだ
    HTTPレスポンスメッセージを返すAPIです。
    """

    authentication_classes = ()
    permission_classes = ()

    def get(self, request, format=None):
        message = {
            "message": "pong",
        }
        serializer = PingSerializer(message)
        return Response(serializer.data)