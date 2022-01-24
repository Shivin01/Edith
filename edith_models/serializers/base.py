from datetime import datetime
from rest_framework import serializers


class TimestampField(serializers.DateTimeField):
    """
    Convert a django datetime to/from timestamp.
    """
    def to_representation(self, value):
        """
        Convert the field to its internal representation (aka timestamp)
        :param value: the DateTime value
        :return: a UTC timestamp integer
        """
        return value.timestamp()

    def to_internal_value(self, value):
        """
        deserialize a timestamp to a DateTime value
        :param value: the timestamp value
        :return: a django DateTime value
        """
        return datetime.fromtimestamp(int(value))


class ReadOnlyTimestampField(serializers.ReadOnlyField):
    def to_representation(self, value):
        return int(value.timestamp() * 1000.)


class BaseSerializer(serializers.ModelSerializer):
    created_at = ReadOnlyTimestampField()
    updated_at = ReadOnlyTimestampField()
