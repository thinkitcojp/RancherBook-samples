FROM python:3.6-alpine

EXPOSE 8000
RUN adduser -D django
USER django
RUN mkdir -p /home/django/.local/bin
ENV PATH=$PATH:/home/django/.local/bin
COPY --chown=django:django requirements.txt ./
RUN pip install --user -r requirements.txt
COPY --chown=django:django ./sampleapp /home/django/sampleapp
WORKDIR /home/django/sampleapp
ENV DJANGO_SETTINGS_MODULE sampleapp.settings_mysql
ENV MYSQL_USER dummy
ENV MYSQL_PASSWORD dummy
ENV MYSQL_DATABASE sampleapp
ENV MYSQL_HOST dummy
ENV MYSQL_PORT 3306

CMD ["gunicorn", "-w" , "2", "-b", "0.0.0.0:8000", "sampleapp.wsgi"]
