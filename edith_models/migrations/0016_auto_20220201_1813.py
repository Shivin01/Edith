# Generated by Django 3.2.11 on 2022-02-01 18:13

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('edith_models', '0015_alter_client_address'),
    ]

    operations = [
        migrations.AlterField(
            model_name='client',
            name='leave_count',
            field=models.IntegerField(default=0),
        ),
        migrations.AlterField(
            model_name='client',
            name='notice_period_count',
            field=models.IntegerField(default=0),
        ),
    ]
