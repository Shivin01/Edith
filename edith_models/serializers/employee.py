from django.contrib.auth import update_session_auth_hash
from builtins import object
from rest_framework import serializers
from rest_auth.registration.serializers import RegisterSerializer

from edith_models.models import Employee, Base
from .base import BaseSerializer, TimestampField


class EmployeeRegisterSerializer(BaseSerializer, RegisterSerializer):
    # password = serializers.CharField(read_only=True)

    # def save(self, request):
    #     self.instance = request.data
    #     self.is_valid(raise_exception=True)
    #     self.validated_data.pop('password1')
    #     self.validated_data.pop('password2')
    #     user = self.create(self.validated_data)
    #     user.save()
    #     return user

    class Meta:
        model = Employee
        fields = "__all__"


class EmployeeSerializer(BaseSerializer):
    password = serializers.CharField(write_only=True)

    class Meta(object):
        model = Employee
        extra_kwargs = {'is_staff': {'read_only': True},
                        'is_admin': {'read_only': True},
                        'is_superuser': {'read_only': True},
                        'is_active': {'read_only': True}}
        fields = '__all__'

    def create(self, validated_data):
        """
        :param validated_data:
        :return: EnterpriseUser instance
        """
        user = super(EmployeeSerializer, self).create(validated_data)
        user.set_password(validated_data['password'])
        user.save()
        return user

    def update(self, instance, validated_data):
        """
        :param instance: Instance of the object
        :param validated_data: Validated data.
        :return: EnterpriseUser instance
        """
        instance = super(
            EmployeeSerializer,
            self).update(
            instance,
            validated_data)
        if 'password' in validated_data:
            instance.set_password(validated_data['password'])
            instance.save()
        if self.context:
            if 'request' in self.context:
                if self.context['request'].user == instance:
                    update_session_auth_hash(self.context['request'], instance)
        return instance

