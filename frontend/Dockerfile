# Use a lightweight Nginx image to serve static files
FROM nginx:1.25-alpine

# Switch to the non-root user provided by the nginx base image
USER nginx

# Copy the static frontend files into the Nginx web root
COPY . /usr/share/nginx/html/

# Copy our custom Nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80