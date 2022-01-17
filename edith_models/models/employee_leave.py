
from django.db import models

from edith_models.models.base import Base
from edith_models.models.employee import Employee
from django.contrib.postgres.fields import ArrayField


# TODO: not required
class EmployeeLeave(Base):
    FIRST_HALF = 0
    SECOND_HALF = 1

    REQUIRED_SCHEMA = {
        "type": "object",
        "properties": {
            "half": {
              "type" : "number",
              "default": FIRST_HALF,
            },
            "date": {
                "type": "number",
                "default": ""
            },
        },
        "additionalProperties": False
    }

    casual = ArrayField(models.DateField())
    half = ArrayField(models.JSONField(default=dict))
    sick = ArrayField(models.DateField())
    employee = models.OneToOneField(Employee, related_name="employee_leave", on_delete=models.CASCADE)



