FROM python:3.9

WORKDIR /app

COPY ./opcua_simulation_server .

RUN pip install asyncua -i https://mirror.sjtu.edu.cn/pypi/web/simple

CMD python server.py

