from builtins import object
from rest_framework import serializers

from edith_models.models import Announcement
from .base import BaseSerializer


class AnnouncementSerializer(BaseSerializer):
    user = serializers.SerializerMethodField()
    client = serializers.SerializerMethodField()

    def get_user(self, obj):
        return {
            'first_name': obj.user.first_name,
            'middle_name': obj.user.middle_name,
            'last_name': obj.user.last_name,
        }

    def get_client(self, obj):
        return obj.client.id

    class Meta(object):
        model = Announcement
        fields = '__all__'
        extra_params = {
            'user': {'read_only': True},
            'client': {'read_only': True}
        }

    def create(self, validated_data):
        """
        Overwriting create method of serializer to populate
        client field of department table
        """
        validated_data['user'] = self.context['request'].user
        validated_data['client'] = self.context['request'].user.client
        return super(AnnouncementSerializer, self).create(validated_data)
