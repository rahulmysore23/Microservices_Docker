from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_http_methods

import face_recognition
import urllib.request
import ast

@csrf_exempt
@require_http_methods(["POST"])
def compare_faces(request):
    response = {}

    json_data = request.body.decode('utf-8')
    data = ast.literal_eval(json_data)
    
    try:
        known_image = face_recognition.load_image_file(urllib.request.urlopen(data['known_image']))
        unknown_image = face_recognition.load_image_file(urllib.request.urlopen(data['unknown_image']))
    except:
        response['Message'] = 'Image/Images are missing'
        return JsonResponse(response, status=400)
    
    try:
        known_encoding = face_recognition.face_encodings(known_image)[0]
        unknown_encoding = face_recognition.face_encodings(unknown_image)[0]
    except IndexError:
        response['Message'] = 'No faces detected.'
        return JsonResponse(response, status=404)

    results = face_recognition.compare_faces([known_encoding], unknown_encoding)
    response['is_face_same'] = str(results[0])

    return JsonResponse(response, status=200)


