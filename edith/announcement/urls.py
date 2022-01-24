from django.urls import re_path, path, include
from rest_framework.routers import DefaultRouter

from announcement.views import NewsFeedViewSet, CelebrationViewSet, HolidayListViewSet


news_feed = NewsFeedViewSet.as_view({
    'get': 'get'
})

celebration = CelebrationViewSet.as_view({
    'get': 'get'
})

# Create a router and register our viewsets with it.
holiday_router = DefaultRouter()
holiday_router.register(r'', HolidayListViewSet)

# The API URLs are now determined automatically by the router.
urlpatterns = [
    re_path(r'^news_feed/$', news_feed, name='news_feed'),
    re_path(r'^celebration/$', celebration, name='celebration'),
    path('holiday_list/', include(holiday_router.urls))
]
