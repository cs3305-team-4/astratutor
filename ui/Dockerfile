FROM node:13.12.0-alpine AS dev

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY . ./

RUN npm install

CMD ["npm", "run", "start"]

FROM node:13.12.0-alpine AS prod

# TODO prod build