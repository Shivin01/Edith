
from builtins import object

from edith_models.models import Resignation
from .base import BaseSerializer, TimestampField


class ResignationSerializer(BaseSerializer):
    start_date_time = TimestampField()
    stop_date_time = TimestampField()

    class Meta(object):
        model = Resignation
        fields = '__all__'
        extra_kwargs = {
            'employee': {'read_only': True}
        }

    def create(self, validated_data):
        """
        Overwriting create method of serializer to populate
        employee field of resignation table
        """
        validated_data['employee'] = self.context['request'].user
        return super(ResignationSerializer, self).create(validated_data)

