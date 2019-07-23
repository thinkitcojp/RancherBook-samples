from rest_framework.views import APIView
from rest_framework import status, permissions
from rest_framework.response import Response
from rest_framework.parsers import JSONParser
from rest.services.task_manipulator import TaskManipulator
from rest.serializers.task_serializer import TaskSerializer

from rest.models.task import Task

from django.contrib.auth.models import User


class TaskView(APIView):
    """
    Task周りのCRUDを管理するためのクラスです。
    """
    # authentication_classes = ()
    # permission_classes = ()

    def get(self, request, task_id=None):
        """
        1) task_idが指定されている場合
        ・対応するidを持つTaskがリクエストしたUserのものである場合は
          当該Task1つのみを当該Taskに含んだを返す。
        ・対応するTaskが存在しない場合は空配列と404 Not Foundを返す
        ・対応するTaskが存在してもリクエストしたUserのものではない場合も
        404 Not Foundとから配列を返す。
        2) task_idが指定されていない場合
        ・リクエストしたUserのTaskの配列を返す(200 OK)

        <返すレスポンスの形式>
        [
            {
                "id": Task.id,
                "title": Task.title,
                "description": Task.description,
                "status": Taskの最新TaskStatus.name
            },...
        ]
        """
        user = request.user
        manipulator = TaskManipulator()

        tasks = []
        if task_id is None:
            ts = manipulator.get_tasks(user=user)
        
            for task, task_status in ts:
                tasks.append({
                    "id": task.id,
                    "title": task.title,
                    "description": task.description,
                    "status": task_status.name,
                })
            

        else:
            task, task_status = manipulator.get_task(user=user, task_id=task_id)

            if task is None or task_status is None:
                # 一方のみがNone出ないことはアプリケーションの仕様上
                # ありえないが念の為
                serializer = TaskSerializer(tasks, many=True)
                return Response(data=serializer.data, status=status.HTTP_404_NOT_FOUND)

            tasks.append({
                "id": task.id,
                "title": task.title,
                "description": task.description,
                "status": task_status.name,
            })
        
        serializer = TaskSerializer(tasks, many=True)
        return Response(data=serializer.data, status=status.HTTP_200_OK)

    # JWTからユーザ情報
    # http://racchai.hatenablog.com/entry/2016/05/08/070000
    def post(self, request):
        """
        メッセージボディとして以下を期待。
        {
            "title": Taskのtitle,
            "description": Taskの内容,
        }
        レスポンスメッセージとして以下を作成する。
        {
            "id": TaskのID,
            "title": Taskの名前,
            "description": Taskの詳細,
            "created_at": Taskの生成日時
        }
        """
        # https://qiita.com/xKxAxKx/items/60e8fb93d6bbeebcf065#%E3%83%AD%E3%82%B0%E3%82%A4%E3%83%B3%E3%83%A6%E3%83%BC%E3%82%B6%E3%81%AE%E6%83%85%E5%A0%B1%E5%8F%96%E5%BE%97
        data = JSONParser().parse(request)
        user = request.user

        manipulator = TaskManipulator()
        task = manipulator.create(
            title=data['title'],
            user=user,
            description=data['description']
        )

        serializer = TaskSerializer(task)
        return Response(serializer.data, status=status.HTTP_201_CREATED)
    
    def patch(self, request, task_id):
        """
        メッセージボディとして以下を期待
        {
            "title": Taskの名前,
            "description": Taskの内容,
            "status": TaskStatusの名前
        }
        レスポンスメッセージとして以下を作成する。
        {
            "id": TaskのID,
            "title": 更新済みTaskの名前,
            "description": 更新済みTaskの詳細,
            "status": 更新済みTaskStatus
        }
        """



        data = JSONParser().parse(request)
        user = request.user

        manipulator = TaskManipulator()
        
        task, task_status = manipulator.update_task(
            user = user,
            task_id=task_id,
            title=data.get("title"),
            description=data.get("description"),
            status=data.get("status")
        )


        if task_status is None:
            if task is None:
                return Response(data=None, status=status.HTTP_404_NOT_FOUND)
            return Response(data=None, status=status.HTTP_400_BAD_REQUEST)
        
        if task is None:
            return Response(data=None, status=status.HTTP_404_NOT_FOUND)
        
        response_task = {
            "id": task.id,
            "title": task.title,
            "description": task.description,
            "status": task_status.name,
        }
        serializer = TaskSerializer(response_task)

        return Response(data=serializer.data, status=status.HTTP_200_OK)

    def delete(self, request, task_id):
        """
        削除対象としたいtask_id以外は、特に不要
        レスポンスメッセージとして以下を期待
        {
            "id": 削除したTaskのID,
            "title": 削除したTaskのtitle,
            "description": 削除したTaskのdescription
        }
        削除に成功した場合は200 OKを返す。
        削除対象のTaskが存在しない場合は404 Not Foundを返す。
        404 Not Foundを返す際は、レスポンスボディは空
        (本当にTaskが存在しない場合、
        Taskに紐づくUserがリクエストしたUserと異なる場合も同様)
        """

        manipulator = TaskManipulator()

        user = request.user

        task, _ = manipulator.get_task(user=user, task_id=task_id)

        if task is None:
            return Response(status=status.HTTP_404_NOT_FOUND)

        response_task = {
            "id": task.id,
            "title": task.title,
            "description": task.description
        }

        task_id = manipulator.delete_task(user=user, task_id=task.id)

        if task_id is None:
            # 基本的にはこれは到達不可能コード
            # (事前に削除対象のタスクが存在し、
            # リクエストしたユーザのものであることを確認済みであるため)
            return Response(status=status.HTTP_404_NOT_FOUND)

        serializer = TaskSerializer(response_task)

        return Response(data=serializer.data, status=status.HTTP_200_OK)



