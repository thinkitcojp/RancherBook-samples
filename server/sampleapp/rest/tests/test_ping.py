from django.urls import reverse
from django.test import TestCase
from rest_framework import status
from rest_framework.test import APIClient

class TestPing(TestCase):
    def test_ping_view(self):
        """
        ping-pong APIにリクエストを投げた際に、
        HTTPステータスコードとして200 OK、
        レスポンスボディとして以下のようなJSONを期待する。
        {
            "message": "pong"
        }
        """
        client = APIClient()

        path = reverse('pingpong')
        response = client.get(path)

        # self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.assertEqual(response.data.get('message'), "pong")