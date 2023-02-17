FROM boxfuse/flyway:latest

RUN ["rm", "-fr", "/flyway/conf"]
COPY flyway.conf /flyway/conf/

RUN ["rm", "-fr", "/flyway/sql"]
COPY *.sql /flyway/sql/

ENTRYPOINT ["flyway", "migrate", "-url=jdbc:postgresql://127.0.0.1:5432/postgres?ssl=false", "-user=postgres", "-password=password"]