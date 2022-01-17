from django.db import models


class Base(models.Model):

    MAX_LENGTH_SMALL = 127
    MAX_LENGTH_MEDIUM = 255
    MAX_LENGTH_LARGE = 511

    APP_LABEL = "edith_models"

    created_at = models.DateTimeField(auto_now_add=True, null=True)
    updated_at = models.DateTimeField(auto_now=True, null=True)

    class Meta:
        abstract = True
        app_label = "edith_models"

    def save(self, *args, **kwargs):
        super().save(*args, **kwargs)
