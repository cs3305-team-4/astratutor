FROM node:13.12.0-alpine AS dev

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY package.json ./
COPY package-lock.json ./

RUN npm install

COPY . ./

CMD ["npm", "run", "start"]

FROM node:13.12.0-alpine AS prod

# TODO prod build