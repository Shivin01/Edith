from django.contrib.auth.backends import ModelBackend, UserModel
from django.core.exceptions import MultipleObjectsReturned
from django.db.models import Q

from edith_models.models import Employee


class EmailBackend(ModelBackend):

    def authenticate(self, request, email=None, username=None, password=None, **kwargs):
        try: #to allow authentication through phone number or any other field, modify the below statement
            user = Employee.objects.get(Q(username__iexact=username) | Q(email__iexact=email))
        except Employee.DoesNotExist:
            UserModel().set_password(password)
        except MultipleObjectsReturned:
            return Employee.objects.filter(email=username).order_by('id').first()
        else:
            if user.check_password(password) and self.user_can_authenticate(user):
                return user

    def get_user(self, user_id):
        try:
            user = Employee.objects.get(pk=user_id)
        except Employee.DoesNotExist:
            return None

        return user if self.user_can_authenticate(user) else None