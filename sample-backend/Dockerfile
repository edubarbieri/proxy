FROM node:14-alpine

ENV PORT=3000

WORKDIR /app

COPY package*.json ./

COPY . .

RUN npm install

ENTRYPOINT [ "npm", "start" ]