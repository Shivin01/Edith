from builtins import object
from rest_framework import serializers

from edith_models.models import Leave
from .base import BaseSerializer, TimestampField


class LeaveSerializer(BaseSerializer):
    start_date_time = TimestampField()
    stop_date_time = TimestampField()
    employee = serializers.SerializerMethodField()
    approved_by = serializers.SerializerMethodField()

    def get_employee(self, obj):
        return {
            'first_name': obj.employee.first_name,
            'middle_name': obj.employee.middle_name,
            'last_name': obj.employee.last_name,
        }

    def get_approved_by(self, obj):
        return {
            'first_name': obj.approved_by.first_name if obj.approved_by else "",
            'middle_name': obj.approved_by.middle_name if obj.approved_by else "",
            'last_name': obj.approved_by.last_name if obj.approved_by else "",
        }

    class Meta(object):
        model = Leave
        fields = '__all__'
        extra_kwargs = {
            'employee': {'read_only': True},
            'approved_by': {'read_only': True}
        }

    def create(self, validated_data):
        """
        Overwriting create method of serializer to populate
        employee field of leave table
        """
        validated_data['employee'] = self.context['request'].user
        return super(LeaveSerializer, self).create(validated_data)

    def update(self, instance, validated_data):
        validated_data['approved_by'] = self.context['request'].user
        return super(LeaveSerializer, self).update(instance, validated_data)
