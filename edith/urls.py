from django.urls import path, include, re_path
from rest_framework_jwt.views import obtain_jwt_token, refresh_jwt_token
from rest_framework_swagger.views import get_swagger_view


urlpatterns = [
    path('rest-auth/', include('rest_auth.urls')),
    path('rest-auth/registration/', include('rest_auth.registration.urls')),
    re_path(r'api/token/$', obtain_jwt_token, name='token_obtain_pair'),
    re_path(r'api/token/refresh/$', refresh_jwt_token, name='token_refresh'),
    path('employees/', include('employee.urls')),
    path('announcements/', include('announcement.urls')),
    path('client/', include('client.urls')),
    re_path(r'^swagger/$', get_swagger_view(title='Edith API')),
]
