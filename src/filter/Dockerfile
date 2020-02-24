FROM python:3.7
ADD ./requirements.txt /tmp/requirements.txt
RUN pip3 install --upgrade pip setuptools && pip3 install --no-cache-dir -r /tmp/requirements.txt
ADD . /src/
WORKDIR /src
CMD [ "python3" , "filter.py" ]
