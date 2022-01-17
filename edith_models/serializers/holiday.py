from builtins import object
from edith_models.models import Holiday
from .base import BaseSerializer


class HolidaySerializer(BaseSerializer):
    class Meta(object):
        model = Holiday
        fields = '__all__'
