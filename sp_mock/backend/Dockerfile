## backend
FROM node:lts-alpine
WORKDIR /usr/app
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 80
# RUN npm run build
CMD ["node", "./server.js"]
