from django.test import TestCase
from rest.models.task_status import TaskStatus

class TestTaskHistory(TestCase):

    fixtures = [
        'task_status',
    ]

    def test_number_of_task_statuses(self):
        """
        server/sampleapp/rest/fixtures/task_status.yaml に記述されている
        4つのレコードが格納されていることを確認する。
        他のテストが失敗する際の切り分けに利用することを期待している。
        """
        statuses = TaskStatus.objects.all()
        self.assertEqual(len(statuses), 4)

    def test_task_status_intended(self):
        """
        server/sampleapp/rest/fixtures/task_status.yaml に記述されている
        4つのレコードが期待された値を持っていることを確認する。
        ┌-----┬------------┐
        │id   │name        │
        ├-----┼------------┤
        │1    │TODO        │
        │2    │RUNNING     │
        │3    │FINISHED    │
        │4    │PENDING     │
        └-----┴------------┘
        """
        todo = TaskStatus.objects.get(pk=1)
        self.assertEqual(todo.name, 'TODO')

        running = TaskStatus.objects.get(pk=2)
        self.assertEqual(running.name, 'RUNNING')

        finished = TaskStatus.objects.get(pk=3)
        self.assertEqual(finished.name, 'FINISHED')

        pending = TaskStatus.objects.get(pk=4)
        self.assertEqual(pending.name, 'PENDING')