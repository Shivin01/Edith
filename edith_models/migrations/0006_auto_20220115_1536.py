# Generated by Django 3.2.11 on 2022-01-15 15:36

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('edith_models', '0005_alter_employee_skills'),
    ]

    operations = [
        migrations.AlterField(
            model_name='employee',
            name='gender',
            field=models.CharField(choices=[('MALE', 'MALE'), ('FEMALE', 'FEMALE'), ('OTHERS', 'OTHERS')], max_length=6, null=True),
        ),
        migrations.AlterField(
            model_name='employee',
            name='skills',
            field=models.ManyToManyField(null=True, related_name='employee', to='edith_models.Skill'),
        ),
    ]
