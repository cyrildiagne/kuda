FROM nvidia/cuda:10.1-base

RUN apt-get update
RUN apt-get install -y --no-install-recommends \
  python3.7 python3-pip \
  && \
  apt-get clean && \
  apt-get autoremove && \
  rm -rf /var/lib/apt/lists/*
RUN pip3 install setuptools


WORKDIR /app

# Install production dependencies.
COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

# Copy local code to the container image.
COPY . ./

ENV PORT 80

# Run the web service using gunicorn.
CMD exec gunicorn --bind :$PORT --workers 1 --threads 8 app:app