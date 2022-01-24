
from django.db import models

from edith_models.models.base import Base
from edith_models.models.employee import Employee


class Bonus(Base):
    name = models.CharField(max_length=Base.MAX_LENGTH_SMALL)
    amount = models.FloatField()
    attachment = models.FileField(upload_to="attachments/", null=True)
    employee = models.ForeignKey(Employee, related_name="bonuses", on_delete=models.CASCADE)
    approval_from = models.ForeignKey(Employee, related_name="approved_bonuses", on_delete=models.SET_NULL, null=True)
    manager = models.ForeignKey(Employee, related_name="employee_bonuses", on_delete=models.SET_NULL, null=True)
