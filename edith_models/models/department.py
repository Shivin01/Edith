
from django.db import models

from edith_models.models.base import Base
from edith_models.models.client import Client
from edith_models.models.employee import Employee


class Department(Base):
    name = models.CharField(unique=True, max_length=Base.MAX_LENGTH_MEDIUM)
    client = models.ForeignKey(Client, on_delete=models.CASCADE, related_name="departments")
    employee = models.ManyToManyField(Employee, related_name="departments")

