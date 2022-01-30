from rest_framework import serializers
from rest_auth.registration.serializers import RegisterSerializer
from typing import List

from edith_models.models import Employee
from .base import BaseSerializer


class EmployeeRegisterSerializer(BaseSerializer, RegisterSerializer):
    password = serializers.CharField(read_only=True)

    def save(self, request):
        self.instance = request.data
        self.is_valid(raise_exception=True)
        self.validated_data.pop('password1')
        self.validated_data.pop('password2')
        user = self.create(self.validated_data)
        user.save()
        return user

    class Meta:
        model = Employee
        fields = "__all__"


class EmployeeMinimalSerializer(BaseSerializer):
    departments = serializers.SerializerMethodField(read_only=True)
    image = serializers.CharField()

    def get_departments(self, obj) -> List[str]:
        departments = []
        for dep in obj.departments.all():
            departments.append(dep.name)
        return departments

    def get_image(self, obj):
        return obj.image if obj.image else ""

    class Meta:
        model = Employee
        fields = (
            'slack_id',
            'username',
            'first_name',
            'middle_name',
            'last_name',
            'designation',
            'phone_number',
            'email',
            'image',
            'gender',
            'skills',
            'departments',
        )


class EmployeeSerializer(BaseSerializer):
    departments = serializers.SerializerMethodField(read_only=True)

    def get_departments(self, obj) -> List[str]:
        departments = []
        for dep in obj.departments.all():
            departments.append(dep.name)
        return departments    

    class Meta:
        model = Employee
        exclude = ('password', )
