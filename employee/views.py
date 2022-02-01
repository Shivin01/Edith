import http
from rest_framework import viewsets
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.exceptions import ValidationError
from django.core.exceptions import ValidationError as DjangoValidationError
from rest_framework import status
from rest_framework import filters
from django_filters.rest_framework import DjangoFilterBackend
from rest_framework import mixins

from edith_models.models import Employee, Resignation, Leave, Bonus, Attendance, HolidayList
from edith_models.serializers import (
    EmployeeSerializer,
    ResignationSerializer,
    LeaveSerializer,
    BonusSerializer,
    EmployeeMinimalSerializer,
    AttendanceSerializer,
)


class EmployeeViewSet(viewsets.ReadOnlyModelViewSet):
    """
    This viewset automatically provides `list` and `retrieve` actions.
    """
    queryset = Employee.objects.all()
    serializer_class = EmployeeSerializer
    permission_classes = (IsAuthenticated, )
    filterset_fields = ['created_at', 'slack_id']
    search_fields = ['employee__name']


class IsAdminUserAuthenticated(IsAuthenticated):
    """
    Allows access only to authenticated users.
    """

    def has_permission(self, request, view):
        return bool(request.user and request.user.is_authenticated and request.user.designation.upper() in ['ADMIN', 'HR', 'MANAGER'])


class EmployeeAdminViewSet(viewsets.ModelViewSet):
    queryset = Employee.objects.all()
    serializer_class = EmployeeSerializer
    permission_classes = (IsAuthenticated, )
    http_method_names = ['delete', 'patch', 'put']


class EmployeeMinimalViewSet(viewsets.ReadOnlyModelViewSet):
    """
    This viewset automatically provides `list` and `retrieve` actions.
    """
    queryset = Employee.objects.all()
    serializer_class = EmployeeMinimalSerializer
    filterset_fields = ['slack_id']
    search_fields = ['employee__name']


class ResignationViewSet(viewsets.ModelViewSet):
    serializer_class = ResignationSerializer
    filter_backends = [filters.SearchFilter, DjangoFilterBackend]
    filterset_fields = ['created_at', 'employee']
    search_fields = ['employee__name']

    def get_queryset(self):
        return Resignation.objects.filter(employee=self.request.user)


class LeaveViewSet(viewsets.ModelViewSet):
    """
    This viewset automatically provides `list` and `retrieve` actions.
    """
    queryset = Leave.objects.all()
    serializer_class = LeaveSerializer
    http_method_names = ['get', 'post']

    def get_queryset(self):
        employees = Employee.objects.filter(client=self.request.user.client)
        return Leave.objects.filter(employee__in=[employee.id for employee in employees]).exclude(approved_by=None)


class BonusViewSet(viewsets.ModelViewSet):
    serializer_class = BonusSerializer
    permission_classes = (IsAuthenticated, )
    filter_backends = [filters.SearchFilter, DjangoFilterBackend]
    filterset_fields = ['name', 'created_at', 'employee']
    search_fields = ['employee__name', 'name']

    def get_queryset(self):
        return Bonus.objects.filter(employee=self.request.user)

    def create(self, request, *args, **kwargs):
        try:
            if 'approval_from' in request.data:
                request.data.pop('approval_from')
            if 'manager' in request.data:
                request.data.pop('manager')
            serializer = self.get_serializer(data=request.data)
            serializer.is_valid(raise_exception=True)
            self.perform_create(serializer)
            return Response(serializer.data, status=status.HTTP_201_CREATED)
        except ValidationError:
            return Response(status=status.HTTP_400_BAD_REQUEST)

    def update(self, request, *args, **kwargs):
        try:
            if 'approval_from' in request.data:
                request.data.pop('approval_from')
            if 'manager' in request.data:
                request.data.pop('manager')
            instance = self.get_object()
            serializer = self.get_serializer(instance,
                                             data=request.data,
                                             partial=True)
            serializer.is_valid(raise_exception=True)
            self.perform_update(serializer)
            return Response(serializer.data, status=status.HTTP_200_OK)
        except ValidationError:
            return Response(status=status.HTTP_400_BAD_REQUEST)


class ApprovalBonusViewSet(viewsets.ModelViewSet):
    serializer_class = BonusSerializer
    permission_classes = (IsAuthenticated, )
    http_method_names = ["get", "patch"]
    filter_backends = [filters.SearchFilter, DjangoFilterBackend]
    filterset_fields = ['name', 'created_at', 'employee', 'approval_from', 'manager']
    search_fields = ['employee__name', 'name', 'approval_from__name', 'manager__name']

    def get_queryset(self):
        employees = Employee.objects.filter(client=self.request.user.client)
        return Bonus.objects.filter(employee__in=[employee.id for employee in employees])

    def update(self, request, *args, **kwargs):
        try:
            if 'approval_from' in request.data:
                request.data.pop('approval_from')
            if 'manager' in request.data:
                request.data.pop('manager')
            request.data['approval_from'] = self.request.user
            instance = self.get_object()
            serializer = self.get_serializer(instance,
                                             data=request.data,
                                             partial=True)
            serializer.is_valid(raise_exception=True)
            self.perform_update(serializer)
        except ValidationError:
            return Response(status=status.HTTP_400_BAD_REQUEST)
        except DjangoValidationError as e:
            return Response({"Message": ";".join(e.messages)}, status=status.HTTP_400_BAD_REQUEST)


class EmployeeAttendanceViewSet(viewsets.ModelViewSet):
    serializer_class = AttendanceSerializer
    permission_classes = (IsAuthenticated, )
    http_method_names = ["get", "post"]
    filterset_fields = ['employee']

    def get_queryset(self):
        return Attendance.objects.filter(employee=self.request.user)

    def create(self, request, *args, **kwargs):
        print(request.data)
        return super(EmployeeAttendanceViewSet, self).create(request, *args, **kwargs)


class AdminAttendanceViewSet(viewsets.ModelViewSet):
    serializer_class = AttendanceSerializer
    permission_classes = (IsAuthenticated, )
    http_method_names = ["get", "post"]
    filterset_fields = ['employee']

    def get_queryset(self):
        employees = Employee.objects.filter(client=self.request.user.client)
        return Attendance.objects.filter(employee__in=[employee.id for employee in employees])


class LeaveApprovalViewSet(viewsets.ModelViewSet):
    serializer_class = LeaveSerializer
    permission_classes = (IsAdminUserAuthenticated, )
    http_method_names = ["get", "patch"]

    def get_queryset(self):
        employees = Employee.objects.filter(client=self.request.user.client)
        return Leave.objects.filter(employee__in=[employee.id for employee in employees], approved_by=None)
