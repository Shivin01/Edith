from django.db import models

from edith_models.models.base import Base
from edith_models.models.employee import Employee


class Resignation(Base):
    start_date_time = models.DateTimeField()
    stop_date_time = models.DateTimeField()
    reason = models.TextField()
    employee = models.OneToOneField(Employee, on_delete=models.SET_NULL, related_name="resignation", null=True)


