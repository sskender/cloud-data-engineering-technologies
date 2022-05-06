FROM mcr.microsoft.com/azure-functions/dotnet:3.0-appservice 
ENV AzureWebJobsScriptRoot=/home/site/wwwroot \
    AzureFunctionsJobHost__Logging__Console__IsEnabled=true

COPY --from=installer-env ["/home/site/wwwroot", "/home/site/wwwroot"]


FROM golang:1.16-alpine

WORKDIR /app/helloworld

COPY . .
RUN go build helloworld.go

EXPOSE 8080

CMD ["./helloworld"]