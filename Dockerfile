FROM scratch

COPY service/migrations/ /migrations
COPY build/my-awesome-birthday-app /my-awesome-birthday-app
CMD ["/my-awesome-birthday-app"]
