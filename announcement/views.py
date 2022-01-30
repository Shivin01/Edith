from rest_framework import viewsets
from rest_framework.permissions import IsAuthenticated
from rest_framework_jwt.authentication import JSONWebTokenAuthentication
from rest_framework.response import Response
from datetime import date
from datetime import timedelta

from edith_models.models import Employee
from edith_models.serializers import AnnouncementSerializer


class CelebrationViewSet(viewsets.ViewSet):
    """
    ViewSet for news feed
    """
    authentication_classes = (JSONWebTokenAuthentication, )
    permission_classes = (IsAuthenticated,)

    def get(self, request, *args, **kwargs):
        """
        :return:
        """
        today = date.today()
        employees = []
        # Employee.objects.filter(client=request.user.client)
        for employee in Employee.objects.all():
            if employee.birth_date == today:
                employees.append({
                    'image': employee.image if employee.image else '',
                    'first_name': employee.first_name,
                    'middle_name': employee.middle_name,
                    'last_name': employee.last_name,
                    'id': employee.id,
                })

        return Response(employees)


class NewsFeedViewSet(viewsets.ViewSet):
    """
    ViewSet for news feed
    """
    authentication_classes = (JSONWebTokenAuthentication, )
    permission_classes = (IsAuthenticated,)

    DAYS_DELTA = 30

    def get(self, request, *args, **kwargs):
        """
        :return:
        """
        today = date.today()
        employees = []
        # Employee.objects.filter(client=request.user.client)
        for employee in Employee.objects.all():
            new_date = employee.joining_date + timedelta(days=self.DAYS_DELTA)
            if new_date > today:
                employees.append({
                    'image': employee.image if employee.image else '',
                    'first_name': employee.first_name,
                    'middle_name': employee.middle_name,
                    'last_name': employee.last_name,
                    'skills': [skill.name for skill in employee.skills.all()],
                    'designation': employee.designation,
                    'id': employee.id,
                    'joining_date': employee.joining_date
                })

        return Response(employees)


class AnnouncementViewSet(viewsets.ModelViewSet):
    """
    Announcement viewset
    """
    authentication_classes = (JSONWebTokenAuthentication, )
    permission_classes = (IsAuthenticated, )
    serializer_class = AnnouncementSerializer
    http_method_names = ["get", "post"]
