from django.db import models

from edith_models.models.base import Base
from edith_models.models.employee import Employee
from edith_models.models.client import Client


class Announcement(Base):
    client = models.ForeignKey(
        Client,
        on_delete=models.CASCADE,
        related_name='announcements')
    user = models.ForeignKey(
        Employee,
        on_delete=models.CASCADE,
        related_name='announcements')
    detail = models.TextField(null=True)
    type = models.CharField(max_length=Base.MAX_LENGTH_SMALL)
