
from builtins import object

from edith_models.models import Attendance
from .base import BaseSerializer


class AttendanceSerializer(BaseSerializer):

    class Meta(object):
        model = Attendance
        fields = '__all__'
        extra_kwargs = {
            'employee': {'read_only': True}
        }

    def create(self, validated_data):
        """
        Overwriting create method of serializer to populate
        employee field of leave table
        """
        print(validated_data)
        validated_data['employee'] = self.context['request'].user
        return super(AttendanceSerializer, self).create(validated_data)

