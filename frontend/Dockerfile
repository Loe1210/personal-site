FROM nginx:1.27-alpine

COPY frontend/nginx/default.conf /etc/nginx/conf.d/default.conf
COPY static/ /usr/share/nginx/html/
COPY pet/ /usr/share/nginx/html/pet/
