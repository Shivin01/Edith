from django.db import models

from edith_models.models.base import Base
from edith_models.models.client import Client


class HolidayList(Base):
    date = models.DateField()
    name = models.CharField(max_length=Base.MAX_LENGTH_SMALL)
    description = models.TextField(null=True)
    year = models.IntegerField()
    client = models.ForeignKey(Client, on_delete=models.CASCADE, related_name="holiday_list")
