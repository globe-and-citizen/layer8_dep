FROM node:21-alpine

WORKDIR /usr/src/resource/frontend

COPY . .

RUN npm install && npm run build

CMD ["npm", "run", "start"]
