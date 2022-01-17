from builtins import object
from edith_models.models import Client
from .base import BaseSerializer


class ClientSerializer(BaseSerializer):
    class Meta(object):
        model = Client
        fields = '__all__'
