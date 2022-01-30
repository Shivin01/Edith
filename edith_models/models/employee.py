from django.core.files.storage import FileSystemStorage
from django.contrib.auth.models import (
    AbstractBaseUser,
    BaseUserManager
)
from django.db import models
from django.utils.translation import gettext_lazy as _
from django.core.validators import RegexValidator

from .base import Base
from .client import Client


class UserManager(BaseUserManager):
    """
    Manager for User Model
    """

    def create_user(self, username, email, password=None, is_active=True, is_staff=False,
                    is_admin=False, is_superuser=False, client=None):
        if not username:
            raise ValueError('User must have username')
        if not password:
            raise ValueError('User must have password')

        user_obj = self.model(
            username=username
        )
        user_obj.email = self.normalize_email(email)
        user_obj.set_password(password)
        user_obj.is_staff = is_staff
        user_obj.is_admin = is_admin
        user_obj.is_active = is_active
        user_obj.is_superuser = is_superuser
        user_obj.client = client
        user_obj.save(using=self._db)
        return user_obj

    def create_staffuser(self, username, email, password=None):
        return self.create_user(
            username,
            email,
            password=password,
            is_staff=True
        )

    def create_superuser(self, username, email, password=None, **kwargs):
        return self.create_user(
            username,
            email,
            password=password,
            is_staff=True,
            is_admin=True,
            is_superuser=True,
            **kwargs
        )


class Employee(AbstractBaseUser, Base):
    """
    Enterprise User Model
    :field: email: Email id of the user
    :field: is_active: Tells us whether this user account should be
            considered active
    :field: is_staff: Tells us whether this user can access the admin site
    :field: is_admin:  Tells us that this user has all permissions without
            explicitly assigning them
    :field: is_superuser: Tells us that this user has all permissions without
            explicitly assigning them (only user's whose email ends with
            quartic.ai can be superusers)
    :field: first_name: First name of the user
    :field: last_name: Last name of the user
    :field: phone_number: Phone number of the user
    """

    GENDER_MALE = 'MALE'
    GENDER_FEMALE = 'FEMALE'
    GENDER_OTHERS = 'OTHERS'

    GENDER_CHOICES = [
        (GENDER_MALE, 'MALE'),
        (GENDER_FEMALE, 'FEMALE'),
        (GENDER_OTHERS, 'OTHERS'),
    ]

    phone_regex = RegexValidator(
        regex=r'^\+?1?\d{9,15}$',
        message="Phone number must be entered in International format: "
                "'+14169058972'. Up to 15 digits allowed.")

    is_active = models.BooleanField(default=True)
    is_staff = models.BooleanField(default=False)
    is_admin = models.BooleanField(default=False)
    is_superuser = models.BooleanField(default=False)

    username = models.CharField(max_length=Base.MAX_LENGTH_MEDIUM, unique=True)
    designation = models.CharField(
        max_length=Base.MAX_LENGTH_SMALL,
        null=True,
        blank=True
    )
    joining_date = models.DateField(null=True)
    first_name = models.CharField(max_length=Base.MAX_LENGTH_MEDIUM)
    middle_name = models.CharField(max_length=Base.MAX_LENGTH_MEDIUM, null=True, blank=True)
    last_name = models.CharField(max_length=Base.MAX_LENGTH_MEDIUM, null=True, blank=True)
    birth_date = models.DateField(null=True)
    phone_number = models.CharField(validators=[phone_regex], max_length=17)
    email = models.EmailField(max_length=Base.MAX_LENGTH_MEDIUM)
    image = models.ImageField(upload_to='profile_images/', storage=FileSystemStorage(),
                              null=True, blank=True)

    emergency_phone_number = models.CharField(validators=[phone_regex], max_length=17,
                                              null=True)
    emergency_email = models.EmailField(max_length=Base.MAX_LENGTH_MEDIUM, null=True)
    address = models.TextField(null=True, blank=True)
    gender = models.CharField(
        max_length=6,
        choices=GENDER_CHOICES,
        null=True
    )
    skills = models.ManyToManyField('Skill', related_name='employee', blank=True)
    client = models.ForeignKey(Client, related_name="employees", on_delete=models.CASCADE, blank=True, null=True)
    resigned_date = models.DateField(null=True, blank=True)
    slack_id = models.CharField(max_length=Base.MAX_LENGTH_MEDIUM, unique=True, null=True, blank=True)
    soft_delete = models.BooleanField(default=False)

    USERNAME_FIELD = 'username'
    REQUIRED_FIELD = []

    objects = UserManager()

    class Meta(Base.Meta):
        verbose_name = _('Employee')
        verbose_name_plural = _('Employees')

    def __str__(self):
        return self.username

    def __unicode__(self):
        return str(self.username)

    def get_full_name(self):
        """
        Gets full name of the user.
        """
        name = str(self.first_name) if self.first_name else ''
        if self.middle_name:
            name += " " + str(self.middle_name)
        if self.last_name:
            name += " " + str(self.last_name)
        return name

    def get_image(self):
        """
        will return the image url.
        """
        if self.image:
            return self.image.url

    # def save(self, *args, **kwargs):
    #     if not self.client:
    #         raise FieldDoesNotExist('Client does not exist for this employee')
    #     super(Employee, self).save(*args, **kwargs)
