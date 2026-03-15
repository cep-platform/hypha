FROM python:3.11-slim
RUN pip install --no-cache-dir "ray"
CMD ["ray", "start", "--block", "--address=127.0.0.1:6379"]
