from django.urls import path
from .views import compare_faces

urlpatterns = [
    path('compare-faces/', compare_faces, name='compare-faces'),
]