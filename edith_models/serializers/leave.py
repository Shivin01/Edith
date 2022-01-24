from builtins import object
from rest_framework import serializers

from edith_models.models import Leave
from .base import BaseSerializer, TimestampField


class LeaveSerializer(BaseSerializer):
    start_date_time = TimestampField()
    stop_date_time = TimestampField()

    class Meta(object):
        model = Leave
        fields = '__all__'
        extra_kwargs = {
            'employee': {'read_only': True}
        }

    def create(self, validated_data):
        """
        Overwriting create method of serializer to populate
        employee field of leave table
        """
        validated_data['employee'] = self.context['request'].user
        return super(LeaveSerializer, self).create(validated_data)
