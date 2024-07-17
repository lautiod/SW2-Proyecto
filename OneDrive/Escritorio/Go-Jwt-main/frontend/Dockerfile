FROM node:18-alpine

EXPOSE 5173

ADD . /frontend
WORKDIR /frontend

COPY package*.json ./
RUN npm i
COPY . .
CMD ["npm", "run", "dev"]