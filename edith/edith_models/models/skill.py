from django.db import models

from edith_models.models.base import Base


class Skill(Base):
    name = models.CharField(unique=True, max_length=Base.MAX_LENGTH_MEDIUM)
    description = models.TextField(null=True, blank=True)
