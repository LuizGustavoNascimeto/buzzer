FROM node:24-alpine3.22

ENV PORT=3000

WORKDIR /frontend-react-js

# Copia APENAS os arquivos de dependências primeiro
# → npm install só roda quando package.json mudar
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

# Agora copia o restante do código
COPY . .

EXPOSE ${PORT}
CMD ["npm", "start"]