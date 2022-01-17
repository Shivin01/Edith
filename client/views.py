from rest_framework import generics, viewsets
from rest_framework.permissions import IsAuthenticated

from edith_models.serializers import ClientSerializer, EmployeeSerializer
from edith_models.models import Client, Employee


class ClientViewSet(viewsets.ModelViewSet):
    """
    Organization view set
    """
    queryset = Client.objects.all()
    serializer_class = ClientSerializer
    permission_classes = (IsAuthenticated,)


class EmployeeListView(generics.ListAPIView):
    """
    User List view.
    """
    permission_classes = (IsAuthenticated,)
    queryset = Employee.objects.all()
    serializer_class = EmployeeSerializer


class UserProfileView(viewsets.ModelViewSet):
    """
    User Profile viewsets.
    """
    serializer_class = serializers.UserProfileSerializer
    permission_classes = (IsAuthenticated, )

    def get_queryset(self):
        """
        Get the queryset.
        :return: user profile instance.
        """
        user = self.request.user
        user_profile, _ = UserProfile.objects.get_or_create(user=user)
        user_profile = UserProfile.objects.filter(user=user)
        return user_profile

    def update(self, request, *args, **kwargs):
        """
        Over writing create method.
        :param request:
        :return:
        """
        try:
            user_profile = UserProfile.objects.get(user=request.user)
            serializer = self.get_serializer(user_profile, data=request.data, partial=True)
            serializer.is_valid(raise_exception=True)
            self.perform_create(serializer)
            data = {
                'first_name': request.data.get('first_name', request.user.first_name),
                'last_name': request.data.get('last_name', request.user.last_name)
            }
            update_serializer = UserSerializer(
                request.user, data=data, partial=True)
            update_serializer.is_valid(raise_exception=True)
            update_serializer.save()
            return Response(serializer.data, status=status.HTTP_201_CREATED)
        except ValidationError:
            return Response(status=status.HTTP_400_BAD_REQUEST)
        except Exception:
            return Response(status=status.HTTP_500_INTERNAL_SERVER_ERROR)
