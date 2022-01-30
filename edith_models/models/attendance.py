from django.db import models

from edith_models.models.base import Base
from edith_models.models.employee import Employee


class Attendance(Base):
    date = models.DateField()
    employee = models.ForeignKey(Employee, related_name="attendances", on_delete=models.CASCADE)
