FROM python:3.9-slim

WORKDIR /app

# Copy requirements and install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy model files and application code
COPY . .

# Expose port
EXPOSE 3005

# Start the application
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "3005"]