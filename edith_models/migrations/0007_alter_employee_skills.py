# Generated by Django 3.2.11 on 2022-01-15 15:37

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('edith_models', '0006_auto_20220115_1536'),
    ]

    operations = [
        migrations.AlterField(
            model_name='employee',
            name='skills',
            field=models.ManyToManyField(related_name='employee', to='edith_models.Skill'),
        ),
    ]
