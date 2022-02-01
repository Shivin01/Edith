from django.db import models

from edith_models.models.base import Base


class Client(Base):
    name = models.CharField(unique=True, max_length=Base.MAX_LENGTH_MEDIUM)
    registered_name = models.CharField(unique=True, max_length=Base.MAX_LENGTH_MEDIUM)
    address = models.TextField(null=True, blank=True)
    leave_count = models.IntegerField(default=0)
    notice_period_count = models.IntegerField(default-0)
