FROM node:18

WORKDIR /app
# копируем package.json и (если есть) package-lock.json
COPY frontend/package.json frontend/package-lock.json* ./

# чистая установка
RUN npm ci

COPY frontend/ ./

RUN npm run build

EXPOSE 3000
CMD ["npm", "start"]
