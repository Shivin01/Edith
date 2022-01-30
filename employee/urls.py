from django.urls import path, include, re_path
from rest_framework.routers import DefaultRouter
from employee import views

# Create a router and register our viewsets with it.
router = DefaultRouter()
router.register(r'info', views.EmployeeViewSet, basename='employee_info')
router.register(r'admin', views.EmployeeAdminViewSet, basename='employee_admin')
router.register(r'minimal_info', views.EmployeeMinimalViewSet, basename='employee_minimal_info')
router.register(r'resignation', views.ResignationViewSet, basename='employee_resignation')
router.register(r'leave', views.LeaveViewSet, basename='employee_leave')
router.register(r'admin_leave', views.LeaveApprovalViewSet, basename='admin_leave')
router.register(r'bonus', views.BonusViewSet, basename='employee_bonus')
router.register(r'bonus_approval', views.ApprovalBonusViewSet, basename='employee_bonus_approval')
router.register(r'employee_attendance', views.EmployeeAttendanceViewSet, basename='employee_attendance')
router.register(r'admin_attendance', views.AdminAttendanceViewSet, basename='admin_attendance')
router.register(r'holiday_list', views.AdminAttendanceViewSet, basename='admin_attendance')


# The API URLs are now determined automatically by the router.
urlpatterns = [
    path('', include(router.urls)),
]
