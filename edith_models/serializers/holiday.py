from builtins import object
from edith_models.models import HolidayList
from .base import BaseSerializer


class HolidayListSerializer(BaseSerializer):
    class Meta(object):
        model = HolidayList
        fields = '__all__'
