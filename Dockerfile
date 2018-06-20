FROM mattparlette/houseparty:latest

COPY bot.py /app

CMD ["python3", "/app/bot.py"]
