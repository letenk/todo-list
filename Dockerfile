# Build a tiny docker image
# Use image alpine
FROM alpine:latest
# Create folder /app
RUN mkdir /app
# Copy binary file from image `builder` to folder `/app` tiny image
COPY todoListApp /app
# Run todoListApp
CMD ["/app/todoListApp"]