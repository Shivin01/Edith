from rest_framework import viewsets
from rest_framework.permissions import IsAuthenticated
from datetime import date

from edith_models.serializers import ClientSerializer, DepartmentSerializer, HolidayListSerializer
from edith_models.models import Client, Department, HolidayList


class ClientViewSet(viewsets.ModelViewSet):
    """
    Organization view set
    """
    queryset = Client.objects.all()
    serializer_class = ClientSerializer
    permission_classes = (IsAuthenticated,)


class DepartmentViewSet(viewsets.ModelViewSet):
    permission_classes = (IsAuthenticated, )
    serializer_class = DepartmentSerializer

    def get_queryset(self):
        return Department.objects.filter(client=self.request.user.client)


class HolidayListViewSet(viewsets.ReadOnlyModelViewSet):
    serializer_class = HolidayListSerializer
    permission_classes = (IsAuthenticated, )

    def get_queryset(self):
        return HolidayList.objects.filter(client=self.request.user.client, year=date.today().year)
