from django.urls import reverse
from django.test import TestCase
from django.contrib.auth.models import User
from rest_framework import status
from rest_framework.test import APIClient

class TestAuth(TestCase):
    def test_auth_view(self):
        """
        JWTによる認証APIにリクエストを投げた際に、
        HTTPステータスコードとして200 OK、
        レスポンスボディとして以下のようなJSONを期待する。
        {
            "token": "何かしらの文字列"
        }
        """

        username = "test"
        password = "testpassword"
        email = "test@example.com"
        user = User.objects.create_user(username=username, email=email, password=password)
        user.save()

        client = APIClient()

        path = reverse('auth')

        data = {
            "username": username,
            "password": password,
        }

        response = client.post(path=path, data=data, format='json')

        # 200 OKがステータスコードとして返されること
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        # レスポンスボディにJSON形式でtokenキーに対する値が存在すること
        self.assertIsNotNone(response.data.get('token'))
        # レスポンスボディのtokenキーに含まれている値が文字列であること(蛇足か?)
        self.assertIsInstance(response.data.get('token'), str)

    def test_auth_view_failure(self):
        """
        JWTによる認証APIにリクエストを投げた際に、
        HTTPステータスコードとして400 Bad Requestを期待する
        """

        username = "test"
        password = "testpassword"
        email = "test@example.com"
        user = User.objects.create_user(username=username, email=email, password=password)
        user.save()

        client = APIClient()

        path = reverse('auth')

        data = {
            "username": username,
            "password": password+"hogehoge", # 絶対に失敗するようにしておく
        }

        response = client.post(path=path, data=data, format='json')

        # 200 OKがステータスコードとして返されること
        self.assertEqual(response.status_code, status.HTTP_400_BAD_REQUEST)