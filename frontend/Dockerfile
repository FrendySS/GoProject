FROM node:18

WORKDIR /app

# копируем только package*.json
COPY frontend/package.json frontend/package-lock.json* ./

# устанавливаем зависимости
RUN npm ci

# копируем весь исходники
COPY frontend/ ./

# билдим Next.js
RUN npm run build

# впереди SSR-сервер Next.js слушает в 3000
EXPOSE 3000

# запускаем продакшен-сервер
CMD ["npm", "start"]
