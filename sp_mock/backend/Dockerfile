FROM node:21-alpine3.17

WORKDIR /app

COPY package.json .
COPY package-lock.json .

RUN npm install

COPY . .

CMD ["node", "server.js"]
