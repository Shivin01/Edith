from builtins import object
from edith_models.models import Skill
from .base import BaseSerializer


class SkillSerializer(BaseSerializer):
    class Meta(object):
        model = Skill
        fields = '__all__'
