# Generated by Django 3.2.11 on 2022-01-30 09:30

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('edith_models', '0014_alter_leave_approved_by'),
    ]

    operations = [
        migrations.AlterField(
            model_name='client',
            name='address',
            field=models.TextField(blank=True, null=True),
        ),
    ]