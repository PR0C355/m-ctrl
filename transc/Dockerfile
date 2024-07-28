FROM python:3.10.6-slim

WORKDIR /app

COPY requirements.txt .

RUN pip install -r requirements.txt

COPY transc.py .

RUN apt-get update
RUN apt-get install ffmpeg -y

CMD ["python", "-m", "transc"]