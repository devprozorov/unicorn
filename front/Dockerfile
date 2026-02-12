FROM node:20-alpine

WORKDIR /app

COPY unicorn/package*.json ./
RUN npm install

COPY unicorn/ .

RUN npm run build

ENV NITRO_PORT=3000
ENV NITRO_HOST=0.0.0.0

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
