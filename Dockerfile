# Start by building the application.
FROM golang:alpine as build

WORKDIR /app
COPY . /app

RUN cd /app && go build -o goapp

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian11

WORKDIR /app
COPY --from=build /app/goapp /app
COPY --from=build /app/.env /app
COPY --from=build /app/conf /app/conf

CMD [ "./goapp" , "-o" , "print"]