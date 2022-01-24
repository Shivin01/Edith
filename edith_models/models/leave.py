from django.db import models

from edith_models.models.base import Base
from edith_models.models.employee import Employee


class Leave(Base):
    LEAVE_FULL_DAY = 'FULL_DAY'
    LEAVE_HALF_DAY = 'HALF_DAY'
    LEAVE_RESTRICTED = 'RESTRICTED'

    LEAVE_CHOICES = [
        (LEAVE_FULL_DAY, 'FULL_DAY'),
        (LEAVE_HALF_DAY, 'HALF_DAY'),
        (LEAVE_RESTRICTED, 'RESTRICTED'),
    ]

    TYPE_CASUAL = 'CASUAL'
    TYPE_SICK = 'SICK'
    TYPE_MATERNITY = 'MATERNITY'

    TYPE_CHOICES = [
        (TYPE_CASUAL, 'CASUAL'),
        (TYPE_SICK, 'SICK'),
        (TYPE_MATERNITY, 'MATERNITY'),
    ]

    kind = models.CharField(
        max_length=Base.MAX_LENGTH_SMALL,
        choices=LEAVE_CHOICES,
        default=LEAVE_FULL_DAY,
    )
    type = models.CharField(
        max_length=Base.MAX_LENGTH_SMALL,
        choices=TYPE_CHOICES,
    )
    start_date_time = models.DateTimeField()
    stop_date_time = models.DateTimeField()
    employee = models.ForeignKey(Employee, related_name="leaves", on_delete=models.CASCADE)


