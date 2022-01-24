from edith_models.models import Bonus
from .base import BaseSerializer


class BonusSerializer(BaseSerializer):

    class Meta:
        model = Bonus
        fields = "__all__"
        extra_params = {
            'employee': {'read_only': True}
        }

    def create(self, validated_data):
        """
        Overwriting create method of serializer to populate
        client field of department table
        """
        validated_data['employee'] = self.context['request'].user
        return super(BonusSerializer, self).create(validated_data)
