from edith_models.models import Department
from .base import BaseSerializer


class DepartmentSerializer(BaseSerializer):

    class Meta:
        model = Department
        fields = "__all__"
        extra_params = {
            'client': {'read_only': True}
        }

    def create(self, validated_data):
        """
        Overwriting create method of serializer to populate
        client field of department table
        """
        validated_data['client'] = self.context['request'].user.client
        return super(DepartmentSerializer, self).create(validated_data)
