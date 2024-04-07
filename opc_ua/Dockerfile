FROM python:3.9

WORKDIR /app

COPY ./opcua_simulation_server .

# RUN pip install asyncua -i https://mirror.sjtu.edu.cn/pypi/web/simple
RUN pip install asyncua

CMD python server.py

