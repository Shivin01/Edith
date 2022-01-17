from builtins import object
from edith_models.models import Announcement
from .base import BaseSerializer


class AnnouncementSerializer(BaseSerializer):
    class Meta(object):
        model = Announcement
        fields = '__all__'
