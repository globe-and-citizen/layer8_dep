## backend
FROM node:lts-alpine
WORKDIR /usr/app
COPY package*.json ./
RUN npm install
COPY . .
# COPY ../../middleware/dist/middleware.wasm ./middleware.wasm
EXPOSE 80
# RUN npm run build
CMD ["node", "./server.js"]
