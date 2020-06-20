FROM python:3.8.3

COPY . /app
WORKDIR /app

RUN apt-get update && apt-get -y install cmake protobuf-compiler
RUN pip install -r requirements.txt

EXPOSE 8000
CMD python manage.py runserver 0.0.0.0:8000