FROM node:13.12.0-alpine AS dev

WORKDIR /app

COPY . ./

RUN npm install

RUN npm install npm-install-if-needed

CMD npx npm-install-if-needed && chmod 777 ./node_modules && npm run start

FROM node:13.12.0-alpine AS prod

# TODO prod build
