FROM node:12.18.3-alpine as buildNode

WORKDIR /app

COPY package-lock.json .
COPY package.json .

RUN yarn install
ADD . .

RUN yarn build --configuration=production

FROM httpd:2.4-alpine
COPY --from=buildNode /app/dist/client /usr/local/apache2/htdocs/
