FROM node:12.18.3-alpine as buildNode

WORKDIR /app

COPY package-lock.json .
COPY package.json .

RUN npm install
COPY src/environments .
ADD . .

RUN npm run build

FROM httpd:2.4-alpine
COPY --from=buildNode /app/dist/client /usr/local/apache2/htdocs/
